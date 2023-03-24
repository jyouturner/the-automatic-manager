package notion

//db is to handle the write to nation database for example the task list
import (
	"context"
	"fmt"
	"time"

	"github.com/dstotijn/go-notion"
	log "github.com/sirupsen/logrus"
)

type TaskService struct {
	DatabaseId string
	client     *notion.Client
}

func NewTaskService(secret string, task_list_database_id string) (*TaskService, error) {
	return &TaskService{
		DatabaseId: task_list_database_id,
		client:     notion.NewClient(secret),
	}, nil
}

//AddTask add a new task to the notiion to-do task list
func (p *TaskService) AddTask(toDoTask Task) (notion.Page, error) {

	dbp := make(notion.DatabasePageProperties)
	//page database properties
	dbp["Status"] = notion.DatabasePageProperty{
		Select: &notion.SelectOptions{
			Name: "To Do",
		},
	}
	dbp["Name"] = notion.DatabasePageProperty{
		Title: []notion.RichText{{
			Text: &notion.Text{
				Content: toDoTask.Title,
			},
		}},
	}
	// customer properties
	for k, m := range toDoTask.CustomProperties {
		dbp[k] = notion.DatabasePageProperty{
			RichText: []notion.RichText{{
				Text: &notion.Text{
					Content: m,
				},
			}},
		}
	}

	return p.client.CreatePage(context.Background(), notion.CreatePageParams{
		ParentType:             notion.ParentTypeDatabase,
		ParentID:               p.DatabaseId,
		DatabasePageProperties: &dbp,
	})
}

//ListTasks fetch all the task from notion to-do task lists
func (p *TaskService) ListTasks() ([]Task, error) {
	var toDoTasks []Task

	response, err := p.client.QueryDatabase(context.Background(), p.DatabaseId, &notion.DatabaseQuery{

		Filter: &notion.DatabaseQueryFilter{
			Property: "Status",
			Select: &notion.SelectDatabaseQueryFilter{
				Equals: "To Do",
			},
		},
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//get the titles from the to do list
	for _, page := range response.Results {

		pageDbProperties := page.Properties.(notion.DatabasePageProperties)
		tds := Task{
			Title:            "",
			CustomProperties: map[string]string{},
		}
		for k, v := range pageDbProperties {
			if k == "Name" {
				if len(v.Title) > 0 {
					title := v.Title[len(v.Title)-1].Text.Content
					tds.Title = title
					log.Debugf("page %s id %s", title, page.ID)
				} else {
					tds.Title = ""
				}

			} else if v.RichText != nil && len(v.RichText) > 0 {
				tds.CustomProperties[k] = v.RichText[len(v.RichText)-1].Text.Content
			}
		}
		toDoTasks = append(toDoTasks, tds)

	}
	return toDoTasks, nil
}

//GetToDoTaskByTitle search the notion to-do task with the given title
func (p *TaskService) GetToDoTaskByTitle(title string) (*notion.Page, error) {
	response, err := p.client.QueryDatabase(context.Background(), p.DatabaseId, &notion.DatabaseQuery{

		Filter: &notion.DatabaseQueryFilter{
			Property: "Status",
			Select: &notion.SelectDatabaseQueryFilter{
				Equals: "To Do",
			},
		},
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//get the titles from the to do list
	page := findPageByTitleFromPages(title, response.Results)
	return page, nil
}

func (p *TaskService) MoveTaskToDoneByTitle(title string) error {
	page, err := p.GetToDoTaskByTitle(title)
	if err != nil {
		return err
	}
	if page == nil {
		return fmt.Errorf("failed to find the page with title %s", title)
	}
	dbp := make(notion.DatabasePageProperties)
	dbp["Status"] = notion.DatabasePageProperty{
		Select: &notion.SelectOptions{
			Name: "Done",
		},
	}
	_, err = p.client.UpdatePage(context.Background(), page.ID, notion.UpdatePageParams{
		DatabasePageProperties: &dbp,
	})
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//CreateToDoTasksFromPageContent will scan the notion page, find the "to_do" and create notion tasks
func (p *TaskService) CreateToDoTasksFromPageContent(page notion.Page) ([]Task, error) {
	tasks := []Task{}

	res, err := p.client.FindBlockChildrenByID(context.Background(), page.ID, &notion.PaginationQuery{})
	if err != nil {
		log.Error(err)
		return tasks, err
	}
	for _, block := range res.Results {

		//log.Debugf("block %s type %s", block.ID, block.Type)
		if block.Type == "to_do" && len(block.ToDo.RichTextBlock.Text) > 0 {
			log.Debugf("to do item: %s", block.ToDo.RichTextBlock.Text[0].Text.Content)
			tasks = append(tasks, Task{
				Title:            block.ToDo.RichTextBlock.Text[0].Text.Content,
				CustomProperties: map[string]string{},
			})
		}
		//if block.Type == "paragraph" && len(block.Paragraph.Text) > 0 {
		//	log.Debug(block.Paragraph.Text[0].Text.Content)
		//}
	}

	return tasks, nil
}

//MovePastTaskToDone check the task if it has property "End", it will move the task to "Done" status if the time has passed.
func (p *TaskService) MovePastTaskToDone(page notion.Page) error {
	pageDbProperties := page.Properties.(notion.DatabasePageProperties)

	for k, v := range pageDbProperties {
		if k == "End" && v.RichText != nil && len(v.RichText) > 0 {
			endTimestamp := v.RichText[len(v.RichText)-1].Text.Content
			log.Debugf("%s", endTimestamp)
			//is it expired?
			es, err := time.Parse("2006-01-02T15:04:05-07:00", endTimestamp)
			if err != nil {
				log.Errorf("failed to parse the end time %v", err)
				return err
			}
			t := time.Now()
			if t.After(es) {
				log.Debugf("%s %s meeting ended, move to done state", t, endTimestamp)
				dbp := make(notion.DatabasePageProperties)
				//page database properties
				dbp["Status"] = notion.DatabasePageProperty{
					Select: &notion.SelectOptions{
						Name: "Done",
					},
				}
				_, err := p.client.UpdatePage(context.Background(), page.ID, notion.UpdatePageParams{
					DatabasePageProperties: &dbp,
				})
				if err != nil {
					log.Error(err)
					return err
				}
			}
		}
	}
	return nil
}

//DoneAndCreate tasks will get all the to-do tasks, check whether the task is done (by the End property) and move them to "Done" state
//It also scans the to-do task pages, if there is any "to-do" item, then create Task for them
func (p *TaskService) DoneAndCreateTasks() error {

	response, err := p.client.QueryDatabase(context.Background(), p.DatabaseId, &notion.DatabaseQuery{

		Filter: &notion.DatabaseQueryFilter{
			Property: "Status",
			Select: &notion.SelectDatabaseQueryFilter{
				Equals: "To Do",
			},
		},
	})
	if err != nil {
		log.Error(err)
		return err
	}
	//get the titles from the to do list
	for _, page := range response.Results {

		err = p.MovePastTaskToDone(page)
		if err != nil {
			log.Error(err)
			return err
		}

		tasks, err := p.CreateToDoTasksFromPageContent(page)
		if err != nil {
			log.Error(err)
			return err
		}
		if len(tasks) > 0 {
			// iterate the to do in reverse order this way we can add the tasks to the list in the same order
			for i := len(tasks) - 1; i >= 0; i-- {
				task := tasks[i]
				//check whether the task already exists in the to-do
				existing_task := findPageByTitleFromPages(task.Title, response.Results)
				if existing_task != nil {
					log.Debug("to do task already exists, skip")
				} else {
					_, err := p.AddTask(task)
					if err != nil {
						log.Error(err)
						//return err
					}

				}
			}

		}
	}
	return nil
}

func findPageByTitleFromPages(title string, pages []notion.Page) *notion.Page {
	for _, page := range pages {
		//check the name - title
		if title == getTitleOfPage(page) {
			return &page
		}
	}
	return nil
}

func getTitleOfPage(page notion.Page) string {
	pageDbProperties := page.Properties.(notion.DatabasePageProperties)
	for k, v := range pageDbProperties {
		if k == "Name" && len(v.Title) > 0 {
			return v.Title[len(v.Title)-1].Text.Content
		}
	}
	return ""
}

//PostToNotion is to post data to Notion through API. Not Really used
/*
func (p *TaskService) PostToNotion(url string, structData interface{}) (string, error) {
	log.Info(structData)
	requestBody, err := json.Marshal(structData)
	log.Info(requestBody)
	if err != nil {
		log.Fatal(err)
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(
		"Content-Type", "application/json",
	)
	request.Header.Set(
		"Authorization", fmt.Sprintf("Bearer %s", p.Secret),
	)
	request.Header.Set(
		"Notion-Version", "2021-08-16",
	)
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body), nil
}
*/

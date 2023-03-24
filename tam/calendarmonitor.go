package automaticmanager

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	google "github.com/jyouturner/automaticmanager/pkg/google"
	"github.com/jyouturner/automaticmanager/pkg/notion"
	log "github.com/sirupsen/logrus"
)

//AddToDoFromCalendar fetch the next 10 calendar events and add to notion to do task. It also add the today journal task.
func AddToDoFromCalendar(googleClient *http.Client, notionClient *notion.TaskService, userCfg UserConfig) {

	//load the group email list
	groupEmailsMap := make(map[string]bool)

	for _, group := range strings.Split(userCfg.Calendar.ExcludeEmailsList, ",") {
		log.Debugf("will skip any event with attendee %s", group)
		groupEmailsMap[group] = true
	}

	google_calendar, err := google.NewCalendarServiceFromClient(googleClient)
	if err != nil {
		log.Fatalf("failed to connect ot google service %v", err)
	}

	events, err := google_calendar.GetNextEvents("primary", 10, groupEmailsMap)
	if err != nil {
		log.Fatalf("failed to get events from google calendar %v", err)
	}
	if len(events) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		//get the Notion to do list
		tasks, err := notionClient.ListTasks()
		if err != nil {
			log.Fatalf("failed to get notion tasks %v", err)
		}
		log.Debugf("found %d events", len(tasks))
		AddTodayTasksFromCalendar(events, notionClient, tasks)

		AddTodayJournalTask(notionClient, tasks)

	}
}

//EventToTask convert a calendar Event to a Task object
func EventToTask(event google.CalendarEvent) *notion.Task {
	return &notion.Task{
		Title: MakeTaskTitle(event),
		CustomProperties: map[string]string{
			"Source":    "Calendar",
			"Location":  event.Location,
			"Id":        event.Id,
			"VideoLink": event.VideoLink,
			"Start":     event.Start,
			"End":       event.End,
		},
	}
}

//MakeTaskTitle return the title of the task based on the calendar calendar event
func MakeTaskTitle(event google.CalendarEvent) string {
	return "Meeting: " + event.Summary
}

//doesTaskExists check whether a event already has the corresponding task in the given list
func doesTaskExist(event google.CalendarEvent, tasks []notion.Task) (bool, *notion.Task) {
	id := event.Id
	//check whether the task already exist
	for index, task := range tasks {
		if task.CustomProperties["Id"] != "" && task.CustomProperties["Id"] == id {
			//skip
			log.Debug("already exists", index)
			return true, &task
		}
	}
	return false, nil
}

//MakeTodayJournalTask create a task for current calendar day
func MakeTodayJournalTask() *notion.Task {
	today := time.Now().Format("2006-01-02")
	now := time.Now()
	end_of_day := time.Date(now.Year(), now.Month(), now.Day(), 23, 0, 0, 0, now.Location())

	return &notion.Task{
		Title: today,
		CustomProperties: map[string]string{
			"End": end_of_day.Format("2006-01-02T15:04:05-07:00"),
		},
	}
}

//doesTodayJournalTaskExists check whether there is already a task for "today"
func doesTodayJournalTaskExist(today_journal *notion.Task, tasks []notion.Task) (bool, *notion.Task) {
	for index, task := range tasks {
		if task.Title == today_journal.Title {
			//skip
			log.Debug("already exists", index)
			return true, &task
		}
	}
	return false, nil
}

//AddTodayJournalTask add the today task to the Notion task list
func AddTodayJournalTask(notion_db *notion.TaskService, tasks []notion.Task) {
	//check whether the today journal task already exists, otherwise, add it to notion
	tsk := MakeTodayJournalTask()
	exists, _ := doesTodayJournalTaskExist(tsk, tasks)
	if !exists {
		_, err := notion_db.AddTask(*tsk)
		if err != nil {
			log.Fatal(err)
		}
	}
}

//addTodayTasksFromCalendar add To Do task to Notion based on given calendar events
func AddTodayTasksFromCalendar(events []google.CalendarEvent, notion_db *notion.TaskService, tasks []notion.Task) {
	for i := len(events) - 1; i >= 0; i-- {
		event := events[i]
		exists, _ := doesTaskExist(event, tasks)

		if !exists {
			tsk := EventToTask(event)
			_, err := notion_db.AddTask(*tsk)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			//update task if necessary
			continue
		}
	}
}

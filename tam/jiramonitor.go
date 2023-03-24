package automaticmanager

import (
	"fmt"
	"net/http"

	atlanssian "github.com/jyouturner/automaticmanager/pkg/atlanssian"
	"github.com/jyouturner/automaticmanager/pkg/notion"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

//JiraToTasks will check the JIRA board and create notion tasks based on the JQL query results
func JiraToTasks(atlanssianClient *http.Client, notionClient *notion.TaskService, userCfg UserConfig) {

	jira := atlanssian.NewJiraPlatformApiClient(atlanssianClient, userCfg.Atlanssian.JiraClound.JiraCloundUrl)

	issues, _, err := jira.SearchIssuesByJQL(userCfg.Atlanssian.JiraClound.JiraMonitor.JQL)
	if err != nil {
		log.Fatalf("Failed to search issues %v", err)
	}
	if issues == nil {
		log.Errorf("no issues found")
		return
	}
	log.Info(issues.Total)
	// iterate each issue

	for _, issue := range issues.Issues {
		id := issue["id"]
		//exclude epic, and subtask

		//log.Debug(id)
		idString, ok := id.(string)
		if !ok {
			log.Errorf("%v is not string", id)
		}
		_, issueStr, err := jira.GetIssueById(idString)
		if err != nil {
			log.Errorf("Failed to get issue %v", err)
		}
		//log.Debug(issue)
		if jira.IsEpic(issueStr) {
			continue
		}
		time := jira.FindDaysOfIssueInStatus(issueStr, "In Dev")
		log.Debug(time)
		flagExists, _ := jira.FlagExists(issueStr, "customfield_10054")
		if time > 5 {
			if !flagExists {
				// flag it and create a notion task
				_, err = jira.FlagIssue(idString, "customfield_10054", "Impediment")
				if err != nil {
					log.Errorf("failed to flag issue %v", err)
				}
				// create task
				_, err = notionClient.AddTask(createJiraTask(issueStr))
				if err != nil {
					log.Errorf("failed to flag issue %v", err)
				}
			}

		} else {
			if flagExists {
				_, err = jira.DeFlagIssue(idString, "customfield_10054", "Impediment")
				if err != nil {
					log.Errorf("failed to deflag issue %v", err)
				}
			}
		}

	}

}

//createJiraTask return a notion Task struct from a JIRA issue string
func createJiraTask(issue string) notion.Task {
	return notion.Task{
		Title: fmt.Sprintf("follow up on issue %s", gjson.Parse(issue).Get("key").String()),
	}
}

//getNotificationService return the Notion client
/*
func getNotionService() (*notion.TaskService, error) {
	key, exists := os.LookupEnv("NOTION_KEY")
	if !exists {
		log.Fatal("missing env varioable NOTION_KEY")
	}
	taskService, err := notion.NewTaskService(key, os.Getenv("NOTION_DATABASE_ID"))
	if err != nil {
		return nil, fmt.Errorf("failed to create notion client %v", err)
	}
	return taskService, nil
}
*/
//NewAtlanssianOauth2Client creates the http client that handles the oauth2 auth with Atlanssian 3-leg auth
/*
func NewAtlanssianOauth2Client(ats oauthext.TokenStorage) (*http.Client, error) {
	scope := strings.Split(os.Getenv("ATLASSIAN_OAUTH2_SCOPE"), ",")
	config := &oauth2.Config{
		ClientID:     os.Getenv("ATLASSIAN_OAUTH2_CLIENT_ID"),
		ClientSecret: os.Getenv("ATLASSIAN_OAUTH2_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("ATLASSIAN_OAUTH2_REDIRECT_URL"),
		Scopes:       scope,
		Endpoint: oauth2.Endpoint{
			AuthURL:  os.Getenv("ATLASSIAN_OAUTH2_AUTH_URL"),
			TokenURL: os.Getenv("ATLASSIAN_OAUTH2_TOKEN_URL"),
		},
	}

	client, err := oauthext.NewClient(ats, config)
	if err != nil {
		log.Errorf("Failed to get oauth2 client, do you have the token file in place?")
		return nil, err
	}
	return client, err
}
*/

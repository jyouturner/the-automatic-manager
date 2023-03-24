package automaticmanager

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	atlanssian "github.com/jyouturner/automaticmanager/pkg/atlanssian"
	"github.com/jyouturner/automaticmanager/pkg/github"
	"github.com/jyouturner/automaticmanager/pkg/google"
	"github.com/jyouturner/automaticmanager/pkg/notion"
	log "github.com/sirupsen/logrus"
)

func MonitorGitHubPullRequests(githubClient *github.GitHubService, atlanssianClient *http.Client, notionClient *notion.TaskService, gmail *google.MailService, userCfg UserConfig) {
	// first get the open PRs
	prs, err := githubClient.SearchOpenPullRequests("jyouturner", "notion-integration")
	if err != nil {
		log.Fatal(err)
	}
	jira := atlanssian.NewJiraPlatformApiClient(atlanssianClient, userCfg.Atlanssian.JiraClound.JiraCloundUrl)

	// check each pr, find the JIRA ticket
	for _, pr := range prs {
		title := strings.TrimLeft(*pr.Title, " ")
		if strings.HasPrefix(title, "[") {
			// jira
			re := regexp.MustCompile("\\[(.*?)\\]")
			jiraIssueKey := re.FindStringSubmatch(title)[1]
			log.Info(jiraIssueKey)
			// now find the jira story
			issue, _, err := jira.GetIssueById(jiraIssueKey)
			if err != nil {
				log.Info("Failed to find the jira issue %v", err)
			}
			if issue == nil {
				log.Errorf("no issue found")
			}
			// find the owner of the issue
			fmt.Println(issue)
			//create a follow up task of the JIRA
			notionClient.AddTask(notion.Task{
				Title: fmt.Sprintf("Follow up of PR - %s", jiraIssueKey),
			})
			// send message to slack through email
			gmail.SendMail("jerry.you@snapdocs.com", fmt.Sprintf("Follow up of PR - %s", jiraIssueKey), "How is the PR?")
			// or send slack through web url
			// https://hooks.slack.com/workflows/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
			//
			message := fmt.Sprintf(`
			{
				"user_email": "jerry.you@snapdocs.com",
				"message": "Follow up of PR - %s"
			}`, jiraIssueKey)
			//send message
			doPost("https://hooks.slack.com/workflows/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", message)
		}
	}
}

func doPost(url string, data string) (string, error) {

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	request, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(
		"Content-Type", "application/json",
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

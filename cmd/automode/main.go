package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/joho/godotenv"
	"github.com/jyouturner/automaticmanager/pkg/github"
	"github.com/jyouturner/automaticmanager/pkg/google"
	"github.com/jyouturner/automaticmanager/pkg/notion"
	automaticmanager "github.com/jyouturner/automaticmanager/tam"
	"github.com/jyouturner/gotoauth"
	"github.com/jyouturner/gotoauth/example/awsserverless"
	log "github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Staring workflow ....")
	time.Sleep(3 * time.Second)
	fmt.Println("ready")
	//reader := bufio.NewReader(os.Stdin)
	fmt.Println("Press anykey to run the calendar monitor job ...")
	//_, _ := reader.ReadString('\n')
	fmt.Scanln()
	//check google calendar
	WatchCalendar()
	//
	fmt.Println("Press anykey to run the notion monitor job ...")
	fmt.Scanln()
	WatchTasks()
	// jira monitor to create stories

	// github

	fmt.Println("Press anykey to run the github monitor job ...")
	fmt.Scanln()
	WatchGithub()
	for {
		fmt.Println("...")
		time.Sleep(10 * time.Second)
	}
}

const AWS_SECRET_NAME_ENV_NAME string = "AWS_SECRET_NAME"

func init() {
	//load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	_, exists := os.LookupEnv(AWS_SECRET_NAME_ENV_NAME)
	if !exists {
		log.Fatalf("missing %s", AWS_SECRET_NAME_ENV_NAME)
	}
	logLevel, exists := os.LookupEnv("LOG_LEVEL")
	if exists {
		level, err := log.ParseLevel(logLevel)
		if err != nil {
			log.Errorf("incorrect LOG_LEVEL %s", level)
		} else {
			log.SetLevel(level)
		}
	}
	_, exists = os.LookupEnv("AWS_SECRET_NAME")
	if !exists {
		log.Fatalf("missing %s", "AWS_SECRET_NAME")
	}

}

func createAwsClient() (*awsserverless.AWSClient, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to create aws session, %v", err)

	}
	awsClient := &awsserverless.AWSClient{
		Config: cfg,
	}
	return awsClient, nil
}

//WatchCalendar is the lambda function to monitor calendar and create to do tasks from the events.
func WatchCalendar() {
	//TODO need to figure out how to handle multiple users
	user := automaticmanager.NewUser("abcde", "12345")

	awsClient, err := createAwsClient()
	if err != nil {
		log.Fatal(err)
	}
	authProvider := automaticmanager.GOOGLE
	awsEnv, err := awsserverless.NewAWSEnvByUser(*awsClient, os.Getenv("AWS_SECRET_NAME"), user.BucketName, user.Id, "")
	if err != nil {
		log.Fatalf("error create aws session %v", err)
	}
	authconfig, err := awsEnv.GetAppOathConfig(authProvider)
	if err != nil {
		log.Fatal("failed to get auth config for %s %v", authProvider, err)
	}
	oauthConfig, err := gotoauth.ConfigFromJSON(authconfig.Secret, strings.Split(automaticmanager.ProviderScope[authProvider], " "))
	if err != nil {
		log.Fatalf("error loading config of auth provider %v", err)
	}
	httpClient, err := gotoauth.NewClient(authconfig.OauthTokenStorage, oauthConfig)

	if err != nil {
		log.Fatalf("error create http client %v", err)
	}
	// load user's config.yml file from s3
	userCfg, err := automaticmanager.GetUserConfigFromS3(*awsClient, user)
	if err != nil {
		log.Fatalf("failed to get user configuration %v", err)
	}
	notionClient, _ := notion.NewTaskService(userCfg.Notion.ApiKey, userCfg.Notion.TaskDatabaseId)
	automaticmanager.AddToDoFromCalendar(httpClient, notionClient, *userCfg)
}

//WatchTasks is the lambda function to monitor Notion to do tasks
func WatchTasks() {
	//TODO need to figure out how to handle multiple users
	user := automaticmanager.NewUser("abcde", "12345")
	log.Debug(user)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	awsClient := &awsserverless.AWSClient{
		Config: cfg,
	}

	// load user's config.yml file from s3
	userCfg, err := automaticmanager.GetUserConfigFromS3(*awsClient, user)
	if err != nil {
		log.Fatalf("failed to get user configuration %v", err)
	}
	notionClient, _ := notion.NewTaskService(userCfg.Notion.ApiKey, userCfg.Notion.TaskDatabaseId)
	automaticmanager.ProcessToDoTasks(notionClient)
	// wait 2 sec
	err = notionClient.MoveTaskToDoneByTitle("Meeting: Discuss the milestone 2 of enote api")
	if err != nil {
		log.Fatalf("failed to move the task to done %v", err)
	}
}

func MonitorGitHubPullRequests(githubClient *github.GitHubService, atlanssianClient *http.Client, notionClient *notion.TaskService, gmail *google.MailService, userCfg automaticmanager.UserConfig) {
	// first get the open PRs
	prs, err := githubClient.SearchOpenPullRequests("jyouturner", "notion-integration")
	if err != nil {
		log.Fatal(err)
	}
	//jira := atlanssian.NewJiraPlatformApiClient(atlanssianClient, userCfg.Atlanssian.JiraClound.JiraCloundUrl)

	// check each pr, find the JIRA ticket
	for _, pr := range prs {
		title := strings.TrimLeft(*pr.Title, " ")
		if strings.HasPrefix(title, "[") {
			// jira

			re := regexp.MustCompile("\\[(.*?)\\]")
			jiraIssueKey := re.FindStringSubmatch(title)[1]
			log.Info(jiraIssueKey)
			// now find the jira story
			/**
			issue, _, err := jira.GetIssueById(jiraIssueKey)
			if err != nil {
				log.Info("Failed to find the jira issue %v", err)
			}
			if issue == nil {
				log.Errorf("no issue found")
			}
			// find the owner of the issue
			fmt.Println(issue)
			*/
			//create a follow up task of the JIRA
			notionClient.AddTask(notion.Task{
				Title: fmt.Sprintf("Follow up of PR - %s", jiraIssueKey),
			})
			// send message to slack through email
			//gmail.SendMail("jerry.you@snapdocs.com", fmt.Sprintf("Follow up of PR - %s", jiraIssueKey), "How is the PR?")
			// or send slack through web url
			// https://hooks.slack.com/workflows/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
			//
			message := fmt.Sprintf(`
			{
				"user_email": "john.dyer@snapdocs.com",
				"message": "Hey, Snapdocs release cut is tomorrow, seems the PR of %s is still in progress?"
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

func WatchGithub() {

	//TODO need to figure out how to handle multiple users
	user := automaticmanager.NewUser("abcde", "12345")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	awsClient := &awsserverless.AWSClient{
		Config: cfg,
	}
	atlassianClient := getAtlassianClient(awsClient, user)
	googleClient := getGoogleClient(awsClient, user)
	gmailService, err := google.NewMailService(googleClient)
	if err != nil {
		log.Fatal(err)
	}
	// load user's config.yml file from s3
	userCfg, err := automaticmanager.GetUserConfigFromS3(*awsClient, user)
	if err != nil {
		log.Fatalf("failed to get user configuration %v", err)
	}
	notionClient, err := notion.NewTaskService(userCfg.Notion.ApiKey, userCfg.Notion.TaskDatabaseId)

	if err != nil {
		log.Fatal(err)
	}
	//userCfg, err := automaticmanager.GetUserConfigFromLocalFile("config/user/config.yml")
	//if err != nil {
	//	log.Fatal(err)
	//}

	githubClient := github.NewGitHubServiceWithAccessToken("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")

	MonitorGitHubPullRequests(githubClient, atlassianClient, notionClient, gmailService, *userCfg)

}

func getAtlassianClient(awsClient *awsserverless.AWSClient, user automaticmanager.User) *http.Client {
	authProvider := automaticmanager.ATLANSSIAN
	awsEnv, err := awsserverless.NewAWSEnvByUser(*awsClient, os.Getenv("AWS_SECRET_NAME"), user.BucketName, user.Id, "")
	if err != nil {
		log.Fatalf("error create aws session %v", err)
	}
	authconfig, err := awsEnv.GetAppOathConfig(authProvider)
	if err != nil {
		log.Fatal("failed to get auth config for %s %v", authProvider, err)
	}
	oauthConfig, err := gotoauth.ConfigFromJSON(authconfig.Secret, strings.Split(automaticmanager.ProviderScope[authProvider], " "))
	if err != nil {
		log.Fatalf("error loading config of auth provider %v", err)
	}
	atlassianClient, err := gotoauth.NewClient(authconfig.OauthTokenStorage, oauthConfig)

	if err != nil {
		log.Fatalf("error create http client %v", err)
	}

	return atlassianClient
}
func getGoogleClient(awsClient *awsserverless.AWSClient, user automaticmanager.User) *http.Client {
	authProvider := automaticmanager.GOOGLE
	awsEnv, err := awsserverless.NewAWSEnvByUser(*awsClient, os.Getenv("AWS_SECRET_NAME"), user.BucketName, user.Id, "")
	if err != nil {
		log.Fatalf("error create aws session %v", err)
	}
	authconfig, err := awsEnv.GetAppOathConfig(authProvider)
	if err != nil {
		log.Fatal("failed to get auth config for %s %v", authProvider, err)
	}
	oauthConfig, err := gotoauth.ConfigFromJSON(authconfig.Secret, strings.Split(automaticmanager.ProviderScope[authProvider], " "))
	if err != nil {
		log.Fatalf("error loading config of auth provider %v", err)
	}
	client, err := gotoauth.NewClient(authconfig.OauthTokenStorage, oauthConfig)

	if err != nil {
		log.Fatalf("error create http client %v", err)
	}
	return client
}

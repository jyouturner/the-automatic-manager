package main

import (
	"context"
	"net/http"
	"os"
	"strings"

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

}

func main() {

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
	log.Println(userCfg.Github)

	githubClient := github.NewGitHubServiceWithAccessToken("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")

	automaticmanager.MonitorGitHubPullRequests(githubClient, atlassianClient, notionClient, gmailService, *userCfg)

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

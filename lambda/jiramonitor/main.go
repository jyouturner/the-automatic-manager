package main

import (
	"context"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/jyouturner/gotoauth"
	"github.com/jyouturner/gotoauth/example/awsserverless"

	"github.com/jyouturner/automaticmanager/pkg/notion"
	automaticmanager "github.com/jyouturner/automaticmanager/tam"
)

const AWS_SECRET_NAME_ENV_NAME string = "AWS_SECRET_NAME"

func init() {
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

func main() {

	lambda.Start(WatchJira)
}

func WatchJira() {
	//TODO need to figure out how to handle multiple users
	user := automaticmanager.NewUser("abcde", "12345")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	awsClient := &awsserverless.AWSClient{
		Config: cfg,
	}

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
	automaticmanager.JiraToTasks(httpClient, notionClient, *userCfg)
}

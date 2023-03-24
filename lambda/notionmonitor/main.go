package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/jyouturner/automaticmanager/pkg/notion"
	automaticmanager "github.com/jyouturner/automaticmanager/tam"
	"github.com/jyouturner/gotoauth/example/awsserverless"
)

func init() {

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

	lambda.Start(WatchTasks)
}

//WatchTasks is the lambda function to monitor Notion to do tasks
func WatchTasks(ctx context.Context, cloundWatchEvent events.CloudWatchEvent) {
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
}

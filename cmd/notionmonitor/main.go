package main

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/jyouturner/automaticmanager/pkg/notion"
	automaticmanager "github.com/jyouturner/automaticmanager/tam"
	"github.com/jyouturner/gotoauth/example/awsserverless"
)

func init() {
	//load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
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
	for {
		WatchTasks()
		time.Sleep(10 * time.Second)
	}

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
}

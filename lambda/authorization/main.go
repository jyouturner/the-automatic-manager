package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"

	lambdahelper "github.com/jyouturner/automaticmanager/lambda"
	automaticmanager "github.com/jyouturner/automaticmanager/tam"
	"github.com/jyouturner/gotoauth"
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
	_, exists = os.LookupEnv("AWS_SECRET_NAME")
	if !exists {
		log.Fatalf("missing %s", "AWS_SECRET_NAME")
	}

	_, exists = os.LookupEnv("OAUTH_PROVIDER")
	if !exists {
		log.Fatalf("missing %s", "OAUTH_PROVIDER")
	}
	_, exists = os.LookupEnv("AUTHORIZED_TO_URL")
	if !exists {
		log.Fatalf("missing %s", "AUTHORIZED_TO_URL")
	}

	_, exists = os.LookupEnv("OAUTH_NOUNCE_BUCKET")
	if !exists {
		log.Fatalf("missing %s", "OAUTH_NOUNCE_BUCKET")
	}

	if automaticmanager.ProviderScope[strings.ToUpper(os.Getenv("OAUTH_PROVIDER"))] == "" {
		log.Fatalf("missing %s in ProviderScope", strings.ToUpper(os.Getenv("OAUTH_PROVIDER")))
	}
}

func main() {
	lambda.Start(Handle)
}

func Handle(ctx context.Context, event json.RawMessage) (lambdahelper.LambdaResponse, error) {
	fmt.Println(string(event))
	eventBodyString := gjson.Get(string(event), "body").String()
	fmt.Println(string(eventBodyString))
	//get user data from body
	/*
		{
		"user": {
			"org_id": "12345",
			"user_id": "abcde"
		}
	*/
	userData := gjson.Get(eventBodyString, "user").String()
	orgUser, err := awsserverless.FromJson([]byte(userData))
	if err != nil {
		log.Errorf("failed to convert json to user meta %v %v", userData, err)
		return lambdahelper.FailureMessage(400, "user data is not expected"), err
	}
	user := automaticmanager.NewUserWithOrg(orgUser)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return lambdahelper.FailureMessage(500, "system error"), fmt.Errorf("system error")

	}

	awsClient := awsserverless.AWSClient{
		Config: cfg,
	}

	nounceState := awsserverless.StateToken{
		User:               orgUser,
		Provider:           strings.ToUpper(os.Getenv("OAUTH_PROVIDER")),
		Scope:              automaticmanager.ProviderScope[strings.ToUpper(os.Getenv("OAUTH_PROVIDER"))],
		SuccessRedirectUrl: os.Getenv("AUTHORIZED_TO_URL"),
	}
	awsEnv, err := awsserverless.NewAWSEnvByUser(awsClient, os.Getenv("AWS_SECRET_NAME"), user.BucketName, orgUser, os.Getenv("OAUTH_NOUNCE_BUCKET"))
	if err != nil {
		return lambdahelper.FailureMessage(500, "could not load oauth config from aws"), err
	}

	authurl, err := gotoauth.GetAuthUrl(nounceState, awsEnv, awsEnv)

	if err != nil {
		log.Errorf("failed to initialize authorization flow for user %v %v", orgUser, err)
		return lambdahelper.FailureMessage(500, "failed to get auth url"), err
	}
	return lambdahelper.Success(*authurl), nil
}

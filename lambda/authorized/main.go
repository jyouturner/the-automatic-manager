package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

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

	_, exists = os.LookupEnv("OAUTH_NOUNCE_BUCKET")
	if !exists {
		log.Fatalf("missing %s", "OAUTH_NOUNCE_BUCKET")
	}
}

func main() {
	lambda.Start(Handle)
}

//Handle function can handle the oauth redirect payload from the oauth provider, the http request parameter is like
//state=...&code=...&scope=https://www.googleapis.com/auth/drive.metadata.readonly
func Handle(ctx context.Context, event json.RawMessage) (lambdahelper.LambdaResponse, error) {
	fmt.Println(string(event))
	qsp := gjson.Get(string(event), "queryStringParameters")
	if !qsp.IsObject() {
		//wrong
		return lambdahelper.FailureMessage(500, "failed to get access token"), fmt.Errorf("no query string parameters found")
	}

	authcode := qsp.Get("code").String()
	nounce := qsp.Get("state").String()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return lambdahelper.FailureMessage(500, "system error"), fmt.Errorf("system error")

	}
	awsClient := awsserverless.AWSClient{
		Config: cfg,
	}

	// get the auth state data from the nounce
	b, err := awsClient.S3Get(os.Getenv("OAUTH_NOUNCE_BUCKET"), nounce)
	if err != nil {
		return lambdahelper.FailureMessage(500, "could not find the matching nounce"), err
	}
	stateData := awsserverless.StateTokenFromBytes(b)
	user := automaticmanager.NewUserWithOrg(stateData.User)
	awsEnv, err := awsserverless.NewAWSEnvByUser(awsClient, os.Getenv("AWS_SECRET_NAME"), user.BucketName, stateData.User, os.Getenv("OAUTH_NOUNCE_BUCKET"))
	if err != nil {
		return lambdahelper.FailureMessage(500, "could not load oauth config from aws"), err
	}
	err = gotoauth.Exchange(authcode, nounce, awsEnv, awsEnv)
	if err != nil {
		log.Error(err)
		return lambdahelper.FailureMessage(500, "failed to get access token"), err
	}
	successRedirectUrl := stateData.SuccessRedirectUrl
	if successRedirectUrl == "" {
		successRedirectUrl = "http://localhost"
	}
	return lambdahelper.Redirect(301, successRedirectUrl), nil
}

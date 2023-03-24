package integrationtesting

/*

 */

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jyouturner/gotoauth/example/awsserverless"
)

func getAWSClient(profile string, t *testing.T) awsserverless.AWSClient {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile))
	if err != nil {
		t.Errorf("failed to create aws session, %v", err)

	}
	awsClient := awsserverless.AWSClient{
		Config: cfg,
	}
	if err != nil {
		t.Error("could not create aws client")
	}
	return awsClient
}

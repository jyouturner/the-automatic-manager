package integrationtesting

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	automaticmanager "github.com/jyouturner/automaticmanager/tam"
	"github.com/jyouturner/gotoauth"
	"github.com/jyouturner/gotoauth/example/awsserverless"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

//do execute the http request, and check the returned status
func do(client http.Client, request *http.Request) (string, error) {

	request.Header.Add("Content-Type", "application/json")
	//request.Header.Add("Authorization", p.AuthorizationHeader)

	resp, err := client.Do(request)
	if err != nil {
		log.Error(err)
		return "", fmt.Errorf("failed to make request %s %v", request.URL, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("filed to read response from %s %v", request.URL, err)
	}
	return string(body), nil
}

func TestConflurnceAWS_GetUserSpacesWithBasicAuth(t *testing.T) {

	user := automaticmanager.NewUser("basic", "tam-org-12345")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	awsClient := &awsserverless.AWSClient{
		Config: cfg,
	}

	authProvider := automaticmanager.ATLANSSIAN
	awsEnv, err := awsserverless.NewAWSEnvByUser(*awsClient, os.Getenv("AWS_SECRET_NAME"), os.Getenv("TOKEN_BUCKET"), user.Id, os.Getenv("NOUNCE_TOKEN_BUCKET"))
	if err != nil {
		log.Fatalf("error create aws session %v", err)
	}
	authconfig, err := awsEnv.GetAppOathConfig(authProvider)
	if err != nil {
		log.Fatalf("failed to get auth config for %s %v", authProvider, err)
	}
	oauthConfig, err := gotoauth.ConfigFromJSON(authconfig.Secret, strings.Split(automaticmanager.ProviderScope[authProvider], " "))
	if err != nil {
		log.Fatalf("error loading config of auth provider %v", err)
	}
	httpClient, err := gotoauth.NewClient(authconfig.OauthTokenStorage, oauthConfig)

	if err != nil {
		log.Fatalf("error create http client %v", err)
	}

	//do the rest api call
	request, err := http.NewRequest("GET", "https://snapdocs-eng.atlassian.net/wiki/rest/api/user/current", nil)
	if err != nil {
		log.Fatal(err)
	}

	body, err := do(*httpClient, request)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}

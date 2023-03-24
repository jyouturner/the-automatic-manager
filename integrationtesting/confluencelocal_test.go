package integrationtesting

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jyouturner/automaticmanager/pkg/atlanssian"
	automaticmanager "github.com/jyouturner/automaticmanager/tam"
	"github.com/jyouturner/gotoauth/example/local"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func TestConflurnceLocal_GetUserBasicAuth(t *testing.T) {
	ats := local.LocalTokenStorage{
		TokenFile: "testdata/atlanssian_basicauth.json",
	}
	token, err := ats.LoadToken()
	if err != nil {
		t.Fail()
	}
	userCfg, err := automaticmanager.GetUserConfigFromLocalFile("testdata/config.yml")
	if err != nil {
		log.Fatal(err)
	}
	//client := atlanssian.NewConfluenceApiApiClientBasicAuth(userCfg.Atlanssian.Confluence.ConfluenceUrl, userCfg.Atlanssian.Confluence.BasicAuthUser, userCfg.Atlanssian.Confluence.BasicAuthToken)
	httpClient := oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(token))
	client := atlanssian.ConfluenceApiApiClient{
		Client: httpClient,
		ApiUrl: userCfg.Atlanssian.Confluence.ConfluenceUrl,
	}
	res, err := client.GetUserSpace()
	if err != nil {
		t.Fail()
	}
	fmt.Println(res)
}

func TestConflurnceLocal_GetUserBasicAuth2(t *testing.T) {

	request, err := http.NewRequest("GET", fmt.Sprintf("https://%s/wiki/rest/api/space", os.Getenv("DOMAIN")), nil)
	if err != nil {
		t.Fail()
	}
	user := os.Getenv("USER")
	token := os.Getenv("TOKEN")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, token)))))
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}

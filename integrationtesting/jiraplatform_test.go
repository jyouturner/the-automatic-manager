package integrationtesting

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	atlanssian "github.com/jyouturner/automaticmanager/pkg/atlanssian"
	automaticmanager "github.com/jyouturner/automaticmanager/tam"
	"github.com/jyouturner/gotoauth"
	"github.com/jyouturner/gotoauth/example/local"
	"github.com/tidwall/gjson"
)

func getHttpClient(t *testing.T) *http.Client {
	//use local file to store tokens
	ats := local.LocalTokenStorage{
		TokenFile: "testdata/atlanssian_token.json",
	}
	authconfig, err := local.ConfigFromLocalJsonFile("testdata/atlanssian_secret.json", []string{"offline_access", "read:jira-user", "read:jira-work", "write:jira-work", "read:confluence-user", "write:confluence-content", "read:confluence-content.all", "read:confluence-space.summary"})
	if err != nil {
		t.Errorf("failed to create the oauth %v", err)
	}
	client, err := gotoauth.NewClient(ats, authconfig)
	if err != nil {
		t.Errorf("failed to create the client %v", err)
	}
	return client

}

func TestJiraPlatformApiClient_GetIssueById(t *testing.T) {
	userCfg, err := automaticmanager.GetUserConfigFromLocalFile("testdata/config.yml")
	if err != nil {
		log.Fatal(err)
	}
	jira := atlanssian.NewJiraPlatformApiClient(getHttpClient(t), userCfg.Atlanssian.JiraClound.JiraCloundUrl)
	_, got1, err := jira.GetIssueById("56276")
	if err != nil {
		t.Errorf("JiraPlatformApiClient.GetIssueById() error = %v", err)
	}
	//log.Info(got)
	//log.Info(got1)
	fmt.Printf("gjson.Parse(got1): %v\n", gjson.Parse(got1))

	//check whether flag exists
	flagExists, err := jira.FlagExists(got1, "customfield_10054")
	if err != nil {
		t.Error("failed to inspect the issue")
	}
	fmt.Printf("flag exists %v", flagExists)
}

func TestJiraPlatformApiClient_GetIssueByKey(t *testing.T) {
	userCfg, err := automaticmanager.GetUserConfigFromLocalFile("testdata/config.yml")
	if err != nil {
		log.Fatal(err)
	}
	jira := atlanssian.NewJiraPlatformApiClient(getHttpClient(t), userCfg.Atlanssian.JiraClound.JiraCloundUrl)
	_, got1, err := jira.GetIssueById("SCON-521")
	if err != nil {
		t.Errorf("JiraPlatformApiClient.GetIssueById() error = %v", err)
	}
	//log.Info(got)
	//log.Info(got1)
	fmt.Printf("gjson.Parse(got1): %v\n", gjson.Parse(got1))
}

func TestJiraPlatformApiClient_FlagIssue(t *testing.T) {
	userCfg, err := automaticmanager.GetUserConfigFromLocalFile("testdata/config.yml")
	if err != nil {
		log.Fatal(err)
	}
	jira := atlanssian.NewJiraPlatformApiClient(getHttpClient(t), userCfg.Atlanssian.JiraClound.JiraCloundUrl)
	if _, err := jira.FlagIssue("56276", "customfield_10054", "Impediment"); err != nil {
		t.Errorf("JiraPlatformApiClient.FlagIssue() error = %v", err)
	}

}

func TestJiraPlatformApiClient_GetIssueEditMetaData(t *testing.T) {
	userCfg, err := automaticmanager.GetUserConfigFromLocalFile("testdata/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	jira := atlanssian.NewJiraPlatformApiClient(getHttpClient(t), userCfg.Atlanssian.JiraClound.JiraCloundUrl)
	got, err := jira.GetIssueEditMeta("56276")
	if err != nil {
		t.Errorf("JiraPlatformApiClient.GetIssueEditMeta() error = %v", err)
	}
	fmt.Printf("gjson.Parse(got): %v\n", gjson.Parse(got))

}

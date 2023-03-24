package integrationtesting

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/jyouturner/automaticmanager/pkg/github"
	automaticmanager "github.com/jyouturner/automaticmanager/tam"
)

/**
func TestGitHubService_ListRepo(t *testing.T) {
	//use local file to store tokens
	ats := gotoauth.LocalTokenStorage{
		TokenFile: "testdata/github_token.json",
	}
	authconfig, err := gotoauth.ConfigFromLocalJsonFile("testdata/github.json", []string{"repo"})
	if err != nil {
		t.Errorf("failed to create the oauth %v", err)
	}
	client, err := gotoauth.NewClient(ats, authconfig)
	if err != nil {
		t.Errorf("failed to create the client %v", err)
	}
	s := github.NewGitHubService(client)

	events, err := s.ListRepo()
	if err != nil {
		t.Error(err)
	}

	for _, event := range events {
		log.Println(event)
	}
}
*/
func TestGitHubService_ListRepo2(t *testing.T) {
	//use local file to store tokens
	userCfg, err := automaticmanager.GetUserConfigFromLocalFile("testdata/config.yml")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(userCfg.Github)
	s := github.NewGitHubServiceWithAccessToken(userCfg.Github.AccessToken)

	events, err := s.ListRepo()
	if err != nil {
		t.Error(err)
	}

	for _, event := range events {
		log.Println(event)
	}
}

func TestGitHubService_GetOpenPrs(t *testing.T) {
	//use local file to store tokens
	userCfg, err := automaticmanager.GetUserConfigFromLocalFile("testdata/config.yml")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(userCfg.Github)
	s := github.NewGitHubServiceWithAccessToken(userCfg.Github.AccessToken)

	prs, err := s.SearchOpenPullRequests("jyouturner", "notion-integration")
	if err != nil {
		t.Error(err)
	}

	for _, pr := range prs {
		log.Println(pr.ID)
		title := strings.TrimLeft(*pr.Title, " ")
		log.Println(title)
		re := regexp.MustCompile("\\[(.*?)\\]")
		match := re.FindStringSubmatch(title)
		fmt.Println(match[1])
	}
}

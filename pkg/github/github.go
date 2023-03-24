package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

type GitHubService struct {
	Client *github.Client
}

func NewGitHubService(httpClient *http.Client) *GitHubService {
	client := github.NewClient(httpClient)
	return &GitHubService{

		Client: client,
	}
}

func NewGitHubServiceWithAccessToken(token string) *GitHubService {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)
	return &GitHubService{

		Client: client,
	}
}

func (p GitHubService) ListRepo() ([]*github.Repository, error) {
	reps, _, err := p.Client.Repositories.ListAll(context.Background(), &github.RepositoryListAllOptions{})
	if err != nil {
		return nil, fmt.
			Errorf("failed to list repos %v", err)
	}
	for _, repo := range reps {
		fmt.Println(*repo.Name)
	}
	return reps, nil
}

func (p GitHubService) SearchOpenPullRequests(owner string, repo string) ([]*github.PullRequest, error) {
	prs, _, err := p.Client.PullRequests.List(context.Background(), owner, repo, &github.PullRequestListOptions{})
	if err != nil {
		return nil, fmt.
			Errorf("failed to list repos %v", err)
	}
	//iterate each pr and filter
	for _, pr := range prs {
		if *pr.State == "open" {
			fmt.Println(*pr.Title)
		}
	}
	return prs, nil
}

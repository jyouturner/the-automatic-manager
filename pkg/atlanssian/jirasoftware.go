//
// To create personal api access token, go to https://id.atlassian.com/manage-profile/security/api-tokens
//
package atlanssian

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

//JiraSoftwareApiClient can request the Jira Software Cloud Rest API at https://developer.atlanssian.com/cloud/jira/software/rest/intro/
type JiraSoftwareApiClient struct {
	ApiUrl              string
	Client              *http.Client
	AuthorizationHeader string
}

//NewJiraSoftwareApiClient returns a client of the Jira Software with given basic authentication credential
func NewJiraSoftwareApiClient(apiUrl, user string, token string) *JiraSoftwareApiClient {
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	return &JiraSoftwareApiClient{
		ApiUrl:              apiUrl,
		Client:              client,
		AuthorizationHeader: fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, token)))),
	}

}

//do execute the http request, and check the returned status
func (p *JiraSoftwareApiClient) do(request *http.Request) (string, error) {

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", p.AuthorizationHeader)

	resp, err := p.Client.Do(request)
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

//doGet execute the GET request and check the return status code
func (p *JiraSoftwareApiClient) doGet(url string) (string, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := p.do(request)
	if err != nil {
		return "", err
	}
	return resp, nil
}

//GetBoard return the JIRA board
func (p *JiraSoftwareApiClient) GetBoard(id string) (string, error) {
	url := p.ApiUrl + fmt.Sprintf("/rest/agile/1.0/board/%s", id)

	return p.doGet(url)
}

//GetIssuesOfBoard return the issues of the JIRA board
func (p *JiraSoftwareApiClient) GetIssuesOfBoard(id string) (string, error) {
	url := p.ApiUrl + fmt.Sprintf("/rest/agile/1.0/board/%s/issue", id)
	return p.doGet(url)
}

//GetIssue return the issue by given id
func (p *JiraSoftwareApiClient) GetIssue(id string) (string, error) {
	url := p.ApiUrl + fmt.Sprintf("/rest/agile/1.0/issue/%s", id)
	return p.doGet(url)
}

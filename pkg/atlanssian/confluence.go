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

type ConfluenceApiApiClient struct {
	ApiUrl              string
	Client              *http.Client
	AuthorizationHeader string
}

//func NewConfluenceApiApiClientAccessToken(apiUri string, client *http.Client, accessToken string) *ConfluenceApiApiClient {
//
//}

//NewConfluenceApiApiClient returns a client with given basic authentication credential
func NewConfluenceApiApiClientBasicAuth(apiUrl, user string, token string) *ConfluenceApiApiClient {
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	return &ConfluenceApiApiClient{
		ApiUrl:              apiUrl,
		Client:              client,
		AuthorizationHeader: fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, token)))),
	}

}

//do execute the http request, and check the returned status
func (p *ConfluenceApiApiClient) do(request *http.Request) (string, error) {

	request.Header.Add("Content-Type", "application/json")
	//request.Header.Add("Authorization", p.AuthorizationHeader)

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
func (p *ConfluenceApiApiClient) doGet(url string) (string, error) {
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

//get user space
func (p *ConfluenceApiApiClient) GetCurrentUser() (string, error) {
	url := p.ApiUrl + "/wiki/rest/api/user/current"
	return p.doGet(url)
}

//get space
func (p *ConfluenceApiApiClient) GetUserSpace() (string, error) {
	url := p.ApiUrl + "/wiki/rest/api/space"
	return p.doGet(url)
}

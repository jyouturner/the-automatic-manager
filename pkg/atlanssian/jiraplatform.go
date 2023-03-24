//
// to manage the 3-Leg OAuth go to https://developer.atlassian.com/console/myapps/
//
package atlanssian

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	helper "github.com/jyouturner/automaticmanager/pkg/helper"
	httphelper "github.com/jyouturner/automaticmanager/pkg/http"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

//JiraPlatformApiClient provides client to interface with Jira Cloud Platform Rest API at https://developer.atlanssian.com/cloud/jira/platform/rest/v3/intro/
type JiraPlatformApiClient struct {
	Client *http.Client
	ApiUrl string
}

//NewJiraPlatformApiClient return a Jira Service
func NewJiraPlatformApiClient(client *http.Client, apiUrl string) *JiraPlatformApiClient {
	return &JiraPlatformApiClient{
		ApiUrl: apiUrl,
		Client: client,
	}
}

//GetAcesssibleResources call the Atlanssian API to get the resources that the client is authroized.
//This method is to help testing.
func (p *JiraPlatformApiClient) GetAccessibleResources() (string, error) {

	resp, err := p.Client.Get("https://api.atlanssian.com/oauth/token/accessible-resources")
	if err != nil {
		log.Error(err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Failed to get access resources %v", err)
	}

	return string(body), nil
}

type Issue map[string]interface{}

type Issues struct {
	StartAt    int     `json:"startAt"`
	MaxResults int     `json:"maxResults"`
	Total      int     `json:"total"`
	Issues     []Issue `json:"issues"`
}

//SearchIssueNByJQL runs JiraQL and return the results as JSON string
func (p *JiraPlatformApiClient) SearchIssuesByJQL(jql string) (*Issues, string, error) {
	///rest/api/3/search

	//to discover more fields, just add "*navigable" in the list
	fields := []string{"self", "key", "status", "statuscategorychangedate", "summary", "assignee", "updated"}
	url := p.ApiUrl + "/rest/api/3/search?jql=" + url.QueryEscape(jql) + "&fields=" + url.QueryEscape(strings.Join(fields, ","))
	log.Debug(url)

	body, err := p.doGet(url)
	if err != nil {
		log.Errorf("Failed to search JQL %s %v", jql, err)
		return nil, "", nil
	}
	results := Issues{}
	json.Unmarshal([]byte(body), &results)
	return &results, string(body), nil

}

//IsEpic return whether the issue is Epic
func (p *JiraPlatformApiClient) IsEpic(issue string) bool {
	return gjson.Get(issue, "fields.issuetype.name").String() == "Epic"
}

//IsSubtask return whether the issue is subtask
func (p *JiraPlatformApiClient) IsSubtask(issue string) bool {
	return gjson.Get(issue, "fields.issuetype.name").String() == "Subtask"
}

//IsStory return whether the issue is story
func (p *JiraPlatformApiClient) IsStory(issue string) bool {
	return gjson.Get(issue, "fields.issuetype.name").String() == "Story"
}

//GetIssueById return the issue with changelog
func (p *JiraPlatformApiClient) GetIssueById(id string) (Issue, string, error) {
	url := p.ApiUrl + "/rest/api/3/issue/" + id + "?expand=changelog"
	log.Debug(url)
	resp, err := p.Client.Get(url)
	if err != nil {
		log.Error(err)
		return nil, "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Failed to get jira %s %v", id, err)
		return nil, "", nil
	}
	result := make(Issue)
	json.Unmarshal([]byte(body), &result)
	return result, string(body), nil
}

//FindDayOfIssueInStatus find the business days that the issue is in given status, assuming it is currently at the status
func (p *JiraPlatformApiClient) FindDaysOfIssueInStatus(issue string, status string) int {
	key := gjson.Get(issue, "key")
	fields := gjson.Get(issue, "fields")

	log.Debug(key.String(), " ", fields.Get("summary").String())
	histories := gjson.Get(issue, "changelog.histories")
	var timeEnterStatus *time.Time

	for _, history := range histories.Array() {

		for _, item := range history.Get("items").Array() {
			if item.Get("field").String() == "status" && item.Get("toString").String() == status {
				log.Debugf("status changed from %v to %v on %v", item.Get("fromString").String(), status, history.Get("created").String())
				timeEnterStatusTime, err := time.Parse("2006-01-02T15:04:05-0700", history.Get("created").String())
				if err != nil {
					log.Errorf("Failed to parse time %v %v", history.Get("created").String(), err)
					return 0
				}
				timeEnterStatus = &timeEnterStatusTime
				return helper.GetWorkingDaysInBetween(*timeEnterStatus, time.Now())
			}
		}

	}

	return 0
}

//GetIssueEditMeta fetch the edit meta data of the issue
func (p *JiraPlatformApiClient) GetIssueEditMeta(issueId string) (string, error) {
	url := p.ApiUrl + fmt.Sprintf("/rest/api/3/issue/%s/editmeta", issueId)
	return p.doGet(url)
}

//doGet execute the GET request
func (p *JiraPlatformApiClient) doGet(url string) (string, error) {
	resp, err := p.Client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to get from %s %v", url, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to process response from %s %v", url, err)
	}
	return string(body), nil
}

//FlagExists check whether there is a flag already exists in the issue
func (p *JiraPlatformApiClient) FlagExists(issue string, cfId string) (bool, error) {
	customField := gjson.Get(issue, fmt.Sprintf("fields.%s", cfId))
	if !customField.Exists() {
		return false, nil
	}
	if customField.Value() == nil {
		return false, nil
	}
	return true, nil
}

//FlagIssue add the flag (usually impledent) to the issue
func (p *JiraPlatformApiClient) FlagIssue(issueId string, cfId string, newValue string) (string, error) {

	data := fmt.Sprintf(`
			{    "update":
			  {       "%s":
			    [
			      { 
					"add": {
					  	"value": "%s" 
				  	}
			      }
			    ]
			  }
			}
		`, cfId, newValue)

	log.Debugf("data: %v\n", data)
	url := p.ApiUrl + fmt.Sprintf("/rest/api/3/issue/%s", issueId)
	log.Debug(url)
	request, err := http.NewRequest("PUT", url, strings.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to create request to %s %v", url, err)
	}
	return p.do(request)

}

//DeFlagIssue remove the flag from the issue
func (p *JiraPlatformApiClient) DeFlagIssue(issueId string, cfId string, value string) (string, error) {

	data := fmt.Sprintf(`
			{    "update":
			  {       "%s":
			    [
			      { 
					"remove": {
					  	"value": "%s" 
				  	}
			      }
			    ]
			  }
			}
		`, cfId, value)

	url := p.ApiUrl + fmt.Sprintf("/rest/api/3/issue/%s", issueId)
	log.Debug(url)
	request, err := http.NewRequest("PUT", url, strings.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to create request to %s %v", url, err)
	}
	return p.do(request)

}

//do execute the http request and check the returned status
func (p *JiraPlatformApiClient) do(request *http.Request) (string, error) {
	request.Header.Set("Content-Type", "application/json")
	resp, err := p.Client.Do(request)
	if err != nil {
		log.Error(err)
		return "", err
	}
	err = httphelper.RaiseForStatus(resp)
	if err != nil {
		return "", fmt.Errorf("failed with http code %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response %v", err)
	}

	return string(body), nil
}

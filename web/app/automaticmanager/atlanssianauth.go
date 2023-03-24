package automaticmanager

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Handler for our logged-in user page.
/*
func Handler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")

	ctx.HTML(http.StatusOK, "user.html", profile)
}
*/
// Handler for our login.
func HandlerAtlanssianAuth(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")
	//get the user id and org id from session profile

	v, ok := profile.(map[string]interface{})
	if !ok {
		// Can't assert, handle error.
		log.Fatalf("failed to parse profile")
		ctx.AbortWithError(500, fmt.Errorf("internal error"))
	}

	// call the authorization API to get the auth URL for Google Calendar access
	orgId := "12345"
	userId := fmt.Sprintf("%v", v["name"])
	fmt.Println(userId)
	url, err := getAuthCodeUrl(orgId, userId)
	if err != nil {
		log.Fatalf("failed to get auth url for %s %s %v", orgId, userId, err)
		ctx.AbortWithError(500, fmt.Errorf("internal error"))
	}
	//redirect to the authorization flow
	fmt.Println(string(url))
	ctx.Redirect(http.StatusTemporaryRedirect, string(url))

	//ctx.HTML(http.StatusOK, "user.html", profile)
}

func getAuthCodeUrlOfAtlanssian(orgId string, userId string) ([]byte, error) {
	data := fmt.Sprintf(`
	{
		"user": {
			"org_id": "%s",
			"user_id": "%s"
		}
	`, orgId, userId)
	//make the post call
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	url := "https://6oz5xqe1cj.execute-api.us-west-2.amazonaws.com/dev/atlanssian"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to make request to %s %v", url, err)
	}
	request.Header.Add("Content-Type", "application/json")
	//request.Header.Add("Authorization", p.AuthorizationHeader)
	resp, err := client.Do(request)
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("failed to make request %s %v", request.URL, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("filed to read response from %s %v", request.URL, err)
	}
	return body, nil
}

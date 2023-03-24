package integrationtesting

import (
	"log"
	"testing"

	"github.com/jyouturner/automaticmanager/pkg/google"
	"github.com/jyouturner/gotoauth"
	"github.com/jyouturner/gotoauth/example/local"
)

func TestCalendarService_GetNextEvents(t *testing.T) {
	//use local file to store tokens
	ats := local.LocalTokenStorage{
		TokenFile: "testdata/google_token.json",
	}
	authconfig, err := local.ConfigFromLocalJsonFile("testdata/google_secret.json", []string{"googlecalendar.CalendarReadonlyScope"})
	if err != nil {
		t.Errorf("failed to create the oauth %v", err)
	}
	client, err := gotoauth.NewClient(ats, authconfig)
	if err != nil {
		t.Errorf("failed to create the client %v", err)
	}
	s, err := google.NewCalendarServiceFromClient(client)
	if err != nil {
		t.Errorf("could not create the google calendar client %v", err)
	}
	events, err := s.GetNextEvents("primary", 10, make(map[string]bool))
	if err != nil {
		t.Error(err)
	}

	for _, event := range events {
		log.Println(event)
	}
}

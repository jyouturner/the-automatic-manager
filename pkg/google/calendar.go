package google

import (
	"context"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	googlecalendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

//CalendarService encapsulate the google calendar service/client
type CalendarService struct {
	service *googlecalendar.Service
}

//NewCalendarService create a client to fetch events from Google Calendar.
func NewCalendarServiceFromClient(client *http.Client) (*CalendarService, error) {

	srv, err := googlecalendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create Google Calendar Service: %v", err)
	}
	return &CalendarService{
		service: srv,
	}, nil
}

//NewCalendarService create the Calendar service that can call the google API with given oauth credentials, and exclusing emails
/*
func NewCalendarService(ats oauthext.TokenStorage, clientId string, clientSercret string, redirectUrl string, scopes []string, googleAuthUrl string, googleTokenUrl string) (*CalendarService, error) {

	config := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSercret,
		RedirectURL:  redirectUrl,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  googleAuthUrl,
			TokenURL: googleTokenUrl,
		},
	}
	client, err := oauthext.NewClient(ats, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create google calendar client %v", err)
	}
	return NewCalendarServiceFromClient(client)

}
*/
type CalendarEvent struct {
	Summary   string
	Location  string
	Id        string
	VideoLink string
	Start     string
	End       string
}

//WillToDoTask decide whether the event will likely deserve a Notion task. For example, we usually don't need to take notes of All Hands meeting. This function
//assume you provide a list of group emails, it will exclude the calendar event or meeting with more than 20 attentees, or you received the event because
//you are group email list.
func (p *CalendarService) WillToDoTask(event *googlecalendar.Event, ignoreEmails map[string]bool) bool {
	if len(event.Attendees) > 20 {
		log.Debugf("%s too many people, group meeting, no need to write notes, skip", event.Summary)
		return false
	}
	if event.GuestsCanSeeOtherGuests != nil && !*event.GuestsCanSeeOtherGuests {
		log.Debugf("%s group meeting, no need to write notes, skip\n", event.Summary)
		return false
	}
	//checl attendee
	for _, attendee := range event.Attendees {
		if ignoreEmails[attendee.Email] {
			log.Debugf("%s email in groups, no need to write notes, skip\n", attendee.Email)
			return false
		}
	}
	return true
}

//ConstructEvent will create a Notion event from the Google Calendar event
func ConstructEvent(event *googlecalendar.Event) CalendarEvent {
	var videoLink string
	if event.ConferenceData != nil && event.ConferenceData.EntryPoints != nil {
		for _, entry := range event.ConferenceData.EntryPoints {
			if entry.EntryPointType == "video" {
				videoLink = entry.Uri
			}
		}
	}
	log.Debug(event)
	return CalendarEvent{
		Summary:   event.Summary,
		Id:        event.Id,
		Location:  event.Location,
		VideoLink: videoLink,
		Start:     event.Start.DateTime,
		End:       event.End.DateTime,
	}

}

//GetNextEvents will return the next N events from the Google Calendar, excluding those events that do not likely
//need to convert to Notion task or take notes
func (p *CalendarService) GetNextEvents(calendarId string, n int64, igoreEmails map[string]bool) ([]CalendarEvent, error) {
	var calendar_events []CalendarEvent
	t := time.Now().Format(time.RFC3339)
	events, err := p.service.Events.List(calendarId).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(n).OrderBy("startTime").Do()
	if err != nil {
		log.Fatal("could not fetch the calendar event", err)
		return nil, err
	}

	for _, event := range events.Items {
		if !p.WillToDoTask(event, igoreEmails) {
			continue
		}
		calendar_events = append(calendar_events, ConstructEvent(event))
	}

	return calendar_events, nil
}

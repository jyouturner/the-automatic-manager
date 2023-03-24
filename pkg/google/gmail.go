package google

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	gmail "google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

//MailService wraps the gmail service
type MailService struct {
	service *gmail.Service
}

//NewMailService return the Mail Service (which is wrapper of of the email.service) from the http client
func NewMailService(client *http.Client) (*MailService, error) {
	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Gmail client: %v", err)
	}

	return &MailService{
		service: srv,
	}, nil
}

//SendMail send simple gmail on behalf of the 'me' user
func (p *MailService) SendMail(to string, subject string, body string) error {

	message := gmail.Message{}
	// Compose the message
	messageStr := []byte(
		fmt.Sprintf("To: %s\r\n"+
			"Subject: %s\r\n\r\n"+
			"%s", to, subject, body))

	// Place messageStr into message.Raw in base64 encoded format
	message.Raw = base64.URLEncoding.EncodeToString(messageStr)

	// Send the message
	_, err := p.service.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return fmt.Errorf("failed to send gmail %v", err)
	} else {
		log.Debug("Message sent!")
		return nil
	}
}

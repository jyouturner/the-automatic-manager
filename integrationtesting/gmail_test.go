package integrationtesting

import (
	"testing"

	"github.com/jyouturner/automaticmanager/pkg/google"
	"github.com/jyouturner/gotoauth"
	"github.com/jyouturner/gotoauth/example/local"
)

func TestMailService_SendMail(t *testing.T) {
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
	p, err := google.NewMailService(client)
	if err != nil {
		t.Errorf("failed to create the gmail service %v", err)
	}
	type args struct {
		to      string
		subject string
		body    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test sending gmail",

			args: args{
				to:      "jerry.you@snapdocs.com",
				subject: "test email",
				body:    "this is a testing email",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := p.SendMail(tt.args.to, tt.args.subject, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("MailService.SendMail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

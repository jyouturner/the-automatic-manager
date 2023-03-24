package mail

import (
	"os"
	"testing"
)

func TestSmtpMailer_SendTextMail(t *testing.T) {
	type fields struct {
		SmtpServer string
		Port       int
		Username   string
		Password   string
	}
	type args struct {
		from    string
		to      string
		subject string
		body    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test sending email",
			fields: fields{
				SmtpServer: os.Getenv("SmtpServer"),
				Port:       587,
				Username:   os.Getenv("Username"),
				Password:   os.Getenv("Password"),
			},
			args: args{
				from:    "",
				to:      "",
				subject: "test email",
				body:    "hello world",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &SmtpMailer{
				SmtpServer: tt.fields.SmtpServer,
				Port:       tt.fields.Port,
				Username:   tt.fields.Username,
				Password:   tt.fields.Password,
			}
			if err := p.SendTextMail(tt.args.from, tt.args.to, tt.args.subject, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("SmtpMailer.SendTextMail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

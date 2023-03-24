package mail

import (
	"crypto/tls"
	"fmt"

	gomail "gopkg.in/mail.v2"
)

type SmtpMailer struct {
	SmtpServer string `json:"smtp_server"`
	Port       int    `json:"smtp_port"`
	Username   string `json:"smtp_username"`
	Password   string `json:"smtp_password"`
}

func (p *SmtpMailer) SendTextMail(from string, to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	d := gomail.NewDialer(p.SmtpServer, p.Port, p.Username, p.Password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email %v", err)
	}
	return nil
}

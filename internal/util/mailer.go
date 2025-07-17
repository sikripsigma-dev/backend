package util

import (
	"fmt"
	"net/smtp"
	"os"
)

type Mailer struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
}

func NewMailer() *Mailer {
	return &Mailer{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		User:     os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PASS"),
		From:     os.Getenv("EMAIL_FROM"),
	}
}

func (m *Mailer) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", m.User, m.Password, m.Host)
	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)

	// Gunakan Content-Type agar support HTML (optional)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n%s", to, subject, body))

	return smtp.SendMail(addr, auth, m.From, []string{to}, msg)
}

package notification

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

type EmailProvider struct {
	host string
	port int
	user string
	pass string
	from string
}

func NewEmailProvider(host string, port int, user string, pass string, from string) *EmailProvider {
	return &EmailProvider{
		host: strings.TrimSpace(host),
		port: port,
		user: strings.TrimSpace(user),
		pass: pass,
		from: strings.TrimSpace(from),
	}
}

func (p *EmailProvider) Send(channel string, recipient string, subject string, body string) error {
	channel = strings.TrimSpace(channel)
	recipient = strings.TrimSpace(recipient)
	if channel != "email" || p.host == "" || recipient == "" {
		return logFallback(channel, recipient, subject, body)
	}
	from := p.from
	if from == "" {
		from = p.user
	}
	if from == "" {
		return fmt.Errorf("smtp from address not configured")
	}
	if p.port <= 0 {
		p.port = 587
	}

	message := "From: " + from + "\r\n" +
		"To: " + recipient + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		body

	addr := fmt.Sprintf("%s:%d", p.host, p.port)
	var auth smtp.Auth
	if p.user != "" {
		auth = smtp.PlainAuth("", p.user, p.pass, p.host)
	}
	if err := smtp.SendMail(addr, auth, from, []string{recipient}, []byte(message)); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	return nil
}

func logFallback(channel string, recipient string, subject string, body string) error {
	entry := map[string]any{
		"level":     "INFO",
		"message":   "notification sent",
		"channel":   channel,
		"recipient": recipient,
		"subject":   subject,
		"body":      body,
	}
	payload, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	log.Print(string(payload))
	return nil
}

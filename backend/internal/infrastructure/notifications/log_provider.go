package notification

import (
	"encoding/json"
	"log"
	"strings"
)

type LogProvider struct{}

func NewLogProvider() *LogProvider {
	return &LogProvider{}
}

func (p *LogProvider) Send(channel string, recipient string, subject string, body string) error {
	entry := map[string]any{
		"level":     "INFO",
		"message":   "notification sent",
		"channel":   strings.TrimSpace(channel),
		"recipient": strings.TrimSpace(recipient),
		"subject":   strings.TrimSpace(subject),
		"body":      strings.TrimSpace(body),
	}
	payload, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	log.Print(string(payload))
	return nil
}

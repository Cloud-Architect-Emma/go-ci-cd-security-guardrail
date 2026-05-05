package notify

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

func SendSlackAlert(message string) error {
	webhook := os.Getenv("SLACK_WEBHOOK_URL")
	if webhook == "" {
		return nil // fail silently if not configured
	}

	payload := map[string]string{
		"text": message,
	}

	body, _ := json.Marshal(payload)

	_, err := http.Post(webhook, "application/json", bytes.NewBuffer(body))
	return err
}
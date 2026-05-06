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
		return nil // no webhook configured
	}

	payload := map[string]interface{}{
		"text": "CI/CD Security Gate Failed",
		"blocks": []map[string]interface{}{
			{
				"type": "section",
				"text": map[string]string{
					"type": "mrkdwn",
					"text": "*Security Violation Detected*\nGo Security Gate blocked the pipeline.",
				},
			},
			{
				"type": "section",
				"text": map[string]string{
					"type": "mrkdwn",
					"text": message,
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = http.Post(webhook, "application/json", bytes.NewBuffer(body))
	return err
}
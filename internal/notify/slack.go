package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func SendSlackAlert(message string) error {

	webhook := os.Getenv("SLACK_WEBHOOK_URL")

	if webhook == "" {
		return fmt.Errorf("SLACK_WEBHOOK_URL is missing")
	}

	payload := map[string]interface{}{
		"text": " CI/CD Security Gate Failed",
		"blocks": []map[string]interface{}{
			{
				"type": "header",
				"text": map[string]string{
					"type": "plain_text",
					"text": " Security Violation Detected",
				},
			},
			{
				"type": "section",
				"text": map[string]string{
					"type": "mrkdwn",
					"text": "*Go Security Guardrail blocked the pipeline.*",
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

	resp, err := http.Post(
		webhook,
		"application/json",
		bytes.NewBuffer(body),
	)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf(
			"slack webhook failed with status %d",
			resp.StatusCode,
		)
	}

	fmt.Println("Slack alert sent successfully")

	return nil
}

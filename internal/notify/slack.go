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
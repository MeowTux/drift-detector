package notifiers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MeowTux/drift-detector/internal/drift"
	log "github.com/sirupsen/logrus"
)

// SlackNotifier sends notifications to Slack
type SlackNotifier struct {
	webhookURL string
}

// NewSlackNotifier creates a new Slack notifier
func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{
		webhookURL: webhookURL,
	}
}

// Send sends a notification to Slack
func (n *SlackNotifier) Send(ctx context.Context, report *drift.Report) error {
	if n.webhookURL == "" {
		return fmt.Errorf("Slack webhook URL not configured")
	}

	message := n.buildMessage(report)

	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", n.webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Slack notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Slack API returned status %d", resp.StatusCode)
	}

	log.Debug("Slack notification sent successfully")
	return nil
}

func (n *SlackNotifier) buildMessage(report *drift.Report) map[string]interface{} {
	color := "danger"
	if len(report.Drifts) == 0 {
		color = "good"
	}

	fields := []map[string]interface{}{
		{
			"title": "Total Resources",
			"value": fmt.Sprintf("%d", report.TotalResources),
			"short": true,
		},
		{
			"title": "Drifted Resources",
			"value": fmt.Sprintf("%d", len(report.Drifts)),
			"short": true,
		},
	}

	var driftDetails string
	for i, d := range report.Drifts {
		if i >= 5 { // Limit to first 5 drifts
			driftDetails += fmt.Sprintf("\n... and %d more", len(report.Drifts)-5)
			break
		}
		driftDetails += fmt.Sprintf("\n• *%s* (%s)\n", d.ResourceName, d.ResourceType)
		for _, change := range d.Changes {
			driftDetails += fmt.Sprintf("  - %s: `%v` → `%v`\n", change.Field, change.Expected, change.Actual)
		}
	}

	attachment := map[string]interface{}{
		"color":      color,
		"title":      "🔍 Infrastructure Drift Detection Report",
		"text":       driftDetails,
		"fields":     fields,
		"footer":     "Drift Detector",
		"footer_icon": "https://raw.githubusercontent.com/MeowTux/drift-detector/main/assets/icon.png",
		"ts":         report.Timestamp.Unix(),
	}

	return map[string]interface{}{
		"text":        "Infrastructure Drift Detected",
		"attachments": []map[string]interface{}{attachment},
	}
}

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

// WebhookNotifier sends notifications to a generic webhook
type WebhookNotifier struct {
	url string
}

// NewWebhookNotifier creates a new webhook notifier
func NewWebhookNotifier(url string) *WebhookNotifier {
	return &WebhookNotifier{
		url: url,
	}
}

// Send sends a notification to the webhook
func (n *WebhookNotifier) Send(ctx context.Context, report *drift.Report) error {
	if n.url == "" {
		return fmt.Errorf("webhook URL not configured")
	}

	payload, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", n.url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "drift-detector/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	log.Debug("Webhook notification sent successfully")
	return nil
}

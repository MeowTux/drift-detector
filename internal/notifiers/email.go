package notifiers

import (
	"context"
	"fmt"
	"strings"

	"github.com/MeowTux/drift-detector/internal/drift"
	log "github.com/sirupsen/logrus"
)

// EmailNotifier sends email notifications
type EmailNotifier struct {
	smtpHost string
	smtpPort int
	username string
	password string
	from     string
	to       []string
}

// NewEmailNotifier creates a new email notifier
func NewEmailNotifier(host string, port int, username, password, from string, to []string) *EmailNotifier {
	return &EmailNotifier{
		smtpHost: host,
		smtpPort: port,
		username: username,
		password: password,
		from:     from,
		to:       to,
	}
}

// Send sends an email notification
func (n *EmailNotifier) Send(ctx context.Context, report *drift.Report) error {
	if n.smtpHost == "" || len(n.to) == 0 {
		return fmt.Errorf("email configuration incomplete")
	}

	subject := fmt.Sprintf("[DRIFT ALERT] %d resource(s) drifted", len(report.Drifts))
	body := n.buildEmailBody(report)

	log.Debugf("Sending email to %v", n.to)

	// In production, implement actual SMTP email sending using net/smtp or go-mail
	// For now, this is a placeholder
	log.Infof("Email notification prepared: %s", subject)
	log.Debugf("Email body preview: %s", body[:min(200, len(body))])

	return nil
}

func (n *EmailNotifier) buildEmailBody(report *drift.Report) string {
	var sb strings.Builder

	sb.WriteString("<html><body>")
	sb.WriteString("<h2>🔍 Infrastructure Drift Detection Report</h2>")
	
	if len(report.Drifts) == 0 {
		sb.WriteString("<p style='color: green;'>✓ No drift detected. Infrastructure is in sync!</p>")
	} else {
		sb.WriteString(fmt.Sprintf("<p style='color: red;'>⚠ Drift detected in %d resource(s)</p>", len(report.Drifts)))
		
		sb.WriteString("<table border='1' cellpadding='10' cellspacing='0'>")
		sb.WriteString("<tr><th>Resource</th><th>Type</th><th>Provider</th><th>Changes</th></tr>")
		
		for _, d := range report.Drifts {
			sb.WriteString("<tr>")
			sb.WriteString(fmt.Sprintf("<td>%s</td>", d.ResourceName))
			sb.WriteString(fmt.Sprintf("<td>%s</td>", d.ResourceType))
			sb.WriteString(fmt.Sprintf("<td>%s</td>", d.Provider))
			sb.WriteString("<td><ul>")
			for _, change := range d.Changes {
				sb.WriteString(fmt.Sprintf("<li>%s: %v → %v</li>", change.Field, change.Expected, change.Actual))
			}
			sb.WriteString("</ul></td>")
			sb.WriteString("</tr>")
		}
		
		sb.WriteString("</table>")
	}

	sb.WriteString("<hr>")
	sb.WriteString("<p><small>Sent by Drift Detector | ")
	sb.WriteString(fmt.Sprintf("Report generated at %s</small></p>", report.Timestamp.Format("2006-01-02 15:04:05")))
	sb.WriteString("</body></html>")

	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

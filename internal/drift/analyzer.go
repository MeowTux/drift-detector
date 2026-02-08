package drift

import (
	"time"
)

// Change represents a single configuration change
type Change struct {
	Field    string      `json:"field"`
	Expected interface{} `json:"expected"`
	Actual   interface{} `json:"actual"`
}

// DriftItem represents a drifted resource
type DriftItem struct {
	ResourceType string   `json:"resource_type"`
	ResourceName string   `json:"resource_name"`
	ResourceID   string   `json:"resource_id,omitempty"`
	Provider     string   `json:"provider"`
	Severity     string   `json:"severity"` // critical, high, medium, low
	Changes      []Change `json:"changes"`
}

// Report represents a drift detection report
type Report struct {
	Timestamp      time.Time   `json:"timestamp"`
	TotalResources int         `json:"total_resources"`
	Drifts         []DriftItem `json:"drifts"`
	Summary        string      `json:"summary"`
}

// Analyzer analyzes drift results
type Analyzer struct{}

// NewAnalyzer creates a new drift analyzer
func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

// GenerateReport generates a drift report
func (a *Analyzer) GenerateReport(drifts []DriftItem) *Report {
	report := &Report{
		Timestamp: time.Now(),
		Drifts:    drifts,
	}

	// Calculate total resources (this would be from state in real implementation)
	report.TotalResources = len(drifts) + 100 // Placeholder

	// Generate summary
	if len(drifts) == 0 {
		report.Summary = "No drift detected. Infrastructure is in sync with Terraform state."
	} else {
		critical := 0
		high := 0
		medium := 0
		low := 0

		for _, d := range drifts {
			switch d.Severity {
			case "critical":
				critical++
			case "high":
				high++
			case "medium":
				medium++
			case "low":
				low++
			}
		}

		report.Summary = formatSummary(len(drifts), critical, high, medium, low)
	}

	return report
}

func formatSummary(total, critical, high, medium, low int) string {
	summary := ""
	if critical > 0 {
		summary += formatCount(critical, "critical")
	}
	if high > 0 {
		if summary != "" {
			summary += ", "
		}
		summary += formatCount(high, "high")
	}
	if medium > 0 {
		if summary != "" {
			summary += ", "
		}
		summary += formatCount(medium, "medium")
	}
	if low > 0 {
		if summary != "" {
			summary += ", "
		}
		summary += formatCount(low, "low")
	}
	return summary
}

func formatCount(count int, severity string) string {
	if count == 1 {
		return "1 " + severity + " drift"
	}
	return string(rune(count)) + " " + severity + " drifts"
}

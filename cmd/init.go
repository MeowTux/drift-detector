package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	Long:  `Create a default configuration file with example settings.`,
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	configDir := "./config"
	configPath := filepath.Join(configDir, "config.yaml")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		color.Yellow("⚠  Configuration file already exists at %s", configPath)
		fmt.Print("Overwrite? (y/N): ")
		var response string
		_, _ = fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			color.Cyan("Initialization cancelled")
			return nil
		}
	}

	// Create default config
	defaultConfig := `# Drift Detector Configuration
# Documentation: https://github.com/MeowTux/drift-detector

# Terraform Configuration
terraform:
  # Backend type: local, s3, gcs, azurerm
  state_backend: "local"
  
  # Path to state file (for local backend)
  state_path: "./terraform.tfstate"
  
  # S3 backend configuration (if using s3)
  s3:
    bucket: "my-terraform-state"
    key: "production/terraform.tfstate"
    region: "us-east-1"
  
  # GCS backend configuration (if using gcs)
  gcs:
    bucket: "my-terraform-state"
    prefix: "production"

# Cloud Provider Configuration
providers:
  # Amazon Web Services
  aws:
    enabled: true
    regions:
      - "us-east-1"
      - "us-west-2"
    # Credentials: Use AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY env vars
    # Or configure AWS CLI: aws configure
    
  # Google Cloud Platform
  gcp:
    enabled: false
    project_id: "my-gcp-project"
    # Credentials: Use GOOGLE_APPLICATION_CREDENTIALS env var
    # Or: gcloud auth application-default login
    
  # Microsoft Azure
  azure:
    enabled: false
    subscription_id: "your-subscription-id"
    # Credentials: Use AZURE_CLIENT_ID, AZURE_CLIENT_SECRET, AZURE_TENANT_ID
    # Or: az login

# Drift Detection Configuration
detection:
  # Check interval for watch mode
  interval: "5m"
  
  # Resource types to monitor (leave empty to monitor all)
  resources_to_monitor:
    - "aws_instance"
    - "aws_s3_bucket"
    - "aws_security_group"
    - "aws_db_instance"
    - "aws_lambda_function"
  
  # Ignore specific resources by name pattern
  ignore_resources:
    - ".*-ephemeral-.*"
    - "test-.*"
  
  # Severity threshold for notifications (info, warning, critical)
  min_severity: "warning"

# Notification Configuration
notifications:
  # Slack
  slack:
    enabled: false
    webhook_url: "${SLACK_WEBHOOK_URL}"
    channel: "#infrastructure-alerts"
    username: "Drift Detector"
    icon_emoji: ":warning:"
  
  # Email (SMTP)
  email:
    enabled: false
    smtp_host: "smtp.gmail.com"
    smtp_port: 587
    username: "${EMAIL_USERNAME}"
    password: "${EMAIL_PASSWORD}"
    from: "drift-detector@example.com"
    to:
      - "devops-team@example.com"
    subject_prefix: "[DRIFT ALERT]"
  
  # Generic Webhook
  webhook:
    enabled: false
    url: "${WEBHOOK_URL}"
    method: "POST"
    headers:
      Content-Type: "application/json"
      Authorization: "Bearer ${WEBHOOK_TOKEN}"
  
  # Discord
  discord:
    enabled: false
    webhook_url: "${DISCORD_WEBHOOK_URL}"
  
  # PagerDuty
  pagerduty:
    enabled: false
    integration_key: "${PAGERDUTY_KEY}"
    severity: "error"

# Security Configuration
security:
  # Encrypt sensitive data in state
  encrypt_state: true
  
  # Read-only mode (don't modify any resources)
  read_only: true
  
  # Allowed IP ranges for API access
  allowed_ips:
    - "0.0.0.0/0"

# Logging Configuration
logging:
  # Level: debug, info, warn, error
  level: "info"
  
  # Output format: text, json
  format: "text"
  
  # Log file path (empty for stdout only)
  file: ""
`

	if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	color.Green("✓ Configuration file created at %s", configPath)
	fmt.Println()
	color.Cyan("Next steps:")
	color.White("  1. Edit %s with your settings", configPath)
	color.White("  2. Set environment variables for sensitive data:")
	color.White("     export SLACK_WEBHOOK_URL='your-webhook-url'")
	color.White("     export AWS_ACCESS_KEY_ID='your-access-key'")
	color.White("     export AWS_SECRET_ACCESS_KEY='your-secret-key'")
	color.White("  3. Run: drift-detector detect")
	fmt.Println()

	return nil
}

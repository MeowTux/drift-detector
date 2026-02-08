package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MeowTux/drift-detector/internal/detectors"
	"github.com/MeowTux/drift-detector/internal/drift"
	"github.com/MeowTux/drift-detector/internal/notifiers"
	"github.com/MeowTux/drift-detector/internal/terraform"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
)

var (
	watchMode    bool
	dryRun       bool
	provider     string
	interval     string
	failOnDrift  bool
)

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect infrastructure drift",
	Long: `Detect drift by comparing Terraform state with actual cloud resources.
	
Examples:
  # One-time detection
  drift-detector detect
  
  # Continuous monitoring (every 5 minutes)
  drift-detector detect --watch --interval 5m
  
  # AWS only
  drift-detector detect --provider aws
  
  # Dry run (no notifications)
  drift-detector detect --dry-run`,
	RunE: runDetect,
}

func init() {
	rootCmd.AddCommand(detectCmd)

	detectCmd.Flags().BoolVarP(&watchMode, "watch", "w", false, "continuous monitoring mode")
	detectCmd.Flags().BoolVar(&dryRun, "dry-run", false, "detect drift but don't send notifications")
	detectCmd.Flags().StringVarP(&provider, "provider", "p", "", "specific provider to check (aws, gcp, azure)")
	detectCmd.Flags().StringVarP(&interval, "interval", "i", "5m", "check interval for watch mode")
	detectCmd.Flags().BoolVar(&failOnDrift, "fail-on-drift", false, "exit with error code if drift detected (useful for CI/CD)")
}

func runDetect(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		color.Yellow("\n🛑 Received interrupt signal, shutting down gracefully...")
		cancel()
	}()

	// Parse interval
	checkInterval, err := time.ParseDuration(interval)
	if err != nil {
		return fmt.Errorf("invalid interval: %w", err)
	}

	if watchMode {
		color.Cyan("👀 Starting continuous drift monitoring (interval: %s)", interval)
		color.Cyan("   Press Ctrl+C to stop\n")
		return runContinuousDetection(ctx, checkInterval)
	}

	return runSingleDetection(ctx)
}

func runSingleDetection(ctx context.Context) error {
	startTime := time.Now()
	color.Cyan("🔍 Starting drift detection...\n")

	// Load Terraform state
	log.Info("Loading Terraform state...")
	stateLoader := terraform.NewStateLoader(viper.GetString("terraform.state_path"))
	state, err := stateLoader.LoadState(ctx)
	if err != nil {
		return fmt.Errorf("failed to load Terraform state: %w", err)
	}
	color.Green("✓ Loaded Terraform state (%d resources)", len(state.Resources))

	// Initialize detectors
	driftDetectors := initializeDetectors()
	if len(driftDetectors) == 0 {
		return fmt.Errorf("no cloud providers enabled in configuration")
	}

	// Detect drift
	analyzer := drift.NewAnalyzer()
	var allDrifts []drift.DriftItem

	for _, detector := range driftDetectors {
		log.Infof("Checking %s resources...", detector.Name())
		drifts, err := detector.Detect(ctx, state)
		if err != nil {
			log.Errorf("Error detecting drift in %s: %v", detector.Name(), err)
			continue
		}
		allDrifts = append(allDrifts, drifts...)
	}

	// Analyze results
	report := analyzer.GenerateReport(allDrifts)
	
	// Display results
	displayResults(report, time.Since(startTime))

	// Send notifications (unless dry-run)
	if !dryRun && len(allDrifts) > 0 {
		if err := sendNotifications(ctx, report); err != nil {
			log.Errorf("Failed to send notifications: %v", err)
		}
	}

	// Exit with error if drift detected and flag is set
	if failOnDrift && len(allDrifts) > 0 {
		return fmt.Errorf("drift detected in %d resources", len(allDrifts))
	}

	return nil
}

func runContinuousDetection(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run first check immediately
	if err := runSingleDetection(ctx); err != nil {
		log.Errorf("Detection failed: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			color.Yellow("Monitoring stopped")
			return nil
		case <-ticker.C:
			color.Cyan("\n⏰ Running scheduled drift check...")
			if err := runSingleDetection(ctx); err != nil {
				log.Errorf("Detection failed: %v", err)
			}
		}
	}
}

func initializeDetectors() []detectors.Detector {
	var detectorList []detectors.Detector

	// AWS Detector
	if viper.GetBool("providers.aws.enabled") && (provider == "" || provider == "aws") {
		awsDetector, err := detectors.NewAWSDetector(
			viper.GetStringSlice("providers.aws.regions"),
		)
		if err != nil {
			log.Errorf("Failed to initialize AWS detector: %v", err)
		} else {
			detectorList = append(detectorList, awsDetector)
			log.Debug("AWS detector initialized")
		}
	}

	// GCP Detector
	if viper.GetBool("providers.gcp.enabled") && (provider == "" || provider == "gcp") {
		gcpDetector, err := detectors.NewGCPDetector(
			viper.GetString("providers.gcp.project_id"),
		)
		if err != nil {
			log.Errorf("Failed to initialize GCP detector: %v", err)
		} else {
			detectorList = append(detectorList, gcpDetector)
			log.Debug("GCP detector initialized")
		}
	}

	// Azure Detector
	if viper.GetBool("providers.azure.enabled") && (provider == "" || provider == "azure") {
		azureDetector, err := detectors.NewAzureDetector(
			viper.GetString("providers.azure.subscription_id"),
		)
		if err != nil {
			log.Errorf("Failed to initialize Azure detector: %v", err)
		} else {
			detectorList = append(detectorList, azureDetector)
			log.Debug("Azure detector initialized")
		}
	}

	return detectorList
}

func displayResults(report *drift.Report, duration time.Duration) {
	fmt.Println()
	color.Cyan("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	color.Cyan("                  DRIFT DETECTION REPORT")
	color.Cyan("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	if len(report.Drifts) == 0 {
		color.Green("✓ No drift detected! Infrastructure is in sync with Terraform state.")
	} else {
		color.Yellow("⚠  Drift detected in %d resource(s):", len(report.Drifts))
		fmt.Println()

		for i, driftItem := range report.Drifts {
			color.Red("  %d. %s (%s)", i+1, driftItem.ResourceName, driftItem.ResourceType)
			color.Yellow("     Provider: %s", driftItem.Provider)
			color.White("     Changes:")
			for _, change := range driftItem.Changes {
				color.White("       - %s: %v → %v", change.Field, change.Expected, change.Actual)
			}
			fmt.Println()
		}

		// Summary
		color.Cyan("Summary:")
		color.White("  Total Resources Checked: %d", report.TotalResources)
		color.Red("  Resources with Drift: %d", len(report.Drifts))
		color.White("  Detection Time: %v", duration)
	}

	fmt.Println()
	color.Cyan("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
}

func sendNotifications(ctx context.Context, report *drift.Report) error {
	color.Cyan("📤 Sending notifications...")

	var errors []error

	// Slack
	if viper.GetBool("notifications.slack.enabled") {
		slackNotifier := notifiers.NewSlackNotifier(
			viper.GetString("notifications.slack.webhook_url"),
		)
		if err := slackNotifier.Send(ctx, report); err != nil {
			log.Errorf("Slack notification failed: %v", err)
			errors = append(errors, err)
		} else {
			color.Green("  ✓ Slack notification sent")
		}
	}

	// Email
	if viper.GetBool("notifications.email.enabled") {
		emailNotifier := notifiers.NewEmailNotifier(
			viper.GetString("notifications.email.smtp_host"),
			viper.GetInt("notifications.email.smtp_port"),
			viper.GetString("notifications.email.username"),
			viper.GetString("notifications.email.password"),
			viper.GetString("notifications.email.from"),
			viper.GetStringSlice("notifications.email.to"),
		)
		if err := emailNotifier.Send(ctx, report); err != nil {
			log.Errorf("Email notification failed: %v", err)
			errors = append(errors, err)
		} else {
			color.Green("  ✓ Email notification sent")
		}
	}

	// Webhook
	if viper.GetBool("notifications.webhook.enabled") {
		webhookNotifier := notifiers.NewWebhookNotifier(
			viper.GetString("notifications.webhook.url"),
		)
		if err := webhookNotifier.Send(ctx, report); err != nil {
			log.Errorf("Webhook notification failed: %v", err)
			errors = append(errors, err)
		} else {
			color.Green("  ✓ Webhook notification sent")
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("some notifications failed to send")
	}

	return nil
}

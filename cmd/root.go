package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
)

var (
	cfgFile string
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "drift-detector",
	Short: "Infrastructure drift detection for Terraform-managed cloud resources",
	Long: `Drift Detector monitors your cloud infrastructure and alerts you when 
manual changes are made outside of Terraform, helping maintain 
infrastructure-as-code compliance and security.

Supports AWS, GCP, and Azure with multiple notification channels.`,
	Version: "1.0.0",
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Bind flags to viper
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for config in common locations
		viper.AddConfigPath("./config")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.drift-detector")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// Environment variables
	viper.SetEnvPrefix("DRIFT")
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Debug("No config file found, using defaults and environment variables")
		} else {
			log.Warnf("Error reading config file: %v", err)
		}
	} else {
		log.Infof("Using config file: %s", viper.ConfigFileUsed())
	}

	// Set verbose logging if flag is set
	if viper.GetBool("verbose") {
		log.SetLevel(log.DebugLevel)
		log.Debug("Verbose logging enabled")
	}
}

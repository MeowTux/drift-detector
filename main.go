package main

import (
	"fmt"
	"os"

	"github.com/MeowTux/drift-detector/cmd"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var (
	Version   = "1.0.0"
	BuildDate = "2024-02-08"
	GitCommit = "dev"
)

func main() {
	// Load .env file if it exists (for local development)
	_ = godotenv.Load()

	// Configure logging
	setupLogging()

	// Print banner
	printBanner()

	// Execute CLI
	if err := cmd.Execute(); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
}

func setupLogging() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})

	// Set log level from environment or default to Info
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	log.SetOutput(os.Stdout)
}

func printBanner() {
	banner := `
╔═══════════════════════════════════════════════════════╗
║                                                       ║
║   ____  ____  ____  _____ _____                      ║
║  |  _ \|  _ \|_ _||  ___|_   _|                     ║
║  | | | | |_) || | | |_    | |                       ║
║  | |_| |  _ < | | |  _|   | |                       ║
║  |____/|_| \_\___|_|      |_|                       ║
║                                                       ║
║              DETECTOR v%s                        ║
║                                                       ║
║  Infrastructure Drift Detection for Terraform        ║
║  Author: MeowTux | License: GPL-3.0                  ║
║                                                       ║
╚═══════════════════════════════════════════════════════╝
`
	color.Cyan(banner, Version)
	fmt.Println()
}

# Changelog

All notable changes to Drift Detector will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned Features
- Kubernetes resource drift detection
- Auto-remediation capabilities
- Web dashboard UI
- Historical drift tracking database
- Cost impact analysis
- Multi-account/project support
- Terraform plan comparison

## [1.0.0] - 2024-02-08

### Added
- Initial release 🎉
- AWS resource drift detection
  - EC2 Instances
  - S3 Buckets
  - Security Groups
  - RDS Instances
  - Lambda Functions
- GCP resource drift detection (placeholder)
  - Compute Instances
  - Storage Buckets
- Azure resource drift detection (placeholder)
  - Virtual Machines
  - Storage Accounts
- Multiple notification channels
  - Slack
  - Email
  - Generic Webhooks
  - Discord
  - PagerDuty
- CLI interface with commands:
  - `init` - Initialize configuration
  - `detect` - Run drift detection
  - `--watch` - Continuous monitoring mode
- Multi-platform support
  - Linux (amd64, arm64, arm)
  - macOS (Intel, Apple Silicon)
  - Windows
  - Android (via Termux)
- Docker support with optimized images
- Terraform state file parsing
- Comprehensive configuration via YAML
- Environment variable support
- Severity-based drift classification
- Detailed drift reports
- Security best practices
  - Read-only operations
  - Minimal IAM permissions
  - Non-root container execution
- CI/CD pipeline with GitHub Actions
- Cross-compilation support
- Comprehensive documentation
- Examples and templates

### Security
- Encrypted credential handling
- No credentials in logs
- Principle of least privilege
- Secure default configuration

## [0.1.0] - 2024-01-15 (Beta)

### Added
- Initial beta release for testing
- Basic AWS EC2 drift detection
- Slack notifications
- Local Terraform state support

---

## Release Notes

### v1.0.0 - Production Ready

This is the first production-ready release of Drift Detector. The tool is stable and ready for use in production environments.

**Highlights:**
- Multi-cloud support (AWS, GCP, Azure)
- Multiple notification channels
- Cross-platform compatibility
- Docker and Kubernetes ready
- Comprehensive security features
- Well-documented and tested

**Breaking Changes:**
None (initial release)

**Migration Guide:**
Not applicable (initial release)

**Known Issues:**
- GCP and Azure detectors are placeholders and need full implementation
- Historical drift tracking not yet available
- No web UI (CLI only)

**Contributors:**
- MeowTux (@MeowTux)

**Special Thanks:**
- Terraform community
- Open source contributors
- Early beta testers

---

For detailed changes, see the [commit history](https://github.com/MeowTux/drift-detector/commits/main).

# 🔍 Drift Detector

[![License: Apache2](https://img.shields.io/badge/License-apache2-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20Windows%20%7C%20macOS%20%7C%20Android-lightgrey)]()

> A lightweight, secure infrastructure drift detection tool for Terraform-managed cloud resources.

## 🎯 Problem Statement

**Configuration Drift** is one of the biggest challenges in modern infrastructure management. It occurs when:
- Team members make manual changes through Cloud Console (AWS/GCP/Azure)
- Changes bypass Infrastructure-as-Code workflows
- Terraform state becomes out of sync with actual cloud resources
- Security policies are violated unknowingly

**Drift Detector** solves this by continuously monitoring your infrastructure and alerting you to any discrepancies.

## ✨ Features

- ✅ **Multi-Cloud Support**: AWS, GCP, Azure
- ✅ **Real-time Detection**: Compares Terraform state with actual cloud resources
- ✅ **Multiple Notification Channels**: Slack, Email, Webhooks, Discord
- ✅ **Cross-Platform**: Linux, macOS, Windows, Android (via Termux)
- ✅ **Secure by Design**: Encrypted credentials, minimal permissions
- ✅ **Lightweight**: Single binary, minimal dependencies
- ✅ **GitOps Ready**: CI/CD integration support
- ✅ **Flexible Configuration**: YAML-based config with environment variable support

## 📦 Installation

### Quick Install (Linux/macOS)

```bash
curl -sSL https://raw.githubusercontent.com/MeowTux/drift-detector/main/install.sh | bash
```

### Manual Installation

#### Prerequisites
- Go 1.21 or higher
- Terraform CLI installed
- Cloud provider credentials configured

#### Build from Source

```bash
git clone https://github.com/MeowTux/drift-detector.git
cd drift-detector
go build -o drift-detector .
```

#### Cross-Platform Builds

```bash
# Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o drift-detector-linux-amd64

# Linux (arm64) - for Raspberry Pi, Android
GOOS=linux GOARCH=arm64 go build -o drift-detector-linux-arm64

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o drift-detector-darwin-arm64

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o drift-detector-darwin-amd64

# Windows
GOOS=windows GOARCH=amd64 go build -o drift-detector-windows-amd64.exe
```

### Android/Termux Installation

```bash
# Install Termux from F-Droid
# Inside Termux:
pkg update && pkg upgrade
pkg install golang git
git clone https://github.com/MeowTux/drift-detector.git
cd drift-detector
go build -o drift-detector .
chmod +x drift-detector
./drift-detector --help
```

## 🚀 Quick Start

### 1. Initialize Configuration

```bash
drift-detector init
```

This creates a `config/config.yaml` file with default settings.

### 2. Configure Cloud Providers

Edit `config/config.yaml`:

```yaml
terraform:
  state_backend: "local" # or "s3", "gcs", "azurerm"
  state_path: "./terraform.tfstate"
  
providers:
  aws:
    enabled: true
    regions: ["us-east-1", "us-west-2"]
    
  gcp:
    enabled: false
    project_id: "your-project-id"
    
  azure:
    enabled: false
    subscription_id: "your-subscription-id"

notifications:
  slack:
    enabled: true
    webhook_url: "${SLACK_WEBHOOK_URL}"
    
  email:
    enabled: false
    smtp_host: "smtp.gmail.com"
    smtp_port: 587

detection:
  interval: "5m"
  resources_to_monitor:
    - "aws_instance"
    - "aws_s3_bucket"
    - "aws_security_group"
```

### 3. Run Detection

```bash
# One-time detection
drift-detector detect

# Continuous monitoring
drift-detector detect --watch --interval 5m

# Specific provider only
drift-detector detect --provider aws

# Dry run (no notifications)
drift-detector detect --dry-run
```

## 📖 Usage Examples

### Basic Drift Detection

```bash
# Detect drift across all enabled providers
drift-detector detect

# Output:
# ✓ Loaded Terraform state (127 resources)
# ⚠ Drift detected in 3 resources:
#   - aws_instance.web-server-01: tags modified
#   - aws_s3_bucket.data-lake: encryption disabled
#   - aws_security_group.app-sg: rule added (port 22)
# 📤 Notifications sent to Slack
```

### Watch Mode (Continuous Monitoring)

```bash
drift-detector detect --watch --interval 10m
```

### CI/CD Integration

```yaml
# .github/workflows/drift-check.yml
name: Infrastructure Drift Check
on:
  schedule:
    - cron: '0 */6 * * *' # Every 6 hours
    
jobs:
  drift-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run Drift Detector
        env:
          AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
          AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK_URL }}
        run: |
          go run . detect --fail-on-drift
```

### Docker Usage

```bash
docker run -v $(pwd)/config:/app/config \
  -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
  -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
  meowtux/drift-detector:latest detect
```

## 🔐 Security Best Practices

1. **Credential Management**
   - Use environment variables for sensitive data
   - Never commit credentials to Git
   - Use IAM roles when running in cloud environments

2. **Minimum Permissions**
   ```json
   {
     "Version": "2012-10-17",
     "Statement": [{
       "Effect": "Allow",
       "Action": [
         "ec2:Describe*",
         "s3:GetBucketLocation",
         "s3:GetBucketVersioning",
         "rds:Describe*"
       ],
       "Resource": "*"
     }]
   }
   ```

3. **Encrypted Storage**
   - State files are read-only
   - Notifications use TLS/SSL
   - Sensitive config values support encryption

## 📊 Supported Resources

### AWS
- EC2 Instances, Security Groups, VPCs
- S3 Buckets, IAM Roles/Policies
- RDS Instances, Lambda Functions
- EKS Clusters, Load Balancers

### GCP
- Compute Instances, Firewall Rules
- Cloud Storage Buckets
- GKE Clusters, Cloud Functions

### Azure
- Virtual Machines, Network Security Groups
- Storage Accounts, Resource Groups
- AKS Clusters, App Services

## 🔔 Notification Channels

- **Slack**: Rich formatted messages with drift details
- **Email**: HTML emails with resource comparisons
- **Webhooks**: Custom HTTP endpoints
- **Discord**: Channel notifications
- **PagerDuty**: Critical drift alerts
- **Custom**: Plugin system for extensibility

## 🛠️ Configuration Reference

See [examples/config.yaml](examples/config.yaml) for full configuration options.

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the Apache2 LICENSE - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Terraform team for the amazing IaC tool
- Cloud provider SDKs
- Open source community

## 📞 Support

- 🐛 [Report Issues](https://github.com/MeowTux/drift-detector/issues)
- 💬 [Discussions](https://github.com/MeowTux/drift-detector/discussions)
- 📧 Email: i don't have any comersial email😂

## 🗺️ Roadmap

- [ ] Kubernetes drift detection
- [ ] Auto-remediation capabilities
- [ ] Web dashboard
- [ ] Drift history tracking
- [ ] Cost impact analysis
- [ ] Multi-account support

---

**Made with ❤️ by MeowTux**

**Star ⭐ this repository if you find it useful!**

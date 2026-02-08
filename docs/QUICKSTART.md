# Quick Start Guide

Get Drift Detector up and running in 5 minutes!

## Prerequisites

- Terraform state file (local or remote)
- Cloud provider credentials (AWS/GCP/Azure)
- Go 1.21+ (if building from source)

## Installation

### macOS/Linux (Quick Install)

```bash
curl -sSL https://raw.githubusercontent.com/MeowTux/drift-detector/main/install.sh | bash
```

### Manual Installation

```bash
# Download latest release
wget https://github.com/MeowTux/drift-detector/releases/latest/download/drift-detector-linux-amd64

# Make executable
chmod +x drift-detector-linux-amd64
sudo mv drift-detector-linux-amd64 /usr/local/bin/drift-detector

# Verify
drift-detector --version
```

### Build from Source

```bash
git clone https://github.com/MeowTux/drift-detector.git
cd drift-detector
make build
```

## Configuration

### Step 1: Initialize

```bash
drift-detector init
```

This creates `config/config.yaml` with default settings.

### Step 2: Configure Cloud Provider

Edit `config/config.yaml`:

```yaml
terraform:
  state_backend: "local"
  state_path: "./terraform.tfstate"

providers:
  aws:
    enabled: true
    regions: ["us-east-1"]
```

### Step 3: Set Credentials

```bash
# AWS
export AWS_ACCESS_KEY_ID="your_key"
export AWS_SECRET_ACCESS_KEY="your_secret"

# Slack (optional)
export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/..."
```

## Usage

### One-Time Detection

```bash
drift-detector detect
```

Output:
```
🔍 Starting drift detection...
✓ Loaded Terraform state (127 resources)
⚠ Drift detected in 2 resources:
  1. aws_s3_bucket.data-lake (aws_s3_bucket)
     - encryption: enabled → disabled
  2. aws_instance.web-server (aws_instance)
     - tags.Environment: production → staging
```

### Continuous Monitoring

```bash
drift-detector detect --watch --interval 5m
```

### Docker

```bash
docker run -v $(pwd)/config:/app/config \
  -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
  -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
  meowtux/drift-detector:latest detect
```

## Common Scenarios

### Scenario 1: Basic AWS Monitoring

```yaml
# config/config.yaml
terraform:
  state_path: "./terraform.tfstate"

providers:
  aws:
    enabled: true
    regions: ["us-east-1"]

notifications:
  slack:
    enabled: true
    webhook_url: "${SLACK_WEBHOOK_URL}"
```

```bash
drift-detector detect --watch --interval 10m
```

### Scenario 2: CI/CD Integration

```yaml
# .github/workflows/drift-check.yml
name: Drift Check
on:
  schedule:
    - cron: '0 */6 * * *'

jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run Drift Detector
        run: |
          curl -sSL https://raw.githubusercontent.com/MeowTux/drift-detector/main/install.sh | bash
          drift-detector detect --fail-on-drift
        env:
          AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
          AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
```

### Scenario 3: Multi-Region AWS

```yaml
providers:
  aws:
    enabled: true
    regions:
      - "us-east-1"
      - "us-west-2"
      - "eu-west-1"
```

### Scenario 4: Dry Run (Testing)

```bash
drift-detector detect --dry-run
```

No notifications are sent, only detection results are shown.

## Troubleshooting

### Error: "Failed to load Terraform state"

**Solution:** Check your `state_path` in config.yaml

```bash
# Verify state file exists
ls -la ./terraform.tfstate

# Or set correct path
drift-detector detect --config /path/to/config.yaml
```

### Error: "AWS credentials not found"

**Solution:** Set environment variables

```bash
export AWS_ACCESS_KEY_ID="your_key"
export AWS_SECRET_ACCESS_KEY="your_secret"
```

Or use AWS CLI configuration:
```bash
aws configure
```

### Error: "No drift detected but I made changes"

**Possible causes:**
1. Resource type not monitored (check `resources_to_monitor` in config)
2. Detection not covering all attributes
3. Terraform state not updated

**Solution:** Enable verbose logging

```bash
drift-detector detect -v
```

## Next Steps

1. **Add more cloud providers** - Enable GCP or Azure
2. **Set up notifications** - Configure Slack, email, or webhooks
3. **Automate checks** - Set up cron or CI/CD integration
4. **Monitor continuously** - Run in watch mode
5. **Secure your setup** - Review [Security Guide](docs/SECURITY.md)

## Resources

- [Full Documentation](README.md)
- [Security Best Practices](docs/SECURITY.md)
- [Android Setup](docs/ANDROID.md)
- [Contributing Guide](CONTRIBUTING.md)

## Getting Help

- 📖 [Documentation](https://github.com/MeowTux/drift-detector)
- 💬 [Discussions](https://github.com/MeowTux/drift-detector/discussions)
- 🐛 [Issues](https://github.com/MeowTux/drift-detector/issues)

---

**Ready to detect drift? Let's go! 🚀**

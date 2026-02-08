# Drift Detector on Android (Termux)

This guide explains how to run Drift Detector on Android devices using Termux.

## Prerequisites

1. Install **Termux** from F-Droid (NOT Google Play Store)
   - Download: https://f-droid.org/packages/com.termux/

2. Open Termux and update packages:
   ```bash
   pkg update && pkg upgrade
   ```

## Installation

### Option 1: Quick Install

```bash
# Install dependencies
pkg install golang git

# Clone and build
git clone https://github.com/MeowTux/drift-detector.git
cd drift-detector
go build -o drift-detector .

# Make executable
chmod +x drift-detector

# Test installation
./drift-detector --version
```

### Option 2: Download Pre-built Binary

```bash
# Install wget
pkg install wget

# Download ARM64 binary
wget https://github.com/MeowTux/drift-detector/releases/latest/download/drift-detector-linux-arm64

# Rename and make executable
mv drift-detector-linux-arm64 drift-detector
chmod +x drift-detector

# Test
./drift-detector --version
```

## Setup

### 1. Initialize Configuration

```bash
./drift-detector init
```

### 2. Edit Configuration

Use nano or vim to edit the config:

```bash
pkg install nano
nano config/config.yaml
```

### 3. Set Up Credentials

Create environment variables:

```bash
# Create .env file
nano .env

# Add your credentials:
export AWS_ACCESS_KEY_ID="your_key"
export AWS_SECRET_ACCESS_KEY="your_secret"
export SLACK_WEBHOOK_URL="your_webhook"

# Load environment
source .env
```

## Running Drift Detector

### One-time Check

```bash
./drift-detector detect
```

### Continuous Monitoring

```bash
./drift-detector detect --watch --interval 5m
```

### Background Process

Use `nohup` or `tmux`:

```bash
# Install tmux
pkg install tmux

# Start tmux session
tmux new -s drift

# Run detector
./drift-detector detect --watch --interval 10m

# Detach: Ctrl+B, then D
# Reattach: tmux attach -t drift
```

## Tips for Android

### Battery Optimization

- Add Termux to battery optimization whitelist
- Use longer check intervals (10m-30m) to save battery
- Consider running only when plugged in

### Automation

Create a script to auto-start on device boot:

```bash
#!/data/data/com.termux/files/usr/bin/bash
cd ~/drift-detector
source .env
./drift-detector detect --watch --interval 15m
```

### Storage Access

Grant storage permissions for Terraform state files:

```bash
termux-setup-storage
```

### Notifications

You can receive notifications via:
- Termux API notifications
- Slack/Discord webhooks
- Email

Example using Termux API:

```bash
pkg install termux-api

# In your notification script:
termux-notification --title "Drift Detected" --content "3 resources drifted"
```

## Troubleshooting

### Out of Memory

Reduce check frequency or resource monitoring scope in config.

### Network Issues

Ensure Termux has network permissions and internet access.

### Crashes

Check logs and reduce log level:
```bash
export LOG_LEVEL=error
```

## Performance Considerations

For optimal performance on mobile:
- Monitor fewer resources
- Use longer intervals (15m+)
- Limit to one cloud provider
- Reduce notification frequency

## Example Mobile Use Cases

1. **On-call monitoring**: Get alerts on your phone
2. **Remote administration**: Check drift while traveling
3. **Quick spot checks**: Run ad-hoc checks from anywhere
4. **Learning**: Experiment with Terraform on the go

## Security Notes

⚠️ **Important Security Considerations:**

- Never store credentials in plain text
- Use environment variables
- Lock your device with strong password/biometrics
- Consider using cloud-based secrets management
- Enable full-disk encryption on your device

## Resources

- Termux Wiki: https://wiki.termux.com
- F-Droid: https://f-droid.org
- Drift Detector Docs: https://github.com/MeowTux/drift-detector

---

**Happy Monitoring from Your Phone! 📱**

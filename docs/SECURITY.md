# Security Best Practices for Drift Detector

## Overview

Security is paramount when dealing with cloud infrastructure monitoring. This guide outlines best practices for securely deploying and using Drift Detector.

## 🔐 Credential Management

### DO ✅

1. **Use Environment Variables**
   ```bash
   export AWS_ACCESS_KEY_ID="your_key"
   export AWS_SECRET_ACCESS_KEY="your_secret"
   export SLACK_WEBHOOK_URL="your_webhook"
   ```

2. **Use Cloud IAM Roles** (Recommended)
   - EC2 Instance Roles
   - ECS Task Roles
   - GCP Service Accounts
   - Azure Managed Identities

3. **Encrypt Sensitive Data**
   - Use encrypted storage for state files
   - Enable encryption at rest for config files
   - Use secrets management services

4. **Rotate Credentials Regularly**
   - Set up automatic rotation (90 days)
   - Use temporary credentials when possible
   - Audit credential usage

### DON'T ❌

1. **Never Commit Credentials to Git**
   ```bash
   # Good: .gitignore includes
   .env
   config/config.yaml
   *.tfvars
   ```

2. **Never Use Root/Admin Credentials**
   - Create dedicated service accounts
   - Apply principle of least privilege

3. **Never Share Credentials**
   - Use separate credentials per environment
   - Don't reuse credentials across tools

## 🛡️ IAM Permissions

### Minimum AWS Permissions

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "DriftDetectorReadOnly",
      "Effect": "Allow",
      "Action": [
        "ec2:Describe*",
        "s3:GetBucketLocation",
        "s3:GetBucketVersioning",
        "s3:GetBucketEncryption",
        "s3:GetBucketPublicAccessBlock",
        "rds:Describe*",
        "lambda:GetFunction",
        "lambda:ListFunctions",
        "elasticloadbalancing:Describe*",
        "autoscaling:Describe*"
      ],
      "Resource": "*"
    },
    {
      "Sid": "TerraformStateAccess",
      "Effect": "Allow",
      "Action": [
        "s3:GetObject"
      ],
      "Resource": "arn:aws:s3:::my-terraform-state/*"
    }
  ]
}
```

### GCP Minimum Permissions

```yaml
roles:
  - compute.viewer
  - storage.objectViewer
  - iam.roleViewer
```

### Azure Minimum Permissions

- Reader role at subscription/resource group level
- Storage Blob Data Reader for state files

## 🔒 Network Security

### Deployment in Private Networks

1. **VPC/VNet Deployment**
   ```bash
   # Run in private subnet
   # No public IP needed
   # Use VPC endpoints for AWS API access
   ```

2. **Firewall Rules**
   - Restrict outbound to cloud APIs only
   - Block unnecessary inbound traffic
   - Use security groups/NSGs

3. **Proxy Configuration**
   ```bash
   export HTTP_PROXY="http://proxy:8080"
   export HTTPS_PROXY="http://proxy:8080"
   export NO_PROXY="169.254.169.254" # AWS metadata
   ```

## 🔍 Audit and Monitoring

### Enable Logging

```yaml
logging:
  level: "info"
  file: "/var/log/drift-detector/app.log"
  format: "json"
```

### Audit Trail

1. **Log All Detections**
   - Timestamp
   - User/Role
   - Resources checked
   - Drift found

2. **Monitor Access**
   - CloudTrail (AWS)
   - Cloud Logging (GCP)
   - Activity Logs (Azure)

3. **Alert on Anomalies**
   - Unexpected API calls
   - Failed authentication
   - Unusual patterns

## 🛠️ Secure Configuration

### Configuration File Security

```bash
# Set proper permissions
chmod 600 config/config.yaml
chown drift-user:drift-group config/config.yaml

# Encrypt sensitive sections
# Use tools like:
# - AWS KMS
# - GCP Secret Manager
# - Azure Key Vault
# - HashiCorp Vault
```

### Docker Security

```dockerfile
# Run as non-root
USER drift

# Read-only filesystem
docker run --read-only ...

# Drop capabilities
docker run --cap-drop=ALL ...

# Security scanning
docker scan meowtux/drift-detector:latest
```

### Kubernetes Security

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: drift-detector
spec:
  securityContext:
    runAsNonRoot: true
    runAsUser: 1000
    fsGroup: 1000
    seccompProfile:
      type: RuntimeDefault
  containers:
  - name: drift-detector
    image: meowtux/drift-detector:latest
    securityContext:
      allowPrivilegeEscalation: false
      readOnlyRootFilesystem: true
      capabilities:
        drop:
        - ALL
```

## 🚨 Incident Response

### If Credentials Are Compromised

1. **Immediate Actions**
   ```bash
   # Revoke compromised credentials
   aws iam delete-access-key --access-key-id AKIA...
   
   # Rotate all secrets
   # Review CloudTrail logs for unauthorized access
   # Change all webhook URLs
   ```

2. **Investigation**
   - Check all resource modifications
   - Review notification history
   - Audit state file access

3. **Prevention**
   - Enable MFA on accounts
   - Implement IP whitelisting
   - Use temporary credentials

## 📋 Security Checklist

- [ ] Credentials stored in environment variables or secrets manager
- [ ] IAM policies follow principle of least privilege
- [ ] No credentials in Git repository
- [ ] Config files have restricted permissions (600)
- [ ] Running as non-root user
- [ ] Network restricted to necessary endpoints
- [ ] Logging enabled and monitored
- [ ] Regular credential rotation scheduled
- [ ] State files encrypted at rest
- [ ] HTTPS/TLS for all communications
- [ ] Security scanning in CI/CD pipeline
- [ ] Incident response plan documented
- [ ] Regular security audits scheduled

## 🔗 Additional Resources

- [AWS Security Best Practices](https://aws.amazon.com/security/best-practices/)
- [GCP Security Best Practices](https://cloud.google.com/security/best-practices)
- [Azure Security Best Practices](https://docs.microsoft.com/en-us/azure/security/fundamentals/best-practices-and-patterns)
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [CIS Benchmarks](https://www.cisecurity.org/cis-benchmarks/)

## 📧 Report Security Issues

If you discover a security vulnerability, please email:
**security@meowtux.dev**

Do NOT create a public GitHub issue for security vulnerabilities.

---

**Security is everyone's responsibility. Stay vigilant! 🛡️**

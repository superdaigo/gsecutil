# Enabling Data Access Audit Logs for Secret Manager

This document explains how to enable Data Access audit logs for Google Cloud Secret Manager, which are required for the `gsecutil auditlog` command to show comprehensive audit information.

## Overview

By default, Google Cloud only logs Admin Activity events (like creating or deleting secrets). To see who accessed secret values and when, you need to enable **Data Access audit logs** for the Secret Manager API.

## Why Enable Data Access Audit Logs?

Data Access audit logs capture:
- **Secret value access** (when someone reads a secret)
- **Secret metadata access** (when someone lists secrets or gets secret info)
- **Who** performed the action (user email or service account)
- **When** the action occurred (timestamp)
- **What** resource was accessed (specific secret name)

Without Data Access audit logs enabled, `gsecutil auditlog` will only show:
- Secret creation/deletion (Admin Activity logs)
- Limited metadata operations
- No information about actual secret value access

## Enabling Data Access Audit Logs

### Method 1: Using Google Cloud Console (Recommended for beginners)

1. **Navigate to IAM & Admin > Audit Logs**
   - Go to the [Google Cloud Console](https://console.cloud.google.com/)
   - Select your project
   - Navigate to **IAM & Admin** → **Audit Logs**

2. **Find Secret Manager API**
   - In the audit logs configuration page, find **Secret Manager API**
   - You can use the search box to filter for "Secret Manager"

3. **Enable Data Access Logs**
   - Click on **Secret Manager API**
   - Check the boxes for:
     - ✅ **Admin Read** (recommended)
     - ✅ **Data Read** (required for secret access logs)
     - ✅ **Data Write** (recommended for secret creation/update logs)
   - Click **Save**

### Method 2: Using gcloud CLI

You can enable audit logs using the gcloud command-line tool:

```bash
# Get current IAM policy
gcloud projects get-iam-policy PROJECT_ID --format=json > policy.json
```

Edit the `policy.json` file to add the audit configuration. Add this section to the policy:

```json
{
  "auditConfigs": [
    {
      "service": "secretmanager.googleapis.com",
      "auditLogConfigs": [
        {
          "logType": "ADMIN_READ"
        },
        {
          "logType": "DATA_READ"
        },
        {
          "logType": "DATA_WRITE"
        }
      ]
    }
  ],
  "bindings": [
    // ... existing bindings ...
  ],
  "etag": "...",
  "version": 1
}
```

Then apply the updated policy:

```bash
# Apply the updated policy
gcloud projects set-iam-policy PROJECT_ID policy.json
```

### Method 3: Using Terraform

If you're managing your infrastructure with Terraform, you can enable audit logs using the `google_project_iam_audit_config` resource:

```hcl
resource "google_project_iam_audit_config" "secret_manager_audit" {
  project = var.project_id
  service = "secretmanager.googleapis.com"

  audit_log_config {
    log_type = "ADMIN_READ"
  }

  audit_log_config {
    log_type = "DATA_READ"
  }

  audit_log_config {
    log_type = "DATA_WRITE"
  }
}
```

## Audit Log Types Explained

| Log Type | What It Captures | Recommended |
|----------|------------------|-------------|
| **ADMIN_READ** | Reading secret metadata, listing secrets | ✅ Yes |
| **DATA_READ** | Accessing secret values | ✅ **Required** |
| **DATA_WRITE** | Creating/updating secret values | ✅ Yes |

## Cost Considerations

⚠️ **Important**: Enabling Data Access audit logs will increase your Cloud Logging costs, as these logs are more verbose than Admin Activity logs.

**Cost optimization tips:**
- Consider enabling audit logs only for production projects where security monitoring is critical
- Use log retention policies to automatically delete old audit logs
- Set up log exclusion filters if you need to reduce specific log types

**Estimated impact:**
- Small projects (< 100 secrets, low access): ~$5-20/month additional logging costs
- Medium projects (100-1000 secrets, moderate access): ~$20-100/month
- Large projects (> 1000 secrets, high access): ~$100+/month

## Verifying Audit Log Configuration

After enabling audit logs, you can verify the configuration:

### Using gcloud CLI:
```bash
# Check current audit log configuration
gcloud projects get-iam-policy PROJECT_ID --format="value(auditConfigs)"
```

### Using gsecutil:
```bash
# Test that audit logs are working (may take 10-15 minutes for logs to appear)
gsecutil auditlog your-secret-name --days 1
```

### Manual verification with Cloud Logging:
1. Go to **Cloud Logging** → **Logs Explorer**
2. Use this filter:
   ```
   protoPayload.serviceName="secretmanager.googleapis.com"
   protoPayload.methodName="google.cloud.secretmanager.v1.SecretManagerService.AccessSecretVersion"
   ```
3. Access a secret value and check if logs appear within 10-15 minutes

## Troubleshooting

### "No audit log entries found" message

If `gsecutil auditlog` shows no results:

1. **Check if audit logs are enabled** (see verification steps above)
2. **Wait for log propagation** - audit logs can take 10-15 minutes to appear
3. **Verify permissions** - ensure your account has `logging.logEntries.list` permission
4. **Check the time range** - use `--days` flag to expand the search window

### Permission errors

If you get permission errors when enabling audit logs:

- You need the **Project IAM Admin** role (`roles/resourcemanager.projectIamAdmin`)
- Or the **Security Admin** role (`roles/iam.securityAdmin`)
- Or custom permissions: `resourcemanager.projects.getIamPolicy` and `resourcemanager.projects.setIamPolicy`

### High logging costs

If audit logging costs are too high:

1. **Create exclusion filters** in Cloud Logging to exclude certain log types
2. **Set up log retention policies** to delete logs after a specific period
3. **Use log-based metrics** instead of storing all raw logs
4. **Consider enabling audit logs only for critical secrets** using resource-specific configurations

## Security Best Practices

1. **Monitor who enables/disables audit logs** - Changes to audit configuration should be monitored
2. **Restrict access to audit logs** - Only security teams should have access to audit log data
3. **Set up alerting** - Create alerts for suspicious secret access patterns
4. **Regular auditing** - Periodically review audit logs for unauthorized access

## Example Usage After Setup

Once Data Access audit logs are enabled, you can use the full power of `gsecutil auditlog`:

```bash
# Show who accessed a secret in the last 7 days
gsecutil auditlog my-secret

# Show detailed access for the last 30 days
gsecutil auditlog my-secret --days 30

# Get JSON output for programmatic processing
gsecutil auditlog my-secret --format json

# Limit results to most recent 10 entries
gsecutil auditlog my-secret --limit 10
```

## Resources

- [Google Cloud Audit Logs Documentation](https://cloud.google.com/logging/docs/audit)
- [Secret Manager Security Best Practices](https://cloud.google.com/secret-manager/docs/best-practices)
- [Cloud Logging Pricing](https://cloud.google.com/stackdriver/pricing)
- [IAM Audit Configuration](https://cloud.google.com/resource-manager/docs/audit-logging)

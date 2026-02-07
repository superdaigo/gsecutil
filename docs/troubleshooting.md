# Troubleshooting Guide

Common issues and solutions when using `gsecutil`.

## Table of Contents

- [Installation Issues](#installation-issues)
- [Authentication Issues](#authentication-issues)
- [Permission Issues](#permission-issues)
- [Project Configuration Issues](#project-configuration-issues)
- [Secret Access Issues](#secret-access-issues)
- [Clipboard Issues](#clipboard-issues)
- [Configuration File Issues](#configuration-file-issues)
- [Audit Log Issues](#audit-log-issues)
- [Version Management Issues](#version-management-issues)
- [Network and Connection Issues](#network-and-connection-issues)

---

## Installation Issues

### "gcloud command not found"

**Problem:** Running `gsecutil` commands results in "gcloud command not found" error.

**Solution:**
1. Install Google Cloud SDK:
   ```bash
   # macOS (Homebrew)
   brew install --cask google-cloud-sdk
   
   # Linux
   curl https://sdk.cloud.google.com | bash
   exec -l $SHELL
   
   # Windows
   # Download installer from https://cloud.google.com/sdk/docs/install
   ```

2. Verify installation:
   ```bash
   gcloud version
   ```

3. Ensure `gcloud` is in your PATH:
   ```bash
   which gcloud  # Unix/macOS
   where gcloud  # Windows
   ```

### "Permission denied" when installing binary

**Problem:** Cannot move binary to `/usr/local/bin/` or system directory.

**Solution:**
```bash
# Use sudo for system directories
sudo mv gsecutil /usr/local/bin/

# Or install to user directory
mkdir -p ~/bin
mv gsecutil ~/bin/
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bashrc  # or ~/.zshrc
source ~/.bashrc
```

---

## Authentication Issues

### "There was a problem refreshing your current auth tokens"

**Problem:**
```
Error: gsecutil: Error executing gcloud command:
────────────────────────────────────────
ERROR: (gcloud.secrets.list) There was a problem refreshing your current auth tokens: 
('invalid_grant: Bad Request', {'error': 'invalid_grant', 'error_description': 'Bad Request'})
────────────────────────────────────────
```

**Solution:**
1. Re-authenticate with gcloud:
   ```bash
   gcloud auth login
   ```

2. If using service account:
   ```bash
   gcloud auth activate-service-account --key-file=service-account.json
   ```

3. For application default credentials:
   ```bash
   gcloud auth application-default login
   ```

### "You do not currently have an active account selected"

**Problem:** No gcloud account is active.

**Solution:**
1. List available accounts:
   ```bash
   gcloud auth list
   ```

2. Set active account:
   ```bash
   gcloud config set account YOUR_EMAIL@example.com
   ```

3. Or login with new account:
   ```bash
   gcloud auth login
   ```

---

## Permission Issues

### "Permission denied" on secret operations

**Problem:**
```
ERROR: (gcloud.secrets.get) PERMISSION_DENIED: Permission 'secretmanager.versions.access' denied
```

**Solution:**
1. Check your IAM roles:
   ```bash
   gcloud projects get-iam-policy YOUR_PROJECT_ID \
     --flatten="bindings[].members" \
     --filter="bindings.members:user:YOUR_EMAIL"
   ```

2. Required roles for different operations:
   - **Read secrets**: `roles/secretmanager.secretAccessor`
   - **Create/update secrets**: `roles/secretmanager.secretVersionAdder`
   - **Manage versions**: `roles/secretmanager.secretVersionManager`
   - **Full access**: `roles/secretmanager.admin`

3. Request access from project owner or add role:
   ```bash
   gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
     --member="user:YOUR_EMAIL" \
     --role="roles/secretmanager.secretAccessor"
   ```

### "Secret Manager API has not been used"

**Problem:**
```
ERROR: (gcloud.secrets.list) Secret Manager API has not been used in project PROJECT_ID
```

**Solution:**
1. Enable Secret Manager API:
   ```bash
   gcloud services enable secretmanager.googleapis.com --project=YOUR_PROJECT_ID
   ```

2. Verify API is enabled:
   ```bash
   gcloud services list --enabled --project=YOUR_PROJECT_ID | grep secretmanager
   ```

---

## Project Configuration Issues

### "Failed to find attribute [project]"

**Problem:**
```
ERROR: (gcloud.secrets.create) Error parsing [secret].
Failed to find attribute [project].
```

**Solution:**
Choose one of these methods to set project:

1. **Using CLI flag:**
   ```bash
   gsecutil create my-secret --project YOUR_PROJECT_ID
   ```

2. **Using environment variable:**
   ```bash
   export GSECUTIL_PROJECT=YOUR_PROJECT_ID
   gsecutil create my-secret
   ```

3. **Using configuration file:**
   ```yaml
   # ~/.config/gsecutil/gsecutil.conf
   project: "YOUR_PROJECT_ID"
   ```

4. **Using gcloud default:**
   ```bash
   gcloud config set project YOUR_PROJECT_ID
   ```

### Wrong project being used

**Problem:** Commands are using incorrect project.

**Solution:**
Check project resolution order (highest to lowest priority):
1. `--project` CLI flag
2. Configuration file (`~/.config/gsecutil/gsecutil.conf`)
3. `GSECUTIL_PROJECT` environment variable
4. gcloud default project

Debug current configuration:
```bash
# Check gcloud default
gcloud config get-value project

# Check environment variable
echo $GSECUTIL_PROJECT

# Check config file
gsecutil config show
```

---

## Secret Access Issues

### "Resource not found"

**Problem:**
```
ERROR: (gcloud.secrets.versions.access) NOT_FOUND: Secret [SECRET_NAME] not found
```

**Solution:**
1. Verify secret exists:
   ```bash
   gsecutil list
   ```

2. Check if using prefix in config:
   ```bash
   gsecutil config show
   ```
   If prefix is configured (e.g., `team-`), use the short name:
   ```bash
   # Correct (if prefix is "team-")
   gsecutil get database-password
   
   # Incorrect (don't include prefix manually)
   gsecutil get team-database-password
   ```

3. Verify correct project:
   ```bash
   gsecutil list --project YOUR_PROJECT_ID
   ```

### Cannot access specific version

**Problem:**
```
ERROR: (gcloud.secrets.versions.access) Secret version [VERSION] not found
```

**Solution:**
1. Check available versions:
   ```bash
   gsecutil describe SECRET_NAME --show-versions
   ```

2. Verify version state (must be ENABLED):
   ```bash
   gsecutil describe SECRET_NAME --show-versions
   ```

3. Use version number, not alias:
   ```bash
   # Correct
   gsecutil get my-secret --version 3
   
   # May not work with aliases
   gsecutil get my-secret --version stable
   ```

---

## Clipboard Issues

### Clipboard not working on Linux

**Problem:** `--clipboard` flag fails or shows warning.

**Solution:**
1. Install clipboard utilities:
   ```bash
   # Ubuntu/Debian
   sudo apt-get install xclip xsel
   
   # Fedora/RHEL
   sudo dnf install xclip xsel
   
   # Arch
   sudo pacman -S xclip xsel
   ```

2. For Wayland (instead of X11):
   ```bash
   sudo apt-get install wl-clipboard
   ```

3. Verify you have a graphical environment (X11 or Wayland).

### Clipboard fails on headless server

**Problem:** No clipboard available on SSH/headless server.

**Solution:**
Don't use `--clipboard` on headless servers. Instead:
```bash
# Pipe to file
gsecutil get my-secret > secret.txt

# Use in command directly
MY_SECRET=$(gsecutil get my-secret)
echo $MY_SECRET
```

---

## Configuration File Issues

### "Failed to parse config file"

**Problem:**
```
Warning: failed to parse config file: yaml: line 5: ...
```

**Solution:**
1. Validate YAML syntax:
   ```bash
   gsecutil config validate
   ```

2. Check for common YAML errors:
   - Incorrect indentation (use spaces, not tabs)
   - Missing colons
   - Unquoted special characters

3. Example of correct YAML:
   ```yaml
   project: "my-project"
   prefix: "team-"
   list:
     attributes:
       - title
       - owner
   ```

### Configuration not being used

**Problem:** Configuration file exists but settings aren't applied.

**Solution:**
1. Verify file location:
   ```bash
   # Default location
   ls -l ~/.config/gsecutil/gsecutil.conf
   ```

2. Check file is valid:
   ```bash
   gsecutil config show
   ```

3. Specify config explicitly:
   ```bash
   gsecutil --config /path/to/config.yaml list
   ```

### "Prefix cannot contain spaces"

**Problem:**
```
ERROR: configuration validation failed: prefix cannot contain spaces
```

**Solution:**
Remove spaces from prefix in config file:
```yaml
# Incorrect
prefix: "my team "

# Correct
prefix: "my-team-"
prefix: "myteam_"
```

---

## Audit Log Issues

### "Audit log entries not found"

**Problem:** `gsecutil auditlog` returns no results even though secrets were accessed.

**Solution:**
1. Enable Data Access audit logs for Secret Manager:
   - Go to [Cloud Console Audit Logs](https://console.cloud.google.com/iam-admin/audit)
   - Find "Secret Manager API"
   - Enable "Admin Read", "Data Read", and "Data Write"
   - See [docs/audit-logging.md](audit-logging.md) for detailed setup

2. Wait for logs to propagate (can take a few minutes).

3. Check correct time range:
   ```bash
   # Last 30 days
   gsecutil auditlog my-secret --days 30
   ```

### "Permission denied" on audit logs

**Problem:**
```
ERROR: (gcloud.logging.read) PERMISSION_DENIED: User does not have permission to access logs
```

**Solution:**
1. Required role: `roles/logging.viewer` or higher

2. Add logging viewer role:
   ```bash
   gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
     --member="user:YOUR_EMAIL" \
     --role="roles/logging.viewer"
   ```

---

## Version Management Issues

### "Would exceed free tier limit"

**Problem:** Prompted about exceeding 6 active versions when creating/updating.

**Solution:**
1. **Recommended:** Disable old versions (stay within free tier):
   - Answer `y` or `yes` when prompted
   - Old versions will be automatically disabled
   - Latest version is always preserved

2. **Proceed anyway** (may incur charges):
   - Answer `n` or `no` when prompted
   - You'll be charged for versions beyond 6

3. **Bypass check** (if you know what you're doing):
   ```bash
   gsecutil create my-secret --data "value" --force
   gsecutil update my-secret --data "value" --force
   ```

### Cannot access disabled version

**Problem:**
```
ERROR: Secret version is DISABLED and cannot be accessed
```

**Solution:**
1. Enable the version:
   ```bash
   gcloud secrets versions enable VERSION --secret=SECRET_NAME
   ```

2. Or access a different version:
   ```bash
   gsecutil describe my-secret --show-versions
   gsecutil get my-secret --version 5  # Use ENABLED version
   ```

---

## Network and Connection Issues

### "Failed to connect" or timeout errors

**Problem:** Commands hang or timeout.

**Solution:**
1. Check internet connectivity:
   ```bash
   ping -c 3 google.com
   ```

2. Verify Google Cloud APIs are accessible:
   ```bash
   curl https://secretmanager.googleapis.com/
   ```

3. Check firewall/proxy settings:
   - Ensure outbound HTTPS (443) is allowed
   - Configure proxy if needed:
     ```bash
     export HTTP_PROXY=http://proxy.example.com:8080
     export HTTPS_PROXY=http://proxy.example.com:8080
     ```

4. Test gcloud connectivity:
   ```bash
   gcloud info --run-diagnostics
   ```

### "SSL certificate verification failed"

**Problem:** SSL/TLS errors when connecting to Google APIs.

**Solution:**
1. Update CA certificates:
   ```bash
   # Ubuntu/Debian
   sudo apt-get update && sudo apt-get install ca-certificates
   
   # macOS
   brew install ca-certificates
   ```

2. Update gcloud SDK:
   ```bash
   gcloud components update
   ```

---

## Debug Mode

For detailed error information, enable gcloud verbose output:

```bash
# Set gcloud to debug mode
export CLOUDSDK_CORE_VERBOSITY=debug

# Run your command
gsecutil list

# Disable debug mode
unset CLOUDSDK_CORE_VERBOSITY
```

## Getting Help

If you're still experiencing issues:

1. **Check command help:**
   ```bash
   gsecutil --help
   gsecutil COMMAND --help
   ```

2. **Validate your setup:**
   ```bash
   # Check gcloud authentication
   gcloud auth list
   
   # Check project configuration
   gcloud config get-value project
   
   # Test Secret Manager access
   gcloud secrets list --limit 1
   ```

3. **Review logs:**
   - Check `~/.config/gcloud/logs/` for gcloud logs
   - Enable debug mode as shown above

4. **Create an issue:**
   - Visit [GitHub Issues](https://github.com/superdaigo/gsecutil/issues)
   - Include error message, command used, and `gsecutil --version`

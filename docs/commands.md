# Command Reference

Complete reference for all `gsecutil` commands.

## Table of Contents

- [Secret Management](#secret-management)
  - [create](#create) - Create a new secret
  - [get](#get) - Retrieve a secret value
  - [update](#update) - Update an existing secret
  - [delete](#delete) - Delete a secret
  - [list](#list) - List all secrets
  - [describe](#describe) - Show secret details
- [Bulk Operations](#bulk-operations)
  - [import](#import) - Import secrets from CSV
  - [export](#export) - Export secrets to CSV
- [Configuration](#configuration)
  - [config init](#config-init) - Initialize configuration
  - [config show](#config-show) - Show configuration
  - [config validate](#config-validate) - Validate configuration
  - [config import](#config-import) - Import configuration
- [Access Management](#access-management)
  - [access list](#access-list) - List access permissions
  - [access grant](#access-grant) - Grant access
  - [access revoke](#access-revoke) - Revoke access
  - [access project](#access-project) - Show project permissions
- [Audit Logs](#audit-logs)
  - [auditlog](#auditlog) - View audit logs

---

## Secret Management

### create

Create a new secret in Google Secret Manager.

**Usage:**
```bash
gsecutil create SECRET_NAME [flags]
```

**Flags:**
- `-d, --data` - Secret data to store
- `--data-file` - Path to file containing secret data
- `--labels` - Labels to apply (format: key=value)
- `-f, --force` - Force creation without version limit checks

**Examples:**
```bash
# Interactive input (secure prompt)
gsecutil create database-password

# From command line
gsecutil create api-key -d "sk-1234567890"

# From file
gsecutil create config --data-file ./config.json

# From stdin
echo "secret-value" | gsecutil create my-secret --data-file -

# With labels
gsecutil create api-key -d "sk-123" --labels env=prod,team=backend
```

**Version Management:**
The free tier allows up to 6 active secret versions. If creating a secret that already exists would exceed this limit, you'll be prompted to disable old versions or proceed anyway.

---

### get

Retrieve a secret value from Google Secret Manager.

**Usage:**
```bash
gsecutil get SECRET_NAME [flags]
```

**Flags:**
- `-v, --version` - Version number to retrieve (default: latest)
- `-c, --clipboard` - Copy secret value to clipboard
- `-m, --show-metadata` - Show version metadata (version, state, created time)

**Examples:**
```bash
# Get latest version
gsecutil get database-password

# Get specific version
gsecutil get api-key --version 3

# Copy to clipboard
gsecutil get api-key --clipboard

# Show metadata
gsecutil get api-key --show-metadata

# Combine options
gsecutil get api-key -v 2 -c -m
```

---

### update

Update an existing secret by creating a new version.

**Usage:**
```bash
gsecutil update SECRET_NAME [flags]
```

**Flags:**
- `-d, --data` - New secret data
- `--data-file` - Path to file containing new secret data
- `-f, --force` - Force update without version limit checks

**Examples:**
```bash
# Interactive input
gsecutil update database-password

# From command line
gsecutil update api-key -d "new-secret-value"

# From file
gsecutil update config --data-file ./new-config.json

# Force update (skip version check)
gsecutil update api-key -d "new-value" --force
```

---

### delete

Delete a secret permanently from Google Secret Manager.

**Usage:**
```bash
gsecutil delete SECRET_NAME [flags]
```

**Flags:**
- `-f, --force` - Force deletion without confirmation

**Examples:**
```bash
# With confirmation prompt
gsecutil delete old-secret

# Force delete (no prompt)
gsecutil delete old-secret --force
```

---

### list

List all secrets in the project.

**Usage:**
```bash
gsecutil list [flags]
```

**Flags:**
- `--filter` - Filter expression for Secret Manager labels
- `--filter-attributes` - Filter by config attributes (format: key=value,key2=value2)
- `--format` - Output format (json, yaml, table)
- `--limit` - Maximum number of secrets to list
- `--no-labels` - Hide labels in output
- `--principal` - List secrets accessible by this principal
- `--show-attributes` - Comma-separated attributes to display from config

**Examples:**
```bash
# List all secrets
gsecutil list

# Filter by Secret Manager label
gsecutil list --filter "labels.env=prod"

# Filter by config attributes
gsecutil list --filter-attributes "environment=production,owner=backend-team"

# Show specific attributes
gsecutil list --show-attributes "title,owner,environment"

# List with limit
gsecutil list --limit 10

# List secrets accessible by user
gsecutil list --principal user:alice@example.com

# JSON output
gsecutil list --format json
```

---

### describe

Get detailed information about a secret.

**Usage:**
```bash
gsecutil describe SECRET_NAME [flags]
```

**Flags:**
- `-v, --show-versions` - Show detailed version information
- `--format` - Output format (json, yaml)

**Examples:**
```bash
# Basic description
gsecutil describe database-password

# With version history
gsecutil describe database-password --show-versions

# JSON output
gsecutil describe database-password --format json
```

**Information Displayed:**
- Basic metadata (name, creation time, ETag)
- Config attributes (from configuration file)
- Default version information
- Replication strategy
- Labels and annotations
- Version aliases
- Expiration and rotation settings
- Pub/Sub topics

---

## Bulk Operations

### import

Import secrets from a CSV file.

**Usage:**
```bash
gsecutil import <csv-file> [flags]
```

**Flags:**
- `--dry-run` - Preview changes without executing
- `--update` - Update existing secrets only
- `--upsert` - Create or update secrets (upsert mode)
- `--update-config` - Update configuration file with metadata from CSV

**Examples:**
```bash
# Import new secrets (skip existing)
gsecutil import secrets.csv

# Preview without changes
gsecutil import secrets.csv --dry-run

# Update existing secrets
gsecutil import secrets.csv --update

# Create or update (upsert)
gsecutil import secrets.csv --upsert

# Update config file with metadata
gsecutil import secrets.csv --upsert --update-config
```

**CSV Format:**
- Required columns: `name`, `value` (for creation)
- Optional columns: `title`, `label:<key>`, custom attributes
- Supports Excel multi-line cells

**See Also:** [CSV Operations Guide](csv-operations.md) for detailed documentation.

---

### export

Export secrets to CSV format.

**Usage:**
```bash
gsecutil export [output-file] [flags]
```

**Flags:**
- `-o, --output` - Output file path (default: stdout)
- `--with-values` - Include secret values in export
- `--filter` - Filter secrets by label

**Examples:**
```bash
# Export to file (without values)
gsecutil export -o secrets.csv

# Export to stdout
gsecutil export

# Export with secret values
gsecutil export --with-values -o backup.csv

# Export filtered secrets
gsecutil export --filter env=production -o prod-secrets.csv
```

**See Also:** [CSV Operations Guide](csv-operations.md) for detailed documentation.

---

## Configuration

### config init

Initialize configuration file interactively.

**Usage:**
```bash
gsecutil config init [flags]
```

**Flags:**
- `-o, --output` - Output path (default: ~/.config/gsecutil/gsecutil.conf)
- `-f, --force` - Overwrite existing configuration

**Example:**
```bash
# Interactive setup
gsecutil config init

# Custom output path
gsecutil config init --output /path/to/config.yaml

# Overwrite existing
gsecutil config init --force
```

---

### config show

Show configuration file contents.

**Usage:**
```bash
gsecutil config show [config-file] [flags]
```

**Flags:**
- `-c, --show-credentials` - Show credentials table

**Examples:**
```bash
# Show default config
gsecutil config show

# Show specific config file
gsecutil config show /path/to/config.yaml

# Show with credentials table
gsecutil config show --show-credentials
```

---

### config validate

Validate configuration file for correctness.

**Usage:**
```bash
gsecutil config validate [config-file] [flags]
```

**Flags:**
- `-v, --verbose` - Show detailed validation results

**Examples:**
```bash
# Validate default config
gsecutil config validate

# Validate specific file
gsecutil config validate /path/to/config.yaml

# Verbose output
gsecutil config validate --verbose
```

---

### config import

Import configuration from an existing file.

**Usage:**
```bash
gsecutil config import <source-file> [flags]
```

**Flags:**
- `-o, --output` - Output path (default: ~/.config/gsecutil/gsecutil.conf)
- `-f, --force` - Overwrite existing configuration

**Examples:**
```bash
# Import to default location
gsecutil config import team-config.yaml

# Import to custom location
gsecutil config import team-config.yaml --output ~/.config/gsecutil/prod.conf

# Force overwrite
gsecutil config import team-config.yaml --force
```

---

## Access Management

### access list

List principals with access to a secret.

**Usage:**
```bash
gsecutil access list <secret> [flags]
```

**Flags:**
- `--include-project` - Include project-level permissions

**Examples:**
```bash
# List secret-level access
gsecutil access list my-secret

# Include project-level permissions
gsecutil access list my-secret --include-project
```

---

### access grant

Grant access to a principal for a secret.

**Usage:**
```bash
gsecutil access grant <secret> --principal <principal> [flags]
```

**Flags:**
- `--principal` - Principal to grant access (required)
- `--role` - Role to grant (default: roles/secretmanager.secretAccessor)

**Principal Formats:**
- `user:email@domain.com`
- `group:group@domain.com`
- `serviceAccount:sa@project.iam.gserviceaccount.com`
- `domain:domain.com`

**Available Roles:**
- `roles/secretmanager.secretAccessor` - Can access secret values
- `roles/secretmanager.viewer` - Can view metadata only
- `roles/secretmanager.admin` - Full access
- `roles/secretmanager.secretVersionManager` - Can create/destroy versions
- `roles/secretmanager.secretVersionAdder` - Can add new versions

**Examples:**
```bash
# Grant default access (secretAccessor)
gsecutil access grant my-secret --principal user:alice@example.com

# Grant specific role
gsecutil access grant my-secret \
  --principal user:alice@example.com \
  --role roles/secretmanager.viewer

# Grant to service account
gsecutil access grant my-secret \
  --principal serviceAccount:app@project.iam.gserviceaccount.com
```

---

### access revoke

Revoke access from a principal for a secret.

**Usage:**
```bash
gsecutil access revoke <secret> --principal <principal> [flags]
```

**Flags:**
- `--principal` - Principal to revoke access from (required)
- `--role` - Role to revoke (default: roles/secretmanager.secretAccessor)

**Examples:**
```bash
# Revoke default access
gsecutil access revoke my-secret --principal user:bob@example.com

# Revoke specific role
gsecutil access revoke my-secret \
  --principal user:bob@example.com \
  --role roles/secretmanager.viewer
```

---

### access project

Show project-level Secret Manager permissions.

**Usage:**
```bash
gsecutil access project
```

**Example:**
```bash
gsecutil access project
```

---

## Audit Logs

### auditlog

Show audit log entries for secrets.

**Usage:**
```bash
gsecutil auditlog [SECRET_NAME] [flags]
```

**Flags:**
- `--days` - Number of days to look back (default: 7)
- `--limit` - Maximum number of entries (default: 100)
- `--format` - Output format (table, json)
- `--principal` - Filter by principal (supports partial matching)
- `--operation` - Filter by operation (comma-separated)

**Available Operations:**
- `ACCESS` - Reading secret values
- `CREATE` - Creating new secrets
- `UPDATE` - Creating new secret versions
- `DELETE` - Deleting secrets
- `GET_METADATA` - Getting secret/version metadata
- `LIST` - Listing secrets
- `UPDATE_METADATA` - Updating secret metadata
- `DESTROY_VERSION` - Destroying specific versions
- `DISABLE_VERSION` - Disabling specific versions
- `ENABLE_VERSION` - Enabling specific versions

**Examples:**
```bash
# Show all Secret Manager audit logs
gsecutil auditlog

# Show logs for specific secret (supports partial matching)
gsecutil auditlog database-password

# Filter by principal
gsecutil auditlog --principal alice

# Filter by operation
gsecutil auditlog --operation ACCESS,CREATE

# Combine filters
gsecutil auditlog db --principal admin --operation UPDATE

# Last 30 days
gsecutil auditlog my-secret --days 30

# JSON output
gsecutil auditlog my-secret --format json

# Limit results
gsecutil auditlog --limit 50
```

**Note:** Requires Data Access audit logs to be enabled for Secret Manager API. See [docs/audit-logging.md](audit-logging.md) for setup instructions.

---

## Global Flags

These flags are available for all commands:

- `-p, --project` - Google Cloud project ID
- `--config` - Configuration file path (default: ~/.config/gsecutil/gsecutil.conf)
- `-h, --help` - Show help for command

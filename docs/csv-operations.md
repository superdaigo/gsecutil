# CSV Import/Export Operations

Complete guide for CSV-based bulk operations with `gsecutil`.

## Table of Contents

- [Overview](#overview)
- [Export Command](#export-command)
- [Import Command](#import-command)
- [CSV Format](#csv-format)
- [Common Workflows](#common-workflows)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

---

## Overview

The CSV import/export functionality allows you to:

- **Export** secrets with metadata to CSV for backup or documentation
- **Import** secrets in bulk from CSV files
- **Migrate** secrets between projects or environments
- **Manage metadata** (titles, labels, attributes) alongside secrets
- **Work with Excel** - full support for multi-line cells

**Key Features:**
- Excel-compatible CSV format
- Multi-line value support (certificates, configs, etc.)
- Metadata integration with configuration file
- Dry-run mode for safe testing
- Flexible update modes (create-only, update-only, upsert)

---

## Export Command

Export secrets and their metadata to CSV format.

### Usage

```bash
gsecutil export [output-file] [flags]
```

### Flags

- `-o, --output <file>` - Output file path (default: stdout)
- `--with-values` - Include secret values in export (⚠️ use with caution)
- `--filter <label=value>` - Filter secrets by label

### Examples

#### Basic Export

```bash
# Export to file (without values)
gsecutil export -o secrets.csv

# Export to stdout
gsecutil export

# Export with secret values
gsecutil export --with-values -o backup.csv
```

#### Filtered Export

```bash
# Export only production secrets
gsecutil export --filter env=production -o prod-secrets.csv

# Export with multiple filters
gsecutil export --filter env=production --filter team=backend -o filtered.csv
```

### Output Format

The exported CSV includes:

| Column | Description | Always Present |
|--------|-------------|----------------|
| `name` | Secret name | ✓ |
| `value` | Secret value | Only with `--with-values` |
| `title` | Title from config | ✓ |
| `label:<key>` | Labels (e.g., `label:env`) | If labels exist |
| Custom columns | Config attributes (e.g., `owner`) | If attributes exist |

**Example CSV:**

```csv
name,title,label:env,label:team,owner,rotation_days
myapp-db-password,Database Password,production,backend,alice,30
myapp-api-key,API Key,production,frontend,bob,90
```

---

## Import Command

Import secrets from CSV files with flexible update modes.

### Usage

```bash
gsecutil import <csv-file> [flags]
```

### Flags

- `--dry-run` - Preview changes without executing
- `--update` - Update existing secrets only
- `--upsert` - Create new secrets and update existing ones
- `--update-config` - Save titles and attributes to configuration file

### Update Modes

| Mode | Behavior | Use Case |
|------|----------|----------|
| **Default** | Create new secrets, skip existing | Initial import |
| `--update` | Update existing, skip non-existent | Value updates only |
| `--upsert` | Create or update all | Full synchronization |

### Examples

#### Basic Import

```bash
# Import new secrets (skip existing)
gsecutil import secrets.csv

# Preview without changes
gsecutil import secrets.csv --dry-run
```

#### Update Modes

```bash
# Update existing secrets only
gsecutil import secrets.csv --update

# Create or update (upsert)
gsecutil import secrets.csv --upsert

# Update and save metadata to config
gsecutil import secrets.csv --upsert --update-config
```

---

## CSV Format

### Required Columns

- **`name`** - Secret name (required)
- **`value`** - Secret value (required for creation)

### Optional Columns

- **`title`** - Secret title (saved to config with `--update-config`)
- **`label:<key>`** - Labels applied to secrets (e.g., `label:env`, `label:team`)
- **Custom columns** - Any other column becomes a config attribute

### Multi-line Values

Excel-style multi-line cells are fully supported:

```csv
name,value,title
myapp-cert,"-----BEGIN CERTIFICATE-----
MIIBkTCB+wIJAKHHCgVZU6K9...
-----END CERTIFICATE-----",TLS Certificate
myapp-config,"database:
  host: localhost
  port: 5432",App Config
```

### Validation

The import command validates:
- ✅ Required columns (`name` must exist)
- ✅ No duplicate column names (case-insensitive)
- ✅ No empty column names
- ✅ Proper CSV format

**Error Example:**

```bash
$ gsecutil import bad.csv
Error: CSV header contains duplicate column names: 'owner' (columns 3, 4)
```

---

## Common Workflows

### 1. Initial Bulk Import

```bash
# Create secrets.csv with your data
cat > secrets.csv << 'EOF'
name,value,title,label:env,owner
myapp-db-prod,secret123,Production DB,production,alice
myapp-db-stage,secret456,Staging DB,staging,alice
myapp-api-key,key789,API Key,production,bob
EOF

# Import with config update
gsecutil import secrets.csv --update-config
```

### 2. Secret Migration Between Projects

```bash
# Step 1: Export from source project
gsecutil export --with-values --project source-project -o export.csv

# Step 2: Import to destination project
gsecutil import export.csv --project dest-project --upsert
```

### 3. Bulk Secret Updates

```bash
# Step 1: Export current state
gsecutil export -o current.csv

# Step 2: Edit values in Excel/text editor

# Step 3: Preview changes
gsecutil import current.csv --update --dry-run

# Step 4: Apply changes
gsecutil import current.csv --update
```

### 4. Metadata Management (Config-Only)

```bash
# Create metadata CSV (no values)
cat > metadata.csv << 'EOF'
name,title,owner,rotation_days,label:env
myapp-db-password,Database Password,alice,30,production
myapp-api-key,API Key,bob,90,production
EOF

# Update config only
gsecutil import metadata.csv --update-config
```

### 5. Environment Sync

```bash
# Export staging secrets
gsecutil export --filter env=staging --with-values -o staging.csv

# Modify for production (change names, labels)
sed 's/staging/production/g' staging.csv > production.csv

# Import to production
gsecutil import production.csv --upsert
```

### 6. Documentation Export

```bash
# Export without values for documentation
gsecutil export -o secrets-inventory.csv

# Add to version control (safe - no sensitive data)
git add secrets-inventory.csv
```

### 7. Backup and Restore

```bash
# Backup
gsecutil export --with-values -o backup-$(date +%Y%m%d).csv

# Store securely (encrypted storage recommended)
gpg --encrypt backup-20260207.csv

# Restore if needed
gpg --decrypt backup-20260207.csv.gpg | gsecutil import /dev/stdin --upsert
```

---

## Best Practices

### Security

1. **⚠️ Never commit CSVs with `--with-values` to version control**
   ```bash
   # Add to .gitignore
   echo "*-with-values.csv" >> .gitignore
   echo "backup-*.csv" >> .gitignore
   ```

2. **Always use `--dry-run` first**
   ```bash
   gsecutil import secrets.csv --dry-run
   # Review output, then run without --dry-run
   ```

3. **Store exported CSVs securely**
   ```bash
   # Encrypt sensitive exports
   gsecutil export --with-values | gpg --encrypt > backup.csv.gpg
   ```

4. **Use metadata-only exports for documentation**
   ```bash
   gsecutil export -o inventory.csv  # Safe to share
   ```

### Workflow

1. **Maintain consistent labels**
   ```csv
   name,value,label:env,label:team,label:app
   secret1,val1,production,backend,myapp
   ```

2. **Use `--update-config` to keep metadata in sync**
   ```bash
   gsecutil import secrets.csv --upsert --update-config
   ```

3. **Version your metadata CSVs**
   ```bash
   # Safe to commit (no values)
   gsecutil export -o docs/secrets-$(date +%Y%m%d).csv
   ```

### Excel Tips

1. **Multi-line cells**: Use `Alt+Enter` (Windows) or `Control+Option+Enter` (Mac)

2. **Save format**: Use "CSV UTF-8" for best compatibility

3. **Avoid special characters in column names**: Stick to alphanumeric and underscore

4. **Quote cells with commas/newlines**: Excel does this automatically

---

## Troubleshooting

### Import Errors

#### "CSV must have 'name' column"

**Cause**: CSV is missing required `name` column.

**Fix**: Ensure first row has `name` column (case-insensitive).

```csv
name,value,title
secret1,value1,Title 1
```

---

#### "CSV header contains duplicate column names"

**Cause**: Multiple columns with same name (case-insensitive).

**Fix**: Remove or rename duplicate columns.

```csv
# ❌ Bad
name,value,owner,Owner
# ✓ Good
name,value,owner,created_by
```

---

#### "CSV header contains empty column name"

**Cause**: Unnamed column in CSV.

**Fix**: Name all columns or remove empty ones.

---

#### "Row X has Y columns, expected Z"

**Cause**: Inconsistent column count.

**Fix**: Ensure all rows have same number of columns. Check for unquoted commas.

---

### Multi-line Value Issues

#### Multi-line values not importing correctly

**Symptoms**: Line breaks lost or CSV parsing fails.

**Fix**: 
1. Ensure cells with line breaks are quoted
2. Use Excel "CSV UTF-8" format
3. Verify file encoding

```bash
# Check file encoding
file -I secrets.csv
# Should show: text/csv; charset=utf-8
```

---

### Configuration Issues

#### "Config file was not created"

**Cause**: Permission issues or invalid path.

**Fix**: Check config directory permissions.

```bash
# Ensure config directory exists and is writable
mkdir -p ~/.config/gsecutil
chmod 755 ~/.config/gsecutil
```

---

#### Config not updated after import

**Cause**: Missing `--update-config` flag.

**Fix**: Add the flag.

```bash
gsecutil import secrets.csv --upsert --update-config
```

---

### GCP Issues

#### "Secrets not created in GCP"

**Troubleshooting:**

1. **Check project ID**
   ```bash
   gsecutil config show
   ```

2. **Verify authentication**
   ```bash
   gcloud auth list
   ```

3. **Ensure Secret Manager API is enabled**
   ```bash
   gcloud services enable secretmanager.googleapis.com
   ```

4. **Check IAM permissions**
   ```bash
   gcloud projects get-iam-policy PROJECT_ID
   ```

---

## Sample Data

Sample CSV files are available in `examples/csv/`:

- **`sample-import.csv`** - Complete example with all columns
- **`sample-multiline.csv`** - Multi-line values (certificates, configs)
- **`sample-metadata-only.csv`** - Metadata without values

See [Examples](../examples/README.md) for detailed descriptions and usage examples.

---

## Related Documentation

- [Command Reference](commands.md) - Full command documentation
- [Configuration Guide](configuration.md) - Config file format
- [Troubleshooting](troubleshooting.md) - General troubleshooting
- [Examples](../examples/README.md) - Sample files and usage examples

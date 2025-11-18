# gsecutil Configuration

gsecutil supports configuration files to make team secret management easier and more organized. Configuration files allow you to:

- Set default Google Cloud project
- Define secret name prefixes for team organization
- Document credentials with rich metadata
- Share team standards and conventions

## Configuration File Locations

### Default Locations

gsecutil looks for configuration files in the following locations by default:

| Platform | Default Path |
|----------|-------------|
| Linux | `$HOME/.config/gsecutil/gsecutil.conf` |
| macOS | `$HOME/.config/gsecutil/gsecutil.conf` |
| Windows | `%USERPROFILE%\.config\gsecutil\gsecutil.conf` |

### Custom Location

You can specify a custom configuration file using the `--config` flag:

```bash
gsecutil --config /path/to/custom/gsecutil.conf list
gsecutil --config C:\team\secrets\gsecutil.conf get my-secret
```

## Configuration Priority

Settings are resolved in the following order (highest to lowest priority):

1. **Command line parameters** - `gsecutil --project my-project get secret`
2. **Configuration file** - Settings in `gsecutil.conf`
3. **Environment variables** - `GSECUTIL_PROJECT`
4. **gcloud CLI default** - Output of `gcloud config get-value project`

## Configuration Format

Configuration files use YAML format and support the following sections:

### Basic Configuration

```yaml
# Google Cloud Project (required)
project: "my-team-project-123"

# Secret name prefix (optional but recommended)
# Only secrets with this prefix will be managed
prefix: "team-shared-"

# List command configuration
list:
  # Attributes to display in list output by default
  attributes:
    - title
    - owner
    - environment
```

### Credential Documentation

```yaml
credentials:
  - name: "team-shared-db-prod"
    title: "Production Database"
    description: "Main application database credentials"
    environment: "production"
    owner: "backend-team"
    rotation_schedule: "quarterly"
    # Add any custom attributes your team needs
    database_type: "postgresql"
    contact: "backend@company.com"
```

## How Prefix Filtering Works

When a `prefix` is specified in the configuration:

### ✅ Commands that respect prefix filtering:
- `gsecutil list` - Only shows secrets with the specified prefix
- `gsecutil get <secret>` - Works with both prefixed and full names
- `gsecutil create <secret>` - Automatically adds prefix if not provided
- `gsecutil update <secret>` - Works with prefixed secrets
- `gsecutil delete <secret>` - Works with prefixed secrets
- `gsecutil describe <secret>` - Works with prefixed secrets
- `gsecutil access <command> <secret>` - Works with prefixed secrets

### ⚠️ Commands that ignore prefix filtering:
- `gsecutil auditlog` - Shows audit logs for all secrets (security requirement)

### 📋 Commands that show config attributes:
- `gsecutil describe <secret>` - Shows all attributes defined in config file for the secret
- `gsecutil list` - Shows attributes based on `list.attributes` config or `--show-attributes` parameter

### Examples with prefix "team-shared-":

```bash
# These commands are equivalent when prefix is configured:
gsecutil get db-prod
gsecutil get team-shared-db-prod

# Creating secrets automatically adds prefix:
gsecutil create api-key  # Creates "team-shared-api-key"

# List only shows team secrets:
gsecutil list  # Only shows secrets starting with "team-shared-"
```

## List Command Configuration

The `list` command can be customized in two ways:

### 1. Attribute Display

Control which attributes are shown in the list output:

```bash
# Show specific attributes (overrides config file)
gsecutil list --show-attributes title,owner,environment

# Show only title and description
gsecutil list --show-attributes title,description

# Show many attributes
gsecutil list --show-attributes title,description,owner,environment,sensitive_level
```

**Default behavior:**
- If config file has `credentials` section: shows `title` by default
- If no config file or no `credentials`: shows only secret names
- Config file `list.attributes` section overrides default
- `--show-attributes` parameter overrides everything

## Describe Command Integration

The `describe` command automatically shows all attributes defined in the configuration file:

```bash
# Shows Secret Manager metadata + all config attributes
gsecutil describe team-db-prod

# Example output:
# Name: projects/my-project/secrets/team-db-prod
# Created: 2025-01-15T10:30:00Z
# Labels:
#   managed_by: gsecutil
#   team: backend
#
# Config Attributes:
#   Title: Production Database
#   Description: PostgreSQL master database
#   Environment: production
#   Owner: backend-team
#   Contact: backend-team@company.com
#   Rotation Schedule: quarterly
#   Sensitive Level: critical
```

**Describe behavior:**
- Shows standard Secret Manager metadata (name, created time, labels, etc.)
- Shows **all** attributes defined for the secret in config file
- If secret not in config: shows only Secret Manager metadata
- No CLI parameters needed - always shows full attribute set

### 2. Attribute Filtering

Filter secrets based on attributes defined in your configuration:

```bash
# Filter by environment (only shows secrets with this attribute value)
gsecutil list --filter-attributes environment=production

# Filter by owner
gsecutil list --filter-attributes owner=backend-team

# Filter by multiple attributes
gsecutil list --filter-attributes environment=production,sensitive_level=critical

# Combine filtering with custom display
gsecutil list --filter-attributes environment=production --show-attributes title,owner
```

## Common Attributes

While you can define any custom attributes, here are some commonly used ones:

- `title` - Human-readable name
- `description` - What this credential is for
- `environment` - production, staging, development
- `owner` - Team or person responsible
- `contact` - Email or Slack channel for questions
- `rotation_schedule` - How often to rotate (monthly, quarterly, etc.)
- `sensitive_level` - low, medium, high, critical
- `category` - database, api_key, service_account, etc.
- `vendor` - Third-party service name
- `compliance_requirements` - PCI-DSS, SOX, HIPAA, etc.

## Team Workflow

### 1. Create Team Configuration

Create a `gsecutil.conf` file for your team:

```yaml
project: "your-team-project"
prefix: "team-"
credentials:
  - name: "team-db-prod"
    title: "Production Database"
    description: "Main application database"
    owner: "backend-team"
```

### 2. Share with Team

Commit the configuration file to your team's repository or share it through your preferred method. Team members can either:

- Place it in their default config location
- Use `--config` flag to specify the location
- Set up a team alias: `alias team-gsecutil="gsecutil --config /team/gsecutil.conf"`

### 3. Use Consistently

All team members now have:
- Same project configuration
- Consistent secret naming (with prefix)
- Shared documentation of what each secret is for
- Ability to filter and organize secrets
- Rich documentation via `describe` command
- Consistent attribute display in `list` command

## Example Workflows

### Small Team Setup

```yaml
project: "startup-secrets"
prefix: "app-"
credentials:
  - name: "app-db"
    title: "Database Password"
    owner: "dev-team"
  - name: "app-stripe"
    title: "Stripe API Key"
    owner: "dev-team"
```

### Multi-Environment Team

```yaml
project: "company-prod"
prefix: "myteam-"
credentials:
  - name: "myteam-db-prod"
    title: "Production DB"
    environment: "production"
    owner: "backend"
  - name: "myteam-db-staging"
    title: "Staging DB"
    environment: "staging"
    owner: "backend"
```

### Enterprise Setup

```yaml
project: "enterprise-secrets"
prefix: "platform-"
defaults:
  labels:
    managed_by: "gsecutil"
    team: "platform"
credentials:
  - name: "platform-k8s-prod"
    title: "Kubernetes Service Account"
    environment: "production"
    owner: "devops-team"
    contact: "devops@company.com"
    sensitive_level: "critical"
    compliance_requirements: ["SOX"]
```

## Security Considerations

- **Configuration files are not encrypted** - Don't put actual secret values in them
- **Use version control carefully** - Configuration files can be committed to repos since they contain no secrets
- **Prefix separation** - Use prefixes to avoid conflicts with other teams' secrets
- **Access control** - Ensure team members have appropriate IAM permissions for the configured project

## Troubleshooting

### Configuration not found
```bash
# Check if config file exists
ls -la ~/.config/gsecutil/gsecutil.conf

# Use custom location
gsecutil --config /path/to/config.conf list
```

### Prefix filtering not working
```bash
# Check current configuration
gsecutil config show

# List all secrets (ignoring prefix)
gsecutil list --no-prefix-filter
```

### Attributes not showing in list output
```bash
# Check if config has credentials section
gsecutil config show

# Force showing attributes
gsecutil list --show-attributes title,owner,environment

# Check what attributes are available for filtering
gsecutil list --show-attributes title --filter-attributes environment=production
```

### Project not detected
```bash
# Check project resolution order
gsecutil config debug

# Override with command line
gsecutil --project my-project list
```

### Config attributes not showing in describe
```bash
# Check if secret is in config file
gsecutil config show | grep -A5 "team-db-prod"

# Verify secret name matches exactly (including prefix)
gsecutil list --show-attributes title | grep "team-db-prod"

# Check if config file is being loaded
gsecutil config debug
```

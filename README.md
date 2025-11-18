# gsecutil - Google Secret Manager Utility

🚀 A simplified command-line wrapper for Google Secret Manager with configuration file support and team-friendly features. `gsecutil` provides convenient commands for common secret operations, making it easier for small teams to manage passwords and credentials using Google Cloud's Secret Manager without needing a dedicated password management tool.

## 🌍 Language Versions

This README is available in multiple languages:

- **English** - [README.md](README.md) (current)
- **日本語** - [README.ja.md](README.ja.md)
- **中文** - [README.zh.md](README.zh.md)
- **Español** - [README.es.md](README.es.md)
- **हिंदी** - [README.hi.md](README.hi.md)
- **Português** - [README.pt.md](README.pt.md)

> **Note**: All non-English versions are machine-translated. For the most accurate and up-to-date information, please refer to the English version.

## 🎯 Who is this for?

`gsecutil` is designed for **small teams** who:

- ✅ Already use Google Cloud services
- ✅ Need to share passwords/credentials among team members
- ✅ Want something better than shared spreadsheets or plain text files
- ✅ Cannot afford dedicated password management tools
- ✅ Are comfortable with command-line tools

### 💰 Cost-effective solution

Many password managers are designed for personal use or require expensive team subscriptions. `gsecutil` leverages Google Cloud's Secret Manager, which offers:

- **Free tier**: Up to 6 active secret versions and 10,000 access operations per month at no cost
- **Pay-as-you-use**: Only pay for what you store and access beyond the free tier (\$0.06 per additional active version)
- **No per-user licensing**: Share access among team members without per-seat costs

### ⚠️ What gsecutil is NOT

- **Not a full-featured password manager** - lacks browser extensions, auto-fill, password generation, etc.
- **Not a replacement for enterprise tools** - missing advanced features like SCIM provisioning, detailed reporting, etc.
- **Not suitable for personal use** - requires Google Cloud setup and is overkill for individual users
- **Not a complete Secret Manager interface** - only covers common use cases, not all Google Secret Manager features

## ✨ What gsecutil offers

### 🔐 **Basic Secret Operations**
- **Simple CRUD commands**: Create, read, update, delete secrets with easy-to-remember commands
- **Version access**: Get specific versions of secrets when needed
- **Cross-platform** support (Linux, macOS, Windows with ARM64 support)
- **Clipboard integration** - copy secret values directly to clipboard for convenience
- **Multiple input methods** - interactive prompts, inline data, or file-based loading

### 🛡️ **Basic Access Management** *(NEW in v1.0.0)*
- **IAM policy viewing** - see who has access to secrets
- **Permission checking** - understand access at secret and project levels
- **Simple access control** - grant/revoke access for users, groups, and service accounts
- **Basic IAM analysis** - identify common access patterns

### 📊 **Audit Log Access**
- **Audit log viewing** - see who accessed secrets and when
- **Basic filtering** - filter by secret name, user, or operation type
- **Simple reporting** - basic audit trail for security review

### 🎯 **Team-Friendly**
- **Consistent commands** - unified approach across all operations
- **Good error messages** - helpful feedback when things go wrong
- **Multiple output formats** - JSON, YAML, or human-readable table output
- **Free to use** - leverages Google Cloud's free tier for small teams

### ⚙️ **Configuration & Team Management** *(NEW in v1.1.0)*
- **Configuration file support** - YAML configuration file for team settings and metadata
- **Prefix functionality** - transparent prefix handling for team secret organization
- **Enhanced list command** - display custom attributes and metadata from config
- **Enhanced describe command** - show config attributes alongside Secret Manager metadata
- **Team metadata** - document secret ownership, environment, descriptions, and custom attributes
- **Flexible configuration** - works with or without config file (fully backward compatible)

## Prerequisites

- [Google Cloud SDK (gcloud)](https://cloud.google.com/sdk/docs/install) installed and authenticated
- Google Cloud project with Secret Manager API enabled
- Appropriate IAM permissions for Secret Manager operations

## Installation

### Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/superdaigo/gsecutil/releases):

| Platform | Architecture | Download |
|----------|--------------|----------|
| Linux | x64 | `gsecutil-linux-amd64-v{version}` |
| Linux | ARM64 | `gsecutil-linux-arm64-v{version}` |
| macOS | Intel | `gsecutil-darwin-amd64-v{version}` |
| macOS | Apple Silicon | `gsecutil-darwin-arm64-v{version}` |
| Windows | x64 | `gsecutil-windows-amd64-v{version}.exe` |

**After Download:** Rename the binary for consistent usage:

```bash
# Linux/macOS example:
mv gsecutil-linux-amd64-v1.1.0 gsecutil
chmod +x gsecutil

# Windows example (PowerShell/Command Prompt):
ren gsecutil-windows-amd64-v1.1.0.exe gsecutil.exe
```

This allows you to use `gsecutil` consistently regardless of version.

### Install with Go

```bash
go install github.com/superdaigo/gsecutil@latest
```

### Build from Source

For comprehensive build instructions, see [BUILD.md](BUILD.md).

**Quick build:**
```bash
git clone https://github.com/superdaigo/gsecutil.git
cd gsecutil

# Build for current platform
make build
# OR
./build.sh          # Linux/macOS
.\build.ps1         # Windows

# Build for all platforms
make build-all
# OR
./build.sh all      # Linux/macOS
.\build.ps1 all     # Windows
```

## Usage

### Global Options

- `-p, --project`: Google Cloud project ID (can also be set via `GSECUTIL_PROJECT` environment variable)
- `--config`: Path to configuration file (default: `~/.config/gsecutil/gsecutil.conf`)

### Commands

#### Get Secret

Retrieve a secret value from Google Secret Manager. By default, gets the latest version, but you can specify any version:

```bash
# Get latest version of a secret
gsecutil get my-secret

# Get specific version (useful for rollbacks, debugging, or accessing historical values)
gsecutil get my-secret --version 1
gsecutil get my-secret -v 3

# Get secret and copy to clipboard
gsecutil get my-secret --clipboard

# Get specific version with clipboard
gsecutil get my-secret --version 2 --clipboard

# Get secret with version metadata (version, created time, state)
gsecutil get my-secret --show-metadata
gsecutil get my-secret -v 1 --show-metadata    # Older version with metadata

# Get secret from specific project
gsecutil get my-secret --project my-gcp-project
```

**Version Support:**
- 🔄 **Latest version**: Default behavior when no `--version` is specified
- 📅 **Historical versions**: Access any previous version by number (e.g., `--version 1`, `--version 2`)
- 🔍 **Version metadata**: Use `--show-metadata` to see version details (creation time, state, ETag)
- 📋 **Clipboard support**: Works with any version using `--clipboard`

#### Create Secret

Create a new secret in Google Secret Manager:

```bash
# Create secret interactively (secure prompt)
gsecutil create my-secret

# Create secret with inline data
gsecutil create my-secret --data "my-secret-value"

# Create secret from file
gsecutil create my-secret --data-file ./secret.txt

# Create secret with labels
gsecutil create my-secret --labels env=prod,team=backend

# Force creation without version management (may exceed free tier)
gsecutil create my-secret --data "value" --force
```

**🆓 Free Tier Version Management:**

Google Cloud Secret Manager's free tier allows up to 6 active secret versions. When creating a secret that already exists (adding a new version), `gsecutil` will:

1. **Check version count**: If the secret already has 6 active versions
2. **Show current versions**: Display all active versions with creation dates
3. **Offer choices**:
   - **Disable old versions** to stay within free tier (recommended)
   - **Proceed anyway** keeping all versions (may incur charges)
4. **Protect default version**: Never disables the current default/latest version
5. **Bypass with --force**: Skip the check entirely if you know what you're doing

#### Update Secret

Update an existing secret by creating a new version:

```bash
# Update secret interactively
gsecutil update my-secret

# Update with inline data
gsecutil update my-secret --data "new-secret-value"

# Update from file
gsecutil update my-secret --data-file ./new-secret.txt

# Force update without version management (may exceed free tier)
gsecutil update my-secret --data "new-value" --force
```

**🆓 Free Tier Version Management:**

Before updating a secret, `gsecutil` automatically checks if adding a new version would exceed the free tier limit of 6 active versions. If so, you'll see:

```
Secret 'my-secret' currently has 6 active versions.
The Google Cloud Secret Manager free tier allows up to 6 active versions.
Adding a new version would exceed this limit.

Current active versions (oldest first):
  Version 1 - Created: 2023-01-01T10:00:00Z
  Version 2 - Created: 2023-01-02T10:00:00Z (default - will be preserved)
  [...]

Choose an option:
  y/yes: Disable old versions to stay within free tier (recommended)
  N/no:  Proceed anyway, keeping all versions (may incur charges)
Your choice (y/N):
```

#### Delete Secret

Delete a secret permanently:

```bash
# Delete with confirmation prompt
gsecutil delete my-secret

# Force delete without confirmation
gsecutil delete my-secret --force
```

#### List Secrets

List all secrets in a project with enhanced configuration support:

```bash
# List all secrets
gsecutil list

# List with config attributes (shows title, owner, environment from config file)
gsecutil list --show-attributes "title,owner,environment"

# Filter by configuration attributes
gsecutil list --filter-attributes "environment=production,owner=backend-team"

# List with custom format
gsecutil list --format json

# List with Secret Manager label filter
gsecutil list --filter "labels.env=prod"

# List with limit
gsecutil list --limit 10

# List secrets accessible by a specific principal
gsecutil list --principal user:alice@example.com
```

**Configuration Features** *(NEW in v1.1.0)*:
- **--show-attributes**: Display custom attributes from configuration file (inserted after NAME, before built-in fields like LABELS and CREATED)
- **--filter-attributes**: Filter secrets by configuration attributes
- **Automatic attribute display**: Shows default attributes when configuration file is present
- **Prefix filtering**: Automatically filters by configured prefix
- **Built-in field preservation**: Secret Manager built-in fields (LABELS, CREATED) are always shown alongside custom attributes

#### Describe Secret

Get information about a secret including metadata, labels, default version, and other details:

```bash
# Describe secret with basic information
gsecutil describe my-secret

# Describe with information about all versions
gsecutil describe my-secret --show-versions

# Describe with JSON output (raw gcloud format)
gsecutil describe my-secret --format json
```

**Information Displayed:**
- Basic metadata (name, creation time, ETag)
- **Config attributes** *(NEW in v1.1.0)* - title, owner, environment, and custom attributes from configuration file
- **Default version** (current active version, state, creation time)
- **Replication strategy** (automatic multi-region or user-managed)
- **Labels** (sorted alphabetically for organization)
- **Tags/Annotations** (additional metadata, sorted alphabetically)
- Version aliases (if configured)
- Expiration and rotation settings (if configured)
- Pub/Sub topic integrations (if configured)

#### Audit Log

View audit log entries for secrets to see who accessed them, when, and what operations were performed.
Supports flexible filtering with partial matching for both secret names and usernames:

```bash
# Show all Secret Manager audit logs
gsecutil auditlog

# Show logs for secrets matching a pattern (partial match)
gsecutil auditlog my-secret

# Show logs for a specific user (partial match)
gsecutil auditlog --principal john

# Filter by specific operations (single or multiple)
gsecutil auditlog --operation ACCESS
gsecutil auditlog --operation ACCESS,CREATE,DELETE

# Combine all filters for precise results
gsecutil auditlog db --principal admin --operation GET_METADATA,LIST

# Show audit log for the last 30 days
gsecutil auditlog my-secret --days 30

# Show audit log with JSON output
gsecutil auditlog my-secret --format json

# Limit the number of entries returned
gsecutil auditlog --limit 20
```

**Key Features:**
- **Optional secret name**: Show all Secret Manager logs when no secret name is provided
- **Partial matching**: Both secret names and principals support partial/substring matching
- **🔍 Operations filtering**: Filter by specific operation types (ACCESS, CREATE, UPDATE, DELETE, etc.)
- **Flexible filtering**: Combine secret name, principal, and operation filters for precise results
- **Case-insensitive**: All partial matching is case-insensitive
- **Multiple operations**: Specify multiple operations separated by commas

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

**Note**: The `auditlog` command requires Data Access audit logs to be enabled for the Secret Manager API. See [docs/audit-logging.md](docs/audit-logging.md) for detailed setup instructions.

#### Access Management
*(Introduced in v1.0.0)*

Basic IAM access management for secrets:

```bash
# List principals with access to a secret
gsecutil access list my-secret

# Include project-level permissions in the check
gsecutil access list my-secret --include-project

# Grant access to a user
gsecutil access grant my-secret --principal user:alice@example.com

# Grant specific role to a service account
gsecutil access grant my-secret \
  --principal serviceAccount:app@project.iam.gserviceaccount.com \
  --role roles/secretmanager.viewer

# Revoke access from a user
gsecutil access revoke my-secret --principal user:bob@example.com

# Show project-level Secret Manager permissions
gsecutil access project

# List secrets a user can access
gsecutil list --principal user:alice@example.com
```

**Access Management Features:**
- **Basic policy checking**: Shows who has access to secrets
- **IAM condition display**: Shows conditional access rules (when present)
- **Common principal types**: Users, groups, service accounts
- **Simple role management**: Grant/revoke access with common Secret Manager roles
- **Project-level awareness**: Shows broader permissions that might affect secrets

**Available Roles:**
- `roles/secretmanager.secretAccessor` - Can access secret values (default)
- `roles/secretmanager.viewer` - Can view metadata only
- `roles/secretmanager.admin` - Full access to secrets
- `roles/secretmanager.secretVersionManager` - Can create/destroy versions
- `roles/secretmanager.secretVersionAdder` - Can add new versions

### Examples

#### Basic Usage

```bash
# Set your project (optional - can be done via flag each time)
export GSECUTIL_PROJECT=my-project-id

# Create a new secret
gsecutil create database-password
# Enter secret value: [hidden input]

# Retrieve the secret
gsecutil get database-password

# Copy secret to clipboard for pasting elsewhere
gsecutil get database-password --clipboard

# Get secret with version and timestamp information
gsecutil get database-password --show-metadata

# Update the secret
gsecutil update database-password --data "new-password"

# List all secrets
gsecutil list

# Get secret metadata with version history
gsecutil describe database-password --show-versions

# View audit log to see who accessed the secret
gsecutil audit database-password
```

### Example Output

#### Get secret with metadata

```bash
$ gsecutil get my-api-key --show-metadata
Secret: my-api-key
Version: projects/my-project/secrets/my-api-key/versions/3
State: ENABLED
Created: 2025-09-04T12:35:42Z
ETag: "abc123def456"
---
Secret Value: sk-1234567890abcdef
```

#### Describe secret with version history

```bash
$ gsecutil describe my-api-key --show-versions
Name: projects/my-project/secrets/my-api-key
Created: 2025-09-01T10:00:00Z
ETag: "xyz789"
Labels:
  env: production
  team: backend

--- Versions ---

Version: 3
  State: ENABLED
  Created: 2025-09-04T12:35:42Z
  ETag: "abc123def456"

Version: 2
  State: DISABLED
  Created: 2025-09-02T14:20:15Z
  ETag: "def456ghi789"

Version: 1
  State: DESTROYED
  Created: 2025-09-01T10:00:00Z
  Destroy Time: 2025-09-03T09:15:30Z
  ETag: "ghi789jkl012"
```

#### Audit log output

```bash
$ gsecutil audit my-api-key
Audit log entries for secret 'my-api-key' (last 7 days):

TIMESTAMP            OPERATION                      USER                                     RESOURCE
------------------------------------------------------------------------------------------------------------------------
2025-09-04 14:30:15  ACCESS                         user@company.com                        .../secrets/my-api-key/versions/3
2025-09-04 12:35:42  CREATE                         user@company.com                        .../secrets/my-api-key
2025-09-03 09:15:30  UPDATE                         service-account@project.iam.gserviceac  .../secrets/my-api-key
2025-09-02 16:20:10  GET_METADATA                   admin@company.com                        .../secrets/my-api-key

Total entries: 4
```

#### Advanced Usage

```bash
# Create secret from environment variable
echo "$DB_PASSWORD" | gsecutil create db-password --data-file -

# Create secret with metadata
gsecutil create api-key \
  --data "sk-1234567890" \
  --labels env=production,service=api,team=backend

# List production secrets only
gsecutil list --filter "labels.env=production"

# Get secret for specific version
gsecutil get api-key --version 3

# Bulk operations using shell scripting
for secret in $(gsecutil list --format="value(name)"); do
  echo "Describing $secret:"
  gsecutil describe "$secret"
done

# Audit multiple secrets for security review
for secret in critical-api-key database-password; do
  echo "Audit log for $secret:"
  gsecutil audit "$secret" --days 30
  echo
done
```

#### CI/CD Integration

```bash
# In CI/CD pipelines, use --force to avoid prompts
gsecutil delete old-secret --force

# Use JSON output for parsing in scripts
SECRET_DATA=$(gsecutil get my-secret --format json | jq -r .data)

# Create secrets from files in deployment
gsecutil create app-config --data-file config.json
```

## Configuration Files
*(Introduced in v1.1.0)*

`gsecutil` supports YAML configuration files to streamline team workflows and add metadata to secrets. Configuration is completely optional - all existing workflows continue to work without any configuration file.

### Configuration File Location

By default, `gsecutil` looks for a configuration file at:
- **Linux/macOS**: `~/.config/gsecutil/gsecutil.conf`
- **Windows**: `%USERPROFILE%\.config\gsecutil\gsecutil.conf`

You can specify a custom config file with the `--config` flag:

```bash
gsecutil --config /path/to/team-config.conf list
```

### Configuration File Format

```yaml
# gsecutil configuration file
# Google Cloud Project (can be overridden by CLI --project flag)
project: "my-team-project-123"

# Secret name prefix (optional but recommended for teams)
# Automatically added to user input and filtered in list commands
prefix: "team-shared-"

# List command configuration
list:
  # Default attributes to show (can be overridden by --show-attributes)
  attributes:
    - title
    - owner
    - environment

# Team metadata for secrets
credentials:
  - name: "database-password"  # Secret name (without prefix)
    title: "Production Database Password"
    description: "MySQL root password for production database"
    environment: "production"
    owner: "backend-team"
    rotation_schedule: "quarterly"

  - name: "api-key"
    title: "External API Key"
    description: "Production API key for payment processing"
    environment: "production"
    owner: "api-team"
    sensitive_level: "high"
```

### Configuration Benefits

- **Team Documentation**: Add titles, descriptions, owners, and custom metadata
- **Prefix Management**: Automatically handle team prefixes (e.g., `team-shared-database-password`)
- **Enhanced List Output**: Show custom attributes alongside built-in Secret Manager fields
- **Filtering**: Filter by configuration attributes (`--filter-attributes "environment=prod"`)
- **Shared Configuration**: Commit config files to version control for team sharing

### Usage with Configuration

```bash
# User types short name, gsecutil adds prefix automatically
gsecutil get database-password  # Actually accesses "team-shared-database-password"

# List shows custom attributes from config alongside built-in fields
gsecutil list
# NAME                           TITLE                        OWNER         ENVIRONMENT  LABELS         CREATED
# team-shared-database-password  Production Database Password backend-team  production   env=prod       2023-06-15
# team-shared-api-key           External API Key             api-team      production   env=staging    2023-07-20

# Describe shows config attributes along with Google Secret Manager metadata
gsecutil describe database-password
# Name: projects/my-project/secrets/team-shared-database-password
# Created: 2025-09-04T13:30:11Z
#
# Config Attributes:
#   title: Production Database Password
#   description: MySQL root password for production database
#   environment: production
#   owner: backend-team
#   rotation_schedule: quarterly
#
# Default Version: 4
# ...
```

For complete configuration examples, see the [`examples/`](examples/) directory.

## Configuration

### Environment Variables

- `GSECUTIL_PROJECT`: Default project ID (overridden by `--project` flag)

### Authentication

`gsecutil` uses the same authentication as `gcloud`. Ensure you're authenticated:

```bash
# Authenticate with gcloud
gcloud auth login

# Set default project
gcloud config set project YOUR_PROJECT_ID

# For service accounts (in CI/CD)
gcloud auth activate-service-account --key-file=service-account.json
```

### Shell Completion

`gsecutil` supports shell autocompletion for bash, zsh, fish, and PowerShell. This enables tab completion for commands, flags, and options, making the CLI more user-friendly.

#### Setup Instructions

**Bash:**
```bash
# Temporary (current session only)
source <(gsecutil completion bash)

# Permanent installation (requires bash-completion package)
# System-wide (requires sudo)
sudo gsecutil completion bash > /etc/bash_completion.d/gsecutil

# User-local installation
gsecutil completion bash > ~/.local/share/bash-completion/completions/gsecutil

# Or add to ~/.bashrc for automatic loading
echo 'source <(gsecutil completion bash)' >> ~/.bashrc
```

**Zsh:**
```bash
# Temporary (current session only)
source <(gsecutil completion zsh)

# Permanent installation
gsecutil completion zsh > "${fpath[1]}/_gsecutil"

# Or add to ~/.zshrc for automatic loading
echo 'source <(gsecutil completion zsh)' >> ~/.zshrc
```

**Fish:**
```bash
# Temporary (current session only)
gsecutil completion fish | source

# Permanent installation
gsecutil completion fish > ~/.config/fish/completions/gsecutil.fish
```

**PowerShell:**
```powershell
# Add to PowerShell profile
gsecutil completion powershell | Out-String | Invoke-Expression

# Or save to profile for automatic loading
gsecutil completion powershell >> $PROFILE
```

#### Features

Once installed, shell completion provides:
- **Command completion**: Tab to complete `gsecutil` subcommands (`get`, `create`, `list`, etc.)
- **Flag completion**: Tab to complete flags like `--project`, `--version`, `--clipboard`
- **Smart suggestions**: Context-aware completions based on the current command
- **Help text**: Brief descriptions for commands and flags (where supported)

#### Usage Example

```bash
# Type and press Tab to see available commands
gsecutil <Tab>
# Shows: access, auditlog, completion, create, delete, describe, get, help, list, update

# Type partial command and press Tab to complete
gsecutil des<Tab>
# Completes to: gsecutil describe

# Tab completion works for flags too
gsecutil get my-secret --<Tab>
# Shows: --clipboard, --project, --show-metadata, --version
```

**Note**: You may need to restart your shell or source your shell configuration file for completion to take effect.

## Security & Best Practices

### Security Features

- **No persistent storage**: Secret values are never logged or stored by `gsecutil`
- **Secure input**: Interactive prompts use hidden password input
- **OS-native clipboard**: Clipboard operations use secure OS-native APIs
- **gcloud delegation**: All operations delegate to authenticated `gcloud` CLI

### Best Practices

- **Use `--force` carefully**: Always review before using `--force` in automated environments
- **Environment variables**: Set `GSECUTIL_PROJECT` to avoid repetitive `--project` flags
- **Version control**: Use specific secret versions in production (`--version N`)
- **Audit regularly**: Monitor secret access with `gsecutil audit secret-name`
- **Rotate secrets**: Regular secret rotation using `gsecutil update`

## Troubleshooting

### Common Issues

1. **"gcloud command not found"**
   - Ensure Google Cloud SDK is installed and `gcloud` is in your PATH

2. **Authentication errors**
   - Run `gcloud auth login` to authenticate
   - Verify project access: `gcloud config get-value project`

3. **Permission denied errors**
   - Ensure your account has the necessary IAM roles:
     - `roles/secretmanager.admin` (for all operations)
     - `roles/secretmanager.secretAccessor` (for read operations)
     - `roles/secretmanager.secretVersionManager` (for create/update operations)

4. **Clipboard not working**
   - Ensure you have a graphical environment (for Linux)
   - On headless servers, clipboard operations may fail gracefully

### Debug Mode

Add verbose output to gcloud commands by setting:

```bash
export CLOUDSDK_CORE_VERBOSITY=debug
```

## Documentation

- **[BUILD.md](BUILD.md)** - Comprehensive build instructions for all platforms
- **[docs/configuration.md](docs/configuration.md)** - Detailed configuration file reference and examples
- **[WARP.md](WARP.md)** - Development guidance for WARP AI terminal integration
- **README.md** - This file, usage and overview

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Related Projects

- [Google Cloud SDK](https://cloud.google.com/sdk)
- [Secret Manager Documentation](https://cloud.google.com/secret-manager/docs)
- [Cobra CLI Framework](https://github.com/spf13/cobra)

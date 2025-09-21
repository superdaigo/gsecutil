# gsecutil - Google Secret Manager Utility

A command-line utility that provides a simple wrapper around the `gcloud` CLI for managing Google Secret Manager secrets. `gsecutil` offers simplified commands for getting, creating, updating, deleting, listing, and describing secrets, with the added convenience of copying secret values directly to your clipboard.

## Features

- **Simple wrapper** around `gcloud` CLI for Google Secret Manager operations
- **Cross-platform** support (Linux, macOS, Windows)
- **Clipboard integration** - copy secret values directly to clipboard
- **Version metadata** - show version numbers, creation times, and states
- **Version history** - view all versions with detailed timestamps
- **Audit logging** - view who accessed secrets, when, and what operations were performed
- **Interactive secret input** with hidden password prompts
- **File-based secret input** for loading secrets from files
- **Comprehensive command set**: get, create, update, delete, list, describe, auditlog
- **Flexible output formatting** (JSON, YAML, table)
- **Project-aware** with global project flag support

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
mv gsecutil-linux-amd64-v0.3.0 gsecutil
chmod +x gsecutil

# Windows example (PowerShell/Command Prompt):
ren gsecutil-windows-amd64-v0.3.0.exe gsecutil.exe
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

- `-p, --project`: Google Cloud project ID (can also be set via `GOOGLE_CLOUD_PROJECT` environment variable)

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
```

#### Update Secret

Update an existing secret by creating a new version:

```bash
# Update secret interactively
gsecutil update my-secret

# Update with inline data
gsecutil update my-secret --data "new-secret-value"

# Update from file
gsecutil update my-secret --data-file ./new-secret.txt
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

List all secrets in a project:

```bash
# List all secrets
gsecutil list

# List with custom format
gsecutil list --format json

# List with filter
gsecutil list --filter "labels.env=prod"

# List with limit
gsecutil list --limit 10
```

#### Describe Secret

Get comprehensive information about a secret including metadata, labels, tags, default version, and replication strategy:

```bash
# Describe secret with comprehensive information
gsecutil describe my-secret

# Describe with detailed information about all versions
gsecutil describe my-secret --show-versions

# Describe with JSON output (raw gcloud format)
gsecutil describe my-secret --format json
```

**Enhanced Information Displayed:**
- Basic metadata (name, creation time, ETag)
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
gsecutil auditlog --user john

# Filter by specific operations (single or multiple)
gsecutil auditlog --operations ACCESS
gsecutil auditlog --operations ACCESS,CREATE,DELETE

# Combine all filters for precise results
gsecutil auditlog db --user admin --operations GET_METADATA,LIST

# Show audit log for the last 30 days
gsecutil auditlog my-secret --days 30

# Show audit log with JSON output
gsecutil auditlog my-secret --format json

# Limit the number of entries returned
gsecutil auditlog --limit 20
```

**Key Features:**
- **Optional secret name**: Show all Secret Manager logs when no secret name is provided
- **Partial matching**: Both secret names and usernames support partial/substring matching
- **🔍 Operations filtering**: Filter by specific operation types (ACCESS, CREATE, UPDATE, DELETE, etc.)
- **Flexible filtering**: Combine secret name, user, and operations filters for precise results
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

### Examples

#### Basic Usage

```bash
# Set your project (optional - can be done via flag each time)
export GOOGLE_CLOUD_PROJECT=my-project-id

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

## Configuration

### Environment Variables

- `GOOGLE_CLOUD_PROJECT`: Default project ID (overridden by `--project` flag)

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

## Building from Source

For comprehensive build instructions including cross-platform builds, CI/CD integration, and troubleshooting, see **[BUILD.md](BUILD.md)**.

### Quick Development Setup

```bash
# Clone and build
git clone https://github.com/superdaigo/gsecutil.git
cd gsecutil

# Install dependencies and build
make deps && make build

# Run tests and validation
make test && make vet && make fmt

# Quick development build and test
make dev
```

### Build Methods

1. **Makefile** (Linux/macOS/WSL): `make build`
2. **Bash Script** (Linux/macOS): `./build.sh`
3. **PowerShell Script** (Windows): `.\build.ps1`
4. **Manual Go**: `go build -o build/gsecutil .`

## System Requirements

### Runtime Requirements

- **Google Cloud SDK**: `gcloud` CLI installed and in PATH
- **Authentication**: Google Cloud SDK authenticated (`gcloud auth login`)
- **API Access**: Secret Manager API enabled in your Google Cloud project
- **Permissions**: Appropriate IAM roles (see [Troubleshooting](#troubleshooting))

### Build Requirements

- **Go**: Version 1.21 or later
- **Make**: Optional, for using Makefile (Linux/macOS/WSL)
- **Git**: For cloning the repository

See [BUILD.md](BUILD.md) for detailed build instructions.

## Security & Best Practices

### Security Features

- **No persistent storage**: Secret values are never logged or stored by `gsecutil`
- **Secure input**: Interactive prompts use hidden password input
- **OS-native clipboard**: Clipboard operations use secure OS-native APIs
- **gcloud delegation**: All operations delegate to authenticated `gcloud` CLI

### Best Practices

- **Use `--force` carefully**: Always review before using `--force` in automated environments
- **Environment variables**: Set `GOOGLE_CLOUD_PROJECT` to avoid repetitive `--project` flags
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
- **[WARP.md](WARP.md)** - Development guidance for WARP AI terminal integration
- **README.md** - This file, usage and overview

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. **Install pre-commit hooks** (see below)
4. Make your changes
5. Add tests if applicable
6. Commit your changes: `git commit -am 'Add amazing feature'` (pre-commit hooks will run automatically)
7. Push to the branch: `git push origin feature/amazing-feature`
8. Submit a pull request

### Pre-commit Hooks

This project uses [pre-commit](https://pre-commit.com/) to ensure code quality. The hooks automatically run on every commit:

**Setup:**
```bash
# Install pre-commit (Ubuntu/Debian)
sudo apt install pre-commit

# Install hooks in your local repository
pre-commit install

# Run hooks manually on all files (optional)
pre-commit run --all-files
```

**What the hooks do:**
- **Format code**: `go fmt`
- **Static analysis**: `go vet`
- **Dependency management**: `go mod tidy`
- **Run tests**: `go test ./cmd`
- **File checks**: Remove trailing whitespace, fix line endings, etc.

The hooks will run automatically on `git commit`. If any hook fails, the commit will be blocked until issues are fixed.

### Development Workflow

```bash
# Setup development environment
make deps

# Install pre-commit hooks (first time only)
pre-commit install

# Make your changes, then:
make fmt vet test     # Manual quality checks (optional - pre-commit handles this)
make dev             # Quick build and test
```

See [BUILD.md](BUILD.md) for detailed development and build instructions.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Related Projects

- [Google Cloud SDK](https://cloud.google.com/sdk)
- [Secret Manager Documentation](https://cloud.google.com/secret-manager/docs)
- [Cobra CLI Framework](https://github.com/spf13/cobra)

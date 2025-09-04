# gsecutil - Google Secret Manager Utility

A command-line utility that provides a simple wrapper around the `gcloud` CLI for managing Google Secret Manager secrets. `gsecutil` offers simplified commands for getting, creating, updating, deleting, listing, and describing secrets, with the added convenience of copying secret values directly to your clipboard.

## Features

- **Simple wrapper** around `gcloud` CLI for Google Secret Manager operations
- **Cross-platform** support (Linux, macOS, Windows)
- **Clipboard integration** - copy secret values directly to clipboard
- **Version metadata** - show version numbers, creation times, and states
- **Version history** - view all versions with detailed timestamps
- **Interactive secret input** with hidden password prompts
- **File-based secret input** for loading secrets from files
- **Comprehensive command set**: get, create, update, delete, list, describe
- **Flexible output formatting** (JSON, YAML, table)
- **Project-aware** with global project flag support

## Prerequisites

- [Google Cloud SDK (gcloud)](https://cloud.google.com/sdk/docs/install) installed and authenticated
- Google Cloud project with Secret Manager API enabled
- Appropriate IAM permissions for Secret Manager operations

## Installation

### Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/daigo/gsecutil/releases):

- **Linux (amd64)**: `gsecutil-linux-amd64`
- **Linux (arm64)**: `gsecutil-linux-arm64`
- **macOS (Intel)**: `gsecutil-darwin-amd64`
- **macOS (Apple Silicon)**: `gsecutil-darwin-arm64`
- **Windows**: `gsecutil-windows-amd64.exe`

### Build from Source

```bash
git clone https://github.com/daigo/gsecutil.git
cd gsecutil
make build
# or for all platforms:
make build-all
```

### Install with Go

```bash
go install github.com/daigo/gsecutil@latest
```

## Usage

### Global Options

- `-p, --project`: Google Cloud project ID (can also be set via `GOOGLE_CLOUD_PROJECT` environment variable)

### Commands

#### Get Secret

Retrieve a secret value from Google Secret Manager:

```bash
# Get latest version of a secret
gsecutil get my-secret

# Get specific version
gsecutil get my-secret --version 2

# Get secret and copy to clipboard
gsecutil get my-secret --clipboard

# Get secret with version metadata (version, created time, state)
gsecutil get my-secret --show-metadata

# Get secret from specific project
gsecutil get my-secret --project my-gcp-project
```

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

Get detailed information about a secret:

```bash
# Describe secret
gsecutil describe my-secret

# Describe with detailed version information
gsecutil describe my-secret --show-versions

# Describe with JSON output
gsecutil describe my-secret --format json
```

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

## Building

### Development

```bash
# Install dependencies
make deps

# Build for current platform
make build

# Run tests
make test

# Format code
make fmt

# Run development version
make dev
```

### Cross-platform Builds

```bash
# Build for all platforms
make build-all

# Or use the build script
./build.sh
```

### Available Make Targets

- `make all` - Build for all platforms
- `make build` - Build for current platform
- `make build-linux` - Build for Linux amd64
- `make build-windows` - Build for Windows amd64
- `make build-darwin` - Build for macOS amd64
- `make clean` - Clean build directory
- `make test` - Run tests
- `make fmt` - Format code
- `make install` - Install locally

## Requirements

### Runtime Requirements

- `gcloud` CLI installed and in PATH
- Google Cloud SDK authenticated
- Secret Manager API enabled in your Google Cloud project

### Build Requirements

- Go 1.21 or later
- Make (optional, for using Makefile)

## Security Notes

- Secret values are never logged or stored by `gsecutil`
- Interactive prompts use hidden input for security
- Clipboard operations are performed securely using OS-native APIs
- Always use `--force` flag carefully in automated environments

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

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Run `make fmt` and `make test`
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Related Projects

- [Google Cloud SDK](https://cloud.google.com/sdk)
- [Secret Manager Documentation](https://cloud.google.com/secret-manager/docs)
- [Cobra CLI Framework](https://github.com/spf13/cobra)

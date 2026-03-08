# gsecutil - Google Secret Manager Utility

🚀 A simplified command-line wrapper for Google Secret Manager with configuration file support and team-friendly features.

## 🌍 Language Versions

- **English** - [README.md](README.md) (current)
- **日本語** - [README.ja.md](README.ja.md)
- **中文** - [README.zh.md](README.zh.md)
- **Español** - [README.es.md](README.es.md)
- **हिंदी** - [README.hi.md](README.hi.md)
- **Português** - [README.pt.md](README.pt.md)

> **Note**: All non-English versions are machine-translated. For the most accurate information, refer to the English version.

## Quick Start

### Installation

Download the latest binary for your platform from the [releases page](https://github.com/superdaigo/gsecutil/releases):

```bash
# macOS Apple Silicon
curl -L https://github.com/superdaigo/gsecutil/releases/latest/download/gsecutil-darwin-arm64 -o gsecutil
chmod +x gsecutil
sudo mv gsecutil /usr/local/bin/

# macOS Intel
curl -L https://github.com/superdaigo/gsecutil/releases/latest/download/gsecutil-darwin-amd64 -o gsecutil
chmod +x gsecutil
sudo mv gsecutil /usr/local/bin/

# Linux
curl -L https://github.com/superdaigo/gsecutil/releases/latest/download/gsecutil-linux-amd64 -o gsecutil
chmod +x gsecutil
sudo mv gsecutil /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/superdaigo/gsecutil/releases/latest/download/gsecutil-windows-amd64.exe" -OutFile "gsecutil.exe"
# Move to a directory in your PATH, e.g., C:\Windows\System32
Move-Item gsecutil.exe C:\Windows\System32\gsecutil.exe
```

Or install with Go:
```bash
go install github.com/superdaigo/gsecutil@latest
```

### Prerequisites

- [Google Cloud SDK (gcloud)](https://cloud.google.com/sdk/docs/install) installed and authenticated
- Google Cloud project with Secret Manager API enabled

### Authentication

```bash
# Authenticate with gcloud
gcloud auth login

# Set default project
gcloud config set project YOUR_PROJECT_ID

# Or set environment variable
export GSECUTIL_PROJECT=YOUR_PROJECT_ID
```

## Basic Usage

### Create a Secret
```bash
# Interactive input
gsecutil create database-password

# From command line
gsecutil create api-key -d "sk-1234567890"

# From file
gsecutil create config --data-file ./config.json
```

### Get a Secret
```bash
# Get latest version
gsecutil get database-password

# Copy to clipboard
gsecutil get api-key --clipboard

# Get specific version
gsecutil get api-key --version 2
```

### List Secrets
```bash
# List all secrets
gsecutil list

# Filter by label
gsecutil list --filter "labels.env=prod"
```

### Update a Secret
```bash
# Interactive input
gsecutil update database-password

# From command line
gsecutil update api-key -d "new-secret-value"
```

### Delete a Secret
```bash
gsecutil delete old-secret
```

## Configuration

gsecutil supports configuration files for project-specific settings. Config files are searched in this order:

1. `--config` flag (if specified)
2. Current directory: `gsecutil.conf`
3. Home directory: `~/.config/gsecutil/gsecutil.conf`

### Example Configuration

```yaml
# Project ID (optional if set via environment or gcloud)
project: "my-project-id"

# Secret name prefix for team organization
prefix: "team-shared-"

# Default attributes to display in list command
list:
  attributes:
    - title
    - owner
    - environment

# Credential metadata (names are bare — prefix is added automatically)
credentials:
  - name: "database-password"    # accesses "team-shared-database-password"
    title: "Production Database Password"
    environment: "production"
    owner: "backend-team"
```

> **Prefix is transparent:** When a prefix is configured, you always use bare names in commands, config, and CSV files. The prefix is added and stripped automatically.

### Quick Start

```bash
# Generate config interactively
gsecutil config init

# Or create a project-specific config
echo 'project: "my-project-123"' > gsecutil.conf
```

For detailed configuration options, see [docs/configuration.md](docs/configuration.md).

## Key Features

- ✅ **Simple CRUD Operations** - Intuitive commands for managing secrets
- ✅ **Clipboard Integration** - Copy secrets directly to clipboard
- ✅ **Version Management** - Access specific versions and manage version lifecycle
- ✅ **Configuration File Support** - Team-friendly metadata and organization
- ✅ **Access Management** - Basic IAM policy management
- ✅ **Audit Logs** - View who accessed secrets and when
- ✅ **Multiple Input Methods** - Interactive, inline, or file-based
- ✅ **Cross-platform** - Linux, macOS, Windows (amd64/arm64)

## Documentation

- **[Configuration Guide](docs/configuration.md)** - Detailed configuration options and examples
- **[Command Reference](docs/commands.md)** - Complete command documentation
- **[Audit Logging Setup](docs/audit-logging.md)** - Enable and use audit logs
- **[Troubleshooting Guide](docs/troubleshooting.md)** - Common issues and solutions
- **[Build Instructions](BUILD.md)** - Build from source
- **[Development Guide](WARP.md)** - Development with WARP AI

## Common Commands

```bash
# Show secret details
gsecutil describe my-secret

# Show version history
gsecutil describe my-secret --show-versions

# View audit logs
gsecutil auditlog my-secret

# Manage access
gsecutil access list my-secret
gsecutil access grant my-secret --principal user:alice@example.com

# Validate configuration
gsecutil config validate

# Show configuration
gsecutil config show
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Related

- [Google Cloud SDK](https://cloud.google.com/sdk)
- [Secret Manager Documentation](https://cloud.google.com/secret-manager/docs)

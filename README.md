# gsecutil - Google Secret Manager Utility

A simplified command-line wrapper for Google Secret Manager that works like a per-project password manager. Store, retrieve, and manage secrets with intuitive commands, clipboard integration, version control, team-friendly configuration files, and audit logging.

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

Download the latest binary for your platform from the [releases page](https://github.com/superdaigo/gsecutil/releases), or install with Go:

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

Each project typically has its own configuration file that stores the project ID, secret naming conventions, and metadata attributes.

### 1. Create a Configuration File

Run the interactive setup to generate a configuration file. This will prompt you for your Google Cloud project ID, secret name prefix, default list attributes, and optional example credentials. The generated file is saved as `gsecutil.conf` in the current directory by default (use `--home` to save to `~/.config/gsecutil/gsecutil.conf` instead).

```bash
gsecutil config init
```

The configuration file is searched in this order:
1. `--config` flag (if specified)
2. Current directory: `gsecutil.conf`
3. Home directory: `~/.config/gsecutil/gsecutil.conf`

### 2. Manage Secrets

```bash
# Create a secret
gsecutil create database-password

# Get the latest version
gsecutil get database-password

# Copy to clipboard
gsecutil get database-password --clipboard

# List all secrets
gsecutil list

# Update a secret
gsecutil update database-password

# Delete a secret
gsecutil delete database-password
```

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

> **Prefix handling:** When a prefix is configured:
> - **Commands and config files**: Use bare names (prefix is added/stripped automatically)
> - **CSV files**: Use full names including the prefix (for import/export compatibility)

For detailed configuration options, see [docs/configuration.md](docs/configuration.md).

## Documentation

- **[Configuration Guide](docs/configuration.md)** - Detailed configuration options and examples
- **[Command Reference](docs/commands.md)** - Complete command documentation
- **[Audit Logging Setup](docs/audit-logging.md)** - Enable and use audit logs
- **[Troubleshooting Guide](docs/troubleshooting.md)** - Common issues and solutions
- **[Build Instructions](BUILD.md)** - Build from source
- **[Development Guide](WARP.md)** - Development with WARP AI

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Related

- [Google Cloud SDK](https://cloud.google.com/sdk)
- [Secret Manager Documentation](https://cloud.google.com/secret-manager/docs)

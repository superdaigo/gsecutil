# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview

`gsecutil` is a Go CLI application that provides a simplified wrapper around the Google Cloud SDK (`gcloud`) for managing Google Secret Manager secrets. It offers streamlined commands for CRUD operations on secrets with additional features like clipboard integration and enhanced metadata display.

## Development Commands

### Building
- `make build` - Build for current platform (outputs to `build/gsecutil`)
- `make build-all` - Build for all supported platforms (Linux, macOS, Windows for amd64/arm64)
- `./build.sh` - Alternative cross-platform build script with colored output
- `make clean` - Clean build artifacts

### Testing and Code Quality
- `make test` - Run all tests with `go test ./...`
- `make fmt` - Format code with `go fmt ./...`
- `make vet` - Run static analysis with `go vet ./...`
- `make lint` - Run golangci-lint (if installed)

### Development Workflow
- `make deps` - Install/update dependencies with `go mod tidy && go mod download`
- `make dev` - Quick development build and show help
- `make install` - Install locally with version info

### Running Locally
```bash
# Build and test a command
make build
./build/gsecutil --help

# Test with local gcloud setup
./build/gsecutil list
./build/gsecutil get some-secret --clipboard
```

## Architecture

### Command Structure
The application uses the [Cobra CLI framework](https://github.com/spf13/cobra) with the following structure:

- **`main.go`** - Entry point, calls `cmd.Execute()`
- **`cmd/root.go`** - Root command configuration with global `--project` flag
- **`cmd/*.go`** - Individual command implementations (get, create, update, delete, list, describe)

### Key Commands and Their Functions
- **`get`** - Retrieves secret values with optional clipboard integration and metadata display
- **`create`** - Creates new secrets with support for interactive input, inline data, or file input
- **`update`** - Updates existing secrets by creating new versions
- **`delete`** - Deletes secrets with confirmation prompts
- **`list`** - Lists secrets with filtering and formatting options  
- **`describe`** - Shows detailed secret metadata with optional version history
- **`audit`** - Shows audit log entries for secret access, including who accessed secrets, when, and what operations were performed

### Shared Utilities (`cmd/clipboard.go`)
Contains core utility functions used across commands:
- **`copyToClipboard()`** - Cross-platform clipboard operations using `github.com/atotto/clipboard`
- **`getSecretInput()`** - Handles secret input from various sources (CLI args, files, interactive prompts)
- **`getSecretVersionInfo()`** - Fetches version metadata via gcloud JSON output
- **`describeSecretWithVersions()`** - Enhanced secret descriptions with version history
- **`SecretVersionInfo`** and **`SecretInfo`** structs - Data models for gcloud JSON responses

### Audit Functionality (`cmd/audit.go`)
Provides audit log querying capabilities:
- **`AuditLogEntry`** struct - Data model for Google Cloud Logging audit log entries
- **`getOperationName()`** - Converts gcloud method names to human-readable operation names
- Uses `gcloud logging read` to query Cloud Audit Logs for Secret Manager events
- Supports filtering by time range, output formatting, and result limiting

### Integration with gcloud
All secret operations are performed by spawning `gcloud` subprocesses:
- Uses `os/exec.Command()` to call gcloud with appropriate arguments
- Parses JSON output for metadata operations
- Handles authentication and project context through gcloud's existing authentication
- Error handling includes parsing stderr from failed gcloud commands

### Dependencies
- **`github.com/spf13/cobra`** - CLI framework for command structure and flag parsing
- **`github.com/atotto/clipboard`** - Cross-platform clipboard operations
- **`golang.org/x/term`** - Secure password input for interactive secret entry

## Key Patterns

### Error Handling
- gcloud command failures are captured from stderr and returned as formatted errors
- Interactive operations (like clipboard) fail gracefully with warnings rather than errors
- Version metadata failures don't block secret retrieval operations

### Security Practices
- Secret values are never logged or stored permanently
- Interactive secret input uses hidden terminal input (`term.ReadPassword`)
- Secrets are passed to gcloud via stdin rather than command-line arguments
- File-based secret input supports stdin (`-`) for pipeline usage

### Cross-Platform Support
- Build system supports Linux, macOS, and Windows for both amd64 and arm64
- Clipboard operations work across all supported platforms
- Build artifacts use platform-specific naming conventions

## Development Notes

### Adding New Commands
1. Create new file in `cmd/` directory (e.g., `cmd/newcommand.go`)
2. Implement cobra command with `RunE` function
3. Register command in `init()` function with `rootCmd.AddCommand()`
4. Follow existing patterns for gcloud integration and error handling
5. Add appropriate flags using cobra's flag system

### Testing Locally
- Requires working gcloud installation and authentication
- Set `GOOGLE_CLOUD_PROJECT` environment variable or use `--project` flag
- Ensure Secret Manager API is enabled in the target GCP project
- Test cross-platform builds using the build script or make targets

### Extending Metadata Support
- JSON response parsing is handled in `cmd/clipboard.go`
- Add new fields to `SecretVersionInfo` or `SecretInfo` structs as needed
- Update display formatting in relevant command files
- All gcloud JSON output should be parsed rather than using text output for reliability

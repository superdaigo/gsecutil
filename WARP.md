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
- **`auditlog`** - Shows audit log entries for secret access, including who accessed secrets, when, and what operations were performed

### Shared Utilities (`cmd/clipboard.go`)
Contains core utility functions used across commands:
- **`copyToClipboard()`** - Cross-platform clipboard operations using `github.com/atotto/clipboard`
- **`getSecretInput()`** - Handles secret input from various sources (CLI args, files, interactive prompts)
- **`getSecretVersionInfo()`** - Fetches version metadata via gcloud JSON output
- **`describeSecretWithVersions()`** - Enhanced secret descriptions with version history
- **`SecretVersionInfo`** and **`SecretInfo`** structs - Data models for gcloud JSON responses

### Audit Log Functionality (`cmd/auditlog.go`)
Provides audit log querying capabilities through the `auditlog` command:
- **`AuditLogEntry`** struct - Data model for Google Cloud Logging audit log entries
- **`getOperationName()`** - Converts gcloud method names to human-readable operation names
- Uses `gcloud logging read` to query Cloud Audit Logs for Secret Manager events
- Supports filtering by time range, output formatting, and result limiting
- Requires Data Access audit logs to be enabled for comprehensive secret access tracking

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

### Git Workflow

**IMPORTANT: Do not commit or push changes without explicit user instruction.**

When making code changes:
1. Make the changes and verify they work (build, test)
2. Show the changes to the user
3. **Wait for explicit confirmation** before committing
4. Only commit when user says "commit" or "push"
5. Always include `Co-Authored-By: Warp <agent@warp.dev>` in commit messages

The pre-commit hook will automatically run:
- Code formatting (`gofmt`)
- Static analysis (`go vet`)
- All tests (`go test`)

If any check fails, the commit will be aborted.

### Adding New Commands
1. Create new file in `cmd/` directory (e.g., `cmd/newcommand.go`)
2. Implement cobra command with `RunE` function
3. Register command in `init()` function with `rootCmd.AddCommand()`
4. Follow existing patterns for gcloud integration and error handling
5. Add appropriate flags using cobra's flag system

### Testing Locally
- Requires working gcloud installation and authentication
- Set `GSECUTIL_PROJECT` environment variable or use `--project` flag
- Ensure Secret Manager API is enabled in the target GCP project
- Test cross-platform builds using the build script or make targets

### Extending Metadata Support
- JSON response parsing is handled in `cmd/clipboard.go`
- Add new fields to `SecretVersionInfo` or `SecretInfo` structs as needed
- Update display formatting in relevant command files
- All gcloud JSON output should be parsed rather than using text output for reliability

## Version Management and Release Process

### Pre-commit Hooks
The repository includes a pre-commit hook that automatically runs before each commit:
- **Location**: `.git/hooks/pre-commit`
- **Checks performed**:
  1. Code formatting with `gofmt -s -w`
  2. Static analysis with `go vet`
  3. All tests with `go test ./...`
- The hook automatically formats code and stages changes if needed
- To bypass the hook (not recommended): `git commit --no-verify`

### Version Updates
Version numbers are stored in the `VERSION` file at the repository root.

**To update the version:**
1. Edit the `VERSION` file directly:
   ```bash
   echo "1.3.0" > VERSION
   ```
2. Commit the version change:
   ```bash
   git add VERSION
   git commit -m "chore: Bump version to 1.3.0
   
   Co-Authored-By: Warp <agent@warp.dev>"
   git push origin main
   ```

### Creating Releases
Use the `scripts/release.sh` script to create and publish releases. This script:
- Validates the version format
- Checks for uncommitted changes
- Verifies you're on the main branch (optional)
- Runs all tests
- Shows changes since the last release
- Creates an annotated git tag
- Pushes the tag to trigger GitHub Actions release workflow

**Usage:**
```bash
# Interactive mode (prompts for version)
./scripts/release.sh

# Specify version directly
./scripts/release.sh 1.3.0

# Pre-release versions
./scripts/release.sh 2.0.0-beta.1
./scripts/release.sh 1.3.0-rc.1
```

**The release script will:**
1. Check prerequisites (clean working directory, on main branch)
2. Show recent tags and changes since the last release
3. Run all tests to ensure quality
4. Prompt for confirmation
5. Create and push the release tag
6. Trigger GitHub Actions to build binaries for all platforms

**GitHub Actions will automatically:**
- Build binaries for Linux, macOS, Windows (amd64/arm64)
- Generate SHA256 checksums
- Create a GitHub release with all artifacts
- Publish the release automatically

**Monitor the release:**
- Actions: `https://github.com/superdaigo/gsecutil/actions`
- Release page: `https://github.com/superdaigo/gsecutil/releases/tag/vX.Y.Z`

### Release Checklist
Before creating a release:
1. ✅ Update `VERSION` file
2. ✅ Update documentation if needed (README, docs/)
3. ✅ Update all language versions of README if changes are significant
4. ✅ Run full test suite: `make test`
5. ✅ Verify builds work: `make build-all`
6. ✅ Commit and push all changes
7. ✅ Run release script: `./scripts/release.sh X.Y.Z`
8. ✅ Monitor GitHub Actions for successful build
9. ✅ Verify release artifacts are published correctly

## Documentation Guidelines

### Multilingual README Policy

**English README (`README.md`) is the source of truth.**

When updating README.md:
1. **Always update English version first** - `README.md`
2. **Machine translate to other languages** - Replace entire content with machine-translated version
3. **Languages to update:**
   - Japanese: `README.ja.md`
   - Chinese: `README.zh.md`
   - Spanish: `README.es.md`
   - Hindi: `README.hi.md`
   - Portuguese: `README.pt.md`

**Translation approach:**
- Use machine translation (e.g., built-in translation tools, LLM translation)
- Completely replace the content - do NOT attempt to manually edit translations
- Ensure all translated versions include the note: "All non-English versions are machine-translated"
- Keep structure and formatting consistent with English version

**When to update translations:**
- Major feature additions or changes
- Installation instructions changes
- Configuration format changes
- Before creating releases (if README changed)

**When NOT to update:**
- Minor typo fixes in English README
- Small clarifications that don't affect usage
- Internal link updates only

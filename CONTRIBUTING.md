# Contributing to gsecutil

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. **Install pre-commit hooks** (see below)
4. Make your changes
5. Add tests if applicable
6. Commit your changes: `git commit -am 'Add amazing feature'` (pre-commit hooks will run automatically)
7. Push to the branch: `git push origin feature/amazing-feature`
8. Submit a pull request

## Pre-commit Hooks

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

## Development Workflow

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
3. **PowerShell Script** (Windows): `.\\build.ps1`
4. **Manual Go**: `go build -o build/gsecutil .`

## System Requirements

### Runtime Requirements

- **Google Cloud SDK**: `gcloud` CLI installed and in PATH
- **Authentication**: Google Cloud SDK authenticated (`gcloud auth login`)
- **API Access**: Secret Manager API enabled in your Google Cloud project
- **Permissions**: Appropriate IAM roles (see main README for troubleshooting)

### Build Requirements

- **Go**: Version 1.21 or later
- **Make**: Optional, for using Makefile (Linux/macOS/WSL)
- **Git**: For cloning the repository

See [BUILD.md](BUILD.md) for detailed build instructions.

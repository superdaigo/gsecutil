# Building gsecutil

This document provides comprehensive instructions for building `gsecutil` from source across different platforms and environments.

## Prerequisites

### Required Tools

- **Go**: Version 1.21 or later
  - Install from: https://golang.org/dl/
  - Verify: `go version`

- **Git**: For cloning the repository
  - Verify: `git --version`

### Optional Tools

- **Make**: For using Makefile targets (Linux/macOS/WSL)
  - Most Linux distributions: `sudo apt-get install build-essential` or equivalent
  - macOS: Install Xcode Command Line Tools (`xcode-select --install`)

- **golangci-lint**: For code linting (optional)
  - Install: https://golangci-lint.run/usage/install/
  - Verify: `golangci-lint --version`

## Quick Start

```bash
# Clone the repository
git clone https://github.com/yourusername/gsecutil.git
cd gsecutil

# Build for current platform
make build
# OR
./build.sh
# OR (Windows)
.\build.ps1

# Run the built binary
./build/gsecutil --help
```

## Build Methods

### Method 1: Using Makefile (Linux/macOS/WSL)

The Makefile provides the most comprehensive build targets:

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Build for specific platforms
make build-linux           # Linux amd64
make build-linux-arm64     # Linux arm64
make build-windows         # Windows amd64
make build-darwin          # macOS amd64
make build-darwin-arm64    # macOS arm64

# Development tasks
make clean                 # Clean build artifacts
make test                  # Run tests
make fmt                   # Format code
make vet                   # Run static analysis
make lint                  # Run golangci-lint
make deps                  # Install/update dependencies
make install               # Install to GOPATH/bin
make dev                   # Quick build and show help
```

### Method 2: Using Build Scripts

#### Linux/macOS (bash script)

```bash
# Default: build for current platform
./build.sh

# Build for all platforms
./build.sh all

# Build for specific platforms
./build.sh linux
./build.sh linux-arm64
./build.sh windows
./build.sh darwin
./build.sh darwin-arm64

# Development tasks
./build.sh clean
./build.sh test
./build.sh fmt
./build.sh vet
./build.sh lint
./build.sh deps
./build.sh install
./build.sh dev
./build.sh help
```

#### Windows (PowerShell script)

```powershell
# Default: build for current platform
.\build.ps1

# Build for all platforms
.\build.ps1 all

# Build for specific platforms
.\build.ps1 linux
.\build.ps1 linux-arm64
.\build.ps1 windows
.\build.ps1 darwin
.\build.ps1 darwin-arm64

# Development tasks
.\build.ps1 clean
.\build.ps1 test
.\build.ps1 fmt
.\build.ps1 vet
.\build.ps1 lint
.\build.ps1 deps
.\build.ps1 install
.\build.ps1 dev
.\build.ps1 help
```

**Note for Windows**: If you encounter execution policy errors:
```powershell
# Allow script execution for current user
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# Or run directly with bypass
powershell -ExecutionPolicy Bypass -File .\build.ps1
```

### Method 3: Manual Go Commands

```bash
# Simple build for current platform
go build -o build/gsecutil .

# Build with version information
go build -ldflags "-X main.Version=1.0.0" -o build/gsecutil .

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o build/gsecutil-linux-amd64 .
GOOS=windows GOARCH=amd64 go build -o build/gsecutil-windows-amd64.exe .
GOOS=darwin GOARCH=amd64 go build -o build/gsecutil-darwin-amd64 .
GOOS=darwin GOARCH=arm64 go build -o build/gsecutil-darwin-arm64 .
GOOS=linux GOARCH=arm64 go build -o build/gsecutil-linux-arm64 .

# Optimized builds (smaller binaries)
go build -ldflags "-X main.Version=1.0.0 -s -w" -o build/gsecutil .
```

## Build Targets and Platforms

### Supported Platforms

| OS      | Architecture | Binary Name                  |
|---------|--------------|------------------------------|
| Linux   | amd64        | `gsecutil-linux-amd64`       |
| Linux   | arm64        | `gsecutil-linux-arm64`       |
| macOS   | amd64        | `gsecutil-darwin-amd64`      |
| macOS   | arm64        | `gsecutil-darwin-arm64`      |
| Windows | amd64        | `gsecutil-windows-amd64.exe` |

### Build Output

All build methods create binaries in the `build/` directory:

```
build/
├── gsecutil                    # Current platform binary
├── gsecutil-linux-amd64        # Linux x64
├── gsecutil-linux-arm64        # Linux ARM64
├── gsecutil-darwin-amd64       # macOS Intel
├── gsecutil-darwin-arm64       # macOS Apple Silicon
└── gsecutil-windows-amd64.exe  # Windows x64
```

## Version Management

### Setting Version at Build Time

```bash
# Using environment variable
export VERSION=1.2.3
make build
# OR
VERSION=1.2.3 ./build.sh

# Using PowerShell
$env:VERSION = "1.2.3"
.\build.ps1
```

### Default Version

If no version is specified, builds default to `1.0.0`.

## Development Workflow

### Standard Development Cycle

```bash
# 1. Install/update dependencies
make deps

# 2. Format code
make fmt

# 3. Run static analysis
make vet

# 4. Run linter (if available)
make lint

# 5. Run tests
make test

# 6. Build for current platform
make build

# 7. Test the build
./build/gsecutil --help

# 8. Build for all platforms (when ready)
make build-all
```

### Using Build Scripts

```bash
# All-in-one development workflow
./build.sh deps && ./build.sh fmt && ./build.sh vet && ./build.sh test && ./build.sh dev
```

### Quick Development Build and Test

```bash
# Build and immediately show help
make dev
# OR
./build.sh dev
# OR
.\build.ps1 dev
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Build
on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21

    - name: Install dependencies
      run: make deps

    - name: Format code
      run: make fmt

    - name: Run vet
      run: make vet

    - name: Run tests
      run: make test

    - name: Build all platforms
      run: make build-all

    - name: Upload artifacts
      uses: actions/upload-artifact@v3
      with:
        name: binaries
        path: build/
```

### Docker Build Example

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.Version=${VERSION:-1.0.0} -s -w" \
    -o gsecutil .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/gsecutil /usr/local/bin/
ENTRYPOINT ["gsecutil"]
```

## Troubleshooting Build Issues

### Common Problems and Solutions

#### 1. "Go command not found"
```bash
# Install Go from https://golang.org/dl/
# Add Go to PATH
export PATH=$PATH:/usr/local/go/bin
```

#### 2. "Make command not found" (Windows)
```bash
# Use PowerShell script instead
.\build.ps1

# Or install make via chocolatey
choco install make
```

#### 3. "Permission denied" on build scripts
```bash
# Make script executable (Linux/macOS)
chmod +x build.sh

# Use PowerShell execution policy (Windows)
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

#### 4. Module download issues
```bash
# Clean module cache and retry
go clean -modcache
go mod download
```

#### 5. Cross-compilation issues
```bash
# Ensure target platform is supported
go tool dist list

# For CGO-dependent packages, may need to disable
CGO_ENABLED=0 GOOS=target_os GOARCH=target_arch go build
```

### Build Environment Verification

```bash
# Check Go installation
go version
go env GOOS GOARCH

# Check module status
go mod verify
go mod tidy

# List available build targets
go tool dist list | grep -E "(linux|darwin|windows)"
```

## Advanced Build Options

### Build Tags

```bash
# Build with specific tags
go build -tags "debug,local" .
```

### Custom LDFLAGS

```bash
# Inject multiple build-time variables
go build -ldflags "
  -X main.Version=1.2.3
  -X main.BuildTime=$(date -u '+%Y-%m-%d_%H:%M:%S')
  -X main.GitCommit=$(git rev-parse HEAD)
  -s -w
" .
```

### Debugging Build Issues

```bash
# Verbose build output
go build -v .

# Show build commands
go build -x .

# Build with debug info (larger binary)
go build -gcflags="all=-N -l" .
```

## Performance Considerations

### Binary Size Optimization

The build scripts use `-ldflags "-s -w"` to reduce binary size:
- `-s`: Strip symbol table and debug information
- `-w`: Strip DWARF debug information

### Build Speed

For faster development builds:
```bash
# Skip optimization for faster builds
go build .

# Use build cache
export GOCACHE=$HOME/.cache/go-build
```

## Platform-Specific Notes

### Linux
- Uses standard Go cross-compilation
- ARM64 builds supported for modern ARM processors
- Static binaries work well in containers

### macOS
- Universal binaries not created (separate amd64/arm64 builds)
- Code signing not included in build process
- May require developer tools for some dependencies

### Windows
- `.exe` extension automatically added
- PowerShell script provides Windows-native experience
- Works with both PowerShell 5.1 and PowerShell 7+

## Getting Help

For build-related issues:

1. Check this document for common solutions
2. Verify prerequisites are installed correctly
3. Try cleaning and rebuilding: `make clean && make build`
4. Check Go environment: `go env`
5. Open an issue with your build environment details

## Build Script Comparison

| Feature | Makefile | build.sh | build.ps1 |
|---------|----------|----------|-----------|
| Platform | Linux/macOS/WSL | Linux/macOS | Windows |
| Colors | No | Yes | Yes |
| Progress | Basic | Enhanced | Enhanced |
| Help | Yes | Yes | Yes |
| All Targets | Yes | Yes | Yes |
| Cross-platform | Yes | Yes | Yes |
| Size reporting | No | Yes | Yes |
| Error handling | Basic | Enhanced | Enhanced |

Choose the method that best fits your development environment!

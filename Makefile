# Makefile for gsecutil

# Binary name
BINARY_NAME=gsecutil

# Version (can be overridden)
VERSION?=1.0.0

# Build directory
BUILD_DIR=build

# Targets
.PHONY: all clean build build-linux build-windows build-darwin build-all test fmt vet

all: build-all

# Clean build directory
clean:
	rm -rf $(BUILD_DIR)
	go clean

# Create build directory
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Build for current platform
build: $(BUILD_DIR)
	go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) .

# Build for Linux (amd64)
build-linux: $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .

# Build for Linux (arm64)
build-linux-arm64: $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .

# Build for Windows (amd64)
build-windows: $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

# Build for macOS (amd64)
build-darwin: $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .

# Build for macOS (arm64)
build-darwin-arm64: $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .

# Build for all platforms
build-all: build-linux build-linux-arm64 build-windows build-darwin build-darwin-arm64

# Run tests
test:
	go test ./...

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Run linter (if golangci-lint is installed)
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod tidy
	go mod download

# Install locally
install:
	go install -ldflags "-X main.Version=$(VERSION)" .

# Development build and run
dev: build
	./$(BUILD_DIR)/$(BINARY_NAME) --help

# Help
help:
	@echo "Available targets:"
	@echo "  all           - Build for all platforms (default)"
	@echo "  build         - Build for current platform"
	@echo "  build-linux   - Build for Linux amd64"
	@echo "  build-linux-arm64 - Build for Linux arm64"
	@echo "  build-windows - Build for Windows amd64"
	@echo "  build-darwin  - Build for macOS amd64"
	@echo "  build-darwin-arm64 - Build for macOS arm64"
	@echo "  build-all     - Build for all platforms"
	@echo "  clean         - Clean build directory"
	@echo "  test          - Run tests"
	@echo "  fmt           - Format code"
	@echo "  vet           - Run go vet"
	@echo "  lint          - Run golangci-lint (if installed)"
	@echo "  deps          - Install dependencies"
	@echo "  install       - Install locally"
	@echo "  dev           - Development build and show help"
	@echo "  help          - Show this help message"

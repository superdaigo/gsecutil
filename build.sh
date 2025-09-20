#!/bin/bash

# Build script for gsecutil
# Builds binaries for multiple platforms
# Compatible with Makefile targets

set -e

BINARY_NAME="gsecutil"
VERSION=${VERSION:-$(cat VERSION 2>/dev/null || echo "1.0.0")}
BUILD_DIR="build"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to show help
show_help() {
    echo -e "${CYAN}Usage: $0 [target]${NC}"
    echo ""
    echo -e "${BLUE}Available targets:${NC}"
    echo -e "  ${GREEN}default${NC}       - Build for current platform"
    echo -e "  ${GREEN}all${NC}           - Build for all platforms"
    echo -e "  ${GREEN}clean${NC}         - Clean build directory"
    echo -e "  ${GREEN}linux${NC}         - Build for Linux amd64"
    echo -e "  ${GREEN}linux-arm64${NC}   - Build for Linux arm64"
    echo -e "  ${GREEN}windows${NC}       - Build for Windows amd64"
    echo -e "  ${GREEN}darwin${NC}        - Build for macOS amd64"
    echo -e "  ${GREEN}darwin-arm64${NC}  - Build for macOS arm64"
    echo -e "  ${GREEN}test${NC}          - Run tests"
    echo -e "  ${GREEN}fmt${NC}           - Format code"
    echo -e "  ${GREEN}vet${NC}           - Run go vet"
    echo -e "  ${GREEN}lint${NC}          - Run golangci-lint (if installed)"
    echo -e "  ${GREEN}deps${NC}          - Install dependencies"
    echo -e "  ${GREEN}install${NC}       - Install locally"
    echo -e "  ${GREEN}dev${NC}           - Development build and show help"
    echo -e "  ${GREEN}help${NC}          - Show this help message"
}

# Function to clean build directory
clean_build_dir() {
    echo -e "${BLUE}Cleaning build directory...${NC}"
    rm -rf $BUILD_DIR
    go clean
    echo -e "${GREEN}✓ Clean completed${NC}"
}

# Function to create build directory if it doesn't exist
create_build_dir() {
    if [ ! -d "$BUILD_DIR" ]; then
        echo -e "${BLUE}Creating build directory...${NC}"
        mkdir -p $BUILD_DIR
        echo -e "${GREEN}✓ Build directory created${NC}"
    fi
}

# Function to build for a specific OS/ARCH
build_for_target() {
    local GOOS=$1
    local GOARCH=$2
    local output_suffix=$3

    echo -e "${YELLOW}Building for $GOOS/$GOARCH...${NC}"

    if [ -z "$output_suffix" ]; then
        output_name="$BINARY_NAME-$GOOS-$GOARCH"
    else
        output_name="$BINARY_NAME$output_suffix"
    fi

    if [ $GOOS = "windows" ]; then
        output_name="$output_name.exe"
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags "-X main.Version=$VERSION -s -w" \
        -o "$BUILD_DIR/$output_name" \
        .

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Built $output_name${NC}"
        return 0
    else
        echo -e "${RED}✗ Failed to build $output_name${NC}"
        return 1
    fi
}

# Define build targets (OS/ARCH combinations)
declare -a all_targets=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

# Function to build for all platforms
build_all() {
    echo -e "${BLUE}Building $BINARY_NAME v$VERSION for all platforms${NC}"
    create_build_dir

    local failed=0

    # Build for each target
    for target in "${all_targets[@]}"; do
        IFS='/' read -ra PARTS <<< "$target"
        GOOS=${PARTS[0]}
        GOARCH=${PARTS[1]}

        build_for_target $GOOS $GOARCH
        if [ $? -ne 0 ]; then
            failed=1
        fi
    done

    if [ $failed -eq 0 ]; then
        show_build_summary
        echo -e "${GREEN}All builds completed successfully!${NC}"
        return 0
    else
        echo -e "${RED}Some builds failed. See above for details.${NC}"
        return 1
    fi
}

# Function to build for current platform
build_current() {
    echo -e "${BLUE}Building $BINARY_NAME v$VERSION for current platform${NC}"
    create_build_dir

    # No GOOS/GOARCH specified means current platform
    go build -ldflags "-X main.Version=$VERSION -s -w" -o "$BUILD_DIR/$BINARY_NAME" .

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Built $BINARY_NAME for current platform${NC}"
        show_build_summary
        return 0
    else
        echo -e "${RED}✗ Failed to build $BINARY_NAME for current platform${NC}"
        return 1
    fi
}

# Function to run tests
run_tests() {
    echo -e "${BLUE}Running tests...${NC}"
    go test ./...
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Tests passed${NC}"
        return 0
    else
        echo -e "${RED}✗ Tests failed${NC}"
        return 1
    fi
}

# Function to format code
format_code() {
    echo -e "${BLUE}Formatting code...${NC}"
    go fmt ./...
    echo -e "${GREEN}✓ Code formatted${NC}"
    return 0
}

# Function to run vet
run_vet() {
    echo -e "${BLUE}Running go vet...${NC}"
    go vet ./...
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Go vet passed${NC}"
        return 0
    else
        echo -e "${RED}✗ Go vet found issues${NC}"
        return 1
    fi
}

# Function to run linter
run_lint() {
    if command -v golangci-lint &> /dev/null; then
        echo -e "${BLUE}Running golangci-lint...${NC}"
        golangci-lint run
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✓ Linting passed${NC}"
            return 0
        else
            echo -e "${RED}✗ Linting found issues${NC}"
            return 1
        fi
    else
        echo -e "${YELLOW}⚠ golangci-lint not installed, skipping lint${NC}"
        return 0
    fi
}

# Function to install dependencies
install_deps() {
    echo -e "${BLUE}Installing dependencies...${NC}"
    go mod tidy
    go mod download
    echo -e "${GREEN}✓ Dependencies installed${NC}"
    return 0
}

# Function to install locally
install_locally() {
    echo -e "${BLUE}Installing $BINARY_NAME v$VERSION locally...${NC}"
    go install -ldflags "-X main.Version=$VERSION" .
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Installed successfully${NC}"
        return 0
    else
        echo -e "${RED}✗ Installation failed${NC}"
        return 1
    fi
}

# Function for development build and run
run_dev() {
    build_current
    if [ $? -eq 0 ]; then
        echo -e "${BLUE}Running $BINARY_NAME help...${NC}"
        "$BUILD_DIR/$BINARY_NAME" --help
    fi
}

# Function to show build summary
show_build_summary() {
    if [ ! -d "$BUILD_DIR" ]; then
        return 0
    fi

    echo -e "${BLUE}"
    echo "Build completed!"
    echo "Binaries are available in the $BUILD_DIR directory:"
    echo -e "${NC}"
    ls -la $BUILD_DIR/

    # Calculate sizes
    echo -e "${BLUE}"
    echo "Binary sizes:"
    echo -e "${NC}"
    for file in $BUILD_DIR/*; do
        if [ -f "$file" ]; then
            size=$(du -h "$file" | cut -f1)
            filename=$(basename "$file")
            echo "  $filename: $size"
        fi
    done
}

# Main logic - process command line argument
if [ $# -eq 0 ]; then
    # Default action: build for current platform
    build_current
exit $?
fi

case "$1" in
    all)
        build_all
        ;;
    clean)
        clean_build_dir
        ;;
    linux)
        create_build_dir
        build_for_target "linux" "amd64"
        show_build_summary
        ;;
    linux-arm64)
        create_build_dir
        build_for_target "linux" "arm64"
        show_build_summary
        ;;
    windows)
        create_build_dir
        build_for_target "windows" "amd64"
        show_build_summary
        ;;
    darwin)
        create_build_dir
        build_for_target "darwin" "amd64"
        show_build_summary
        ;;
    darwin-arm64)
        create_build_dir
        build_for_target "darwin" "arm64"
        show_build_summary
        ;;
    test)
        run_tests
        ;;
    fmt)
        format_code
        ;;
    vet)
        run_vet
        ;;
    lint)
        run_lint
        ;;
    deps)
        install_deps
        ;;
    install)
        install_locally
        ;;
    dev)
        run_dev
        ;;
    help)
        show_help
        ;;
    *)
        echo -e "${RED}Unknown target: $1${NC}"
        show_help
        exit 1
        ;;
esac

exit $?

#!/bin/bash

# Build script for gsecutil
# Builds binaries for multiple platforms

set -e

BINARY_NAME="gsecutil"
VERSION=${VERSION:-"1.0.0"}
BUILD_DIR="build"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Building $BINARY_NAME v$VERSION${NC}"

# Clean and create build directory
rm -rf $BUILD_DIR
mkdir -p $BUILD_DIR

# Define build targets (OS/ARCH combinations)
declare -a targets=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

# Build for each target
for target in "${targets[@]}"; do
    IFS='/' read -ra PARTS <<< "$target"
    GOOS=${PARTS[0]}
    GOARCH=${PARTS[1]}
    
    output_name="$BINARY_NAME-$GOOS-$GOARCH"
    if [ $GOOS = "windows" ]; then
        output_name="$output_name.exe"
    fi
    
    echo -e "${YELLOW}Building for $GOOS/$GOARCH...${NC}"
    
    env GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags "-X main.Version=$VERSION -s -w" \
        -o "$BUILD_DIR/$output_name" \
        .
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Built $output_name${NC}"
    else
        echo -e "${RED}✗ Failed to build $output_name${NC}"
        exit 1
    fi
done

echo -e "${BLUE}"
echo "Build completed successfully!"
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

echo -e "${GREEN}"
echo "All builds completed successfully!"
echo -e "${NC}"

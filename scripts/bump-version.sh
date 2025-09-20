#!/bin/bash

# Version bumping script for gsecutil
# Updates VERSION file and optionally creates a release

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to show help
show_help() {
    echo -e "${BLUE}gsecutil Version Bumping Tool${NC}"
    echo ""
    echo "Usage: $0 [VERSION] [--release]"
    echo ""
    echo "Updates the VERSION file and optionally creates a release."
    echo ""
    echo "Arguments:"
    echo "  VERSION        New version number (e.g., 1.2.3, 2.0.0-beta.1)"
    echo "                 If not provided, will increment patch version"
    echo ""
    echo "Options:"
    echo "  --release      Create and push a git tag to trigger release"
    echo "  --major        Increment major version (X.0.0)"
    echo "  --minor        Increment minor version (X.Y.0)"  
    echo "  --patch        Increment patch version (X.Y.Z) [default]"
    echo "  --dry-run      Show what would be done without making changes"
    echo ""
    echo "Examples:"
    echo "  $0 1.2.3           # Set version to 1.2.3"
    echo "  $0 1.2.3 --release # Set version and create release"
    echo "  $0 --minor         # Increment minor version"
    echo "  $0 --major --release # Increment major and release"
    echo "  $0 --dry-run       # Preview next patch version"
}

# Function to validate version format
validate_version() {
    local version=$1
    if [[ ! $version =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.-]+)?$ ]]; then
        echo -e "${RED}Error: Invalid version format. Expected format: X.Y.Z or X.Y.Z-suffix${NC}"
        echo "Examples: 1.2.3, 2.0.0-beta.1, 1.0.0-rc.1"
        return 1
    fi
    return 0
}

# Function to get current version
get_current_version() {
    if [ -f "VERSION" ]; then
        cat VERSION
    else
        echo "0.0.0"
    fi
}

# Function to increment version
increment_version() {
    local version=$1
    local part=$2
    
    IFS='.' read -ra VERSION_PARTS <<< "$version"
    local major=${VERSION_PARTS[0]:-0}
    local minor=${VERSION_PARTS[1]:-0}
    local patch=${VERSION_PARTS[2]:-0}
    
    # Remove any suffix from patch version
    patch=$(echo "$patch" | sed 's/-.*//')
    
    case $part in
        "major")
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        "minor")
            minor=$((minor + 1))
            patch=0
            ;;
        "patch"|*)
            patch=$((patch + 1))
            ;;
    esac
    
    echo "$major.$minor.$patch"
}

# Function to update VERSION file
update_version_file() {
    local new_version=$1
    local dry_run=$2
    
    if [ "$dry_run" = "true" ]; then
        echo -e "${YELLOW}[DRY RUN] Would update VERSION file to: $new_version${NC}"
    else
        echo "$new_version" > VERSION
        echo -e "${GREEN}✓ Updated VERSION file to: $new_version${NC}"
    fi
}

# Function to test build with new version
test_build() {
    local dry_run=$1
    
    if [ "$dry_run" = "true" ]; then
        echo -e "${YELLOW}[DRY RUN] Would test build with new version${NC}"
    else
        echo -e "${BLUE}Testing build with new version...${NC}"
        if make build >/dev/null 2>&1; then
            echo -e "${GREEN}✓ Build test passed${NC}"
            # Show version from built binary
            if [ -f "build/gsecutil" ]; then
                BUILT_VERSION=$(./build/gsecutil --version 2>/dev/null | grep -o '[0-9]\+\.[0-9]\+\.[0-9]\+' || echo "unknown")
                echo -e "${GREEN}✓ Built binary version: $BUILT_VERSION${NC}"
            fi
        else
            echo -e "${RED}✗ Build test failed${NC}"
            return 1
        fi
    fi
}

# Function to create release
create_release() {
    local version=$1
    local dry_run=$2
    
    if [ "$dry_run" = "true" ]; then
        echo -e "${YELLOW}[DRY RUN] Would create release tag: v$version${NC}"
    else
        echo -e "${BLUE}Creating release for version $version...${NC}"
        ./scripts/release.sh "$version"
    fi
}

# Main function
main() {
    local new_version=""
    local create_release_flag=false
    local increment_type="patch"
    local dry_run=false
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            --release)
                create_release_flag=true
                shift
                ;;
            --major)
                increment_type="major"
                shift
                ;;
            --minor)
                increment_type="minor"
                shift
                ;;
            --patch)
                increment_type="patch"
                shift
                ;;
            --dry-run)
                dry_run=true
                shift
                ;;
            -*)
                echo -e "${RED}Unknown option: $1${NC}"
                show_help
                exit 1
                ;;
            *)
                if [ -z "$new_version" ]; then
                    new_version=$1
                else
                    echo -e "${RED}Too many arguments${NC}"
                    show_help
                    exit 1
                fi
                shift
                ;;
        esac
    done
    
    echo -e "${BLUE}🔧 gsecutil Version Bumping Tool${NC}"
    echo ""
    
    # Get current version
    current_version=$(get_current_version)
    echo -e "Current version: ${YELLOW}$current_version${NC}"
    
    # Determine new version
    if [ -z "$new_version" ]; then
        new_version=$(increment_version "$current_version" "$increment_type")
        echo -e "Auto-incremented ${increment_type} version: ${GREEN}$new_version${NC}"
    else
        # Validate provided version
        if ! validate_version "$new_version"; then
            exit 1
        fi
        echo -e "Target version: ${GREEN}$new_version${NC}"
    fi
    
    # Check if version is the same
    if [ "$current_version" = "$new_version" ]; then
        echo -e "${YELLOW}Version is already $new_version${NC}"
        if [ "$dry_run" = "false" ] && [ "$create_release_flag" = "false" ]; then
            exit 0
        fi
    fi
    
    echo ""
    
    # Update VERSION file
    update_version_file "$new_version" "$dry_run"
    
    # Test build
    test_build "$dry_run"
    
    # Create release if requested
    if [ "$create_release_flag" = "true" ]; then
        echo ""
        create_release "$new_version" "$dry_run"
    fi
    
    if [ "$dry_run" = "false" ]; then
        echo ""
        echo -e "${GREEN}✅ Version bump completed!${NC}"
        echo ""
        echo "Next steps:"
        if [ "$create_release_flag" = "false" ]; then
            echo -e "  ${BLUE}1.${NC} Test your changes"
            echo -e "  ${BLUE}2.${NC} Commit the VERSION file: ${YELLOW}git add VERSION && git commit -m \"Bump version to $new_version\"${NC}"
            echo -e "  ${BLUE}3.${NC} Create release: ${YELLOW}./scripts/release.sh $new_version${NC}"
        else
            echo -e "  ${BLUE}1.${NC} Monitor release: ${YELLOW}https://github.com/$(git remote get-url origin | sed 's/.*github.com[/:]\([^/]*\/[^.]*\).*/\1/')/actions${NC}"
        fi
    fi
}

# Call main function with all arguments
main "$@"

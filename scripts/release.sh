#!/bin/bash

# Release helper script for gsecutil
# Creates and pushes version tags to trigger GitHub Actions release workflow

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to show help
show_help() {
    echo -e "${BLUE}gsecutil Release Helper${NC}"
    echo ""
    echo "Usage: $0 [VERSION]"
    echo ""
    echo "Creates a version tag and pushes it to trigger the release workflow."
    echo ""
    echo "Arguments:"
    echo "  VERSION    Version number (e.g., 1.2.3, 2.0.0-beta.1)"
    echo "             If not provided, will prompt for input"
    echo ""
    echo "Examples:"
    echo "  $0 1.2.3           # Create release v1.2.3"
    echo "  $0 2.0.0-beta.1    # Create pre-release v2.0.0-beta.1"
    echo "  $0                 # Interactive mode"
    echo ""
    echo "Requirements:"
    echo "  - Clean working directory (no uncommitted changes)"
    echo "  - On main branch (recommended)"
    echo "  - All tests passing"
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

# Function to check if we're on main branch
check_branch() {
    local current_branch=$(git rev-parse --abbrev-ref HEAD)
    if [ "$current_branch" != "main" ]; then
        echo -e "${YELLOW}Warning: You're not on the main branch (currently on: $current_branch)${NC}"
        echo -e "It's recommended to create releases from the main branch."
        echo ""
        read -p "Do you want to continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${YELLOW}Aborted by user${NC}"
            exit 1
        fi
    fi
}

# Function to check working directory status
check_working_directory() {
    if ! git diff-index --quiet HEAD --; then
        echo -e "${RED}Error: Working directory is not clean. Please commit or stash your changes first.${NC}"
        echo ""
        echo "Uncommitted changes:"
        git status --porcelain
        exit 1
    fi
}

# Function to check if tag already exists
check_existing_tag() {
    local version=$1
    local tag="v$version"

    if git rev-parse "$tag" >/dev/null 2>&1; then
        echo -e "${RED}Error: Tag '$tag' already exists${NC}"
        echo ""
        echo "Existing tags:"
        git tag -l | grep -E "^v[0-9]" | sort -V | tail -5
        exit 1
    fi
}

# Function to run tests
run_tests() {
    echo -e "${BLUE}Running tests...${NC}"

    if ! make test >/dev/null 2>&1; then
        echo -e "${RED}Error: Tests failed. Please fix the tests before creating a release.${NC}"
        echo ""
        echo "Run 'make test' to see the detailed test output."
        exit 1
    fi

    echo -e "${GREEN}✓ Tests passed${NC}"
}

# Function to show recent tags
show_recent_tags() {
    echo -e "${BLUE}Recent tags:${NC}"
    git tag -l | grep -E "^v[0-9]" | sort -V | tail -5 || echo "No existing version tags found"
    echo ""
}

# Function to show what will be included in the release
show_changes_since_last_tag() {
    local latest_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "")

    if [ -z "$latest_tag" ]; then
        echo -e "${BLUE}Changes since repository creation:${NC}"
        git log --oneline --reverse | head -10
        local total_commits=$(git rev-list --count HEAD)
        if [ $total_commits -gt 10 ]; then
            echo "... and $((total_commits - 10)) more commits"
        fi
    else
        echo -e "${BLUE}Changes since $latest_tag:${NC}"
        local changes=$(git log $latest_tag..HEAD --oneline)
        if [ -z "$changes" ]; then
            echo -e "${YELLOW}No new changes since last tag${NC}"
            echo ""
            read -p "Do you want to create a release anyway? (y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                echo -e "${YELLOW}Aborted by user${NC}"
                exit 1
            fi
        else
            echo "$changes"
        fi
    fi
    echo ""
}

# Main function
main() {
    local version=$1

    # Show help if requested
    if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
        show_help
        exit 0
    fi

    echo -e "${BLUE}🚀 gsecutil Release Helper${NC}"
    echo ""

    # Check prerequisites
    check_working_directory
    check_branch

    # Show recent tags
    show_recent_tags

    # Get version if not provided
    if [ -z "$version" ]; then
        echo -e "${YELLOW}Enter the version number (without 'v' prefix):${NC}"
        read -p "Version: " version

        if [ -z "$version" ]; then
            echo -e "${RED}Error: Version cannot be empty${NC}"
            exit 1
        fi
    fi

    # Validate version format
    if ! validate_version "$version"; then
        exit 1
    fi

    # Check if tag already exists
    check_existing_tag "$version"

    # Show what will be included
    show_changes_since_last_tag

    # Run tests
    run_tests

    local tag="v$version"

    # Confirm release creation
    echo -e "${YELLOW}About to create and push tag: $tag${NC}"
    echo ""
    echo "This will trigger the GitHub Actions release workflow which will:"
    echo "  - Build binaries for all platforms"
    echo "  - Create a GitHub release with download links"
    echo "  - Generate checksums for all binaries"
    echo "  - Publish the release automatically"
    echo ""
    read -p "Are you sure you want to continue? (y/N): " -n 1 -r
    echo

    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}Aborted by user${NC}"
        exit 1
    fi

    # Create and push the tag
    echo -e "${BLUE}Creating tag $tag...${NC}"
    git tag -a "$tag" -m "Release $version"

    echo -e "${BLUE}Pushing tag to origin...${NC}"
    git push origin "$tag"

    echo -e "${GREEN}✅ Success!${NC}"
    echo ""
    echo -e "${GREEN}Tag $tag has been created and pushed.${NC}"
    echo ""
    echo "The GitHub Actions release workflow is now running."
    echo "You can monitor the progress at:"
    echo -e "${BLUE}https://github.com/$(git remote get-url origin | sed 's/.*github.com[/:]\([^/]*\/[^.]*\).*/\1/')/actions${NC}"
    echo ""
    echo "The release will be available at:"
    echo -e "${BLUE}https://github.com/$(git remote get-url origin | sed 's/.*github.com[/:]\([^/]*\/[^.]*\).*/\1/')/releases/tag/$tag${NC}"
    echo ""
    echo -e "${YELLOW}Note: It may take a few minutes for the release to be fully processed.${NC}"
}

# Check if git is available
if ! command -v git >/dev/null 2>&1; then
    echo -e "${RED}Error: git is not installed or not in PATH${NC}"
    exit 1
fi

# Check if we're in a git repository
if ! git rev-parse --git-dir >/dev/null 2>&1; then
    echo -e "${RED}Error: Not in a git repository${NC}"
    exit 1
fi

# Run main function
main "$@"

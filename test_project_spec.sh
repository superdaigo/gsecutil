#!/bin/bash

# Test script to verify project specification methods across all commands
# This script tests --project, -p, and GSECUTIL_PROJECT environment variable

# set -e  # Don't exit on error, we're testing

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

TEST_PROJECT="test-project-12345"
GSECUTIL="./build/gsecutil"

echo -e "${BLUE}Testing Project Specification Methods${NC}"
echo "======================================"
echo ""

# Array of commands to test with their required arguments
# Format: "command:args:description"
declare -a COMMANDS=(
    "get:test-secret:Get secret"
    "create:test-secret --data test:Create secret"
    "update:test-secret --data test:Update secret"
    "delete:test-secret --force:Delete secret"
    "list::List secrets"
    "describe:test-secret:Describe secret"
    "auditlog::Audit log"
    "access list:test-secret:Access list"
)

test_command() {
    local cmd=$1
    local args=$2
    local desc=$3
    local method=$4
    local project_spec=$5

    echo -n "  Testing '$desc' with $method... "

    # Build the command based on method
    case $method in
        "--project")
            full_cmd="$GSECUTIL $cmd $args --project $project_spec 2>&1"
            ;;
        "-p")
            full_cmd="$GSECUTIL $cmd $args -p $project_spec 2>&1"
            ;;
        "env")
            full_cmd="GSECUTIL_PROJECT=$project_spec $GSECUTIL $cmd $args 2>&1"
            ;;
    esac

    # Execute and check if project flag was recognized (not checking actual execution)
    # We're just verifying the flag is parsed correctly
    output=$(eval $full_cmd || true)

    # Check if the command doesn't complain about unknown flag or invalid project format
    if echo "$output" | grep -qiE "(unknown flag|unknown shorthand|invalid flag|unexpected argument)"; then
        echo -e "${RED}FAIL${NC}"
        echo "    Error: $output"
        return 1
    else
        echo -e "${GREEN}PASS${NC}"
        return 0
    fi
}

total_tests=0
passed_tests=0
failed_tests=0

echo -e "${YELLOW}Method 1: Using --project flag${NC}"
echo "--------------------------------"
for cmd_info in "${COMMANDS[@]}"; do
    IFS=':' read -r cmd args desc <<< "$cmd_info"
    ((total_tests++))
    if test_command "$cmd" "$args" "$desc" "--project" "$TEST_PROJECT"; then
        ((passed_tests++))
    else
        ((failed_tests++))
    fi
done
echo ""

echo -e "${YELLOW}Method 2: Using -p short flag${NC}"
echo "--------------------------------"
for cmd_info in "${COMMANDS[@]}"; do
    IFS=':' read -r cmd args desc <<< "$cmd_info"
    ((total_tests++))
    if test_command "$cmd" "$args" "$desc" "-p" "$TEST_PROJECT"; then
        ((passed_tests++))
    else
        ((failed_tests++))
    fi
done
echo ""

echo -e "${YELLOW}Method 3: Using GSECUTIL_PROJECT environment variable${NC}"
echo "--------------------------------------------------------"
for cmd_info in "${COMMANDS[@]}"; do
    IFS=':' read -r cmd args desc <<< "$cmd_info"
    ((total_tests++))
    if test_command "$cmd" "$args" "$desc" "env" "$TEST_PROJECT"; then
        ((passed_tests++))
    else
        ((failed_tests++))
    fi
done
echo ""

echo "======================================"
echo -e "${BLUE}Test Results${NC}"
echo "======================================"
echo "Total tests: $total_tests"
echo -e "${GREEN}Passed: $passed_tests${NC}"
if [ $failed_tests -gt 0 ]; then
    echo -e "${RED}Failed: $failed_tests${NC}"
    exit 1
else
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
fi

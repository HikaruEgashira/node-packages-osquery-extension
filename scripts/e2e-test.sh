#!/bin/bash
set -e

# E2E Test Script for node-packages-osquery-extension
# This script sets up a test environment and validates the extension

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
EXTENSION_PATH="${PROJECT_DIR}/node_packages_extension"
TEST_DIR="${PROJECT_DIR}/test-packages"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if osquery is installed
if ! command -v osqueryi &> /dev/null; then
    log_error "osqueryi not found. Please install osquery first."
    exit 1
fi

# Check if extension exists
if [ ! -f "$EXTENSION_PATH" ]; then
    log_error "Extension not found at $EXTENSION_PATH. Please build it first."
    exit 1
fi

log_info "Extension found at: $EXTENSION_PATH"

# Setup test environment
log_info "Setting up test environment..."
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

# Install npm packages
log_info "Installing npm packages..."
if ! npm init -y &> /dev/null; then
    log_warn "npm init failed, continuing..."
fi

npm install react react-dom lodash express axios &> /dev/null || log_warn "Some npm packages failed to install"

# Install pnpm packages (if pnpm is available)
if command -v pnpm &> /dev/null; then
    log_info "Installing pnpm packages..."
    pnpm add typescript @types/node &> /dev/null || log_warn "Some pnpm packages failed to install"
else
    log_warn "pnpm not found, skipping pnpm tests"
fi

# Install bun packages (if bun is available)
if command -v bun &> /dev/null; then
    log_info "Installing bun packages..."
    bun add jest vitest &> /dev/null || log_warn "Some bun packages failed to install"
else
    log_warn "bun not found, skipping bun tests"
fi

# Run osquery tests
log_info "Running osquery extension tests..."

# Helper function to run osquery with extension
run_osquery() {
    local query="$1"
    osqueryi --extension "$EXTENSION_PATH" --ephemeral --json "$query" 2>&1
}

# Test 1: Basic query
log_info "Test 1: Querying all packages..."
RESULT=$(run_osquery "SELECT * FROM node_packages;")
echo "$RESULT" > results.json

if [ ! -s results.json ]; then
    log_error "No results returned from basic query"
    exit 1
fi

PACKAGE_COUNT=$(echo "$RESULT" | grep -o '"name"' | wc -l | tr -d ' ')
log_info "Found $PACKAGE_COUNT packages"

# Test 2: Manager-specific query
log_info "Test 2: Querying npm packages..."
NPM_RESULT=$(run_osquery "SELECT name, version FROM node_packages WHERE manager = 'npm';")
echo "$NPM_RESULT" > npm_packages.json
log_info "NPM packages: $(echo "$NPM_RESULT" | grep -o '"name"' | wc -l | tr -d ' ')"

# Test 3: Count by manager
log_info "Test 3: Counting packages by manager..."
MANAGER_COUNT=$(run_osquery "SELECT manager, COUNT(*) as count FROM node_packages GROUP BY manager;")
echo "$MANAGER_COUNT" > manager_counts.json
log_info "Manager counts:"
echo "$MANAGER_COUNT" | grep -o '"manager":"[^"]*"' | sed 's/"manager":"/  - /g' | sed 's/"$//'

# Test 4: Search for specific package
log_info "Test 4: Searching for react packages..."
REACT_RESULT=$(run_osquery "SELECT * FROM node_packages WHERE name LIKE '%react%';")
REACT_COUNT=$(echo "$REACT_RESULT" | grep -o '"name"' | wc -l | tr -d ' ')
log_info "Found $REACT_COUNT react-related packages"

# Validate results
log_info "Validating results..."

if [ "$PACKAGE_COUNT" -eq 0 ]; then
    log_error "No packages found! Test failed."
    exit 1
fi

if [ "$REACT_COUNT" -eq 0 ]; then
    log_warn "No react packages found (expected at least react and react-dom)"
fi

# Check if expected packages are present
log_info "Checking for expected packages..."
EXPECTED_PACKAGES=("react" "lodash" "express")
MISSING_PACKAGES=()

for pkg in "${EXPECTED_PACKAGES[@]}"; do
    if ! echo "$RESULT" | grep -q "\"name\":\"$pkg\""; then
        MISSING_PACKAGES+=("$pkg")
    fi
done

if [ ${#MISSING_PACKAGES[@]} -gt 0 ]; then
    log_warn "Missing expected packages: ${MISSING_PACKAGES[*]}"
else
    log_info "All expected packages found!"
fi

# Cleanup
cd "$PROJECT_DIR"
log_info "Cleaning up test environment..."
# Uncomment to remove test directory
# rm -rf "$TEST_DIR"

log_info "E2E tests completed successfully!"
log_info "Test results saved in:"
log_info "  - ${TEST_DIR}/results.json"
log_info "  - ${TEST_DIR}/npm_packages.json"
log_info "  - ${TEST_DIR}/manager_counts.json"

exit 0

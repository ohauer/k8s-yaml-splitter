#!/bin/sh
#
# Test script for k8s-yaml-splitter
#

set -e

BINARY="./bin/k8s-yaml-splitter-linux-amd64"
TESTDATA_DIR="./testdata"
TEST_OUTPUT_BASE="./test-output"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    printf "${GREEN}[INFO]${NC} %s\n" "$1"
}

log_warn() {
    printf "${YELLOW}[WARN]${NC} %s\n" "$1"
}

log_error() {
    printf "${RED}[ERROR]${NC} %s\n" "$1"
}

cleanup() {
    log_info "Cleaning up test directories..."
    rm -rf ${TEST_OUTPUT_BASE}-*
}

check_binary() {
    if [ ! -f "$BINARY" ]; then
        log_error "Binary not found: $BINARY"
        log_info "Run 'make build' first"
        exit 1
    fi
}

test_basic_functionality() {
    log_info "Testing basic functionality..."
    OUTPUT_DIR="${TEST_OUTPUT_BASE}-basic"
    mkdir -p "$OUTPUT_DIR"

    $BINARY "$TESTDATA_DIR/filter-test.yaml" "$OUTPUT_DIR"

    if [ ! -d "$OUTPUT_DIR" ]; then
        log_error "Output directory not created"
        return 1
    fi

    file_count=$(find "$OUTPUT_DIR" -name "*.yaml" | wc -l)
    if [ "$file_count" -ne 5 ]; then
        log_error "Expected 5 files, got $file_count"
        return 1
    fi

    log_info "Basic functionality: PASSED"
}

test_error_handling() {
    log_info "Testing error handling..."
    OUTPUT_DIR="${TEST_OUTPUT_BASE}-errors"
    mkdir -p "$OUTPUT_DIR"

    # This should exit with code 1 due to errors but continue processing
    if $BINARY "$TESTDATA_DIR/with-errors.yaml" "$OUTPUT_DIR" >/dev/null 2>&1; then
        log_warn "Expected non-zero exit code for malformed YAML"
    fi

    # Should still create some files
    file_count=$(find "$OUTPUT_DIR" -name "*.yaml" | wc -l)
    if [ "$file_count" -eq 0 ]; then
        log_error "No files created despite some valid documents"
        return 1
    fi

    log_info "Error handling: PASSED"
}

run_all_tests() {
    log_info "Starting k8s-yaml-splitter tests..."

    check_binary
    cleanup

    test_basic_functionality
    test_error_handling

    cleanup

    log_info "All tests PASSED!"
}

# Run tests
run_all_tests

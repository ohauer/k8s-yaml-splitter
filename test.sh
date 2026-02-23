#!/bin/sh
#
# Test script for k8s-yaml-splitter
# POSIX compliant shell script
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

    $BINARY -f "$TESTDATA_DIR/filter-test.yaml" -d "$OUTPUT_DIR"

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

test_namespace_dirs() {
    log_info "Testing namespace directories..."
    OUTPUT_DIR="${TEST_OUTPUT_BASE}-namespace"

    $BINARY -f "$TESTDATA_DIR/multi-namespace.yaml" -namespace-dirs -d "$OUTPUT_DIR"

    if [ ! -d "$OUTPUT_DIR/frontend" ] || [ ! -d "$OUTPUT_DIR/backend" ] || [ ! -d "$OUTPUT_DIR/cluster-scoped" ]; then
        log_error "Namespace directories not created"
        return 1
    fi

    log_info "Namespace directories: PASSED"
}

test_filtering() {
    log_info "Testing resource filtering..."
    OUTPUT_DIR="${TEST_OUTPUT_BASE}-filter"

    $BINARY -f "$TESTDATA_DIR/filter-test.yaml" -include "Deployment,Service" -d "$OUTPUT_DIR"

    file_count=$(find "$OUTPUT_DIR" -name "*.yaml" | wc -l)
    if [ "$file_count" -ne 2 ]; then
        log_error "Expected 2 files with include filter, got $file_count"
        return 1
    fi

    log_info "Resource filtering: PASSED"
}

test_stdin() {
    log_info "Testing stdin input..."
    OUTPUT_DIR="${TEST_OUTPUT_BASE}-stdin"

    cat "$TESTDATA_DIR/filter-test.yaml" | $BINARY -f - -d "$OUTPUT_DIR"

    file_count=$(find "$OUTPUT_DIR" -name "*.yaml" | wc -l)
    if [ "$file_count" -ne 5 ]; then
        log_error "Expected 5 files from stdin, got $file_count"
        return 1
    fi

    log_info "Stdin input: PASSED"
}

test_json_output() {
    log_info "Testing JSON output..."
    OUTPUT_DIR="${TEST_OUTPUT_BASE}-json"

    $BINARY -f "$TESTDATA_DIR/filter-test.yaml" -o json -d "$OUTPUT_DIR"

    json_count=$(find "$OUTPUT_DIR" -name "*.json" | wc -l)
    if [ "$json_count" -ne 5 ]; then
        log_error "Expected 5 JSON files, got $json_count"
        return 1
    fi

    log_info "JSON output: PASSED"
}

test_error_handling() {
    log_info "Testing error handling..."
    OUTPUT_DIR="${TEST_OUTPUT_BASE}-errors"

    # This should exit with code 1 due to errors but continue processing
    if $BINARY -f "$TESTDATA_DIR/with-errors.yaml" -d "$OUTPUT_DIR" >/dev/null 2>&1; then
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

test_dry_run() {
    log_info "Testing dry run mode..."
    OUTPUT_DIR="${TEST_OUTPUT_BASE}-dryrun"

    $BINARY -f "$TESTDATA_DIR/filter-test.yaml" -dry-run -d "$OUTPUT_DIR" >/dev/null

    if [ -d "$OUTPUT_DIR" ] && [ "$(find "$OUTPUT_DIR" -name "*.yaml" | wc -l)" -gt 0 ]; then
        log_error "Files created in dry-run mode"
        return 1
    fi

    log_info "Dry run mode: PASSED"
}

run_all_tests() {
    log_info "Starting k8s-yaml-splitter tests..."

    check_binary
    cleanup

    test_basic_functionality
    test_namespace_dirs
    test_filtering
    test_stdin
    test_json_output
    test_error_handling
    test_dry_run

    cleanup

    log_info "All tests PASSED!"
}

# Run tests
run_all_tests

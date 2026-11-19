#!/bin/sh
#
# Release helper for k8s-yaml-splitter
#
# Purpose: Verifies the working tree, runs all checks, then creates and pushes
#          an annotated version tag. The GitHub Actions workflow
#          (.github/workflows/go.yml) builds the binaries and publishes the
#          GitHub release with generated release notes when the tag arrives.
#
# Usage:   scripts/release.sh v0.2.0        (release)
#          scripts/release.sh v0.3.0-rc1    (pre-release, detected by "-rc")
#
# Requires: git, make, go, golangci-lint, govulncheck

set -eu

usage() {
    echo "Usage: $0 vMAJOR.MINOR.PATCH[-rcN]" >&2
    exit 2
}

log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') release - $1"
}

main() {
    version="${1:-}"
    [ -n "${version}" ] || usage
    case "${version}" in
        v[0-9]*.[0-9]*.[0-9]*) ;;
        *) usage ;;
    esac

    for tool in git make go; do
        command -v "${tool}" > /dev/null 2>&1 || { log "ERROR: ${tool} not found"; exit 3; }
    done

    branch=$(git rev-parse --abbrev-ref HEAD)
    [ "${branch}" = "main" ] || { log "ERROR: must be on main (currently ${branch})"; exit 1; }

    [ -z "$(git status --porcelain --untracked-files=no)" ] || { log "ERROR: uncommitted changes"; exit 1; }

    git fetch origin main --tags --quiet
    [ "$(git rev-parse HEAD)" = "$(git rev-parse origin/main)" ] || { log "ERROR: main differs from origin/main"; exit 1; }

    if git rev-parse -q --verify "refs/tags/${version}" > /dev/null; then
        log "ERROR: tag ${version} already exists"
        exit 1
    fi

    log "Running checks (lint, vuln, tests)"
    make lint vuln test-all

    log "Creating tag ${version}"
    git tag -a "${version}" -m "Release ${version}"
    git push origin "${version}"

    log "Tag pushed. GitHub Actions publishes the release."
    log "Next: update CHANGELOG.md (move [Unreleased] to ${version}) if not done yet"
}

main "$@"

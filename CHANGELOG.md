# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Security - COMPLETED 2026-02-23
- **AUDIT COMPLETED** - Comprehensive security audit identified 26 issues
- **FIXES IMPLEMENTED** - All critical and medium security issues resolved

#### Critical Fixes (3/3) ✅
- **Fixed path traversal vulnerability** in filename sanitization (CVE-2026-XXXX)
  - Apply `filepath.Base()` before character replacement
  - Prevents attacks like `"../../../etc/passwd"`
- **Added output path validation** to prevent directory escape (CVE-2026-YYYY)
  - Resolve to absolute paths and verify prefix
  - Ensures all output stays within designated directory
- **Implemented resource limits** to prevent DoS attacks (CVE-2026-ZZZZ)
  - Max 10,000 documents per run
  - Max 10MB per document
  - Max 1GB total output
  - Max 100MB input file

#### Medium Fixes (3/3) ✅
- **Added file type validation** - Reject symlinks and special files before writing
- **Implemented atomic file writes** - Use random temp files with atomic rename
- **Added directory validation** - Prevent symlink attacks on directory creation

### Added
- **Version information** - Now shows Version, GitCommit, and BuildDate
- **Subcommand support** - `help` and `version` subcommands
- **Signal handling** - Graceful shutdown on SIGTERM/SIGINT
- **Comprehensive unit tests** - 12 test functions, 49 subtests (53.6% coverage)
- **Makefile targets** - `test-unit`, `test-all`, `coverage`, `lint`
- **Linting configuration** - `.golangci.yml` for code quality

### Changed
- **Build process** - Static linking with CGO_ENABLED=0
- **Error handling** - Use `%w` for proper error wrapping
- **Exit codes** - Standardized exit codes (0, 1, 2)

### Security
- See `docs/development/SECURITY_AUDIT.md` for full audit details
- See `SECURITY.md` for security policy
- See `FIXPLAN.md` for implementation tracking

### Planned for 0.3.0 (Best Practices - Remaining)
- Make scanner buffer size configurable
- Add cleanup on fatal errors flag
- Complete documentation (configuration.md, troubleshooting.md)

---

## [0.1.4] - 2024-10-21

### Added
- **kubectl-style `-f` parameter** - Input files must now be specified with `-f` flag
- **Stdin support** - Use `-f -` to read from stdin for piping from kubectl, helm, etc.
- **Namespace directory organization** - New `-namespace-dirs` flag organizes output by namespace
- **Resource filtering** - Include/exclude specific resource types with `-include` and `-exclude`
- **Enhanced error handling** - Continue processing on errors with detailed reporting
- **JSON output support** - Export resources as JSON with `-o json` (always sorted)
- **Processing statistics** - Summary of processed/skipped/error counts
- **Comprehensive test suite** - Automated tests covering all features
- **POSIX shell test script** - Cross-platform testing with colored output
- **Professional documentation** - Updated README with real examples using testdata

### Changed
- **BREAKING**: Removed legacy positional arguments - must use `-f` parameter
- **Module path** - Updated from `github.com/mintel/k8s-yaml-splitter` to `github.com/ohauer/k8s-yaml-splitter`
- **File naming** - Better sanitization of special characters in resource names
- **Content preservation** - Original YAML formatting preserved when not sorting
- **Build system** - Enhanced Makefile with help, test, and container targets

### Fixed
- **Line count accuracy** - Fixed YAML processing to preserve original content
- **Memory usage** - Optimized for large Kubernetes manifests
- **Error recovery** - Robust handling of malformed YAML documents

### Technical Details
- Go 1.25 compatibility
- Cross-platform binaries (Linux, FreeBSD, Darwin)
- Container support with binary extraction
- Comprehensive `.gitignore` for development

---

## [0.1.3] - Original Fork Base

### Initial Features (from upstream)
- Basic YAML document splitting
- Simple file output
- Dry-run mode
- Basic command-line interface

**Note**: This version represents a complete rewrite and enhancement of the original k8s-yaml-splitter tool, adding professional-grade features for production Kubernetes workflows.

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed
- `-h` printed no options (flags were defined after the help call) and differed from the usage shown for invalid flags or missing arguments. All paths now print the same complete help: `-h`/`--help`/`help` to stdout with exit 0, errors to stderr with exit 2
- Flags after the output directory were ignored and treated as extra directories (e.g. `-f - -d out -s` failed with a usage error). Flags are now accepted anywhere; `--` ends flag parsing. Error messages name the actual problem (missing `-f`, missing or extra `<output-dir>`)

## [0.2.0-rc3] - 2026-11-19

### Fixed
- Release workflow: create release as draft, upload assets, then publish (required because the repository uses immutable releases; v0.2.0-rc2 was published without binaries and has been withdrawn)


### Changed
- Dockerfile builder image pinned by digest; Dependabot also updates Docker base images
- CI hardened: least-privilege `permissions`, `persist-credentials: false`, concurrency, `workflow_dispatch`, golangci-lint step, coverage gate (55%), fuzz via `make fuzz FUZZTIME=10s` (anchored target names, `-run='^$'`)
- Replaced archived `gopkg.in/yaml.v3` with maintained `go.yaml.in/yaml/v3` v3.0.6 (output verified identical for all testdata in yaml, sorted and JSON modes)
- GitHub Actions updated and pinned by commit SHA (checkout v7, setup-go v7, upload-artifact v7, download-artifact v8, action-gh-release v3)
- Container image runs as non-root (UID 65534); verified build on Go 1.27.4
- Replaced placeholder `CVE-2026-XXXX/YYYY/ZZZZ` with internal audit IDs `SA-01/02/03` (no CVEs were assigned)
- Moved `test.sh` to `scripts/test.sh`; removed `NOTES/` (content consolidated into `docs/development/RELEASING.md`)
- Simplified `.gitignore` (no longer ignores all `*.json` or `*secret*`/`*token*` paths)
- Updated Go to 1.27.4 (go.mod, Dockerfile, CI, CONTRIBUTING)
- Applied `go fix` modernizers (`any`, range-over-int); signal handling now uses `signal.NotifyContext` (finishes current document, prints stats, exits 130)
- Migrated `.golangci.yml` to golangci-lint v2 format; fixed all lint findings (checked `encoder.Close` error)
- Cleaned up documentation references to removed tracking files

### Added
- `.github/dependabot.yml` (monthly: GitHub Actions, Go modules)
- `scripts/release.sh` (checks, tag, push) and `docs/development/RELEASING.md`; release workflow now uses generated release notes and marks `-rc` tags as pre-release
- `.editorconfig`, `.gitattributes`, `.vscode/settings.json`
- `make vuln` and CI step running `govulncheck` (pinned v1.8.0)
- `kind: List` support (e.g. `kubectl get -o yaml` output, `*List` kinds with items): each item is split into its own file; filters, namespace dirs, sorting and JSON output apply per item; nesting limited to 2 levels
- Fuzz test `FuzzProcessObjectList`
- Fuzz tests: `FuzzSplitYAMLDocument`, `FuzzFormatOutput`, `FuzzShouldIncludeResource`, `FuzzGetOutputPath` (also wired into `make fuzz` and CI)

---

## [0.2.0-rc1] - 2026-02-23

### Security - COMPLETED 2026-02-23
- **AUDIT COMPLETED** - Comprehensive security audit identified 26 issues
- **FIXES IMPLEMENTED** - All critical and medium security issues resolved

#### Critical Fixes (3/3) ✅
- **Fixed path traversal vulnerability** in filename sanitization (SA-01)
  - Apply `filepath.Base()` before character replacement
  - Prevents attacks like `"../../../etc/passwd"`
- **Added output path validation** to prevent directory escape (SA-02)
  - Resolve to absolute paths and verify prefix
  - Ensures all output stays within designated directory
- **Implemented resource limits** to prevent DoS attacks (SA-03)
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
- See `docs/development/ROADMAP.md` for remaining work

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

# Project Roadmap

**Project:** k8s-yaml-splitter  
**Last Updated:** 2026-11-19  
**Next Review:** 2026-12-19

---

## Version 0.2.0-rc1 - Security Release ✅ (Completed: 2026-02-23)

### What Was Fixed ✅
**Security (6 fixes):**
- Path traversal in `sanitizeFilename()` - SA-01
- Output path validation - SA-02  
- Resource limits (DoS prevention) - SA-03
- File type validation (reject symlinks/special files)
- Atomic writes with random temp files
- Directory validation

**Testing:**
- 12 test functions, 49 subtests, 53.6% coverage
- 3 fuzz functions (1.4M+ executions, no crashes)
- All tests pass with race detector
- Integration tests pass (7 tests)

**Best Practices:**
- Version variables (Version, GitCommit, BuildDate)
- Static builds (CGO_ENABLED=0)
- Subcommands (help, version)
- Signal handling (SIGTERM, SIGINT)
- Linting config (.golangci.yml)
- Error wrapping (%w)
- Exit code constants

### Release Checklist
- [ ] Manual testing with real Kubernetes YAML
- [ ] Commit all changes
- [ ] Tag as v0.2.0-rc1
- [ ] Create GitHub release with security advisory

**Time:** ~3 hours | **Issues Resolved:** 21/26 (81%)

---

## Unreleased - Toolchain & Features (2026-11-19) ✅

- Go 1.27.4 (go.mod, Dockerfile, CI), `go fix` modernizers applied (`any`, range-over-int)
- `kind: List` / `*List` expansion (one file per item)
- Fuzz tests extended to 8 targets
- Signal handling via `signal.NotifyContext`
- `.golangci.yml` migrated to golangci-lint v2, 0 issues
- `govulncheck` in `make vuln` and CI

---

## Version 0.3.0 - Polish Release (Future)

### Remaining Items
- [ ] Review Go security best practices (https://go.dev/doc/security/best-practices) beyond govulncheck
- [ ] Make scanner buffer size configurable (`-max-buffer-size` flag)
- [ ] Add cleanup-on-error flag (`-cleanup-on-error` flag)
- [ ] Add validate subcommand
- [ ] Create docs/configuration.md
- [ ] Create docs/troubleshooting.md
- [ ] Increase test coverage to 80%+

---

## Version 0.4.0+ - Feature Releases (Future)

Feature ideas (see `01_IDEAS/`):
- `-strip` parameter (remove runtime fields)
- `-merge` mode (combine split files)
- `-validate-k8s` (schema validation)
- Progress bar, verbose/quiet modes
- Label selectors, advanced filtering

---

## Maintenance

**Monthly:** Security review, dependency updates  
**Quarterly:** Full security audit, roadmap review

---

**Last Updated:** 2026-11-19  
**Next Review:** 2026-12-19

# Project Roadmap

**Project:** k8s-yaml-splitter  
**Last Updated:** 2026-02-23  
**Next Review:** 2026-03-23

---

## Version 0.2.0-rc1 - Security Release ✅ (Completed: 2026-02-23)

### What Was Fixed ✅
**Security (6 fixes):**
- Path traversal in `sanitizeFilename()` - CVE-2026-XXXX
- Output path validation - CVE-2026-YYYY  
- Resource limits (DoS prevention) - CVE-2026-ZZZZ
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

## Version 0.3.0 - Polish Release (Future)

### Remaining Items (6 tasks, ~2-3 hours)
- [ ] Make scanner buffer size configurable (`-max-buffer-size` flag)
- [ ] Add cleanup-on-error flag (`-cleanup-on-error` flag)
- [ ] Add validate subcommand
- [ ] Create docs/configuration.md
- [ ] Create docs/troubleshooting.md
- [ ] Increase test coverage to 80%+

---

## Version 0.4.0+ - Feature Releases (Future)

See `01_IDEAS.md` for feature ideas:
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

**Last Updated:** 2026-02-23  
**Next Review:** 2026-03-23

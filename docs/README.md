# Documentation Directory

Documentation for k8s-yaml-splitter.

## Structure

```
README.md                         User guide and examples
CHANGELOG.md                      Version history
SECURITY.md                       Security policy
CONTRIBUTING.md                   Contribution guide
docs/
├── README.md                     (this file)
└── development/
    ├── RELEASING.md              Release process (scripts/release.sh + CI)
    ├── ROADMAP.md                Version planning and remaining work
    └── SECURITY_AUDIT.md         Security audit findings (reference)
```

## Where to find what

| Question | Document |
|----------|----------|
| How do I install and use it? | `../README.md` |
| What changed? | `../CHANGELOG.md` |
| How do I report a vulnerability? | `../SECURITY.md` |
| How do I contribute? | `../CONTRIBUTING.md` |
| What is planned? | `development/ROADMAP.md` |
| What was found in the security audit? | `development/SECURITY_AUDIT.md` |
| How do I make a release? | `development/RELEASING.md` |
| Ideas | `../01_IDEAS/` |

## Make targets for developers

| Target | Purpose |
|--------|---------|
| `make test-unit` | Unit tests with race detector |
| `make test` | Integration tests with built binary |
| `make fuzz` | Fuzz tests (30s each) |
| `scripts/release.sh vX.Y.Z` | Check, tag and push a release |
| `make lint` | golangci-lint v2 |
| `make vuln` | govulncheck |

## Conventions

- Keep tracking minimal: history in `CHANGELOG.md`, future work in `development/ROADMAP.md`.
- Documents with dates carry **Last Updated** and **Next Review** (+1 month).
- Use `example.com` placeholders; never commit real hostnames or credentials.

---

**Last Updated:** 2026-11-19  
**Maintained By:** ohauer

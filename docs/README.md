# Documentation Directory

This directory contains development and user documentation for k8s-yaml-splitter.

---

## Directory Structure

```
docs/
├── README.md                    (this file)
├── configuration.md             [TODO] Detailed configuration reference
├── troubleshooting.md           [TODO] Common issues and solutions
└── development/
    ├── SECURITY_AUDIT.md        Security audit findings and CVE details
    └── ROADMAP.md               Version planning and project roadmap
```

---

## Document Guide

### For Users

**Getting Started:**
1. Read `../README.md` for installation and basic usage
2. Check `configuration.md` for detailed flag documentation (TODO)
3. See `troubleshooting.md` for common issues (TODO)

**Security:**
1. Read `../SECURITY.md` for security policy and best practices
2. Check supported versions and known vulnerabilities
3. Follow secure usage guidelines

**Changes:**
1. Read `../CHANGELOG.md` for version history
2. Check `[Unreleased]` section for upcoming changes

### For Developers

**Understanding the Code:**
1. Read `../README.md` for project overview
2. Review `development/ROADMAP.md` for project direction
3. Check `../CONTRIBUTING.md` for contribution guidelines

**Security:**
1. Read `development/SECURITY_AUDIT.md` for security findings
2. Review `../SECURITY.md` for security requirements
3. Follow security best practices in code reviews

**Implementation:**
1. Start with `../IMPLEMENTATION_CHECKLIST.md` for daily tasks
2. Reference `../PLAN.md` for detailed implementation steps
3. Check `../FIXPLAN.md` for issue tracking and status

**Testing:**
1. Run `make test` for integration tests
2. Run `make test-unit` for unit tests (after implementation)
3. Run `make test-security` for security tests (after implementation)
4. Run `make test-fuzz` for fuzzing tests (after implementation)

### For Project Managers

**Planning:**
1. Review `development/ROADMAP.md` for version milestones
2. Check `../FIXPLAN.md` for timeline and rollout strategy
3. Track progress in `../IMPLEMENTATION_CHECKLIST.md`

**Status:**
1. Check issue status in `../FIXPLAN.md` Issue Inventory
2. Review success criteria in `../PLAN.md`
3. Monitor progress percentage in `../IMPLEMENTATION_CHECKLIST.md`

---

## Document Relationships

### Implementation Flow
```
IMPLEMENTATION_CHECKLIST.md → Daily tasks with time estimates
           ↓
       PLAN.md → Detailed implementation steps with code examples
           ↓
    FIXPLAN.md → Issue tracking and status updates
```

### Security Flow
```
SECURITY_AUDIT.md → Technical security findings (CVEs, test cases)
           ↓
    SECURITY.md → User-facing security policy and guidance
           ↓
   CHANGELOG.md → Security fixes in version history
```

### Planning Flow
```
ROADMAP.md → Long-term version planning and milestones
           ↓
  PLAN.md → Detailed phase-by-phase implementation
           ↓
FIXPLAN.md → Week-by-week rollout strategy
```

---

## Documentation Standards

### File Headers
All documentation files should include:
- **Project name**
- **Last updated date**
- **Next review date** (typically +1 month)
- **Purpose/overview**

### Update Schedule
- **Weekly:** Update progress in IMPLEMENTATION_CHECKLIST.md
- **After each phase:** Update FIXPLAN.md issue status
- **Monthly:** Review and update ROADMAP.md
- **On release:** Update CHANGELOG.md and SECURITY.md

### Markdown Standards
- Use ATX-style headers (`#` not `===`)
- Use fenced code blocks with language tags
- Use tables for structured data
- Use checkboxes for task lists
- Use emoji sparingly (🔴🟡🟢 for severity only)

---

## Maintenance

### When to Update

**SECURITY_AUDIT.md:**
- After completing security fixes
- After external security review
- Quarterly security reviews
- When new vulnerabilities discovered

**ROADMAP.md:**
- Monthly reviews
- After major version releases
- When priorities change
- When new features planned

**configuration.md (TODO):**
- When new flags added
- When flag behavior changes
- When examples need updating

**troubleshooting.md (TODO):**
- When new issues reported
- When solutions found
- When FAQs accumulate

### Review Schedule

| Document | Review Frequency | Next Review |
|----------|------------------|-------------|
| SECURITY_AUDIT.md | Monthly | 2026-03-23 |
| ROADMAP.md | Monthly | 2026-03-23 |
| configuration.md | As needed | N/A |
| troubleshooting.md | As needed | N/A |

---

## Contributing to Documentation

### Adding New Documents
1. Follow file header standards
2. Add to this README.md
3. Update document relationships section
4. Add to appropriate flow diagram

### Updating Existing Documents
1. Update "Last Updated" date
2. Update "Next Review" date if applicable
3. Maintain consistent formatting
4. Update cross-references if needed

### Documentation Review Checklist
- [ ] Headers include all required fields
- [ ] Dates are current
- [ ] Cross-references are valid
- [ ] Code examples are tested
- [ ] Formatting is consistent
- [ ] No trailing whitespace
- [ ] Proper markdown syntax

---

## Quick Links

### Root Documentation
- [README.md](../README.md) - User guide and examples
- [SECURITY.md](../SECURITY.md) - Security policy
- [CHANGELOG.md](../CHANGELOG.md) - Version history
- [CONTRIBUTING.md](../CONTRIBUTING.md) - Contribution guide
- [LICENSE](../LICENSE) - License information

### Planning Documents
- [PLAN.md](../PLAN.md) - Detailed implementation plan
- [FIXPLAN.md](../FIXPLAN.md) - Fix tracking and rollout
- [IMPLEMENTATION_CHECKLIST.md](../IMPLEMENTATION_CHECKLIST.md) - Daily tasks
- [ANALYSIS_SUMMARY.md](../ANALYSIS_SUMMARY.md) - Analysis overview

### Development Documents
- [SECURITY_AUDIT.md](development/SECURITY_AUDIT.md) - Security audit
- [ROADMAP.md](development/ROADMAP.md) - Project roadmap

### Ideas & Notes
- [01_IDEAS.md](../01_IDEAS.md) - Feature ideas
- [NOTES/](../NOTES/) - Development notes

---

**Last Updated:** 2026-02-23  
**Maintained By:** ohauer

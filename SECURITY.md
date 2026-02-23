# Security Policy

## Supported Versions

| Version | Supported          | Status |
| ------- | ------------------ | ------ |
| 0.2.x   | :white_check_mark: | Current (security fixes) |
| 0.1.x   | :x:                | Vulnerable - upgrade immediately |
| < 0.1   | :x:                | Not supported |

## Known Vulnerabilities

### Version 0.1.x and Earlier

**Critical vulnerabilities identified on 2026-02-23:**

1. **Path Traversal (CVE-2026-XXXX)** - HIGH
   - Attackers can write files outside output directory
   - Fixed in: 0.2.0
   - Workaround: Run in isolated environment with limited permissions

2. **Resource Exhaustion (CVE-2026-ZZZZ)** - MEDIUM-HIGH
   - No limits on document count or size
   - Fixed in: 0.2.0
   - Workaround: Validate input files before processing

3. **Symlink Following (CVE-2026-YYYY)** - MEDIUM
   - Can write through symlinks to unintended locations
   - Fixed in: 0.2.0
   - Workaround: Ensure output directory has no symlinks

**Recommendation:** Upgrade to version 0.2.0 or later immediately.

## Reporting a Vulnerability

We take security seriously. If you discover a security vulnerability, please follow these steps:

### 1. Do NOT Open a Public Issue

Security vulnerabilities should not be disclosed publicly until a fix is available.

### 2. Report Privately

**Email:** ohauer@users.noreply.github.com  
**Subject:** [SECURITY] k8s-yaml-splitter vulnerability report

**Include:**
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)
- Your contact information (optional)

### 3. Response Timeline

- **24 hours:** Initial acknowledgment
- **72 hours:** Preliminary assessment
- **7 days:** Detailed response with fix timeline
- **30 days:** Security patch released (target)

### 4. Disclosure Policy

- We will work with you to understand and fix the issue
- We will credit you in the security advisory (unless you prefer anonymity)
- We will coordinate disclosure timing with you
- Public disclosure after fix is released and users have time to upgrade

## Security Best Practices

### For Users

#### Running the Tool Safely

1. **Use Latest Version**
   ```bash
   # Check your version
   k8s-yaml-splitter -version
   
   # Upgrade if needed
   curl -L -o k8s-yaml-splitter https://github.com/ohauer/k8s-yaml-splitter/releases/latest/download/k8s-yaml-splitter-linux-amd64
   ```

2. **Validate Input Files**
   ```bash
   # Check file size before processing
   ls -lh input.yaml
   
   # Validate YAML syntax
   yamllint input.yaml
   ```

3. **Use Isolated Output Directories**
   ```bash
   # Create dedicated output directory
   mkdir -p ~/k8s-splits/output
   k8s-yaml-splitter -f input.yaml ~/k8s-splits/output
   ```

4. **Run with Limited Permissions**
   ```bash
   # Run as non-root user
   # Ensure output directory is owned by current user
   chown -R $USER:$USER ~/k8s-splits
   ```

5. **Use Dry Run First**
   ```bash
   # Preview what will be created
   k8s-yaml-splitter -f input.yaml -dry-run output/
   ```

#### Validating Downloads

Always verify checksums after downloading:

```bash
# Download binary and checksums
curl -L -o k8s-yaml-splitter https://github.com/ohauer/k8s-yaml-splitter/releases/latest/download/k8s-yaml-splitter-linux-amd64
curl -L -o checksums.txt https://github.com/ohauer/k8s-yaml-splitter/releases/latest/download/checksums.txt

# Verify checksum
sha256sum -c checksums.txt --ignore-missing

# Only use if verification passes
chmod +x k8s-yaml-splitter
```

#### Untrusted Input

When processing YAML from untrusted sources:

```bash
# 1. Use dry-run mode first
k8s-yaml-splitter -f untrusted.yaml -dry-run /tmp/test-output

# 2. Review what would be created
# 3. Use isolated directory
mkdir -p /tmp/isolated-output
k8s-yaml-splitter -f untrusted.yaml /tmp/isolated-output

# 4. Review output before using
ls -la /tmp/isolated-output
```

### For Developers

#### Secure Development

1. **Run Security Tests**
   ```bash
   make test-security
   make test-fuzz
   ```

2. **Use Race Detector**
   ```bash
   go test -race ./...
   ```

3. **Run Linter**
   ```bash
   make lint
   ```

4. **Review Security Checklist**
   - Path validation for all file operations
   - Resource limits enforced
   - No symlink following
   - Atomic file writes
   - Input validation

#### Code Review Checklist

Before merging code:
- [ ] All tests pass (including `-race`)
- [ ] Linter passes
- [ ] Security implications reviewed
- [ ] No new path operations without validation
- [ ] No new file operations without type checking
- [ ] Resource limits considered
- [ ] Error handling consistent
- [ ] Documentation updated

## Security Features

### Current (Version 0.2.0+)

- ✅ Path traversal prevention
- ✅ Output path validation
- ✅ Resource limits (documents, size, total output)
- ✅ File type validation (reject symlinks/special files)
- ✅ Atomic file writes
- ✅ Directory validation
- ✅ Input file validation

### Planned (Future Versions)

- ⏳ Signature verification for binaries
- ⏳ SBOM (Software Bill of Materials)
- ⏳ Audit logging
- ⏳ Policy enforcement
- ⏳ Sandboxing options

## Threat Model

### Assets
- **Output files:** Kubernetes manifests (may contain sensitive data)
- **Filesystem:** User's filesystem integrity
- **System resources:** Memory, disk, CPU

### Threats

#### 1. Malicious Input Files
**Threat:** Attacker provides crafted YAML to exploit vulnerabilities

**Mitigations:**
- Path traversal prevention
- Resource limits
- Input validation
- Dry-run mode

#### 2. Compromised Output Directory
**Threat:** Attacker pre-creates symlinks in output directory

**Mitigations:**
- File type validation
- Directory validation
- Atomic writes with random temp files

#### 3. Resource Exhaustion
**Threat:** Attacker provides huge files to exhaust resources

**Mitigations:**
- Document count limits
- Document size limits
- Total output size limits
- Input file size limits

#### 4. Race Conditions
**Threat:** Attacker exploits TOCTOU vulnerabilities

**Mitigations:**
- Atomic file writes
- Random temp filenames
- File type validation before write

### Trust Boundaries

```
┌─────────────────────────────────────────┐
│ User (Trusted)                          │
│  - Provides input file                  │
│  - Specifies output directory           │
│  - Runs tool with permissions           │
└─────────────────┬───────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────┐
│ k8s-yaml-splitter (Trust Boundary)      │
│  - Validates all inputs                 │
│  - Enforces resource limits             │
│  - Prevents path traversal              │
│  - Validates file types                 │
└─────────────────┬───────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────┐
│ Filesystem (Untrusted)                  │
│  - May contain symlinks                 │
│  - May have special files               │
│  - May have permission issues           │
└─────────────────────────────────────────┘
```

### Assumptions

**We assume:**
- User has legitimate access to input files
- User has write permissions to output directory
- Operating system is not compromised
- Go runtime is not compromised

**We do NOT assume:**
- Input YAML is well-formed or safe
- Output directory is empty or safe
- Filesystem paths are trustworthy
- Resource availability is unlimited

## Security Contacts

**Primary:** ohauer@users.noreply.github.com  
**Response Time:** 24-72 hours  
**PGP Key:** Not yet available

## Acknowledgments

We thank the following individuals for responsible disclosure:

- *No reports yet*

---

**Last Updated:** 2026-02-23  
**Next Review:** 2026-03-23

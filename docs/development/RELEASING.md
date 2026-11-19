# Releasing

**Last Updated:** 2026-11-19  
**Next Review:** 2026-12-19

Releases are published by GitHub Actions when a version tag is pushed.

## Steps

1. Merge all changes to `main` via pull request (the PR titles become the generated release notes).
2. Move `[Unreleased]` in `CHANGELOG.md` to the new version and commit via PR.
3. On `main`, run:
   ```console
   scripts/release.sh v0.2.0
   ```
   The script checks branch, clean tree and sync with `origin/main`, runs `make lint vuln test-all`, then creates and pushes an annotated tag.
4. The workflow `.github/workflows/go.yml` builds static binaries (linux/darwin/freebsd amd64), creates `checksums.txt` and publishes the GitHub release with generated notes. Tags containing `-rc` are marked as pre-release.

## Notes

- Use semantic versioning: `vMAJOR.MINOR.PATCH`, release candidates `vX.Y.Z-rcN`.
- Label PRs (e.g. `security`) and configure `.github/release.yml` to group entries in the notes.
- For a security release, add a short hand-written summary to the release body on GitHub after publishing.
- Verify downloads with `sha256sum -c checksums.txt`.

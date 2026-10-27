# Security policy

Canonical policy on the docs site: [English](https://docs.gov.skin/en/community/security/) · [中文](https://docs.gov.skin/community/security/)

**Do not disclose security vulnerabilities in public issues.** Use "Private vulnerability reporting" under this repository's Security tab to submit a private security advisory, or email the maintainer at <team@tsinbei.com>.

## Repository-specific dependency baseline

- Currently **zero third-party dependencies** (pure standard library, no `go.sum`).
- New dependencies must pass `govulncheck` in CI; if `klauspost/compress` is pulled in transitively it must be ≥ v1.18.7.

# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |

## Reporting a Vulnerability

If you discover a security vulnerability in this project, please report it responsibly:

1. **Do not** open a public GitHub issue for security vulnerabilities
2. Email the maintainer directly or use [GitHub's private vulnerability reporting](https://github.com/hackIDLE/fedramp-browser/security/advisories/new)
3. Include as much detail as possible:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

## Response Timeline

- **Acknowledgment**: Within 48 hours
- **Initial Assessment**: Within 7 days
- **Resolution Target**: Within 30 days for critical issues

## Security Measures

This project implements the following security practices:

- **Dependency Scanning**: OSV-Scanner runs on every PR and weekly
- **Static Analysis**: CodeQL scans for security vulnerabilities
- **Supply Chain Security**: All GitHub Actions are pinned to SHA hashes
- **SBOM Generation**: Software Bill of Materials included with releases

## Scope

This security policy applies to:
- The fedramp-browser CLI application
- GitHub Actions workflows in this repository
- Release artifacts published to GitHub Releases

This policy does not cover:
- Third-party dependencies (report to upstream maintainers)
- The FedRAMP data source (report to FedRAMP PMO)

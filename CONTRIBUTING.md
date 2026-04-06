# Contributing to FedRAMP Browser

Thank you for your interest in contributing to FedRAMP Browser!

## How to Contribute

### Reporting Bugs

1. Check [existing issues](https://github.com/hackIDLE/fedramp-browser/issues) to avoid duplicates
2. Open a new issue with:
   - Clear, descriptive title
   - Steps to reproduce
   - Expected vs actual behavior
   - Go version and OS

### Suggesting Enhancements

Open an issue describing:
- The enhancement and its use case
- Why it would be useful
- Any implementation ideas

### Pull Requests

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Make your changes
4. Run tests and linting (see below)
5. Commit with clear messages
6. Push to your fork
7. Open a Pull Request

## Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/fedramp-browser.git
cd fedramp-browser

# Build
go build ./...

# Run tests
go test ./...

# Run linter
golangci-lint run
```

## Code Standards

### Style
- Follow standard Go conventions
- Run `golangci-lint run` before committing
- Keep functions focused and small

### Testing
- Add tests for new functionality
- Ensure existing tests pass
- Run `go test ./...` before submitting

### Commits
- Use clear, descriptive commit messages
- Reference issues when applicable (e.g., "Fixes #123")

## Code of Conduct

Be respectful and constructive in all interactions.

## Questions?

Open an issue or start a discussion.

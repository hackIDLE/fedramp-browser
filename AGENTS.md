# Repository Guidelines

## Project Structure & Module Organization
- `main.go`: CLI entrypoint; parses `--refresh` and launches the Bubble Tea app.
- `internal/tui`: UI state, key bindings, list/detail views, and styles.
- `internal/api`: fetches and parses `FRMR.documentation.json` from `FedRAMP/docs`.
- `internal/model`: shared data types (`Document`, `Requirement`, `Definition`, `Indicator`).
- `internal/cache`: local file cache (`~/.cache/fedramp-browser`, 24-hour TTL).
- Tests live next to code as `*_test.go`; fuzz tests are in `internal/api/fuzz_test.go`.
- Automation is in `.github/workflows/` (test, lint, release, security, upstream checks).

## Build, Test, and Development Commands
- `go build ./...`: build all packages.
- `go build -o fedramp-browser .`: build a local executable from repo root.
- `go run .`: run the TUI locally.
- `go run . --refresh`: ignore cache and force fresh upstream fetch.
- `go test ./... -count=1 -timeout 120s`: run full tests without cached results.
- `go vet ./...`: run built-in static analysis.
- `golangci-lint run`: run configured linters (`errcheck`, `gosec`, `staticcheck`, etc.).

## Coding Style & Naming Conventions
- Follow idiomatic Go; keep packages small and purpose-driven.
- Use `gofmt` defaults for formatting (tabs, standard spacing, no manual alignment).
- Use clear exported names (`WithRefresh`, `ParseRequirements`) and lowercase package names.
- Prefer small functions, explicit error handling, and minimal side effects.

## Testing Guidelines
- Keep tests adjacent to implementation files.
- Name tests `TestXxx`; name fuzz tests `FuzzXxx` (existing project pattern).
- Add or update tests for parser changes, filter logic, and view-state transitions.
- Run `go test ./...` and `golangci-lint run` before opening a PR.

## Commit & Pull Request Guidelines
- Follow the repo’s Conventional Commit style: `feat:`, `fix:`, `ci:`, `deps:`, `docs:`, `test:`, `security:`.
- Keep each commit focused on one logical change.
- PRs should include:
  - What changed and why
  - Linked issue(s) when applicable (`Fixes #123`)
  - Verification steps/results (build, test, lint)
  - Screenshots or GIFs for TUI behavior changes

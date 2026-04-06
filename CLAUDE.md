# CLAUDE.md

## Build & Test
- `go build ./...` — build all packages
- `go test ./... -count=1 -timeout 120s` — run all tests (no caching)
- `go vet ./...` — static analysis

## Dependency Management
- Go stdlib vulnerabilities: bump `go` directive in `go.mod`, then `go mod tidy`
- Verify with: `osv-scanner scan --lockfile go.mod` (install: `go install github.com/google/osv-scanner/v2/cmd/osv-scanner@latest`)

## CI Workflows
- **Test** — runs `go test`
- **OSV Scanner** — checks for known vulnerabilities in dependencies and Go stdlib
- **OpenSSF Scorecard** — security best practices score
- **Upstream Check** — monitors upstream FedRAMP data changes
- **Dependency Graph** — GitHub dependency tracking

## Release Process
- Versioning: semver (patch for fixes/security, minor for features)
- Create release: `gh release create v0.x.y --title "v0.x.y" --notes "..." --target main`
- GoReleaser workflow triggers automatically on tag push
- Produces 13 assets: binaries for darwin/linux/windows × amd64/arm64, each with SBOM, plus checksums.txt
- Verify: `gh release view v0.x.y` to confirm all assets attached

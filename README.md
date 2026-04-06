# FedRAMP Browser

[![Go Report Card](https://goreportcard.com/badge/github.com/hackIDLE/fedramp-browser)](https://goreportcard.com/report/github.com/hackIDLE/fedramp-browser)
[![Release](https://img.shields.io/github/v/release/hackIDLE/fedramp-browser)](https://github.com/hackIDLE/fedramp-browser/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/hackIDLE/fedramp-browser/badge)](https://scorecard.dev/viewer/?uri=github.com/hackIDLE/fedramp-browser)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/11619/badge)](https://www.bestpractices.dev/projects/11619)

A terminal user interface for browsing FedRAMP documentation.

## Demo

![FedRAMP TUI demo](demo.gif)

## Features

- **Document Navigator**: Browse all 12 FedRAMP document categories
- **Requirements Search**: Search and filter requirements across all documents
- **Definitions Lookup**: Quick access to FedRAMP terminology
- **Key Security Indicators**: View KSI themes with SP 800-53 control mappings

## Installation

### Homebrew (macOS/Linux)

```bash
brew install hackIDLE/tap/fedramp-browser
```

### Scoop (Windows)

```powershell
scoop bucket add hackidle https://github.com/hackIDLE/scoop-bucket
scoop install fedramp-browser
```

### Download Binary

Download from [Releases](https://github.com/hackIDLE/fedramp-browser/releases):

| Platform | Binary |
|----------|--------|
| macOS (Apple Silicon) | `fedramp-browser-darwin-arm64` |
| macOS (Intel) | `fedramp-browser-darwin-amd64` |
| Linux (x64) | `fedramp-browser-linux-amd64` |
| Linux (ARM64) | `fedramp-browser-linux-arm64` |
| Windows (x64) | `fedramp-browser-windows-amd64.exe` |

### Go Install

```bash
go install github.com/hackIDLE/fedramp-browser@latest
```

### Build from Source

```bash
git clone https://github.com/hackIDLE/fedramp-browser.git
cd fedramp-browser
go build -o fedramp-browser .
```

## Usage

```bash
fedramp-browser
```

### Command Line Options

| Flag | Description |
|------|-------------|
| `--refresh` | Force fresh fetch from GitHub, ignoring cache |

### Caching

Data is cached locally at `~/.cache/fedramp-browser/` with a 24-hour TTL. On subsequent runs, the TUI loads instantly from cache. Use `--refresh` to force a fresh fetch.

### Key Bindings

| Key | Action |
|-----|--------|
| `1` | View Documents |
| `2` | View Requirements |
| `3` | View Definitions |
| `4` | View Key Security Indicators |
| `j/k` or `↑/↓` | Navigate list |
| `Enter` | View details |
| `Esc` or `Backspace` | Go back |
| `/` | Filter/search |
| `m` | Filter MUST requirements (Requirements view) |
| `s` | Filter SHOULD requirements (Requirements view) |
| `x` | Cycle affects filter: All → Providers → Agencies → Assessors → FedRAMP (Requirements view) |
| `f` | Clear filters (Requirements view) |
| `q` | Quit |

## Data Sources

Data is fetched directly from the [FedRAMP/docs](https://github.com/FedRAMP/docs) GitHub repository.

### Documents

| Code | Name |
|------|------|
| FRD | FedRAMP Definitions |
| KSI | Key Security Indicators |
| VDR | Vulnerability Detection & Response |
| UCM | Using Cryptographic Modules |
| RSC | Recommended Secure Configuration |
| ADS | Authorization Data Sharing |
| CCM | Collaborative Continuous Monitoring |
| FSI | FedRAMP Security Inbox |
| ICP | Incident Communications Procedures |
| MAS | Minimum Assessment Scope |
| PVA | Persistent Validation & Assessment |
| SCN | Significant Change Notifications |

## Acknowledgments

- [Charm](https://charm.sh/) for the TUI libraries: [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Lip Gloss](https://github.com/charmbracelet/lipgloss), [Bubbles](https://github.com/charmbracelet/bubbles)
- [GoReleaser](https://goreleaser.com/) for release automation
- [FedRAMP](https://github.com/FedRAMP/docs) for the source documentation

## Disclaimer

This is not an official FedRAMP tool and I'm not officially associated with the GSA/FedRAMP PMO.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

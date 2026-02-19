# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.1] - 2026-02-20

### Added
- TUI test suite: model init, key handling, navigation, view rendering, edge cases
- Makefile with ldflags for local builds (`make build`, `make install`)

### Fixed
- Version command now shows correct version, commit, and build date
- All lint errors resolved (unused vars, unchecked returns, gosimple)
- Cross-platform test compatibility (darwin-specific tests behind build tags)
- Go version compatibility with golangci-lint (go 1.23)

### Changed
- BACKLOG.md updated with completed items checked off

## [0.1.0] - 2026-02-19

### Added
- Interactive TUI dashboard with Bubble Tea
- Port scanning on macOS (lsof) and Linux (ss)
- Process details (CPU, memory, command, start time)
- One-key kill with confirmation dialog
- Port conflict detection with color coding
- Search and filter by port or process name
- Service groups via `~/.portpilot.yaml` config
- CLI commands: `list`, `kill`, `check`, `watch`, `version`
- JSON/table output formats
- Live auto-refresh (configurable interval)
- Sortable columns (press 1-8)
- Help overlay (press ?)
- GoReleaser for cross-platform binary releases
- GitHub Actions CI (lint, test, build on macOS + Ubuntu)
- Comprehensive README, CONTRIBUTING.md, ARCHITECTURE.md
- Issue templates for bugs and feature requests

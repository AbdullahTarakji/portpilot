# BACKLOG — portpilot

## Current Sprint: P0 — Full Feature Build ✅

### Core Engine
- [x] PP-01: Port scanner — list all listening ports with PID, process name, protocol (TCP/UDP), state
- [x] PP-02: Process details — show command, user, CPU/memory usage per port
- [x] PP-03: Process kill — kill any process by port (with confirmation)
- [x] PP-04: Port conflict detection — detect multiple processes binding same port
- [x] PP-05: Cross-platform support — macOS (lsof) + Linux (ss/netstat)

### TUI Dashboard
- [x] PP-06: Main table view — sortable columns (port, PID, process, protocol, state, CPU, memory)
- [x] PP-07: Live refresh — auto-update every 1-2 seconds
- [x] PP-08: Search/filter — filter by port number, process name, or state
- [x] PP-09: One-key kill — press 'k' to kill selected process (with confirm dialog)
- [x] PP-10: Detail panel — press Enter to see full process info
- [x] PP-11: Color coding — green for normal, red for conflicts, yellow for high resource usage
- [x] PP-12: Help overlay — press '?' for keybindings

### CLI Mode
- [x] PP-13: `portpilot list` — non-interactive list of ports (for scripting/piping)
- [x] PP-14: `portpilot kill <port>` — kill process on port directly
- [x] PP-15: `portpilot check <port>` — check if port is in use (exit code for scripts)
- [x] PP-16: `portpilot watch` — watch mode with streaming output
- [x] PP-17: JSON/table output formats for CLI mode

### Service Groups
- [x] PP-18: Config file (~/.portpilot.yaml) for service groups
- [x] PP-19: Tag ports as "frontend", "backend", "database", etc.
- [x] PP-20: Group view in TUI — toggle between all ports and grouped view

### Polish
- [x] PP-21: Proper error handling (permission errors, OS-specific edge cases)
- [x] PP-22: Installation support — go install, brew formula, binary releases
- [x] PP-23: Man page / --help with examples

## P1 — Documentation & CI ✅
- [x] PP-24: README.md — full docs with GIF demo, install, usage, config
- [x] PP-25: CONTRIBUTING.md
- [x] PP-26: CHANGELOG.md
- [x] PP-27: GitHub Actions CI — lint, test, build, release
- [x] PP-28: GitHub issue templates (bug report + feature request)
- [x] PP-29: docs/ARCHITECTURE.md
- [x] PP-30: goreleaser config for cross-platform binaries

## P2 — Nice to Have
- [ ] PP-31: Docker container port mapping display
- [ ] PP-32: Port forwarding shortcuts
- [ ] PP-33: Notification when a watched port becomes available
- [ ] PP-34: Homebrew tap

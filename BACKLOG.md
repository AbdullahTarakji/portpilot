# BACKLOG — portpilot

## Current Sprint: P0 — Full Feature Build

### Core Engine
- [ ] PP-01: Port scanner — list all listening ports with PID, process name, protocol (TCP/UDP), state
- [ ] PP-02: Process details — show command, user, CPU/memory usage per port
- [ ] PP-03: Process kill — kill any process by port (with confirmation)
- [ ] PP-04: Port conflict detection — detect multiple processes binding same port
- [ ] PP-05: Cross-platform support — macOS (lsof) + Linux (ss/netstat)

### TUI Dashboard
- [ ] PP-06: Main table view — sortable columns (port, PID, process, protocol, state, CPU, memory)
- [ ] PP-07: Live refresh — auto-update every 1-2 seconds
- [ ] PP-08: Search/filter — filter by port number, process name, or state
- [ ] PP-09: One-key kill — press 'k' to kill selected process (with confirm dialog)
- [ ] PP-10: Detail panel — press Enter to see full process info
- [ ] PP-11: Color coding — green for normal, red for conflicts, yellow for high resource usage
- [ ] PP-12: Help overlay — press '?' for keybindings

### CLI Mode
- [ ] PP-13: `portpilot list` — non-interactive list of ports (for scripting/piping)
- [ ] PP-14: `portpilot kill <port>` — kill process on port directly
- [ ] PP-15: `portpilot check <port>` — check if port is in use (exit code for scripts)
- [ ] PP-16: `portpilot watch` — watch mode with streaming output
- [ ] PP-17: JSON/table output formats for CLI mode

### Service Groups
- [ ] PP-18: Config file (~/.portpilot.yaml) for service groups
- [ ] PP-19: Tag ports as "frontend", "backend", "database", etc.
- [ ] PP-20: Group view in TUI — toggle between all ports and grouped view

### Polish
- [ ] PP-21: Proper error handling (permission errors, OS-specific edge cases)
- [ ] PP-22: Installation support — go install, brew formula, binary releases
- [ ] PP-23: Man page / --help with examples

## P1 — Documentation & CI
- [ ] PP-24: README.md — full docs with GIF demo, install, usage, config
- [ ] PP-25: CONTRIBUTING.md
- [ ] PP-26: CHANGELOG.md
- [ ] PP-27: GitHub Actions CI — lint, test, build, release
- [ ] PP-28: GitHub issue templates (bug report + feature request)
- [ ] PP-29: docs/ARCHITECTURE.md
- [ ] PP-30: goreleaser config for cross-platform binaries

## P2 — Nice to Have
- [ ] PP-31: Docker container port mapping display
- [ ] PP-32: Port forwarding shortcuts
- [ ] PP-33: Notification when a watched port becomes available
- [ ] PP-34: Homebrew tap

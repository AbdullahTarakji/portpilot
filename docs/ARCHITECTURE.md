# Architecture

## Overview

PortPilot is a Go CLI/TUI application built with a clean separation of concerns:

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   cmd/       │────▶│  internal/   │────▶│  OS Tools    │
│  (Cobra CLI) │     │  scanner/    │     │  (lsof/ss)   │
│              │     │  process/    │     │              │
│              │     │  config/     │     │              │
│              │────▶│  tui/        │     │              │
│              │     │  (Bubble Tea)│     │              │
└──────────────┘     └──────────────┘     └──────────────┘
```

## Components

### Scanner (`internal/scanner/`)
Platform-specific port discovery using native OS tools.

- **Interface:** `Scanner` with `Scan() ([]PortInfo, error)`
- **macOS:** Parses `lsof -iTCP -iUDP -nP -sTCP:LISTEN`
- **Linux:** Parses `ss -tulnp`
- **Enrichment:** Gets CPU/memory via `ps -p <pid> -o %cpu,%mem,lstart,command`
- Uses Go build tags (`//go:build darwin`, `//go:build linux`) for platform dispatch

### Process Manager (`internal/process/`)
Process lifecycle operations — primarily killing processes with configurable signals.

- `Kill(pid int, sig os.Signal) error`
- `ParseSignal(name string) (os.Signal, error)` — maps SIGTERM, SIGKILL, etc.
- Safety: Never kills PID 0 or 1

### TUI (`internal/tui/`)
Interactive terminal dashboard using the Elm architecture (Bubble Tea).

- **Model:** Holds state (ports list, selected row, filter text, view mode)
- **Update:** Handles key events, tick events, scan results
- **View:** Renders table, detail panel, help overlay
- Auto-refreshes via `tea.Tick` every N seconds

### Config (`internal/config/`)
Optional YAML configuration from `~/.portpilot.yaml`.

- Service groups with port assignments and colors
- Refresh interval
- System port visibility toggle
- Graceful fallback to defaults when no config exists

## Data Flow

1. **Scan:** Scanner runs OS command → parses output → returns `[]PortInfo`
2. **Enrich:** Each port's PID gets CPU/memory stats via `ps`
3. **Display:** TUI renders the port list with sorting/filtering/coloring
4. **Action:** User can kill processes, which sends signals via `process.Kill()`

## Platform Support

| Feature | macOS | Linux |
|---------|-------|-------|
| Port scan | `lsof` | `ss` |
| Process stats | `ps` | `ps` |
| Kill | `syscall.Kill` | `syscall.Kill` |
| TUI | ✅ | ✅ |

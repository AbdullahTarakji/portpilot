<p align="center">
  <h1 align="center">🚀 PortPilot</h1>
  <p align="center">
    <strong>A beautiful TUI + CLI for managing ports, processes, and dev services</strong>
  </p>
  <p align="center">
    <a href="https://github.com/AbdullahTarakji/portpilot/actions/workflows/ci.yml"><img src="https://github.com/AbdullahTarakji/portpilot/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
    <a href="https://github.com/AbdullahTarakji/portpilot/releases"><img src="https://img.shields.io/github/v/release/AbdullahTarakji/portpilot" alt="Release"></a>
    <a href="https://github.com/AbdullahTarakji/portpilot/blob/main/LICENSE"><img src="https://img.shields.io/github/license/AbdullahTarakji/portpilot" alt="License"></a>
    <a href="https://goreportcard.com/report/github.com/AbdullahTarakji/portpilot"><img src="https://goreportcard.com/badge/github.com/AbdullahTarakji/portpilot" alt="Go Report Card"></a>
  </p>
</p>

---

**Ever typed `lsof -i :3000 | grep LISTEN` for the hundredth time?**

PortPilot gives you a single command to see everything running on your ports — with a beautiful interactive TUI, one-key process killing, conflict detection, and more.

## ✨ Features

- 📊 **Interactive TUI** — Real-time dashboard of all listening ports
- 🔍 **Search & Filter** — Find ports by number or process name instantly
- ⚡ **One-Key Kill** — Select a process, press `k`, confirm, done
- 🚨 **Conflict Detection** — Highlights when multiple processes fight for the same port
- 🎨 **Color Coded** — Red for conflicts, yellow for high resource usage, green for normal
- 📋 **CLI Mode** — Scriptable commands for automation (`list`, `kill`, `check`, `watch`)
- 🏷️ **Service Groups** — Tag ports as "frontend", "backend", "database" via config
- 🔄 **Live Refresh** — Auto-updates every 2 seconds
- 📦 **JSON Output** — Pipe to `jq`, scripts, or other tools
- 🖥️ **Cross-Platform** — macOS and Linux support

## 📦 Installation

### Go Install
```bash
go install github.com/AbdullahTarakji/portpilot/cmd/portpilot@latest
```

### From Source
```bash
git clone https://github.com/AbdullahTarakji/portpilot.git
cd portpilot
go build -o portpilot ./cmd/portpilot
sudo mv portpilot /usr/local/bin/
```

### Binary Releases
Download pre-built binaries from [Releases](https://github.com/AbdullahTarakji/portpilot/releases).

## 🚀 Quick Start

```bash
# Launch the interactive TUI
portpilot

# List all listening ports
portpilot list

# Check if port 3000 is in use
portpilot check 3000

# Kill whatever is on port 8080
portpilot kill 8080
```

## 📖 Usage

### Interactive TUI

Just run `portpilot` with no arguments:

```
$ portpilot
```

This opens a real-time dashboard showing all listening ports:

```
🚀 PortPilot — mike@macbook — 8 ports — 12 connections

 PORT   PROTO  PID    PROCESS     USER  CPU%   MEM%  STATE
 ────   ─────  ───    ───────     ────  ────   ────  ─────
 3000   TCP    12345  node        mike   2.1    1.3  LISTEN
 3001   TCP    12346  node        mike   0.5    0.8  LISTEN
 5173   TCP    12400  vite        mike   1.2    0.9  LISTEN
 5432   TCP    3125   postgres    mike   0.0    0.1  LISTEN
 6379   TCP    2882   redis-ser   mike   0.0    0.0  LISTEN
 8080   TCP    14500  Python      mike   0.1    0.2  LISTEN
 27017  TCP    9800   mongod      mike   0.3    2.1  LISTEN

 🔍 Filter: _                    Last refresh: 20:15:03
 [k]ill  [/]filter  [Enter]details  [g]roups  [?]help  [q]uit
```

#### TUI Keybindings

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate rows |
| `Enter` | View process details |
| `k` | Kill selected process |
| `/` | Enter search/filter mode |
| `g` | Toggle service group view |
| `1`-`8` | Sort by column |
| `r` | Force refresh |
| `?` | Show help overlay |
| `q` / `Ctrl+C` | Quit |

### CLI Commands

#### `portpilot list` — List Ports

```bash
# Table output (default)
portpilot list

# JSON output
portpilot list --json

# Filter by port
portpilot list --port 3000

# Filter by process name
portpilot list --process node
```

Example output:
```
$ portpilot list
PORT   PROTO  PID    PROCESS     USER  CPU%  MEM%  STATE
----   -----  ---    -------     ----  ----  ----  -----
3000   TCP    12345  node        mike  2.1   1.3   LISTEN
5432   TCP    3125   postgres    mike  0.0   0.1   LISTEN
6379   TCP    2882   redis-ser   mike  0.0   0.0   LISTEN
```

#### `portpilot kill <port>` — Kill Process

```bash
# Kill with confirmation prompt
portpilot kill 3000
# > Kill "node" (PID 12345) on port 3000? [y/N]

# Force kill (skip confirmation)
portpilot kill 3000 --force

# Send specific signal
portpilot kill 3000 --signal SIGKILL
```

#### `portpilot check <port>` — Check Port Availability

```bash
# Check if port is free (exit code 0 = free, 1 = in use)
portpilot check 3000
# > Port 3000 is in use by "node" (PID 12345, TCP)

portpilot check 9999
# > Port 9999 is free

# Use in scripts
if portpilot check 3000 2>/dev/null; then
  echo "Port 3000 is available"
else
  echo "Port 3000 is taken!"
fi
```

#### `portpilot watch` — Watch Mode

```bash
# Watch all ports (refreshes every 2s)
portpilot watch

# Watch specific port
portpilot watch --port 3000

# Custom refresh interval
portpilot watch --interval 5
```

## ⚙️ Configuration

Create `~/.portpilot.yaml` to customize behavior:

```yaml
# Group ports by service type
groups:
  frontend:
    ports: [3000, 3001, 5173, 8080]
    color: blue
  backend:
    ports: [4000, 8000, 9000]
    color: green
  database:
    ports: [5432, 3306, 27017, 6379]
    color: yellow

# Auto-refresh interval in seconds (default: 2)
refresh_interval: 2

# Show system/root ports (default: false)
show_system_ports: false
```

Press `g` in the TUI to toggle the group view, which labels ports by their service group.

## 🏗️ Tech Stack

- **Language:** [Go](https://go.dev/) — Fast, cross-platform, single binary
- **TUI Framework:** [Bubble Tea](https://github.com/charmbracelet/bubbletea) — Elm-architecture TUI
- **Styling:** [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Terminal CSS
- **CLI Framework:** [Cobra](https://github.com/spf13/cobra) — Industry-standard Go CLI
- **Port Scanning:** `lsof` (macOS) / `ss` (Linux) — Native OS tools, no root required

## 📁 Project Structure

```
portpilot/
├── cmd/portpilot/
│   └── main.go              # Entry point, Cobra commands
├── internal/
│   ├── scanner/
│   │   ├── types.go          # PortInfo struct
│   │   ├── scanner.go        # Scanner interface + shared utils
│   │   ├── darwin.go          # macOS scanner (lsof)
│   │   ├── linux.go           # Linux scanner (ss)
│   │   └── scanner_test.go    # Scanner tests
│   ├── tui/
│   │   ├── app.go             # Main TUI model (Bubble Tea)
│   │   ├── table.go           # Port table component
│   │   ├── detail.go          # Process detail panel
│   │   ├── help.go            # Help overlay
│   │   └── styles.go          # Lip Gloss styles
│   ├── process/
│   │   ├── process.go         # Kill, signal handling
│   │   └── process_test.go    # Process tests
│   └── config/
│       ├── config.go          # YAML config parsing
│       └── config_test.go     # Config tests
├── .github/
│   ├── workflows/
│   │   ├── ci.yml             # Lint + test + build
│   │   └── release.yml        # GoReleaser on tags
│   └── ISSUE_TEMPLATE/
├── .goreleaser.yaml
├── BACKLOG.md
├── CHANGELOG.md
├── CONTRIBUTING.md
├── LICENSE
└── README.md
```

## 🤝 Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

1. Fork the repo
2. Create a feature branch from `develop`
3. Make your changes with tests
4. Submit a PR

## 📄 License

[MIT](LICENSE) — use it however you want.

## 🌟 Star History

If PortPilot saves you time, consider giving it a ⭐!

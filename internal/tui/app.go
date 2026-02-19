// Package tui implements the interactive terminal UI for portpilot.
package tui

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/AbdullahTarakji/portpilot/internal/config"
	"github.com/AbdullahTarakji/portpilot/internal/process"
	"github.com/AbdullahTarakji/portpilot/internal/scanner"
)

// viewMode tracks the current UI state.
type viewMode int

const (
	viewTable viewMode = iota
	viewDetail
	viewHelp
	viewConfirmKill
)

// Model is the main bubbletea model for the TUI.
type Model struct {
	ports       []scanner.PortInfo
	scanner     scanner.Scanner
	config      *config.Config
	width       int
	height      int
	cursor      int
	sortCol     sortOrder
	filter      string
	filterMode  bool
	view        viewMode
	showGroups  bool
	lastRefresh time.Time
	statusMsg   string
	err         error
	hostname    string
}

type tickMsg time.Time

type scanResultMsg struct {
	ports []scanner.PortInfo
	err   error
}

// New creates a new TUI model.
func New(s scanner.Scanner, cfg *config.Config) Model {
	hostname, _ := os.Hostname()
	return Model{
		scanner:  s,
		config:   cfg,
		sortCol:  sortOrder{column: 0, asc: true},
		hostname: hostname,
	}
}

// Run starts the TUI application.
func Run(s scanner.Scanner, cfg *config.Config) error {
	m := New(s, cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func doScan(s scanner.Scanner) tea.Cmd {
	return func() tea.Msg {
		ports, err := s.Scan()
		return scanResultMsg{ports: ports, err: err}
	}
}

func tickCmd(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		doScan(m.scanner),
		tickCmd(time.Duration(m.config.RefreshInterval)*time.Second),
	)
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		return m, tea.Batch(
			doScan(m.scanner),
			tickCmd(time.Duration(m.config.RefreshInterval)*time.Second),
		)

	case scanResultMsg:
		if msg.err != nil {
			m.err = msg.err
			m.statusMsg = fmt.Sprintf("Scan error: %v", msg.err)
		} else {
			m.ports = msg.ports
			m.lastRefresh = time.Now()
			m.err = nil
			// Ensure cursor is in bounds
			filtered := filterPorts(m.ports, m.filter)
			if m.cursor >= len(filtered) {
				m.cursor = max(0, len(filtered)-1)
			}
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global keys
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	}

	switch m.view {
	case viewConfirmKill:
		return m.handleConfirmKillKey(msg)
	case viewHelp:
		return m.handleHelpKey(msg)
	case viewDetail:
		return m.handleDetailKey(msg)
	default:
		if m.filterMode {
			return m.handleFilterKey(msg)
		}
		return m.handleTableKey(msg)
	}
}

func (m Model) handleTableKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	filtered := filterPorts(m.ports, m.filter)

	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "?":
		m.view = viewHelp
		return m, nil
	case "/":
		m.filterMode = true
		return m, nil
	case "r":
		m.statusMsg = "Refreshing..."
		return m, doScan(m.scanner)
	case "g":
		m.showGroups = !m.showGroups
		return m, nil
	case "k":
		if len(filtered) > 0 && m.cursor < len(filtered) {
			m.view = viewConfirmKill
		}
		return m, nil
	case "enter":
		if len(filtered) > 0 && m.cursor < len(filtered) {
			m.view = viewDetail
		}
		return m, nil
	case "up":
		if m.cursor > 0 {
			m.cursor--
		}
		return m, nil
	case "down", "j":
		if m.cursor < len(filtered)-1 {
			m.cursor++
		}
		return m, nil
	case "1", "2", "3", "4", "5", "6", "7", "8":
		col := int(msg.String()[0] - '1')
		if m.sortCol.column == col {
			m.sortCol.asc = !m.sortCol.asc
		} else {
			m.sortCol.column = col
			m.sortCol.asc = true
		}
		return m, nil
	case "esc":
		if m.filter != "" {
			m.filter = ""
			m.cursor = 0
		}
		return m, nil
	}

	return m, nil
}

func (m Model) handleFilterKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.filterMode = false
		return m, nil
	case "enter":
		m.filterMode = false
		m.cursor = 0
		return m, nil
	case "backspace":
		if len(m.filter) > 0 {
			m.filter = m.filter[:len(m.filter)-1]
			m.cursor = 0
		}
		return m, nil
	default:
		if len(msg.String()) == 1 {
			m.filter += msg.String()
			m.cursor = 0
		}
		return m, nil
	}
}

func (m Model) handleDetailKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "enter", "q":
		m.view = viewTable
	}
	return m, nil
}

func (m Model) handleHelpKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "?", "esc", "q":
		m.view = viewTable
	}
	return m, nil
}

func (m Model) handleConfirmKillKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	filtered := filterPorts(m.ports, m.filter)
	sorted := sortPorts(filtered, m.sortCol)

	switch msg.String() {
	case "y", "Y":
		if m.cursor < len(sorted) {
			p := sorted[m.cursor]
			if err := process.Kill(p.PID, syscall.SIGTERM); err != nil {
				m.statusMsg = fmt.Sprintf("Failed to kill PID %d: %v", p.PID, err)
			} else {
				m.statusMsg = fmt.Sprintf("Killed PID %d (%s) on port %d", p.PID, p.ProcessName, p.Port)
			}
		}
		m.view = viewTable
		return m, doScan(m.scanner)
	case "n", "N", "esc":
		m.view = viewTable
		m.statusMsg = "Kill cancelled"
	}
	return m, nil
}

// View renders the UI.
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var sections []string

	// Header
	sections = append(sections, m.renderHeader())

	switch m.view {
	case viewHelp:
		sections = append(sections, renderHelp(m.width))
	case viewDetail:
		filtered := filterPorts(m.ports, m.filter)
		sorted := sortPorts(filtered, m.sortCol)
		if m.cursor < len(sorted) {
			sections = append(sections, renderDetail(sorted[m.cursor].PID, m.width))
		}
	case viewConfirmKill:
		sections = append(sections, renderTable(m.ports, m.cursor, m.sortCol, m.filter, m.showGroups, m.config, m.width))
		filtered := filterPorts(m.ports, m.filter)
		sorted := sortPorts(filtered, m.sortCol)
		if m.cursor < len(sorted) {
			p := sorted[m.cursor]
			confirm := confirmStyle.Render(fmt.Sprintf(
				"Kill process %q (PID %d) on port %d?\n\n  [y] Yes   [n] No",
				p.ProcessName, p.PID, p.Port,
			))
			sections = append(sections, confirm)
		}
	default:
		// Search bar
		if m.filterMode || m.filter != "" {
			search := searchStyle.Render("Filter: ") + searchInputStyle.Render(m.filter)
			if m.filterMode {
				search += "█"
			}
			sections = append(sections, search)
		}

		sections = append(sections, renderTable(m.ports, m.cursor, m.sortCol, m.filter, m.showGroups, m.config, m.width))
	}

	// Status bar
	sections = append(sections, m.renderStatusBar())

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	return content
}

func (m Model) renderHeader() string {
	filtered := filterPorts(m.ports, m.filter)
	title := titleStyle.Render("PortPilot")
	stats := headerStyle.Render(fmt.Sprintf(
		"%s │ %d ports │ %d shown",
		m.hostname, len(m.ports), len(filtered),
	))
	return lipgloss.JoinHorizontal(lipgloss.Top, title, stats)
}

func (m Model) renderStatusBar() string {
	var left string
	if m.statusMsg != "" {
		left = m.statusMsg
	} else if m.err != nil {
		left = fmt.Sprintf("Error: %v", m.err)
	} else {
		left = fmt.Sprintf("Last refresh: %s", m.lastRefresh.Format("15:04:05"))
	}

	hints := []string{
		statusKeyStyle.Render("?") + " help",
		statusKeyStyle.Render("/") + " filter",
		statusKeyStyle.Render("k") + " kill",
		statusKeyStyle.Render("q") + " quit",
	}
	right := strings.Join(hints, "  ")

	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right) - 4
	if gap < 1 {
		gap = 1
	}

	return statusBarStyle.Render(left + strings.Repeat(" ", gap) + right)
}

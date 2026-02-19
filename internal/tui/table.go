package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/AbdullahTarakji/portpilot/internal/config"
	"github.com/AbdullahTarakji/portpilot/internal/scanner"
)

type column struct {
	title string
	width int
}

var columns = []column{
	{"Port", 7},
	{"Proto", 6},
	{"PID", 8},
	{"Process", 16},
	{"User", 12},
	{"CPU%", 7},
	{"Mem%", 7},
	{"State", 8},
}

type sortOrder struct {
	column int
	asc    bool
}

// renderTable renders the port table with the current state.
func renderTable(ports []scanner.PortInfo, cursor int, sortCol sortOrder, filter string, showGroups bool, cfg *config.Config, width int) string {
	filtered := filterPorts(ports, filter)
	sorted := sortPorts(filtered, sortCol)

	// Detect conflicts (same port, different PID)
	conflicts := findConflicts(sorted)

	// Calculate dynamic process column width
	remainingWidth := width - 4 // borders/padding
	fixedWidth := 0
	for i, c := range columns {
		if i != 3 { // skip Process column
			fixedWidth += c.width + 2 // +2 for padding
		}
	}
	if showGroups {
		fixedWidth += 12 // group column
	}
	processWidth := remainingWidth - fixedWidth
	if processWidth < 10 {
		processWidth = 10
	}
	if processWidth > 30 {
		processWidth = 30
	}

	// Header
	var headerCells []string
	for i, c := range columns {
		w := c.width
		if i == 3 {
			w = processWidth
		}
		title := c.title
		if i == sortCol.column {
			if sortCol.asc {
				title += " ▲"
			} else {
				title += " ▼"
			}
		}
		headerCells = append(headerCells, tableHeaderStyle.Width(w).Render(title))
	}
	if showGroups {
		headerCells = append(headerCells, tableHeaderStyle.Width(10).Render("Group"))
	}
	header := lipgloss.JoinHorizontal(lipgloss.Top, headerCells...)

	// Rows
	var rows []string
	for i, p := range sorted {
		isSelected := i == cursor
		isConflict := conflicts[p.Port]
		isHighCPU := p.CPU > 50
		isHighMem := p.Mem > 10
		isSystem := p.PID > 0 && p.PID < 100

		var cells []string
		values := []string{
			fmt.Sprintf("%d", p.Port),
			p.Protocol,
			fmt.Sprintf("%d", p.PID),
			truncate(p.ProcessName, processWidth),
			truncate(p.User, columns[4].width),
			fmt.Sprintf("%.1f", p.CPU),
			fmt.Sprintf("%.1f", p.Mem),
			p.State,
		}

		for j, v := range values {
			w := columns[j].width
			if j == 3 {
				w = processWidth
			}
			cell := lipgloss.NewStyle().Width(w).Padding(0, 1).Render(v)
			cells = append(cells, cell)
		}

		if showGroups {
			groupName := cfg.GroupForPort(p.Port)
			groupColor := cfg.GroupColor(groupName)
			groupCell := groupLabelStyle(groupColor).Width(10).Padding(0, 1).Render(groupName)
			cells = append(cells, groupCell)
		}

		row := lipgloss.JoinHorizontal(lipgloss.Top, cells...)

		// Apply row-level styling
		switch {
		case isSelected:
			row = selectedRowStyle.Render(row)
		case isConflict:
			row = conflictStyle.Render(row)
		case isHighCPU || isHighMem:
			row = warningStyle.Render(row)
		case isSystem:
			row = dimStyle.Render(row)
		default:
			row = healthyStyle.Render(row)
		}

		rows = append(rows, row)
	}

	var parts []string
	parts = append(parts, header)
	parts = append(parts, strings.Repeat("─", width-2))
	parts = append(parts, rows...)

	if len(sorted) == 0 {
		msg := "No ports found"
		if filter != "" {
			msg = fmt.Sprintf("No ports matching %q", filter)
		}
		parts = append(parts, dimStyle.Padding(1, 2).Render(msg))
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func filterPorts(ports []scanner.PortInfo, filter string) []scanner.PortInfo {
	if filter == "" {
		return ports
	}
	lower := strings.ToLower(filter)
	var result []scanner.PortInfo
	for _, p := range ports {
		portStr := fmt.Sprintf("%d", p.Port)
		if strings.Contains(portStr, lower) ||
			strings.Contains(strings.ToLower(p.ProcessName), lower) ||
			strings.Contains(strings.ToLower(p.User), lower) ||
			strings.Contains(strings.ToLower(p.Command), lower) {
			result = append(result, p)
		}
	}
	return result
}

func sortPorts(ports []scanner.PortInfo, so sortOrder) []scanner.PortInfo {
	sorted := make([]scanner.PortInfo, len(ports))
	copy(sorted, ports)

	sort.Slice(sorted, func(i, j int) bool {
		var less bool
		switch so.column {
		case 0:
			less = sorted[i].Port < sorted[j].Port
		case 1:
			less = sorted[i].Protocol < sorted[j].Protocol
		case 2:
			less = sorted[i].PID < sorted[j].PID
		case 3:
			less = strings.ToLower(sorted[i].ProcessName) < strings.ToLower(sorted[j].ProcessName)
		case 4:
			less = strings.ToLower(sorted[i].User) < strings.ToLower(sorted[j].User)
		case 5:
			less = sorted[i].CPU < sorted[j].CPU
		case 6:
			less = sorted[i].Mem < sorted[j].Mem
		case 7:
			less = sorted[i].State < sorted[j].State
		default:
			less = sorted[i].Port < sorted[j].Port
		}
		if !so.asc {
			return !less
		}
		return less
	})

	return sorted
}

func findConflicts(ports []scanner.PortInfo) map[int]bool {
	portPIDs := make(map[int]map[int]bool)
	for _, p := range ports {
		if _, ok := portPIDs[p.Port]; !ok {
			portPIDs[p.Port] = make(map[int]bool)
		}
		portPIDs[p.Port][p.PID] = true
	}

	conflicts := make(map[int]bool)
	for port, pids := range portPIDs {
		if len(pids) > 1 {
			conflicts[port] = true
		}
	}
	return conflicts
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

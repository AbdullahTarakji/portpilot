package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/AbdullahTarakji/portpilot/internal/process"
)

func renderDetail(pid int, width int) string {
	details, err := process.GetDetails(pid)
	if err != nil {
		return detailBorderStyle.Width(width - 4).Render(
			fmt.Sprintf("Error getting details for PID %d: %v", pid, err),
		)
	}

	rows := []struct {
		key   string
		value string
	}{
		{"PID", fmt.Sprintf("%d", details.PID)},
		{"Parent PID", fmt.Sprintf("%d", details.ParentPID)},
		{"Name", details.Name},
		{"User", details.User},
		{"CPU", fmt.Sprintf("%.1f%%", details.CPU)},
		{"Memory", fmt.Sprintf("%.1f%%", details.Mem)},
		{"Started", details.StartTime.Format("2006-01-02 15:04:05")},
		{"Command", details.Command},
	}

	var lines []string
	lines = append(lines, titleStyle.Render("Process Details"))
	lines = append(lines, "")

	for _, r := range rows {
		line := lipgloss.JoinHorizontal(lipgloss.Top,
			detailKeyStyle.Render(r.key+":"),
			detailValueStyle.Render(r.value),
		)
		lines = append(lines, line)
	}

	lines = append(lines, "")
	lines = append(lines, dimStyle.Render("Press Esc or Enter to close"))

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return detailBorderStyle.Width(width - 4).Render(content)
}

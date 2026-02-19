package tui

import "github.com/charmbracelet/lipgloss"

type helpEntry struct {
	key  string
	desc string
}

var helpEntries = []helpEntry{
	{"1-8", "Sort by column (toggle asc/desc)"},
	{"/", "Search / filter by port or process"},
	{"Esc", "Clear search / close panel"},
	{"Enter", "View process details"},
	{"k", "Kill selected process"},
	{"r", "Manual refresh"},
	{"g", "Toggle group view"},
	{"?", "Toggle this help"},
	{"q", "Quit"},
	{"Up/Down", "Navigate rows"},
	{"j/k nav", "Vim-style navigation (when not killing)"},
}

func renderHelp(width int) string {
	var rows []string
	for _, e := range helpEntries {
		row := lipgloss.JoinHorizontal(lipgloss.Top,
			helpKeyStyle.Render(e.key),
			helpDescStyle.Render(e.desc),
		)
		rows = append(rows, row)
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("Keybindings"),
		"",
	)

	for _, r := range rows {
		content = lipgloss.JoinVertical(lipgloss.Left, content, r)
	}

	content = lipgloss.JoinVertical(lipgloss.Left, content, "", dimStyle.Render("Press ? or Esc to close"))

	styled := helpStyle.Width(width - 4).Render(content)
	return styled
}

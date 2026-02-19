package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	colorRed     = lipgloss.Color("#FF5555")
	colorGreen   = lipgloss.Color("#50FA7B")
	colorYellow  = lipgloss.Color("#F1FA8C")
	colorBlue    = lipgloss.Color("#6272A4")
	colorMagenta = lipgloss.Color("#FF79C6")
	colorCyan    = lipgloss.Color("#8BE9FD")
	colorWhite   = lipgloss.Color("#F8F8F2")
	colorDim     = lipgloss.Color("#6272A4")
	colorBgAlt = lipgloss.Color("#44475A")

	// Header
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorCyan).
			Padding(0, 1)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorMagenta).
			Padding(0, 1)

	// Status bar
	statusBarStyle = lipgloss.NewStyle().
			Foreground(colorDim).
			Padding(0, 1)

	statusKeyStyle = lipgloss.NewStyle().
			Foreground(colorCyan).
			Bold(true)

	// Table
	tableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorWhite).
				Background(colorBgAlt).
				Padding(0, 1)

	selectedRowStyle = lipgloss.NewStyle().
				Background(colorBgAlt).
				Foreground(colorWhite).
				Bold(true)

	// Row coloring
	conflictStyle = lipgloss.NewStyle().
			Background(colorRed).
			Foreground(colorWhite)

	warningStyle = lipgloss.NewStyle().
			Foreground(colorYellow)

	healthyStyle = lipgloss.NewStyle().
			Foreground(colorGreen)

	dimStyle = lipgloss.NewStyle().
			Foreground(colorDim)

	// Search
	searchStyle = lipgloss.NewStyle().
			Foreground(colorCyan).
			Bold(true).
			Padding(0, 1)

	searchInputStyle = lipgloss.NewStyle().
				Foreground(colorWhite)

	// Detail view
	detailKeyStyle = lipgloss.NewStyle().
			Foreground(colorCyan).
			Bold(true).
			Width(14).
			Align(lipgloss.Right)

	detailValueStyle = lipgloss.NewStyle().
				Foreground(colorWhite).
				PaddingLeft(2)

	detailBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorBlue).
				Padding(1, 2)

	// Help overlay
	helpStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorMagenta).
			Padding(1, 2)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(colorCyan).
			Bold(true).
			Width(12)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(colorWhite)

	// Confirmation dialog
	confirmStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorRed).
			Padding(1, 2).
			Foreground(colorWhite)

	// Group labels
	groupColors = map[string]lipgloss.Color{
		"red":     colorRed,
		"green":   colorGreen,
		"yellow":  colorYellow,
		"blue":    colorBlue,
		"magenta": colorMagenta,
		"cyan":    colorCyan,
	}
)

func groupLabelStyle(color string) lipgloss.Style {
	c, ok := groupColors[color]
	if !ok {
		c = colorDim
	}
	return lipgloss.NewStyle().Foreground(c).Bold(true)
}

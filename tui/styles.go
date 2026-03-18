package tui

import "github.com/charmbracelet/lipgloss"

var (
	barStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 2)

	currentStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86"))

	dimmedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
	scanStatusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))
	scanWarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))
	selectedTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("86")).
				BorderLeft(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("86")).
				PaddingLeft(1)

	selectedDescStyle = selectedTitleStyle.
				Foreground(lipgloss.Color("86"))

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")).
			Padding(0, 1)
	emptyStyle = lipgloss.NewStyle().
			Padding(1, 3).
			Foreground(lipgloss.Color("240"))
)

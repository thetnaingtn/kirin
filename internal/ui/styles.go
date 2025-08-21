package ui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	InputStyle = lipgloss.NewStyle().
			Padding(1, 2)

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))
)

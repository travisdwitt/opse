package ui

import "github.com/charmbracelet/lipgloss"

var (
	SidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 1)

	SidebarFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("252")).
				Padding(1, 1)

	CategoryStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("3"))

	ItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	ItemSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("3")).
				Bold(true)

	// Prominent selection style for saved rolls modal
	SrSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("0")).
				Background(lipgloss.Color("3")).
				Bold(true).
				Padding(0, 1)

	ShortcutStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	LogStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	LogStyleFocused = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("252")).
			Padding(0, 1)

	ResultBlockStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62")).
				Padding(0, 1).
				MarginTop(1).MarginBottom(1)

	ResultLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("62")).
				Bold(true)

	SuitRedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	SuitWhiteStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	DimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	InputStyleFocused = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("252")).
				Padding(0, 1)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("3")).
			Padding(0, 1)

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	PortraitBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62"))
)

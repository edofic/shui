package main

import "github.com/charmbracelet/lipgloss"

var (
	// Border styles
	BorderColor        = lipgloss.Color("240")
	FocusedBorderColor = lipgloss.Color("62")

	// Status bar colors
	SuccessColor = lipgloss.Color("42")
	ErrorColor   = lipgloss.Color("196")
	InfoColor    = lipgloss.Color("39")

	// Editor styles
	EditorStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(0, 1)

	EditorFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(FocusedBorderColor).
				Padding(0, 1)

	// Output styles
	OutputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(0, 1)

	// Status bar styles
	StatusBarStyle = lipgloss.NewStyle().
			Padding(0, 1)

	StatusSuccessStyle = StatusBarStyle.
				Foreground(SuccessColor)

	StatusErrorStyle = StatusBarStyle.
				Foreground(ErrorColor)

	StatusInfoStyle = StatusBarStyle.
			Foreground(InfoColor)

	// Title styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	// Help text style
	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

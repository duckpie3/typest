package ui

import "charm.land/lipgloss/v2"

// Colors

var (
	fg     = lipgloss.Color("7")
	bg     = lipgloss.Color("0")
	muted  = lipgloss.Color("8")
	err    = lipgloss.Color("1")
	accent = lipgloss.Color("4")
)

// Styles
var (
	ErrorStyle = lipgloss.NewStyle().
			Foreground(err)

	UntypedStyle = lipgloss.NewStyle().
			Foreground(muted)

	TypedStyle = lipgloss.NewStyle().
			Foreground(fg)

	CursorStyle = lipgloss.NewStyle().
			Foreground(bg).
			Background(fg)

	TestStyle = lipgloss.NewStyle().
			Width(80).
			Align(lipgloss.Left)

	InputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(fg).
			PaddingLeft(1)

	ControlGuideStyle = lipgloss.NewStyle().
				Foreground(muted).
				Align(lipgloss.Left)

	GraphStyle = lipgloss.NewStyle().
			Foreground(accent)
)

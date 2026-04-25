package ui

import "charm.land/lipgloss/v2"

var ErrorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("1"))

var UntypedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("8"))

var TypedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("7"))

var CursorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("0")).
	Background(lipgloss.Color("7"))

var TestStyle = lipgloss.NewStyle().
	Width(80).
	Align(lipgloss.Left)

var InputStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	PaddingLeft(1)

var GraphStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("4"))

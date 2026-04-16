package main

import "charm.land/lipgloss/v2"

var errorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("1"))

var untypedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("8"))

var typedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("7"))

var cursorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("0")).
	Background(lipgloss.Color("7"))

var testStyle = lipgloss.NewStyle().
	Width(64).
	Align(lipgloss.Left)

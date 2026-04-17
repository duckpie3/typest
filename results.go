package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type resultsModel struct {
	wpm      int
	nextTest bool
}

func NewResultsModel(wpm int) resultsModel {
	return resultsModel{wpm: wpm, nextTest: false}
}

func (m resultsModel) Init() tea.Cmd {
	return nil
}

func (m resultsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "space":
			m.nextTest = true
			return m, nil
		}
	}
	return m, nil
}

func (m resultsModel) View() tea.View {
	content := fmt.Sprintf("WPM: %d\n\nPress spacebar for next test", m.wpm)
	content = lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, content)

	return tea.View{Content: content, AltScreen: true}
}

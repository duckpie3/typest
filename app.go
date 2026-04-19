package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
)

var termWidth int
var termHeight int

type appModel struct {
	currentModel tea.Model
}

func NewAppModel() appModel {
	return appModel{currentModel: NewTypingTestModel()}
}

func (m appModel) Init() tea.Cmd {
	return m.currentModel.Init()
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch _msg := msg.(type) {
	case tea.WindowSizeMsg:
		termWidth = _msg.Width
		termHeight = _msg.Height
	case tea.KeyPressMsg:
		switch _msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.currentModel, cmd = m.currentModel.Update(msg)
	switch currentModel := m.currentModel.(type) {
	case typingTestModel:
		if currentModel.finished {
			m.currentModel = NewResultsModel(currentModel.wpm)
		}
	case resultsModel:
		if currentModel.nextTest {
			m.currentModel = NewTypingTestModel()
		}
	}

	return m, cmd
}

func (m appModel) View() tea.View {
	return m.currentModel.View()
}

func main() {
	m := NewAppModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

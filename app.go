package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
)

var termWidth int
var termHeight int

type appModel struct {
	screen tea.Model
}

func NewAppModel() appModel {
	return appModel{screen: NewTypingTestModel()}
}

func (m appModel) Init() tea.Cmd {
	return m.screen.Init()
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch _msg := msg.(type) {
	case tea.WindowSizeMsg:
		termWidth = _msg.Width
		termHeight = _msg.Height
	}

	var cmd tea.Cmd
	m.screen, cmd = m.screen.Update(msg)
	switch screen := m.screen.(type) {
	case typingTestModel:
		if screen.finished {
			m.screen = NewResultsModel(screen.wpm)
		}
	case resultsModel:
		if screen.nextTest {
			m.screen = NewTypingTestModel()
		}
	}
	return m, cmd
}

func (m appModel) View() tea.View {
	return m.screen.View()
}

func main() {
	m := NewAppModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

package app

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/duckpie3/typest/internal/loader"
	"github.com/duckpie3/typest/internal/results"
	"github.com/duckpie3/typest/internal/typing"
)

type Model struct {
	currentModel tea.Model
	quotesData   *loader.QuotesData
	wordsData    *loader.WordsData
	width        int
	height       int
}

func New() Model {
	quotes, err := loader.LoadQuotes("assets/quotes.json")
	if err != nil {
		log.Fatal(err)
	}
	words, err := loader.LoadWords("assets/words.json")
	if err != nil {
		log.Fatal(err)
	}

	return Model{currentModel: typing.New(), quotesData: quotes, wordsData: words}
}

func (m Model) Init() tea.Cmd {
	return m.currentModel.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch _msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = _msg.Width
		m.height = _msg.Height
	case tea.KeyPressMsg:
		switch _msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.currentModel, cmd = m.currentModel.Update(msg)
	switch currentModel := m.currentModel.(type) {
	case typing.Model:
		if currentModel.Done {
			m.currentModel = results.New(currentModel.Stats)
			seeded, _ := m.currentModel.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
			m.currentModel = seeded
		}
	case results.Model:
		if currentModel.NextTest {
			m.currentModel = typing.New()
			seeded, _ := m.currentModel.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
			m.currentModel = seeded
		}
	}

	return m, cmd
}

func (m Model) View() tea.View {
	return m.currentModel.View()
}

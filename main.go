package main

import (
	"log"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	testWords      []string
	testWordsView  []string
	testPosition   int
	inputView      string
	cursorPosition int
	charsStack     []string
	testView       string
	inputModel     textinput.Model
}

func NewModel() model {
	ti := textinput.New()
	ti.Focus()

	text := "I am putting myself to the fullest possible use, which is all I think that any conscious entity can ever hope to do."
	words := strings.Split(text, " ")
	wordsView := make([]string, len(words))
	for i := range words {
		words[i] += " "
		wordsView[i] = untypedStyle.Render(words[i])
	}

	return model{
		testWords:      words,
		testWordsView:  wordsView,
		testPosition:   0,
		inputView:      "",
		cursorPosition: 0,
		charsStack:     []string{},
		testView:       "",
		inputModel:     ti,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "left", "right":
			return m, nil
		case "space":
			if m.testWords[m.testPosition] == m.inputModel.Value()+" " {
				if m.testPosition+1 >= len(m.testWords) {
					return m, tea.Quit
				}
				m.testWordsView[m.testPosition] = typedStyle.Render(m.testWords[m.testPosition])
				m.testPosition += 1
				cursor := cursorStyle.Render(string(m.testWords[m.testPosition][0]))
				m.testWordsView[m.testPosition] = cursor + untypedStyle.Render(m.testWords[m.testPosition][1:])
				m.inputView = ""
				m.inputModel.SetValue("")
				m.cursorPosition = 0
				m.charsStack = []string{}
				return m, nil
			}

		}
	}
	var cmd tea.Cmd
	m.inputModel, cmd = m.inputModel.Update(msg)
	currentWord := m.testWords[m.testPosition]
	var lastTypedChar byte

	if m.cursorPosition < m.inputModel.Position() { // User enters a character
		lastTypedChar = m.inputModel.Value()[m.cursorPosition]
	} else if m.cursorPosition > m.inputModel.Position() { // User deletes a character
		top := m.charsStack[len(m.charsStack)-1]
		m.charsStack = m.charsStack[:len(m.charsStack)-1]
		m.inputView = m.inputView[:len(m.inputView)-len(top)]
	}

	var cursor string
	if m.inputModel.Position() >= len(currentWord) {
		cursor = cursorStyle.Render(" ")
		if m.cursorPosition < m.inputModel.Position() { // User entered a character
			lastTypedCharStyled := errorStyle.Render(string(lastTypedChar))
			m.charsStack = append(m.charsStack, lastTypedCharStyled)
			m.inputView += lastTypedCharStyled
		}
		m.testWordsView[m.testPosition] = m.inputView + cursor
	} else {
		if m.cursorPosition < m.inputModel.Position() {
			var lastTypedCharStyled string
			if lastTypedChar == currentWord[m.cursorPosition] {
				lastTypedCharStyled = typedStyle.Render(string(lastTypedChar))
			} else {
				lastTypedCharStyled = errorStyle.Render(string(lastTypedChar))
			}
			m.charsStack = append(m.charsStack, lastTypedCharStyled)
			m.inputView += lastTypedCharStyled
		}
		cursor = cursorStyle.Render(string(currentWord[m.inputModel.Position()]))
		m.testWordsView[m.testPosition] = m.inputView + cursor + untypedStyle.Render(currentWord[m.inputModel.Position()+1:])
	}

	m.cursorPosition = m.inputModel.Position()
	return m, cmd
}

func (m model) View() tea.View {

	m.testView = ""
	for _, w := range m.testWordsView {
		m.testView += w
	}
	var s string
	s = m.testView
	s = testStyle.Render(s)
	s += "\n"
	s += m.inputModel.View()
	s += "\n"
	return tea.View{Content: s}
}

func main() {
	m := NewModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

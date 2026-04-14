package main

import (
	"log"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	text             string
	words            []string
	wordsView        []string
	wordPosition     int
	currentInputView string
	cursorPosition   int
	charsStack       []string
	textView         string
	inputModel       textinput.Model
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
		text:             text,
		words:            words,
		wordsView:        wordsView,
		wordPosition:     0,
		currentInputView: "",
		cursorPosition:   0,
		charsStack:       []string{},
		textView:         "",
		inputModel:       ti,
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
			if m.words[m.wordPosition] == m.inputModel.Value()+" " {
				if m.wordPosition+1 >= len(m.words) {
					return m, tea.Quit
				}
				m.wordsView[m.wordPosition] = typedStyle.Render(m.words[m.wordPosition])
				m.wordPosition += 1
				cursor := cursorStyle.Render(string(m.words[m.wordPosition][0]))
				m.wordsView[m.wordPosition] = cursor + untypedStyle.Render(m.words[m.wordPosition][1:])
				m.currentInputView = ""
				m.inputModel.SetValue("")
				m.cursorPosition = 0
				m.charsStack = []string{}
				return m, nil
			}

		}
	}
	var cmd tea.Cmd
	m.inputModel, cmd = m.inputModel.Update(msg)
	currentWord := m.words[m.wordPosition]
	var lastTypedChar byte

	if m.cursorPosition < m.inputModel.Position() { // User enters a character
		lastTypedChar = m.inputModel.Value()[m.cursorPosition]
	} else if m.cursorPosition > m.inputModel.Position() { // User deletes a character
		top := m.charsStack[len(m.charsStack)-1]
		m.charsStack = m.charsStack[:len(m.charsStack)-1]
		m.currentInputView = m.currentInputView[:len(m.currentInputView)-len(top)]
	}

	var cursor string
	if m.inputModel.Position() >= len(currentWord) {
		cursor = cursorStyle.Render(" ")
		if m.cursorPosition < m.inputModel.Position() { // User entered a character
			lastTypedCharStyled := errorStyle.Render(string(lastTypedChar))
			m.charsStack = append(m.charsStack, lastTypedCharStyled)
			m.currentInputView += lastTypedCharStyled
		}
		m.wordsView[m.wordPosition] = m.currentInputView + cursor
	} else {
		if m.cursorPosition < m.inputModel.Position() {
			var lastTypedCharStyled string
			if lastTypedChar == currentWord[m.cursorPosition] {
				lastTypedCharStyled = typedStyle.Render(string(lastTypedChar))
			} else {
				lastTypedCharStyled = errorStyle.Render(string(lastTypedChar))
			}
			m.charsStack = append(m.charsStack, lastTypedCharStyled)
			m.currentInputView += lastTypedCharStyled
		}
		cursor = cursorStyle.Render(string(currentWord[m.inputModel.Position()]))
		m.wordsView[m.wordPosition] = m.currentInputView + cursor + untypedStyle.Render(currentWord[m.inputModel.Position()+1:])
	}

	m.cursorPosition = m.inputModel.Position()
	return m, cmd
}

func (m model) View() tea.View {

	m.textView = ""
	for _, w := range m.wordsView {
		m.textView += w
	}
	var s string
	s = m.textView
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

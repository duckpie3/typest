package main

import (
	"log"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type wpmDataPoint struct {
	time float64
	wpm  int
}

type testStats struct {
	characters  int
	wpm         int
	wpmData     []wpmDataPoint
	greatestwpm int
	startTime   time.Time
	elapsedTime float64
}

type typingTestModel struct {
	testWords      []string
	testWordsView  []string
	testPosition   int
	inputView      string
	cursorPosition int
	charsStack     []string
	testView       string
	inputModel     textinput.Model
	started        bool
	finished       bool
	typedChars     int
	stats          testStats
}

var nextSecond float64

func (m *typingTestModel) nextWord() {
	m.testWordsView[m.testPosition] = typedStyle.Render(m.testWords[m.testPosition])
	m.testPosition += 1
	cursor := cursorStyle.Render(string(m.testWords[m.testPosition][0]))
	m.testWordsView[m.testPosition] = cursor + untypedStyle.Render(m.testWords[m.testPosition][1:])
	m.inputView = ""
	m.inputModel.SetValue("")
	m.cursorPosition = 0
	m.charsStack = []string{}
}

func NewTypingTestModel() typingTestModel {
	ti := textinput.New()
	ti.Focus()
	nextSecond = 1
	data, err := loadQuotes("quotes.json")
	if err != nil {
		panic(err)
	}
	quote := data.RandomQuote()
	words := strings.Split(quote.Text, " ")
	wordsView := make([]string, len(words))
	for i := range words {
		words[i] += " "
		wordsView[i] = untypedStyle.Render(words[i])
	}

	return typingTestModel{
		testWords:      words,
		testWordsView:  wordsView,
		testPosition:   0,
		inputView:      "",
		cursorPosition: 0,
		charsStack:     []string{},
		testView:       "",
		inputModel:     ti,
		started:        false,
		finished:       false,
		typedChars:     0,
		stats:          testStats{characters: quote.Length},
	}
}

func (m typingTestModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m typingTestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.PasteMsg:
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "left", "right", "ctrl+v":
			return m, nil
		case "space":
			if m.testPosition < len(m.testWords)-1 && m.testWords[m.testPosition] == m.inputModel.Value()+" " {
				if m.testPosition+1 >= len(m.testWords) {
					return m, tea.Quit
				}
				m.nextWord()
				return m, nil
			} else {
				return m, nil
			}
		}
	}
	var cmd tea.Cmd
	m.inputModel, cmd = m.inputModel.Update(msg)
	currentWord := m.testWords[m.testPosition]
	var lastTypedChar byte

	if m.cursorPosition < m.inputModel.Position() { // User enters a character
		m.typedChars++
		if !m.started {
			m.started = true
			m.stats.startTime = time.Now()
		}
		lastTypedChar = m.inputModel.Value()[m.cursorPosition]
	} else if m.cursorPosition > m.inputModel.Position() { // User deletes a character
		m.typedChars--
		if len(m.charsStack) < 1 {
			log.Fatal("Tried to delete a character that didn't exist.")
		}
		top := m.charsStack[len(m.charsStack)-1]
		m.charsStack = m.charsStack[:len(m.charsStack)-1]
		m.inputView = m.inputView[:len(m.inputView)-len(top)]
	}

	var cursor string
	if m.inputModel.Position() >= len(currentWord) { // Out of bounds
		cursor = cursorStyle.Render(" ")
		if m.cursorPosition < m.inputModel.Position() { // User entered a character
			lastTypedCharStyled := errorStyle.Render(string(lastTypedChar))
			m.charsStack = append(m.charsStack, lastTypedCharStyled)
			m.inputView += lastTypedCharStyled
		}
		m.testWordsView[m.testPosition] = m.inputView + cursor
	} else { // Not out of bounds
		if m.cursorPosition < m.inputModel.Position() {
			var lastTypedCharStyled string
			if lastTypedChar == currentWord[m.cursorPosition] {
				lastTypedCharStyled = typedStyle.Render(string(lastTypedChar))
			} else {
				lastTypedCharStyled = errorStyle.Render(string(currentWord[m.cursorPosition]))
			}
			m.charsStack = append(m.charsStack, lastTypedCharStyled)
			m.inputView += lastTypedCharStyled
		}
		cursor = cursorStyle.Render(string(currentWord[m.inputModel.Position()]))
		m.testWordsView[m.testPosition] = m.inputView + cursor + untypedStyle.Render(currentWord[m.inputModel.Position()+1:])
	}

	if m.started == true {
		m.stats.elapsedTime = time.Since(m.stats.startTime).Seconds()
		if m.stats.elapsedTime > 1 {
			m.stats.wpm = int(float64(m.typedChars) / 5 / (m.stats.elapsedTime / 60))
			if m.stats.wpm > m.stats.greatestwpm {
				m.stats.greatestwpm = m.stats.wpm
			}
		}
		// Add performance data point
		if m.stats.elapsedTime > nextSecond {
			nextSecond = m.stats.elapsedTime + 1
			m.stats.wpmData = append(m.stats.wpmData, wpmDataPoint{time: m.stats.elapsedTime, wpm: m.stats.wpm})
		}
	}

	// When on the last word, check if it's correct so there is no need Pleaseto enter space
	if m.testPosition == len(m.testWords)-1 && len(m.inputModel.Value()) == len(currentWord)-1 {
		if m.testWords[m.testPosition] == m.inputModel.Value()+" " {
			m.stats.wpmData = append(m.stats.wpmData, wpmDataPoint{time: m.stats.elapsedTime, wpm: m.stats.wpm})
			m.finished = true
			return m, nil
		}
	}
	m.cursorPosition = m.inputModel.Position()
	return m, cmd
}

func (m typingTestModel) View() tea.View {

	m.testView = ""
	for _, w := range m.testWordsView {
		m.testView += w
	}

	content := testStyle.Render(m.testView) + "\n\n" + m.inputModel.View()
	s := lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, content)

	return tea.View{Content: s, AltScreen: true}
}

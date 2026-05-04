package typing

import (
	"log"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/duckpie3/typest/internal/loader"
	"github.com/duckpie3/typest/internal/ui"
)

type Model struct {
	testWords      []string
	testWordsView  []string
	testPosition   int
	inputView      string
	cursorPosition int
	charsStack     []string
	testView       string
	inputModel     textinput.Model
	started        bool
	Done           bool
	typedChars     int
	Stats          TestStats
	width          int
	height         int
}

var nextSecond float64

func (m *Model) clearCurrentWord() {
	if m.testPosition >= len(m.testWords) {
		return
	}

	deletedChars := len(m.inputModel.Value())
	if deletedChars > m.typedChars {
		m.typedChars = 0
	} else {
		m.typedChars -= deletedChars
	}

	m.inputModel.SetValue("")
	m.inputView = ""
	m.cursorPosition = 0
	m.charsStack = []string{}

	currentWord := m.testWords[m.testPosition]
	if len(currentWord) == 0 {
		m.testWordsView[m.testPosition] = ui.CursorStyle.Render(" ")
		return
	}

	m.testWordsView[m.testPosition] = ui.CursorStyle.Render(string(currentWord[0])) + ui.TypedStyle.Render(currentWord[1:])
}

func (m *Model) nextWord() {
	m.testWordsView[m.testPosition] = ui.TypedStyle.Render(m.testWords[m.testPosition])
	m.testPosition += 1
	cursor := ui.CursorStyle.Render(string(m.testWords[m.testPosition][0]))
	m.testWordsView[m.testPosition] = cursor + ui.UntypedStyle.Render(m.testWords[m.testPosition][1:])
	m.inputView = ""
	m.inputModel.SetValue("")
	m.cursorPosition = 0
	m.charsStack = []string{}
}

func New() Model {
	ti := textinput.New()
	ti.Focus()
	ti.SetWidth(32)
	ti.Prompt = ""
	ti.Placeholder = "Type the above word here"
	nextSecond = 1
	data, err := loader.LoadQuotes("assets/quotes.json")
	if err != nil {
		panic(err)
	}
	quote := data.RandomQuote()
	words := strings.Split(quote.Text, " ")
	wordsView := make([]string, len(words))
	for i := range words {
		words[i] += " "
		wordsView[i] = ui.UntypedStyle.Render(words[i])
	}

	return Model{
		testWords:      words,
		testWordsView:  wordsView,
		testPosition:   0,
		inputView:      "",
		cursorPosition: 0,
		charsStack:     []string{},
		testView:       "",
		inputModel:     ti,
		started:        false,
		Done:           false,
		typedChars:     0,
		Stats:          TestStats{Characters: quote.Length},
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.PasteMsg:
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "left", "right", "ctrl+v":
			return m, nil
		case "ctrl+backspace", "alt+backspace", "ctrl+w":
			m.clearCurrentWord()
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
			m.Stats.startTime = time.Now()
			m.inputModel.Placeholder = ""
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
		cursor = ui.CursorStyle.Render(" ")
		if m.cursorPosition < m.inputModel.Position() { // User entered a character
			lastTypedCharStyled := ui.ErrorStyle.Render(string(lastTypedChar))
			m.charsStack = append(m.charsStack, lastTypedCharStyled)
			m.inputView += lastTypedCharStyled
		}
		m.testWordsView[m.testPosition] = m.inputView + cursor
	} else { // Not out of bounds
		if m.cursorPosition < m.inputModel.Position() {
			var lastTypedCharStyled string
			if lastTypedChar == currentWord[m.cursorPosition] {
				lastTypedCharStyled = ui.TypedStyle.Render(string(lastTypedChar))
			} else {
				lastTypedCharStyled = ui.ErrorStyle.Render(string(currentWord[m.cursorPosition]))
			}
			m.charsStack = append(m.charsStack, lastTypedCharStyled)
			m.inputView += lastTypedCharStyled
		}
		cursor = ui.CursorStyle.Render(string(currentWord[m.inputModel.Position()]))
		m.testWordsView[m.testPosition] = m.inputView + cursor + ui.UntypedStyle.Render(currentWord[m.inputModel.Position()+1:])
	}

	if m.started == true {
		m.Stats.ElapsedTime = time.Since(m.Stats.startTime).Seconds()
		if m.Stats.ElapsedTime > 1 {
			m.Stats.Wpm = int(float64(m.typedChars) / 5 / (m.Stats.ElapsedTime / 60))
			if m.Stats.Wpm > m.Stats.Greatestwpm {
				m.Stats.Greatestwpm = m.Stats.Wpm
			}
		}
		// Add performance data point
		if m.Stats.ElapsedTime > nextSecond {
			nextSecond = m.Stats.ElapsedTime + 1
			m.Stats.WpmData = append(m.Stats.WpmData, wpmDataPoint{Time: m.Stats.ElapsedTime, Wpm: m.Stats.Wpm})
		}
	}

	// When on the last word, check if it's correct so there is no need Pleaseto enter space
	if m.testPosition == len(m.testWords)-1 && len(m.inputModel.Value()) == len(currentWord)-1 {
		if m.testWords[m.testPosition] == m.inputModel.Value()+" " {
			m.Stats.WpmData = append(m.Stats.WpmData, wpmDataPoint{Time: m.Stats.ElapsedTime, Wpm: m.Stats.Wpm})
			m.Done = true
			return m, nil
		}
	}
	m.cursorPosition = m.inputModel.Position()
	return m, cmd
}

func (m Model) View() tea.View {

	m.testView = ""
	for _, w := range m.testWordsView {
		m.testView += w
	}

	content := ui.TestStyle.Render(m.testView) + "\n\n\n" + ui.InputStyle.Render(m.inputModel.View())
	s := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)

	return tea.View{Content: s, AltScreen: true}
}

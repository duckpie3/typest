package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/NimbleMarkets/ntcharts/v2/canvas"
	lc "github.com/NimbleMarkets/ntcharts/v2/linechart"
)

type resultsModel struct {
	stats    testStats
	graph    lc.Model
	nextTest bool
}

func (m *resultsModel) buildGraph() {
	for i := range len(m.stats.wpmData) - 1 {
		point1 := canvas.Float64Point{X: m.stats.wpmData[i].time, Y: float64(m.stats.wpmData[i].wpm)}
		point2 := canvas.Float64Point{X: m.stats.wpmData[i+1].time, Y: float64(m.stats.wpmData[i+1].wpm)}
		m.graph.DrawBrailleLine(point1, point2)
	}
	m.graph.DrawXYAxisAndLabel()
}

func NewResultsModel(stats testStats) resultsModel {
	width := 130
	height := 10
	graph := lc.New(width, height, 1.0, stats.elapsedTime, 0.0, float64(stats.greatestwpm))
	return resultsModel{stats: stats, graph: graph, nextTest: false}
}

func (m resultsModel) Init() tea.Cmd {
	return nil
}

func (m resultsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "space":
			m.nextTest = true
			return m, nil
		}
	}
	return m, nil
}

func (m resultsModel) View() tea.View {
	m.buildGraph()
	content := fmt.Sprintf("WPM: %d\n\n%v\n\ntime: %.2fs\tcharacters: %d\n\nPress spacebar for next test", m.stats.wpm, m.graph.View(), m.stats.elapsedTime, m.stats.characters)
	content = lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, content)

	return tea.View{Content: content, AltScreen: true}
}

package results

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/NimbleMarkets/ntcharts/v2/canvas"
	lc "github.com/NimbleMarkets/ntcharts/v2/linechart"
	"github.com/duckpie3/typest/internal/typing"
	"github.com/duckpie3/typest/internal/ui"
)

type Model struct {
	stats    typing.TestStats
	graph    lc.Model
	NextTest bool
	width    int
	height   int
}

func (m *Model) buildGraph() {
	m.graph.Style = ui.GraphStyle
	m.graph.DrawXYAxisAndLabel()
	for i := range len(m.stats.WpmData) - 1 {
		point1 := canvas.Float64Point{X: m.stats.WpmData[i].Time, Y: float64(m.stats.WpmData[i].Wpm)}
		point2 := canvas.Float64Point{X: m.stats.WpmData[i+1].Time, Y: float64(m.stats.WpmData[i+1].Wpm)}
		m.graph.DrawBrailleLine(point1, point2)
	}
}

func New(stats typing.TestStats) Model {
	width := 100
	height := 10
	graph := lc.New(width, height, 1.0, stats.ElapsedTime, 0.0, float64(stats.Greatestwpm))
	return Model{stats: stats, graph: graph, NextTest: false}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "space":
			m.NextTest = true
			return m, nil
		}
	}
	return m, nil
}

func (m Model) View() tea.View {
	m.buildGraph()
	content := fmt.Sprintf("WPM: %d\n\n%v\n\ntime: %.2fs\tcharacters: %d\n\nPress spacebar for next test", m.stats.Wpm, m.graph.View(), m.stats.ElapsedTime, m.stats.Characters)
	content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)

	return tea.View{Content: content, AltScreen: true}
}

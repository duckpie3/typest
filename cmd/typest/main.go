package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/duckpie3/typest/internal/app"
)

func main() {
	m := app.New()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

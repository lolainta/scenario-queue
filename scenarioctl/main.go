package main

import (
	"log"

	"github.com/joho/godotenv"

	tea "github.com/charmbracelet/bubbletea"

	"scenarioctl/app"
	"scenarioctl/db"
	"scenarioctl/pages"
	"scenarioctl/repo"
)

func main() {
	godotenv.Load()

	pool := db.Connect()
	defer pool.Close()

	r := repo.New(pool)

	a := app.New([]app.NavItem{
		{Name: "Tasks", New: func() app.Page { return pages.NewTaskPage(r) }},
		{Name: "Plans", New: func() app.Page { return pages.NewPlanPage(r) }},
		{Name: "AVs", New: func() app.Page { return pages.NewAVPage(r) }},
		{Name: "Simulators", New: func() app.Page { return pages.NewSimulatorPage(r) }},
		{Name: "Scenarios", New: func() app.Page { return pages.NewScenarios(r) }},
		{Name: "Maps", New: func() app.Page { return pages.NewMapPage(r) }},
		{Name: "Samplers", New: func() app.Page { return pages.NewSamplerPage(r) }},
	})

	p := tea.NewProgram(a, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

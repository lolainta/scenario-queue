package pages

import (
	"context"
	"fmt"

	"scenarioctl/app"
	"scenarioctl/repo"

	"github.com/charmbracelet/bubbles/table"
)

func NewScenarios(r *repo.Repo) app.Page {
	var tp *TablePage
	tp = NewTablePage(
		"Scenarios",
		[]table.Column{
			{Title: "ID", Width: 6},
			{Title: "Title", Width: 20},
			{Title: "Scenario Path", Width: 30},
			{Title: "Goal Config", Width: 30},
		},
		func(ctx context.Context) ([]table.Row, error) {
			scenarios, err := r.ListScenarios(ctx)
			if err != nil {
				return nil, err
			}

			var out []table.Row
			for _, scenario := range scenarios {
				title := ""
				if scenario.Title != nil && *scenario.Title != "" {
					title = *scenario.Title
				}
				scenarioPath := scenario.ScenarioPath
				if scenario.ScenarioPath != "" {
					scenarioPath = scenario.ScenarioPath
				}
				goalConfig := ""
				if scenario.GoalConfig != nil && *scenario.GoalConfig != "" {
					goalConfig = *scenario.GoalConfig
				}
				out = append(out, table.Row{
					fmt.Sprint(scenario.ID),
					title,
					scenarioPath,
					goalConfig,
				})
			}
			return out, nil
		},
	).(*TablePage)

	tp.WithCRUD(&CRUDCallbacks{
		OnCreate: func() error {
			tp.StartForm(4, []string{"Title", "Scenario Path", "Goal Config"}, -1, func(values []string) error {
				var title *string
				if values[0] != "" {
					title = &values[0]
				}
				var goalConfig *string
				if values[2] != "" {
					goalConfig = &values[2]
				}
				_, err := r.CreateScenario(context.Background(), title, values[1], goalConfig)
				return err
			})
			return nil
		},
		OnUpdate: func(rowIndex int) error {
			if rowIndex >= len(tp.currentRows) {
				return nil
			}
			row := tp.currentRows[rowIndex]
			var id int
			fmt.Sscanf(row[0], "%d", &id)

			tp.StartForm(4, []string{"Title", "Scenario Path", "Goal Config"}, id, func(values []string) error {
				var title *string
				if values[0] != "" {
					title = &values[0]
				}
				var goalConfig *string
				if values[2] != "" {
					goalConfig = &values[2]
				}
				return r.UpdateScenario(context.Background(), id, title, values[1], goalConfig)
			})
			return nil
		},
		OnDelete: func(rowIndex int) error {
			if rowIndex >= len(tp.currentRows) {
				return nil
			}
			row := tp.currentRows[rowIndex]
			var id int
			fmt.Sscanf(row[0], "%d", &id)
			return r.DeleteScenario(context.Background(), id)
		},
	})

	return tp
}

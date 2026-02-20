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
			{Title: "Description", Width: 30},
			{Title: "Scenario Path", Width: 30},
			{Title: "Param Path", Width: 30},
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
				description := ""
				if scenario.Description != nil && *scenario.Description != "" {
					description = *scenario.Description
				}
				paramPath := ""
				if scenario.ParamPath != nil && *scenario.ParamPath != "" {
					paramPath = *scenario.ParamPath
				}
				out = append(out, table.Row{
					fmt.Sprint(scenario.ID),
					title,
					description,
					scenario.ScenarioPath,
					paramPath,
				})
			}
			return out, nil
		},
	).(*TablePage)

	tp.WithCRUD(&CRUDCallbacks{
		OnCreate: func() error {
			tp.StartForm(4, []string{"Title", "Description", "Scenario Path", "Param Path"}, -1, func(values []string) error {
				var title *string
				if values[0] != "" {
					title = &values[0]
				}
				var description *string
				if values[1] != "" {
					description = &values[1]
				}
				var paramPath *string
				if values[3] != "" {
					paramPath = &values[3]
				}
				_, err := r.CreateScenario(context.Background(), title, description, values[2], paramPath)
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

			tp.StartForm(4, []string{"Title", "Description", "Scenario Path", "Param Path"}, id, func(values []string) error {
				var title *string
				if values[0] != "" {
					title = &values[0]
				}
				var description *string
				if values[1] != "" {
					description = &values[1]
				}
				var paramPath *string
				if values[3] != "" {
					paramPath = &values[3]
				}
				return r.UpdateScenario(context.Background(), id, title, description, values[2], paramPath)
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

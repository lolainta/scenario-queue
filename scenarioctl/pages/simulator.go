package pages

import (
	"context"
	"fmt"

	"scenarioctl/app"
	"scenarioctl/repo"

	"github.com/charmbracelet/bubbles/table"
)

func NewSimulatorPage(r *repo.Repo) app.Page {
	var tp *TablePage
	tp = NewTablePage(
		"Simulators",
		[]table.Column{
			{Title: "ID", Width: 6},
			{Title: "Name", Width: 20},
			{Title: "Image Path", Width: 25},
			{Title: "Config Path", Width: 25},
			{Title: "NV Runtime", Width: 12},
		},
		func(ctx context.Context) ([]table.Row, error) {
			rows, err := r.ListSimulators(ctx)
			if err != nil {
				return nil, err
			}

			var out []table.Row
			for _, sim := range rows {
				nvRuntime := "false"
				if sim.NvRuntime {
					nvRuntime = "true"
				}
				out = append(out, table.Row{
					fmt.Sprint(sim.ID),
					sim.Name,
					sim.ImagePath,
					sim.ConfigPath,
					nvRuntime,
				})
			}
			return out, nil
		},
	).(*TablePage)

	tp.WithCRUD(&CRUDCallbacks{
		OnCreate: func() error {
			tp.StartFormWithDefs([]FieldDef{
				{Label: "Name", Type: FieldTypeText},
				{Label: "Image Path", Type: FieldTypeText},
				{Label: "Config Path", Type: FieldTypeText},
				{Label: "NV Runtime", Type: FieldTypeSelect, Options: []SelectOption{{Label: "false", Value: "false"}, {Label: "true", Value: "true"}}},
			}, -1, func(values []string) error {
				nvRuntime := values[3] == "true"
				_, err := r.CreateSimulator(context.Background(), values[0], values[1], values[2], nvRuntime)
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

			tp.StartFormWithDefs([]FieldDef{
				{Label: "Name", Type: FieldTypeText},
				{Label: "Image Path", Type: FieldTypeText},
				{Label: "Config Path", Type: FieldTypeText},
				{Label: "NV Runtime", Type: FieldTypeSelect, Options: []SelectOption{{Label: "false", Value: "false"}, {Label: "true", Value: "true"}}},
			}, id, func(values []string) error {
				nvRuntime := values[3] == "true"
				return r.UpdateSimulator(context.Background(), id, values[0], values[1], values[2], nvRuntime)
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
			return r.DeleteSimulator(context.Background(), id)
		},
	})

	return tp
}

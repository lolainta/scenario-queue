package pages

import (
	"context"
	"fmt"

	"scenarioctl/app"
	"scenarioctl/repo"

	"github.com/charmbracelet/bubbles/table"
)

func NewSamplerPage(r *repo.Repo) app.Page {
	var tp *TablePage
	tp = NewTablePage(
		"Samplers",
		[]table.Column{
			{Title: "ID", Width: 6},
			{Title: "Name", Width: 20},
			{Title: "Module Path", Width: 30},
			{Title: "Config Path", Width: 30},
		},
		func(ctx context.Context) ([]table.Row, error) {
			rows, err := r.ListSamplers(ctx)
			if err != nil {
				return nil, err
			}

			var out []table.Row
			for _, sampler := range rows {
				configPath := ""
				if sampler.ConfigPath != nil {
					configPath = *sampler.ConfigPath
				}
				out = append(out, table.Row{
					fmt.Sprint(sampler.ID),
					sampler.Name,
					sampler.ModulePath,
					configPath,
				})
			}
			return out, nil
		},
	).(*TablePage)

	tp.WithCRUD(&CRUDCallbacks{
		OnCreate: func() error {
			tp.StartForm(3, []string{"Name", "Module Path", "Config Path"}, -1, func(values []string) error {
				var configPath *string
				if values[2] != "" {
					configPath = &values[2]
				}
				_, err := r.CreateSampler(context.Background(), values[0], values[1], configPath)
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

			tp.StartForm(3, []string{"Name", "Module Path", "Config Path"}, id, func(values []string) error {
				var configPath *string
				if values[2] != "" {
					configPath = &values[2]
				}
				return r.UpdateSampler(context.Background(), id, values[0], values[1], configPath)
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
			return r.DeleteSampler(context.Background(), id)
		},
	})

	return tp
}

package pages

import (
	"context"
	"fmt"

	"scenarioctl/app"
	"scenarioctl/repo"

	"github.com/charmbracelet/bubbles/table"
)

func NewAVPage(r *repo.Repo) app.Page {
	var tp *TablePage
	tp = NewTablePage(
		"AVs",
		[]table.Column{
			{Title: "ID", Width: 6},
			{Title: "Name", Width: 15},
			{Title: "Image Path", Width: 25},
			{Title: "Config Path", Width: 25},
		},
		func(ctx context.Context) ([]table.Row, error) {
			rows, err := r.ListAV(ctx)
			if err != nil {
				return nil, err
			}

			var out []table.Row
			for _, av := range rows {
				out = append(out, table.Row{
					fmt.Sprint(av.ID),
					av.Name,
					av.ImagePath,
					av.ConfigPath,
				})
			}
			return out, nil
		},
	).(*TablePage)

	tp.WithCRUD(&CRUDCallbacks{
		OnCreate: func() error {
			tp.StartForm(3, []string{"Name", "Image Path", "Config Path"}, -1, func(values []string) error {
				_, err := r.CreateAV(context.Background(), values[0], values[1], values[2], false)
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

			tp.StartForm(3, []string{"Name", "Image Path", "Config Path"}, id, func(values []string) error {
				return r.UpdateAV(context.Background(), id, values[0], values[1], values[2], false)
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
			return r.DeleteAV(context.Background(), id)
		},
	})

	return tp
}

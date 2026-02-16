package pages

import (
	"context"
	"fmt"

	"scenarioctl/app"
	"scenarioctl/repo"

	"github.com/charmbracelet/bubbles/table"
)

func NewMapPage(r *repo.Repo) app.Page {
	var tp *TablePage
	tp = NewTablePage(
		"Maps",
		[]table.Column{
			{Title: "ID", Width: 6},
			{Title: "Name", Width: 20},
			{Title: "XODR Path", Width: 30},
			{Title: "OSM Path", Width: 30},
		},
		func(ctx context.Context) ([]table.Row, error) {
			rows, err := r.ListMaps(ctx)
			if err != nil {
				return nil, err
			}

			var out []table.Row
			for _, m := range rows {
				xodrPath := ""
				if m.XodrPath != nil {
					xodrPath = *m.XodrPath
				}
				osmPath := ""
				if m.OsmPath != nil {
					osmPath = *m.OsmPath
				}
				out = append(out, table.Row{
					fmt.Sprint(m.ID),
					m.Name,
					xodrPath,
					osmPath,
				})
			}
			return out, nil
		},
	).(*TablePage)

	tp.WithCRUD(&CRUDCallbacks{
		OnCreate: func() error {
			tp.StartForm(3, []string{"Name", "XODR Path", "OSM Path"}, -1, func(values []string) error {
				var xodrPath *string
				if values[1] != "" {
					xodrPath = &values[1]
				}
				var osmPath *string
				if values[2] != "" {
					osmPath = &values[2]
				}
				_, err := r.CreateMap(context.Background(), values[0], xodrPath, osmPath)
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

			tp.StartForm(3, []string{"Name", "XODR Path", "OSM Path"}, id, func(values []string) error {
				var xodrPath *string
				if values[1] != "" {
					xodrPath = &values[1]
				}
				var osmPath *string
				if values[2] != "" {
					osmPath = &values[2]
				}
				return r.UpdateMap(context.Background(), id, values[0], xodrPath, osmPath)
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
			return r.DeleteMap(context.Background(), id)
		},
	})

	return tp
}

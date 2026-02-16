package pages

import (
	"context"
	"fmt"
	"strconv"

	"scenarioctl/app"
	"scenarioctl/repo"

	"github.com/charmbracelet/bubbles/table"
)

func NewPlanPage(r *repo.Repo) app.Page {
	var tp *TablePage
	tp = NewTablePage(
		"Plans",
		[]table.Column{
			{Title: "ID", Width: 6},
			{Title: "Name", Width: 20},
			{Title: "Map", Width: 20},
			{Title: "Scenario", Width: 30},
		},
		func(ctx context.Context) ([]table.Row, error) {
			rows, err := r.ListPlans(ctx)
			if err != nil {
				return nil, err
			}

			var out []table.Row
			for _, plan := range rows {
				out = append(out, table.Row{
					fmt.Sprint(plan.ID),
					plan.Name,
					plan.Map,
					plan.Scenario,
				})
			}
			return out, nil
		},
	).(*TablePage)

	tp.WithCRUD(&CRUDCallbacks{
		OnCreate: func() error {
			// Fetch available maps and scenarios for selection
			ctx := context.Background()
			maps, err := r.ListMaps(ctx)
			if err != nil {
				return err
			}
			scenarios, err := r.ListScenarios(ctx)
			if err != nil {
				return err
			}

			// Build options for select fields (display name + value ID)
			var mapOptions []string
			for _, m := range maps {
				mapOptions = append(mapOptions, m.Name, strconv.Itoa(m.ID))
			}
			var scenarioOptions []string
			for _, s := range scenarios {
				title := ""
				if s.Title != nil {
					title = *s.Title
				}
				scenarioOptions = append(scenarioOptions, title, strconv.Itoa(s.ID))
			}

			// Create form with text field for name and select fields for map/scenario
			fieldDefs := []FieldDef{
				{Label: "Name", Type: FieldTypeText},
				{Label: "Map", Type: FieldTypeSelect, Options: mapOptions},
				{Label: "Scenario", Type: FieldTypeSelect, Options: scenarioOptions},
			}

			tp.StartFormWithDefs(fieldDefs, -1, func(values []string) error {
				mapID, _ := strconv.Atoi(values[1])
				scenarioID, _ := strconv.Atoi(values[2])
				_, err := r.CreatePlan(context.Background(), values[0], mapID, scenarioID)
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

			// Fetch available maps and scenarios for selection
			ctx := context.Background()
			maps, err := r.ListMaps(ctx)
			if err != nil {
				return err
			}
			scenarios, err := r.ListScenarios(ctx)
			if err != nil {
				return err
			}

			// Build options for select fields
			var mapOptions []string
			currentMapIdx := 0
			for i, m := range maps {
				mapOptions = append(mapOptions, m.Name, strconv.Itoa(m.ID))
				if m.Name == row[2] {
					currentMapIdx = i
				}
			}
			var scenarioOptions []string
			currentScenarioIdx := 0
			for i, s := range scenarios {
				title := ""
				if s.Title != nil {
					title = *s.Title
				}
				scenarioOptions = append(scenarioOptions, title, strconv.Itoa(s.ID))
				if title == row[3] {
					currentScenarioIdx = i
				}
			}

			// Create form with pre-populated data
			fieldDefs := []FieldDef{
				{Label: "Name", Type: FieldTypeText},
				{Label: "Map", Type: FieldTypeSelect, Options: mapOptions},
				{Label: "Scenario", Type: FieldTypeSelect, Options: scenarioOptions},
			}

			tp.StartFormWithDefs(fieldDefs, id, func(values []string) error {
				mapID, _ := strconv.Atoi(values[1])
				scenarioID, _ := strconv.Atoi(values[2])
				return r.UpdatePlan(context.Background(), id, values[0], mapID, scenarioID)
			})

			// Pre-populate text field and select indices
			if tp.form != nil && len(tp.form.fields) > 0 {
				tp.form.fields[0].SetValue(row[1]) // Name
				tp.form.selectIndex[1] = currentMapIdx
				tp.form.selectIndex[2] = currentScenarioIdx
				if currentMapIdx < len(mapOptions) {
					tp.form.fields[1].SetValue(mapOptions[currentMapIdx*2])
				}
				if currentScenarioIdx < len(scenarioOptions) {
					tp.form.fields[2].SetValue(scenarioOptions[currentScenarioIdx*2])
				}
			}

			return nil
		},
		OnDelete: func(rowIndex int) error {
			if rowIndex >= len(tp.currentRows) {
				return nil
			}
			row := tp.currentRows[rowIndex]
			var id int
			fmt.Sscanf(row[0], "%d", &id)
			return r.DeletePlan(context.Background(), id)
		},
	})

	return tp
}

package pages

import (
	"context"
	"fmt"
	"strconv"

	"scenarioctl/app"
	"scenarioctl/repo"

	"github.com/charmbracelet/bubbles/table"
)

func NewTaskPage(r *repo.Repo) app.Page {
	var tp *TablePage
	tp = NewTablePage(
		"Tasks",
		[]table.Column{
			{Title: "ID", Width: 6},
			{Title: "Worker ID", Width: 10},
			{Title: "Status", Width: 12},
			{Title: "Plan", Width: 15},
			{Title: "AV", Width: 15},
			{Title: "Simulator", Width: 15},
		},
		func(ctx context.Context) ([]table.Row, error) {
			tasks, err := r.ListTasks(ctx)
			if err != nil {
				return nil, err
			}

			var out []table.Row
			for _, task := range tasks {
				workerID := ""
				if task.WorkerID != nil {
					workerID = strconv.Itoa(*task.WorkerID)
				}
				out = append(out, table.Row{
					fmt.Sprint(task.ID),
					task.Status,
					task.Plan,
					task.AV,
					task.Simulator,
					workerID,
				})
			}
			return out, nil
		},
	).(*TablePage)

	tp.WithCRUD(&CRUDCallbacks{
		OnCreate: func() error {
			ctx := context.Background()

			// Fetch all required data for selections
			plans, err := r.ListPlans(ctx)
			if err != nil {
				return err
			}
			avs, err := r.ListAV(ctx)
			if err != nil {
				return err
			}
			simulators, err := r.ListSimulators(ctx)
			if err != nil {
				return err
			}
			samplers, err := r.ListSamplers(ctx)
			if err != nil {
				return err
			}

			// Build options for select fields
			var planOptions []SelectOption
			for _, p := range plans {
				planOptions = append(planOptions, SelectOption{Label: p.Name, Value: strconv.Itoa(p.ID)})
			}

			var avOptions []SelectOption
			for _, a := range avs {
				avOptions = append(avOptions, SelectOption{Label: a.Name, Value: strconv.Itoa(a.ID)})
			}

			var simulatorOptions []SelectOption
			for _, s := range simulators {
				simulatorOptions = append(simulatorOptions, SelectOption{Label: s.Name, Value: strconv.Itoa(s.ID)})
			}

			var samplerOptions []SelectOption
			for _, s := range samplers {
				samplerOptions = append(samplerOptions, SelectOption{Label: s.Name, Value: strconv.Itoa(s.ID)})
			}

			// Create form with only required fields (no status, no worker_id)
			// Status defaults to "pending", Worker ID defaults to null
			fieldDefs := []FieldDef{
				{Label: "Plan", Type: FieldTypeSelect, Options: planOptions},
				{Label: "AV", Type: FieldTypeSelect, Options: avOptions},
				{Label: "Simulator", Type: FieldTypeSelect, Options: simulatorOptions},
				{Label: "Sampler", Type: FieldTypeSelect, Options: samplerOptions},
			}

			tp.StartFormWithDefs(fieldDefs, -1, func(values []string) error {
				planID, _ := strconv.Atoi(values[0])
				avID, _ := strconv.Atoi(values[1])
				simulatorID, _ := strconv.Atoi(values[2])
				samplerID, _ := strconv.Atoi(values[3])
				status := "pending" // Default status
				var workerID *int   // Defaults to null

				_, err := r.CreateTask(context.Background(), planID, avID, simulatorID, samplerID, status, workerID)
				return err
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
			return r.DeleteTask(context.Background(), id)
		},
		OnUpdate: func(rowIndex int) error {
			if rowIndex >= len(tp.currentRows) {
				return nil
			}
			row := tp.currentRows[rowIndex]
			var id int
			fmt.Sscanf(row[0], "%d", &id)

			ctx := context.Background()

			// Fetch all tasks to get current values
			allTasks, err := r.ListTasks(ctx)
			if err != nil {
				return err
			}

			var currentTask *repo.TaskRow
			for i := range allTasks {
				if allTasks[i].ID == id {
					currentTask = &allTasks[i]
					break
				}
			}
			if currentTask == nil {
				return fmt.Errorf("task not found")
			}

			// Fetch all required data for selections
			plans, err := r.ListPlans(ctx)
			if err != nil {
				return err
			}
			avs, err := r.ListAV(ctx)
			if err != nil {
				return err
			}
			simulators, err := r.ListSimulators(ctx)
			if err != nil {
				return err
			}
			samplers, err := r.ListSamplers(ctx)
			if err != nil {
				return err
			}

			// Build options for select fields
			var planOptions []SelectOption
			var currentPlanIdx int
			for i, p := range plans {
				planOptions = append(planOptions, SelectOption{Label: p.Name, Value: strconv.Itoa(p.ID)})
				if p.Name == currentTask.Plan {
					currentPlanIdx = i
				}
			}

			var avOptions []SelectOption
			var currentAVIdx int
			for i, a := range avs {
				avOptions = append(avOptions, SelectOption{Label: a.Name, Value: strconv.Itoa(a.ID)})
				if a.Name == currentTask.AV {
					currentAVIdx = i
				}
			}

			var simulatorOptions []SelectOption
			var currentSimulatorIdx int
			for i, s := range simulators {
				simulatorOptions = append(simulatorOptions, SelectOption{Label: s.Name, Value: strconv.Itoa(s.ID)})
				if s.Name == currentTask.Simulator {
					currentSimulatorIdx = i
				}
			}

			var samplerOptions []SelectOption
			var currentSamplerIdx int
			for _, s := range samplers {
				samplerOptions = append(samplerOptions, SelectOption{Label: s.Name, Value: strconv.Itoa(s.ID)})
				// Note: samplers don't have a display name in TaskRow, so we skip current matching
			}

			// Status options and index
			statusOptions := []SelectOption{
				{Label: "Pending", Value: "pending"},
				{Label: "In Progress", Value: "in_progress"},
				{Label: "Completed", Value: "completed"},
				{Label: "Failed", Value: "failed"},
			}
			var currentStatusIdx int
			statusValues := []string{"pending", "in_progress", "completed", "failed"}
			for i, s := range statusValues {
				if s == currentTask.Status {
					currentStatusIdx = i
					break
				}
			}

			// Worker ID field
			currentWorkerID := ""
			if currentTask.WorkerID != nil {
				currentWorkerID = strconv.Itoa(*currentTask.WorkerID)
			}

			// Create form with all editable fields
			fieldDefs := []FieldDef{
				{Label: "Plan", Type: FieldTypeSelect, Options: planOptions},
				{Label: "AV", Type: FieldTypeSelect, Options: avOptions},
				{Label: "Simulator", Type: FieldTypeSelect, Options: simulatorOptions},
				{Label: "Sampler", Type: FieldTypeSelect, Options: samplerOptions},
				{Label: "Status", Type: FieldTypeSelect, Options: statusOptions},
				{Label: "Worker ID", Type: FieldTypeText},
			}

			tp.StartFormWithDefs(fieldDefs, id, func(values []string) error {
				planID, _ := strconv.Atoi(values[0])
				avID, _ := strconv.Atoi(values[1])
				simulatorID, _ := strconv.Atoi(values[2])
				samplerID, _ := strconv.Atoi(values[3])
				status := values[4]

				var workerID *int
				if values[5] != "" {
					id, _ := strconv.Atoi(values[5])
					workerID = &id
				}

				return r.UpdateTask(context.Background(), id, planID, avID, simulatorID, samplerID, status, workerID)
			})

			// Pre-populate form with current values
			if tp.form != nil {
				tp.form.selectIndex[0] = currentPlanIdx
				tp.form.selectIndex[1] = currentAVIdx
				tp.form.selectIndex[2] = currentSimulatorIdx
				tp.form.selectIndex[3] = currentSamplerIdx
				tp.form.selectIndex[4] = currentStatusIdx

				if len(tp.form.fields) > 0 {
					tp.form.fields[0].SetValue(planOptions[currentPlanIdx].Label)
				}
				if len(tp.form.fields) > 1 {
					tp.form.fields[1].SetValue(avOptions[currentAVIdx].Label)
				}
				if len(tp.form.fields) > 2 {
					tp.form.fields[2].SetValue(simulatorOptions[currentSimulatorIdx].Label)
				}
				if len(tp.form.fields) > 3 {
					tp.form.fields[3].SetValue(samplerOptions[currentSamplerIdx].Label)
				}
				if len(tp.form.fields) > 4 {
					tp.form.fields[4].SetValue(statusOptions[currentStatusIdx].Label)
				}
				if len(tp.form.fields) > 5 {
					tp.form.fields[5].SetValue(currentWorkerID)
				}
			}

			return nil
		},
	})

	return tp
}

package repo

import "context"

type PlanRow struct {
	ID       int
	Name     string
	Map      string
	Scenario string
}

func (r *Repo) ListPlans(ctx context.Context) ([]PlanRow, error) {
	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.name, m.name, s.title
		FROM plan p
		JOIN map m ON p.map_id = m.id
		JOIN scenario s ON p.scenario_id = s.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []PlanRow
	for rows.Next() {
		var x PlanRow
		rows.Scan(&x.ID, &x.Name, &x.Map, &x.Scenario)
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r *Repo) CreatePlan(ctx context.Context, name string, mapID int, scenarioID int) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, `
		INSERT INTO plan (name, map_id, scenario_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`, name, mapID, scenarioID).Scan(&id)
	return id, err
}

func (r *Repo) UpdatePlan(ctx context.Context, id int, name string, mapID int, scenarioID int) error {
	_, err := r.db.Exec(ctx, `
		UPDATE plan
		SET name = $1, map_id = $2, scenario_id = $3
		WHERE id = $4
	`, name, mapID, scenarioID, id)
	return err
}

func (r *Repo) DeletePlan(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM plan WHERE id = $1`, id)
	return err
}

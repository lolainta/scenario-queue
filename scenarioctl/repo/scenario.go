package repo

import (
	"context"
)

type Scenario struct {
	ID           int
	Title        *string
	Description  *string
	ScenarioPath string
	ParamPath    *string
}

func (r *Repo) ListScenarios(ctx context.Context) ([]Scenario, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, title, description, scenario_path, param_path
		FROM scenario
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Scenario
	for rows.Next() {
		var s Scenario
		err := rows.Scan(
			&s.ID,
			&s.Title,
			&s.Description,
			&s.ScenarioPath,
			&s.ParamPath,
		)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}

	return out, rows.Err()
}

func (r *Repo) CreateScenario(ctx context.Context, title *string, description *string, scenarioPath string, paramPath *string) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, `
		INSERT INTO scenario (title, description, scenario_path, param_path)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, title, description, scenarioPath, paramPath).Scan(&id)
	return id, err
}

func (r *Repo) UpdateScenario(ctx context.Context, id int, title *string, description *string, scenarioPath string, paramPath *string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE scenario
		SET title = $1, description = $2, scenario_path = $3, param_path = $4
		WHERE id = $5
	`, title, description, scenarioPath, paramPath, id)
	return err
}

func (r *Repo) DeleteScenario(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM scenario WHERE id = $1`, id)
	return err
}

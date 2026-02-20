package repo

import (
	"context"
)

type Scenario struct {
	ID           int
	Title        *string
	ScenarioPath string
	GoalConfig   *string
}

func (r *Repo) ListScenarios(ctx context.Context) ([]Scenario, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, title, scenario_path, goal_config
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
			&s.ScenarioPath,
			&s.GoalConfig,
		)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}

	return out, rows.Err()
}

func (r *Repo) CreateScenario(ctx context.Context, title *string, scenarioPath string, goalConfig *string) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, `
		INSERT INTO scenario (title, scenario_path, goal_config)
		VALUES ($1, $2, $3)
		RETURNING id
	`, title, scenarioPath, goalConfig).Scan(&id)
	return id, err
}

func (r *Repo) UpdateScenario(ctx context.Context, id int, title *string, scenarioPath string, goalConfig *string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE scenario
		SET title = $1, scenario_path = $2, goal_config = $3
		WHERE id = $4
	`, title, scenarioPath, goalConfig, id)
	return err
}

func (r *Repo) DeleteScenario(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM scenario WHERE id = $1`, id)
	return err
}

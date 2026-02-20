package repo

import "context"

type SimulatorRow struct {
	ID         int
	Name       string
	ImagePath  string
	ConfigPath string
	NvRuntime  bool
}

func (r *Repo) ListSimulators(ctx context.Context) ([]SimulatorRow, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, image_path, config_path, nv_runtime FROM simulator ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []SimulatorRow
	for rows.Next() {
		var x SimulatorRow
		rows.Scan(&x.ID, &x.Name, &x.ImagePath, &x.ConfigPath, &x.NvRuntime)
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r *Repo) CreateSimulator(ctx context.Context, name, imagePath, configPath string, nvRuntime bool) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, `
		INSERT INTO simulator (name, image_path, config_path, nv_runtime)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, name, imagePath, configPath, nvRuntime).Scan(&id)
	return id, err
}

func (r *Repo) UpdateSimulator(ctx context.Context, id int, name, imagePath, configPath string, nvRuntime bool) error {
	_, err := r.db.Exec(ctx, `
		UPDATE simulator
		SET name = $1, image_path = $2, config_path = $3, nv_runtime = $4
		WHERE id = $5
	`, name, imagePath, configPath, nvRuntime, id)
	return err
}

func (r *Repo) DeleteSimulator(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM simulator WHERE id = $1`, id)
	return err
}

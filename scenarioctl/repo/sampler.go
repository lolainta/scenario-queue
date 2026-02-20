package repo

import "context"

type SamplerRow struct {
	ID         int
	Name       string
	ModulePath string
	ConfigPath *string
}

func (r *Repo) ListSamplers(ctx context.Context) ([]SamplerRow, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, module_path, config_path FROM sampler ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []SamplerRow
	for rows.Next() {
		var x SamplerRow
		rows.Scan(&x.ID, &x.Name, &x.ModulePath, &x.ConfigPath)
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r *Repo) CreateSampler(ctx context.Context, name, modulePath string, configPath *string) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, `
		INSERT INTO sampler (name, module_path, config_path)
		VALUES ($1, $2, $3)
		RETURNING id
	`, name, modulePath, configPath).Scan(&id)
	return id, err
}

func (r *Repo) UpdateSampler(ctx context.Context, id int, name, modulePath string, configPath *string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE sampler
		SET name = $1, module_path = $2, config_path = $3
		WHERE id = $4
	`, name, modulePath, configPath, id)
	return err
}

func (r *Repo) DeleteSampler(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM sampler WHERE id = $1`, id)
	return err
}

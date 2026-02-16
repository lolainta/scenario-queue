package repo

import "context"

type MapRow struct {
	ID       int
	Name     string
	XodrPath *string
	OsmPath  *string
}

func (r *Repo) ListMaps(ctx context.Context) ([]MapRow, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, xodr_path, osm_path FROM map ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []MapRow
	for rows.Next() {
		var x MapRow
		rows.Scan(&x.ID, &x.Name, &x.XodrPath, &x.OsmPath)
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r *Repo) CreateMap(ctx context.Context, name string, xodrPath *string, osmPath *string) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, `
		INSERT INTO map (name, xodr_path, osm_path)
		VALUES ($1, $2, $3)
		RETURNING id
	`, name, xodrPath, osmPath).Scan(&id)
	return id, err
}

func (r *Repo) UpdateMap(ctx context.Context, id int, name string, xodrPath *string, osmPath *string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE map
		SET name = $1, xodr_path = $2, osm_path = $3
		WHERE id = $4
	`, name, xodrPath, osmPath, id)
	return err
}

func (r *Repo) DeleteMap(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM map WHERE id = $1`, id)
	return err
}

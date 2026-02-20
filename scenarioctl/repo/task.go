package repo

import "context"

type TaskRow struct {
	ID        int
	Status    string
	Plan      string
	AV        string
	Simulator string
	WorkerID  *int
}

func (r *Repo) ListTasks(ctx context.Context) ([]TaskRow, error) {
	rows, err := r.db.Query(ctx, `
		SELECT t.id, t.status, p.name, a.name, s.name, t.worker_id
		FROM task t
		JOIN plan p ON t.plan_id = p.id
		JOIN av a ON t.av_id = a.id
		JOIN simulator s ON t.simulator_id = s.id
		ORDER BY t.id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []TaskRow
	for rows.Next() {
		var x TaskRow
		rows.Scan(&x.ID, &x.Status, &x.Plan, &x.AV, &x.Simulator, &x.WorkerID)
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r *Repo) UpdateTaskStatus(ctx context.Context, id int, status string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE task
		SET status = $1
		WHERE id = $2
	`, status, id)
	return err
}

func (r *Repo) UpdateTask(ctx context.Context, id int, planID, avID, simulatorID, samplerID int, status string, workerID *int) error {
	_, err := r.db.Exec(ctx, `
		UPDATE task
		SET plan_id = $1, av_id = $2, simulator_id = $3, sampler_id = $4, status = $5, worker_id = $6
		WHERE id = $7
	`, planID, avID, simulatorID, samplerID, status, workerID, id)
	return err
}

func (r *Repo) CreateTask(ctx context.Context, planID, avID, simulatorID, samplerID int, status string, workerID *int) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, `
		INSERT INTO task (plan_id, av_id, simulator_id, sampler_id, status, worker_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, planID, avID, simulatorID, samplerID, status, workerID).Scan(&id)
	return id, err
}

func (r *Repo) DeleteTask(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM task WHERE id = $1`, id)
	return err
}

package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"erpagendamentos/backend/internal/domain/job"
	"github.com/google/uuid"
)

type JobRepository struct {
	db *sql.DB
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) Enqueue(j job.Job) (job.Job, error) {
	if j.ID == "" {
		j.ID = uuid.NewString()
	}
	if j.Status == "" {
		j.Status = "pending"
	}
	if j.MaxAttempts == 0 {
		j.MaxAttempts = 5
	}
	now := time.Now().UTC()
	var runAt time.Time
	if j.RunAt != "" {
		if t, err := time.Parse(time.RFC3339, j.RunAt); err == nil {
			runAt = t
		}
	}
	if runAt.IsZero() {
		runAt = now
	}

	payload := j.Payload
	if payload == "" {
		payload = "{}"
	}

	err := r.db.QueryRow(`
		INSERT INTO catalog.jobs (id, tenant_id, type, payload, status, attempts, max_attempts, run_at, last_error, created_at, updated_at)
		VALUES ($1, $2, $3, $4::jsonb, $5, $6, $7, $8, $9, $10, $11)
		RETURNING created_at, updated_at
	`, j.ID, j.TenantID, j.Type, payload, j.Status, j.Attempts, j.MaxAttempts, runAt, j.LastError, now, now).Scan(&j.CreatedAt, &j.UpdatedAt)
	if err != nil {
		return job.Job{}, fmt.Errorf("enqueue job: %w", err)
	}
	j.RunAt = runAt.Format(time.RFC3339)
	return j, nil
}

func (r *JobRepository) FetchPending(limit int) ([]job.Job, error) {
	if limit <= 0 {
		limit = 10
	}
	rows, err := r.db.Query(`
		SELECT id, tenant_id, type, payload::text, status, attempts, max_attempts, run_at, last_error, created_at, updated_at
		FROM catalog.jobs
		WHERE status IN ('pending', 'failed') AND attempts < max_attempts AND run_at <= NOW()
		ORDER BY run_at ASC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("fetch pending jobs: %w", err)
	}
	defer rows.Close()

	var jobs []job.Job
	for rows.Next() {
		var j job.Job
		var runAt, created, updated time.Time
		if err := rows.Scan(&j.ID, &j.TenantID, &j.Type, &j.Payload, &j.Status, &j.Attempts, &j.MaxAttempts, &runAt, &j.LastError, &created, &updated); err != nil {
			return nil, err
		}
		j.RunAt = runAt.Format(time.RFC3339)
		j.CreatedAt = created.Format(time.RFC3339)
		j.UpdatedAt = updated.Format(time.RFC3339)
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func (r *JobRepository) MarkRunning(id string) error {
	now := time.Now().UTC()
	result, err := r.db.Exec(`
		UPDATE catalog.jobs
		SET status = 'running', attempts = attempts + 1, updated_at = $1
		WHERE id = $2
	`, now, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("job not found")
	}
	return nil
}

func (r *JobRepository) MarkCompleted(id string) error {
	now := time.Now().UTC()
	_, err := r.db.Exec(`
		UPDATE catalog.jobs
		SET status = 'completed', updated_at = $1
		WHERE id = $2
	`, now, id)
	return err
}

func (r *JobRepository) MarkFailed(id string, errMessage string, retryAfterSeconds int) error {
	now := time.Now().UTC()
	runAt := now.Add(time.Duration(retryAfterSeconds) * time.Second)
	_, err := r.db.Exec(`
		UPDATE catalog.jobs
		SET status = 'failed', last_error = $1, run_at = $2, updated_at = $3
		WHERE id = $4
	`, errMessage, runAt, now, id)
	return err
}

func (r *JobRepository) CountByStatus() (map[string]int, error) {
	rows, err := r.db.Query(`
		SELECT status, COUNT(*)
		FROM catalog.jobs
		GROUP BY status
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err == nil {
			counts[status] = count
		}
	}
	return counts, nil
}


package memory

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"erpagendamentos/backend/internal/domain/job"
)

type JobRepository struct {
	mu   sync.RWMutex
	jobs map[string]job.Job
}

func NewJobRepository() *JobRepository {
	return &JobRepository{
		jobs: make(map[string]job.Job),
	}
}

func (r *JobRepository) Enqueue(j job.Job) (job.Job, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if j.ID == "" {
		j.ID = fmt.Sprintf("job_%d", time.Now().UnixNano())
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if j.RunAt == "" {
		j.RunAt = now
	}
	if j.Status == "" {
		j.Status = "pending"
	}
	if j.MaxAttempts == 0 {
		j.MaxAttempts = 5
	}
	j.CreatedAt = now
	j.UpdatedAt = now
	r.jobs[j.ID] = j
	return j, nil
}

func (r *JobRepository) FetchPending(limit int) ([]job.Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	now := time.Now().UTC().Format(time.RFC3339)
	var result []job.Job
	for _, j := range r.jobs {
		if (j.Status == "pending" || j.Status == "failed") && j.Attempts < j.MaxAttempts && j.RunAt <= now {
			result = append(result, j)
			if limit > 0 && len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func (r *JobRepository) MarkRunning(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	j, found := r.jobs[id]
	if !found {
		return errors.New("job not found")
	}
	j.Status = "running"
	j.Attempts++
	j.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	r.jobs[id] = j
	return nil
}

func (r *JobRepository) MarkCompleted(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	j, found := r.jobs[id]
	if !found {
		return errors.New("job not found")
	}
	j.Status = "completed"
	j.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	r.jobs[id] = j
	return nil
}

func (r *JobRepository) MarkFailed(id string, errMessage string, retryAfterSeconds int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	j, found := r.jobs[id]
	if !found {
		return errors.New("job not found")
	}
	j.Status = "failed"
	j.LastError = errMessage
	j.RunAt = time.Now().Add(time.Duration(retryAfterSeconds) * time.Second).UTC().Format(time.RFC3339)
	j.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	r.jobs[id] = j
	return nil
}

func (r *JobRepository) CountByStatus() (map[string]int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	counts := make(map[string]int)
	for _, j := range r.jobs {
		counts[j.Status]++
	}
	return counts, nil
}


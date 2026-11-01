package job

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"erpagendamentos/backend/internal/domain/job"
)

type HandlerFunc func(ctx context.Context, j job.Job) error

type Queue struct {
	repo     job.Repository
	handlers map[string]HandlerFunc
	mu       sync.RWMutex
	stopCh   chan struct{}
	wg       sync.WaitGroup
	running  bool
}

func NewQueue(repo job.Repository) *Queue {
	return &Queue{
		repo:     repo,
		handlers: make(map[string]HandlerFunc),
		stopCh:   make(chan struct{}),
	}
}

func (q *Queue) RegisterHandler(jobType string, handler HandlerFunc) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.handlers[jobType] = handler
}

func (q *Queue) Enqueue(tenantID string, jobType string, payload any, delay time.Duration) (job.Job, error) {
	var payloadBytes []byte
	var err error
	if payload != nil {
		switch p := payload.(type) {
		case string:
			payloadBytes = []byte(p)
		default:
			payloadBytes, err = json.Marshal(p)
			if err != nil {
				return job.Job{}, fmt.Errorf("marshal job payload: %w", err)
			}
		}
	} else {
		payloadBytes = []byte("{}")
	}

	runAt := time.Now().Add(delay).UTC().Format(time.RFC3339)
	return q.repo.Enqueue(job.Job{
		TenantID: tenantID,
		Type:     jobType,
		Payload:  string(payloadBytes),
		RunAt:    runAt,
	})
}

func (q *Queue) Start(concurrency int, pollInterval time.Duration) {
	q.mu.Lock()
	if q.running {
		q.mu.Unlock()
		return
	}
	q.running = true
	q.mu.Unlock()

	if concurrency <= 0 {
		concurrency = 2
	}
	if pollInterval <= 0 {
		pollInterval = 5 * time.Second
	}

	jobCh := make(chan job.Job, concurrency*2)

	// Start workers
	for i := 0; i < concurrency; i++ {
		q.wg.Add(1)
		go func(workerID int) {
			defer q.wg.Done()
			for {
				select {
				case <-q.stopCh:
					return
				case j, ok := <-jobCh:
					if !ok {
						return
					}
					q.processJob(j)
				}
			}
		}(i)
	}

	// Start poller
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-q.stopCh:
				close(jobCh)
				return
			case <-ticker.C:
				jobs, err := q.repo.FetchPending(concurrency * 2)
				if err != nil {
					continue
				}
				for _, j := range jobs {
					if err := q.repo.MarkRunning(j.ID); err == nil {
						select {
						case jobCh <- j:
						case <-q.stopCh:
							close(jobCh)
							return
						}
					}
				}
			}
		}
	}()
}

func (q *Queue) Stop() {
	q.mu.Lock()
	if !q.running {
		q.mu.Unlock()
		return
	}
	q.running = false
	close(q.stopCh)
	q.mu.Unlock()

	q.wg.Wait()
}

func (q *Queue) processJob(j job.Job) {
	q.mu.RLock()
	handler, exists := q.handlers[j.Type]
	q.mu.RUnlock()

	if !exists {
		_ = q.repo.MarkFailed(j.ID, fmt.Sprintf("no handler registered for job type %s", j.Type), 60)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err := handler(ctx, j); err != nil {
		log.Printf(`{"level":"WARN","message":"job execution failed","job_id":"%s","job_type":"%s","attempts":%d,"error":"%v"}`, j.ID, j.Type, j.Attempts, err)
		retryBackoff := (j.Attempts + 1) * 30
		_ = q.repo.MarkFailed(j.ID, err.Error(), retryBackoff)
	} else {
		_ = q.repo.MarkCompleted(j.ID)
	}
}

func (q *Queue) CountByStatus() (map[string]int, error) {
	return q.repo.CountByStatus()
}


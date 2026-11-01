package job_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"erpagendamentos/backend/internal/application/job"
	jobdomain "erpagendamentos/backend/internal/domain/job"
	"erpagendamentos/backend/internal/infrastructure/persistence/memory"
)

func TestJobQueue(t *testing.T) {
	repo := memory.NewJobRepository()
	queue := job.NewQueue(repo)

	var executed int64
	queue.RegisterHandler("test_task", func(ctx context.Context, j jobdomain.Job) error {
		atomic.AddInt64(&executed, 1)
		return nil
	})

	// Enqueue job with 0 delay
	_, err := queue.Enqueue("tenant-1", "test_task", map[string]string{"key": "value"}, 0)
	if err != nil {
		t.Fatalf("enqueue error: %v", err)
	}

	queue.Start(2, 50*time.Millisecond)
	time.Sleep(200 * time.Millisecond)
	queue.Stop()

	if atomic.LoadInt64(&executed) == 0 {
		t.Fatalf("expected job to be executed at least once")
	}

	counts, err := queue.CountByStatus()
	if err != nil {
		t.Fatalf("count error: %v", err)
	}
	if counts["completed"] == 0 {
		t.Fatalf("expected completed count > 0, got %v", counts)
	}
}


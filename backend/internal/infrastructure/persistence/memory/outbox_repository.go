package memory

import (
	"sync"
	"time"

	outboxdomain "erpagendamentos/backend/internal/domain/outbox"
)

type OutboxRepository struct {
	mu     sync.RWMutex
	events []outboxdomain.OutboxEvent
}

func NewOutboxRepository() *OutboxRepository {
	return &OutboxRepository{
		events: make([]outboxdomain.OutboxEvent, 0),
	}
}

func (r *OutboxRepository) Save(event outboxdomain.OutboxEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, event)
	return nil
}

func (r *OutboxRepository) ListPending(limit int) ([]outboxdomain.OutboxEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var pending []outboxdomain.OutboxEvent
	for _, e := range r.events {
		if e.Status == "pending" {
			pending = append(pending, e)
			if len(pending) >= limit {
				break
			}
		}
	}
	return pending, nil
}

func (r *OutboxRepository) MarkPublished(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	for i := range r.events {
		if r.events[i].ID == id {
			r.events[i].Status = "published"
			r.events[i].ProcessedAt = &now
			return nil
		}
	}
	return nil
}

func (r *OutboxRepository) MarkFailed(id string, errMsg string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	for i := range r.events {
		if r.events[i].ID == id {
			r.events[i].Status = "failed"
			r.events[i].RetryCount++
			r.events[i].ErrorMessage = errMsg
			r.events[i].ProcessedAt = &now
			return nil
		}
	}
	return nil
}

func (r *OutboxRepository) CountByStatus() (map[string]int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	counts := map[string]int{
		"pending":   0,
		"published": 0,
		"failed":    0,
	}
	for _, e := range r.events {
		counts[e.Status]++
	}
	return counts, nil
}

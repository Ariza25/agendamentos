package outbox

import (
	"context"
	"time"

	"github.com/google/uuid"
	outboxdomain "erpagendamentos/backend/internal/domain/outbox"
)

type EventDispatcher func(tenantID string, topic string, payload map[string]any) error

type UseCase struct {
	repo       outboxdomain.Repository
	dispatcher EventDispatcher
}

func NewUseCase(repo outboxdomain.Repository, dispatcher EventDispatcher) *UseCase {
	return &UseCase{
		repo:       repo,
		dispatcher: dispatcher,
	}
}

func (u *UseCase) Publish(tenantID string, topic string, payload map[string]any) error {
	if u == nil || u.repo == nil {
		return nil
	}
	event := outboxdomain.OutboxEvent{
		ID:        uuid.NewString(),
		TenantID:  tenantID,
		Topic:     topic,
		Payload:   payload,
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	return u.repo.Save(event)
}

func (u *UseCase) DispatchPending(limit int) (int, error) {
	if u == nil || u.repo == nil {
		return 0, nil
	}
	events, err := u.repo.ListPending(limit)
	if err != nil {
		return 0, err
	}

	dispatched := 0
	for _, e := range events {
		var dispatchErr error
		if u.dispatcher != nil {
			dispatchErr = u.dispatcher(e.TenantID, e.Topic, e.Payload)
		}
		if dispatchErr == nil {
			_ = u.repo.MarkPublished(e.ID)
			dispatched++
		} else {
			_ = u.repo.MarkFailed(e.ID, dispatchErr.Error())
		}
	}
	return dispatched, nil
}

func (u *UseCase) Stats() (map[string]int, error) {
	if u == nil || u.repo == nil {
		return map[string]int{"pending": 0, "published": 0, "failed": 0}, nil
	}
	return u.repo.CountByStatus()
}

// StartWorker starts a lightweight background goroutine that polls pending outbox events
func (u *UseCase) StartWorker(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = 5 * time.Second
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_, _ = u.DispatchPending(25)
			}
		}
	}()
}

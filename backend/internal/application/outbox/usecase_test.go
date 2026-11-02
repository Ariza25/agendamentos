package outbox

import (
	"testing"

	"erpagendamentos/backend/internal/infrastructure/persistence/memory"
)

func TestOutboxPublishAndDispatch(t *testing.T) {
	repo := memory.NewOutboxRepository()
	var dispatchedEvents []string

	uc := NewUseCase(repo, func(tenantID string, topic string, payload map[string]any) error {
		dispatchedEvents = append(dispatchedEvents, topic)
		return nil
	})

	err := uc.Publish("tenant-1", "appointment.created", map[string]any{"id": "apt-123"})
	if err != nil {
		t.Fatalf("publish error: %v", err)
	}

	err = uc.Publish("tenant-1", "checkout.completed", map[string]any{"id": "chk-456"})
	if err != nil {
		t.Fatalf("publish error: %v", err)
	}

	stats, _ := uc.Stats()
	if stats["pending"] != 2 {
		t.Fatalf("expected 2 pending events, got %d", stats["pending"])
	}

	dispatched, err := uc.DispatchPending(10)
	if err != nil {
		t.Fatalf("dispatch error: %v", err)
	}
	if dispatched != 2 {
		t.Fatalf("expected 2 dispatched events, got %d", dispatched)
	}
	if len(dispatchedEvents) != 2 {
		t.Fatalf("expected 2 dispatched topics, got %d", len(dispatchedEvents))
	}

	statsAfter, _ := uc.Stats()
	if statsAfter["pending"] != 0 || statsAfter["published"] != 2 {
		t.Fatalf("expected 0 pending and 2 published, got: %+v", statsAfter)
	}
}

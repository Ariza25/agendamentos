package calendar_test

import (
	"testing"

	"erpagendamentos/backend/internal/application/calendar"
	calendardomain "erpagendamentos/backend/internal/domain/calendar"
	"erpagendamentos/backend/internal/infrastructure/persistence/memory"
)

func TestCalendarUseCase(t *testing.T) {
	repo := memory.NewCalendarRepository()
	uc := calendar.NewUseCase(repo)

	conn, err := uc.ConnectCalendar(calendardomain.Connection{
		TenantID:           "tenant-1",
		ProfessionalID:     "prof-1",
		Provider:           "google",
		ExternalCalendarID: "primary",
	})
	if err != nil {
		t.Fatalf("connect error: %v", err)
	}
	if conn.ID == "" || conn.Provider != "google" {
		t.Fatalf("unexpected conn: %+v", conn)
	}

	err = uc.SyncAppointment("tenant-1", "prof-1", calendardomain.SyncEvent{
		ID:       "evt-1",
		Title:    "Atendimento Corte",
		StartsAt: "2026-11-01T10:00:00Z",
		EndsAt:   "2026-11-01T10:30:00Z",
	})
	if err != nil {
		t.Fatalf("sync error: %v", err)
	}

	fetched, err := uc.GetConnection("tenant-1", "prof-1")
	if err != nil {
		t.Fatalf("get conn error: %v", err)
	}
	if fetched.LastSyncedAt == "" {
		t.Fatalf("expected last synced at set, got empty")
	}
}


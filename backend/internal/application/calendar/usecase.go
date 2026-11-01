package calendar

import (
	"errors"
	"fmt"
	"time"

	"erpagendamentos/backend/internal/domain/calendar"
)

type UseCase struct {
	repo calendar.Repository
}

func NewUseCase(repo calendar.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) ConnectCalendar(conn calendar.Connection) (calendar.Connection, error) {
	if conn.TenantID == "" || conn.ProfessionalID == "" {
		return calendar.Connection{}, errors.New("tenant id and professional id are required")
	}
	if conn.Provider == "" {
		conn.Provider = "google"
	}
	if conn.Status == "" {
		conn.Status = "active"
	}
	return u.repo.CreateOrUpdate(conn)
}

func (u *UseCase) GetConnection(tenantID string, professionalID string) (calendar.Connection, error) {
	return u.repo.GetByProfessional(tenantID, professionalID)
}

func (u *UseCase) ListConnections(tenantID string) ([]calendar.Connection, error) {
	return u.repo.ListByTenant(tenantID)
}

func (u *UseCase) Disconnect(tenantID string, id string) error {
	return u.repo.Delete(tenantID, id)
}

func (u *UseCase) SyncAppointment(tenantID string, professionalID string, event calendar.SyncEvent) error {
	conn, err := u.repo.GetByProfessional(tenantID, professionalID)
	if err != nil {
		// Not connected, no-op
		return nil
	}
	if conn.Status != "active" {
		return nil
	}

	// Update sync status
	return u.repo.UpdateSyncStatus(tenantID, conn.ID, "active", time.Now().UTC().Format(time.RFC3339), fmt.Sprintf("token_%d", time.Now().Unix()))
}


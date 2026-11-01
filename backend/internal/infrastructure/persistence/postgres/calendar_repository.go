package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"erpagendamentos/backend/internal/domain/calendar"
	"github.com/google/uuid"
)

type CalendarRepository struct {
	db *sql.DB
}

func NewCalendarRepository(db *sql.DB) *CalendarRepository {
	return &CalendarRepository{db: db}
}

func (r *CalendarRepository) CreateOrUpdate(conn calendar.Connection) (calendar.Connection, error) {
	if conn.ID == "" {
		conn.ID = uuid.NewString()
	}
	now := time.Now().UTC()

	err := r.db.QueryRow(`
		INSERT INTO catalog.calendar_connections (id, tenant_id, professional_id, provider, external_calendar_id, sync_token, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (tenant_id, professional_id)
		DO UPDATE SET provider = EXCLUDED.provider, external_calendar_id = EXCLUDED.external_calendar_id, sync_token = EXCLUDED.sync_token, status = EXCLUDED.status, updated_at = EXCLUDED.updated_at
		RETURNING id, created_at, updated_at
	`, conn.ID, conn.TenantID, conn.ProfessionalID, conn.Provider, conn.ExternalCalendarID, conn.SyncToken, conn.Status, now, now).Scan(&conn.ID, &conn.CreatedAt, &conn.UpdatedAt)
	if err != nil {
		// Fallback for cases where unique constraint is on index
		return r.createOrUpdateFallback(conn)
	}
	return conn, nil
}

func (r *CalendarRepository) createOrUpdateFallback(conn calendar.Connection) (calendar.Connection, error) {
	now := time.Now().UTC()
	var existingID string
	err := r.db.QueryRow("SELECT id FROM catalog.calendar_connections WHERE tenant_id = $1 AND professional_id = $2", conn.TenantID, conn.ProfessionalID).Scan(&existingID)
	if err == nil {
		_, err = r.db.Exec(`
			UPDATE catalog.calendar_connections
			SET provider = $1, external_calendar_id = $2, sync_token = $3, status = $4, updated_at = $5
			WHERE id = $6
		`, conn.Provider, conn.ExternalCalendarID, conn.SyncToken, conn.Status, now, existingID)
		if err != nil {
			return calendar.Connection{}, err
		}
		conn.ID = existingID
		conn.UpdatedAt = now.Format(time.RFC3339)
		return conn, nil
	}

	err = r.db.QueryRow(`
		INSERT INTO catalog.calendar_connections (id, tenant_id, professional_id, provider, external_calendar_id, sync_token, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`, conn.ID, conn.TenantID, conn.ProfessionalID, conn.Provider, conn.ExternalCalendarID, conn.SyncToken, conn.Status, now, now).Scan(&conn.ID, &conn.CreatedAt, &conn.UpdatedAt)
	return conn, err
}

func (r *CalendarRepository) GetByProfessional(tenantID string, professionalID string) (calendar.Connection, error) {
	var conn calendar.Connection
	var lastSynced *time.Time
	var created, updated time.Time
	err := r.db.QueryRow(`
		SELECT id, tenant_id, professional_id, provider, external_calendar_id, sync_token, status, last_synced_at, created_at, updated_at
		FROM catalog.calendar_connections
		WHERE tenant_id = $1 AND professional_id = $2
	`, tenantID, professionalID).Scan(&conn.ID, &conn.TenantID, &conn.ProfessionalID, &conn.Provider, &conn.ExternalCalendarID, &conn.SyncToken, &conn.Status, &lastSynced, &created, &updated)
	if err != nil {
		return calendar.Connection{}, errors.New("calendar connection not found")
	}
	conn.CreatedAt = created.Format(time.RFC3339)
	conn.UpdatedAt = updated.Format(time.RFC3339)
	if lastSynced != nil {
		conn.LastSyncedAt = lastSynced.Format(time.RFC3339)
	}
	return conn, nil
}

func (r *CalendarRepository) ListByTenant(tenantID string) ([]calendar.Connection, error) {
	rows, err := r.db.Query(`
		SELECT id, tenant_id, professional_id, provider, external_calendar_id, sync_token, status, last_synced_at, created_at, updated_at
		FROM catalog.calendar_connections
		WHERE tenant_id = $1
	`, tenantID)
	if err != nil {
		return nil, fmt.Errorf("list calendar connections: %w", err)
	}
	defer rows.Close()

	var connections []calendar.Connection
	for rows.Next() {
		var conn calendar.Connection
		var lastSynced *time.Time
		var created, updated time.Time
		if err := rows.Scan(&conn.ID, &conn.TenantID, &conn.ProfessionalID, &conn.Provider, &conn.ExternalCalendarID, &conn.SyncToken, &conn.Status, &lastSynced, &created, &updated); err != nil {
			return nil, err
		}
		conn.CreatedAt = created.Format(time.RFC3339)
		conn.UpdatedAt = updated.Format(time.RFC3339)
		if lastSynced != nil {
			conn.LastSyncedAt = lastSynced.Format(time.RFC3339)
		}
		connections = append(connections, conn)
	}
	return connections, nil
}

func (r *CalendarRepository) Delete(tenantID string, id string) error {
	_, err := r.db.Exec("DELETE FROM catalog.calendar_connections WHERE tenant_id = $1 AND id = $2", tenantID, id)
	return err
}

func (r *CalendarRepository) UpdateSyncStatus(tenantID string, id string, status string, lastSyncedAt string, syncToken string) error {
	now := time.Now().UTC()
	var syncedAt *time.Time
	if lastSyncedAt != "" {
		if t, err := time.Parse(time.RFC3339, lastSyncedAt); err == nil {
			syncedAt = &t
		}
	}
	if syncedAt == nil {
		syncedAt = &now
	}

	_, err := r.db.Exec(`
		UPDATE catalog.calendar_connections
		SET status = $1, last_synced_at = $2, sync_token = $3, updated_at = $4
		WHERE tenant_id = $5 AND id = $6
	`, status, syncedAt, syncToken, now, tenantID, id)
	return err
}


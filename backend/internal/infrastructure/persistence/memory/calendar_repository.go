package memory

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"erpagendamentos/backend/internal/domain/calendar"
)

type CalendarRepository struct {
	mu          sync.RWMutex
	connections map[string]calendar.Connection
}

func NewCalendarRepository() *CalendarRepository {
	return &CalendarRepository{
		connections: make(map[string]calendar.Connection),
	}
}

func (r *CalendarRepository) CreateOrUpdate(conn calendar.Connection) (calendar.Connection, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := conn.TenantID + ":" + conn.ProfessionalID
	now := time.Now().UTC().Format(time.RFC3339)
	if existing, found := r.connections[key]; found {
		conn.ID = existing.ID
		conn.CreatedAt = existing.CreatedAt
		conn.UpdatedAt = now
	} else {
		if conn.ID == "" {
			conn.ID = fmt.Sprintf("calconn_%d", time.Now().UnixNano())
		}
		conn.CreatedAt = now
		conn.UpdatedAt = now
	}
	r.connections[key] = conn
	return conn, nil
}

func (r *CalendarRepository) GetByProfessional(tenantID string, professionalID string) (calendar.Connection, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	conn, found := r.connections[tenantID+":"+professionalID]
	if !found {
		return calendar.Connection{}, errors.New("calendar connection not found")
	}
	return conn, nil
}

func (r *CalendarRepository) ListByTenant(tenantID string) ([]calendar.Connection, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []calendar.Connection
	for _, conn := range r.connections {
		if conn.TenantID == tenantID {
			result = append(result, conn)
		}
	}
	return result, nil
}

func (r *CalendarRepository) Delete(tenantID string, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for k, conn := range r.connections {
		if conn.TenantID == tenantID && conn.ID == id {
			delete(r.connections, k)
			return nil
		}
	}
	return nil
}

func (r *CalendarRepository) UpdateSyncStatus(tenantID string, id string, status string, lastSyncedAt string, syncToken string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for k, conn := range r.connections {
		if conn.TenantID == tenantID && conn.ID == id {
			conn.Status = status
			conn.LastSyncedAt = lastSyncedAt
			conn.SyncToken = syncToken
			conn.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
			r.connections[k] = conn
			return nil
		}
	}
	return errors.New("connection not found")
}


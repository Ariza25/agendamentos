package memory

import (
	"strings"
	"sync"

	auditdomain "erpagendamentos/backend/internal/domain/audit"
)

type AuditRepository struct {
	mu   sync.RWMutex
	logs []auditdomain.AuditLog
}

func NewAuditRepository() *AuditRepository {
	return &AuditRepository{
		logs: make([]auditdomain.AuditLog, 0),
	}
}

func (r *AuditRepository) Record(entry auditdomain.AuditLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Prepend so latest is first
	r.logs = append([]auditdomain.AuditLog{entry}, r.logs...)
	// Cap in-memory logs to 2000
	if len(r.logs) > 2000 {
		r.logs = r.logs[:2000]
	}
	return nil
}

func (r *AuditRepository) List(tenantID string, limit int, offset int) ([]auditdomain.AuditLog, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tenantID = strings.TrimSpace(tenantID)
	var filtered []auditdomain.AuditLog

	for _, l := range r.logs {
		if tenantID == "" || tenantID == "default" || l.TenantID == tenantID {
			filtered = append(filtered, l)
		}
	}

	total := len(filtered)
	if offset >= total {
		return []auditdomain.AuditLog{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return filtered[offset:end], total, nil
}

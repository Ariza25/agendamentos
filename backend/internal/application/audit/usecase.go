package audit

import (
	"time"

	"github.com/google/uuid"
	auditdomain "erpagendamentos/backend/internal/domain/audit"
)

type UseCase struct {
	repo auditdomain.Repository
}

func NewUseCase(repo auditdomain.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) Log(entry auditdomain.AuditLog) error {
	if u == nil || u.repo == nil {
		return nil
	}
	if entry.ID == "" {
		entry.ID = uuid.NewString()
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now()
	}
	return u.repo.Record(entry)
}

func (u *UseCase) List(tenantID string, limit, offset int) ([]auditdomain.AuditLog, int, error) {
	if u == nil || u.repo == nil {
		return []auditdomain.AuditLog{}, 0, nil
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return u.repo.List(tenantID, limit, offset)
}


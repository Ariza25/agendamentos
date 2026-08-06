package memory

import (
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/notification"
)

type NotificationRepository struct {
	mu   sync.Mutex
	logs []notification.Log
}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{}
}

func (r *NotificationRepository) Create(input notification.Log) (notification.Log, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item := input
	item.ID = uuid.NewString()
	item.CreatedAt = time.Now().Format(time.RFC3339)
	if item.SentAt == "" {
		item.SentAt = item.CreatedAt
	}
	r.logs = append([]notification.Log{item}, r.logs...)
	return item, nil
}

func (r *NotificationRepository) ListByTenant(tenantID string, limit int) ([]notification.Log, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if limit <= 0 || limit > 200 {
		limit = 50
	}
	tenantID = strings.TrimSpace(tenantID)
	items := []notification.Log{}
	for _, item := range r.logs {
		if item.TenantID != tenantID {
			continue
		}
		items = append(items, item)
		if len(items) >= limit {
			break
		}
	}
	return items, nil
}

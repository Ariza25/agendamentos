package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/notification"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(input notification.Log) (notification.Log, error) {
	item := input
	item.ID = uuid.NewString()
	item.CreatedAt = time.Now().Format(time.RFC3339)
	if item.SentAt == "" {
		item.SentAt = item.CreatedAt
	}
	if _, err := r.db.Exec(`
		INSERT INTO catalog.notification_logs (id, tenant_id, appointment_id, type, channel, recipient, subject, body, status, error)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, item.ID, item.TenantID, item.AppointmentID, item.Type, item.Channel, item.Recipient, item.Subject, item.Body, item.Status, item.Error); err != nil {
		return notification.Log{}, fmt.Errorf("insert notification log: %w", err)
	}
	return item, nil
}

func (r *NotificationRepository) ListByTenant(tenantID string, limit int) ([]notification.Log, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.db.Query(`
		SELECT id, tenant_id, appointment_id, type, channel, recipient, subject, body, status, error, created_at, sent_at
		FROM catalog.notification_logs
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, strings.TrimSpace(tenantID), limit)
	if err != nil {
		return nil, fmt.Errorf("query notification logs: %w", err)
	}
	defer rows.Close()

	items := []notification.Log{}
	for rows.Next() {
		var item notification.Log
		var createdAt, sentAt time.Time
		if err := rows.Scan(&item.ID, &item.TenantID, &item.AppointmentID, &item.Type, &item.Channel, &item.Recipient, &item.Subject, &item.Body, &item.Status, &item.Error, &createdAt, &sentAt); err != nil {
			return nil, fmt.Errorf("scan notification log: %w", err)
		}
		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.SentAt = sentAt.Format(time.RFC3339)
		items = append(items, item)
	}
	return items, rows.Err()
}

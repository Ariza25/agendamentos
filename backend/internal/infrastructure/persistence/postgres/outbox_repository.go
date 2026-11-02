package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	outboxdomain "erpagendamentos/backend/internal/domain/outbox"
)

type OutboxRepository struct {
	db *sql.DB
}

func NewOutboxRepository(db *sql.DB) *OutboxRepository {
	return &OutboxRepository{db: db}
}

func (r *OutboxRepository) Save(event outboxdomain.OutboxEvent) error {
	payloadJSON, err := json.Marshal(event.Payload)
	if err != nil {
		payloadJSON = []byte("{}")
	}

	query := `
		INSERT INTO catalog.outbox_events (
			id, tenant_id, topic, payload, status, retry_count, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = r.db.Exec(
		query,
		event.ID,
		event.TenantID,
		event.Topic,
		payloadJSON,
		"pending",
		event.RetryCount,
		event.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("save outbox event: %w", err)
	}

	// Trigger PostgreSQL LISTEN/NOTIFY for zero-cost cross-instance reactive notification
	notifyMsg, _ := json.Marshal(map[string]any{
		"id":        event.ID,
		"tenant_id": event.TenantID,
		"topic":     event.Topic,
	})
	_, _ = r.db.Exec(fmt.Sprintf("SELECT pg_notify('erpagendamentos_events', %s)", quoteLiteral(string(notifyMsg))))

	return nil
}

func (r *OutboxRepository) ListPending(limit int) ([]outboxdomain.OutboxEvent, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `
		SELECT id, tenant_id, topic, payload, status, retry_count, error_message, created_at, processed_at
		FROM catalog.outbox_events
		WHERE status = 'pending'
		ORDER BY created_at ASC
		LIMIT $1
	`
	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("list pending outbox events: %w", err)
	}
	defer rows.Close()

	var events []outboxdomain.OutboxEvent
	for rows.Next() {
		var e outboxdomain.OutboxEvent
		var rawPayload []byte
		var processedAt sql.NullTime
		if err := rows.Scan(
			&e.ID,
			&e.TenantID,
			&e.Topic,
			&rawPayload,
			&e.Status,
			&e.RetryCount,
			&e.ErrorMessage,
			&e.CreatedAt,
			&processedAt,
		); err != nil {
			return nil, fmt.Errorf("scan outbox event: %w", err)
		}
		if len(rawPayload) > 0 {
			_ = json.Unmarshal(rawPayload, &e.Payload)
		}
		if processedAt.Valid {
			e.ProcessedAt = &processedAt.Time
		}
		events = append(events, e)
	}
	return events, nil
}

func (r *OutboxRepository) MarkPublished(id string) error {
	now := time.Now()
	_, err := r.db.Exec(`
		UPDATE catalog.outbox_events
		SET status = 'published', processed_at = $1
		WHERE id = $2
	`, now, id)
	return err
}

func (r *OutboxRepository) MarkFailed(id string, errMsg string) error {
	now := time.Now()
	_, err := r.db.Exec(`
		UPDATE catalog.outbox_events
		SET status = 'failed', retry_count = retry_count + 1, error_message = $1, processed_at = $2
		WHERE id = $3
	`, errMsg, now, id)
	return err
}

func (r *OutboxRepository) CountByStatus() (map[string]int, error) {
	rows, err := r.db.Query(`
		SELECT status, COUNT(*)
		FROM catalog.outbox_events
		GROUP BY status
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := map[string]int{
		"pending":   0,
		"published": 0,
		"failed":    0,
	}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err == nil {
			counts[status] = count
		}
	}
	return counts, nil
}

func quoteLiteral(s string) string {
	return "'" + s + "'"
}


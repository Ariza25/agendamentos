package outbox

import "time"

type OutboxEvent struct {
	ID           string         `json:"id"`
	TenantID     string         `json:"tenantId"`
	Topic        string         `json:"topic"`
	Payload      map[string]any `json:"payload"`
	Status       string         `json:"status"` // "pending", "published", "failed"
	RetryCount   int            `json:"retryCount"`
	ErrorMessage string         `json:"errorMessage,omitempty"`
	CreatedAt    time.Time      `json:"createdAt"`
	ProcessedAt  *time.Time     `json:"processedAt,omitempty"`
}

type Repository interface {
	Save(event OutboxEvent) error
	ListPending(limit int) ([]OutboxEvent, error)
	MarkPublished(id string) error
	MarkFailed(id string, errMsg string) error
	CountByStatus() (map[string]int, error)
}


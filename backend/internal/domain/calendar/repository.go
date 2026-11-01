package calendar

type Connection struct {
	ID                 string `json:"id"`
	TenantID           string `json:"tenantId"`
	ProfessionalID     string `json:"professionalId"`
	Provider           string `json:"provider"` // "google", "outlook", "apple"
	ExternalCalendarID string `json:"externalCalendarId"`
	SyncToken          string `json:"syncToken"`
	Status             string `json:"status"` // "active", "error", "disconnected"
	LastSyncedAt       string `json:"lastSyncedAt,omitempty"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}

type SyncEvent struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartsAt    string `json:"startsAt"`
	EndsAt      string `json:"endsAt"`
	ExternalID  string `json:"externalId,omitempty"`
}

type Repository interface {
	CreateOrUpdate(conn Connection) (Connection, error)
	GetByProfessional(tenantID string, professionalID string) (Connection, error)
	ListByTenant(tenantID string) ([]Connection, error)
	Delete(tenantID string, id string) error
	UpdateSyncStatus(tenantID string, id string, status string, lastSyncedAt string, syncToken string) error
}


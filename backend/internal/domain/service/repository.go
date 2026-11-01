package service

type Service struct {
	ID              string   `json:"id"`
	TenantID        string   `json:"tenantId"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	DurationMinutes int      `json:"durationMinutes"`
	PriceCents      int      `json:"priceCents"`
	Color           string   `json:"color"`
	Active          bool     `json:"active"`
	MaxCapacity     int      `json:"maxCapacity,omitempty"`
	ProfessionalIDs []string `json:"professionalIds,omitempty"`
	CreatedAt       string   `json:"createdAt"`
	UpdatedAt       string   `json:"updatedAt"`
}

type Repository interface {
	Create(input Service) (Service, error)
	Update(input Service) (Service, error)
	GetByID(tenantID string, id string) (Service, error)
	ListByTenant(tenantID string) ([]Service, error)
	SetActive(tenantID string, id string, active bool) (Service, error)
	SetProfessionals(tenantID string, serviceID string, professionalIDs []string) error
	GetProfessionalIDs(tenantID string, serviceID string) ([]string, error)
	Delete(tenantID string, id string) error
}

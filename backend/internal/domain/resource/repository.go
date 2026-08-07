package resource

type Resource struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenantId"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Capacity  int    `json:"capacity"`
	Active    bool   `json:"active"`
	Notes     string `json:"notes"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type Repository interface {
	Create(input Resource) (Resource, error)
	Update(input Resource) (Resource, error)
	GetByID(tenantID string, id string) (Resource, error)
	ListByTenant(tenantID string) ([]Resource, error)
	SetActive(tenantID string, id string, active bool) (Resource, error)
	Delete(tenantID string, id string) error
}

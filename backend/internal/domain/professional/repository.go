package professional

type Professional struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenantId"`
	UserID    string `json:"userId"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Color                string  `json:"color"`
	CommissionPercentage float64 `json:"commissionPercentage,omitempty"`
	Active               bool    `json:"active"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type Repository interface {
	Create(input Professional) (Professional, error)
	Update(input Professional) (Professional, error)
	GetByID(tenantID string, id string) (Professional, error)
	ListByTenant(tenantID string) ([]Professional, error)
	SetActive(tenantID string, id string, active bool) (Professional, error)
	Delete(tenantID string, id string) error
}

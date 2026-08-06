package customer

type Customer struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenantId"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Birthday  string `json:"birthday"`
	Notes     string `json:"notes"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type Repository interface {
	Create(input Customer) (Customer, error)
	Update(input Customer) (Customer, error)
	GetByID(tenantID string, id string) (Customer, error)
	FindByPhone(tenantID string, phone string) (Customer, error)
	ListByTenant(tenantID string, search string) ([]Customer, error)
	Delete(tenantID string, id string) error
}

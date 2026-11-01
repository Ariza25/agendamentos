package tenant

type Tenant struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Segment   string `json:"segment,omitempty"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	DeletedAt string `json:"-"`
}

type Repository interface {
	Create(input Tenant) (Tenant, error)
	GetByID(id string) (Tenant, error)
	GetBySlug(slug string) (Tenant, error)
	List() ([]Tenant, error)
	Update(input Tenant) (Tenant, error)
	SetActive(id string, active bool) (Tenant, error)
	Delete(id string) error
}

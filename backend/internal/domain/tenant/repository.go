package tenant

type Tenant struct {
	ID                    string   `json:"id"`
	Name                  string   `json:"name"`
	Slug                  string   `json:"slug"`
	Segment               string   `json:"segment,omitempty"`
	Active                bool     `json:"active"`
	Plan                  string   `json:"plan"`               // "trial", "starter", "pro", "enterprise", "custom"
	SubscriptionStatus    string   `json:"subscriptionStatus"` // "trialing", "active", "past_due", "canceled"
	TrialEndsAt           string   `json:"trialEndsAt,omitempty"`
	SubscriptionExpiresAt string   `json:"subscriptionExpiresAt,omitempty"`
	ActiveModules         []string `json:"activeModules"`
	MaxProfessionals      int      `json:"maxProfessionals"`
	CreatedAt             string   `json:"createdAt"`
	UpdatedAt             string   `json:"updatedAt"`
	DeletedAt             string   `json:"-"`
}

type UpdateSubscriptionInput struct {
	Plan                  string   `json:"plan"`
	SubscriptionStatus    string   `json:"subscriptionStatus"`
	TrialEndsAt           string   `json:"trialEndsAt,omitempty"`
	SubscriptionExpiresAt string   `json:"subscriptionExpiresAt,omitempty"`
	ActiveModules         []string `json:"activeModules"`
	MaxProfessionals      int      `json:"maxProfessionals"`
}

type Repository interface {
	Create(input Tenant) (Tenant, error)
	GetByID(id string) (Tenant, error)
	GetBySlug(slug string) (Tenant, error)
	List() ([]Tenant, error)
	Update(input Tenant) (Tenant, error)
	SetActive(id string, active bool) (Tenant, error)
	UpdateSubscription(id string, input UpdateSubscriptionInput) (Tenant, error)
	Delete(id string) error
}

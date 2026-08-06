package commerce

type Package struct {
	ID            string `json:"id"`
	TenantID      string `json:"tenantId"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	PriceCents    int    `json:"priceCents"`
	SessionsCount int    `json:"sessionsCount"`
	DurationDays  int    `json:"durationDays"`
	Active        bool   `json:"active"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

type PackageCredit struct {
	ID            string `json:"id"`
	TenantID      string `json:"tenantId"`
	PackageID     string `json:"packageId,omitempty"`
	PackageName   string `json:"packageName,omitempty"`
	CustomerID    string `json:"customerId"`
	CustomerName  string `json:"customerName,omitempty"`
	SessionsTotal int    `json:"sessionsTotal"`
	SessionsUsed  int    `json:"sessionsUsed"`
	ExpiresAt     string `json:"expiresAt,omitempty"`
	CreatedAt     string `json:"createdAt"`
}

type GiftCard struct {
	ID           string `json:"id"`
	TenantID     string `json:"tenantId"`
	Code         string `json:"code"`
	ValueCents   int    `json:"valueCents"`
	BalanceCents int    `json:"balanceCents"`
	CustomerID   string `json:"customerId,omitempty"`
	CustomerName string `json:"customerName,omitempty"`
	Status       string `json:"status"`
	ExpiresAt    string `json:"expiresAt,omitempty"`
	CreatedAt    string `json:"createdAt"`
}

type Membership struct {
	ID                string `json:"id"`
	TenantID          string `json:"tenantId"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	PriceCents        int    `json:"priceCents"`
	BillingPeriodDays int    `json:"billingPeriodDays"`
	Active            bool   `json:"active"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}

type MemberPlan struct {
	ID           string `json:"id"`
	TenantID     string `json:"tenantId"`
	MembershipID string `json:"membershipId,omitempty"`
	Membership   string `json:"membershipName,omitempty"`
	CustomerID   string `json:"customerId"`
	CustomerName string `json:"customerName,omitempty"`
	Status       string `json:"status"`
	StartedAt    string `json:"startedAt"`
	ExpiresAt    string `json:"expiresAt,omitempty"`
	CreatedAt    string `json:"createdAt"`
}

type PackageRepository interface {
	Create(input Package) (Package, error)
	Update(input Package) (Package, error)
	GetByID(tenantID string, id string) (Package, error)
	ListByTenant(tenantID string) ([]Package, error)
	SetActive(tenantID string, id string, active bool) (Package, error)
	Delete(tenantID string, id string) error
}

type PackageCreditRepository interface {
	Create(input PackageCredit) (PackageCredit, error)
	GetByID(tenantID string, id string) (PackageCredit, error)
	ListByCustomer(customerID string) ([]PackageCredit, error)
	ListByTenant(tenantID string) ([]PackageCredit, error)
	IncrementUsed(tenantID string, id string) (PackageCredit, error)
}

type GiftCardRepository interface {
	Create(input GiftCard) (GiftCard, error)
	GetByID(tenantID string, id string) (GiftCard, error)
	GetByCode(tenantID string, code string) (GiftCard, error)
	ListByTenant(tenantID string) ([]GiftCard, error)
	Update(input GiftCard) (GiftCard, error)
	Delete(tenantID string, id string) error
}

type MembershipRepository interface {
	Create(input Membership) (Membership, error)
	Update(input Membership) (Membership, error)
	GetByID(tenantID string, id string) (Membership, error)
	ListByTenant(tenantID string) ([]Membership, error)
	SetActive(tenantID string, id string, active bool) (Membership, error)
	Delete(tenantID string, id string) error
}

type MemberPlanRepository interface {
	Create(input MemberPlan) (MemberPlan, error)
	GetByID(tenantID string, id string) (MemberPlan, error)
	ListByCustomer(customerID string) ([]MemberPlan, error)
	ListByTenant(tenantID string) ([]MemberPlan, error)
	UpdateStatus(tenantID string, id string, status string) (MemberPlan, error)
}

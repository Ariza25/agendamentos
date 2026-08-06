package memory

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/commerce"
)

type PackageRepository struct {
	mu       sync.Mutex
	packages map[string]commerce.Package
}

func NewPackageRepository() *PackageRepository {
	return &PackageRepository{
		packages: make(map[string]commerce.Package),
	}
}

func (r *PackageRepository) Create(input commerce.Package) (commerce.Package, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item := input
	item.ID = uuid.NewString()
	item.CreatedAt = time.Now().Format(time.RFC3339)
	item.UpdatedAt = item.CreatedAt
	r.packages[item.ID] = item
	return item, nil
}

func (r *PackageRepository) Update(input commerce.Package) (commerce.Package, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.packages[strings.TrimSpace(input.ID)]
	if !exists || current.TenantID != strings.TrimSpace(input.TenantID) {
		return commerce.Package{}, fmt.Errorf("package not found")
	}
	current.Name = input.Name
	current.Description = input.Description
	current.PriceCents = input.PriceCents
	current.SessionsCount = input.SessionsCount
	current.DurationDays = input.DurationDays
	current.Active = input.Active
	current.UpdatedAt = time.Now().Format(time.RFC3339)
	r.packages[current.ID] = current
	return current, nil
}

func (r *PackageRepository) GetByID(tenantID string, id string) (commerce.Package, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.packages[strings.TrimSpace(id)]
	if !exists || item.TenantID != strings.TrimSpace(tenantID) {
		return commerce.Package{}, fmt.Errorf("package not found")
	}
	return item, nil
}

func (r *PackageRepository) ListByTenant(tenantID string) ([]commerce.Package, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID = strings.TrimSpace(tenantID)
	items := []commerce.Package{}
	for _, item := range r.packages {
		if item.TenantID == tenantID {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
	})
	return items, nil
}

func (r *PackageRepository) SetActive(tenantID string, id string, active bool) (commerce.Package, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.packages[strings.TrimSpace(id)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return commerce.Package{}, fmt.Errorf("package not found")
	}
	current.Active = active
	current.UpdatedAt = time.Now().Format(time.RFC3339)
	r.packages[current.ID] = current
	return current, nil
}

func (r *PackageRepository) Delete(tenantID string, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.packages[strings.TrimSpace(id)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return fmt.Errorf("package not found")
	}
	delete(r.packages, current.ID)
	return nil
}

type PackageCreditRepository struct {
	mu      sync.Mutex
	credits map[string]commerce.PackageCredit
}

func NewPackageCreditRepository() *PackageCreditRepository {
	return &PackageCreditRepository{
		credits: make(map[string]commerce.PackageCredit),
	}
}

func (r *PackageCreditRepository) Create(input commerce.PackageCredit) (commerce.PackageCredit, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item := input
	item.ID = uuid.NewString()
	item.CreatedAt = time.Now().Format(time.RFC3339)
	r.credits[item.ID] = item
	return item, nil
}

func (r *PackageCreditRepository) GetByID(tenantID string, id string) (commerce.PackageCredit, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.credits[strings.TrimSpace(id)]
	if !exists || item.TenantID != strings.TrimSpace(tenantID) {
		return commerce.PackageCredit{}, fmt.Errorf("package credit not found")
	}
	return item, nil
}

func (r *PackageCreditRepository) ListByCustomer(customerID string) ([]commerce.PackageCredit, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	customerID = strings.TrimSpace(customerID)
	items := []commerce.PackageCredit{}
	for _, item := range r.credits {
		if item.CustomerID == customerID {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt > items[j].CreatedAt
	})
	return items, nil
}

func (r *PackageCreditRepository) ListByTenant(tenantID string) ([]commerce.PackageCredit, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID = strings.TrimSpace(tenantID)
	items := []commerce.PackageCredit{}
	for _, item := range r.credits {
		if item.TenantID == tenantID {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt > items[j].CreatedAt
	})
	return items, nil
}

func (r *PackageCreditRepository) IncrementUsed(tenantID string, id string) (commerce.PackageCredit, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.credits[strings.TrimSpace(id)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return commerce.PackageCredit{}, fmt.Errorf("package credit not found")
	}
	if current.SessionsUsed >= current.SessionsTotal {
		return commerce.PackageCredit{}, fmt.Errorf("package credit exhausted")
	}
	current.SessionsUsed++
	r.credits[current.ID] = current
	return current, nil
}

type GiftCardRepository struct {
	mu        sync.Mutex
	giftCards map[string]commerce.GiftCard
}

func NewGiftCardRepository() *GiftCardRepository {
	return &GiftCardRepository{
		giftCards: make(map[string]commerce.GiftCard),
	}
}

func (r *GiftCardRepository) Create(input commerce.GiftCard) (commerce.GiftCard, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item := input
	item.ID = uuid.NewString()
	item.CreatedAt = time.Now().Format(time.RFC3339)
	if item.Status == "" {
		item.Status = "active"
	}
	if item.BalanceCents == 0 {
		item.BalanceCents = item.ValueCents
	}
	r.giftCards[item.ID] = item
	return item, nil
}

func (r *GiftCardRepository) GetByID(tenantID string, id string) (commerce.GiftCard, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.giftCards[strings.TrimSpace(id)]
	if !exists || item.TenantID != strings.TrimSpace(tenantID) {
		return commerce.GiftCard{}, fmt.Errorf("gift card not found")
	}
	return item, nil
}

func (r *GiftCardRepository) GetByCode(tenantID string, code string) (commerce.GiftCard, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	code = strings.ToUpper(strings.TrimSpace(code))
	tenantID = strings.TrimSpace(tenantID)
	for _, item := range r.giftCards {
		if item.TenantID == tenantID && strings.ToUpper(item.Code) == code {
			return item, nil
		}
	}
	return commerce.GiftCard{}, fmt.Errorf("gift card not found")
}

func (r *GiftCardRepository) ListByTenant(tenantID string) ([]commerce.GiftCard, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID = strings.TrimSpace(tenantID)
	items := []commerce.GiftCard{}
	for _, item := range r.giftCards {
		if item.TenantID == tenantID {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt > items[j].CreatedAt
	})
	return items, nil
}

func (r *GiftCardRepository) Update(input commerce.GiftCard) (commerce.GiftCard, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.giftCards[strings.TrimSpace(input.ID)]
	if !exists || current.TenantID != strings.TrimSpace(input.TenantID) {
		return commerce.GiftCard{}, fmt.Errorf("gift card not found")
	}
	current.BalanceCents = input.BalanceCents
	current.CustomerID = input.CustomerID
	current.Status = input.Status
	r.giftCards[current.ID] = current
	return current, nil
}

func (r *GiftCardRepository) Delete(tenantID string, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.giftCards[strings.TrimSpace(id)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return fmt.Errorf("gift card not found")
	}
	delete(r.giftCards, current.ID)
	return nil
}

type MembershipRepository struct {
	mu          sync.Mutex
	memberships map[string]commerce.Membership
}

func NewMembershipRepository() *MembershipRepository {
	return &MembershipRepository{
		memberships: make(map[string]commerce.Membership),
	}
}

func (r *MembershipRepository) Create(input commerce.Membership) (commerce.Membership, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item := input
	item.ID = uuid.NewString()
	item.CreatedAt = time.Now().Format(time.RFC3339)
	item.UpdatedAt = item.CreatedAt
	if item.BillingPeriodDays <= 0 {
		item.BillingPeriodDays = 30
	}
	r.memberships[item.ID] = item
	return item, nil
}

func (r *MembershipRepository) Update(input commerce.Membership) (commerce.Membership, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.memberships[strings.TrimSpace(input.ID)]
	if !exists || current.TenantID != strings.TrimSpace(input.TenantID) {
		return commerce.Membership{}, fmt.Errorf("membership not found")
	}
	current.Name = input.Name
	current.Description = input.Description
	current.PriceCents = input.PriceCents
	if input.BillingPeriodDays > 0 {
		current.BillingPeriodDays = input.BillingPeriodDays
	}
	current.Active = input.Active
	current.UpdatedAt = time.Now().Format(time.RFC3339)
	r.memberships[current.ID] = current
	return current, nil
}

func (r *MembershipRepository) GetByID(tenantID string, id string) (commerce.Membership, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.memberships[strings.TrimSpace(id)]
	if !exists || item.TenantID != strings.TrimSpace(tenantID) {
		return commerce.Membership{}, fmt.Errorf("membership not found")
	}
	return item, nil
}

func (r *MembershipRepository) ListByTenant(tenantID string) ([]commerce.Membership, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID = strings.TrimSpace(tenantID)
	items := []commerce.Membership{}
	for _, item := range r.memberships {
		if item.TenantID == tenantID {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
	})
	return items, nil
}

func (r *MembershipRepository) SetActive(tenantID string, id string, active bool) (commerce.Membership, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.memberships[strings.TrimSpace(id)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return commerce.Membership{}, fmt.Errorf("membership not found")
	}
	current.Active = active
	current.UpdatedAt = time.Now().Format(time.RFC3339)
	r.memberships[current.ID] = current
	return current, nil
}

func (r *MembershipRepository) Delete(tenantID string, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.memberships[strings.TrimSpace(id)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return fmt.Errorf("membership not found")
	}
	delete(r.memberships, current.ID)
	return nil
}

type MemberPlanRepository struct {
	mu    sync.Mutex
	plans map[string]commerce.MemberPlan
}

func NewMemberPlanRepository() *MemberPlanRepository {
	return &MemberPlanRepository{
		plans: make(map[string]commerce.MemberPlan),
	}
}

func (r *MemberPlanRepository) Create(input commerce.MemberPlan) (commerce.MemberPlan, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item := input
	item.ID = uuid.NewString()
	item.CreatedAt = time.Now().Format(time.RFC3339)
	if item.Status == "" {
		item.Status = "active"
	}
	if item.StartedAt == "" {
		item.StartedAt = item.CreatedAt
	}
	r.plans[item.ID] = item
	return item, nil
}

func (r *MemberPlanRepository) GetByID(tenantID string, id string) (commerce.MemberPlan, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.plans[strings.TrimSpace(id)]
	if !exists || item.TenantID != strings.TrimSpace(tenantID) {
		return commerce.MemberPlan{}, fmt.Errorf("member plan not found")
	}
	return item, nil
}

func (r *MemberPlanRepository) ListByCustomer(customerID string) ([]commerce.MemberPlan, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	customerID = strings.TrimSpace(customerID)
	items := []commerce.MemberPlan{}
	for _, item := range r.plans {
		if item.CustomerID == customerID {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt > items[j].CreatedAt
	})
	return items, nil
}

func (r *MemberPlanRepository) ListByTenant(tenantID string) ([]commerce.MemberPlan, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID = strings.TrimSpace(tenantID)
	items := []commerce.MemberPlan{}
	for _, item := range r.plans {
		if item.TenantID == tenantID {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt > items[j].CreatedAt
	})
	return items, nil
}

func (r *MemberPlanRepository) UpdateStatus(tenantID string, id string, status string) (commerce.MemberPlan, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.plans[strings.TrimSpace(id)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return commerce.MemberPlan{}, fmt.Errorf("member plan not found")
	}
	current.Status = status
	r.plans[current.ID] = current
	return current, nil
}

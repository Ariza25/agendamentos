package memory

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"erpagendamentos/backend/internal/domain/tenant"
)

type TenantRepository struct {
	mu      sync.Mutex
	tenants map[string]tenant.Tenant
	bySlug  map[string]string
}

func NewTenantRepository() *TenantRepository {
	return &TenantRepository{
		tenants: make(map[string]tenant.Tenant),
		bySlug:  make(map[string]string),
	}
}

func defaultModules() []string {
	return []string{
		"agenda", "appointments", "customers", "services", "professionals",
		"finance", "commissions", "inventory", "pos_cashregister", "notifications", "reviews",
	}
}

func (r *TenantRepository) Create(input tenant.Tenant) (tenant.Tenant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	slug := strings.ToLower(strings.TrimSpace(input.Slug))
	if slug == "" {
		slug = strings.ToLower(strings.TrimSpace(input.ID))
	}
	if existingID, exists := r.bySlug[slug]; exists && existingID != input.ID {
		return tenant.Tenant{}, fmt.Errorf("slug already used")
	}
	item := input
	item.ID = strings.TrimSpace(input.ID)
	if item.ID == "" {
		return tenant.Tenant{}, fmt.Errorf("tenant id is required")
	}
	item.Slug = slug
	if item.Segment == "" {
		item.Segment = "general"
	}
	if item.Plan == "" {
		item.Plan = "trial"
	}
	if item.SubscriptionStatus == "" {
		item.SubscriptionStatus = "trialing"
	}
	if len(item.ActiveModules) == 0 {
		item.ActiveModules = defaultModules()
	}
	if item.MaxProfessionals <= 0 {
		item.MaxProfessionals = 5
	}
	item.CreatedAt = time.Now().Format(time.RFC3339)
	item.UpdatedAt = item.CreatedAt
	r.tenants[item.ID] = item
	r.bySlug[slug] = item.ID
	return item, nil
}

func (r *TenantRepository) GetByID(id string) (tenant.Tenant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.tenants[strings.TrimSpace(id)]
	if !exists {
		return tenant.Tenant{}, fmt.Errorf("tenant not found")
	}
	return item, nil
}

func (r *TenantRepository) GetBySlug(slug string) (tenant.Tenant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id, exists := r.bySlug[strings.ToLower(strings.TrimSpace(slug))]
	if !exists {
		return tenant.Tenant{}, fmt.Errorf("tenant not found")
	}
	return r.tenants[id], nil
}

func (r *TenantRepository) List() ([]tenant.Tenant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	items := make([]tenant.Tenant, 0, len(r.tenants))
	for _, item := range r.tenants {
		items = append(items, item)
	}
	return items, nil
}

func (r *TenantRepository) Update(input tenant.Tenant) (tenant.Tenant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.tenants[input.ID]
	if !exists {
		return tenant.Tenant{}, fmt.Errorf("tenant not found")
	}
	slug := strings.ToLower(strings.TrimSpace(input.Slug))
	if slug == "" {
		slug = current.Slug
	}
	if existingID, ok := r.bySlug[slug]; ok && existingID != input.ID {
		return tenant.Tenant{}, fmt.Errorf("slug already used")
	}
	delete(r.bySlug, strings.ToLower(strings.TrimSpace(current.Slug)))
	current.Name = strings.TrimSpace(input.Name)
	current.Slug = slug
	if strings.TrimSpace(input.Segment) != "" {
		current.Segment = strings.TrimSpace(input.Segment)
	}
	current.Active = input.Active
	current.UpdatedAt = time.Now().Format(time.RFC3339)
	r.tenants[current.ID] = current
	r.bySlug[slug] = current.ID
	return current, nil
}

func (r *TenantRepository) SetActive(id string, active bool) (tenant.Tenant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.tenants[id]
	if !exists {
		return tenant.Tenant{}, fmt.Errorf("tenant not found")
	}
	current.Active = active
	current.UpdatedAt = time.Now().Format(time.RFC3339)
	r.tenants[id] = current
	return current, nil
}

func (r *TenantRepository) UpdateSubscription(id string, input tenant.UpdateSubscriptionInput) (tenant.Tenant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.tenants[id]
	if !exists {
		return tenant.Tenant{}, fmt.Errorf("tenant not found")
	}
	current.Plan = input.Plan
	current.SubscriptionStatus = input.SubscriptionStatus
	current.TrialEndsAt = input.TrialEndsAt
	current.SubscriptionExpiresAt = input.SubscriptionExpiresAt
	if len(input.ActiveModules) > 0 {
		current.ActiveModules = input.ActiveModules
	}
	if input.MaxProfessionals > 0 {
		current.MaxProfessionals = input.MaxProfessionals
	}
	current.UpdatedAt = time.Now().Format(time.RFC3339)
	r.tenants[id] = current
	return current, nil
}

func (r *TenantRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.tenants[id]
	if !exists {
		return fmt.Errorf("tenant not found")
	}
	delete(r.tenants, id)
	delete(r.bySlug, strings.ToLower(strings.TrimSpace(current.Slug)))
	return nil
}

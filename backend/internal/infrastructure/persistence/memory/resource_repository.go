package memory

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/resource"
)

type ResourceRepository struct {
	mu        sync.Mutex
	resources map[string]resource.Resource
}

func NewResourceRepository() *ResourceRepository {
	return &ResourceRepository{
		resources: make(map[string]resource.Resource),
	}
}

func (r *ResourceRepository) Create(input resource.Resource) (resource.Resource, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item := input
	item.ID = uuid.NewString()
	item.Name = strings.TrimSpace(item.Name)
	item.Type = strings.TrimSpace(item.Type)
	item.Notes = strings.TrimSpace(item.Notes)
	item.CreatedAt = time.Now().Format(time.RFC3339)
	item.UpdatedAt = item.CreatedAt
	r.resources[item.ID] = item
	return item, nil
}

func (r *ResourceRepository) Update(input resource.Resource) (resource.Resource, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.resources[input.ID]
	if !exists {
		return resource.Resource{}, fmt.Errorf("resource not found")
	}
	if strings.TrimSpace(input.Name) != "" {
		current.Name = input.Name
	}
	if strings.TrimSpace(input.Type) != "" {
		current.Type = input.Type
	}
	if input.Capacity > 0 {
		current.Capacity = input.Capacity
	}
	current.Notes = strings.TrimSpace(input.Notes)
	current.Active = input.Active
	current.UpdatedAt = time.Now().Format(time.RFC3339)
	r.resources[current.ID] = current
	return current, nil
}

func (r *ResourceRepository) GetByID(tenantID string, id string) (resource.Resource, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.resources[strings.TrimSpace(id)]
	if !exists || item.TenantID != strings.TrimSpace(tenantID) {
		return resource.Resource{}, fmt.Errorf("resource not found")
	}
	return item, nil
}

func (r *ResourceRepository) ListByTenant(tenantID string) ([]resource.Resource, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID = strings.TrimSpace(tenantID)
	items := []resource.Resource{}
	for _, item := range r.resources {
		if item.TenantID == tenantID {
			items = append(items, item)
		}
	}
	return items, nil
}

func (r *ResourceRepository) SetActive(tenantID string, id string, active bool) (resource.Resource, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.resources[strings.TrimSpace(id)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return resource.Resource{}, fmt.Errorf("resource not found")
	}
	current.Active = active
	current.UpdatedAt = time.Now().Format(time.RFC3339)
	r.resources[current.ID] = current
	return current, nil
}

func (r *ResourceRepository) Delete(tenantID string, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.resources[strings.TrimSpace(id)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return fmt.Errorf("resource not found")
	}
	delete(r.resources, current.ID)
	return nil
}

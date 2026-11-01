package memory

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/service"
)

type ServiceRepository struct {
	mu       sync.Mutex
	services map[string]service.Service
}

func NewServiceRepository() *ServiceRepository {
	return &ServiceRepository{
		services: make(map[string]service.Service),
	}
}

func (r *ServiceRepository) Create(input service.Service) (service.Service, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item := input
	item.ID = uuid.NewString()
	item.Name = strings.TrimSpace(item.Name)
	item.Description = strings.TrimSpace(item.Description)
	if item.MaxCapacity <= 0 {
		item.MaxCapacity = 1
	}
	item.CreatedAt = time.Now().Format(time.RFC3339)
	item.UpdatedAt = item.CreatedAt
	r.services[item.ID] = item
	return item, nil
}

func (r *ServiceRepository) Update(input service.Service) (service.Service, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.services[input.ID]
	if !exists {
		return service.Service{}, fmt.Errorf("service not found")
	}
	if strings.TrimSpace(input.Name) != "" {
		current.Name = input.Name
	}
	current.Description = strings.TrimSpace(input.Description)
	if input.DurationMinutes > 0 {
		current.DurationMinutes = input.DurationMinutes
	}
	current.PriceCents = input.PriceCents
	if strings.TrimSpace(input.Color) != "" {
		current.Color = input.Color
	}
	if input.MaxCapacity > 0 {
		current.MaxCapacity = input.MaxCapacity
	}
	current.Active = input.Active
	current.UpdatedAt = time.Now().Format(time.RFC3339)
	r.services[current.ID] = current
	return current, nil
}

func (r *ServiceRepository) GetByID(tenantID string, id string) (service.Service, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.services[strings.TrimSpace(id)]
	if !exists || item.TenantID != strings.TrimSpace(tenantID) {
		return service.Service{}, fmt.Errorf("service not found")
	}
	return item, nil
}

func (r *ServiceRepository) ListByTenant(tenantID string) ([]service.Service, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID = strings.TrimSpace(tenantID)
	items := []service.Service{}
	for _, item := range r.services {
		if item.TenantID == tenantID {
			items = append(items, item)
		}
	}
	return items, nil
}

func (r *ServiceRepository) SetActive(tenantID string, id string, active bool) (service.Service, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.services[strings.TrimSpace(id)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return service.Service{}, fmt.Errorf("service not found")
	}
	current.Active = active
	current.UpdatedAt = time.Now().Format(time.RFC3339)
	r.services[current.ID] = current
	return current, nil
}

func (r *ServiceRepository) SetProfessionals(tenantID string, serviceID string, professionalIDs []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.services[strings.TrimSpace(serviceID)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return fmt.Errorf("service not found")
	}
	ids := []string{}
	for _, id := range professionalIDs {
		id = strings.TrimSpace(id)
		if id != "" {
			ids = append(ids, id)
		}
	}
	current.ProfessionalIDs = ids
	r.services[current.ID] = current
	return nil
}

func (r *ServiceRepository) GetProfessionalIDs(tenantID string, serviceID string) ([]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.services[strings.TrimSpace(serviceID)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return nil, fmt.Errorf("service not found")
	}
	return current.ProfessionalIDs, nil
}

func (r *ServiceRepository) Delete(tenantID string, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.services[strings.TrimSpace(id)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return fmt.Errorf("service not found")
	}
	delete(r.services, current.ID)
	return nil
}

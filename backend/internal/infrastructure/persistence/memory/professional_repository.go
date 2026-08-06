package memory

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/professional"
)

type ProfessionalRepository struct {
	mu            sync.Mutex
	professionals map[string]professional.Professional
}

func NewProfessionalRepository() *ProfessionalRepository {
	return &ProfessionalRepository{
		professionals: make(map[string]professional.Professional),
	}
}

func (r *ProfessionalRepository) Create(input professional.Professional) (professional.Professional, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item := input
	item.ID = uuid.NewString()
	item.Name = strings.TrimSpace(item.Name)
	item.Email = strings.ToLower(strings.TrimSpace(item.Email))
	item.CreatedAt = time.Now().Format(time.RFC3339)
	item.UpdatedAt = item.CreatedAt
	r.professionals[item.ID] = item
	return item, nil
}

func (r *ProfessionalRepository) Update(input professional.Professional) (professional.Professional, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.professionals[input.ID]
	if !exists {
		return professional.Professional{}, fmt.Errorf("professional not found")
	}
	if strings.TrimSpace(input.Name) != "" {
		current.Name = input.Name
	}
	if strings.TrimSpace(input.Email) != "" {
		current.Email = strings.ToLower(strings.TrimSpace(input.Email))
	}
	current.Phone = strings.TrimSpace(input.Phone)
	if strings.TrimSpace(input.Color) != "" {
		current.Color = input.Color
	}
	if strings.TrimSpace(input.UserID) != "" {
		current.UserID = input.UserID
	}
	current.Active = input.Active
	current.UpdatedAt = time.Now().Format(time.RFC3339)
	r.professionals[current.ID] = current
	return current, nil
}

func (r *ProfessionalRepository) GetByID(tenantID string, id string) (professional.Professional, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.professionals[strings.TrimSpace(id)]
	if !exists || item.TenantID != strings.TrimSpace(tenantID) {
		return professional.Professional{}, fmt.Errorf("professional not found")
	}
	return item, nil
}

func (r *ProfessionalRepository) ListByTenant(tenantID string) ([]professional.Professional, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID = strings.TrimSpace(tenantID)
	items := []professional.Professional{}
	for _, item := range r.professionals {
		if item.TenantID == tenantID {
			items = append(items, item)
		}
	}
	return items, nil
}

func (r *ProfessionalRepository) SetActive(tenantID string, id string, active bool) (professional.Professional, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.professionals[strings.TrimSpace(id)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return professional.Professional{}, fmt.Errorf("professional not found")
	}
	current.Active = active
	current.UpdatedAt = time.Now().Format(time.RFC3339)
	r.professionals[current.ID] = current
	return current, nil
}

func (r *ProfessionalRepository) Delete(tenantID string, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.professionals[strings.TrimSpace(id)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return fmt.Errorf("professional not found")
	}
	delete(r.professionals, current.ID)
	return nil
}

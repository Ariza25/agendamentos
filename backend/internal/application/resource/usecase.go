package resource

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	domain "erpagendamentos/backend/internal/domain/resource"
)

var (
	ErrResourceNotFound = errors.New("resource not found")
	ErrInvalidResource  = errors.New("invalid resource")
)

type UseCase struct {
	repository domain.Repository
}

func NewUseCase(repository domain.Repository) *UseCase {
	return &UseCase{repository: repository}
}

func (u *UseCase) ListByTenant(tenantID string) ([]domain.Resource, error) {
	return u.repository.ListByTenant(strings.TrimSpace(tenantID))
}

func (u *UseCase) GetByID(tenantID string, id string) (domain.Resource, error) {
	item, err := u.repository.GetByID(strings.TrimSpace(tenantID), strings.TrimSpace(id))
	if err != nil {
		return domain.Resource{}, ErrResourceNotFound
	}
	return item, nil
}

func (u *UseCase) Create(input domain.Resource) (domain.Resource, error) {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.Name = strings.TrimSpace(input.Name)
	input.Type = strings.TrimSpace(input.Type)
	input.Notes = strings.TrimSpace(input.Notes)
	if input.TenantID == "" || input.Name == "" {
		return domain.Resource{}, ErrInvalidResource
	}
	if input.Type == "" {
		input.Type = "room"
	}
	if input.Capacity <= 0 {
		input.Capacity = 1
	}
	input.ID = uuid.NewString()
	created, err := u.repository.Create(input)
	if err != nil {
		return domain.Resource{}, err
	}
	return created, nil
}

func (u *UseCase) Update(input domain.Resource) (domain.Resource, error) {
	tenantID := strings.TrimSpace(input.TenantID)
	id := strings.TrimSpace(input.ID)
	current, err := u.repository.GetByID(tenantID, id)
	if err != nil {
		return domain.Resource{}, ErrResourceNotFound
	}
	input.TenantID = tenantID
	input.CreatedAt = current.CreatedAt
	if strings.TrimSpace(input.Name) == "" {
		input.Name = current.Name
	}
	if strings.TrimSpace(input.Type) == "" {
		input.Type = current.Type
	}
	if input.Capacity <= 0 {
		input.Capacity = current.Capacity
	}
	input.Notes = strings.TrimSpace(input.Notes)
	updated, err := u.repository.Update(input)
	if err != nil {
		return domain.Resource{}, err
	}
	return updated, nil
}

func (u *UseCase) SetActive(tenantID string, id string, active bool) (domain.Resource, error) {
	item, err := u.repository.SetActive(strings.TrimSpace(tenantID), strings.TrimSpace(id), active)
	if err != nil {
		return domain.Resource{}, ErrResourceNotFound
	}
	return item, nil
}

func (u *UseCase) Delete(tenantID string, id string) error {
	if err := u.repository.Delete(strings.TrimSpace(tenantID), strings.TrimSpace(id)); err != nil {
		return ErrResourceNotFound
	}
	return nil
}

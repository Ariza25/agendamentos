package service

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	domain "erpagendamentos/backend/internal/domain/service"
)

var (
	ErrServiceNotFound = errors.New("service not found")
	ErrInvalidService  = errors.New("invalid service")
)

type UseCase struct {
	repository domain.Repository
}

func NewUseCase(repository domain.Repository) *UseCase {
	return &UseCase{repository: repository}
}

func (u *UseCase) ListByTenant(tenantID string) ([]domain.Service, error) {
	return u.repository.ListByTenant(strings.TrimSpace(tenantID))
}

func (u *UseCase) GetByID(tenantID string, id string) (domain.Service, error) {
	item, err := u.repository.GetByID(strings.TrimSpace(tenantID), strings.TrimSpace(id))
	if err != nil {
		return domain.Service{}, ErrServiceNotFound
	}
	return item, nil
}

func (u *UseCase) Create(input domain.Service) (domain.Service, error) {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)
	input.Color = strings.TrimSpace(input.Color)
	if input.TenantID == "" || input.Name == "" {
		return domain.Service{}, ErrInvalidService
	}
	if input.DurationMinutes <= 0 {
		input.DurationMinutes = 30
	}
	if input.PriceCents < 0 {
		input.PriceCents = 0
	}
	input.ID = uuid.NewString()
	if input.Color == "" {
		input.Color = "#8B5CF6"
	}
	return u.repository.Create(input)
}

func (u *UseCase) Update(input domain.Service) (domain.Service, error) {
	tenantID := strings.TrimSpace(input.TenantID)
	id := strings.TrimSpace(input.ID)
	current, err := u.repository.GetByID(tenantID, id)
	if err != nil {
		return domain.Service{}, ErrServiceNotFound
	}
	input.TenantID = tenantID
	input.CreatedAt = current.CreatedAt
	if strings.TrimSpace(input.Name) == "" {
		input.Name = current.Name
	}
	if strings.TrimSpace(input.Description) == "" {
		input.Description = current.Description
	}
	if strings.TrimSpace(input.Color) == "" {
		input.Color = current.Color
	}
	if input.DurationMinutes <= 0 {
		input.DurationMinutes = current.DurationMinutes
	}
	if input.PriceCents < 0 {
		input.PriceCents = current.PriceCents
	}
	return u.repository.Update(input)
}

func (u *UseCase) SetActive(tenantID string, id string, active bool) (domain.Service, error) {
	item, err := u.repository.SetActive(strings.TrimSpace(tenantID), strings.TrimSpace(id), active)
	if err != nil {
		return domain.Service{}, ErrServiceNotFound
	}
	return item, nil
}

func (u *UseCase) Delete(tenantID string, id string) error {
	if err := u.repository.Delete(strings.TrimSpace(tenantID), strings.TrimSpace(id)); err != nil {
		return ErrServiceNotFound
	}
	return nil
}

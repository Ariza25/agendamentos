package tenant

import (
	"errors"
	"regexp"
	"strings"

	"github.com/google/uuid"

	domain "erpagendamentos/backend/internal/domain/tenant"
)

var (
	ErrTenantNotFound   = errors.New("tenant not found")
	ErrInvalidTenant    = errors.New("invalid tenant")
	ErrSlugAlreadyUsed  = errors.New("tenant slug already used")
	ErrLastActiveTenant = errors.New("cannot delete the last active tenant")
)

var slugPattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,60}[a-z0-9])?$`)

type UseCase struct {
	repository domain.Repository
}

func NewUseCase(repository domain.Repository) *UseCase {
	return &UseCase{repository: repository}
}

func (u *UseCase) ListTenants() ([]domain.Tenant, error) {
	return u.repository.List()
}

func (u *UseCase) GetTenant(id string) (domain.Tenant, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.Tenant{}, ErrTenantNotFound
	}
	tenant, err := u.repository.GetByID(id)
	if err != nil {
		return domain.Tenant{}, ErrTenantNotFound
	}
	return tenant, nil
}

func (u *UseCase) CreateTenant(input domain.Tenant) (domain.Tenant, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Slug = strings.TrimSpace(strings.ToLower(input.Slug))

	if input.Name == "" {
		return domain.Tenant{}, ErrInvalidTenant
	}
	if input.Slug == "" {
		return domain.Tenant{}, ErrInvalidTenant
	}
	if !slugPattern.MatchString(input.Slug) {
		return domain.Tenant{}, ErrInvalidTenant
	}
	if _, err := u.repository.GetBySlug(input.Slug); err == nil {
		return domain.Tenant{}, ErrSlugAlreadyUsed
	}

	input.ID = uuid.NewString()
	input.Active = true
	return u.repository.Create(input)
}

func (u *UseCase) UpdateTenant(input domain.Tenant) (domain.Tenant, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Slug = strings.TrimSpace(strings.ToLower(input.Slug))
	if strings.TrimSpace(input.ID) == "" {
		return domain.Tenant{}, ErrTenantNotFound
	}
	if input.Name == "" {
		return domain.Tenant{}, ErrInvalidTenant
	}
	if input.Slug != "" && !slugPattern.MatchString(input.Slug) {
		return domain.Tenant{}, ErrInvalidTenant
	}

	current, err := u.repository.GetByID(input.ID)
	if err != nil {
		return domain.Tenant{}, ErrTenantNotFound
	}
	if input.Slug == "" {
		input.Slug = current.Slug
	}
	if existing, err := u.repository.GetBySlug(input.Slug); err == nil && existing.ID != input.ID {
		return domain.Tenant{}, ErrSlugAlreadyUsed
	}

	input.Active = current.Active
	return u.repository.Update(input)
}

func (u *UseCase) SetActive(id string, active bool) (domain.Tenant, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.Tenant{}, ErrTenantNotFound
	}
	if _, err := u.repository.GetByID(id); err != nil {
		return domain.Tenant{}, ErrTenantNotFound
	}
	return u.repository.SetActive(id, active)
}

func (u *UseCase) DeleteTenant(id string) error {
	id = strings.TrimSpace(id)
	if id == "" || id == "default" {
		return ErrInvalidTenant
	}
	if _, err := u.repository.GetByID(id); err != nil {
		return ErrTenantNotFound
	}
	return u.repository.Delete(id)
}

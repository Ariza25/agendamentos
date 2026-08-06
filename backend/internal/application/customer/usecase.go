package customer

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	domain "erpagendamentos/backend/internal/domain/customer"
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrInvalidCustomer  = errors.New("invalid customer")
)

type UseCase struct {
	repository domain.Repository
}

func NewUseCase(repository domain.Repository) *UseCase {
	return &UseCase{repository: repository}
}

func (u *UseCase) ListByTenant(tenantID string, search string) ([]domain.Customer, error) {
	return u.repository.ListByTenant(strings.TrimSpace(tenantID), strings.TrimSpace(search))
}

func (u *UseCase) GetByID(tenantID string, id string) (domain.Customer, error) {
	item, err := u.repository.GetByID(strings.TrimSpace(tenantID), strings.TrimSpace(id))
	if err != nil {
		return domain.Customer{}, ErrCustomerNotFound
	}
	return item, nil
}

func (u *UseCase) Create(input domain.Customer) (domain.Customer, error) {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.Name = strings.TrimSpace(input.Name)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.Birthday = strings.TrimSpace(input.Birthday)
	input.Notes = strings.TrimSpace(input.Notes)
	if input.TenantID == "" || input.Name == "" {
		return domain.Customer{}, ErrInvalidCustomer
	}
	input.ID = uuid.NewString()
	return u.repository.Create(input)
}

func (u *UseCase) Update(input domain.Customer) (domain.Customer, error) {
	tenantID := strings.TrimSpace(input.TenantID)
	id := strings.TrimSpace(input.ID)
	current, err := u.repository.GetByID(tenantID, id)
	if err != nil {
		return domain.Customer{}, ErrCustomerNotFound
	}
	input.TenantID = tenantID
	input.CreatedAt = current.CreatedAt
	if strings.TrimSpace(input.Name) == "" {
		input.Name = current.Name
	}
	input.Phone = strings.TrimSpace(input.Phone)
	if strings.TrimSpace(input.Phone) == "" {
		input.Phone = current.Phone
	}
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if strings.TrimSpace(input.Email) == "" {
		input.Email = current.Email
	}
	if strings.TrimSpace(input.Birthday) == "" {
		input.Birthday = current.Birthday
	}
	if strings.TrimSpace(input.Notes) == "" {
		input.Notes = current.Notes
	}
	return u.repository.Update(input)
}

func (u *UseCase) Delete(tenantID string, id string) error {
	if err := u.repository.Delete(strings.TrimSpace(tenantID), strings.TrimSpace(id)); err != nil {
		return ErrCustomerNotFound
	}
	return nil
}

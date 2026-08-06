package professional

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	domain "erpagendamentos/backend/internal/domain/professional"
)

var (
	ErrProfessionalNotFound = errors.New("professional not found")
	ErrInvalidProfessional  = errors.New("invalid professional")
)

type UseCase struct {
	repository domain.Repository
}

func NewUseCase(repository domain.Repository) *UseCase {
	return &UseCase{repository: repository}
}

func (u *UseCase) ListByTenant(tenantID string) ([]domain.Professional, error) {
	return u.repository.ListByTenant(strings.TrimSpace(tenantID))
}

func (u *UseCase) GetByID(tenantID string, id string) (domain.Professional, error) {
	professional, err := u.repository.GetByID(strings.TrimSpace(tenantID), strings.TrimSpace(id))
	if err != nil {
		return domain.Professional{}, ErrProfessionalNotFound
	}
	return professional, nil
}

func (u *UseCase) Create(input domain.Professional) (domain.Professional, error) {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.Phone = strings.TrimSpace(input.Phone)
	input.Color = strings.TrimSpace(input.Color)
	if input.TenantID == "" || input.Name == "" {
		return domain.Professional{}, ErrInvalidProfessional
	}
	input.ID = uuid.NewString()
	if input.Color == "" {
		input.Color = "#3B82F6"
	}
	return u.repository.Create(input)
}

func (u *UseCase) Update(input domain.Professional) (domain.Professional, error) {
	tenantID := strings.TrimSpace(input.TenantID)
	id := strings.TrimSpace(input.ID)
	current, err := u.repository.GetByID(tenantID, id)
	if err != nil {
		return domain.Professional{}, ErrProfessionalNotFound
	}
	input.TenantID = tenantID
	input.CreatedAt = current.CreatedAt
	if strings.TrimSpace(input.Name) == "" {
		input.Name = current.Name
	}
	if strings.TrimSpace(input.Color) == "" {
		input.Color = current.Color
	}
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if strings.TrimSpace(input.Email) == "" {
		input.Email = current.Email
	}
	if strings.TrimSpace(input.Phone) == "" {
		input.Phone = current.Phone
	}
	return u.repository.Update(input)
}

func (u *UseCase) SetActive(tenantID string, id string, active bool) (domain.Professional, error) {
	professional, err := u.repository.SetActive(strings.TrimSpace(tenantID), strings.TrimSpace(id), active)
	if err != nil {
		return domain.Professional{}, ErrProfessionalNotFound
	}
	return professional, nil
}

func (u *UseCase) Delete(tenantID string, id string) error {
	if err := u.repository.Delete(strings.TrimSpace(tenantID), strings.TrimSpace(id)); err != nil {
		return ErrProfessionalNotFound
	}
	return nil
}

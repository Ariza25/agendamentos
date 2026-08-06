package configuration

import (
	"errors"
	"strings"

	domain "erpagendamentos/backend/internal/domain/configuration"
)

var (
	ErrInvalidConfiguration = errors.New("invalid configuration")
)

type UseCase struct {
	repository domain.Repository
}

func NewUseCase(repository domain.Repository) *UseCase {
	return &UseCase{repository: repository}
}

func (u *UseCase) Get(tenantID string) (domain.Configuration, error) {
	return u.repository.Get(strings.TrimSpace(tenantID))
}

func (u *UseCase) Update(input domain.Configuration) (domain.Configuration, error) {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.EstablishmentName = strings.TrimSpace(input.EstablishmentName)
	if input.TenantID == "" {
		return domain.Configuration{}, ErrInvalidConfiguration
	}
	if input.AppointmentIntervalMinutes <= 0 {
		input.AppointmentIntervalMinutes = 30
	}
	if input.DefaultDurationMinutes <= 0 {
		input.DefaultDurationMinutes = 30
	}
	if input.AdvanceNoticeHours < 0 {
		input.AdvanceNoticeHours = 0
	}
	return u.repository.Update(input)
}

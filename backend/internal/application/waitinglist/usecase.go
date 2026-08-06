package waitinglist

import (
	"errors"
	"strings"
	"time"

	domain "erpagendamentos/backend/internal/domain/waitinglist"
)

var (
	ErrWaitingListNotFound = errors.New("waiting list entry not found")
	ErrInvalidWaitingList  = errors.New("invalid waiting list entry")
)

type UseCase struct {
	repository domain.Repository
}

func NewUseCase(repository domain.Repository) *UseCase {
	return &UseCase{repository: repository}
}

func (u *UseCase) Create(input domain.Entry) (domain.Entry, error) {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.ServiceID = strings.TrimSpace(input.ServiceID)
	input.ProfessionalID = strings.TrimSpace(input.ProfessionalID)
	input.CustomerID = strings.TrimSpace(input.CustomerID)
	input.Name = strings.TrimSpace(input.Name)
	input.Phone = strings.TrimSpace(input.Phone)
	if input.TenantID == "" || input.Name == "" || input.Phone == "" {
		return domain.Entry{}, ErrInvalidWaitingList
	}
	if input.Status == "" {
		input.Status = domain.StatusWaiting
	}
	return u.repository.Create(input)
}

func (u *UseCase) ListByTenant(tenantID string, status string) ([]domain.Entry, error) {
	return u.repository.ListByTenant(strings.TrimSpace(tenantID), strings.TrimSpace(status))
}

func (u *UseCase) GetByID(tenantID string, id string) (domain.Entry, error) {
	item, err := u.repository.GetByID(strings.TrimSpace(tenantID), strings.TrimSpace(id))
	if err != nil {
		return domain.Entry{}, ErrWaitingListNotFound
	}
	return item, nil
}

func (u *UseCase) UpdateStatus(tenantID string, id string, status string) (domain.Entry, error) {
	tenantID = strings.TrimSpace(tenantID)
	id = strings.TrimSpace(id)
	status = strings.TrimSpace(status)
	if status != domain.StatusWaiting && status != domain.StatusNotified && status != domain.StatusConverted && status != domain.StatusCancelled {
		return domain.Entry{}, ErrInvalidWaitingList
	}
	item, err := u.repository.UpdateStatus(tenantID, id, status)
	if err != nil {
		return domain.Entry{}, ErrWaitingListNotFound
	}
	return item, nil
}

func (u *UseCase) ListWaitingForSlot(tenantID string, serviceID string, professionalID string) ([]domain.Entry, error) {
	return u.repository.ListWaitingForSlot(
		strings.TrimSpace(tenantID),
		strings.TrimSpace(serviceID),
		strings.TrimSpace(professionalID),
		"oldest",
	)
}

func (u *UseCase) MarkNotified(tenantID string, id string) (domain.Entry, error) {
	item, err := u.repository.UpdateStatus(strings.TrimSpace(tenantID), strings.TrimSpace(id), domain.StatusNotified)
	if err != nil {
		return domain.Entry{}, ErrWaitingListNotFound
	}
	return item, nil
}

func (u *UseCase) CreatedAt() string {
	return time.Now().UTC().Format(time.RFC3339)
}

package notification

import (
	"errors"
	"strings"
	"time"

	domain "erpagendamentos/backend/internal/domain/notification"
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
	ErrNotificationFailed   = errors.New("notification failed")
)

type UseCase struct {
	repository domain.Repository
	provider   domain.Provider
}

func NewUseCase(repository domain.Repository, provider domain.Provider) *UseCase {
	return &UseCase{repository: repository, provider: provider}
}

func (u *UseCase) ListByTenant(tenantID string, limit int) ([]domain.Log, error) {
	return u.repository.ListByTenant(strings.TrimSpace(tenantID), limit)
}

func (u *UseCase) Log(input domain.Log) (domain.Log, error) {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.Recipient = strings.TrimSpace(input.Recipient)
	if input.TenantID == "" || input.Recipient == "" {
		return domain.Log{}, ErrNotificationNotFound
	}
	if input.Status == "" {
		input.Status = "sent"
	}
	input.SentAt = time.Now().UTC().Format(time.RFC3339)
	return u.repository.Create(input)
}

func (u *UseCase) Send(input domain.Log) (domain.Log, error) {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.Recipient = strings.TrimSpace(input.Recipient)
	input.Subject = strings.TrimSpace(input.Subject)
	input.Body = strings.TrimSpace(input.Body)
	input.Type = strings.TrimSpace(input.Type)
	if input.Channel == "" {
		input.Channel = "email"
	}
	input.Channel = strings.TrimSpace(input.Channel)
	if input.TenantID == "" || input.Recipient == "" || input.Subject == "" || input.Body == "" {
		return domain.Log{}, ErrNotificationNotFound
	}

	input.Status = "sent"
	if u.provider != nil {
		if err := u.provider.Send(input.Channel, input.Recipient, input.Subject, input.Body); err != nil {
			input.Status = "failed"
			input.Error = err.Error()
		}
	}
	input.SentAt = time.Now().UTC().Format(time.RFC3339)
	return u.repository.Create(input)
}

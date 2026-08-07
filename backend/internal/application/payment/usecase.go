package payment

import (
	"errors"
	"strings"
	"time"

	domain "erpagendamentos/backend/internal/domain/payment"
)

var (
	ErrPaymentNotFound = errors.New("payment not found")
	ErrInvalidPayment  = errors.New("invalid payment")
)

type UseCase struct {
	repository domain.Repository
	provider   domain.Provider
}

func NewUseCase(repository domain.Repository, provider domain.Provider) *UseCase {
	return &UseCase{repository: repository, provider: provider}
}

func (u *UseCase) ListByTenant(tenantID string, limit int) ([]domain.Payment, error) {
	return u.repository.ListByTenant(strings.TrimSpace(tenantID), limit)
}

func (u *UseCase) ListByRange(tenantID string, from string, to string) ([]domain.Payment, error) {
	return u.repository.ListByTenantRange(strings.TrimSpace(tenantID), strings.TrimSpace(from), strings.TrimSpace(to))
}

func (u *UseCase) GetByID(tenantID string, id string) (domain.Payment, error) {
	item, err := u.repository.GetByID(strings.TrimSpace(tenantID), strings.TrimSpace(id))
	if err != nil {
		return domain.Payment{}, ErrPaymentNotFound
	}
	return item, nil
}

func (u *UseCase) GetByAppointment(tenantID string, appointmentID string) (domain.Payment, error) {
	item, err := u.repository.GetByAppointment(strings.TrimSpace(tenantID), strings.TrimSpace(appointmentID))
	if err != nil {
		return domain.Payment{}, ErrPaymentNotFound
	}
	return item, nil
}

func (u *UseCase) Create(input domain.Payment) (domain.Payment, error) {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.AppointmentID = strings.TrimSpace(input.AppointmentID)
	input.Description = strings.TrimSpace(input.Description)
	if input.TenantID == "" || input.AmountCents <= 0 {
		return domain.Payment{}, ErrInvalidPayment
	}
	if input.Method == "" {
		input.Method = string(domain.MethodPix)
	}
	if input.Method != string(domain.MethodPix) && input.Method != string(domain.MethodCard) {
		return domain.Payment{}, ErrInvalidPayment
	}
	if input.Status == "" {
		input.Status = string(domain.StatusPending)
	}
	return u.repository.Create(input)
}

func (u *UseCase) CreatePaymentIntent(tenantID string, input domain.Payment) (domain.Payment, error) {
	created, err := u.Create(input)
	if err != nil {
		return domain.Payment{}, err
	}
	if u.provider == nil {
		return created, nil
	}

	result, err := u.provider.CreatePayment(domain.IntentInput{
		AmountCents: created.AmountCents,
		Description: created.Description,
	})
	if err != nil {
		return domain.Payment{}, err
	}
	return u.repository.UpdateStatus(created.TenantID, created.ID, result.Status, result.ProviderPaymentID)
}

func (u *UseCase) UpdateStatus(tenantID string, id string, status string, providerPaymentID string) (domain.Payment, error) {
	item, err := u.repository.UpdateStatus(strings.TrimSpace(tenantID), strings.TrimSpace(id), strings.TrimSpace(status), strings.TrimSpace(providerPaymentID))
	if err != nil {
		return domain.Payment{}, ErrPaymentNotFound
	}
	return item, nil
}

func (u *UseCase) MarkPaid(tenantID string, id string) (domain.Payment, error) {
	item, err := u.repository.UpdateStatus(strings.TrimSpace(tenantID), strings.TrimSpace(id), string(domain.StatusPaid), "")
	if err != nil {
		return domain.Payment{}, ErrPaymentNotFound
	}
	return item, nil
}

func (u *UseCase) Refund(tenantID string, id string) (domain.Payment, error) {
	tenantID = strings.TrimSpace(tenantID)
	id = strings.TrimSpace(id)
	current, err := u.repository.GetByID(tenantID, id)
	if err != nil {
		return domain.Payment{}, ErrPaymentNotFound
	}
	if current.Status != string(domain.StatusPaid) {
		return domain.Payment{}, ErrInvalidPayment
	}
	if u.provider != nil && current.ProviderPaymentID != "" {
		if err := u.provider.Refund(current.ProviderPaymentID); err != nil {
			return domain.Payment{}, err
		}
	}
	return u.repository.UpdateStatus(tenantID, id, string(domain.StatusRefunded), current.ProviderPaymentID)
}

func (u *UseCase) SyncProviderStatus(tenantID string, id string) (domain.Payment, error) {
	tenantID = strings.TrimSpace(tenantID)
	id = strings.TrimSpace(id)
	current, err := u.repository.GetByID(tenantID, id)
	if err != nil {
		return domain.Payment{}, ErrPaymentNotFound
	}
	if u.provider == nil || current.ProviderPaymentID == "" {
		return current, nil
	}
	providerStatus, err := u.provider.GetPaymentStatus(current.ProviderPaymentID)
	if err != nil {
		return domain.Payment{}, err
	}
	if providerStatus == current.Status {
		return current, nil
	}
	normalized := normalizeStatus(providerStatus)
	if normalized == "" || normalized == current.Status {
		return current, nil
	}
	return u.repository.UpdateStatus(tenantID, id, normalized, current.ProviderPaymentID)
}

func normalizeStatus(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "succeeded", "paid", "approved", "completed", "captured":
		return string(domain.StatusPaid)
	case "pending", "requires_payment_method", "processing":
		return string(domain.StatusPending)
	case "failed", "canceled", "cancelled", "declined", "requires_action":
		return string(domain.StatusFailed)
	case "refunded", "reversed":
		return string(domain.StatusRefunded)
	default:
		return ""
	}
}

func (u *UseCase) SetPaidAt(tenantID string, id string, paidAt string) (domain.Payment, error) {
	item, err := u.repository.UpdateStatus(strings.TrimSpace(tenantID), strings.TrimSpace(id), string(domain.StatusPaid), "")
	if err != nil {
		return domain.Payment{}, ErrPaymentNotFound
	}
	if trimmed := strings.TrimSpace(paidAt); trimmed != "" {
		item.PaidAt = trimmed
	}
	return item, nil
}

func (u *UseCase) Now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

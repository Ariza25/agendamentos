package memory

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/payment"
)

type PaymentRepository struct {
	mu       sync.Mutex
	payments map[string]payment.Payment
}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{
		payments: make(map[string]payment.Payment),
	}
}

func (r *PaymentRepository) Create(input payment.Payment) (payment.Payment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item := input
	item.ID = uuid.NewString()
	item.CreatedAt = time.Now().Format(time.RFC3339)
	if item.Status == string(payment.StatusPaid) {
		item.PaidAt = item.CreatedAt
	}
	r.payments[item.ID] = item
	return item, nil
}

func (r *PaymentRepository) GetByID(tenantID string, id string) (payment.Payment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.payments[strings.TrimSpace(id)]
	if !exists || item.TenantID != strings.TrimSpace(tenantID) {
		return payment.Payment{}, fmt.Errorf("payment not found")
	}
	return item, nil
}

func (r *PaymentRepository) GetByAppointment(tenantID string, appointmentID string) (payment.Payment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID = strings.TrimSpace(tenantID)
	appointmentID = strings.TrimSpace(appointmentID)
	latest := payment.Payment{}
	found := false
	for _, item := range r.payments {
		if item.TenantID == tenantID && item.AppointmentID == appointmentID {
			if !found || item.CreatedAt > latest.CreatedAt {
				latest = item
				found = true
			}
		}
	}
	if !found {
		return payment.Payment{}, fmt.Errorf("payment not found")
	}
	return latest, nil
}

func (r *PaymentRepository) GetByProviderPaymentID(providerPaymentID string) (payment.Payment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	providerPaymentID = strings.TrimSpace(providerPaymentID)
	for _, item := range r.payments {
		if item.ProviderPaymentID == providerPaymentID {
			return item, nil
		}
	}
	return payment.Payment{}, fmt.Errorf("payment not found")
}

func (r *PaymentRepository) ListByTenant(tenantID string, limit int) ([]payment.Payment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if limit <= 0 || limit > 200 {
		limit = 50
	}
	tenantID = strings.TrimSpace(tenantID)
	items := []payment.Payment{}
	for _, item := range r.payments {
		if item.TenantID == tenantID {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt > items[j].CreatedAt
	})
	if len(items) > limit {
		items = items[:limit]
	}
	return items, nil
}

func (r *PaymentRepository) ListByTenantRange(tenantID string, from string, to string) ([]payment.Payment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID = strings.TrimSpace(tenantID)
	from = strings.TrimSpace(from)
	to = strings.TrimSpace(to)
	items := []payment.Payment{}
	for _, item := range r.payments {
		if item.TenantID != tenantID {
			continue
		}
		if from != "" && item.CreatedAt < from {
			continue
		}
		if to != "" && item.CreatedAt > to {
			continue
		}
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt > items[j].CreatedAt
	})
	return items, nil
}

func (r *PaymentRepository) UpdateStatus(tenantID string, id string, status string, providerPaymentID string) (payment.Payment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.payments[strings.TrimSpace(id)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return payment.Payment{}, fmt.Errorf("payment not found")
	}
	current.Status = status
	if strings.TrimSpace(providerPaymentID) != "" {
		current.ProviderPaymentID = providerPaymentID
	}
	if status == string(payment.StatusPaid) {
		current.PaidAt = time.Now().Format(time.RFC3339)
	}
	r.payments[current.ID] = current
	return current, nil
}

func (r *PaymentRepository) UpdateIntent(tenantID string, id string, result payment.IntentResult) (payment.Payment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.payments[strings.TrimSpace(id)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return payment.Payment{}, fmt.Errorf("payment not found")
	}
	current.Status = result.Status
	current.ProviderPaymentID = result.ProviderPaymentID
	current.PixCopiaECola = result.PixCopiaECola
	current.PixQRCodeBase64 = result.PixQRCodeBase64
	current.ProviderData = result.ProviderData
	if result.Status == string(payment.StatusPaid) {
		current.PaidAt = time.Now().Format(time.RFC3339)
	}
	r.payments[current.ID] = current
	return current, nil
}

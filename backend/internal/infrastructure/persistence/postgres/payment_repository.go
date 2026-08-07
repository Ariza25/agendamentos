package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/payment"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(input payment.Payment) (payment.Payment, error) {
	item := input
	item.ID = uuid.NewString()
	item.CreatedAt = time.Now().Format(time.RFC3339)
	if _, err := r.db.Exec(`
		INSERT INTO catalog.payments (id, tenant_id, appointment_id, amount_cents, method, provider, provider_payment_id, status, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, item.ID, item.TenantID, item.AppointmentID, item.AmountCents, item.Method, item.Provider, item.ProviderPaymentID, item.Status, item.Description); err != nil {
		return payment.Payment{}, fmt.Errorf("insert payment: %w", err)
	}
	return r.GetByID(item.TenantID, item.ID)
}

func (r *PaymentRepository) GetByID(tenantID string, id string) (payment.Payment, error) {
	return scanPaymentRow(r.db.QueryRow(`
		SELECT id, tenant_id, appointment_id, amount_cents, method, provider, provider_payment_id, status, description, paid_at, created_at
		FROM catalog.payments
		WHERE id = $1 AND tenant_id = $2
	`, strings.TrimSpace(id), strings.TrimSpace(tenantID)))
}

func (r *PaymentRepository) GetByAppointment(tenantID string, appointmentID string) (payment.Payment, error) {
	return scanPaymentRow(r.db.QueryRow(`
		SELECT id, tenant_id, appointment_id, amount_cents, method, provider, provider_payment_id, status, description, paid_at, created_at
		FROM catalog.payments
		WHERE tenant_id = $1 AND appointment_id = $2
		ORDER BY created_at DESC
		LIMIT 1
	`, strings.TrimSpace(tenantID), strings.TrimSpace(appointmentID)))
}

func (r *PaymentRepository) ListByTenant(tenantID string, limit int) ([]payment.Payment, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.db.Query(`
		SELECT id, tenant_id, appointment_id, amount_cents, method, provider, provider_payment_id, status, description, paid_at, created_at
		FROM catalog.payments
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, strings.TrimSpace(tenantID), limit)
	if err != nil {
		return nil, fmt.Errorf("query payments: %w", err)
	}
	defer rows.Close()

	items := []payment.Payment{}
	for rows.Next() {
		item, err := scanPaymentRows(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PaymentRepository) ListByTenantRange(tenantID string, from string, to string) ([]payment.Payment, error) {
	query := `
		SELECT id, tenant_id, appointment_id, amount_cents, method, provider, provider_payment_id, status, description, paid_at, created_at
		FROM catalog.payments
		WHERE tenant_id = $1`
	args := []any{strings.TrimSpace(tenantID)}
	if from = strings.TrimSpace(from); from != "" {
		args = append(args, from)
		query += fmt.Sprintf(" AND created_at >= $%d", len(args))
	}
	if to = strings.TrimSpace(to); to != "" {
		args = append(args, to)
		query += fmt.Sprintf(" AND created_at <= $%d", len(args))
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query payments: %w", err)
	}
	defer rows.Close()

	items := []payment.Payment{}
	for rows.Next() {
		item, err := scanPaymentRows(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PaymentRepository) UpdateStatus(tenantID string, id string, status string, providerPaymentID string) (payment.Payment, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.payments
		SET status = $3,
		    provider_payment_id = COALESCE(NULLIF($4, ''), provider_payment_id),
		    paid_at = CASE WHEN $3 = 'paid' THEN NOW() ELSE paid_at END
		WHERE id = $1 AND tenant_id = $2
	`, id, tenantID, status, providerPaymentID)
	if err != nil {
		return payment.Payment{}, fmt.Errorf("update payment status: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return payment.Payment{}, fmt.Errorf("check payment update: %w", err)
	}
	if affected == 0 {
		return payment.Payment{}, fmt.Errorf("payment not found")
	}
	return r.GetByID(tenantID, id)
}

type paymentScanner interface {
	Scan(dest ...any) error
}

func scanPaymentRow(scanner paymentScanner) (payment.Payment, error) {
	return scanPaymentRows(scanner)
}

func scanPaymentRows(scanner paymentScanner) (payment.Payment, error) {
	var item payment.Payment
	var paidAt, createdAt sql.NullTime
	if err := scanner.Scan(&item.ID, &item.TenantID, &item.AppointmentID, &item.AmountCents, &item.Method, &item.Provider, &item.ProviderPaymentID, &item.Status, &item.Description, &paidAt, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			return payment.Payment{}, fmt.Errorf("payment not found")
		}
		return payment.Payment{}, fmt.Errorf("scan payment: %w", err)
	}
	if paidAt.Valid {
		item.PaidAt = paidAt.Time.UTC().Format(time.RFC3339)
	}
	item.CreatedAt = createdAt.Time.Format(time.RFC3339)
	return item, nil
}

package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/cashregister"
)

type CashRegisterRepository struct {
	db *sql.DB
}

func NewCashRegisterRepository(db *sql.DB) *CashRegisterRepository {
	return &CashRegisterRepository{db: db}
}

func (r *CashRegisterRepository) Create(item cashregister.CashRegister) (cashregister.CashRegister, error) {
	if item.ID == "" {
		item.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	err := r.db.QueryRow(`
		INSERT INTO catalog.cash_registers (id, tenant_id, opened_by_user_id, status, initial_balance_cents, system_balance_cents, opened_at, notes)
		VALUES ($1, $2, $3, 'open', $4, $4, $5, $6)
		RETURNING opened_at
	`, item.ID, item.TenantID, item.OpenedByUserID, item.InitialBalanceCents, now, item.Notes).Scan(&now)
	if err != nil {
		return cashregister.CashRegister{}, fmt.Errorf("create cash register: %w", err)
	}
	item.OpenedAt = now.Format(time.RFC3339)
	item.Status = "open"
	item.SystemBalanceCents = item.InitialBalanceCents
	return item, nil
}

func (r *CashRegisterRepository) Update(item cashregister.CashRegister) (cashregister.CashRegister, error) {
	var closedAt *time.Time
	if item.ClosedAt != "" {
		if t, err := time.Parse(time.RFC3339, item.ClosedAt); err == nil {
			closedAt = &t
		}
	}
	if item.Status == "closed" && closedAt == nil {
		now := time.Now().UTC()
		closedAt = &now
	}

	_, err := r.db.Exec(`
		UPDATE catalog.cash_registers
		SET closed_by_user_id = $1, status = $2, closing_balance_cents = $3,
		    system_balance_cents = $4, difference_cents = $5, closed_at = $6, notes = $7
		WHERE tenant_id = $8 AND id = $9
	`, item.ClosedByUserID, item.Status, item.ClosingBalanceCents, item.SystemBalanceCents, item.DifferenceCents, closedAt, item.Notes, item.TenantID, item.ID)
	if err != nil {
		return cashregister.CashRegister{}, fmt.Errorf("update cash register: %w", err)
	}
	return r.GetByID(item.TenantID, item.ID)
}

func (r *CashRegisterRepository) GetByID(tenantID string, id string) (cashregister.CashRegister, error) {
	var item cashregister.CashRegister
	var openedAt time.Time
	var closedAt *time.Time
	err := r.db.QueryRow(`
		SELECT id, tenant_id, opened_by_user_id, closed_by_user_id, status,
		       initial_balance_cents, closing_balance_cents, system_balance_cents, difference_cents,
		       opened_at, closed_at, notes
		FROM catalog.cash_registers
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, id).Scan(
		&item.ID, &item.TenantID, &item.OpenedByUserID, &item.ClosedByUserID, &item.Status,
		&item.InitialBalanceCents, &item.ClosingBalanceCents, &item.SystemBalanceCents, &item.DifferenceCents,
		&openedAt, &closedAt, &item.Notes,
	)
	if err != nil {
		return cashregister.CashRegister{}, err
	}
	item.OpenedAt = openedAt.Format(time.RFC3339)
	if closedAt != nil {
		item.ClosedAt = closedAt.Format(time.RFC3339)
	}
	return item, nil
}

func (r *CashRegisterRepository) GetCurrentOpen(tenantID string) (cashregister.CashRegister, error) {
	var item cashregister.CashRegister
	var openedAt time.Time
	var closedAt *time.Time
	err := r.db.QueryRow(`
		SELECT id, tenant_id, opened_by_user_id, closed_by_user_id, status,
		       initial_balance_cents, closing_balance_cents, system_balance_cents, difference_cents,
		       opened_at, closed_at, notes
		FROM catalog.cash_registers
		WHERE tenant_id = $1 AND status = 'open'
		ORDER BY opened_at DESC
		LIMIT 1
	`, tenantID).Scan(
		&item.ID, &item.TenantID, &item.OpenedByUserID, &item.ClosedByUserID, &item.Status,
		&item.InitialBalanceCents, &item.ClosingBalanceCents, &item.SystemBalanceCents, &item.DifferenceCents,
		&openedAt, &closedAt, &item.Notes,
	)
	if err != nil {
		return cashregister.CashRegister{}, err
	}
	item.OpenedAt = openedAt.Format(time.RFC3339)
	if closedAt != nil {
		item.ClosedAt = closedAt.Format(time.RFC3339)
	}
	return item, nil
}

func (r *CashRegisterRepository) List(tenantID string, limit int) ([]cashregister.CashRegister, error) {
	if limit <= 0 || limit > 100 {
		limit = 30
	}
	rows, err := r.db.Query(`
		SELECT id, tenant_id, opened_by_user_id, closed_by_user_id, status,
		       initial_balance_cents, closing_balance_cents, system_balance_cents, difference_cents,
		       opened_at, closed_at, notes
		FROM catalog.cash_registers
		WHERE tenant_id = $1
		ORDER BY opened_at DESC
		LIMIT $2
	`, tenantID, limit)
	if err != nil {
		return nil, fmt.Errorf("list cash registers: %w", err)
	}
	defer rows.Close()

	var list []cashregister.CashRegister
	for rows.Next() {
		var item cashregister.CashRegister
		var openedAt time.Time
		var closedAt *time.Time
		if err := rows.Scan(
			&item.ID, &item.TenantID, &item.OpenedByUserID, &item.ClosedByUserID, &item.Status,
			&item.InitialBalanceCents, &item.ClosingBalanceCents, &item.SystemBalanceCents, &item.DifferenceCents,
			&openedAt, &closedAt, &item.Notes,
		); err != nil {
			return nil, err
		}
		item.OpenedAt = openedAt.Format(time.RFC3339)
		if closedAt != nil {
			item.ClosedAt = closedAt.Format(time.RFC3339)
		}
		list = append(list, item)
	}
	return list, nil
}

func (r *CashRegisterRepository) CreateEntry(entry cashregister.CashRegisterEntry) (cashregister.CashRegisterEntry, error) {
	if entry.ID == "" {
		entry.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	err := r.db.QueryRow(`
		INSERT INTO catalog.cash_register_entries (id, tenant_id, cash_register_id, type, amount_cents, reason, created_by_user_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at
	`, entry.ID, entry.TenantID, entry.CashRegisterID, entry.Type, entry.AmountCents, strings.TrimSpace(entry.Reason), entry.CreatedByUserID, now).Scan(&now)
	if err != nil {
		return cashregister.CashRegisterEntry{}, fmt.Errorf("create cash register entry: %w", err)
	}
	entry.CreatedAt = now.Format(time.RFC3339)
	return entry, nil
}

func (r *CashRegisterRepository) ListEntries(tenantID string, cashRegisterID string) ([]cashregister.CashRegisterEntry, error) {
	rows, err := r.db.Query(`
		SELECT id, tenant_id, cash_register_id, type, amount_cents, reason, created_by_user_id, created_at
		FROM catalog.cash_register_entries
		WHERE tenant_id = $1 AND cash_register_id = $2
		ORDER BY created_at ASC
	`, tenantID, cashRegisterID)
	if err != nil {
		return nil, fmt.Errorf("list cash register entries: %w", err)
	}
	defer rows.Close()

	var list []cashregister.CashRegisterEntry
	for rows.Next() {
		var entry cashregister.CashRegisterEntry
		var createdAt time.Time
		if err := rows.Scan(&entry.ID, &entry.TenantID, &entry.CashRegisterID, &entry.Type, &entry.AmountCents, &entry.Reason, &entry.CreatedByUserID, &createdAt); err != nil {
			return nil, err
		}
		entry.CreatedAt = createdAt.Format(time.RFC3339)
		list = append(list, entry)
	}
	return list, nil
}


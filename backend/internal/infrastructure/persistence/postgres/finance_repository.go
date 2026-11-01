package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"erpagendamentos/backend/internal/domain/finance"
	"github.com/google/uuid"
)

type FinanceRepository struct {
	db *sql.DB
}

func NewFinanceRepository(db *sql.DB) *FinanceRepository {
	return &FinanceRepository{db: db}
}

func (r *FinanceRepository) Create(tx finance.Transaction) (finance.Transaction, error) {
	if tx.ID == "" {
		tx.ID = uuid.NewString()
	}
	if tx.DueDate == "" {
		tx.DueDate = time.Now().Format("2006-01-02")
	}
	now := time.Now().UTC()
	var paidAt *time.Time
	if tx.Status == "paid" {
		paidAt = &now
	}

	err := r.db.QueryRow(`
		INSERT INTO catalog.financial_transactions (id, tenant_id, type, category, amount_cents, description, reference_id, due_date, paid_at, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING created_at, updated_at
	`, tx.ID, tx.TenantID, tx.Type, tx.Category, tx.AmountCents, tx.Description, tx.ReferenceID, tx.DueDate, paidAt, tx.Status, now, now).Scan(&tx.CreatedAt, &tx.UpdatedAt)
	if err != nil {
		return finance.Transaction{}, fmt.Errorf("create transaction: %w", err)
	}
	if paidAt != nil {
		tx.PaidAt = paidAt.Format(time.RFC3339)
	}
	return tx, nil
}

func (r *FinanceRepository) Update(tx finance.Transaction) (finance.Transaction, error) {
	now := time.Now().UTC()
	var paidAt *time.Time
	if tx.Status == "paid" {
		if tx.PaidAt != "" {
			if t, err := time.Parse(time.RFC3339, tx.PaidAt); err == nil {
				paidAt = &t
			}
		}
		if paidAt == nil {
			paidAt = &now
		}
	}

	result, err := r.db.Exec(`
		UPDATE catalog.financial_transactions
		SET type = $1, category = $2, amount_cents = $3, description = $4, reference_id = $5, due_date = $6, paid_at = $7, status = $8, updated_at = $9
		WHERE tenant_id = $10 AND id = $11
	`, tx.Type, tx.Category, tx.AmountCents, tx.Description, tx.ReferenceID, tx.DueDate, paidAt, tx.Status, now, tx.TenantID, tx.ID)
	if err != nil {
		return finance.Transaction{}, fmt.Errorf("update transaction: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return finance.Transaction{}, errors.New("transaction not found")
	}
	tx.UpdatedAt = now.Format(time.RFC3339)
	if paidAt != nil {
		tx.PaidAt = paidAt.Format(time.RFC3339)
	}
	return tx, nil
}

func (r *FinanceRepository) GetByID(tenantID string, id string) (finance.Transaction, error) {
	var tx finance.Transaction
	var paidAt *time.Time
	var created, updated time.Time
	err := r.db.QueryRow(`
		SELECT id, tenant_id, type, category, amount_cents, description, reference_id, due_date, paid_at, status, created_at, updated_at
		FROM catalog.financial_transactions
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, id).Scan(&tx.ID, &tx.TenantID, &tx.Type, &tx.Category, &tx.AmountCents, &tx.Description, &tx.ReferenceID, &tx.DueDate, &paidAt, &tx.Status, &created, &updated)
	if err != nil {
		return finance.Transaction{}, fmt.Errorf("get transaction: %w", err)
	}
	tx.CreatedAt = created.Format(time.RFC3339)
	tx.UpdatedAt = updated.Format(time.RFC3339)
	if paidAt != nil {
		tx.PaidAt = paidAt.Format(time.RFC3339)
	}
	return tx, nil
}

func (r *FinanceRepository) List(tenantID string, txType string, status string, fromDate string, toDate string) ([]finance.Transaction, error) {
	query := `
		SELECT id, tenant_id, type, category, amount_cents, description, reference_id, due_date, paid_at, status, created_at, updated_at
		FROM catalog.financial_transactions
		WHERE tenant_id = $1
	`
	args := []any{tenantID}
	if txType != "" {
		args = append(args, txType)
		query += fmt.Sprintf(" AND type = $%d", len(args))
	}
	if status != "" {
		args = append(args, status)
		query += fmt.Sprintf(" AND status = $%d", len(args))
	}
	if fromDate != "" {
		args = append(args, fromDate)
		query += fmt.Sprintf(" AND due_date >= $%d", len(args))
	}
	if toDate != "" {
		args = append(args, toDate)
		query += fmt.Sprintf(" AND due_date <= $%d", len(args))
	}
	query += " ORDER BY due_date DESC, created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list transactions: %w", err)
	}
	defer rows.Close()

	var transactions []finance.Transaction
	for rows.Next() {
		var tx finance.Transaction
		var paidAt *time.Time
		var created, updated time.Time
		if err := rows.Scan(&tx.ID, &tx.TenantID, &tx.Type, &tx.Category, &tx.AmountCents, &tx.Description, &tx.ReferenceID, &tx.DueDate, &paidAt, &tx.Status, &created, &updated); err != nil {
			return nil, err
		}
		tx.CreatedAt = created.Format(time.RFC3339)
		tx.UpdatedAt = updated.Format(time.RFC3339)
		if paidAt != nil {
			tx.PaidAt = paidAt.Format(time.RFC3339)
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}

func (r *FinanceRepository) Delete(tenantID string, id string) error {
	_, err := r.db.Exec("DELETE FROM catalog.financial_transactions WHERE tenant_id = $1 AND id = $2", tenantID, id)
	return err
}

func (r *FinanceRepository) GetSummary(tenantID string, fromDate string, toDate string) (finance.Summary, error) {
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN type = 'income' AND status = 'paid' THEN amount_cents ELSE 0 END), 0) AS total_income,
			COALESCE(SUM(CASE WHEN type = 'expense' AND status = 'paid' THEN amount_cents ELSE 0 END), 0) AS total_expense,
			COALESCE(SUM(CASE WHEN type = 'income' AND status = 'pending' THEN amount_cents ELSE 0 END), 0) AS pending_income,
			COALESCE(SUM(CASE WHEN type = 'expense' AND status = 'pending' THEN amount_cents ELSE 0 END), 0) AS pending_expense
		FROM catalog.financial_transactions
		WHERE tenant_id = $1
	`
	args := []any{tenantID}
	if fromDate != "" {
		args = append(args, fromDate)
		query += fmt.Sprintf(" AND due_date >= $%d", len(args))
	}
	if toDate != "" {
		args = append(args, toDate)
		query += fmt.Sprintf(" AND due_date <= $%d", len(args))
	}

	var summary finance.Summary
	err := r.db.QueryRow(query, args...).Scan(&summary.TotalIncomeCents, &summary.TotalExpenseCents, &summary.PendingIncomeCents, &summary.PendingExpenseCents)
	if err != nil {
		return summary, fmt.Errorf("get financial summary: %w", err)
	}
	summary.NetProfitCents = summary.TotalIncomeCents - summary.TotalExpenseCents

	// Categories breakdown
	catQuery := `
		SELECT category, type, COALESCE(SUM(amount_cents), 0) AS amount_cents
		FROM catalog.financial_transactions
		WHERE tenant_id = $1
	`
	if fromDate != "" {
		catQuery += fmt.Sprintf(" AND due_date >= $%d", len(args)) // re-uses args if length matches
	}
	catQuery += " GROUP BY category, type ORDER BY amount_cents DESC"

	catRows, err := r.db.Query(catQuery, args...)
	if err == nil {
		defer catRows.Close()
		for catRows.Next() {
			var cat finance.CategoryTotal
			if err := catRows.Scan(&cat.Category, &cat.Type, &cat.AmountCents); err == nil {
				summary.ByCategory = append(summary.ByCategory, cat)
			}
		}
	}

	return summary, nil
}


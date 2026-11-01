package memory

import (
	"errors"
	"sync"
	"time"

	"erpagendamentos/backend/internal/domain/finance"
	"github.com/google/uuid"
)

type FinanceRepository struct {
	mu           sync.RWMutex
	transactions map[string]finance.Transaction
}

func NewFinanceRepository() *FinanceRepository {
	return &FinanceRepository{
		transactions: make(map[string]finance.Transaction),
	}
}

func (r *FinanceRepository) Create(tx finance.Transaction) (finance.Transaction, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if tx.ID == "" {
		tx.ID = uuid.NewString()
	}
	now := time.Now().UTC().Format(time.RFC3339)
	tx.CreatedAt = now
	tx.UpdatedAt = now
	if tx.DueDate == "" {
		tx.DueDate = time.Now().Format("2006-01-02")
	}
	if tx.Status == "paid" && tx.PaidAt == "" {
		tx.PaidAt = now
	}
	r.transactions[tx.TenantID+":"+tx.ID] = tx
	return tx, nil
}

func (r *FinanceRepository) Update(tx finance.Transaction) (finance.Transaction, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := tx.TenantID + ":" + tx.ID
	existing, found := r.transactions[key]
	if !found {
		return finance.Transaction{}, errors.New("transaction not found")
	}
	tx.CreatedAt = existing.CreatedAt
	tx.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if tx.Status == "paid" && tx.PaidAt == "" {
		tx.PaidAt = tx.UpdatedAt
	}
	r.transactions[key] = tx
	return tx, nil
}

func (r *FinanceRepository) GetByID(tenantID string, id string) (finance.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tx, found := r.transactions[tenantID+":"+id]
	if !found {
		return finance.Transaction{}, errors.New("transaction not found")
	}
	return tx, nil
}

func (r *FinanceRepository) List(tenantID string, txType string, status string, fromDate string, toDate string) ([]finance.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []finance.Transaction
	for _, tx := range r.transactions {
		if tx.TenantID == tenantID {
			if txType != "" && tx.Type != txType {
				continue
			}
			if status != "" && tx.Status != status {
				continue
			}
			if fromDate != "" && tx.DueDate < fromDate {
				continue
			}
			if toDate != "" && tx.DueDate > toDate {
				continue
			}
			result = append(result, tx)
		}
	}
	return result, nil
}

func (r *FinanceRepository) Delete(tenantID string, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.transactions, tenantID+":"+id)
	return nil
}

func (r *FinanceRepository) GetSummary(tenantID string, fromDate string, toDate string) (finance.Summary, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var summary finance.Summary
	categories := make(map[string]*finance.CategoryTotal)

	for _, tx := range r.transactions {
		if tx.TenantID == tenantID {
			if fromDate != "" && tx.DueDate < fromDate {
				continue
			}
			if toDate != "" && tx.DueDate > toDate {
				continue
			}

			catKey := tx.Category + ":" + tx.Type
			cat, exists := categories[catKey]
			if !exists {
				cat = &finance.CategoryTotal{
					Category: tx.Category,
					Type:     tx.Type,
				}
				categories[catKey] = cat
			}
			cat.AmountCents += int64(tx.AmountCents)

			if tx.Status == "paid" {
				if tx.Type == "income" {
					summary.TotalIncomeCents += int64(tx.AmountCents)
				} else if tx.Type == "expense" {
					summary.TotalExpenseCents += int64(tx.AmountCents)
				}
			} else if tx.Status == "pending" {
				if tx.Type == "income" {
					summary.PendingIncomeCents += int64(tx.AmountCents)
				} else if tx.Type == "expense" {
					summary.PendingExpenseCents += int64(tx.AmountCents)
				}
			}
		}
	}
	summary.NetProfitCents = summary.TotalIncomeCents - summary.TotalExpenseCents

	for _, c := range categories {
		summary.ByCategory = append(summary.ByCategory, *c)
	}

	return summary, nil
}

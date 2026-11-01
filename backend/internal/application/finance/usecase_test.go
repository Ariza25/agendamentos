package finance_test

import (
	"testing"

	"erpagendamentos/backend/internal/application/finance"
	financedomain "erpagendamentos/backend/internal/domain/finance"
	"erpagendamentos/backend/internal/infrastructure/persistence/memory"
)

func TestFinanceUseCase(t *testing.T) {
	repo := memory.NewFinanceRepository()
	uc := finance.NewUseCase(repo)

	// Income: appointment payment
	err := uc.RecordAppointmentIncome("tenant-1", "appt-1", 15000, "Haircut & Styling")
	if err != nil {
		t.Fatalf("record income error: %v", err)
	}

	// Expense: rent
	_, err = uc.Create(financedomain.Transaction{
		TenantID:    "tenant-1",
		Type:        "expense",
		Category:    "rent",
		AmountCents: 5000,
		Description: "Office Rent",
		Status:      "paid",
	})
	if err != nil {
		t.Fatalf("create expense error: %v", err)
	}

	// Expense: pending supplies bill
	_, err = uc.Create(financedomain.Transaction{
		TenantID:    "tenant-1",
		Type:        "expense",
		Category:    "supplies",
		AmountCents: 2000,
		Description: "Products Supplier Invoice",
		Status:      "pending",
	})
	if err != nil {
		t.Fatalf("create pending expense error: %v", err)
	}

	// Summary
	summary, err := uc.GetSummary("tenant-1", "", "")
	if err != nil {
		t.Fatalf("summary error: %v", err)
	}
	t.Logf("Summary: %+v", summary)
	if summary.TotalIncomeCents != 15000 {
		t.Fatalf("expected income 15000, got %d", summary.TotalIncomeCents)
	}
	if summary.TotalExpenseCents != 5000 {
		t.Fatalf("expected paid expense 5000, got %d", summary.TotalExpenseCents)
	}
	if summary.NetProfitCents != 10000 {
		t.Fatalf("expected net profit 10000, got %d", summary.NetProfitCents)
	}
	if summary.PendingExpenseCents != 2000 {
		t.Fatalf("expected pending expense 2000, got %d", summary.PendingExpenseCents)
	}
}

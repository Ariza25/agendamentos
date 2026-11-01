package commission_test

import (
	"testing"

	"erpagendamentos/backend/internal/application/commission"
	commissiondomain "erpagendamentos/backend/internal/domain/commission"
	"erpagendamentos/backend/internal/infrastructure/persistence/memory"
)

func TestCommissionUseCase(t *testing.T) {
	repo := memory.NewCommissionRepository()
	uc := commission.NewUseCase(repo)

	// Create global 40% rule
	_, err := uc.CreateRule(commissiondomain.CommissionRule{
		TenantID:   "tenant-1",
		Percentage: 40.0,
	})
	if err != nil {
		t.Fatalf("create rule error: %v", err)
	}

	// Specific rule for prof-1 on srv-1: 50% + 1000 cents
	_, err = uc.CreateRule(commissiondomain.CommissionRule{
		TenantID:         "tenant-1",
		ProfessionalID:   "prof-1",
		ServiceID:        "srv-1",
		Percentage:       50.0,
		FixedAmountCents: 1000,
	})
	if err != nil {
		t.Fatalf("create specific rule error: %v", err)
	}

	// Calculate commission for prof-1 on srv-1 (price 10,000 cents)
	rec, err := uc.CalculateAndRecord("tenant-1", "appt-1", "prof-1", "Ana", "srv-1", "Corte", 10000)
	if err != nil {
		t.Fatalf("calculate error: %v", err)
	}
	// 50% of 10000 = 5000 + 1000 = 6000 cents
	if rec.CommissionAmountCents != 6000 {
		t.Fatalf("expected 6000 cents commission, got %d", rec.CommissionAmountCents)
	}
	if rec.Status != "pending" {
		t.Fatalf("expected pending status, got %s", rec.Status)
	}

	// Summaries
	summaries, err := uc.GetSummaries("tenant-1")
	if err != nil || len(summaries) == 0 {
		t.Fatalf("get summaries error: %v, summaries: %+v", err, summaries)
	}
	if summaries[0].PendingCents != 6000 {
		t.Fatalf("expected 6000 pending, got %d", summaries[0].PendingCents)
	}

	// Settle
	err = uc.SettleCommissions("tenant-1", "prof-1", []string{rec.ID})
	if err != nil {
		t.Fatalf("settle error: %v", err)
	}

	afterSettle, err := uc.GetSummaries("tenant-1")
	if err != nil || len(afterSettle) == 0 {
		t.Fatalf("get summaries after settle error: %v", err)
	}
	if afterSettle[0].PendingCents != 0 || afterSettle[0].PaidCents != 6000 {
		t.Fatalf("expected 0 pending and 6000 paid, got %+v", afterSettle[0])
	}
}


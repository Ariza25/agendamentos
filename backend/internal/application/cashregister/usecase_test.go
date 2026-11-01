package cashregister_test

import (
	"testing"

	"erpagendamentos/backend/internal/application/cashregister"
	"erpagendamentos/backend/internal/infrastructure/persistence/memory"
)

func TestCashRegisterFlow(t *testing.T) {
	regRepo := memory.NewCashRegisterRepository()
	finRepo := memory.NewFinanceRepository()
	uc := cashregister.NewUseCase(regRepo, finRepo)

	tenantID := "test-tenant-1"
	userID := "user-123"

	// 1. Initial state: no shift open
	_, err := uc.GetCurrentShift(tenantID)
	if err == nil {
		t.Fatalf("expected error when no shift is open, got nil")
	}

	// 2. Open Shift
	shift, err := uc.OpenShift(tenantID, userID, 10000, "Abertura com R$ 100")
	if err != nil {
		t.Fatalf("failed to open shift: %v", err)
	}
	if shift.InitialBalanceCents != 10000 || shift.Status != "open" {
		t.Fatalf("unexpected shift initial balance or status: %+v", shift)
	}

	// Cannot open second shift
	_, err = uc.OpenShift(tenantID, userID, 5000, "Second open")
	if err != cashregister.ErrShiftAlreadyOpen {
		t.Fatalf("expected ErrShiftAlreadyOpen, got: %v", err)
	}

	// 3. Suprimento (aporte)
	entrySup, err := uc.RecordEntry(tenantID, userID, "suprimento", 2000, "Troco extra")
	if err != nil {
		t.Fatalf("failed to record suprimento: %v", err)
	}
	if entrySup.AmountCents != 2000 {
		t.Fatalf("unexpected entry amount: %d", entrySup.AmountCents)
	}

	// 4. Sangria (retirada)
	entrySangria, err := uc.RecordEntry(tenantID, userID, "sangria", 1500, "Pagamento fornecedor pão")
	if err != nil {
		t.Fatalf("failed to record sangria: %v", err)
	}
	if entrySangria.AmountCents != 1500 {
		t.Fatalf("unexpected sangria amount: %d", entrySangria.AmountCents)
	}

	// 5. Check Current Shift Summary
	summary, err := uc.GetCurrentShift(tenantID)
	if err != nil {
		t.Fatalf("failed to get current shift: %v", err)
	}
	// Expected drawer: 10000 (initial) + 2000 (suprimento) - 1500 (sangria) = 10500
	if summary.CurrentDrawerCents != 10500 {
		t.Fatalf("expected drawer cents 10500, got: %d", summary.CurrentDrawerCents)
	}
	if summary.TotalSangriasCents != 1500 || summary.TotalSuprimentosCents != 2000 {
		t.Fatalf("unexpected entries totals in summary: %+v", summary)
	}

	// 6. Close Shift
	closed, err := uc.CloseShift(tenantID, userID, 10500, "Caixa conferido sem quebra")
	if err != nil {
		t.Fatalf("failed to close shift: %v", err)
	}
	if closed.Status != "closed" || closed.DifferenceCents != 0 {
		t.Fatalf("unexpected closed shift: %+v", closed)
	}
}


package stock_test

import (
	"testing"

	"erpagendamentos/backend/internal/application/stock"
	stockdomain "erpagendamentos/backend/internal/domain/stock"
	"erpagendamentos/backend/internal/infrastructure/persistence/memory"
)

func TestStockUseCase(t *testing.T) {
	repo := memory.NewStockRepository()
	uc := stock.NewUseCase(repo)

	// Create product
	prod, err := uc.CreateProduct(stockdomain.Product{
		TenantID:      "tenant-1",
		Name:          "Shampoo 1L",
		SKU:           "SHP-001",
		Category:      "Hair Care",
		PriceCents:    5000,
		CostCents:     2500,
		StockQuantity: 10,
		MinStockAlert: 2,
		IsSupply:      false,
		Unit:          "un",
		Active:        true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prod.ID == "" || prod.Name != "Shampoo 1L" {
		t.Fatalf("unexpected product: %+v", prod)
	}

	// Stock movement (in)
	mov, err := uc.CreateMovement(stockdomain.StockMovement{
		TenantID:  "tenant-1",
		ProductID: prod.ID,
		Type:      "in",
		Quantity:  5,
	})
	if err != nil {
		t.Fatalf("movement error: %v", err)
	}
	if mov.Quantity != 5 {
		t.Fatalf("unexpected movement: %+v", mov)
	}

	// Verify updated stock
	updated, err := uc.GetProduct("tenant-1", prod.ID)
	if err != nil {
		t.Fatalf("get product error: %v", err)
	}
	if updated.StockQuantity != 15 {
		t.Fatalf("expected stock 15, got %v", updated.StockQuantity)
	}

	// Service supplies & auto deduction
	err = uc.SetServiceSupplies("tenant-1", "srv-1", []stockdomain.ServiceSupply{
		{ProductID: prod.ID, QuantityConsumed: 2},
	})
	if err != nil {
		t.Fatalf("set supplies error: %v", err)
	}

	err = uc.DeductSuppliesForAppointment("tenant-1", "appt-1", "srv-1")
	if err != nil {
		t.Fatalf("deduct error: %v", err)
	}

	afterDeduct, err := uc.GetProduct("tenant-1", prod.ID)
	if err != nil {
		t.Fatalf("get product error: %v", err)
	}
	if afterDeduct.StockQuantity != 13 {
		t.Fatalf("expected stock 13 after deduction, got %v", afterDeduct.StockQuantity)
	}
}


package checkout_test

import (
	"testing"

	"erpagendamentos/backend/internal/application/checkout"
	appointmentdomain "erpagendamentos/backend/internal/domain/appointment"
	professionaldomain "erpagendamentos/backend/internal/domain/professional"
	servicedomain "erpagendamentos/backend/internal/domain/service"
	stockdomain "erpagendamentos/backend/internal/domain/stock"
	"erpagendamentos/backend/internal/infrastructure/persistence/memory"
)

func TestCheckoutFlow(t *testing.T) {
	apptRepo := memory.NewAppointmentRepository()
	servRepo := memory.NewServiceRepository()
	profRepo := memory.NewProfessionalRepository()
	stockRepo := memory.NewStockRepository()
	commRepo := memory.NewCommissionRepository()
	finRepo := memory.NewFinanceRepository()
	cashRepo := memory.NewCashRegisterRepository()

	uc := checkout.NewUseCase(apptRepo, servRepo, profRepo, stockRepo, commRepo, finRepo, cashRepo)

	tenantID := "tenant-pos-1"

	// 1. Setup professional with 50% commission
	prof, err := profRepo.Create(professionaldomain.Professional{
		TenantID:             tenantID,
		Name:                 "Marcos Barbeiro",
		CommissionPercentage: 50,
		Active:               true,
	})
	if err != nil {
		t.Fatalf("create prof: %v", err)
	}

	// 2. Setup service (R$ 50,00)
	serv, err := servRepo.Create(servicedomain.Service{
		TenantID:        tenantID,
		Name:            "Corte Degradê",
		PriceCents:      5000,
		DurationMinutes: 30,
		Active:          true,
	})
	if err != nil {
		t.Fatalf("create service: %v", err)
	}

	// 3. Setup product in stock (Pomada R$ 40,00, stock 10)
	prod, err := stockRepo.CreateProduct(stockdomain.Product{
		TenantID:      tenantID,
		Name:          "Pomada Modeladora Matte",
		PriceCents:    4000,
		StockQuantity: 10,
		Active:        true,
	})
	if err != nil {
		t.Fatalf("create product: %v", err)
	}

	// 4. Setup scheduled appointment
	appt, err := apptRepo.Create(appointmentdomain.Appointment{
		TenantID:       tenantID,
		ProfessionalID: prof.ID,
		ServiceID:      serv.ID,
		CustomerID:     "cust-1",
		Status:         "scheduled",
		PaymentStatus:  "pending",
	})
	if err != nil {
		t.Fatalf("create appointment: %v", err)
	}

	// 5. Execute Checkout: 1 Service (50.00) + 1 Product (40.00) - 5.00 discount = 85.00
	res, err := uc.Checkout(tenantID, checkout.CheckoutInput{
		AppointmentID: appt.ID,
		Products: []checkout.ProductItem{
			{ProductID: prod.ID, Quantity: 1},
		},
		DiscountCents: 500,
		PaymentMethod: "pix",
		Notes:         "Cliente satisfeito",
	})
	if err != nil {
		t.Fatalf("checkout failed: %v", err)
	}

	if res.Appointment.Status != "completed" || res.Appointment.PaymentStatus != "paid" {
		t.Fatalf("appointment not updated to completed/paid: %+v", res.Appointment)
	}
	if res.ServiceTotalCents != 5000 {
		t.Fatalf("expected service total 5000, got: %d", res.ServiceTotalCents)
	}
	if res.ProductsTotalCents != 4000 {
		t.Fatalf("expected products total 4000, got: %d", res.ProductsTotalCents)
	}
	if res.TotalChargedCents != 8500 {
		t.Fatalf("expected total charged 8500, got: %d", res.TotalChargedCents)
	}
	// Commission: 50% of 5000 = 2500
	if res.CommissionCents != 2500 {
		t.Fatalf("expected commission 2500, got: %d", res.CommissionCents)
	}

	// Verify stock was decremented
	updatedProd, err := stockRepo.GetProductByID(tenantID, prod.ID)
	if err != nil || updatedProd.StockQuantity != 9 {
		t.Fatalf("expected stock quantity 9, got: %f", updatedProd.StockQuantity)
	}

	// Verify commission record was created
	records, err := commRepo.ListRecords(tenantID, prof.ID, "pending")
	if err != nil || len(records) != 1 || records[0].CommissionAmountCents != 2500 {
		t.Fatalf("expected 1 commission record of 2500, got: %+v", records)
	}
}


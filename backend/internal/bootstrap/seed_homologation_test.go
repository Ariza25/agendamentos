package bootstrap

import (
	"testing"
)

func TestSeedHomologationData(t *testing.T) {
	container := NewContainer()
	defer container.Close()

	// Check Tenant
	tenant, err := container.tenantUseCase.GetTenant("default")
	if err != nil {
		t.Fatalf("expected tenant default to exist: %v", err)
	}
	if tenant.Name != "Barbearia & Studio Dom Navalha" {
		t.Errorf("expected tenant name 'Barbearia & Studio Dom Navalha', got '%s'", tenant.Name)
	}
	if tenant.Slug != "dom-navalha" {
		t.Errorf("expected slug 'dom-navalha', got '%s'", tenant.Slug)
	}

	// Check Customers
	custs, err := container.customerUseCase.ListByTenant("default", "")
	if err != nil {
		t.Fatalf("expected customers list to succeed: %v", err)
	}
	if len(custs) < 5 {
		t.Errorf("expected at least 5 customers, got %d", len(custs))
	}

	// Check Professionals
	pros, err := container.professionalUseCase.ListByTenant("default")
	if err != nil {
		t.Fatalf("expected professionals list to succeed: %v", err)
	}
	if len(pros) < 3 {
		t.Errorf("expected at least 3 professionals, got %d", len(pros))
	}

	// Check Services
	svcs, err := container.serviceUseCase.ListByTenant("default")
	if err != nil {
		t.Fatalf("expected services list to succeed: %v", err)
	}
	if len(svcs) < 5 {
		t.Errorf("expected at least 5 services, got %d", len(svcs))
	}

	// Check Appointments
	appts, err := container.appointmentUseCase.ListByRange("default", "", "", "", "", "")
	if err != nil {
		t.Fatalf("expected appointments list to succeed: %v", err)
	}
	if len(appts) < 10 {
		t.Errorf("expected at least 10 appointments, got %d", len(appts))
	}

	// Check Products
	products, err := container.stockUseCase.ListProducts("default", false)
	if err != nil {
		t.Fatalf("expected products list to succeed: %v", err)
	}
	if len(products) < 4 {
		t.Errorf("expected at least 4 products, got %d", len(products))
	}

	// Check Cash Register
	shift, err := container.cashRegisterUseCase.GetCurrentShift("default")
	if err != nil {
		t.Fatalf("expected open cash register shift: %v", err)
	}
	if shift.CashRegister.Status != "open" {
		t.Errorf("expected open status, got '%s'", shift.CashRegister.Status)
	}

	// Check Reviews
	reviews, err := container.reviewUseCase.ListByTenant("default", "")
	if err != nil {
		t.Fatalf("expected reviews list to succeed: %v", err)
	}
	if len(reviews) < 3 {
		t.Errorf("expected at least 3 reviews, got %d", len(reviews))
	}
}

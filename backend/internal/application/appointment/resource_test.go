package appointment

import (
	"testing"

	domain "erpagendamentos/backend/internal/domain/appointment"
	resourcedomain "erpagendamentos/backend/internal/domain/resource"
	"erpagendamentos/backend/internal/infrastructure/persistence/memory"
)

func newResourceUseCase() (*UseCase, resourcedomain.Repository) {
	appointmentRepo := memory.NewAppointmentRepository()
	resourceRepo := memory.NewResourceRepository()
	useCase := NewUseCase(appointmentRepo, resourceRepo)
	return useCase, resourceRepo
}

func newAppointment(tenantID string, professionalID string, startsAt string, endsAt string, resourceID string) domain.Appointment {
	return domain.Appointment{
		TenantID:       tenantID,
		ProfessionalID: professionalID,
		ServiceID:      "",
		ResourceID:     resourceID,
		CustomerID:     "customer-1",
		StartsAt:       startsAt,
		EndsAt:         endsAt,
		Status:         StatusScheduled,
	}
}

func TestCreateAutoAssignsAvailableResource(t *testing.T) {
	useCase, resourceRepo := newResourceUseCase()

	_, _ = resourceRepo.Create(resourcedomain.Resource{TenantID: "default", Name: "Sala 1", Type: "room", Capacity: 1, Active: true})
	_, _ = resourceRepo.Create(resourcedomain.Resource{TenantID: "default", Name: "Sala 2", Type: "room", Capacity: 1, Active: true})

	created, err := useCase.Create(newAppointment("default", "professional-1", "2026-01-05T10:00:00Z", "2026-01-05T10:30:00Z", ""))
	if err != nil {
		t.Fatalf("create appointment: %v", err)
	}
	if created.ResourceID == "" {
		t.Fatalf("expected an auto-assigned resource, got empty")
	}

	second, err := useCase.Create(newAppointment("default", "professional-2", "2026-01-05T10:00:00Z", "2026-01-05T10:30:00Z", ""))
	if err != nil {
		t.Fatalf("create second appointment: %v", err)
	}
	if second.ResourceID == "" {
		t.Fatalf("expected an auto-assigned resource for the second appointment, got empty")
	}
	if second.ResourceID == created.ResourceID {
		t.Fatalf("expected a different resource after the first was booked, got %s", second.ResourceID)
	}
}

func TestCreateRejectsBusyExplicitResource(t *testing.T) {
	useCase, resourceRepo := newResourceUseCase()

	resource, _ := resourceRepo.Create(resourcedomain.Resource{TenantID: "default", Name: "Sala 1", Type: "room", Capacity: 1, Active: true})

	if _, err := useCase.Create(newAppointment("default", "professional-1", "2026-01-05T10:00:00Z", "2026-01-05T10:30:00Z", resource.ID)); err != nil {
		t.Fatalf("create first appointment: %v", err)
	}
	_, err := useCase.Create(newAppointment("default", "professional-2", "2026-01-05T10:15:00Z", "2026-01-05T10:45:00Z", resource.ID))
	if err == nil {
		t.Fatalf("expected conflict when booking a busy resource")
	}
	if err != ErrTimeConflict {
		t.Fatalf("expected ErrTimeConflict, got %v", err)
	}
}

func TestCreateRejectsUnknownResource(t *testing.T) {
	useCase, _ := newResourceUseCase()

	_, err := useCase.Create(newAppointment("default", "professional-1", "2026-01-05T10:00:00Z", "2026-01-05T10:30:00Z", "missing-resource"))
	if err == nil {
		t.Fatalf("expected error for unknown resource")
	}
	if err != ErrInvalidAppointment {
		t.Fatalf("expected ErrInvalidAppointment, got %v", err)
	}
}

func TestUpdateAutoReassignsFreedResource(t *testing.T) {
	useCase, resourceRepo := newResourceUseCase()

	resource, _ := resourceRepo.Create(resourcedomain.Resource{TenantID: "default", Name: "Sala 1", Type: "room", Capacity: 1, Active: true})
	created, err := useCase.Create(newAppointment("default", "professional-1", "2026-01-05T10:00:00Z", "2026-01-05T10:30:00Z", ""))
	if err != nil {
		t.Fatalf("create appointment: %v", err)
	}
	if created.ResourceID != resource.ID {
		t.Fatalf("expected resource %s, got %s", resource.ID, created.ResourceID)
	}

	conflicting, err := useCase.Create(newAppointment("default", "professional-2", "2026-01-05T10:10:00Z", "2026-01-05T10:40:00Z", resource.ID))
	if err == nil || err != ErrTimeConflict {
		t.Fatalf("expected conflict on busy resource, got %v", err)
	}
	_ = conflicting

	updated, err := useCase.Create(newAppointment("default", "professional-3", "2026-01-05T10:05:00Z", "2026-01-05T10:20:00Z", ""))
	if err != nil {
		t.Fatalf("create later appointment: %v", err)
	}
	if updated.ResourceID != "" {
		t.Fatalf("expected no resource left available, got %s", updated.ResourceID)
	}
}

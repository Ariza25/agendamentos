package booking

import (
	"errors"
	"strings"
	"testing"
	"time"

	appappointment "erpagendamentos/backend/internal/application/appointment"
	configurationdomain "erpagendamentos/backend/internal/domain/configuration"
	professionaldomain "erpagendamentos/backend/internal/domain/professional"
	servicedomain "erpagendamentos/backend/internal/domain/service"
	tenantdomain "erpagendamentos/backend/internal/domain/tenant"
	"erpagendamentos/backend/internal/infrastructure/persistence/memory"
)

func newTestUseCase() *UseCase {
	tenantRepo := memory.NewTenantRepository()
	configurationRepo := memory.NewConfigurationRepository()
	serviceRepo := memory.NewServiceRepository()
	professionalRepo := memory.NewProfessionalRepository()
	customerRepo := memory.NewCustomerRepository()

	_, _ = tenantRepo.Create(tenantdomain.Tenant{
		ID: "default", Name: "Teste", Slug: "default", Active: true,
	})
	_, _ = professionalRepo.Create(professionaldomain.Professional{
		TenantID: "default", Name: "Dra. Ana", Active: true,
	})
	_, _ = serviceRepo.Create(servicedomain.Service{
		TenantID: "default", Name: "Corte", DurationMinutes: 30, Active: true,
	})

	appointmentRepo := memory.NewAppointmentRepository()
	appointments := appappointment.NewUseCase(appointmentRepo, memory.NewResourceRepository())
	return NewUseCase(tenantRepo, configurationRepo, serviceRepo, professionalRepo, memory.NewResourceRepository(), customerRepo, appointments)
}

func TestCreateBookingGeneratesManagementCode(t *testing.T) {
	useCase := newTestUseCase()
	future := time.Now().Add(24 * time.Hour)
	day := future.Format("2006-01-02")
	_, _ = useCase.configurationRepo.Update(configurationdomain.Configuration{
		TenantID:                   "default",
		EstablishmentName:          "Teste",
		AllowOnlineBooking:         true,
		WorkingHours:               []configurationdomain.WorkingHours{{DayOfWeek: int(future.Weekday()), OpenTime: "09:00", CloseTime: "18:00", Enabled: true}},
		AppointmentIntervalMinutes: 30,
		DefaultDurationMinutes:     30,
	})

	professional, err := useCase.professionalRepo.ListByTenant("default")
	if err != nil || len(professional) == 0 {
		t.Fatalf("no professionals: %v", err)
	}

	result, err := useCase.CreateBooking("default", BookingInput{
		ServiceID:      "",
		ProfessionalID: professional[0].ID,
		StartsAt:       day + "T10:00",
		Name:           "Maria",
		Phone:          "11999998888",
	})
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}
	if strings.TrimSpace(result.ManagementCode) == "" {
		t.Fatalf("expected management code, got empty")
	}

	loaded, err := useCase.GetByManagementCode("default", result.ManagementCode)
	if err != nil {
		t.Fatalf("get by management code: %v", err)
	}
	if loaded.ID != result.Appointment.ID {
		t.Fatalf("expected appointment %s, got %s", result.Appointment.ID, loaded.ID)
	}

	if _, err := useCase.GetByManagementCode("default", "INVALID-XX"); err == nil {
		t.Fatalf("expected error for invalid code")
	}
}

func TestCancelPublic(t *testing.T) {
	useCase := newTestUseCase()
	future := time.Now().Add(48 * time.Hour)
	day := future.Format("2006-01-02")
	_, _ = useCase.configurationRepo.Update(configurationdomain.Configuration{
		TenantID:                   "default",
		EstablishmentName:          "Teste",
		AllowOnlineBooking:         true,
		WorkingHours:               []configurationdomain.WorkingHours{{DayOfWeek: int(future.Weekday()), OpenTime: "09:00", CloseTime: "18:00", Enabled: true}},
		AppointmentIntervalMinutes: 30,
		DefaultDurationMinutes:     30,
	})

	professional, _ := useCase.professionalRepo.ListByTenant("default")

	result, err := useCase.CreateBooking("default", BookingInput{
		ProfessionalID: professional[0].ID,
		StartsAt:       day + "T11:00",
		Name:           "João",
		Phone:          "11999997777",
	})
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}

	cancelled, err := useCase.CancelPublic("default", result.ManagementCode, "mudança de planos")
	if err != nil {
		t.Fatalf("cancel booking: %v", err)
	}
	if cancelled.Status != "cancelled" {
		t.Fatalf("expected cancelled status, got %s", cancelled.Status)
	}

	if _, err := useCase.CancelPublic("default", result.ManagementCode, ""); err == nil {
		t.Fatalf("expected error when cancelling an already cancelled appointment")
	}
}

func TestAvailableSlotsRespectsClosedDays(t *testing.T) {
	useCase := newTestUseCase()
	day := time.Now().Add(72 * time.Hour)
	_, _ = useCase.configurationRepo.Update(configurationdomain.Configuration{
		TenantID:                   "default",
		EstablishmentName:          "Teste",
		AllowOnlineBooking:         true,
		WorkingHours:               []configurationdomain.WorkingHours{{DayOfWeek: int(day.Weekday()), OpenTime: "09:00", CloseTime: "18:00", Enabled: true}},
		AppointmentIntervalMinutes: 30,
		DefaultDurationMinutes:     30,
		ClosedDays:                 []configurationdomain.ClosedDay{{Date: day.Format("2006-01-02"), Name: "Feriado", AllDay: true}},
	})

	professional, _ := useCase.professionalRepo.ListByTenant("default")
	availability, err := useCase.Availability("default", day.Format("2006-01-02"), "", professional[0].ID)
	if err != nil {
		t.Fatalf("availability: %v", err)
	}
	slots, ok := availability["slots"].([]Slot)
	if !ok {
		t.Fatalf("expected slots key")
	}
	if len(slots) != 0 {
		t.Fatalf("expected no slots on closed day, got %d", len(slots))
	}
}

func TestManagementCodeNormalization(t *testing.T) {
	raw := "AB-CDE-12"
	hash := hashManagementCode(raw)
	other := hashManagementCode(strings.ToLower(raw))
	if hash != other {
		t.Fatalf("hash should be case-insensitive")
	}
	if hashManagementCode(raw) == hashManagementCode("ZZ-CDE-12") {
		t.Fatalf("different codes should produce different hashes")
	}
}

func TestPortalLookup(t *testing.T) {
	useCase := newTestUseCase()
	future := time.Now().Add(24 * time.Hour)
	day := future.Format("2006-01-02")
	_, _ = useCase.configurationRepo.Update(configurationdomain.Configuration{
		TenantID:                   "default",
		EstablishmentName:          "Teste",
		AllowOnlineBooking:         true,
		WorkingHours:               []configurationdomain.WorkingHours{{DayOfWeek: int(future.Weekday()), OpenTime: "09:00", CloseTime: "18:00", Enabled: true}},
		AppointmentIntervalMinutes: 30,
		DefaultDurationMinutes:     30,
	})

	professional, _ := useCase.professionalRepo.ListByTenant("default")

	result, err := useCase.CreateBooking("default", BookingInput{
		ProfessionalID: professional[0].ID,
		ServiceID:      "",
		StartsAt:       day + "T10:00",
		Name:           "Portal Cliente",
		Phone:          "11988887777",
	})
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}

	portal, err := useCase.PortalLookup("default", "11988887777")
	if err != nil {
		t.Fatalf("portal lookup: %v", err)
	}
	if portal.Customer.Phone != "11988887777" {
		t.Fatalf("expected phone match, got %s", portal.Customer.Phone)
	}
	if len(portal.Appointments) != 1 {
		t.Fatalf("expected 1 appointment, got %d", len(portal.Appointments))
	}
	if portal.Appointments[0].ID != result.Appointment.ID {
		t.Fatalf("expected appointment %s, got %s", result.Appointment.ID, portal.Appointments[0].ID)
	}
	if portal.Appointments[0].ProfessionalName != "Dra. Ana" {
		t.Fatalf("expected professional name, got %q", portal.Appointments[0].ProfessionalName)
	}

	if _, err := useCase.PortalLookup("default", "11900001111"); err == nil {
		t.Fatalf("expected error for unknown phone")
	}
}

func TestPortalCancelValidatesOwnership(t *testing.T) {
	useCase := newTestUseCase()
	future := time.Now().Add(48 * time.Hour)
	day := future.Format("2006-01-02")
	_, _ = useCase.configurationRepo.Update(configurationdomain.Configuration{
		TenantID:                   "default",
		EstablishmentName:          "Teste",
		AllowOnlineBooking:         true,
		WorkingHours:               []configurationdomain.WorkingHours{{DayOfWeek: int(future.Weekday()), OpenTime: "09:00", CloseTime: "18:00", Enabled: true}},
		AppointmentIntervalMinutes: 30,
		DefaultDurationMinutes:     30,
	})

	professional, _ := useCase.professionalRepo.ListByTenant("default")

	result, err := useCase.CreateBooking("default", BookingInput{
		ProfessionalID: professional[0].ID,
		StartsAt:       day + "T11:00",
		Name:           "Dona",
		Phone:          "11966665555",
	})
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}

	cancelled, err := useCase.PortalCancel("default", "11966665555", result.Appointment.ID, "imprevisto")
	if err != nil {
		t.Fatalf("portal cancel: %v", err)
	}
	if cancelled.Status != "cancelled" {
		t.Fatalf("expected cancelled, got %s", cancelled.Status)
	}

	if _, err := useCase.PortalCancel("default", "11900001111", result.Appointment.ID, ""); !errors.Is(err, ErrCustomerNotFound) {
		t.Fatalf("expected customer not found, got %v", err)
	}
	if _, err := useCase.PortalCancel("default", "11966665555", "00000000-0000-0000-0000-000000000000", ""); !errors.Is(err, ErrAppointmentAccess) {
		t.Fatalf("expected appointment access error, got %v", err)
	}
	if _, err := useCase.PortalCancel("default", "11966665555", result.Appointment.ID, ""); !errors.Is(err, ErrAppointmentDone) {
		t.Fatalf("expected appointment done error, got %v", err)
	}
}

func TestPortalReschedule(t *testing.T) {
	useCase := newTestUseCase()
	future := time.Now().Add(72 * time.Hour)
	day := future.Format("2006-01-02")
	_, _ = useCase.configurationRepo.Update(configurationdomain.Configuration{
		TenantID:                   "default",
		EstablishmentName:          "Teste",
		AllowOnlineBooking:         true,
		WorkingHours:               []configurationdomain.WorkingHours{{DayOfWeek: int(future.Weekday()), OpenTime: "09:00", CloseTime: "18:00", Enabled: true}},
		AppointmentIntervalMinutes: 30,
		DefaultDurationMinutes:     30,
	})

	professional, _ := useCase.professionalRepo.ListByTenant("default")

	result, err := useCase.CreateBooking("default", BookingInput{
		ProfessionalID: professional[0].ID,
		StartsAt:       day + "T09:00",
		Name:           "Sr. Fulano",
		Phone:          "11955554444",
	})
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}

	updated, err := useCase.PortalReschedule("default", "11955554444", result.Appointment.ID, day+"T14:00")
	if err != nil {
		t.Fatalf("portal reschedule: %v", err)
	}
	if updated.StartsAt[:16] != day+"T14:00" && updated.StartsAt[:10] != day {
		t.Fatalf("expected rescheduled time, got %s", updated.StartsAt)
	}

	if _, err := useCase.PortalReschedule("default", "11955554444", result.Appointment.ID, day+"T14:30"); err != nil {
		t.Fatalf("portal reschedule of an active appointment should be allowed, got %v", err)
	}
}

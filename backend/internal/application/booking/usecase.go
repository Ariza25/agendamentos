package booking

import (
	"errors"
	"fmt"
	"strings"
	"time"

	appappointment "erpagendamentos/backend/internal/application/appointment"
	appointmentdomain "erpagendamentos/backend/internal/domain/appointment"
	configurationdomain "erpagendamentos/backend/internal/domain/configuration"
	customerdomain "erpagendamentos/backend/internal/domain/customer"
	professionaldomain "erpagendamentos/backend/internal/domain/professional"
	servicedomain "erpagendamentos/backend/internal/domain/service"
	tenantdomain "erpagendamentos/backend/internal/domain/tenant"
)

var (
	ErrTenantNotFound    = errors.New("tenant not found")
	ErrBookingDisabled   = errors.New("online booking disabled")
	ErrInvalidBooking    = errors.New("invalid booking")
	ErrSlotUnavailable   = errors.New("slot unavailable")
	ErrNoAvailableSlots  = errors.New("no available slots")
	ErrNoActiveProviders = errors.New("no active professionals")
)

const (
	createdByOnlineBooking = "booking"
)

type Slot struct {
	StartsAt string `json:"startsAt"`
	EndsAt   string `json:"endsAt"`
}

type ProfessionalSlots struct {
	ProfessionalID   string `json:"professionalId"`
	ProfessionalName string `json:"professionalName"`
	Slots            []Slot `json:"slots"`
}

type BookingInput struct {
	ServiceID      string `json:"serviceId"`
	ProfessionalID string `json:"professionalId"`
	StartsAt       string `json:"startsAt"`
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	Notes          string `json:"notes"`
}

type BookingResult struct {
	Appointment appointmentdomain.Appointment `json:"appointment"`
	Customer    customerdomain.Customer       `json:"customer"`
}

type UseCase struct {
	tenantRepo        tenantdomain.Repository
	configurationRepo configurationdomain.Repository
	serviceRepo       servicedomain.Repository
	professionalRepo  professionaldomain.Repository
	customerRepo      customerdomain.Repository
	appointments      *appappointment.UseCase
}

func NewUseCase(
	tenantRepo tenantdomain.Repository,
	configurationRepo configurationdomain.Repository,
	serviceRepo servicedomain.Repository,
	professionalRepo professionaldomain.Repository,
	customerRepo customerdomain.Repository,
	appointments *appappointment.UseCase,
) *UseCase {
	return &UseCase{
		tenantRepo:        tenantRepo,
		configurationRepo: configurationRepo,
		serviceRepo:       serviceRepo,
		professionalRepo:  professionalRepo,
		customerRepo:      customerRepo,
		appointments:      appointments,
	}
}

func (u *UseCase) ResolveTenantBySlug(slug string) (tenantdomain.Tenant, error) {
	tenant, err := u.tenantRepo.GetBySlug(strings.TrimSpace(slug))
	if err != nil {
		return tenantdomain.Tenant{}, ErrTenantNotFound
	}
	if !tenant.Active {
		return tenantdomain.Tenant{}, ErrTenantNotFound
	}
	return tenant, nil
}

func (u *UseCase) PublicInfo(slug string) (map[string]any, error) {
	tenant, err := u.ResolveTenantBySlug(slug)
	if err != nil {
		return nil, err
	}
	configuration, err := u.configurationRepo.Get(tenant.ID)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"tenant":        tenant,
		"configuration": configuration,
	}, nil
}

func (u *UseCase) ListActiveServices(slug string) ([]servicedomain.Service, error) {
	tenant, err := u.ResolveTenantBySlug(slug)
	if err != nil {
		return nil, err
	}
	services, err := u.serviceRepo.ListByTenant(tenant.ID)
	if err != nil {
		return nil, err
	}
	active := make([]servicedomain.Service, 0, len(services))
	for _, service := range services {
		if service.Active {
			active = append(active, service)
		}
	}
	return active, nil
}

func (u *UseCase) ListActiveProfessionals(slug string) ([]professionaldomain.Professional, error) {
	tenant, err := u.ResolveTenantBySlug(slug)
	if err != nil {
		return nil, err
	}
	professionals, err := u.professionalRepo.ListByTenant(tenant.ID)
	if err != nil {
		return nil, err
	}
	active := make([]professionaldomain.Professional, 0, len(professionals))
	for _, professional := range professionals {
		if professional.Active {
			active = append(active, professional)
		}
	}
	return active, nil
}

func (u *UseCase) Availability(slug string, date string, serviceID string, professionalID string) (map[string]any, error) {
	tenant, err := u.ResolveTenantBySlug(slug)
	if err != nil {
		return nil, err
	}
	day, err := time.Parse("2006-01-02", strings.TrimSpace(date))
	if err != nil {
		return nil, ErrInvalidBooking
	}
	duration, err := u.resolveDuration(tenant.ID, strings.TrimSpace(serviceID))
	if err != nil {
		return nil, err
	}
	configuration, err := u.configurationRepo.Get(tenant.ID)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(professionalID) != "" {
		slots, err := u.availableSlots(tenant.ID, day, duration, configuration, strings.TrimSpace(professionalID))
		if err != nil {
			return nil, err
		}
		return map[string]any{"professionalId": strings.TrimSpace(professionalID), "slots": slots}, nil
	}

	professionals, err := u.ListActiveProfessionals(slug)
	if err != nil {
		return nil, err
	}
	grouped := make([]ProfessionalSlots, 0, len(professionals))
	for _, professional := range professionals {
		slots, err := u.availableSlots(tenant.ID, day, duration, configuration, professional.ID)
		if err != nil {
			return nil, err
		}
		grouped = append(grouped, ProfessionalSlots{
			ProfessionalID:   professional.ID,
			ProfessionalName: professional.Name,
			Slots:            slots,
		})
	}
	return map[string]any{"professionals": grouped}, nil
}

func (u *UseCase) CreateBooking(slug string, input BookingInput) (BookingResult, error) {
	tenant, err := u.ResolveTenantBySlug(slug)
	if err != nil {
		return BookingResult{}, err
	}
	input.ServiceID = strings.TrimSpace(input.ServiceID)
	input.ProfessionalID = strings.TrimSpace(input.ProfessionalID)
	input.Name = strings.TrimSpace(input.Name)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Email = strings.TrimSpace(input.Email)
	input.Notes = strings.TrimSpace(input.Notes)

	if input.ProfessionalID == "" || input.Name == "" || input.Phone == "" {
		return BookingResult{}, ErrInvalidBooking
	}

	configuration, err := u.configurationRepo.Get(tenant.ID)
	if err != nil {
		return BookingResult{}, err
	}
	if !configuration.AllowOnlineBooking {
		return BookingResult{}, ErrBookingDisabled
	}

	duration, err := u.resolveDuration(tenant.ID, input.ServiceID)
	if err != nil {
		return BookingResult{}, err
	}

	startsAt, err := parseRFC3339(input.StartsAt)
	if err != nil {
		return BookingResult{}, ErrInvalidBooking
	}
	endsAt := startsAt.Add(time.Duration(duration) * time.Minute)

	if err := u.validateSlotInWorkingHours(startsAt, endsAt, configuration); err != nil {
		return BookingResult{}, err
	}
	if !startsAt.After(time.Now().Add(-5 * time.Minute)) {
		return BookingResult{}, ErrSlotUnavailable
	}

	customer, err := u.customerRepo.FindByPhone(tenant.ID, input.Phone)
	if err != nil {
		customer, err = u.customerRepo.Create(customerdomain.Customer{
			TenantID: tenant.ID,
			Name:     input.Name,
			Phone:    input.Phone,
			Email:    input.Email,
			Notes:    input.Notes,
		})
		if err != nil {
			return BookingResult{}, fmt.Errorf("create booking customer: %w", err)
		}
	}

	appointment, err := u.appointments.Create(appointmentdomain.Appointment{
		TenantID:       tenant.ID,
		ProfessionalID: input.ProfessionalID,
		ServiceID:      input.ServiceID,
		CustomerID:     customer.ID,
		StartsAt:       startsAt.Format(time.RFC3339),
		EndsAt:         endsAt.Format(time.RFC3339),
		Status:         appappointment.StatusScheduled,
		Notes:          input.Notes,
		CreatedBy:      createdByOnlineBooking,
	})
	if err != nil {
		if errors.Is(err, appappointment.ErrTimeConflict) {
			return BookingResult{}, ErrSlotUnavailable
		}
		return BookingResult{}, err
	}

	return BookingResult{Appointment: appointment, Customer: customer}, nil
}

func (u *UseCase) resolveDuration(tenantID string, serviceID string) (int, error) {
	if serviceID != "" {
		service, err := u.serviceRepo.GetByID(tenantID, serviceID)
		if err != nil {
			return 0, ErrInvalidBooking
		}
		if !service.Active || service.DurationMinutes <= 0 {
			return 0, ErrInvalidBooking
		}
		return service.DurationMinutes, nil
	}
	configuration, err := u.configurationRepo.Get(tenantID)
	if err != nil {
		return 0, ErrInvalidBooking
	}
	if configuration.DefaultDurationMinutes <= 0 {
		return 30, nil
	}
	return configuration.DefaultDurationMinutes, nil
}

func (u *UseCase) availableSlots(tenantID string, day time.Time, duration int, configuration configurationdomain.Configuration, professionalID string) ([]Slot, error) {
	if _, err := u.professionalRepo.GetByID(tenantID, professionalID); err != nil {
		return nil, ErrInvalidBooking
	}

	workingHours, err := workingHoursFor(configuration, int(day.Weekday()))
	if err != nil {
		return []Slot{}, nil
	}

	open, err := parseClock(day, workingHours.OpenTime)
	if err != nil {
		return []Slot{}, nil
	}
	closeTime, err := parseClock(day, workingHours.CloseTime)
	if err != nil {
		return []Slot{}, nil
	}

	interval := configuration.AppointmentIntervalMinutes
	if interval <= 0 {
		interval = 30
	}
	minStart := time.Now().Add(time.Duration(configuration.AdvanceNoticeHours) * time.Hour)

	from := localMidnight(day)
	to := localMidnight(day.AddDate(0, 0, 1))
	appointments, err := u.appointments.ListByRange(tenantID, from.Format(time.RFC3339), to.Format(time.RFC3339), professionalID, "", "")
	if err != nil {
		return nil, err
	}

	busy := make([]Slot, 0, len(appointments))
	for _, appointment := range appointments {
		if appointment.Status == appappointment.StatusCancelled || appointment.Status == appappointment.StatusNoShow {
			continue
		}
		busy = append(busy, Slot{StartsAt: appointment.StartsAt, EndsAt: appointment.EndsAt})
	}

	slots := make([]Slot, 0)
	for start := open; ; start = start.Add(time.Duration(interval) * time.Minute) {
		end := start.Add(time.Duration(duration) * time.Minute)
		if end.After(closeTime) {
			break
		}
		if start.Before(minStart) {
			continue
		}
		if overlapsAnyBusy(start, end, busy) {
			continue
		}
		slots = append(slots, Slot{
			StartsAt: start.UTC().Format(time.RFC3339),
			EndsAt:   end.UTC().Format(time.RFC3339),
		})
	}
	return slots, nil
}

func (u *UseCase) validateSlotInWorkingHours(startsAt time.Time, endsAt time.Time, configuration configurationdomain.Configuration) error {
	workingHours, err := workingHoursFor(configuration, int(startsAt.Local().Weekday()))
	if err != nil {
		return ErrSlotUnavailable
	}
	local := startsAt.Local()
	open, err := parseClock(local, workingHours.OpenTime)
	if err != nil {
		return ErrSlotUnavailable
	}
	closeTime, err := parseClock(local, workingHours.CloseTime)
	if err != nil {
		return ErrSlotUnavailable
	}
	if startsAt.Before(open) || endsAt.After(closeTime) {
		return ErrSlotUnavailable
	}
	return nil
}

func workingHoursFor(configuration configurationdomain.Configuration, weekday int) (configurationdomain.WorkingHours, error) {
	for _, hours := range configuration.WorkingHours {
		if hours.DayOfWeek == weekday {
			if !hours.Enabled {
				return configurationdomain.WorkingHours{}, fmt.Errorf("closed on weekday %d", weekday)
			}
			return hours, nil
		}
	}
	return configurationdomain.WorkingHours{}, fmt.Errorf("no working hours for weekday %d", weekday)
}

func parseClock(day time.Time, clock string) (time.Time, error) {
	value := strings.TrimSpace(clock)
	if value == "" {
		return time.Time{}, fmt.Errorf("empty clock")
	}
	return time.ParseInLocation("2006-01-02 15:04", day.Format("2006-01-02")+" "+value, time.Local)
}

func localMidnight(day time.Time) time.Time {
	year, month, date := day.Date()
	return time.Date(year, month, date, 0, 0, 0, 0, time.Local)
}

func overlapsAnyBusy(start time.Time, end time.Time, busy []Slot) bool {
	for _, interval := range busy {
		busyStart, err := parseRFC3339(interval.StartsAt)
		if err != nil {
			continue
		}
		busyEnd, err := parseRFC3339(interval.EndsAt)
		if err != nil {
			continue
		}
		if start.Before(busyEnd) && end.After(busyStart) {
			return true
		}
	}
	return false
}

func parseRFC3339(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, fmt.Errorf("empty time")
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		parsed, err = time.Parse("2006-01-02T15:04", value)
	}
	return parsed, err
}

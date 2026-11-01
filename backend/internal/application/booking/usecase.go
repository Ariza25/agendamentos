package booking

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"

	appappointment "erpagendamentos/backend/internal/application/appointment"
	appnotification "erpagendamentos/backend/internal/application/notification"
	appointmentdomain "erpagendamentos/backend/internal/domain/appointment"
	configurationdomain "erpagendamentos/backend/internal/domain/configuration"
	customerdomain "erpagendamentos/backend/internal/domain/customer"
	notificationdomain "erpagendamentos/backend/internal/domain/notification"
	professionaldomain "erpagendamentos/backend/internal/domain/professional"
	resourcedomain "erpagendamentos/backend/internal/domain/resource"
	servicedomain "erpagendamentos/backend/internal/domain/service"
	tenantdomain "erpagendamentos/backend/internal/domain/tenant"
)

var (
	ErrTenantNotFound      = errors.New("tenant not found")
	ErrBookingDisabled     = errors.New("online booking disabled")
	ErrInvalidBooking      = errors.New("invalid booking")
	ErrSlotUnavailable     = errors.New("slot unavailable")
	ErrNoAvailableSlots    = errors.New("no available slots")
	ErrNoActiveProviders   = errors.New("no active professionals")
	ErrInvalidCode         = errors.New("invalid management code")
	ErrAppointmentDone     = errors.New("appointment can no longer be managed online")
	ErrCancellationTooLate = errors.New("cancellation window has passed")
	ErrCustomerNotFound    = errors.New("customer not found for the given phone")
	ErrAppointmentAccess   = errors.New("appointment not found for this customer")
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
	ServiceID       string   `json:"serviceId"`
	ServiceIDs      []string `json:"serviceIds,omitempty"`
	ProfessionalID  string   `json:"professionalId"`
	StartsAt        string   `json:"startsAt"`
	Name            string   `json:"name"`
	Phone           string   `json:"phone"`
	Email           string   `json:"email"`
	Notes           string   `json:"notes"`
	RecurrenceRule  string   `json:"recurrenceRule,omitempty"`
	RecurrenceCount int      `json:"recurrenceCount,omitempty"`
	Timezone        string   `json:"timezone,omitempty"`
}

type BookingResult struct {
	Appointment    appointmentdomain.Appointment `json:"appointment"`
	Customer       customerdomain.Customer       `json:"customer"`
	ManagementCode string                        `json:"managementCode"`
}

type PortalCustomer struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type PortalAppointment struct {
	ID                 string `json:"id"`
	ServiceID          string `json:"serviceId"`
	ServiceName        string `json:"serviceName"`
	ProfessionalID     string `json:"professionalId"`
	ProfessionalName   string `json:"professionalName"`
	ResourceID         string `json:"resourceId,omitempty"`
	ResourceName       string `json:"resourceName,omitempty"`
	StartsAt           string `json:"startsAt"`
	EndsAt             string `json:"endsAt"`
	Status             string `json:"status"`
	PriceCents         int    `json:"priceCents"`
	CancellationReason string `json:"cancellationReason,omitempty"`
	Source             string `json:"source,omitempty"`
}

type PortalResult struct {
	Customer     PortalCustomer      `json:"customer"`
	Appointments []PortalAppointment `json:"appointments"`
}

type UseCase struct {
	tenantRepo        tenantdomain.Repository
	configurationRepo configurationdomain.Repository
	serviceRepo       servicedomain.Repository
	professionalRepo  professionaldomain.Repository
	resourceRepo      resourcedomain.Repository
	customerRepo      customerdomain.Repository
	appointments      *appappointment.UseCase
	notifications     *appnotification.UseCase
}

func NewUseCase(
	tenantRepo tenantdomain.Repository,
	configurationRepo configurationdomain.Repository,
	serviceRepo servicedomain.Repository,
	professionalRepo professionaldomain.Repository,
	resourceRepo resourcedomain.Repository,
	customerRepo customerdomain.Repository,
	appointments *appappointment.UseCase,
	notifications *appnotification.UseCase,
) *UseCase {
	return &UseCase{
		tenantRepo:        tenantRepo,
		configurationRepo: configurationRepo,
		serviceRepo:       serviceRepo,
		professionalRepo:  professionalRepo,
		resourceRepo:      resourceRepo,
		customerRepo:      customerRepo,
		appointments:      appointments,
		notifications:     notifications,
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
		slots, err := u.availableSlots(tenant.ID, day, duration, configuration, strings.TrimSpace(professionalID), strings.TrimSpace(serviceID))
		if err != nil {
			return nil, err
		}
		return map[string]any{"professionalId": strings.TrimSpace(professionalID), "slots": slots}, nil
	}

	allowed := map[string]bool{}
	if strings.TrimSpace(serviceID) != "" {
		if service, err := u.serviceRepo.GetByID(tenant.ID, strings.TrimSpace(serviceID)); err == nil {
			for _, professionalID := range service.ProfessionalIDs {
				allowed[professionalID] = true
			}
		}
	}

	professionals, err := u.ListActiveProfessionals(slug)
	if err != nil {
		return nil, err
	}
	grouped := make([]ProfessionalSlots, 0, len(professionals))
	for _, professional := range professionals {
		if len(allowed) > 0 && !allowed[professional.ID] {
			continue
		}
		slots, err := u.availableSlots(tenant.ID, day, duration, configuration, professional.ID, strings.TrimSpace(serviceID))
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
	if len(input.ServiceIDs) > 0 {
		comboDuration := 0
		for _, sID := range input.ServiceIDs {
			if d, err := u.resolveDuration(tenant.ID, sID); err == nil && d > 0 {
				comboDuration += d
			}
		}
		if comboDuration > 0 {
			duration = comboDuration
		}
		if input.ServiceID == "" && len(input.ServiceIDs) > 0 {
			input.ServiceID = input.ServiceIDs[0]
		}
	} else if err != nil {
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

	seriesID := ""
	if input.RecurrenceCount > 1 {
		seriesID = uuid.NewString()
	}

	managementCode := generateManagementCode()
	appointment, err := u.appointments.Create(appointmentdomain.Appointment{
		TenantID:           tenant.ID,
		ProfessionalID:     input.ProfessionalID,
		ServiceID:          input.ServiceID,
		ServiceIDs:         input.ServiceIDs,
		CustomerID:         customer.ID,
		StartsAt:           startsAt.Format(time.RFC3339),
		EndsAt:             endsAt.Format(time.RFC3339),
		Status:             appappointment.StatusScheduled,
		Notes:              input.Notes,
		CreatedBy:          createdByOnlineBooking,
		ManagementCodeHash: hashManagementCode(managementCode),
		SeriesID:           seriesID,
		RecurrenceRule:     input.RecurrenceRule,
		Timezone:           input.Timezone,
		Source:             "booking",
	})
	if err != nil {
		if errors.Is(err, appappointment.ErrTimeConflict) {
			return BookingResult{}, ErrSlotUnavailable
		}
		return BookingResult{}, err
	}

	// Create subsequent recurring appointments if requested
	if input.RecurrenceCount > 1 {
		for i := 1; i < input.RecurrenceCount; i++ {
			recStarts := startsAt.AddDate(0, 0, 7*i)
			recEnds := endsAt.AddDate(0, 0, 7*i)
			subCode := generateManagementCode()
			_, _ = u.appointments.Create(appointmentdomain.Appointment{
				TenantID:           tenant.ID,
				ProfessionalID:     input.ProfessionalID,
				ServiceID:          input.ServiceID,
				ServiceIDs:         input.ServiceIDs,
				CustomerID:         customer.ID,
				StartsAt:           recStarts.Format(time.RFC3339),
				EndsAt:             recEnds.Format(time.RFC3339),
				Status:             appappointment.StatusScheduled,
				Notes:              input.Notes,
				CreatedBy:          createdByOnlineBooking,
				ManagementCodeHash: hashManagementCode(subCode),
				SeriesID:           seriesID,
				RecurrenceRule:     input.RecurrenceRule,
				Timezone:           input.Timezone,
				Source:             "booking",
			})
		}
	}

	result := BookingResult{Appointment: appointment, Customer: customer, ManagementCode: managementCode}
	u.notifyCustomer(tenant, appointment, managementCode, string(notificationdomain.TypeConfirmation), appnotification.RenderConfirmation)
	return result, nil
}

func (u *UseCase) GetByManagementCode(slug string, code string) (appointmentdomain.Appointment, error) {
	tenant, err := u.ResolveTenantBySlug(slug)
	if err != nil {
		return appointmentdomain.Appointment{}, err
	}
	appointment, err := u.appointments.GetByManagementCode(tenant.ID, hashManagementCode(strings.ToUpper(strings.TrimSpace(code))))
	if err != nil {
		return appointmentdomain.Appointment{}, ErrInvalidCode
	}
	return appointment, nil
}

func (u *UseCase) CancelPublic(slug string, code string, reason string) (appointmentdomain.Appointment, error) {
	tenant, err := u.ResolveTenantBySlug(slug)
	if err != nil {
		return appointmentdomain.Appointment{}, err
	}
	appointment, err := u.appointments.GetByManagementCode(tenant.ID, hashManagementCode(strings.ToUpper(strings.TrimSpace(code))))
	if err != nil {
		return appointmentdomain.Appointment{}, ErrInvalidCode
	}
	if appointment.Status == appappointment.StatusCancelled || appointment.Status == appappointment.StatusCompleted || appointment.Status == appappointment.StatusNoShow {
		return appointmentdomain.Appointment{}, ErrAppointmentDone
	}
	configuration, err := u.configurationRepo.Get(tenant.ID)
	if err != nil {
		return appointmentdomain.Appointment{}, err
	}
	if configuration.CancellationPolicyHours > 0 {
		startsAt, err := parseRFC3339(appointment.StartsAt)
		if err != nil {
			return appointmentdomain.Appointment{}, ErrInvalidCode
		}
		if time.Now().After(startsAt.Add(-time.Duration(configuration.CancellationPolicyHours) * time.Hour)) {
			return appointmentdomain.Appointment{}, ErrCancellationTooLate
		}
	}
	input := appointment
	input.Status = appappointment.StatusCancelled
	input.CancellationReason = strings.TrimSpace(reason)
	updated, err := u.appointments.Update(input)
	if err != nil {
		return appointmentdomain.Appointment{}, err
	}
	u.notifyCustomer(tenant, updated, code, string(notificationdomain.TypeCancellation), appnotification.RenderCancellation)
	return updated, nil
}

func (u *UseCase) ReschedulePublic(slug string, code string, newStartsAt string) (appointmentdomain.Appointment, error) {
	tenant, err := u.ResolveTenantBySlug(slug)
	if err != nil {
		return appointmentdomain.Appointment{}, err
	}
	appointment, err := u.appointments.GetByManagementCode(tenant.ID, hashManagementCode(strings.ToUpper(strings.TrimSpace(code))))
	if err != nil {
		return appointmentdomain.Appointment{}, ErrInvalidCode
	}
	if appointment.Status == appappointment.StatusCancelled || appointment.Status == appappointment.StatusCompleted || appointment.Status == appappointment.StatusNoShow {
		return appointmentdomain.Appointment{}, ErrAppointmentDone
	}

	configuration, err := u.configurationRepo.Get(tenant.ID)
	if err != nil {
		return appointmentdomain.Appointment{}, err
	}
	if !configuration.AllowOnlineBooking {
		return appointmentdomain.Appointment{}, ErrBookingDisabled
	}

	duration := appointmentDurationMinutes(appointment)
	startsAt, err := parseRFC3339(newStartsAt)
	if err != nil {
		return appointmentdomain.Appointment{}, ErrInvalidBooking
	}
	endsAt := startsAt.Add(time.Duration(duration) * time.Minute)

	if err := u.validateSlotInWorkingHours(startsAt, endsAt, configuration); err != nil {
		return appointmentdomain.Appointment{}, ErrSlotUnavailable
	}
	if !startsAt.After(time.Now().Add(-5 * time.Minute)) {
		return appointmentdomain.Appointment{}, ErrSlotUnavailable
	}

	input := appointment
	input.StartsAt = startsAt.Format(time.RFC3339)
	input.EndsAt = endsAt.Format(time.RFC3339)
	updated, err := u.appointments.Update(input)
	if err != nil {
		if errors.Is(err, appappointment.ErrTimeConflict) {
			return appointmentdomain.Appointment{}, ErrSlotUnavailable
		}
		return appointmentdomain.Appointment{}, err
	}
	u.notifyCustomer(tenant, updated, code, string(notificationdomain.TypeReschedule), appnotification.RenderReschedule)
	return updated, nil
}

func (u *UseCase) PortalLookup(slug string, phone string) (PortalResult, error) {
	tenant, err := u.ResolveTenantBySlug(slug)
	if err != nil {
		return PortalResult{}, err
	}
	customer, err := u.customerRepo.FindByPhone(tenant.ID, strings.TrimSpace(phone))
	if err != nil {
		return PortalResult{}, ErrCustomerNotFound
	}

	appointments, err := u.appointments.ListByRange(
		tenant.ID,
		time.Now().AddDate(-2, 0, 0).Format(time.RFC3339),
		time.Now().AddDate(0, 0, 730).Format(time.RFC3339),
		"",
		customer.ID,
		"",
	)
	if err != nil {
		return PortalResult{}, err
	}

	portalAppointments := make([]PortalAppointment, 0, len(appointments))
	for _, item := range appointments {
		portalAppointments = append(portalAppointments, u.toPortalAppointment(tenant.ID, item))
	}
	return PortalResult{
		Customer: PortalCustomer{
			ID:    customer.ID,
			Name:  customer.Name,
			Phone: customer.Phone,
			Email: customer.Email,
		},
		Appointments: portalAppointments,
	}, nil
}

func (u *UseCase) PortalCancel(slug string, phone string, appointmentID string, reason string) (appointmentdomain.Appointment, error) {
	tenant, err := u.ResolveTenantBySlug(slug)
	if err != nil {
		return appointmentdomain.Appointment{}, err
	}
	customer, err := u.resolvePortalCustomer(tenant.ID, phone)
	if err != nil {
		return appointmentdomain.Appointment{}, err
	}
	appointment, err := u.appointments.GetByID(tenant.ID, strings.TrimSpace(appointmentID))
	if err != nil || appointment.CustomerID != customer.ID {
		return appointmentdomain.Appointment{}, ErrAppointmentAccess
	}
	if appointment.Status == appappointment.StatusCancelled || appointment.Status == appappointment.StatusCompleted || appointment.Status == appappointment.StatusNoShow {
		return appointmentdomain.Appointment{}, ErrAppointmentDone
	}
	configuration, err := u.configurationRepo.Get(tenant.ID)
	if err != nil {
		return appointmentdomain.Appointment{}, err
	}
	if configuration.CancellationPolicyHours > 0 {
		startsAt, err := parseRFC3339(appointment.StartsAt)
		if err != nil {
			return appointmentdomain.Appointment{}, ErrAppointmentAccess
		}
		if time.Now().After(startsAt.Add(-time.Duration(configuration.CancellationPolicyHours) * time.Hour)) {
			return appointmentdomain.Appointment{}, ErrCancellationTooLate
		}
	}
	input := appointment
	input.Status = appappointment.StatusCancelled
	input.CancellationReason = strings.TrimSpace(reason)
	updated, err := u.appointments.Update(input)
	if err != nil {
		return appointmentdomain.Appointment{}, err
	}
	u.notifyCustomer(tenant, updated, "", string(notificationdomain.TypeCancellation), appnotification.RenderCancellation)
	return updated, nil
}

func (u *UseCase) PortalReschedule(slug string, phone string, appointmentID string, newStartsAt string) (appointmentdomain.Appointment, error) {
	tenant, err := u.ResolveTenantBySlug(slug)
	if err != nil {
		return appointmentdomain.Appointment{}, err
	}
	customer, err := u.resolvePortalCustomer(tenant.ID, phone)
	if err != nil {
		return appointmentdomain.Appointment{}, err
	}
	appointment, err := u.appointments.GetByID(tenant.ID, strings.TrimSpace(appointmentID))
	if err != nil || appointment.CustomerID != customer.ID {
		return appointmentdomain.Appointment{}, ErrAppointmentAccess
	}
	if appointment.Status == appappointment.StatusCancelled || appointment.Status == appappointment.StatusCompleted || appointment.Status == appappointment.StatusNoShow {
		return appointmentdomain.Appointment{}, ErrAppointmentDone
	}

	configuration, err := u.configurationRepo.Get(tenant.ID)
	if err != nil {
		return appointmentdomain.Appointment{}, err
	}
	if !configuration.AllowOnlineBooking {
		return appointmentdomain.Appointment{}, ErrBookingDisabled
	}

	duration := appointmentDurationMinutes(appointment)
	startsAt, err := parseRFC3339(newStartsAt)
	if err != nil {
		return appointmentdomain.Appointment{}, ErrInvalidBooking
	}
	endsAt := startsAt.Add(time.Duration(duration) * time.Minute)

	if err := u.validateSlotInWorkingHours(startsAt, endsAt, configuration); err != nil {
		return appointmentdomain.Appointment{}, ErrSlotUnavailable
	}
	if !startsAt.After(time.Now().Add(-5 * time.Minute)) {
		return appointmentdomain.Appointment{}, ErrSlotUnavailable
	}

	input := appointment
	input.StartsAt = startsAt.Format(time.RFC3339)
	input.EndsAt = endsAt.Format(time.RFC3339)
	updated, err := u.appointments.Update(input)
	if err != nil {
		if errors.Is(err, appappointment.ErrTimeConflict) {
			return appointmentdomain.Appointment{}, ErrSlotUnavailable
		}
		return appointmentdomain.Appointment{}, err
	}
	u.notifyCustomer(tenant, updated, "", string(notificationdomain.TypeReschedule), appnotification.RenderReschedule)
	return updated, nil
}

func (u *UseCase) resolvePortalCustomer(tenantID string, phone string) (customerdomain.Customer, error) {
	customer, err := u.customerRepo.FindByPhone(tenantID, strings.TrimSpace(phone))
	if err != nil {
		return customerdomain.Customer{}, ErrCustomerNotFound
	}
	return customer, nil
}

func (u *UseCase) notifyCustomer(
	tenant tenantdomain.Tenant,
	appointment appointmentdomain.Appointment,
	code string,
	notificationType string,
	build func(appnotification.TemplateData) (string, string),
) {
	if u == nil || u.notifications == nil {
		return
	}
	configuration, err := u.configurationRepo.Get(tenant.ID)
	if err != nil || !configuration.NotificationsEnabled {
		return
	}
	customer, err := u.customerRepo.GetByID(tenant.ID, appointment.CustomerID)
	if err != nil {
		return
	}

	channel := strings.TrimSpace(configuration.NotificationChannel)
	if channel == "" {
		channel = "whatsapp"
	}
	recipient := strings.TrimSpace(customer.Phone)
	if channel == "email" {
		recipient = strings.TrimSpace(customer.Email)
	}
	if recipient == "" {
		return
	}

	serviceName := "Agendamento"
	if appointment.ServiceID != "" {
		if service, err := u.serviceRepo.GetByID(tenant.ID, appointment.ServiceID); err == nil {
			serviceName = service.Name
		}
	}
	professionalName := "Equipe"
	if appointment.ProfessionalID != "" {
		if professional, err := u.professionalRepo.GetByID(tenant.ID, appointment.ProfessionalID); err == nil {
			professionalName = professional.Name
		}
	}

	data := appnotification.TemplateData{
		EstablishmentName: tenant.Name,
		CustomerName:      customer.Name,
		ServiceName:       serviceName,
		ProfessionalName:  professionalName,
		StartsAt:          formatAppointmentForTemplate(appointment.StartsAt),
		ManagementCode:    strings.TrimSpace(code),
		ManageURL:         appnotification.PublicBookingOrigin() + "/#/manage/" + tenant.Slug,
		PortalURL:         appnotification.PublicBookingOrigin() + "/#/portal/" + tenant.Slug,
	}
	subject, body := build(data)
	if _, err := u.notifications.Send(notificationdomain.Log{
		TenantID:      tenant.ID,
		AppointmentID: appointment.ID,
		Type:          notificationType,
		Channel:       channel,
		Recipient:     recipient,
		Subject:       subject,
		Body:          body,
	}); err != nil {
		log.Printf(`{"level":"WARN","message":"notification send failed","type":"%s","appointment":"%s","error":"%v"}`, notificationType, appointment.ID, err)
	}
}

func formatAppointmentForTemplate(value string) string {
	parsed, err := parseRFC3339(value)
	if err != nil {
		return strings.TrimSpace(value)
	}
	return parsed.Local().Format("02/01/2006 às 15:04")
}

func (u *UseCase) toPortalAppointment(tenantID string, item appointmentdomain.Appointment) PortalAppointment {
	portal := PortalAppointment{
		ID:                 item.ID,
		ServiceID:          item.ServiceID,
		ProfessionalID:     item.ProfessionalID,
		ResourceID:         item.ResourceID,
		StartsAt:           item.StartsAt,
		EndsAt:             item.EndsAt,
		Status:             item.Status,
		CancellationReason: item.CancellationReason,
		Source:             item.Source,
	}
	if service, err := u.serviceRepo.GetByID(tenantID, item.ServiceID); err == nil {
		portal.ServiceName = service.Name
		portal.PriceCents = service.PriceCents
	}
	if professional, err := u.professionalRepo.GetByID(tenantID, item.ProfessionalID); err == nil {
		portal.ProfessionalName = professional.Name
	}
	if item.ResourceID != "" && u.resourceRepo != nil {
		if resource, err := u.resourceRepo.GetByID(tenantID, item.ResourceID); err == nil {
			portal.ResourceName = resource.Name
		}
	}
	return portal
}

func appointmentDurationMinutes(item appointmentdomain.Appointment) int {
	startsAt, startErr := parseRFC3339(item.StartsAt)
	endsAt, endErr := parseRFC3339(item.EndsAt)
	if startErr != nil || endErr != nil || !endsAt.After(startsAt) {
		return 30
	}
	minutes := int(endsAt.Sub(startsAt).Minutes())
	if minutes <= 0 {
		return 30
	}
	return minutes
}

func generateManagementCode() string {
	const letters = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	var builder strings.Builder
	for index := 0; index < 4; index++ {
		builder.WriteByte(letters[rand.Intn(len(letters))])
	}
	code := builder.String()
	return fmt.Sprintf("%s-%s-%s", code[:2], code[2:3], code[3:])
}

func hashManagementCode(code string) string {
	sum := sha256.Sum256([]byte(strings.ToUpper(strings.TrimSpace(code))))
	return hex.EncodeToString(sum[:])
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

func (u *UseCase) availableSlots(tenantID string, day time.Time, duration int, configuration configurationdomain.Configuration, professionalID string, serviceID string) ([]Slot, error) {
	if _, err := u.professionalRepo.GetByID(tenantID, professionalID); err != nil {
		return nil, ErrInvalidBooking
	}

	maxCapacity := 1
	if serviceID != "" {
		if svc, err := u.serviceRepo.GetByID(tenantID, serviceID); err == nil && svc.MaxCapacity > 1 {
			maxCapacity = svc.MaxCapacity
		}
	}

	if isClosedDay(configuration, day) {
		return []Slot{}, nil
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
	buffer := configuration.BufferMinutes
	if buffer < 0 {
		buffer = 0
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
		busyStart, err := parseRFC3339(appointment.StartsAt)
		if err != nil {
			continue
		}
		busyEnd, err := parseRFC3339(appointment.EndsAt)
		if err != nil {
			continue
		}
		busy = append(busy, Slot{
			StartsAt: busyStart.Add(-time.Duration(buffer) * time.Minute).UTC().Format(time.RFC3339),
			EndsAt:   busyEnd.Add(time.Duration(buffer) * time.Minute).UTC().Format(time.RFC3339),
		})
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
		if countOverlappingSlots(start, end, busy) >= maxCapacity {
			continue
		}
		slots = append(slots, Slot{
			StartsAt: start.UTC().Format(time.RFC3339),
			EndsAt:   end.UTC().Format(time.RFC3339),
		})
	}
	return slots, nil
}

func countOverlappingSlots(start time.Time, end time.Time, busy []Slot) int {
	count := 0
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
			count++
		}
	}
	return count
}

func (u *UseCase) validateSlotInWorkingHours(startsAt time.Time, endsAt time.Time, configuration configurationdomain.Configuration) error {
	if isClosedDay(configuration, startsAt) {
		return ErrSlotUnavailable
	}
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
	if local.Before(open) || endsAt.Local().After(closeTime) {
		return ErrSlotUnavailable
	}
	return nil
}

func isClosedDay(configuration configurationdomain.Configuration, day time.Time) bool {
	target := day.Format("2006-01-02")
	for _, closed := range configuration.ClosedDays {
		if strings.TrimSpace(closed.Date) == target {
			return true
		}
	}
	return false
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
		parsed, err = time.ParseInLocation("2006-01-02T15:04", value, time.Local)
	}
	return parsed, err
}

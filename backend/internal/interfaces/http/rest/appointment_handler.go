package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	appappointment "erpagendamentos/backend/internal/application/appointment"
	appcommission "erpagendamentos/backend/internal/application/commission"
	appfinance "erpagendamentos/backend/internal/application/finance"
	appprofessional "erpagendamentos/backend/internal/application/professional"
	appreview "erpagendamentos/backend/internal/application/review"
	appstock "erpagendamentos/backend/internal/application/stock"
	appointmentdomain "erpagendamentos/backend/internal/domain/appointment"
	servicedomain "erpagendamentos/backend/internal/domain/service"
)

type AppointmentHandler struct {
	useCase             *appappointment.UseCase
	reviewUseCase       *appreview.UseCase
	stockUseCase        *appstock.UseCase
	commissionUseCase   *appcommission.UseCase
	financeUseCase      *appfinance.UseCase
	professionalUseCase *appprofessional.UseCase
	serviceRepo         servicedomain.Repository
	eventBroker         *EventBroker
}

func NewAppointmentHandler(
	useCase *appappointment.UseCase,
	reviewUseCase *appreview.UseCase,
	stockUseCase *appstock.UseCase,
	commissionUseCase *appcommission.UseCase,
	financeUseCase *appfinance.UseCase,
	professionalUseCase *appprofessional.UseCase,
	serviceRepo servicedomain.Repository,
	eventBroker *EventBroker,
) *AppointmentHandler {
	return &AppointmentHandler{
		useCase:             useCase,
		reviewUseCase:       reviewUseCase,
		stockUseCase:        stockUseCase,
		commissionUseCase:   commissionUseCase,
		financeUseCase:      financeUseCase,
		professionalUseCase: professionalUseCase,
		serviceRepo:         serviceRepo,
		eventBroker:         eventBroker,
	}
}

func (h *AppointmentHandler) onStatusCompleted(appt appointmentdomain.Appointment) {
	if h.reviewUseCase != nil && appt.CustomerID != "" {
		_, _ = h.reviewUseCase.CreateInvitation(appt.TenantID, appt.ID, appt.CustomerID)
	}
	if h.stockUseCase != nil && appt.ServiceID != "" {
		_ = h.stockUseCase.DeductSuppliesForAppointment(appt.TenantID, appt.ID, appt.ServiceID)
	}
	var servicePrice int
	var serviceName string
	if h.serviceRepo != nil && appt.ServiceID != "" {
		if svc, err := h.serviceRepo.GetByID(appt.TenantID, appt.ServiceID); err == nil {
			servicePrice = svc.PriceCents
			serviceName = svc.Name
		}
	}
	var profName string
	if h.professionalUseCase != nil && appt.ProfessionalID != "" {
		if prof, err := h.professionalUseCase.GetByID(appt.TenantID, appt.ProfessionalID); err == nil {
			profName = prof.Name
		}
	}
	if h.commissionUseCase != nil && appt.ProfessionalID != "" {
		_, _ = h.commissionUseCase.CalculateAndRecord(appt.TenantID, appt.ID, appt.ProfessionalID, profName, appt.ServiceID, serviceName, servicePrice)
	}
	if h.financeUseCase != nil && servicePrice > 0 {
		desc := "Atendimento: " + serviceName
		if desc == "Atendimento: " {
			desc = "Atendimento Concluído"
		}
		_ = h.financeUseCase.RecordAppointmentIncome(appt.TenantID, appt.ID, servicePrice, desc)
	}
}

func (h *AppointmentHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	appointments, err := h.useCase.ListByRange(
		resolveRequestTenantID(r, nil),
		query.Get("from"),
		query.Get("to"),
		query.Get("professionalId"),
		query.Get("customerId"),
		query.Get("status"),
	)
	if err != nil {
		writeInternalError(w, "Could not load the appointments.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"appointments": appointments})
}

func (h *AppointmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	appointment, err := h.useCase.GetByID(resolveRequestTenantID(r, nil), r.PathValue("appointmentID"))
	if err != nil {
		writeNotFoundError(w, "not_found", "Appointment not found.")
		return
	}
	writeJSON(w, http.StatusOK, appointment)
}

func (h *AppointmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input appointmentdomain.Appointment
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the appointment data.")
		return
	}
	input.TenantID = resolveRequestTenantID(r, nil)
	appointment, err := h.useCase.Create(input)
	if err != nil {
		writeAppointmentError(w, err)
		return
	}
	if h.eventBroker != nil {
		h.eventBroker.Publish(appointment.TenantID, "appointment.created", appointment)
	}
	writeJSON(w, http.StatusCreated, appointment)
}

func (h *AppointmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	var input appointmentdomain.Appointment
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the appointment data.")
		return
	}
	input.TenantID = resolveRequestTenantID(r, nil)
	input.ID = r.PathValue("appointmentID")
	appointment, err := h.useCase.Update(input)
	if err != nil {
		writeAppointmentError(w, err)
		return
	}
	if h.eventBroker != nil {
		h.eventBroker.Publish(appointment.TenantID, "appointment.updated", appointment)
	}
	writeJSON(w, http.StatusOK, appointment)
}

func (h *AppointmentHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Provide the new appointment status.")
		return
	}
	appointment, err := h.useCase.UpdateStatus(
		resolveRequestTenantID(r, nil),
		r.PathValue("appointmentID"),
		input.Status,
		resolveActorTenantID(r, nil),
	)
	if err != nil {
		writeAppointmentError(w, err)
		return
	}
	if appointment.Status == appappointment.StatusCompleted {
		h.onStatusCompleted(appointment)
	}
	if h.eventBroker != nil {
		h.eventBroker.Publish(appointment.TenantID, "appointment."+appointment.Status, appointment)
	}
	writeJSON(w, http.StatusOK, appointment)
}

func (h *AppointmentHandler) Summary(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	summary, err := h.useCase.Summary(resolveRequestTenantID(r, nil), date)
	if err != nil {
		writeBadRequestError(w, "invalid_input", "Provide the summary date (format YYYY-MM-DD).")
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

func (h *AppointmentHandler) Analytics(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	analytics, err := h.useCase.Analytics(
		resolveRequestTenantID(r, nil),
		query.Get("from"),
		query.Get("to"),
	)
	if err != nil {
		writeInternalError(w, "Could not load the analytics.")
		return
	}
	writeJSON(w, http.StatusOK, analytics)
}

func (h *AppointmentHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
	appointment, err := h.useCase.CheckIn(
		resolveRequestTenantID(r, nil),
		r.PathValue("appointmentID"),
		resolveActorTenantID(r, nil),
	)
	if err != nil {
		writeAppointmentError(w, err)
		return
	}
	if h.eventBroker != nil {
		h.eventBroker.Publish(appointment.TenantID, "appointment.checked_in", appointment)
	}
	writeJSON(w, http.StatusOK, appointment)
}

func (h *AppointmentHandler) MarkNoShow(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FeeCents int `json:"feeCents"`
	}
	_ = json.NewDecoder(r.Body).Decode(&input)
	appointment, err := h.useCase.MarkNoShow(
		resolveRequestTenantID(r, nil),
		r.PathValue("appointmentID"),
		resolveActorTenantID(r, nil),
		input.FeeCents,
	)
	if err != nil {
		writeAppointmentError(w, err)
		return
	}
	if h.eventBroker != nil {
		h.eventBroker.Publish(appointment.TenantID, "appointment.no_show", appointment)
	}
	writeJSON(w, http.StatusOK, appointment)
}

func (h *AppointmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	id := r.PathValue("appointmentID")
	if err := h.useCase.Delete(tenantID, id); err != nil {
		writeAppointmentError(w, err)
		return
	}
	if h.eventBroker != nil {
		h.eventBroker.Publish(tenantID, "appointment.deleted", map[string]string{"id": id})
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeAppointmentError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, appappointment.ErrAppointmentNotFound):
		writeNotFoundError(w, "not_found", "Appointment not found.")
	case errors.Is(err, appappointment.ErrInvalidAppointment):
		writeBadRequestError(w, "invalid_input", "Invalid appointment data.")
	case errors.Is(err, appappointment.ErrTimeConflict):
		writeConflictError(w, "time_conflict", "An appointment already exists for this time.")
	case errors.Is(err, appappointment.ErrInvalidStatusChange):
		writeConflictError(w, "invalid_status_change", "Invalid status transition for this appointment.")
	default:
		writeInternalError(w, "Could not process the appointment.")
	}
}

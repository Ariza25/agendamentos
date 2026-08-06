package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	appappointment "erpagendamentos/backend/internal/application/appointment"
	appointmentdomain "erpagendamentos/backend/internal/domain/appointment"
)

type AppointmentHandler struct {
	useCase *appappointment.UseCase
}

func NewAppointmentHandler(useCase *appappointment.UseCase) *AppointmentHandler {
	return &AppointmentHandler{useCase: useCase}
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
	writeJSON(w, http.StatusOK, appointment)
}

func (h *AppointmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.useCase.Delete(resolveRequestTenantID(r, nil), r.PathValue("appointmentID")); err != nil {
		writeAppointmentError(w, err)
		return
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

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
		writeInternalError(w, "Nao foi possivel carregar os agendamentos.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"appointments": appointments})
}

func (h *AppointmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	appointment, err := h.useCase.GetByID(resolveRequestTenantID(r, nil), r.PathValue("appointmentID"))
	if err != nil {
		writeNotFoundError(w, "not_found", "Agendamento nao encontrado.")
		return
	}
	writeJSON(w, http.StatusOK, appointment)
}

func (h *AppointmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input appointmentdomain.Appointment
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Nao foi possivel ler os dados do agendamento.")
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
		writeInvalidPayload(w, "Nao foi possivel ler os dados do agendamento.")
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
		writeInvalidPayload(w, "Informe o novo status do agendamento.")
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
		writeBadRequestError(w, "invalid_input", "Informe a data do resumo (formato AAAA-MM-DD).")
		return
	}
	writeJSON(w, http.StatusOK, summary)
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
		writeNotFoundError(w, "not_found", "Agendamento nao encontrado.")
	case errors.Is(err, appappointment.ErrInvalidAppointment):
		writeBadRequestError(w, "invalid_input", "Dados do agendamento invalidos.")
	case errors.Is(err, appappointment.ErrTimeConflict):
		writeConflictError(w, "time_conflict", "Ja existe um agendamento para este horario.")
	case errors.Is(err, appappointment.ErrInvalidStatusChange):
		writeConflictError(w, "invalid_status_change", "Transicao de status invalida para este agendamento.")
	default:
		writeInternalError(w, "Nao foi possivel processar o agendamento.")
	}
}

package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	appbooking "erpagendamentos/backend/internal/application/booking"
)

type BookingHandler struct {
	useCase *appbooking.UseCase
}

func NewBookingHandler(useCase *appbooking.UseCase) *BookingHandler {
	return &BookingHandler{useCase: useCase}
}

func (h *BookingHandler) Info(w http.ResponseWriter, r *http.Request) {
	info, err := h.useCase.PublicInfo(r.PathValue("slug"))
	if err != nil {
		writeBookingError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, info)
}

func (h *BookingHandler) Services(w http.ResponseWriter, r *http.Request) {
	services, err := h.useCase.ListActiveServices(r.PathValue("slug"))
	if err != nil {
		writeBookingError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"services": services})
}

func (h *BookingHandler) Professionals(w http.ResponseWriter, r *http.Request) {
	professionals, err := h.useCase.ListActiveProfessionals(r.PathValue("slug"))
	if err != nil {
		writeBookingError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"professionals": professionals})
}

func (h *BookingHandler) Availability(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	availability, err := h.useCase.Availability(
		r.PathValue("slug"),
		query.Get("date"),
		query.Get("serviceId"),
		query.Get("professionalId"),
	)
	if err != nil {
		writeBookingError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, availability)
}

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input appbooking.BookingInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Nao foi possivel ler os dados do agendamento.")
		return
	}
	result, err := h.useCase.CreateBooking(r.PathValue("slug"), input)
	if err != nil {
		writeBookingError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func writeBookingError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, appbooking.ErrTenantNotFound):
		writeNotFoundError(w, "not_found", "Estabelecimento nao encontrado.")
	case errors.Is(err, appbooking.ErrBookingDisabled):
		writeAPIError(w, http.StatusForbidden, "booking_disabled", "O agendamento online esta desativado para este estabelecimento.")
	case errors.Is(err, appbooking.ErrSlotUnavailable):
		writeConflictError(w, "time_conflict", "Este horario nao esta mais disponivel.")
	case errors.Is(err, appbooking.ErrInvalidBooking):
		writeBadRequestError(w, "invalid_input", "Dados do agendamento invalidos.")
	case errors.Is(err, appbooking.ErrNoAvailableSlots):
		writeBadRequestError(w, "no_available_slots", "Nao ha horarios disponiveis para esta data.")
	default:
		writeInternalError(w, "Nao foi possivel processar o agendamento online.")
	}
}

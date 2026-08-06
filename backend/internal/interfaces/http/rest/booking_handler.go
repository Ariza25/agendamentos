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
		writeInvalidPayload(w, "Could not read the appointment data.")
		return
	}
	result, err := h.useCase.CreateBooking(r.PathValue("slug"), input)
	if err != nil {
		writeBookingError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func (h *BookingHandler) GetByManagementCode(w http.ResponseWriter, r *http.Request) {
	appointment, err := h.useCase.GetByManagementCode(r.PathValue("slug"), r.PathValue("code"))
	if err != nil {
		writeBookingError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, appointment)
}

func (h *BookingHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Reason string `json:"reason"`
	}
	_ = json.NewDecoder(r.Body).Decode(&input)
	appointment, err := h.useCase.CancelPublic(r.PathValue("slug"), r.PathValue("code"), input.Reason)
	if err != nil {
		writeBookingError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, appointment)
}

func (h *BookingHandler) Reschedule(w http.ResponseWriter, r *http.Request) {
	var input struct {
		StartsAt string `json:"startsAt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Provide the new appointment time.")
		return
	}
	appointment, err := h.useCase.ReschedulePublic(r.PathValue("slug"), r.PathValue("code"), input.StartsAt)
	if err != nil {
		writeBookingError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, appointment)
}

func writeBookingError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, appbooking.ErrTenantNotFound):
		writeNotFoundError(w, "not_found", "Establishment not found.")
	case errors.Is(err, appbooking.ErrBookingDisabled):
		writeAPIError(w, http.StatusForbidden, "booking_disabled", "Online booking is disabled for this establishment.")
	case errors.Is(err, appbooking.ErrSlotUnavailable):
		writeConflictError(w, "time_conflict", "This time slot is no longer available.")
	case errors.Is(err, appbooking.ErrInvalidBooking):
		writeBadRequestError(w, "invalid_input", "Invalid appointment data.")
	case errors.Is(err, appbooking.ErrNoAvailableSlots):
		writeBadRequestError(w, "no_available_slots", "No slots available for this date.")
	case errors.Is(err, appbooking.ErrInvalidCode):
		writeNotFoundError(w, "not_found", "Appointment not found. Check the code and try again.")
	case errors.Is(err, appbooking.ErrAppointmentDone):
		writeConflictError(w, "appointment_done", "This appointment can no longer be managed online.")
	case errors.Is(err, appbooking.ErrCancellationTooLate):
		writeConflictError(w, "cancellation_window", "The cancellation window for this appointment has passed.")
	default:
		writeInternalError(w, "Could not process the online appointment.")
	}
}

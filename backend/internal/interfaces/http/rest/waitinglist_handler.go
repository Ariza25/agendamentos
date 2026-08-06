package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	appbooking "erpagendamentos/backend/internal/application/booking"
	appwaitinglist "erpagendamentos/backend/internal/application/waitinglist"
	waitinglistdomain "erpagendamentos/backend/internal/domain/waitinglist"
)

type WaitingListHandler struct {
	useCase        *appwaitinglist.UseCase
	bookingUseCase *appbooking.UseCase
}

func NewWaitingListHandler(useCase *appwaitinglist.UseCase, bookingUseCase *appbooking.UseCase) *WaitingListHandler {
	return &WaitingListHandler{useCase: useCase, bookingUseCase: bookingUseCase}
}

func (h *WaitingListHandler) List(w http.ResponseWriter, r *http.Request) {
	entries, err := h.useCase.ListByTenant(resolveRequestTenantID(r, nil), r.URL.Query().Get("status"))
	if err != nil {
		writeInternalError(w, "Could not load the waiting list.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"entries": entries})
}

func (h *WaitingListHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Provide the new waiting list status.")
		return
	}
	entry, err := h.useCase.UpdateStatus(
		resolveRequestTenantID(r, nil),
		r.PathValue("entryID"),
		input.Status,
	)
	if err != nil {
		writeWaitingListError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, entry)
}

func (h *WaitingListHandler) PublicSignup(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ServiceID      string `json:"serviceId"`
		ProfessionalID string `json:"professionalId"`
		Name           string `json:"name"`
		Phone          string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the waiting list data.")
		return
	}
	if h.bookingUseCase == nil {
		writeInternalError(w, "Could not process the waiting list.")
		return
	}
	tenant, err := h.bookingUseCase.ResolveTenantBySlug(r.PathValue("slug"))
	if err != nil {
		writeWaitingListError(w, err)
		return
	}
	entry, err := h.useCase.Create(waitinglistdomain.Entry{
		TenantID:       tenant.ID,
		ServiceID:      input.ServiceID,
		ProfessionalID: input.ProfessionalID,
		Name:           input.Name,
		Phone:          input.Phone,
	})
	if err != nil {
		writeWaitingListError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, entry)
}

func writeWaitingListError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, appwaitinglist.ErrWaitingListNotFound):
		writeNotFoundError(w, "not_found", "Waiting list entry not found.")
	case errors.Is(err, appwaitinglist.ErrInvalidWaitingList):
		writeBadRequestError(w, "invalid_input", "Invalid waiting list data.")
	case errors.Is(err, appbooking.ErrTenantNotFound):
		writeNotFoundError(w, "not_found", "Establishment not found.")
	default:
		writeInternalError(w, "Could not process the waiting list.")
	}
}

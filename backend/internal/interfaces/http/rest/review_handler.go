package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	appbooking "erpagendamentos/backend/internal/application/booking"
	appreview "erpagendamentos/backend/internal/application/review"
)

type ReviewHandler struct {
	useCase        *appreview.UseCase
	bookingUseCase *appbooking.UseCase
}

func NewReviewHandler(useCase *appreview.UseCase, bookingUseCase *appbooking.UseCase) *ReviewHandler {
	return &ReviewHandler{useCase: useCase, bookingUseCase: bookingUseCase}
}

func (h *ReviewHandler) List(w http.ResponseWriter, r *http.Request) {
	reviews, err := h.useCase.ListByTenant(resolveRequestTenantID(r, nil), r.URL.Query().Get("status"))
	if err != nil {
		writeInternalError(w, "Could not load the reviews.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"reviews": reviews})
}

func (h *ReviewHandler) GetPublic(w http.ResponseWriter, r *http.Request) {
	if h.bookingUseCase == nil {
		writeInternalError(w, "Could not load the review.")
		return
	}
	tenant, err := h.bookingUseCase.ResolveTenantBySlug(r.PathValue("slug"))
	if err != nil {
		writeReviewError(w, err)
		return
	}
	review, err := h.useCase.GetByID(tenant.ID, r.PathValue("reviewID"))
	if err != nil {
		writeReviewError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, review)
}

func (h *ReviewHandler) SubmitPublic(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the review data.")
		return
	}
	if h.bookingUseCase == nil {
		writeInternalError(w, "Could not submit the review.")
		return
	}
	tenant, err := h.bookingUseCase.ResolveTenantBySlug(r.PathValue("slug"))
	if err != nil {
		writeReviewError(w, err)
		return
	}
	review, err := h.useCase.Submit(tenant.ID, r.PathValue("reviewID"), input.Rating, input.Comment)
	if err != nil {
		writeReviewError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, review)
}

func writeReviewError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, appbooking.ErrTenantNotFound):
		writeNotFoundError(w, "not_found", "Establishment not found.")
	case errors.Is(err, appreview.ErrReviewNotFound):
		writeNotFoundError(w, "not_found", "Review not found.")
	case errors.Is(err, appreview.ErrInvalidReview):
		writeBadRequestError(w, "invalid_input", "Invalid review data.")
	case errors.Is(err, appreview.ErrReviewSubmitted):
		writeConflictError(w, "review_submitted", "This review was already submitted.")
	default:
		writeInternalError(w, "Could not process the review.")
	}
}

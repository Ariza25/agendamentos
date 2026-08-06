package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	apppayment "erpagendamentos/backend/internal/application/payment"
	paymentdomain "erpagendamentos/backend/internal/domain/payment"
)

type PaymentHandler struct {
	useCase *apppayment.UseCase
}

func NewPaymentHandler(useCase *apppayment.UseCase) *PaymentHandler {
	return &PaymentHandler{useCase: useCase}
}

func (h *PaymentHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	payments, err := h.useCase.ListByTenant(resolveRequestTenantID(r, nil), limit)
	if err != nil {
		writeInternalError(w, "Could not load the payments.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"payments": payments})
}

func (h *PaymentHandler) Get(w http.ResponseWriter, r *http.Request) {
	payment, err := h.useCase.GetByID(resolveRequestTenantID(r, nil), r.PathValue("paymentID"))
	if err != nil {
		writePaymentError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, payment)
}

func (h *PaymentHandler) MarkPaid(w http.ResponseWriter, r *http.Request) {
	payment, err := h.useCase.MarkPaid(resolveRequestTenantID(r, nil), r.PathValue("paymentID"))
	if err != nil {
		writePaymentError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, payment)
}

func (h *PaymentHandler) Refund(w http.ResponseWriter, r *http.Request) {
	payment, err := h.useCase.Refund(resolveRequestTenantID(r, nil), r.PathValue("paymentID"))
	if err != nil {
		writePaymentError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, payment)
}

func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		AppointmentID string `json:"appointmentId"`
		AmountCents   int    `json:"amountCents"`
		Method        string `json:"method"`
		Description   string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the payment data.")
		return
	}
	payment, err := h.useCase.CreatePaymentIntent(resolveRequestTenantID(r, nil), paymentdomain.Payment{
		TenantID:      resolveRequestTenantID(r, nil),
		AppointmentID: input.AppointmentID,
		AmountCents:   input.AmountCents,
		Method:        input.Method,
		Description:   input.Description,
	})
	if err != nil {
		writePaymentError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, payment)
}

func writePaymentError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apppayment.ErrPaymentNotFound):
		writeNotFoundError(w, "not_found", "Payment not found.")
	case errors.Is(err, apppayment.ErrInvalidPayment):
		writeBadRequestError(w, "invalid_input", "Invalid payment data.")
	default:
		writeInternalError(w, "Could not process the payment.")
	}
}

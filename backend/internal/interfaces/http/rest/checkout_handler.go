package rest

import (
	"encoding/json"
	"net/http"

	appcheckout "erpagendamentos/backend/internal/application/checkout"
)

type CheckoutHandler struct {
	useCase     *appcheckout.UseCase
	eventBroker *EventBroker
}

func NewCheckoutHandler(useCase *appcheckout.UseCase, eventBroker *EventBroker) *CheckoutHandler {
	return &CheckoutHandler{
		useCase:     useCase,
		eventBroker: eventBroker,
	}
}

func (h *CheckoutHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)

	var input appcheckout.CheckoutInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Invalid checkout payload")
		return
	}

	result, err := h.useCase.Checkout(tenantID, input)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if h.eventBroker != nil {
		h.eventBroker.Publish(tenantID, "appointment.completed", result.Appointment)
		h.eventBroker.Publish(tenantID, "checkout.completed", result)
	}

	writeJSON(w, http.StatusOK, result)
}

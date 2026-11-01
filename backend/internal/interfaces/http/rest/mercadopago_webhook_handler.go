package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	appappointment "erpagendamentos/backend/internal/application/appointment"
	appfinance "erpagendamentos/backend/internal/application/finance"
	apppayment "erpagendamentos/backend/internal/application/payment"
	financedomain "erpagendamentos/backend/internal/domain/finance"
	paymentdomain "erpagendamentos/backend/internal/domain/payment"
)

type MercadoPagoWebhookHandler struct {
	paymentUseCase     *apppayment.UseCase
	appointmentUseCase *appappointment.UseCase
	financeUseCase     *appfinance.UseCase
	eventBroker        *EventBroker
}

func NewMercadoPagoWebhookHandler(
	paymentUseCase *apppayment.UseCase,
	appointmentUseCase *appappointment.UseCase,
	financeUseCase *appfinance.UseCase,
	eventBroker *EventBroker,
) *MercadoPagoWebhookHandler {
	return &MercadoPagoWebhookHandler{
		paymentUseCase:     paymentUseCase,
		appointmentUseCase: appointmentUseCase,
		financeUseCase:     financeUseCase,
		eventBroker:        eventBroker,
	}
}

func (h *MercadoPagoWebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	paymentID := r.URL.Query().Get("data.id")
	if paymentID == "" {
		paymentID = r.URL.Query().Get("id")
	}

	if paymentID == "" && r.Body != nil {
		var body struct {
			Action string `json:"action"`
			Type   string `json:"type"`
			Data   struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body.Data.ID != "" {
			paymentID = body.Data.ID
		}
	}

	paymentID = strings.TrimSpace(paymentID)
	if paymentID == "" {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ignored", "reason": "no payment id"})
		return
	}

	// Find the payment record
	paymentRecord, err := h.paymentUseCase.GetByProviderPaymentID(paymentID)
	if err != nil {
		// Could be a test webhook or payment not found
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "notice": "payment record not found locally"})
		return
	}

	// Sync with Mercado Pago provider
	updated, err := h.paymentUseCase.SyncProviderStatus(paymentRecord.TenantID, paymentRecord.ID)
	if err != nil {
		writeInternalError(w, err.Error())
		return
	}

	if updated.Status == string(paymentdomain.StatusPaid) {
		// Update appointment payment status if linked
		if updated.AppointmentID != "" && h.appointmentUseCase != nil {
			appt, err := h.appointmentUseCase.GetByID(updated.TenantID, updated.AppointmentID)
			if err == nil {
				appt.PaymentStatus = "paid"
				_, _ = h.appointmentUseCase.Update(appt)
			}
		}

		// Record finance income if not already registered
		if h.financeUseCase != nil {
			_, _ = h.financeUseCase.Create(financedomain.Transaction{
				TenantID:    updated.TenantID,
				Type:        "income",
				Category:    "Serviços",
				Description: fmt.Sprintf("Pagamento Pix aprovado (%s)", updated.Description),
				AmountCents: updated.AmountCents,
				DueDate:     updated.PaidAt,
				PaidAt:      updated.PaidAt,
				Status:      "paid",
				ReferenceID: updated.ID,
			})
		}

		// Broadcast SSE
		if h.eventBroker != nil {
			h.eventBroker.Publish(updated.TenantID, "payment.approved", updated)
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "paymentId": updated.ID})
}

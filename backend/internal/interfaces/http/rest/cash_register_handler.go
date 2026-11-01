package rest

import (
	"encoding/json"
	"net/http"

	appcash "erpagendamentos/backend/internal/application/cashregister"
)

type CashRegisterHandler struct {
	useCase *appcash.UseCase
}

func NewCashRegisterHandler(useCase *appcash.UseCase) *CashRegisterHandler {
	return &CashRegisterHandler{useCase: useCase}
}

func (h *CashRegisterHandler) Current(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	summary, err := h.useCase.GetCurrentShift(tenantID)
	if err != nil {
		if err == appcash.ErrNoOpenShift {
			writeJSON(w, http.StatusOK, map[string]any{"status": "closed", "open": false})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "open",
		"open":    true,
		"summary": summary,
	})
}

func (h *CashRegisterHandler) Open(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	var userID string
	if claims, ok := resolveAuthClaims(r, nil); ok {
		userID = claims.UserID
	}

	var input struct {
		InitialBalanceCents int    `json:"initialBalanceCents"`
		Notes               string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Invalid cash register payload")
		return
	}

	shift, err := h.useCase.OpenShift(tenantID, userID, input.InitialBalanceCents, input.Notes)
	if err != nil {
		if err == appcash.ErrShiftAlreadyOpen {
			writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, shift)
}

func (h *CashRegisterHandler) Entry(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	var userID string
	if claims, ok := resolveAuthClaims(r, nil); ok {
		userID = claims.UserID
	}

	var input struct {
		Type        string `json:"type"` // "sangria", "suprimento"
		AmountCents int    `json:"amountCents"`
		Reason      string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Invalid entry payload")
		return
	}

	entry, err := h.useCase.RecordEntry(tenantID, userID, input.Type, input.AmountCents, input.Reason)
	if err != nil {
		if err == appcash.ErrNoOpenShift {
			writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, entry)
}

func (h *CashRegisterHandler) Close(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	var userID string
	if claims, ok := resolveAuthClaims(r, nil); ok {
		userID = claims.UserID
	}

	var input struct {
		ClosingBalanceCents int    `json:"closingBalanceCents"`
		Notes               string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Invalid close payload")
		return
	}

	shift, err := h.useCase.CloseShift(tenantID, userID, input.ClosingBalanceCents, input.Notes)
	if err != nil {
		if err == appcash.ErrNoOpenShift {
			writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, shift)
}

func (h *CashRegisterHandler) History(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	shifts, err := h.useCase.ListShifts(tenantID, 30)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"shifts": shifts})
}

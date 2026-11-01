package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"erpagendamentos/backend/internal/application/commission"
	commissiondomain "erpagendamentos/backend/internal/domain/commission"
)

type CommissionHandler struct {
	commissionUseCase *commission.UseCase
}

func NewCommissionHandler(commissionUseCase *commission.UseCase) *CommissionHandler {
	return &CommissionHandler{commissionUseCase: commissionUseCase}
}

func (h *CommissionHandler) ListRules(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	rules, err := h.commissionUseCase.ListRules(tenantID)
	if err != nil {
		writeInternalError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"rules": rules})
}

func (h *CommissionHandler) CreateRule(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	var rule commissiondomain.CommissionRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		writeInvalidPayload(w, "invalid request body")
		return
	}
	rule.TenantID = tenantID

	created, err := h.commissionUseCase.CreateRule(rule)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_rule", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *CommissionHandler) UpdateRule(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	id := r.PathValue("ruleID")
	if id == "" {
		id = strings.TrimPrefix(r.URL.Path, "/api/v1/admin/commissions/rules/")
	}

	var rule commissiondomain.CommissionRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		writeInvalidPayload(w, "invalid request body")
		return
	}
	rule.TenantID = tenantID
	rule.ID = id

	updated, err := h.commissionUseCase.UpdateRule(rule)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "update_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *CommissionHandler) DeleteRule(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	id := r.PathValue("ruleID")
	if id == "" {
		id = strings.TrimPrefix(r.URL.Path, "/api/v1/admin/commissions/rules/")
	}

	if err := h.commissionUseCase.DeleteRule(tenantID, id); err != nil {
		writeInternalError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *CommissionHandler) ListRecords(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	professionalID := r.URL.Query().Get("professionalId")
	status := r.URL.Query().Get("status")

	records, err := h.commissionUseCase.ListRecords(tenantID, professionalID, status)
	if err != nil {
		writeInternalError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"records": records})
}

func (h *CommissionHandler) SettleCommissions(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	var payload struct {
		ProfessionalID string   `json:"professionalId"`
		RecordIDs      []string `json:"recordIds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeInvalidPayload(w, "invalid request body")
		return
	}

	if err := h.commissionUseCase.SettleCommissions(tenantID, payload.ProfessionalID, payload.RecordIDs); err != nil {
		writeInternalError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *CommissionHandler) GetSummaries(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	summaries, err := h.commissionUseCase.GetSummaries(tenantID)
	if err != nil {
		writeInternalError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"summaries": summaries})
}

func (h *CommissionHandler) Payout(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	var input commission.PayoutInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "invalid request body")
		return
	}
	if claims, ok := resolveAuthClaims(r, nil); ok {
		input.CreatedByUserID = claims.UserID
	}

	payout, err := h.commissionUseCase.ExecutePayout(tenantID, input)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "payout_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, payout)
}

func (h *CommissionHandler) ListPayouts(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	professionalID := r.URL.Query().Get("professionalId")
	payouts, err := h.commissionUseCase.ListPayouts(tenantID, professionalID)
	if err != nil {
		writeInternalError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"payouts": payouts})
}

package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"erpagendamentos/backend/internal/application/finance"
	financedomain "erpagendamentos/backend/internal/domain/finance"
)

type FinanceHandler struct {
	financeUseCase *finance.UseCase
}

func NewFinanceHandler(financeUseCase *finance.UseCase) *FinanceHandler {
	return &FinanceHandler{financeUseCase: financeUseCase}
}

func (h *FinanceHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	txType := r.URL.Query().Get("type")
	status := r.URL.Query().Get("status")
	fromDate := r.URL.Query().Get("fromDate")
	toDate := r.URL.Query().Get("toDate")

	transactions, err := h.financeUseCase.List(tenantID, txType, status, fromDate, toDate)
	if err != nil {
		writeInternalError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"transactions": transactions})
}

func (h *FinanceHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	var tx financedomain.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		writeInvalidPayload(w, "invalid request body")
		return
	}
	tx.TenantID = tenantID

	created, err := h.financeUseCase.Create(tx)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_transaction", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *FinanceHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	id := r.PathValue("txID")
	if id == "" {
		id = strings.TrimPrefix(r.URL.Path, "/api/v1/admin/finance/transactions/")
	}

	var tx financedomain.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		writeInvalidPayload(w, "invalid request body")
		return
	}
	tx.TenantID = tenantID
	tx.ID = id

	updated, err := h.financeUseCase.Update(tx)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "update_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *FinanceHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	id := r.PathValue("txID")
	if id == "" {
		id = strings.TrimPrefix(r.URL.Path, "/api/v1/admin/finance/transactions/")
	}

	if err := h.financeUseCase.Delete(tenantID, id); err != nil {
		writeInternalError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *FinanceHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	fromDate := r.URL.Query().Get("fromDate")
	toDate := r.URL.Query().Get("toDate")

	summary, err := h.financeUseCase.GetSummary(tenantID, fromDate, toDate)
	if err != nil {
		writeInternalError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

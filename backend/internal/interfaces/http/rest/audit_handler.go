package rest

import (
	"net/http"
	"strconv"

	appaudit "erpagendamentos/backend/internal/application/audit"
)

type AuditHandler struct {
	useCase *appaudit.UseCase
}

func NewAuditHandler(useCase *appaudit.UseCase) *AuditHandler {
	return &AuditHandler{useCase: useCase}
}

func (h *AuditHandler) List(w http.ResponseWriter, r *http.Request) {
	if h.useCase == nil {
		writeJSON(w, http.StatusOK, map[string]any{"logs": []any{}, "total": 0})
		return
	}

	limit := 50
	if rawLimit := r.URL.Query().Get("limit"); rawLimit != "" {
		if val, err := strconv.Atoi(rawLimit); err == nil && val > 0 && val <= 100 {
			limit = val
		}
	}

	offset := 0
	if rawOffset := r.URL.Query().Get("offset"); rawOffset != "" {
		if val, err := strconv.Atoi(rawOffset); err == nil && val >= 0 {
			offset = val
		}
	}

	tenantID := r.URL.Query().Get("tenantId")
	claims, ok := resolveAuthClaims(r, nil)
	if ok && claims.TenantID != "default" && claims.TenantID != "" {
		tenantID = claims.TenantID
	}

	logs, total, err := h.useCase.List(tenantID, limit, offset)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "audit_list_failed", "Não foi possível carregar a trilha de auditoria.")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"logs":   logs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

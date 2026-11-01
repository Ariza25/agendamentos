package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"erpagendamentos/backend/internal/application/calendar"
	calendardomain "erpagendamentos/backend/internal/domain/calendar"
)

type CalendarSyncHandler struct {
	calendarUseCase *calendar.UseCase
}

func NewCalendarSyncHandler(calendarUseCase *calendar.UseCase) *CalendarSyncHandler {
	return &CalendarSyncHandler{calendarUseCase: calendarUseCase}
}

func (h *CalendarSyncHandler) ListConnections(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	connections, err := h.calendarUseCase.ListConnections(tenantID)
	if err != nil {
		writeInternalError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"connections": connections})
}

func (h *CalendarSyncHandler) Connect(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	var conn calendardomain.Connection
	if err := json.NewDecoder(r.Body).Decode(&conn); err != nil {
		writeInvalidPayload(w, "invalid request body")
		return
	}
	conn.TenantID = tenantID

	created, err := h.calendarUseCase.ConnectCalendar(conn)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "connect_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, created)
}

func (h *CalendarSyncHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	id := r.PathValue("connID")
	if id == "" {
		id = strings.TrimPrefix(r.URL.Path, "/api/v1/admin/calendar/connections/")
	}

	if err := h.calendarUseCase.Disconnect(tenantID, id); err != nil {
		writeInternalError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

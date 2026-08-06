package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	appnotification "erpagendamentos/backend/internal/application/notification"
	notificationdomain "erpagendamentos/backend/internal/domain/notification"
)

type NotificationHandler struct {
	useCase *appnotification.UseCase
}

func NewNotificationHandler(useCase *appnotification.UseCase) *NotificationHandler {
	return &NotificationHandler{useCase: useCase}
}

func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	logs, err := h.useCase.ListByTenant(resolveRequestTenantID(r, nil), limit)
	if err != nil {
		writeInternalError(w, "Could not load the notifications.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"notifications": logs})
}

func (h *NotificationHandler) SendTest(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Channel   string `json:"channel"`
		Recipient string `json:"recipient"`
		Subject   string `json:"subject"`
		Body      string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the notification data.")
		return
	}
	logEntry, err := h.useCase.Send(notificationdomain.Log{
		TenantID:  resolveRequestTenantID(r, nil),
		Channel:   input.Channel,
		Recipient: input.Recipient,
		Subject:   input.Subject,
		Body:      input.Body,
		Type:      "test",
	})
	if err != nil {
		writeBadRequestError(w, "invalid_input", "Invalid notification data.")
		return
	}
	if logEntry.Status == "failed" {
		writeAPIError(w, http.StatusBadGateway, "notification_failed", "Could not deliver the notification.")
		return
	}
	writeJSON(w, http.StatusCreated, logEntry)
}

func writeNotificationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, appnotification.ErrNotificationNotFound):
		writeNotFoundError(w, "not_found", "Notification not found.")
	default:
		writeInternalError(w, "Could not process the notification.")
	}
}

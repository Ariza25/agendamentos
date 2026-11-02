package rest

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
)

type errorResponse struct {
	Error       string            `json:"error"`
	Code        string            `json:"code"`
	FieldErrors map[string]string `json:"fieldErrors,omitempty"`
}

func writeCSV(w http.ResponseWriter, filename string, headers []string, rows [][]string) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("\uFEFF"))
	writer := csv.NewWriter(w)
	_ = writer.Write(headers)
	for _, row := range rows {
		_ = writer.Write(row)
	}
	writer.Flush()
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeRawJSON(w http.ResponseWriter, status int, payload []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(payload)
}

func writeAPIError(w http.ResponseWriter, status int, code string, message string) {
	writeJSON(w, status, errorResponse{
		Error: message,
		Code:  code,
	})
}

func writeBadRequestError(w http.ResponseWriter, code string, message string) {
	writeAPIError(w, http.StatusBadRequest, code, message)
}

func writeConflictError(w http.ResponseWriter, code string, message string) {
	writeAPIError(w, http.StatusConflict, code, message)
}

func writeNotFoundError(w http.ResponseWriter, code string, message string) {
	writeAPIError(w, http.StatusNotFound, code, message)
}

func writeInternalError(w http.ResponseWriter, message string) {
	writeAPIError(w, http.StatusInternalServerError, "internal_error", message)
}

func writeValidationError(w http.ResponseWriter, message string, fieldErrors map[string]string) {
	writeJSON(w, http.StatusBadRequest, errorResponse{
		Error:       message,
		Code:        "invalid_input",
		FieldErrors: fieldErrors,
	})
}

func writeInvalidPayload(w http.ResponseWriter, message string) {
	writeBadRequestError(w, "invalid_payload", message)
}

func writeUnauthorized(w http.ResponseWriter) {
	writeAPIError(w, http.StatusUnauthorized, "unauthorized", "Sessão expirada ou usuário não autenticado.")
}

func writeForbidden(w http.ResponseWriter) {
	writeAPIError(w, http.StatusForbidden, "forbidden", "Você não tem permissão para acessar este recurso.")
}

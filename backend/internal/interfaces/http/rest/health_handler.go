package rest

import (
	"net/http"
)

func health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func readiness(checker func() error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if checker == nil {
			writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
			return
		}

		if err := checker(); err != nil {
			logError("readiness check failed", map[string]any{
				"alert":      "readiness_failed",
				"request_id": requestIDFromContext(r.Context()),
				"error":      err.Error(),
			})
			writeAPIError(w, http.StatusServiceUnavailable, "service_unavailable", "Servico temporariamente indisponivel.")
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	}
}

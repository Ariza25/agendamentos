package rest

import (
	"net/http"
	"time"
)

func health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func readiness(checker func() error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if checker == nil {
			writeJSON(w, http.StatusOK, map[string]any{
				"status":    "ready",
				"database":  "memory",
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			})
			return
		}

		if err := checker(); err != nil {
			logError("readiness check failed", map[string]any{
				"alert":      "readiness_failed",
				"request_id": requestIDFromContext(r.Context()),
				"error":      err.Error(),
			})
			writeAPIError(w, http.StatusServiceUnavailable, "service_unavailable", "Database check failed: "+err.Error())
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"status":    "ready",
			"database":  "connected",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	}
}

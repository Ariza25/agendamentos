package rest

import (
	"net/http"
	"os"
	"strconv"
	"strings"
)

const defaultMaxRequestBodyBytes int64 = 2 << 20

func withRequestBodyLimit(next http.Handler) http.Handler {
	maxBytes := resolveMaxRequestBodyBytes()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if shouldLimitRequestBody(r) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		}
		next.ServeHTTP(w, r)
	})
}

func shouldLimitRequestBody(r *http.Request) bool {
	switch r.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return r.Body != nil
	default:
		return false
	}
}

func resolveMaxRequestBodyBytes() int64 {
	raw := strings.TrimSpace(os.Getenv("MAX_REQUEST_BODY_BYTES"))
	if raw == "" {
		return defaultMaxRequestBodyBytes
	}
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || value <= 0 {
		return defaultMaxRequestBodyBytes
	}
	return value
}

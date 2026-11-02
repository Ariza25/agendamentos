package rest

import (
	"net/http"
	"net/url"
	"os"
	"strings"
)

var defaultAllowedOrigins = []string{}

func withCORS(next http.Handler) http.Handler {
	allowedOrigins := resolveAllowedOrigins()

	allowedMethods := "GET, POST, PUT, PATCH, DELETE, OPTIONS"
	allowedHeaders := strings.Join([]string{
		"Accept",
		"Authorization",
		"Content-Type",
		"X-Request-Id",
		"X-Tenant-ID",
		"X-CSRF-Token",
		"X-Client-ID",
		"Idempotency-Key",
		"X-Idempotency-Key",
	}, ", ")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := strings.TrimSpace(r.Header.Get("Origin"))
		originAllowed := origin != "" && isOriginAllowed(origin, allowedOrigins)

		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		w.Header().Add("Vary", "Access-Control-Request-Headers")

		if originAllowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
			w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set(
				"Access-Control-Expose-Headers",
				"X-Request-Id, Retry-After, Idempotency-Key, X-Idempotency-Key, X-Idempotency-Replayed, X-RateLimit-Limit, X-RateLimit-Remaining",
			)
			w.Header().Set("Access-Control-Max-Age", "86400")
		}

		if r.Method == http.MethodOptions {
			if origin == "" {
				http.Error(w, "missing origin", http.StatusBadRequest)
				return
			}

			if !originAllowed {
				http.Error(w, "origin not allowed", http.StatusForbidden)
				return
			}

			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func resolveAllowedOrigins() []string {
	raw := strings.TrimSpace(os.Getenv("ALLOWED_ORIGINS"))
	if raw == "" {
		raw = strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))
	}
	if raw == "" {
		return defaultAllowedOrigins
	}

	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		origin := strings.TrimSpace(part)
		if origin == "" {
			continue
		}
		origins = append(origins, origin)
	}

	if len(origins) == 0 {
		return defaultAllowedOrigins
	}

	return origins
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if isLocalDevelopmentOrigin(origin) {
		return true
	}

	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin == "*" || origin == allowedOrigin {
			return true
		}
	}
	return false
}

func isLocalDevelopmentOrigin(origin string) bool {
	if isAuthCookieProductionEnv() {
		return false
	}
	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}

	hostname := parsed.Hostname()
	return hostname == "localhost" || hostname == "127.0.0.1"
}

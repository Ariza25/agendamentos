package rest

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

const csrfHeaderName = "X-CSRF-Token"

func withCSRFProtection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !requiresCSRFProtection(r) {
			next.ServeHTTP(w, r)
			return
		}

		cookieToken := csrfTokenFromCookie(r)
		headerToken := strings.TrimSpace(r.Header.Get(csrfHeaderName))
		if cookieToken == "" || headerToken == "" || subtle.ConstantTimeCompare([]byte(cookieToken), []byte(headerToken)) != 1 {
			writeAPIError(w, http.StatusForbidden, "csrf_validation_failed", "A sessao de seguranca expirou. Atualize a pagina e tente novamente.")
			return
		}

		origin := strings.TrimSpace(r.Header.Get("Origin"))
		if origin != "" && !isOriginAllowed(origin, resolveAllowedOrigins()) {
			writeForbidden(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func requiresCSRFProtection(r *http.Request) bool {
	if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
		return false
	}
	if authTokenFromCookie(r) == "" {
		return false
	}
	path := strings.TrimSpace(r.URL.Path)
	if strings.HasPrefix(path, "/api/v1/auth/login") ||
		strings.HasPrefix(path, "/api/v1/auth/password-reset") {
		return false
	}
	return strings.HasPrefix(path, "/api/v1/admin/") ||
		path == "/api/v1/auth/logout" ||
		path == "/api/v1/auth/profile" ||
		path == "/api/v1/auth/password" ||
		strings.HasPrefix(path, "/api/v1/auth/sessions") ||
		strings.HasPrefix(path, "/api/v1/auth/mfa")
}

package rest

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	authCookieName = "erpagendamentos_session"
	csrfCookieName = "erpagendamentos_csrf"
)

func setAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:        resolvedAuthCookieName(),
		Value:       strings.TrimSpace(token),
		Path:        "/",
		Domain:      authCookieDomain(),
		MaxAge:      int((7 * 24 * time.Hour).Seconds()),
		HttpOnly:    true,
		Secure:      authCookieSecure(),
		SameSite:    authCookieSameSite(),
		Partitioned: authCookiePartitioned(),
	})
}

func setSessionCookies(w http.ResponseWriter, token string) string {
	setAuthCookie(w, token)
	csrfToken := newCSRFToken()
	http.SetCookie(w, &http.Cookie{
		Name:        resolvedCSRFCookieName(),
		Value:       csrfToken,
		Path:        "/",
		Domain:      authCookieDomain(),
		MaxAge:      int((7 * 24 * time.Hour).Seconds()),
		HttpOnly:    false,
		Secure:      authCookieSecure(),
		SameSite:    authCookieSameSite(),
		Partitioned: authCookiePartitioned(),
	})
	return csrfToken
}

func clearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:        resolvedAuthCookieName(),
		Value:       "",
		Path:        "/",
		Domain:      authCookieDomain(),
		MaxAge:      -1,
		HttpOnly:    true,
		Secure:      authCookieSecure(),
		SameSite:    authCookieSameSite(),
		Partitioned: authCookiePartitioned(),
	})
	http.SetCookie(w, &http.Cookie{
		Name:        resolvedCSRFCookieName(),
		Value:       "",
		Path:        "/",
		Domain:      authCookieDomain(),
		MaxAge:      -1,
		HttpOnly:    false,
		Secure:      authCookieSecure(),
		SameSite:    authCookieSameSite(),
		Partitioned: authCookiePartitioned(),
	})
}

func authTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(resolvedAuthCookieName())
	if err != nil {
		return ""
	}
	return strings.TrimSpace(cookie.Value)
}

func csrfTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(resolvedCSRFCookieName())
	if err != nil {
		return ""
	}
	return strings.TrimSpace(cookie.Value)
}

func newCSRFToken() string {
	value := make([]byte, 32)
	if _, err := rand.Read(value); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return base64.RawURLEncoding.EncodeToString(value)
}

func authCookieSecure() bool {
	raw := strings.ToLower(strings.TrimSpace(os.Getenv("AUTH_COOKIE_SECURE")))
	if raw == "1" || raw == "true" || raw == "yes" {
		return true
	}
	if raw == "0" || raw == "false" || raw == "no" {
		return false
	}
	return isAuthCookieProductionEnv()
}

func authCookieSameSite() http.SameSite {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("AUTH_COOKIE_SAMESITE"))) {
	case "strict":
		return http.SameSiteStrictMode
	case "lax":
		return http.SameSiteLaxMode
	default:
		if isAuthCookieProductionEnv() {
			return http.SameSiteNoneMode
		}
		return http.SameSiteLaxMode
	}
}

func authCookieDomain() string {
	return strings.TrimSpace(os.Getenv("AUTH_COOKIE_DOMAIN"))
}

func authCookiePartitioned() bool {
	raw := strings.ToLower(strings.TrimSpace(os.Getenv("AUTH_COOKIE_PARTITIONED")))
	if raw == "1" || raw == "true" || raw == "yes" {
		return authCookieSecure()
	}
	if raw == "0" || raw == "false" || raw == "no" {
		return false
	}
	return authCookieSecure() && authCookieSameSite() == http.SameSiteNoneMode && authCookieDomain() == ""
}

func resolvedAuthCookieName() string {
	if isAuthCookieProductionEnv() && authCookieDomain() == "" {
		return "__Host-" + authCookieName
	}
	return authCookieName
}

func resolvedCSRFCookieName() string {
	if isAuthCookieProductionEnv() && authCookieDomain() == "" {
		return "__Host-" + csrfCookieName
	}
	return csrfCookieName
}

func isAuthCookieProductionEnv() bool {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("ERPAGENDAMENTOS_ENV")))
	return env == "production" || env == "prod"
}

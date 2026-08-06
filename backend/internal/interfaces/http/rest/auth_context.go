package rest

import (
	"context"
	"net/http"
	"os"
	"strings"

	appauth "erpagendamentos/backend/internal/application/auth"
	authdomain "erpagendamentos/backend/internal/domain/auth"
)

type authClaimsContextKey struct{}

func requestWithAuthClaims(r *http.Request, claims authdomain.Claims) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), authClaimsContextKey{}, claims))
}

func resolveAuthClaims(r *http.Request, authUseCase *appauth.UseCase) (authdomain.Claims, bool) {
	if claims, ok := r.Context().Value(authClaimsContextKey{}).(authdomain.Claims); ok && claims.UserID != "" {
		return claims, true
	}
	if authUseCase == nil {
		return authdomain.Claims{}, false
	}

	token := authTokenFromCookie(r)
	if token == "" && bearerAuthAllowed() {
		header := strings.TrimSpace(r.Header.Get("Authorization"))
		if strings.HasPrefix(header, "Bearer ") {
			token = strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		}
	}
	if token == "" {
		return authdomain.Claims{}, false
	}

	claims, err := authUseCase.ParseToken(token)
	if err != nil {
		return authdomain.Claims{}, false
	}

	return claims, true
}

func bearerAuthAllowed() bool {
	raw := strings.ToLower(strings.TrimSpace(os.Getenv("ALLOW_BEARER_AUTH")))
	if raw == "true" || raw == "1" || raw == "yes" {
		return true
	}
	return !isAuthCookieProductionEnv()
}

func withProtectedNamespaces(authUseCase *appauth.UseCase, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimSpace(r.URL.Path)
		if !strings.HasPrefix(path, "/api/v1/admin/") {
			next.ServeHTTP(w, r)
			return
		}
		claims, ok := resolveAuthClaims(r, authUseCase)
		if !ok {
			writeUnauthorized(w)
			return
		}
		if !isAdminRole(claims.Role) {
			writeForbidden(w)
			return
		}
		next.ServeHTTP(w, requestWithAuthClaims(r, claims))
	})
}

func resolveTenantID(r *http.Request, authUseCase *appauth.UseCase) string {
	if claims, ok := resolveAuthClaims(r, authUseCase); ok {
		if account, err := authUseCase.Me(claims.UserID); err == nil && strings.TrimSpace(account.TenantID) != "" {
			return account.TenantID
		}
		if strings.TrimSpace(claims.TenantID) != "" {
			return claims.TenantID
		}
	}
	return "default"
}

func resolveRequestTenantID(r *http.Request, authUseCase *appauth.UseCase) string {
	if claims, ok := resolveAuthClaims(r, authUseCase); ok {
		if isAdminRole(claims.Role) {
			if tenantID := strings.TrimSpace(claims.TenantID); tenantID != "" {
				return tenantID
			}
			return "default"
		}
		if tenantID := strings.TrimSpace(claims.TenantID); tenantID != "" {
			return tenantID
		}
	}
	if value := strings.TrimSpace(r.Header.Get("X-Tenant-ID")); value != "" {
		return value
	}
	if value := strings.TrimSpace(r.URL.Query().Get("tenantId")); value != "" {
		return value
	}
	if value := strings.TrimSpace(r.URL.Query().Get("tenant")); value != "" {
		return value
	}
	return "default"
}

func requireAdmin(authUseCase *appauth.UseCase, next http.HandlerFunc) http.HandlerFunc {
	return requireRoles(authUseCase, []string{"owner", "admin", "receptionist", "professional"}, next)
}

func requireRoles(authUseCase *appauth.UseCase, allowedRoles []string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := resolveAuthClaims(r, authUseCase)
		if !ok {
			writeUnauthorized(w)
			return
		}
		account, err := authUseCase.Me(claims.UserID)
		if err != nil {
			writeUnauthorized(w)
			return
		}
		claims.Role = account.Role
		claims.TenantID = account.TenantID

		if !hasRole(claims.Role, allowedRoles) {
			writeForbidden(w)
			return
		}

		next(w, requestWithAuthClaims(r, claims))
	}
}

func hasRole(role string, allowedRoles []string) bool {
	for _, allowed := range allowedRoles {
		if role == allowed {
			return true
		}
	}
	return false
}

func isAdminRole(role string) bool {
	switch role {
	case "owner", "admin", "receptionist", "professional":
		return true
	default:
		return false
	}
}

func isPrivilegedAdminRole(role string) bool {
	return role == "admin" || role == "owner"
}

func requireOwnerOrAdmin(authUseCase *appauth.UseCase, next http.HandlerFunc) http.HandlerFunc {
	return requireRoles(authUseCase, []string{"owner", "admin"}, next)
}

func requirePrivilegedAdmin(authUseCase *appauth.UseCase, next http.HandlerFunc) http.HandlerFunc {
	return requireRoles(authUseCase, []string{"owner", "admin"}, next)
}

func requireMasterAdmin(authUseCase *appauth.UseCase, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := resolveAuthClaims(r, authUseCase)
		if !ok {
			writeUnauthorized(w)
			return
		}
		account, err := authUseCase.Me(claims.UserID)
		if err != nil {
			writeUnauthorized(w)
			return
		}
		if account.Role != "admin" || (strings.TrimSpace(account.TenantID) != "" && strings.TrimSpace(account.TenantID) != "default") {
			writeForbidden(w)
			return
		}
		next(w, r)
	}
}

func resolveActorTenantID(r *http.Request, authUseCase *appauth.UseCase) string {
	if claims, ok := resolveAuthClaims(r, authUseCase); ok {
		if authUseCase != nil {
			if account, err := authUseCase.Me(claims.UserID); err == nil {
				return account.TenantID
			}
		}
		return claims.TenantID
	}
	return ""
}

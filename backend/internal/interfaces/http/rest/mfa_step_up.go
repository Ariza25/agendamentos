package rest

import (
	"net/http"
	"os"
	"strings"
	"time"

	appauth "erpagendamentos/backend/internal/application/auth"
)

func withPrivilegedStepUp(authUseCase *appauth.UseCase, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !mfaEnforcementEnabled() || !isSensitiveAdminMutation(r) {
			next.ServeHTTP(w, r)
			return
		}
		claims, ok := resolveAuthClaims(r, authUseCase)
		if !ok || !isPrivilegedAdminRole(claims.Role) {
			next.ServeHTTP(w, r)
			return
		}
		account, err := authUseCase.Me(claims.UserID)
		if err != nil || !account.MFAEnabled {
			writeAPIError(w, http.StatusPreconditionRequired, "mfa_setup_required", "Ative a autenticacao em duas etapas para realizar esta acao.")
			return
		}
		if claims.MFAAt == 0 || time.Since(time.Unix(claims.MFAAt, 0)) > 15*time.Minute {
			writeAPIError(w, http.StatusPreconditionRequired, "mfa_step_up_required", "Confirme o codigo do autenticador para continuar.")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func mfaEnforcementEnabled() bool {
	raw := strings.ToLower(strings.TrimSpace(os.Getenv("FEATURE_PRIVILEGED_MFA_ENFORCEMENT")))
	if raw != "" {
		return raw == "1" || raw == "true" || raw == "yes"
	}
	return isAuthCookieProductionEnv()
}

func isSensitiveAdminMutation(r *http.Request) bool {
	if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
		return false
	}
	path := strings.TrimSpace(r.URL.Path)
	if path == "/api/v1/auth/mfa/step-up" || strings.HasPrefix(path, "/api/v1/auth/mfa") {
		return false
	}
	return strings.HasPrefix(path, "/api/v1/admin/users") ||
		strings.HasPrefix(path, "/api/v1/admin/tenants") ||
		strings.HasPrefix(path, "/api/v1/admin/configuration")
}

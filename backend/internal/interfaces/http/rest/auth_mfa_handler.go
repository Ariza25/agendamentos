package rest

import (
	"encoding/json"
	"net/http"

	appauth "erpagendamentos/backend/internal/application/auth"
)

type authMFAHandler struct {
	useCase *appauth.UseCase
}

func newAuthMFAHandler(useCase *appauth.UseCase) authMFAHandler {
	return authMFAHandler{useCase: useCase}
}

func (h authMFAHandler) status(w http.ResponseWriter, r *http.Request) {
	claims, ok := resolveAuthClaims(r, h.useCase)
	if !ok {
		writeUnauthorized(w)
		return
	}
	status, err := h.useCase.MFAStatus(claims.UserID)
	if err != nil {
		writeForbidden(w)
		return
	}
	writeJSON(w, http.StatusOK, status)
}

func (h authMFAHandler) beginSetup(w http.ResponseWriter, r *http.Request) {
	claims, ok := resolveAuthClaims(r, h.useCase)
	if !ok {
		writeUnauthorized(w)
		return
	}
	secret, uri, err := h.useCase.BeginMFASetup(claims.UserID)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "mfa_setup_failed", "Could not start the setup.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"secret": secret, "otpauthUri": uri})
}

func (h authMFAHandler) enable(w http.ResponseWriter, r *http.Request) {
	claims, ok := resolveAuthClaims(r, h.useCase)
	if !ok {
		writeUnauthorized(w)
		return
	}
	var input struct {
		Secret string `json:"secret"`
		Code   string `json:"code"`
	}
	if json.NewDecoder(r.Body).Decode(&input) != nil {
		writeInvalidPayload(w, "Provide the authenticator code.")
		return
	}
	if err := h.useCase.EnableMFA(claims.UserID, input.Secret, input.Code); err != nil {
		logWarn("mfa enable failed", map[string]any{"alert": "auth_suspicious_activity", "user_id": claims.UserID, "tenant_id": claims.TenantID})
		writeAPIError(w, http.StatusBadRequest, "invalid_mfa_code", "Invalid verification code.")
		return
	}
	logInfo("mfa enabled", map[string]any{"user_id": claims.UserID, "tenant_id": claims.TenantID})
	clearAuthCookie(w)
	writeJSON(w, http.StatusOK, map[string]any{"enabled": true, "reauthenticationRequired": true})
}

func (h authMFAHandler) disable(w http.ResponseWriter, r *http.Request) {
	claims, ok := resolveAuthClaims(r, h.useCase)
	if !ok {
		writeUnauthorized(w)
		return
	}
	var input struct {
		Password string `json:"password"`
		Code     string `json:"code"`
	}
	if json.NewDecoder(r.Body).Decode(&input) != nil {
		writeInvalidPayload(w, "Provide password and authenticator code.")
		return
	}
	if err := h.useCase.DisableMFA(claims.UserID, input.Password, input.Code); err != nil {
		logWarn("mfa disable failed", map[string]any{"alert": "auth_suspicious_activity", "user_id": claims.UserID, "tenant_id": claims.TenantID})
		writeAPIError(w, http.StatusBadRequest, "invalid_mfa_confirmation", "Invalid password or verification code.")
		return
	}
	logInfo("mfa disabled", map[string]any{"user_id": claims.UserID, "tenant_id": claims.TenantID})
	clearAuthCookie(w)
	writeJSON(w, http.StatusOK, map[string]any{"enabled": false, "reauthenticationRequired": true})
}

func (h authMFAHandler) stepUp(w http.ResponseWriter, r *http.Request) {
	claims, ok := resolveAuthClaims(r, h.useCase)
	if !ok {
		writeUnauthorized(w)
		return
	}
	var input struct {
		Code string `json:"code"`
	}
	if json.NewDecoder(r.Body).Decode(&input) != nil {
		writeInvalidPayload(w, "Provide the authenticator code.")
		return
	}
	session, err := h.useCase.StepUpMFA(claims.UserID, input.Code)
	if err != nil {
		logWarn("mfa step-up failed", map[string]any{"alert": "auth_suspicious_activity", "user_id": claims.UserID, "tenant_id": claims.TenantID})
		writeAPIError(w, http.StatusBadRequest, "invalid_mfa_code", "Invalid verification code.")
		return
	}
	logInfo("mfa step-up completed", map[string]any{"user_id": claims.UserID, "tenant_id": claims.TenantID})
	csrfToken := setSessionCookies(w, session.Token)
	writeJSON(w, http.StatusOK, map[string]any{"csrfToken": csrfToken, "verified": true})
}

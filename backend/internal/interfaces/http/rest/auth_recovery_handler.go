package rest

import (
	"encoding/json"
	"net/http"

	appauth "erpagendamentos/backend/internal/application/auth"
	authdomain "erpagendamentos/backend/internal/domain/auth"
)

type authPasswordRecoveryHandler struct {
	useCase *appauth.UseCase
}

func (h *AuthHandler) passwordRecovery() authPasswordRecoveryHandler {
	return authPasswordRecoveryHandler{useCase: h.useCase}
}

func (h authPasswordRecoveryHandler) requestLink(w http.ResponseWriter, r *http.Request) {
	var input authdomain.PasswordResetRequestInput
	if json.NewDecoder(r.Body).Decode(&input) != nil {
		writeInvalidPayload(w, "Could not read the reset data.")
		return
	}
	token, err := h.useCase.RequestPasswordReset(input)
	if err != nil {
		logError("password reset request failed", map[string]any{"error": err.Error()})
		writeAPIError(w, http.StatusServiceUnavailable, "auth_service_unavailable", "Could not process the request right now.")
		return
	}
	payload := map[string]any{"sent": true}
	if token != "" && exposeAuthTokenForTests() {
		payload["token"] = token
	}
	writeJSON(w, http.StatusAccepted, payload)
}

func (h authPasswordRecoveryHandler) requestCode(w http.ResponseWriter, r *http.Request) {
	var input authdomain.PasswordResetRequestInput
	if json.NewDecoder(r.Body).Decode(&input) != nil {
		writeInvalidPayload(w, "Could not read the reset data.")
		return
	}
	allowed, err := h.useCase.AllowPasswordResetCodeRequest(input.Email, requestClientIP(r))
	if err != nil {
		logError("password reset code rate limit failed", map[string]any{"error": err.Error()})
		writeAPIError(w, http.StatusServiceUnavailable, "auth_service_unavailable", "Could not process the request right now.")
		return
	}
	if !allowed {
		w.Header().Set("Retry-After", "900")
		writeAPIError(w, http.StatusTooManyRequests, "rate_limited", "Too many attempts. Wait a few minutes and try again.")
		return
	}
	code, err := h.useCase.RequestPasswordResetCode(input)
	if err != nil {
		logError("password reset code request failed", map[string]any{"error": err.Error()})
		writeAPIError(w, http.StatusServiceUnavailable, "auth_service_unavailable", "Could not process the request right now.")
		return
	}
	payload := map[string]any{"sent": true}
	if code != "" && exposeAuthTokenForTests() {
		payload["code"] = code
	}
	writeJSON(w, http.StatusAccepted, payload)
}

func (h authPasswordRecoveryHandler) confirmLink(w http.ResponseWriter, r *http.Request) {
	var input authdomain.PasswordResetConfirmInput
	if json.NewDecoder(r.Body).Decode(&input) != nil {
		writeInvalidPayload(w, "Could not read the new password.")
		return
	}
	if err := h.useCase.ConfirmPasswordReset(input); err != nil {
		if validation, ok := err.(*appauth.ValidationError); ok {
			writeValidationError(w, "Review the password fields and try again.", validation.FieldErrors)
			return
		}
		writeAPIError(w, http.StatusBadRequest, "invalid_reset_token", "Invalid or expired reset link.")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h authPasswordRecoveryHandler) confirmCode(w http.ResponseWriter, r *http.Request) {
	var input authdomain.PasswordResetCodeConfirmInput
	if json.NewDecoder(r.Body).Decode(&input) != nil {
		writeInvalidPayload(w, "Could not read the reset code.")
		return
	}
	allowed, err := h.useCase.AllowPasswordResetCodeAttempt(requestClientIP(r))
	if err != nil {
		logError("password reset code attempt rate limit failed", map[string]any{"error": err.Error()})
		writeAPIError(w, http.StatusServiceUnavailable, "auth_service_unavailable", "Could not process the request right now.")
		return
	}
	if !allowed {
		w.Header().Set("Retry-After", "900")
		writeAPIError(w, http.StatusTooManyRequests, "rate_limited", "Too many attempts. Wait a few minutes and try again.")
		return
	}
	if err := h.useCase.ConfirmPasswordResetCode(input); err != nil {
		if validation, ok := err.(*appauth.ValidationError); ok {
			writeValidationError(w, "Review the fields and try again.", validation.FieldErrors)
			return
		}
		writeAPIError(w, http.StatusBadRequest, "invalid_reset_code", "Invalid or expired code.")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

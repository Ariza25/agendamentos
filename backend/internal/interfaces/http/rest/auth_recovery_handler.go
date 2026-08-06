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
		writeInvalidPayload(w, "Nao foi possivel ler os dados de redefinicao.")
		return
	}
	token, err := h.useCase.RequestPasswordReset(input)
	if err != nil {
		logError("password reset request failed", map[string]any{"error": err.Error()})
		writeAPIError(w, http.StatusServiceUnavailable, "auth_service_unavailable", "Nao foi possivel processar a solicitacao agora.")
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
		writeInvalidPayload(w, "Nao foi possivel ler os dados de redefinicao.")
		return
	}
	allowed, err := h.useCase.AllowPasswordResetCodeRequest(input.Email, requestClientIP(r))
	if err != nil {
		logError("password reset code rate limit failed", map[string]any{"error": err.Error()})
		writeAPIError(w, http.StatusServiceUnavailable, "auth_service_unavailable", "Nao foi possivel processar a solicitacao agora.")
		return
	}
	if !allowed {
		w.Header().Set("Retry-After", "900")
		writeAPIError(w, http.StatusTooManyRequests, "rate_limited", "Muitas tentativas. Aguarde alguns minutos e tente novamente.")
		return
	}
	code, err := h.useCase.RequestPasswordResetCode(input)
	if err != nil {
		logError("password reset code request failed", map[string]any{"error": err.Error()})
		writeAPIError(w, http.StatusServiceUnavailable, "auth_service_unavailable", "Nao foi possivel processar a solicitacao agora.")
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
		writeInvalidPayload(w, "Nao foi possivel ler a nova senha.")
		return
	}
	if err := h.useCase.ConfirmPasswordReset(input); err != nil {
		if validation, ok := err.(*appauth.ValidationError); ok {
			writeValidationError(w, "Revise os campos de senha e tente novamente.", validation.FieldErrors)
			return
		}
		writeAPIError(w, http.StatusBadRequest, "invalid_reset_token", "Link de redefinicao invalido ou expirado.")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h authPasswordRecoveryHandler) confirmCode(w http.ResponseWriter, r *http.Request) {
	var input authdomain.PasswordResetCodeConfirmInput
	if json.NewDecoder(r.Body).Decode(&input) != nil {
		writeInvalidPayload(w, "Nao foi possivel ler o codigo de redefinicao.")
		return
	}
	allowed, err := h.useCase.AllowPasswordResetCodeAttempt(requestClientIP(r))
	if err != nil {
		logError("password reset code attempt rate limit failed", map[string]any{"error": err.Error()})
		writeAPIError(w, http.StatusServiceUnavailable, "auth_service_unavailable", "Nao foi possivel processar a solicitacao agora.")
		return
	}
	if !allowed {
		w.Header().Set("Retry-After", "900")
		writeAPIError(w, http.StatusTooManyRequests, "rate_limited", "Muitas tentativas. Aguarde alguns minutos e tente novamente.")
		return
	}
	if err := h.useCase.ConfirmPasswordResetCode(input); err != nil {
		if validation, ok := err.(*appauth.ValidationError); ok {
			writeValidationError(w, "Revise os campos e tente novamente.", validation.FieldErrors)
			return
		}
		writeAPIError(w, http.StatusBadRequest, "invalid_reset_code", "Codigo invalido ou expirado.")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

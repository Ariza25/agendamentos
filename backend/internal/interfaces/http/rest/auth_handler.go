package rest

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	appauth "erpagendamentos/backend/internal/application/auth"
	authdomain "erpagendamentos/backend/internal/domain/auth"
)

type AuthHandler struct {
	useCase *appauth.UseCase
}

func NewAuthHandler(useCase *appauth.UseCase) *AuthHandler {
	return &AuthHandler{useCase: useCase}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input authdomain.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Nao foi possivel ler os dados de login. Revise os campos enviados.")
		return
	}
	for _, key := range []string{requestClientIP(r), strings.ToLower(strings.TrimSpace(input.Audience)) + ":" + strings.ToLower(strings.TrimSpace(input.Email))} {
		allowed, limitErr := h.useCase.AllowSecurityRequest("login", key, 10, 15*time.Minute)
		if limitErr != nil {
			writeAPIError(w, http.StatusServiceUnavailable, "auth_service_unavailable", "Nao foi possivel autenticar agora.")
			return
		}
		if !allowed {
			w.Header().Set("Retry-After", "900")
			writeAPIError(w, http.StatusTooManyRequests, "rate_limited", "Muitas tentativas de acesso. Aguarde alguns minutos.")
			return
		}
	}

	session, err := h.useCase.Login(input)
	if err != nil {
		if errors.Is(err, appauth.ErrMFARequired) {
			writeAPIError(w, http.StatusPreconditionRequired, "mfa_required", "Informe o codigo do aplicativo autenticador.")
			return
		}
		if validation, ok := err.(*appauth.ValidationError); ok {
			writeValidationError(w, "Revise os campos destacados e tente novamente.", validation.FieldErrors)
			return
		}
		logWarn("admin login failed", map[string]any{
			"alert":      "auth_suspicious_activity",
			"request_id": requestIDFromContext(r.Context()),
			"source":     "auth",
			"error":      err.Error(),
		})
		writeAPIError(w, http.StatusUnauthorized, "invalid_credentials", "Email ou senha incorretos.")
		return
	}

	csrfToken := setSessionCookies(w, session.Token)
	h.trackSession(r, session.Token)
	writeJSON(w, http.StatusOK, browserSessionPayload(session, csrfToken))
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if claims, ok := resolveAuthClaims(r, h.useCase); ok && claims.SessionID != "" {
		_ = h.useCase.RevokeSession(claims.UserID, claims.SessionID)
	}
	clearAuthCookie(w)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := resolveAuthClaims(r, h.useCase)
	if !ok {
		writeUnauthorized(w)
		return
	}

	account, err := h.useCase.Me(claims.UserID)
	if err != nil {
		writeUnauthorized(w)
		return
	}

	writeJSON(w, http.StatusOK, account)
}

func (h *AuthHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	claims, ok := resolveAuthClaims(r, h.useCase)
	if !ok {
		writeUnauthorized(w)
		return
	}
	sessions, err := h.useCase.ListSessions(claims.UserID, claims.SessionID)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "session_list_failed", "Nao foi possivel carregar os dispositivos conectados.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"sessions": sessions})
}

func (h *AuthHandler) RevokeSession(w http.ResponseWriter, r *http.Request) {
	claims, ok := resolveAuthClaims(r, h.useCase)
	if !ok {
		writeUnauthorized(w)
		return
	}
	sessionID := strings.TrimSpace(r.PathValue("sessionID"))
	if sessionID == "" {
		writeInvalidPayload(w, "Informe a sessao que sera encerrada.")
		return
	}
	if err := h.useCase.RevokeSession(claims.UserID, sessionID); err != nil {
		writeAPIError(w, http.StatusInternalServerError, "session_revoke_failed", "Nao foi possivel encerrar a sessao.")
		return
	}
	if sessionID == claims.SessionID {
		clearAuthCookie(w)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) RevokeOtherSessions(w http.ResponseWriter, r *http.Request) {
	claims, ok := resolveAuthClaims(r, h.useCase)
	if !ok {
		writeUnauthorized(w)
		return
	}
	if err := h.useCase.RevokeOtherSessions(claims.UserID, claims.SessionID); err != nil {
		writeAPIError(w, http.StatusInternalServerError, "session_revoke_failed", "Nao foi possivel encerrar as outras sessoes.")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := resolveAuthClaims(r, h.useCase)
	if !ok {
		writeUnauthorized(w)
		return
	}

	var input authdomain.UpdateProfileInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Nao foi possivel ler os dados do perfil. Revise os campos enviados.")
		return
	}
	input.UserID = claims.UserID

	account, err := h.useCase.UpdateProfile(input)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_profile", "Nao foi possivel atualizar o perfil.")
		return
	}

	writeJSON(w, http.StatusOK, account)
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims, ok := resolveAuthClaims(r, h.useCase)
	if !ok {
		writeUnauthorized(w)
		return
	}

	var input authdomain.ChangePasswordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Nao foi possivel ler os dados de senha. Revise os campos enviados.")
		return
	}
	input.UserID = claims.UserID

	if err := h.useCase.ChangePassword(input); err != nil {
		if validation, ok := err.(*appauth.ValidationError); ok {
			writeValidationError(w, "Revise os campos de senha e tente novamente.", validation.FieldErrors)
			return
		}
		writeAPIError(w, http.StatusBadRequest, "password_change_failed", "Nao foi possivel alterar a senha.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	h.passwordRecovery().requestLink(w, r)
}

func (h *AuthHandler) RequestPasswordResetCode(w http.ResponseWriter, r *http.Request) {
	h.passwordRecovery().requestCode(w, r)
}

func (h *AuthHandler) ConfirmPasswordReset(w http.ResponseWriter, r *http.Request) {
	h.passwordRecovery().confirmLink(w, r)
}

func (h *AuthHandler) ConfirmPasswordResetCode(w http.ResponseWriter, r *http.Request) {
	h.passwordRecovery().confirmCode(w, r)
}

func (h *AuthHandler) MFAStatus(w http.ResponseWriter, r *http.Request) {
	newAuthMFAHandler(h.useCase).status(w, r)
}

func (h *AuthHandler) BeginMFASetup(w http.ResponseWriter, r *http.Request) {
	newAuthMFAHandler(h.useCase).beginSetup(w, r)
}

func (h *AuthHandler) EnableMFA(w http.ResponseWriter, r *http.Request) {
	newAuthMFAHandler(h.useCase).enable(w, r)
}

func (h *AuthHandler) DisableMFA(w http.ResponseWriter, r *http.Request) {
	newAuthMFAHandler(h.useCase).disable(w, r)
}

func (h *AuthHandler) StepUpMFA(w http.ResponseWriter, r *http.Request) {
	newAuthMFAHandler(h.useCase).stepUp(w, r)
}

func (h *AuthHandler) trackSession(r *http.Request, token string) {
	claims, err := h.useCase.ParseToken(token)
	if err != nil || claims.SessionID == "" {
		return
	}
	deviceName := strings.TrimSpace(r.UserAgent())
	if len(deviceName) > 180 {
		deviceName = deviceName[:180]
	}
	h.useCase.TrackSession(claims.SessionID, r.Header.Get("X-Client-ID"), deviceName, requestClientIP(r))
}

func (h *AuthHandler) ListAdminUsers(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.useCase.ListAccounts()
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "internal_error", "Nao foi possivel carregar os usuarios.")
		return
	}

	tenantID := ""
	if claims, ok := resolveAuthClaims(r, h.useCase); ok && claims.Role != "admin" && claims.TenantID != "" {
		tenantID = claims.TenantID
	}

	panelUsers := make([]authdomain.Account, 0, len(accounts))
	for _, account := range accounts {
		switch account.Role {
		case "owner", "admin", "receptionist", "professional":
			if tenantID == "" || account.TenantID == tenantID {
				panelUsers = append(panelUsers, account)
			}
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{"users": panelUsers})
}

func (h *AuthHandler) CreateAdminUser(w http.ResponseWriter, r *http.Request) {
	var input authdomain.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Nao foi possivel ler os dados do usuario. Revise os campos enviados.")
		return
	}

	if claims, ok := resolveAuthClaims(r, h.useCase); ok && claims.TenantID != "" {
		input.TenantID = claims.TenantID
	}

	password := secureTemporaryPassword()
	input.Password = password

	account, err := h.useCase.CreateStaffAccount(input)
	if err != nil {
		writeAdminAuthError(w, err)
		return
	}

	var inviteURL string
	if token, err := h.useCase.CreatePasswordResetToken(account.ID, 72*time.Hour); err == nil {
		inviteURL = buildPasswordResetURL(r, token)
	}

	response := map[string]any{
		"user":            account,
		"inviteUrl":       inviteURL,
		"inviteEmailSent": false,
	}
	if exposeAuthTokenForTests() {
		response["temporaryPassword"] = password
	}
	writeJSON(w, http.StatusCreated, response)
}

func (h *AuthHandler) UpdateAdminUserRole(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Role string `json:"role"`
	}
	var input payload
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Nao foi possivel ler o perfil do usuario. Revise os campos enviados.")
		return
	}

	account, err := h.useCase.UpdateAccountRole(resolveActorTenantID(r, h.useCase), r.PathValue("userID"), input.Role)
	if err != nil {
		writeAdminAuthError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, account)
}

func (h *AuthHandler) DeleteAdminUser(w http.ResponseWriter, r *http.Request) {
	if err := h.useCase.DeleteStaffAccount(resolveActorTenantID(r, h.useCase), r.PathValue("userID")); err != nil {
		writeAdminAuthError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeAdminAuthError(w http.ResponseWriter, err error) {
	if validation, ok := err.(*appauth.ValidationError); ok {
		writeValidationError(w, "Revise os campos destacados e tente novamente.", validation.FieldErrors)
		return
	}

	switch {
	case errors.Is(err, appauth.ErrInvalidRole):
		writeBadRequestError(w, "invalid_input", "Perfil de usuario invalido.")
	case errors.Is(err, appauth.ErrAccountNotFound):
		writeNotFoundError(w, "not_found", "Usuario nao encontrado.")
	case errors.Is(err, appauth.ErrLastAdmin):
		writeConflictError(w, "last_admin", "Nao e possivel remover o ultimo acesso administrativo.")
	case errors.Is(err, appauth.ErrEmailAlreadyRegistered):
		writeConflictError(w, "email_already_registered", "Ja existe uma conta cadastrada com este email.")
	default:
		writeInternalError(w, "Falha interna ao processar usuario.")
	}
}

func browserSessionPayload(session authdomain.Session, csrfToken string) map[string]any {
	payload := map[string]any{"user": session.User, "csrfToken": csrfToken}
	if exposeAuthTokenForTests() {
		payload["token"] = session.Token
	}
	return payload
}

func exposeAuthTokenForTests() bool {
	if isAuthCookieProductionEnv() {
		return false
	}
	raw := strings.ToLower(strings.TrimSpace(os.Getenv("EXPOSE_AUTH_TOKEN_FOR_TESTS")))
	return raw == "1" || raw == "true" || raw == "yes"
}

func secureTemporaryPassword() string {
	randomBytes := make([]byte, 18)
	if _, err := rand.Read(randomBytes); err != nil {
		return "Aa1!" + "temporary-password-rotate-now"
	}
	return "Aa1!" + base64.RawURLEncoding.EncodeToString(randomBytes)
}

func buildPasswordResetURL(r *http.Request, token string) string {
	baseURL := strings.TrimRight(os.Getenv("ADMIN_BASE_URL"), "/")
	if baseURL == "" {
		if origin := strings.TrimRight(r.Header.Get("Origin"), "/"); strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1") {
			baseURL = origin
		}
	}
	if baseURL == "" {
		baseURL = "http://localhost:5174"
	}
	return baseURL + "/#/reset-password?token=" + token
}

package rest

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	appauth "erpagendamentos/backend/internal/application/auth"
	apponboarding "erpagendamentos/backend/internal/application/onboarding"
	authdomain "erpagendamentos/backend/internal/domain/auth"
)

type AuthGoogleHandler struct {
	authUseCase       *appauth.UseCase
	onboardingUseCase *apponboarding.UseCase
}

func NewAuthGoogleHandler(authUseCase *appauth.UseCase, onboardingUseCase *apponboarding.UseCase) *AuthGoogleHandler {
	return &AuthGoogleHandler{
		authUseCase:       authUseCase,
		onboardingUseCase: onboardingUseCase,
	}
}

type googleStateData struct {
	Audience  string `json:"aud"`
	ReturnURL string `json:"return_url"`
	CreatedAt int64  `json:"t"`
}

type googleTokenResponse struct {
	AccessToken string `json:"access_token"`
	IdToken     string `json:"id_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type googleUserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Picture       string `json:"picture"`
}

func (h *AuthGoogleHandler) Status(w http.ResponseWriter, r *http.Request) {
	clientID := strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID"))
	writeJSON(w, http.StatusOK, map[string]any{
		"enabled":  clientID != "",
		"clientId": clientID,
	})
}

func (h *AuthGoogleHandler) Login(w http.ResponseWriter, r *http.Request) {
	clientID := strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID"))
	frontendURL := resolveFrontendURL(r)

	if clientID == "" {
		if strings.Contains(r.Header.Get("Accept"), "application/json") {
			writeAPIError(w, http.StatusBadRequest, "google_not_configured", "A integração com o Google ainda não foi configurada. As chaves GOOGLE_CLIENT_ID e GOOGLE_CLIENT_SECRET devem ser adicionadas nas variáveis de ambiente.")
			return
		}
		separator := "?"
		if strings.Contains(frontendURL, "?") {
			separator = "&"
		}
		http.Redirect(w, r, fmt.Sprintf("%s%sgoogle_error=not_configured", frontendURL, separator), http.StatusFound)
		return
	}

	redirectURI := resolveGoogleRedirectURI(r)
	stateObj := googleStateData{
		Audience:  strings.TrimSpace(r.URL.Query().Get("audience")),
		ReturnURL: strings.TrimSpace(r.URL.Query().Get("return_url")),
		CreatedAt: time.Now().Unix(),
	}
	stateJSON, _ := json.Marshal(stateObj)
	encodedState := base64.RawURLEncoding.EncodeToString(stateJSON)

	googleAuthURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&access_type=online&prompt=select_account&state=%s",
		url.QueryEscape(clientID),
		url.QueryEscape(redirectURI),
		url.QueryEscape("openid email profile"),
		url.QueryEscape(encodedState),
	)

	http.Redirect(w, r, googleAuthURL, http.StatusFound)
}

func (h *AuthGoogleHandler) Callback(w http.ResponseWriter, r *http.Request) {
	frontendURL := resolveFrontendURL(r)

	if errParam := r.URL.Query().Get("error"); errParam != "" {
		http.Redirect(w, r, appendQueryParam(frontendURL, "google_error", "cancelled"), http.StatusFound)
		return
	}

	code := strings.TrimSpace(r.URL.Query().Get("code"))
	if code == "" {
		http.Redirect(w, r, appendQueryParam(frontendURL, "google_error", "missing_code"), http.StatusFound)
		return
	}

	clientID := strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID"))
	clientSecret := strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_SECRET"))
	if clientID == "" || clientSecret == "" {
		http.Redirect(w, r, appendQueryParam(frontendURL, "google_error", "not_configured"), http.StatusFound)
		return
	}

	// Parse state if present to retrieve custom return_url
	if stateParam := r.URL.Query().Get("state"); stateParam != "" {
		if stateBytes, err := base64.RawURLEncoding.DecodeString(stateParam); err == nil {
			var stateObj googleStateData
			if json.Unmarshal(stateBytes, &stateObj) == nil && stateObj.ReturnURL != "" {
				frontendURL = stateObj.ReturnURL
			}
		}
	}

	// Exchange authorization code for Google tokens
	tokenURL := "https://oauth2.googleapis.com/token"
	formValues := url.Values{}
	formValues.Set("code", code)
	formValues.Set("client_id", clientID)
	formValues.Set("client_secret", clientSecret)
	formValues.Set("redirect_uri", resolveGoogleRedirectURI(r))
	formValues.Set("grant_type", "authorization_code")

	tokenResp, err := http.PostForm(tokenURL, formValues)
	if err != nil || tokenResp.StatusCode != http.StatusOK {
		http.Redirect(w, r, appendQueryParam(frontendURL, "google_error", "token_exchange_failed"), http.StatusFound)
		return
	}
	defer tokenResp.Body.Close()

	var tokenData googleTokenResponse
	if err := json.NewDecoder(tokenResp.Body).Decode(&tokenData); err != nil || tokenData.AccessToken == "" {
		http.Redirect(w, r, appendQueryParam(frontendURL, "google_error", "invalid_token_response"), http.StatusFound)
		return
	}

	// Retrieve user profile from Google
	userInfo, err := fetchGoogleUserInfo(tokenData.AccessToken)
	if err != nil || userInfo.Email == "" {
		http.Redirect(w, r, appendQueryParam(frontendURL, "google_error", "user_info_failed"), http.StatusFound)
		return
	}

	// Authenticate existing account or sign up new tenant
	session, err := h.authenticateOrRegisterGoogleUser(userInfo)
	if err != nil {
		http.Redirect(w, r, appendQueryParam(frontendURL, "google_error", "auth_failed"), http.StatusFound)
		return
	}

	// Set session cookies
	csrfToken := setSessionCookies(w, session.Token)

	userJSON, _ := json.Marshal(session.User)
	fragment := fmt.Sprintf(
		"#google_auth=success&csrfToken=%s&user=%s&token=%s",
		url.QueryEscape(csrfToken),
		url.QueryEscape(string(userJSON)),
		url.QueryEscape(session.Token),
	)

	http.Redirect(w, r, frontendURL+fragment, http.StatusFound)
}

func (h *AuthGoogleHandler) VerifyIDToken(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Credential string `json:"credential"`
		Audience   string `json:"audience"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Credential) == "" {
		writeInvalidPayload(w, "Informe a credencial do Google para autenticação.")
		return
	}

	// Verify ID token with Google's public tokeninfo endpoint
	tokenInfoURL := "https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(body.Credential)
	resp, err := http.Get(tokenInfoURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		writeAPIError(w, http.StatusUnauthorized, "invalid_credentials", "Não foi possível validar a credencial do Google.")
		return
	}
	defer resp.Body.Close()

	var userInfo googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil || userInfo.Email == "" {
		writeAPIError(w, http.StatusUnauthorized, "invalid_credentials", "A conta Google não forneceu um e-mail válido.")
		return
	}

	session, err := h.authenticateOrRegisterGoogleUser(userInfo)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "signup_failed", err.Error())
		return
	}

	csrfToken := setSessionCookies(w, session.Token)
	writeJSON(w, http.StatusOK, browserSessionPayload(session, csrfToken))
}

func (h *AuthGoogleHandler) authenticateOrRegisterGoogleUser(userInfo googleUserInfo) (authdomain.Session, error) {
	email := strings.ToLower(strings.TrimSpace(userInfo.Email))
	name := strings.TrimSpace(userInfo.Name)
	if name == "" {
		name = "Administrador"
	}

	// Check if an account with this email already exists
	existingAccount, err := h.authUseCase.GetStaffAccountByEmail(email)
	if err == nil {
		// Existing user: generate session
		return h.authUseCase.BuildSessionForAccount(existingAccount)
	}

	// New user: register tenant + owner account with 14-day free trial
	session, _, err := h.onboardingUseCase.SignupWithExternalProvider(name, email, "google")
	return session, err
}

func fetchGoogleUserInfo(accessToken string) (googleUserInfo, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return googleUserInfo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return googleUserInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return googleUserInfo{}, fmt.Errorf("google userinfo error: %s", string(body))
	}

	var info googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return googleUserInfo{}, err
	}
	return info, nil
}

func resolveGoogleRedirectURI(r *http.Request) string {
	if redirect := strings.TrimSpace(os.Getenv("GOOGLE_REDIRECT_URL")); redirect != "" {
		return redirect
	}
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/api/v1/auth/google/callback", scheme, r.Host)
}

func resolveFrontendURL(r *http.Request) string {
	if target := strings.TrimSpace(r.URL.Query().Get("return_url")); target != "" {
		return target
	}
	baseURL := strings.TrimRight(os.Getenv("ADMIN_BASE_URL"), "/")
	if baseURL == "" {
		baseURL = strings.TrimRight(os.Getenv("BASE_URL"), "/")
	}
	if baseURL == "" {
		baseURL = "https://ariza25.github.io/agendamentos/admin"
	}
	return baseURL
}

func appendQueryParam(rawURL, key, val string) string {
	separator := "?"
	if strings.Contains(rawURL, "?") {
		separator = "&"
	}
	return fmt.Sprintf("%s%s%s=%s", rawURL, separator, url.QueryEscape(key), url.QueryEscape(val))
}

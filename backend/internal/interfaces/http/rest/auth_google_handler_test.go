package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	appauth "erpagendamentos/backend/internal/application/auth"
	apponboarding "erpagendamentos/backend/internal/application/onboarding"
	"erpagendamentos/backend/internal/infrastructure/persistence/memory"
	"erpagendamentos/backend/internal/infrastructure/security"
	"erpagendamentos/backend/internal/interfaces/http/rest"
)

func TestAuthGoogleHandlerStatusAndLogin(t *testing.T) {
	tenantRepo := memory.NewTenantRepository()
	authRepo := memory.NewAuthRepository()
	configRepo := memory.NewConfigurationRepository()
	serviceRepo := memory.NewServiceRepository()
	proRepo := memory.NewProfessionalRepository()

	hasher := security.NewBcryptHasher(4)
	tokens := security.NewJWTManager("test-jwt-secret-key-32-chars-long", time.Hour)
	authUseCase := appauth.NewUseCase(authRepo, hasher, tokens)

	onboardingUseCase := apponboarding.NewUseCase(
		tenantRepo,
		authUseCase,
		configRepo,
		serviceRepo,
		proRepo,
	)

	handler := rest.NewAuthGoogleHandler(authUseCase, onboardingUseCase)

	// 1. Status without GOOGLE_CLIENT_ID should be disabled
	os.Unsetenv("GOOGLE_CLIENT_ID")
	req := httptest.NewRequest("GET", "/api/v1/auth/google/status", nil)
	rec := httptest.NewRecorder()
	handler.Status(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var statusResp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &statusResp); err != nil {
		t.Fatalf("failed to parse json: %v", err)
	}
	if statusResp["enabled"] != false {
		t.Errorf("expected enabled=false, got %v", statusResp["enabled"])
	}

	// 2. Status with GOOGLE_CLIENT_ID should be enabled
	os.Setenv("GOOGLE_CLIENT_ID", "test-google-client-id-123")
	rec2 := httptest.NewRecorder()
	handler.Status(rec2, req)
	var statusResp2 map[string]any
	_ = json.Unmarshal(rec2.Body.Bytes(), &statusResp2)
	if statusResp2["enabled"] != true {
		t.Errorf("expected enabled=true, got %v", statusResp2["enabled"])
	}
	if statusResp2["clientId"] != "test-google-client-id-123" {
		t.Errorf("expected clientId=test-google-client-id-123, got %v", statusResp2["clientId"])
	}

	// 3. Login without GOOGLE_CLIENT_ID should return not_configured
	os.Unsetenv("GOOGLE_CLIENT_ID")
	recLoginUnset := httptest.NewRecorder()
	reqLoginUnset := httptest.NewRequest("GET", "/api/v1/auth/google/login", nil)
	reqLoginUnset.Header.Set("Accept", "application/json")
	handler.Login(recLoginUnset, reqLoginUnset)
	if recLoginUnset.Code != http.StatusBadRequest {
		t.Errorf("expected 400 when unconfigured, got %d", recLoginUnset.Code)
	}

	// 4. Login with GOOGLE_CLIENT_ID should redirect to Google
	os.Setenv("GOOGLE_CLIENT_ID", "my-client-id")
	recLogin := httptest.NewRecorder()
	reqLogin := httptest.NewRequest("GET", "/api/v1/auth/google/login", nil)
	handler.Login(recLogin, reqLogin)
	if recLogin.Code != http.StatusFound {
		t.Fatalf("expected 302 redirect, got %d", recLogin.Code)
	}
	location := recLogin.Header().Get("Location")
	if !strings.Contains(location, "accounts.google.com/o/oauth2/v2/auth") {
		t.Errorf("expected Google auth URL, got %s", location)
	}
	if !strings.Contains(location, "client_id=my-client-id") {
		t.Errorf("expected client_id in URL, got %s", location)
	}

	// Cleanup
	os.Unsetenv("GOOGLE_CLIENT_ID")
}


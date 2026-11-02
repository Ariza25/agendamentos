package rest

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestRateLimitAllowsNormalTrafficAndLimitsAbuse(t *testing.T) {
	_ = os.Setenv("RATE_LIMIT_PER_MINUTE", "5")
	defer os.Unsetenv("RATE_LIMIT_PER_MINUTE")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})

	limiterMw := withRateLimit(nil, handler)

	for i := 1; i <= 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/general", nil)
		req.RemoteAddr = "203.0.113.10:12345"
		rec := httptest.NewRecorder()
		limiterMw.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("request %d expected status 200, got %d", i, rec.Code)
		}
		if rec.Header().Get("X-RateLimit-Limit") != "5" {
			t.Fatalf("expected X-RateLimit-Limit: 5, got %s", rec.Header().Get("X-RateLimit-Limit"))
		}
	}

	// 6th request from same IP should be blocked
	req := httptest.NewRequest(http.MethodGet, "/api/v1/general", nil)
	req.RemoteAddr = "203.0.113.10:12345"
	rec := httptest.NewRecorder()
	limiterMw.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 Too Many Requests, got %d", rec.Code)
	}
	if rec.Header().Get("Retry-After") == "" {
		t.Fatalf("expected Retry-After header on 429 response")
	}
}

func TestSensitiveRoutesClassification(t *testing.T) {
	tests := []struct {
		method    string
		path      string
		sensitive bool
	}{
		{"POST", "/api/v1/auth/login", true},
		{"POST", "/api/v1/auth/signup", true},
		{"POST", "/api/v1/onboarding/signup", true},
		{"POST", "/api/v1/auth/mfa/challenge", true},
		{"POST", "/api/v1/public/clinica/appointments", true},
		{"POST", "/api/v1/public/clinica/portal/lookup", true},
		{"POST", "/api/v1/admin/checkout", true},
		{"POST", "/api/v1/admin/cash-register/open", true},
		{"POST", "/api/v1/admin/commissions/settle", true},
		{"GET", "/api/v1/public/clinica/services", false},
		{"GET", "/api/v1/admin/customers", false},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		isSens := isSensitiveRateLimitedRoute(req)
		if isSens != tt.sensitive {
			t.Errorf("path %s %s: expected sensitive=%v, got %v", tt.method, tt.path, tt.sensitive, isSens)
		}
	}
}

func TestPersistentRateLimitConfigMatch(t *testing.T) {
	tests := []struct {
		method string
		path   string
		match  bool
		scope  string
	}{
		{"POST", "/api/v1/auth/login", true, "login-route"},
		{"POST", "/api/v1/auth/signup", true, "auth-signup"},
		{"POST", "/api/v1/public/clinica-x/appointments", true, "public-appointments"},
		{"POST", "/api/v1/public/clinica-x/portal/lookup", true, "portal-lookup"},
		{"POST", "/api/v1/admin/checkout", true, "admin-checkout"},
		{"POST", "/api/v1/admin/cash-register/close", true, "admin-cash-register"},
		{"GET", "/api/v1/admin/customers", false, ""},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		cfg, ok := persistentRateLimitConfig(req)
		if ok != tt.match {
			t.Errorf("%s %s expected match=%v, got %v", tt.method, tt.path, tt.match, ok)
		}
		if ok && cfg.scope != tt.scope {
			t.Errorf("%s %s expected scope=%s, got %s", tt.method, tt.path, tt.scope, cfg.scope)
		}
	}
}

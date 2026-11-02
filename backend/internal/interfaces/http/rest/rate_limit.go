package rest

import (
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	appauth "erpagendamentos/backend/internal/application/auth"
)

const (
	defaultRateLimitPerMinute          = 600
	defaultSensitiveRateLimitPerMinute = 60
)

type rateLimitBucket struct {
	windowStart time.Time
	count       int
}

type memoryRateLimiter struct {
	mu        sync.Mutex
	buckets   map[string]rateLimitBucket
	limit     int
	sensitive int
	window    time.Duration
}

func newMemoryRateLimiter() *memoryRateLimiter {
	return &memoryRateLimiter{
		buckets:   make(map[string]rateLimitBucket),
		limit:     resolvePositiveEnvInt("RATE_LIMIT_PER_MINUTE", defaultRateLimitPerMinute),
		sensitive: resolvePositiveEnvInt("SENSITIVE_RATE_LIMIT_PER_MINUTE", defaultSensitiveRateLimitPerMinute),
		window:    time.Minute,
	}
}

func withRateLimit(authUseCase *appauth.UseCase, next http.Handler) http.Handler {
	limiter := newMemoryRateLimiter()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}
		clientIP := requestClientIP(r)
		routeKey := rateLimitRouteKey(r)

		if config, ok := persistentRateLimitConfig(r); ok && authUseCase != nil {
			identities := []string{"ip:" + clientIP}
			if clientID := normalizedClientID(r.Header.Get("X-Client-ID")); clientID != "" {
				identities = append(identities, "device:"+clientID)
			}
			if claims, authenticated := resolveAuthClaims(r, authUseCase); authenticated {
				identities = append(identities, "user:"+claims.UserID, "tenant:"+claims.TenantID)
			}
			for _, identity := range identities {
				allowed, err := authUseCase.AllowSecurityRequest(config.scope, identity, config.limit, config.window)
				if err != nil {
					writeAPIError(w, http.StatusServiceUnavailable, "rate_limit_unavailable", "Could not validate the request right now.")
					return
				}
				if !allowed {
					w.Header().Set("Retry-After", strconv.Itoa(int(config.window.Seconds())))
					writeRateLimited(w, r, routeKey, clientIP, config.scope)
					return
				}
			}
		}

		allowed, remaining, limit, retryAfter := limiter.allow(clientIP, routeKey, isSensitiveRateLimitedRoute(r))
		if !allowed {
			w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
			writeRateLimited(w, r, routeKey, clientIP, "local-fallback")
			return
		}

		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		next.ServeHTTP(w, r)
	})
}

func normalizedClientID(value string) string {
	value = strings.TrimSpace(value)
	if len(value) < 16 || len(value) > 128 {
		return ""
	}
	for _, char := range value {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && (char < '0' || char > '9') && char != '-' && char != '_' {
			return ""
		}
	}
	return value
}

type persistentRateLimit struct {
	scope  string
	limit  int
	window time.Duration
}

func persistentRateLimitConfig(r *http.Request) (persistentRateLimit, bool) {
	path := strings.TrimSpace(r.URL.Path)
	switch {
	case path == "/api/v1/auth/login" || path == "/api/v1/auth/google":
		return persistentRateLimit{scope: "login-route", limit: 30, window: 15 * time.Minute}, true
	case path == "/api/v1/auth/signup" || path == "/api/v1/onboarding/signup":
		return persistentRateLimit{scope: "auth-signup", limit: 10, window: 15 * time.Minute}, true
	case path == "/api/v1/auth/password-reset/request":
		return persistentRateLimit{scope: "password-reset-request", limit: 10, window: 15 * time.Minute}, true
	case path == "/api/v1/auth/password-reset/confirm":
		return persistentRateLimit{scope: "password-reset-confirm", limit: 30, window: 15 * time.Minute}, true
	case path == "/api/v1/auth/password-reset/code/request":
		return persistentRateLimit{scope: "password-reset-code-request", limit: 10, window: 15 * time.Minute}, true
	case path == "/api/v1/auth/password-reset/code/confirm":
		return persistentRateLimit{scope: "password-reset-code-confirm", limit: 30, window: 15 * time.Minute}, true
	case strings.HasPrefix(path, "/api/v1/auth/mfa"):
		return persistentRateLimit{scope: "auth-mfa", limit: 20, window: 15 * time.Minute}, true
	case strings.HasPrefix(path, "/api/v1/auth/sessions") && r.Method != http.MethodGet:
		return persistentRateLimit{scope: "session-management", limit: 20, window: time.Minute}, true
	case strings.HasPrefix(path, "/api/v1/public/") && strings.HasSuffix(path, "/appointments") && r.Method == http.MethodPost:
		return persistentRateLimit{scope: "public-appointments", limit: 25, window: 15 * time.Minute}, true
	case strings.HasPrefix(path, "/api/v1/public/") && strings.Contains(path, "/portal/lookup") && r.Method == http.MethodPost:
		return persistentRateLimit{scope: "portal-lookup", limit: 15, window: 15 * time.Minute}, true
	case strings.HasPrefix(path, "/api/v1/public/") && strings.Contains(path, "/waiting-list") && r.Method == http.MethodPost:
		return persistentRateLimit{scope: "public-waiting-list", limit: 15, window: 15 * time.Minute}, true
	case strings.HasPrefix(path, "/api/v1/public/") && strings.Contains(path, "/reviews/") && r.Method == http.MethodPost:
		return persistentRateLimit{scope: "public-reviews", limit: 15, window: 15 * time.Minute}, true
	case strings.HasPrefix(path, "/api/v1/admin/checkout") && r.Method != http.MethodGet:
		return persistentRateLimit{scope: "admin-checkout", limit: 60, window: time.Minute}, true
	case strings.HasPrefix(path, "/api/v1/admin/cash-register") && r.Method != http.MethodGet:
		return persistentRateLimit{scope: "admin-cash-register", limit: 60, window: time.Minute}, true
	case strings.HasPrefix(path, "/api/v1/admin/commissions") && r.Method != http.MethodGet:
		return persistentRateLimit{scope: "admin-commissions", limit: 60, window: time.Minute}, true
	case strings.HasPrefix(path, "/api/v1/admin/finance") && r.Method != http.MethodGet:
		return persistentRateLimit{scope: "admin-finance", limit: 60, window: time.Minute}, true
	default:
		return persistentRateLimit{}, false
	}
}

func writeRateLimited(w http.ResponseWriter, r *http.Request, routeKey string, clientIP string, scope string) {
	logWarn("request rate limited", map[string]any{
		"alert":      "rate_limited",
		"request_id": requestIDFromContext(r.Context()),
		"ip":         clientIP,
		"method":     r.Method,
		"path":       r.URL.Path,
		"route_key":  routeKey,
		"scope":      scope,
	})
	writeAPIError(w, http.StatusTooManyRequests, "rate_limited", "Too many attempts. Please wait a moment and try again.")
}

func (l *memoryRateLimiter) allow(clientIP string, routeKey string, sensitive bool) (bool, int, int, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	l.cleanup(now)
	key := clientIP + "|" + routeKey
	bucket := l.buckets[key]
	if bucket.windowStart.IsZero() || now.Sub(bucket.windowStart) >= l.window {
		bucket = rateLimitBucket{windowStart: now}
	}
	bucket.count++
	l.buckets[key] = bucket
	limit := l.limit
	if sensitive {
		limit = l.sensitive
	}
	remaining := limit - bucket.count
	if remaining < 0 {
		remaining = 0
	}
	if bucket.count <= limit {
		return true, remaining, limit, 0
	}
	retryAfter := l.window - now.Sub(bucket.windowStart)
	if retryAfter < time.Second {
		retryAfter = time.Second
	}
	return false, 0, limit, retryAfter
}

func (l *memoryRateLimiter) cleanup(now time.Time) {
	for key, bucket := range l.buckets {
		if now.Sub(bucket.windowStart) > 2*l.window {
			delete(l.buckets, key)
		}
	}
}

func requestClientIP(r *http.Request) string {
	remoteIP := remoteClientIP(r.RemoteAddr)
	if !isTrustedProxy(remoteIP) {
		return remoteIP
	}
	if value := forwardedClientIP(r.Header.Get("X-Forwarded-For")); value != "" {
		return value
	}
	if value := validIPAddress(r.Header.Get("X-Real-IP")); value != "" {
		return value
	}
	return remoteIP
}

func remoteClientIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil && strings.TrimSpace(host) != "" {
		return host
	}
	return strings.TrimSpace(remoteAddr)
}

func isTrustedProxy(clientIP string) bool {
	ip := net.ParseIP(strings.TrimSpace(clientIP))
	if ip == nil {
		return false
	}
	for _, rawCIDR := range strings.Split(os.Getenv("TRUSTED_PROXY_CIDRS"), ",") {
		_, network, err := net.ParseCIDR(strings.TrimSpace(rawCIDR))
		if err == nil && network.Contains(ip) {
			return true
		}
	}
	return false
}

func forwardedClientIP(value string) string {
	parts := strings.Split(value, ",")
	for index := len(parts) - 1; index >= 0; index-- {
		ip := validIPAddress(parts[index])
		if ip != "" && !isTrustedProxy(ip) {
			return ip
		}
	}
	return ""
}

func validIPAddress(value string) string {
	value = strings.TrimSpace(value)
	if net.ParseIP(value) == nil {
		return ""
	}
	return value
}

func rateLimitRouteKey(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" {
		path = "/"
	}
	return r.Method + " " + path
}

func isSensitiveRateLimitedRoute(r *http.Request) bool {
	path := strings.TrimSpace(r.URL.Path)
	switch {
	case strings.HasPrefix(path, "/api/v1/auth/login") || strings.HasPrefix(path, "/api/v1/auth/google"):
		return true
	case strings.HasPrefix(path, "/api/v1/auth/signup") || strings.HasPrefix(path, "/api/v1/onboarding/signup"):
		return true
	case strings.HasPrefix(path, "/api/v1/auth/password-reset"):
		return true
	case strings.HasPrefix(path, "/api/v1/auth/mfa"):
		return true
	case strings.HasPrefix(path, "/api/v1/public/") && r.Method != http.MethodGet:
		return true
	case strings.HasPrefix(path, "/api/v1/admin/checkout") && r.Method != http.MethodGet:
		return true
	case strings.HasPrefix(path, "/api/v1/admin/cash-register") && r.Method != http.MethodGet:
		return true
	case strings.HasPrefix(path, "/api/v1/admin/commissions") && r.Method != http.MethodGet:
		return true
	case strings.HasPrefix(path, "/api/v1/admin/finance") && r.Method != http.MethodGet:
		return true
	default:
		return false
	}
}

func resolvePositiveEnvInt(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}

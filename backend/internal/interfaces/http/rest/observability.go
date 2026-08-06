package rest

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type contextKey string

const requestIDContextKey contextKey = "request_id"

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func withObservability(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := resolveRequestID(r)
		startedAt := time.Now()

		w.Header().Set("X-Request-Id", requestID)
		recorder := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
		ctx := context.WithValue(r.Context(), requestIDContextKey, requestID)
		r = r.WithContext(ctx)

		defer func() {
			if recovered := recover(); recovered != nil {
				logError("panic recovered", map[string]any{
					"alert":      "panic",
					"request_id": requestID,
					"method":     r.Method,
					"path":       r.URL.Path,
					"panic":      fmt.Sprint(recovered),
				})
				writeAPIError(recorder, http.StatusInternalServerError, "internal_error", "Falha interna no servidor.")
			}

			durationMs := time.Since(startedAt).Milliseconds()
			fields := map[string]any{
				"request_id":  requestID,
				"method":      r.Method,
				"path":        r.URL.Path,
				"status_code": recorder.status,
				"duration_ms": durationMs,
				"remote_addr": r.RemoteAddr,
			}
			if recorder.status >= http.StatusInternalServerError {
				fields["alert"] = "http_5xx"
				logError("request completed with server error", fields)
				return
			}
			if r.Method == http.MethodOptions {
				return
			}
		}()

		next.ServeHTTP(recorder, r)
	})
}

func resolveRequestID(r *http.Request) string {
	if requestID := r.Header.Get("X-Request-Id"); requestID != "" {
		return requestID
	}
	if requestID := r.Header.Get("X-Request-ID"); requestID != "" {
		return requestID
	}
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err == nil {
		return hex.EncodeToString(raw[:])
	}
	return time.Now().UTC().Format("20060102150405.000000000")
}

func requestIDFromContext(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDContextKey).(string)
	return requestID
}

func logInfo(message string, fields map[string]any) {
	logStructured("INFO", message, fields)
}

func logWarn(message string, fields map[string]any) {
	logStructured("WARN", message, fields)
}

func logError(message string, fields map[string]any) {
	logStructured("ERROR", message, fields)
}

func logStructured(level string, message string, fields map[string]any) {
	payload := map[string]any{
		"level":   level,
		"message": message,
	}
	for key, value := range fields {
		payload[key] = sanitizeLogValue(key, value)
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		log.Printf(`{"level":"ERROR","message":"failed to encode log entry","error":"%v"}`, err)
		return
	}
	log.Println(string(encoded))
}

func sanitizeLogValue(key string, value any) any {
	if isSensitiveLogKey(key) {
		return "[REDACTED]"
	}
	switch typed := value.(type) {
	case map[string]any:
		sanitized := make(map[string]any, len(typed))
		for childKey, childValue := range typed {
			sanitized[childKey] = sanitizeLogValue(childKey, childValue)
		}
		return sanitized
	case map[string]string:
		sanitized := make(map[string]string, len(typed))
		for childKey, childValue := range typed {
			if isSensitiveLogKey(childKey) {
				sanitized[childKey] = "[REDACTED]"
				continue
			}
			sanitized[childKey] = truncateLogString(childValue)
		}
		return sanitized
	case []any:
		sanitized := make([]any, 0, len(typed))
		for _, item := range typed {
			sanitized = append(sanitized, sanitizeLogValue(key, item))
		}
		return sanitized
	case string:
		return sanitizeLogString(typed)
	default:
		return value
	}
}

func isSensitiveLogKey(key string) bool {
	normalized := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(key), "-", "_"), ".", "_"))
	sensitiveParts := []string{
		"authorization",
		"access_token",
		"api_token",
		"secret",
		"password",
		"passwd",
		"token",
		"credential",
		"card",
		"cpf",
		"cnpj",
		"document",
		"phone",
		"email",
		"payload",
		"body",
		"qr_code",
	}
	for _, part := range sensitiveParts {
		if strings.Contains(normalized, part) {
			return true
		}
	}
	return false
}

func sanitizeLogString(value string) string {
	value = redactSensitiveText(value)
	return truncateLogString(value)
}

func redactSensitiveText(value string) string {
	if value == "" {
		return value
	}
	lower := strings.ToLower(value)
	sensitiveMarkers := []string{
		"authorization",
		"access_token",
		"api_token",
		"secret",
		"password",
		"passwd",
		"credential",
		"card_token",
		"card number",
		"qr_code",
		"qrcode",
		"qr code",
		"payload",
		"request body",
		"response body",
	}
	for _, marker := range sensitiveMarkers {
		if strings.Contains(lower, marker) {
			return "[REDACTED]"
		}
	}
	value = emailLogPattern.ReplaceAllString(value, "[REDACTED_EMAIL]")
	value = cpfLogPattern.ReplaceAllString(value, "[REDACTED_CPF]")
	value = cnpjLogPattern.ReplaceAllString(value, "[REDACTED_CNPJ]")
	value = bearerLogPattern.ReplaceAllString(value, "Bearer [REDACTED]")
	value = longNumberLogPattern.ReplaceAllString(value, "[REDACTED_NUMBER]")
	return value
}

func truncateLogString(value string) string {
	const maxLogStringLength = 512
	if len(value) <= maxLogStringLength {
		return value
	}
	return value[:maxLogStringLength] + "...[TRUNCATED]"
}

var (
	emailLogPattern      = regexp.MustCompile(`(?i)[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}`)
	cpfLogPattern        = regexp.MustCompile(`\b\d{3}\.?\d{3}\.?\d{3}-?\d{2}\b`)
	cnpjLogPattern       = regexp.MustCompile(`\b\d{2}\.?\d{3}\.?\d{3}/?\d{4}-?\d{2}\b`)
	bearerLogPattern     = regexp.MustCompile(`(?i)Bearer\s+[A-Za-z0-9._\-]+`)
	longNumberLogPattern = regexp.MustCompile(`\b\d{13,19}\b`)
)

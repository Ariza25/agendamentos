package rest

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"sync"
	"time"

	appauth "erpagendamentos/backend/internal/application/auth"
)

const (
	defaultIdempotencyTTL   = 24 * time.Hour
	maxIdempotencyBodyBytes = 4 * 1024 * 1024 // 4MB buffer cap
)

type idempotencyEntry struct {
	inFlight   bool
	statusCode int
	headers    http.Header
	body       []byte
	createdAt  time.Time
}

type IdempotencyStore struct {
	mu      sync.RWMutex
	entries map[string]*idempotencyEntry
	ttl     time.Duration
}

func NewIdempotencyStore(ttl time.Duration) *IdempotencyStore {
	if ttl <= 0 {
		ttl = defaultIdempotencyTTL
	}
	return &IdempotencyStore{
		entries: make(map[string]*idempotencyEntry),
		ttl:     ttl,
	}
}

// LockOrGet attempts to reserve the key.
// If an entry is already in flight, returns (entry, false).
// If an entry is completed, returns (entry, false).
// If the key is new or expired, reserves it in flight and returns (entry, true).
func (s *IdempotencyStore) LockOrGet(key string) (*idempotencyEntry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	if existing, ok := s.entries[key]; ok {
		if now.Sub(existing.createdAt) > s.ttl {
			delete(s.entries, key)
		} else {
			return existing, false
		}
	}

	entry := &idempotencyEntry{
		inFlight:  true,
		createdAt: now,
	}
	s.entries[key] = entry
	return entry, true
}

func (s *IdempotencyStore) Save(key string, statusCode int, headers http.Header, body []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entry, ok := s.entries[key]; ok {
		entry.inFlight = false
		entry.statusCode = statusCode
		entry.headers = headers.Clone()
		entry.body = body
		entry.createdAt = time.Now()
	}
}

func (s *IdempotencyStore) Release(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, key)
}

func (s *IdempotencyStore) Cleanup(now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range s.entries {
		if now.Sub(v.createdAt) > s.ttl {
			delete(s.entries, k)
		}
	}
}

type idempotencyRecorder struct {
	http.ResponseWriter
	statusCode  int
	body        bytes.Buffer
	wroteHeader bool
	maxBodyLen  int
}

func (rec *idempotencyRecorder) WriteHeader(code int) {
	if rec.wroteHeader {
		return
	}
	rec.statusCode = code
	rec.wroteHeader = true
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *idempotencyRecorder) Write(b []byte) (int, error) {
	if !rec.wroteHeader {
		rec.WriteHeader(http.StatusOK)
	}
	if rec.maxBodyLen <= 0 || rec.body.Len()+len(b) <= rec.maxBodyLen {
		rec.body.Write(b)
	}
	return rec.ResponseWriter.Write(b)
}

func (rec *idempotencyRecorder) Flush() {
	if flusher, ok := rec.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// withIdempotency intercepts mutating requests that supply an Idempotency-Key.
// If an identical request is already running, it returns HTTP 409 Conflict.
// If an identical request previously succeeded or produced a client error, it replays the response.
// If the server fails with 5xx or panics, the key reservation is released so the client can retry.
func withIdempotency(authUseCase *appauth.UseCase, next http.Handler) http.Handler {
	store := NewIdempotencyStore(defaultIdempotencyTTL)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only mutating methods need idempotency deduplication
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		rawKey := strings.TrimSpace(r.Header.Get("Idempotency-Key"))
		if rawKey == "" {
			rawKey = strings.TrimSpace(r.Header.Get("X-Idempotency-Key"))
		}

		if rawKey == "" {
			next.ServeHTTP(w, r)
			return
		}

		if len(rawKey) > 256 {
			writeAPIError(w, http.StatusBadRequest, "invalid_idempotency_key", "Idempotency-Key must not exceed 256 characters.")
			return
		}

		compoundKey := buildCompoundIdempotencyKey(r, rawKey, authUseCase)

		entry, isNew := store.LockOrGet(compoundKey)
		if !isNew {
			if entry.inFlight {
				w.Header().Set("Retry-After", "2")
				writeAPIError(w, http.StatusConflict, "idempotency_conflict", "A request with this Idempotency-Key is currently in progress. Please retry in a moment.")
				return
			}

			// Replay stored response
			for k, values := range entry.headers {
				if isHopByHopHeader(k) {
					continue
				}
				for _, v := range values {
					w.Header().Add(k, v)
				}
			}
			w.Header().Set("X-Idempotency-Replayed", "true")
			w.Header().Set("Idempotency-Key", rawKey)
			w.WriteHeader(entry.statusCode)
			if len(entry.body) > 0 {
				_, _ = w.Write(entry.body)
			}
			return
		}

		rec := &idempotencyRecorder{
			ResponseWriter: w,
			maxBodyLen:     maxIdempotencyBodyBytes,
		}

		panicked := true
		defer func() {
			if panicked || rec.statusCode >= 500 {
				store.Release(compoundKey)
			}
		}()

		next.ServeHTTP(rec, r)
		panicked = false

		if rec.statusCode == 0 {
			rec.statusCode = http.StatusOK
		}

		if rec.statusCode < 500 {
			store.Save(compoundKey, rec.statusCode, rec.Header(), rec.body.Bytes())
		}
	})
}

func isHopByHopHeader(header string) bool {
	h := strings.ToLower(header)
	return h == "connection" || h == "keep-alive" || h == "proxy-authenticate" ||
		h == "proxy-authorization" || h == "te" || h == "trailer" ||
		h == "transfer-encoding" || h == "upgrade"
}

func buildCompoundIdempotencyKey(r *http.Request, rawKey string, authUseCase *appauth.UseCase) string {
	identity := "ip:" + requestClientIP(r)
	if claims, ok := resolveAuthClaims(r, authUseCase); ok && claims.TenantID != "" {
		identity = "tenant:" + claims.TenantID + ":user:" + claims.UserID
	} else if clientID := normalizedClientID(r.Header.Get("X-Client-ID")); clientID != "" {
		identity = "device:" + clientID
	}

	raw := identity + "|" + r.Method + "|" + r.URL.Path + "|" + rawKey
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}

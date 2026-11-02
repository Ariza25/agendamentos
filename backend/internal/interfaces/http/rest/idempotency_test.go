package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func TestIdempotencyReplaysCompletedResponse(t *testing.T) {
	var callCount int32
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&callCount, 1)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Custom-Header", "value-123")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"order_id": "ord_999",
			"call":     count,
		})
	})

	mw := withIdempotency(nil, handler)

	// First request with Idempotency-Key
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/admin/checkout", strings.NewReader(`{}`))
	req1.Header.Set("Idempotency-Key", "test-key-abc-123")
	req1.RemoteAddr = "192.168.1.1:12345"
	rec1 := httptest.NewRecorder()
	mw.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec1.Code)
	}
	if rec1.Header().Get("X-Idempotency-Replayed") == "true" {
		t.Fatalf("first request should not be marked as replayed")
	}
	var res1 map[string]any
	if err := json.Unmarshal(rec1.Body.Bytes(), &res1); err != nil {
		t.Fatalf("failed to parse json response: %v", err)
	}
	if res1["call"] != float64(1) {
		t.Fatalf("expected call count 1, got %v", res1["call"])
	}

	// Second request with SAME Idempotency-Key
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/admin/checkout", strings.NewReader(`{}`))
	req2.Header.Set("Idempotency-Key", "test-key-abc-123")
	req2.RemoteAddr = "192.168.1.1:12345"
	rec2 := httptest.NewRecorder()
	mw.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusCreated {
		t.Fatalf("expected status 201 on replay, got %d", rec2.Code)
	}
	if rec2.Header().Get("X-Idempotency-Replayed") != "true" {
		t.Fatalf("expected X-Idempotency-Replayed: true header")
	}
	if rec2.Header().Get("X-Custom-Header") != "value-123" {
		t.Fatalf("expected preserved custom header, got %s", rec2.Header().Get("X-Custom-Header"))
	}
	var res2 map[string]any
	if err := json.Unmarshal(rec2.Body.Bytes(), &res2); err != nil {
		t.Fatalf("failed to parse json response: %v", err)
	}
	if res2["call"] != float64(1) {
		t.Fatalf("expected call count to remain 1, got %v", res2["call"])
	}
	if atomic.LoadInt32(&callCount) != 1 {
		t.Fatalf("handler should only have been called once, got %d", atomic.LoadInt32(&callCount))
	}
}

func TestIdempotencyBlocksInFlightConcurrentRequests(t *testing.T) {
	started := make(chan struct{})
	finish := make(chan struct{})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-finish
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})

	mw := withIdempotency(nil, handler)

	var wg sync.WaitGroup
	var rec1 *httptest.ResponseRecorder
	wg.Add(1)
	go func() {
		defer wg.Done()
		req1 := httptest.NewRequest(http.MethodPost, "/api/v1/admin/checkout", strings.NewReader(`{}`))
		req1.Header.Set("Idempotency-Key", "concurrent-key-xyz")
		req1.RemoteAddr = "10.0.0.1:44321"
		rec1 = httptest.NewRecorder()
		mw.ServeHTTP(rec1, req1)
	}()

	<-started

	// Second concurrent request with same key while first is running
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/admin/checkout", strings.NewReader(`{}`))
	req2.Header.Set("Idempotency-Key", "concurrent-key-xyz")
	req2.RemoteAddr = "10.0.0.1:44321"
	rec2 := httptest.NewRecorder()
	mw.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusConflict {
		t.Fatalf("expected HTTP 409 Conflict for in-flight request, got %d", rec2.Code)
	}

	close(finish)
	wg.Wait()

	if rec1.Code != http.StatusOK {
		t.Fatalf("original request expected 200, got %d", rec1.Code)
	}
}

func TestIdempotencyDoesNotCacheServerFailures(t *testing.T) {
	var callCount int32
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&callCount, 1)
		if count == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"transient_failure"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})

	mw := withIdempotency(nil, handler)

	// First request fails with 500
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/admin/checkout", strings.NewReader(`{}`))
	req1.Header.Set("Idempotency-Key", "retryable-key-777")
	req1.RemoteAddr = "10.0.0.2:1234"
	rec1 := httptest.NewRecorder()
	mw.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec1.Code)
	}

	// Retry should be allowed to execute again
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/admin/checkout", strings.NewReader(`{}`))
	req2.Header.Set("Idempotency-Key", "retryable-key-777")
	req2.RemoteAddr = "10.0.0.2:1234"
	rec2 := httptest.NewRecorder()
	mw.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("expected retry to succeed with 200, got %d", rec2.Code)
	}
	if atomic.LoadInt32(&callCount) != 2 {
		t.Fatalf("expected handler to be invoked twice on 500 retry, got %d", atomic.LoadInt32(&callCount))
	}
}

func TestIdempotencyIgnoredOnGetRequests(t *testing.T) {
	var callCount int32
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "data")
	})

	mw := withIdempotency(nil, handler)

	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/admin/checkout", nil)
	req1.Header.Set("Idempotency-Key", "same-key")
	rec1 := httptest.NewRecorder()
	mw.ServeHTTP(rec1, req1)

	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/admin/checkout", nil)
	req2.Header.Set("Idempotency-Key", "same-key")
	rec2 := httptest.NewRecorder()
	mw.ServeHTTP(rec2, req2)

	if atomic.LoadInt32(&callCount) != 2 {
		t.Fatalf("GET requests should not be cached by idempotency key; got %d calls", atomic.LoadInt32(&callCount))
	}
}

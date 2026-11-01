package rest

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
)

type MetricsRegistry struct {
	mu              sync.RWMutex
	requestCounts   map[string]*int64 // key: "method|path|status"
	requestDuration map[string]*int64 // key: "method|path", sum in milliseconds
	broker          *EventBroker
}

var globalMetrics = &MetricsRegistry{
	requestCounts:   make(map[string]*int64),
	requestDuration: make(map[string]*int64),
}

func RegisterMetricsBroker(b *EventBroker) {
	globalMetrics.broker = b
}

func RecordRequestMetric(method string, path string, status int, durationMs int64) {
	key := fmt.Sprintf("%s|%s|%d", method, path, status)
	durKey := fmt.Sprintf("%s|%s", method, path)

	globalMetrics.mu.Lock()
	counter, exists := globalMetrics.requestCounts[key]
	if !exists {
		var initial int64
		counter = &initial
		globalMetrics.requestCounts[key] = counter
	}
	durCounter, durExists := globalMetrics.requestDuration[durKey]
	if !durExists {
		var initialDur int64
		durCounter = &initialDur
		globalMetrics.requestDuration[durKey] = durCounter
	}
	globalMetrics.mu.Unlock()

	atomic.AddInt64(counter, 1)
	atomic.AddInt64(durCounter, durationMs)
}

func HandleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

	globalMetrics.mu.RLock()
	defer globalMetrics.mu.RUnlock()

	fmt.Fprintln(w, "# HELP erpagendamentos_http_requests_total Total number of HTTP requests processed")
	fmt.Fprintln(w, "# TYPE erpagendamentos_http_requests_total counter")
	for key, counter := range globalMetrics.requestCounts {
		var method, path string
		var status int
		fmt.Sscanf(key, "%s|%s|%d", &method, &path, &status)
		val := atomic.LoadInt64(counter)
		fmt.Fprintf(w, "erpagendamentos_http_requests_total{method=%q,path=%q,status=\"%d\"} %d\n", method, path, status, val)
	}

	fmt.Fprintln(w, "# HELP erpagendamentos_http_request_duration_ms_total Total duration of HTTP requests in milliseconds")
	fmt.Fprintln(w, "# TYPE erpagendamentos_http_request_duration_ms_total counter")
	for key, durCounter := range globalMetrics.requestDuration {
		var method, path string
		fmt.Sscanf(key, "%s|%s", &method, &path)
		val := atomic.LoadInt64(durCounter)
		fmt.Fprintf(w, "erpagendamentos_http_request_duration_ms_total{method=%q,path=%q} %d\n", method, path, val)
	}

	if globalMetrics.broker != nil {
		fmt.Fprintln(w, "# HELP erpagendamentos_sse_connected_clients Current number of active SSE client connections")
		fmt.Fprintln(w, "# TYPE erpagendamentos_sse_connected_clients gauge")
		fmt.Fprintf(w, "erpagendamentos_sse_connected_clients %d\n", globalMetrics.broker.ActiveClientsCount(""))
	}
}


package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"erpagendamentos/backend/internal/bootstrap"
)

type lazyAppHandler struct {
	mu        sync.Mutex
	container *bootstrap.Container
	handler   http.Handler
}

func newLazyAppHandler() *lazyAppHandler {
	return &lazyAppHandler{}
}

func (h *lazyAppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Path == "/health" {
		writeLazyJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		return
	}

	h.init().ServeHTTP(w, r)
}

func (h *lazyAppHandler) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.container == nil {
		return nil
	}
	return h.container.Close()
}

func (h *lazyAppHandler) init() http.Handler {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.handler != nil {
		return h.handler
	}
	container := bootstrap.NewContainer()
	h.container = container
	h.handler = container.Router().Handler()
	return h.handler
}

func writeLazyJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("write response: %v", err)
	}
}

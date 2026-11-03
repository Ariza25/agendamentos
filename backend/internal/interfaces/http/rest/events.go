package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Event struct {
	Type      string `json:"type"`
	TenantID  string `json:"tenantId"`
	Data      any    `json:"data"`
	Timestamp string `json:"timestamp"`
}

type EventBroker struct {
	mu      sync.RWMutex
	clients map[string]map[chan Event]struct{} // tenantID -> set of client channels
}

func NewEventBroker() *EventBroker {
	return &EventBroker{
		clients: make(map[string]map[chan Event]struct{}),
	}
}

func (b *EventBroker) Subscribe(tenantID string) chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan Event, 32)
	if _, exists := b.clients[tenantID]; !exists {
		b.clients[tenantID] = make(map[chan Event]struct{})
	}
	b.clients[tenantID][ch] = struct{}{}
	return ch
}

func (b *EventBroker) Unsubscribe(tenantID string, ch chan Event) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if tenantClients, exists := b.clients[tenantID]; exists {
		delete(tenantClients, ch)
		close(ch)
		if len(tenantClients) == 0 {
			delete(b.clients, tenantID)
		}
	}
}

func (b *EventBroker) Publish(tenantID string, eventType string, data any) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	tenantClients, exists := b.clients[tenantID]
	if !exists || len(tenantClients) == 0 {
		return
	}

	event := Event{
		Type:      eventType,
		TenantID:  tenantID,
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	for ch := range tenantClients {
		select {
		case ch <- event:
		default:
			// Client channel full, skip to prevent blocking
		}
	}
}

func (b *EventBroker) ActiveClientsCount(tenantID string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if tenantID != "" {
		return len(b.clients[tenantID])
	}
	total := 0
	for _, clients := range b.clients {
		total += len(clients)
	}
	return total
}

func (b *EventBroker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	tenantID := resolveRequestTenantID(r, nil)
	if tenantID == "" {
		writeAPIError(w, http.StatusBadRequest, "missing_tenant", "Tenant context required for event stream.")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	eventCh := b.Subscribe(tenantID)
	defer b.Unsubscribe(tenantID, eventCh)

	// Send initial connection event
	fmt.Fprintf(w, "event: connected\ndata: {\"tenantId\":%q,\"timestamp\":%q}\n\n", tenantID, time.Now().UTC().Format(time.RFC3339))
	flusher.Flush()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	notify := r.Context().Done()
	for {
		select {
		case <-notify:
			return
		case <-ticker.C:
			fmt.Fprintf(w, ": heartbeat\n\n")
			flusher.Flush()
		case evt, ok := <-eventCh:
			if !ok {
				return
			}
			payload, err := json.Marshal(evt)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", evt.Type, string(payload))
			flusher.Flush()
		}
	}
}

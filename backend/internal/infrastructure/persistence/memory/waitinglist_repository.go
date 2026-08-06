package memory

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/waitinglist"
)

type WaitingListRepository struct {
	mu      sync.Mutex
	entries map[string]waitinglist.Entry
}

func NewWaitingListRepository() *WaitingListRepository {
	return &WaitingListRepository{
		entries: make(map[string]waitinglist.Entry),
	}
}

func (r *WaitingListRepository) Create(input waitinglist.Entry) (waitinglist.Entry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item := input
	item.ID = uuid.NewString()
	item.CreatedAt = time.Now().Format(time.RFC3339)
	if item.Status == "" {
		item.Status = waitinglist.StatusWaiting
	}
	r.entries[item.ID] = item
	return item, nil
}

func (r *WaitingListRepository) GetByID(tenantID string, id string) (waitinglist.Entry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.entries[strings.TrimSpace(id)]
	if !exists || item.TenantID != strings.TrimSpace(tenantID) {
		return waitinglist.Entry{}, fmt.Errorf("waiting list entry not found")
	}
	return item, nil
}

func (r *WaitingListRepository) ListByTenant(tenantID string, status string) ([]waitinglist.Entry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID = strings.TrimSpace(tenantID)
	status = strings.TrimSpace(status)
	items := []waitinglist.Entry{}
	for _, item := range r.entries {
		if item.TenantID != tenantID {
			continue
		}
		if status != "" && item.Status != status {
			continue
		}
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt < items[j].CreatedAt
	})
	return items, nil
}

func (r *WaitingListRepository) UpdateStatus(tenantID string, id string, status string) (waitinglist.Entry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.entries[strings.TrimSpace(id)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return waitinglist.Entry{}, fmt.Errorf("waiting list entry not found")
	}
	current.Status = status
	if status == waitinglist.StatusNotified {
		current.NotifiedAt = time.Now().Format(time.RFC3339)
	}
	r.entries[current.ID] = current
	return current, nil
}

func (r *WaitingListRepository) ListWaitingForSlot(tenantID string, serviceID string, professionalID string, order string) ([]waitinglist.Entry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID = strings.TrimSpace(tenantID)
	serviceID = strings.TrimSpace(serviceID)
	professionalID = strings.TrimSpace(professionalID)
	items := []waitinglist.Entry{}
	for _, item := range r.entries {
		if item.TenantID != tenantID || item.Status != waitinglist.StatusWaiting {
			continue
		}
		if serviceID != "" && item.ServiceID != "" && item.ServiceID != serviceID {
			continue
		}
		if professionalID != "" && item.ProfessionalID != "" && item.ProfessionalID != professionalID {
			continue
		}
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		if order == "oldest" {
			return items[i].CreatedAt < items[j].CreatedAt
		}
		return items[i].CreatedAt > items[j].CreatedAt
	})
	return items, nil
}

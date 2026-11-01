package memory

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/cashregister"
)

type CashRegisterRepository struct {
	mu        sync.Mutex
	registers map[string]cashregister.CashRegister
	entries   map[string][]cashregister.CashRegisterEntry
}

func NewCashRegisterRepository() *CashRegisterRepository {
	return &CashRegisterRepository{
		registers: make(map[string]cashregister.CashRegister),
		entries:   make(map[string][]cashregister.CashRegisterEntry),
	}
}

func (r *CashRegisterRepository) Create(item cashregister.CashRegister) (cashregister.CashRegister, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if item.ID == "" {
		item.ID = uuid.NewString()
	}
	item.OpenedAt = time.Now().Format(time.RFC3339)
	item.Status = "open"
	item.SystemBalanceCents = item.InitialBalanceCents
	r.registers[item.ID] = item
	return item, nil
}

func (r *CashRegisterRepository) Update(item cashregister.CashRegister) (cashregister.CashRegister, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.registers[item.ID]
	if !exists {
		return cashregister.CashRegister{}, errors.New("cash register not found")
	}
	if item.Status == "closed" && item.ClosedAt == "" {
		item.ClosedAt = time.Now().Format(time.RFC3339)
	}
	r.registers[item.ID] = item
	return item, nil
}

func (r *CashRegisterRepository) GetByID(tenantID string, id string) (cashregister.CashRegister, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.registers[id]
	if !exists || item.TenantID != tenantID {
		return cashregister.CashRegister{}, errors.New("cash register not found")
	}
	return item, nil
}

func (r *CashRegisterRepository) GetCurrentOpen(tenantID string) (cashregister.CashRegister, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, item := range r.registers {
		if item.TenantID == tenantID && item.Status == "open" {
			return item, nil
		}
	}
	return cashregister.CashRegister{}, errors.New("no open cash register")
}

func (r *CashRegisterRepository) List(tenantID string, limit int) ([]cashregister.CashRegister, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var list []cashregister.CashRegister
	for _, item := range r.registers {
		if item.TenantID == tenantID {
			list = append(list, item)
		}
	}
	return list, nil
}

func (r *CashRegisterRepository) CreateEntry(entry cashregister.CashRegisterEntry) (cashregister.CashRegisterEntry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if entry.ID == "" {
		entry.ID = uuid.NewString()
	}
	entry.CreatedAt = time.Now().Format(time.RFC3339)
	r.entries[entry.CashRegisterID] = append(r.entries[entry.CashRegisterID], entry)
	return entry, nil
}

func (r *CashRegisterRepository) ListEntries(tenantID string, cashRegisterID string) ([]cashregister.CashRegisterEntry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	list := r.entries[cashRegisterID]
	var res []cashregister.CashRegisterEntry
	for _, e := range list {
		if e.TenantID == tenantID {
			res = append(res, e)
		}
	}
	return res, nil
}


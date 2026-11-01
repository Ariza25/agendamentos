package memory

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"erpagendamentos/backend/internal/domain/stock"
	"github.com/google/uuid"
)

type StockRepository struct {
	mu        sync.RWMutex
	products  map[string]stock.Product
	supplies  map[string][]stock.ServiceSupply // key: tenantID:serviceID
	movements []stock.StockMovement
}

func NewStockRepository() *StockRepository {
	return &StockRepository{
		products: make(map[string]stock.Product),
		supplies: make(map[string][]stock.ServiceSupply),
	}
}

func (r *StockRepository) CreateProduct(input stock.Product) (stock.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if input.ID == "" {
		input.ID = uuid.NewString()
	}
	now := time.Now().UTC().Format(time.RFC3339)
	input.CreatedAt = now
	input.UpdatedAt = now
	key := input.TenantID + ":" + input.ID
	r.products[key] = input
	return input, nil
}

func (r *StockRepository) UpdateProduct(input stock.Product) (stock.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := input.TenantID + ":" + input.ID
	existing, found := r.products[key]
	if !found {
		return stock.Product{}, errors.New("product not found")
	}
	input.CreatedAt = existing.CreatedAt
	input.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	r.products[key] = input
	return input, nil
}

func (r *StockRepository) GetProductByID(tenantID string, id string) (stock.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := tenantID + ":" + id
	p, found := r.products[key]
	if !found {
		return stock.Product{}, errors.New("product not found")
	}
	return p, nil
}

func (r *StockRepository) ListProducts(tenantID string, activeOnly bool) ([]stock.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []stock.Product
	for _, p := range r.products {
		if p.TenantID == tenantID {
			if activeOnly && !p.Active {
				continue
			}
			result = append(result, p)
		}
	}
	return result, nil
}

func (r *StockRepository) DeleteProduct(tenantID string, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.products, tenantID+":"+id)
	return nil
}

func (r *StockRepository) CreateStockMovement(input stock.StockMovement) (stock.StockMovement, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if input.ID == "" {
		input.ID = uuid.NewString()
	}
	input.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	// Adjust product stock
	key := input.TenantID + ":" + input.ProductID
	if p, found := r.products[key]; found {
		switch input.Type {
		case "in", "adjustment_up":
			p.StockQuantity += input.Quantity
		case "out", "consumption", "sale", "adjustment_down":
			p.StockQuantity -= input.Quantity
		case "adjustment":
			p.StockQuantity = input.Quantity
		}
		p.UpdatedAt = input.CreatedAt
		r.products[key] = p
		input.ProductName = p.Name
	}

	r.movements = append(r.movements, input)
	return input, nil
}

func (r *StockRepository) ListStockMovements(tenantID string, productID string, limit int) ([]stock.StockMovement, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []stock.StockMovement
	for i := len(r.movements) - 1; i >= 0; i-- {
		m := r.movements[i]
		if m.TenantID == tenantID {
			if productID != "" && m.ProductID != productID {
				continue
			}
			result = append(result, m)
			if limit > 0 && len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func (r *StockRepository) SetServiceSupplies(tenantID string, serviceID string, supplies []stock.ServiceSupply) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC().Format(time.RFC3339)
	for i := range supplies {
		if supplies[i].ID == "" {
			supplies[i].ID = uuid.NewString()
		}
		supplies[i].TenantID = tenantID
		supplies[i].ServiceID = serviceID
		supplies[i].CreatedAt = now
		if p, found := r.products[tenantID+":"+supplies[i].ProductID]; found {
			supplies[i].ProductName = p.Name
		}
	}
	r.supplies[tenantID+":"+serviceID] = supplies
	return nil
}

func (r *StockRepository) ListServiceSupplies(tenantID string, serviceID string) ([]stock.ServiceSupply, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return append([]stock.ServiceSupply(nil), r.supplies[tenantID+":"+serviceID]...), nil
}

func (r *StockRepository) DeductSuppliesForAppointment(tenantID string, appointmentID string, serviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	supplies := r.supplies[tenantID+":"+serviceID]
	for _, s := range supplies {
		key := tenantID + ":" + s.ProductID
		if p, found := r.products[key]; found {
			p.StockQuantity -= s.QuantityConsumed
			p.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
			r.products[key] = p

			mov := stock.StockMovement{
				ID:            uuid.NewString(),
				TenantID:      tenantID,
				ProductID:     s.ProductID,
				ProductName:   p.Name,
				Type:          "consumption",
				Quantity:      s.QuantityConsumed,
				UnitCostCents: p.CostCents,
				Notes:         fmt.Sprintf("Auto-deduction for appointment %s", appointmentID),
				AppointmentID: appointmentID,
				CreatedAt:     time.Now().UTC().Format(time.RFC3339),
			}
			r.movements = append(r.movements, mov)
		}
	}
	return nil
}

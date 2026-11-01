package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"erpagendamentos/backend/internal/domain/stock"
	"github.com/google/uuid"
)

type StockRepository struct {
	db *sql.DB
}

func NewStockRepository(db *sql.DB) *StockRepository {
	return &StockRepository{db: db}
}

func (r *StockRepository) CreateProduct(input stock.Product) (stock.Product, error) {
	if input.ID == "" {
		input.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	err := r.db.QueryRow(`
		INSERT INTO catalog.products (id, tenant_id, name, sku, category, price_cents, cost_cents, stock_quantity, min_stock_alert, is_supply, unit, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING created_at, updated_at
	`, input.ID, input.TenantID, input.Name, input.SKU, input.Category, input.PriceCents, input.CostCents, input.StockQuantity, input.MinStockAlert, input.IsSupply, input.Unit, input.Active, now, now).Scan(&input.CreatedAt, &input.UpdatedAt)
	if err != nil {
		return stock.Product{}, fmt.Errorf("create product: %w", err)
	}
	return input, nil
}

func (r *StockRepository) UpdateProduct(input stock.Product) (stock.Product, error) {
	now := time.Now().UTC()
	result, err := r.db.Exec(`
		UPDATE catalog.products
		SET name = $1, sku = $2, category = $3, price_cents = $4, cost_cents = $5, stock_quantity = $6, min_stock_alert = $7, is_supply = $8, unit = $9, active = $10, updated_at = $11
		WHERE tenant_id = $12 AND id = $13
	`, input.Name, input.SKU, input.Category, input.PriceCents, input.CostCents, input.StockQuantity, input.MinStockAlert, input.IsSupply, input.Unit, input.Active, now, input.TenantID, input.ID)
	if err != nil {
		return stock.Product{}, fmt.Errorf("update product: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return stock.Product{}, errors.New("product not found")
	}
	input.UpdatedAt = now.Format(time.RFC3339)
	return input, nil
}

func (r *StockRepository) GetProductByID(tenantID string, id string) (stock.Product, error) {
	var p stock.Product
	var created, updated time.Time
	err := r.db.QueryRow(`
		SELECT id, tenant_id, name, sku, category, price_cents, cost_cents, stock_quantity, min_stock_alert, is_supply, unit, active, created_at, updated_at
		FROM catalog.products
		WHERE tenant_id = $1 AND id = $2
	`, tenantID, id).Scan(&p.ID, &p.TenantID, &p.Name, &p.SKU, &p.Category, &p.PriceCents, &p.CostCents, &p.StockQuantity, &p.MinStockAlert, &p.IsSupply, &p.Unit, &p.Active, &created, &updated)
	if err != nil {
		return stock.Product{}, fmt.Errorf("get product: %w", err)
	}
	p.CreatedAt = created.Format(time.RFC3339)
	p.UpdatedAt = updated.Format(time.RFC3339)
	return p, nil
}

func (r *StockRepository) ListProducts(tenantID string, activeOnly bool) ([]stock.Product, error) {
	query := `
		SELECT id, tenant_id, name, sku, category, price_cents, cost_cents, stock_quantity, min_stock_alert, is_supply, unit, active, created_at, updated_at
		FROM catalog.products
		WHERE tenant_id = $1
	`
	if activeOnly {
		query += " AND active = true"
	}
	query += " ORDER BY name ASC"

	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()

	var products []stock.Product
	for rows.Next() {
		var p stock.Product
		var created, updated time.Time
		if err := rows.Scan(&p.ID, &p.TenantID, &p.Name, &p.SKU, &p.Category, &p.PriceCents, &p.CostCents, &p.StockQuantity, &p.MinStockAlert, &p.IsSupply, &p.Unit, &p.Active, &created, &updated); err != nil {
			return nil, err
		}
		p.CreatedAt = created.Format(time.RFC3339)
		p.UpdatedAt = updated.Format(time.RFC3339)
		products = append(products, p)
	}
	return products, nil
}

func (r *StockRepository) DeleteProduct(tenantID string, id string) error {
	_, err := r.db.Exec("DELETE FROM catalog.products WHERE tenant_id = $1 AND id = $2", tenantID, id)
	return err
}

func (r *StockRepository) CreateStockMovement(input stock.StockMovement) (stock.StockMovement, error) {
	if input.ID == "" {
		input.ID = uuid.NewString()
	}
	now := time.Now().UTC()

	tx, err := r.db.Begin()
	if err != nil {
		return stock.StockMovement{}, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Adjust quantity in products
	var delta float64
	switch input.Type {
	case "in", "adjustment_up":
		delta = input.Quantity
	case "out", "consumption", "sale", "adjustment_down":
		delta = -input.Quantity
	case "adjustment":
		// handled specially
	}

	if input.Type == "adjustment" {
		_, err = tx.Exec("UPDATE catalog.products SET stock_quantity = $1, updated_at = $2 WHERE tenant_id = $3 AND id = $4", input.Quantity, now, input.TenantID, input.ProductID)
	} else {
		_, err = tx.Exec("UPDATE catalog.products SET stock_quantity = stock_quantity + $1, updated_at = $2 WHERE tenant_id = $3 AND id = $4", delta, now, input.TenantID, input.ProductID)
	}
	if err != nil {
		return stock.StockMovement{}, fmt.Errorf("update product stock: %w", err)
	}

	err = tx.QueryRow(`
		INSERT INTO catalog.stock_movements (id, tenant_id, product_id, type, quantity, unit_cost_cents, notes, appointment_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at
	`, input.ID, input.TenantID, input.ProductID, input.Type, input.Quantity, input.UnitCostCents, input.Notes, input.AppointmentID, now).Scan(&input.CreatedAt)
	if err != nil {
		return stock.StockMovement{}, fmt.Errorf("insert stock movement: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return stock.StockMovement{}, err
	}
	return input, nil
}

func (r *StockRepository) ListStockMovements(tenantID string, productID string, limit int) ([]stock.StockMovement, error) {
	query := `
		SELECT sm.id, sm.tenant_id, sm.product_id, p.name, sm.type, sm.quantity, sm.unit_cost_cents, sm.notes, sm.appointment_id, sm.created_at
		FROM catalog.stock_movements sm
		JOIN catalog.products p ON sm.product_id = p.id
		WHERE sm.tenant_id = $1
	`
	args := []any{tenantID}
	if productID != "" {
		query += " AND sm.product_id = $2"
		args = append(args, productID)
	}
	query += " ORDER BY sm.created_at DESC"
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list stock movements: %w", err)
	}
	defer rows.Close()

	var movements []stock.StockMovement
	for rows.Next() {
		var m stock.StockMovement
		var created time.Time
		if err := rows.Scan(&m.ID, &m.TenantID, &m.ProductID, &m.ProductName, &m.Type, &m.Quantity, &m.UnitCostCents, &m.Notes, &m.AppointmentID, &created); err != nil {
			return nil, err
		}
		m.CreatedAt = created.Format(time.RFC3339)
		movements = append(movements, m)
	}
	return movements, nil
}

func (r *StockRepository) SetServiceSupplies(tenantID string, serviceID string, supplies []stock.ServiceSupply) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err = tx.Exec("DELETE FROM catalog.service_supplies WHERE tenant_id = $1 AND service_id = $2", tenantID, serviceID); err != nil {
		return fmt.Errorf("delete old supplies: %w", err)
	}

	now := time.Now().UTC()
	for _, s := range supplies {
		id := s.ID
		if id == "" {
			id = uuid.NewString()
		}
		if _, err = tx.Exec(`
			INSERT INTO catalog.service_supplies (id, tenant_id, service_id, product_id, quantity_consumed, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, id, tenantID, serviceID, s.ProductID, s.QuantityConsumed, now); err != nil {
			return fmt.Errorf("insert service supply: %w", err)
		}
	}
	return tx.Commit()
}

func (r *StockRepository) ListServiceSupplies(tenantID string, serviceID string) ([]stock.ServiceSupply, error) {
	rows, err := r.db.Query(`
		SELECT ss.id, ss.tenant_id, ss.service_id, ss.product_id, p.name, ss.quantity_consumed, ss.created_at
		FROM catalog.service_supplies ss
		JOIN catalog.products p ON ss.product_id = p.id
		WHERE ss.tenant_id = $1 AND ss.service_id = $2
	`, tenantID, serviceID)
	if err != nil {
		return nil, fmt.Errorf("list service supplies: %w", err)
	}
	defer rows.Close()

	var supplies []stock.ServiceSupply
	for rows.Next() {
		var s stock.ServiceSupply
		var created time.Time
		if err := rows.Scan(&s.ID, &s.TenantID, &s.ServiceID, &s.ProductID, &s.ProductName, &s.QuantityConsumed, &created); err != nil {
			return nil, err
		}
		s.CreatedAt = created.Format(time.RFC3339)
		supplies = append(supplies, s)
	}
	return supplies, nil
}

func (r *StockRepository) DeductSuppliesForAppointment(tenantID string, appointmentID string, serviceID string) error {
	supplies, err := r.ListServiceSupplies(tenantID, serviceID)
	if err != nil || len(supplies) == 0 {
		return err
	}

	for _, s := range supplies {
		_, err := r.CreateStockMovement(stock.StockMovement{
			TenantID:      tenantID,
			ProductID:     s.ProductID,
			Type:          "consumption",
			Quantity:      s.QuantityConsumed,
			Notes:         fmt.Sprintf("Auto-deduction for appointment %s", appointmentID),
			AppointmentID: appointmentID,
		})
		if err != nil {
			return err
		}
	}
	return nil
}


package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/service"
)

type ServiceRepository struct {
	db *sql.DB
}

func NewServiceRepository(db *sql.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) Create(input service.Service) (service.Service, error) {
	item := input
	item.ID = uuid.NewString()
	item.Name = strings.TrimSpace(item.Name)
	item.Description = strings.TrimSpace(item.Description)
	if _, err := r.db.Exec(`
		INSERT INTO catalog.services (id, tenant_id, name, description, duration_minutes, price_cents, color, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, item.ID, item.TenantID, item.Name, item.Description, item.DurationMinutes, item.PriceCents, item.Color, item.Active); err != nil {
		return service.Service{}, fmt.Errorf("insert service: %w", err)
	}
	return r.GetByID(item.TenantID, item.ID)
}

func (r *ServiceRepository) Update(input service.Service) (service.Service, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.services
		SET name = COALESCE(NULLIF($3, ''), name),
		    description = $4,
		    duration_minutes = $5,
		    price_cents = $6,
		    color = COALESCE(NULLIF($7, ''), color),
		    active = $8,
		    updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`, input.ID, input.TenantID, strings.TrimSpace(input.Name), strings.TrimSpace(input.Description), input.DurationMinutes, input.PriceCents, strings.TrimSpace(input.Color), input.Active)
	if err != nil {
		return service.Service{}, fmt.Errorf("update service: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return service.Service{}, fmt.Errorf("check service update: %w", err)
	}
	if affected == 0 {
		return service.Service{}, fmt.Errorf("service not found")
	}
	return r.GetByID(input.TenantID, input.ID)
}

func (r *ServiceRepository) GetByID(tenantID string, id string) (service.Service, error) {
	item, err := scanServiceRow(r.db.QueryRow(`
		SELECT id, tenant_id, name, description, duration_minutes, price_cents, color, active, created_at, updated_at
		FROM catalog.services
		WHERE id = $1 AND tenant_id = $2
	`, strings.TrimSpace(id), strings.TrimSpace(tenantID)))
	if err != nil {
		return service.Service{}, err
	}
	return item, nil
}

func (r *ServiceRepository) ListByTenant(tenantID string) ([]service.Service, error) {
	rows, err := r.db.Query(`
		SELECT id, tenant_id, name, description, duration_minutes, price_cents, color, active, created_at, updated_at
		FROM catalog.services
		WHERE tenant_id = $1
		ORDER BY name ASC
	`, strings.TrimSpace(tenantID))
	if err != nil {
		return nil, fmt.Errorf("query services: %w", err)
	}
	defer rows.Close()

	items := []service.Service{}
	for rows.Next() {
		item, err := scanServiceRows(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *ServiceRepository) SetActive(tenantID string, id string, active bool) (service.Service, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.services
		SET active = $3, updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`, id, tenantID, active)
	if err != nil {
		return service.Service{}, fmt.Errorf("set service active: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return service.Service{}, fmt.Errorf("check service active: %w", err)
	}
	if affected == 0 {
		return service.Service{}, fmt.Errorf("service not found")
	}
	return r.GetByID(tenantID, id)
}

func (r *ServiceRepository) Delete(tenantID string, id string) error {
	result, err := r.db.Exec(`
		DELETE FROM catalog.services
		WHERE id = $1 AND tenant_id = $2
	`, id, tenantID)
	if err != nil {
		return fmt.Errorf("delete service: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check service delete: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("service not found")
	}
	return nil
}

type serviceScanner interface {
	Scan(dest ...any) error
}

func scanServiceRow(scanner serviceScanner) (service.Service, error) {
	return scanServiceRows(scanner)
}

func scanServiceRows(scanner serviceScanner) (service.Service, error) {
	var item service.Service
	var createdAt, updatedAt time.Time
	if err := scanner.Scan(&item.ID, &item.TenantID, &item.Name, &item.Description, &item.DurationMinutes, &item.PriceCents, &item.Color, &item.Active, &createdAt, &updatedAt); err != nil {
		if err == sql.ErrNoRows {
			return service.Service{}, fmt.Errorf("service not found")
		}
		return service.Service{}, fmt.Errorf("scan service: %w", err)
	}
	item.CreatedAt = createdAt.Format(time.RFC3339)
	item.UpdatedAt = updatedAt.Format(time.RFC3339)
	return item, nil
}

package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/resource"
)

type ResourceRepository struct {
	db *sql.DB
}

func NewResourceRepository(db *sql.DB) *ResourceRepository {
	return &ResourceRepository{db: db}
}

func (r *ResourceRepository) Create(input resource.Resource) (resource.Resource, error) {
	item := input
	item.ID = uuid.NewString()
	item.Name = strings.TrimSpace(item.Name)
	item.Type = strings.TrimSpace(item.Type)
	item.Notes = strings.TrimSpace(item.Notes)
	if _, err := r.db.Exec(`
		INSERT INTO catalog.resources (id, tenant_id, name, type, capacity, active, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, item.ID, item.TenantID, item.Name, item.Type, item.Capacity, item.Active, item.Notes); err != nil {
		return resource.Resource{}, fmt.Errorf("insert resource: %w", err)
	}
	return r.GetByID(item.TenantID, item.ID)
}

func (r *ResourceRepository) Update(input resource.Resource) (resource.Resource, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.resources
		SET name = COALESCE(NULLIF($3, ''), name),
		    type = COALESCE(NULLIF($4, ''), type),
		    capacity = $5,
		    active = $6,
		    notes = $7,
		    updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`, input.ID, input.TenantID, strings.TrimSpace(input.Name), strings.TrimSpace(input.Type), input.Capacity, input.Active, strings.TrimSpace(input.Notes))
	if err != nil {
		return resource.Resource{}, fmt.Errorf("update resource: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return resource.Resource{}, fmt.Errorf("check resource update: %w", err)
	}
	if affected == 0 {
		return resource.Resource{}, fmt.Errorf("resource not found")
	}
	return r.GetByID(input.TenantID, input.ID)
}

func (r *ResourceRepository) GetByID(tenantID string, id string) (resource.Resource, error) {
	return scanResourceRow(r.db.QueryRow(`
		SELECT id, tenant_id, name, type, capacity, active, notes, created_at, updated_at
		FROM catalog.resources
		WHERE id = $1 AND tenant_id = $2
	`, strings.TrimSpace(id), strings.TrimSpace(tenantID)))
}

func (r *ResourceRepository) ListByTenant(tenantID string) ([]resource.Resource, error) {
	rows, err := r.db.Query(`
		SELECT id, tenant_id, name, type, capacity, active, notes, created_at, updated_at
		FROM catalog.resources
		WHERE tenant_id = $1
		ORDER BY name ASC
	`, strings.TrimSpace(tenantID))
	if err != nil {
		return nil, fmt.Errorf("query resources: %w", err)
	}
	defer rows.Close()

	items := []resource.Resource{}
	for rows.Next() {
		item, err := scanResourceRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *ResourceRepository) SetActive(tenantID string, id string, active bool) (resource.Resource, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.resources
		SET active = $3, updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`, id, tenantID, active)
	if err != nil {
		return resource.Resource{}, fmt.Errorf("set resource active: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return resource.Resource{}, fmt.Errorf("check resource active: %w", err)
	}
	if affected == 0 {
		return resource.Resource{}, fmt.Errorf("resource not found")
	}
	return r.GetByID(tenantID, id)
}

func (r *ResourceRepository) Delete(tenantID string, id string) error {
	result, err := r.db.Exec(`
		DELETE FROM catalog.resources
		WHERE id = $1 AND tenant_id = $2
	`, id, tenantID)
	if err != nil {
		return fmt.Errorf("delete resource: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check resource delete: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("resource not found")
	}
	return nil
}

type resourceScanner interface {
	Scan(dest ...any) error
}

func scanResourceRow(scanner resourceScanner) (resource.Resource, error) {
	var item resource.Resource
	var createdAt, updatedAt time.Time
	if err := scanner.Scan(&item.ID, &item.TenantID, &item.Name, &item.Type, &item.Capacity, &item.Active, &item.Notes, &createdAt, &updatedAt); err != nil {
		if err == sql.ErrNoRows {
			return resource.Resource{}, fmt.Errorf("resource not found")
		}
		return resource.Resource{}, fmt.Errorf("scan resource: %w", err)
	}
	item.CreatedAt = createdAt.Format(time.RFC3339)
	item.UpdatedAt = updatedAt.Format(time.RFC3339)
	return item, nil
}

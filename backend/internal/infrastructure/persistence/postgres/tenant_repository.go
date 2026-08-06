package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"erpagendamentos/backend/internal/domain/tenant"
)

type TenantRepository struct {
	db *sql.DB
}

func NewTenantRepository(db *sql.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) Create(input tenant.Tenant) (tenant.Tenant, error) {
	item := input
	if strings.TrimSpace(item.ID) == "" {
		return tenant.Tenant{}, fmt.Errorf("tenant id is required")
	}
	item.Slug = strings.ToLower(strings.TrimSpace(item.Slug))
	if _, err := r.db.Exec(`
		INSERT INTO catalog.tenants (id, name, slug, active)
		VALUES ($1, $2, $3, $4)
	`, item.ID, item.Name, item.Slug, item.Active); err != nil {
		if isTenantIdentityConflict(err) {
			return tenant.Tenant{}, fmt.Errorf("slug already used")
		}
		return tenant.Tenant{}, fmt.Errorf("insert tenant: %w", err)
	}
	return r.GetByID(item.ID)
}

func (r *TenantRepository) GetByID(id string) (tenant.Tenant, error) {
	item, err := scanTenantRow(r.db.QueryRow(`
		SELECT id, name, slug, active, created_at, updated_at, deleted_at
		FROM catalog.tenants
		WHERE id = $1
	`, strings.TrimSpace(id)))
	if err != nil {
		return tenant.Tenant{}, err
	}
	return item, nil
}

func (r *TenantRepository) GetBySlug(slug string) (tenant.Tenant, error) {
	item, err := scanTenantRow(r.db.QueryRow(`
		SELECT id, name, slug, active, created_at, updated_at, deleted_at
		FROM catalog.tenants
		WHERE LOWER(BTRIM(slug)) = LOWER($1)
	`, strings.TrimSpace(slug)))
	if err != nil {
		return tenant.Tenant{}, err
	}
	return item, nil
}

func (r *TenantRepository) List() ([]tenant.Tenant, error) {
	rows, err := r.db.Query(`
		SELECT id, name, slug, active, created_at, updated_at, deleted_at
		FROM catalog.tenants
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("query tenants: %w", err)
	}
	defer rows.Close()

	items := []tenant.Tenant{}
	for rows.Next() {
		item, err := scanTenantRows(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *TenantRepository) Update(input tenant.Tenant) (tenant.Tenant, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.tenants
		SET name = $2, slug = $3, updated_at = NOW()
		WHERE id = $1
	`, input.ID, input.Name, strings.ToLower(strings.TrimSpace(input.Slug)))
	if err != nil {
		if isTenantIdentityConflict(err) {
			return tenant.Tenant{}, fmt.Errorf("slug already used")
		}
		return tenant.Tenant{}, fmt.Errorf("update tenant: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return tenant.Tenant{}, fmt.Errorf("check tenant update: %w", err)
	}
	if affected == 0 {
		return tenant.Tenant{}, fmt.Errorf("tenant not found")
	}
	return r.GetByID(input.ID)
}

func (r *TenantRepository) SetActive(id string, active bool) (tenant.Tenant, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.tenants
		SET active = $2, updated_at = NOW()
		WHERE id = $1
	`, id, active)
	if err != nil {
		return tenant.Tenant{}, fmt.Errorf("set tenant active: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return tenant.Tenant{}, fmt.Errorf("check tenant active: %w", err)
	}
	if affected == 0 {
		return tenant.Tenant{}, fmt.Errorf("tenant not found")
	}
	return r.GetByID(id)
}

func (r *TenantRepository) Delete(id string) error {
	result, err := r.db.Exec(`
		UPDATE catalog.tenants
		SET active = FALSE, deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("delete tenant: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check tenant delete: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("tenant not found")
	}
	return nil
}

type tenantScanner interface {
	Scan(dest ...any) error
}

func scanTenantRow(scanner tenantScanner) (tenant.Tenant, error) {
	return scanTenantRows(scanner)
}

func scanTenantRows(scanner tenantScanner) (tenant.Tenant, error) {
	var item tenant.Tenant
	var createdAt, updatedAt time.Time
	var deletedAt sql.NullTime
	if err := scanner.Scan(&item.ID, &item.Name, &item.Slug, &item.Active, &createdAt, &updatedAt, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return tenant.Tenant{}, fmt.Errorf("tenant not found")
		}
		return tenant.Tenant{}, fmt.Errorf("scan tenant: %w", err)
	}
	item.CreatedAt = createdAt.Format(time.RFC3339)
	item.UpdatedAt = updatedAt.Format(time.RFC3339)
	if deletedAt.Valid {
		item.DeletedAt = deletedAt.Time.Format(time.RFC3339)
	}
	return item, nil
}

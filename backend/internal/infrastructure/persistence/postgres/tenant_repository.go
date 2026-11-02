package postgres

import (
	"database/sql"
	"encoding/json"
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

func defaultModules() []string {
	return []string{
		"agenda", "appointments", "customers", "services", "professionals",
		"finance", "commissions", "inventory", "pos_cashregister", "notifications", "reviews",
	}
}

func (r *TenantRepository) Create(input tenant.Tenant) (tenant.Tenant, error) {
	item := input
	if strings.TrimSpace(item.ID) == "" {
		return tenant.Tenant{}, fmt.Errorf("tenant id is required")
	}
	item.Slug = strings.ToLower(strings.TrimSpace(item.Slug))
	if item.Segment == "" {
		item.Segment = "general"
	}
	if item.Plan == "" {
		item.Plan = "trial"
	}
	if item.SubscriptionStatus == "" {
		item.SubscriptionStatus = "trialing"
	}
	if len(item.ActiveModules) == 0 {
		item.ActiveModules = defaultModules()
	}
	if item.MaxProfessionals <= 0 {
		item.MaxProfessionals = 5
	}
	modulesJSON, _ := json.Marshal(item.ActiveModules)

	if _, err := r.db.Exec(`
		INSERT INTO catalog.tenants (id, name, slug, segment, active, plan, subscription_status, active_modules, max_professionals)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, item.ID, item.Name, item.Slug, item.Segment, item.Active, item.Plan, item.SubscriptionStatus, string(modulesJSON), item.MaxProfessionals); err != nil {
		if isTenantIdentityConflict(err) {
			return tenant.Tenant{}, fmt.Errorf("slug already used")
		}
		return tenant.Tenant{}, fmt.Errorf("insert tenant: %w", err)
	}
	return r.GetByID(item.ID)
}

func (r *TenantRepository) GetByID(id string) (tenant.Tenant, error) {
	item, err := scanTenantRow(r.db.QueryRow(`
		SELECT id, name, slug, COALESCE(segment, 'general'), active,
		       COALESCE(plan, 'trial'), COALESCE(subscription_status, 'trialing'),
		       trial_ends_at, subscription_expires_at,
		       COALESCE(active_modules, '[]'::jsonb), COALESCE(max_professionals, 5),
		       created_at, updated_at, deleted_at
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
		SELECT id, name, slug, COALESCE(segment, 'general'), active,
		       COALESCE(plan, 'trial'), COALESCE(subscription_status, 'trialing'),
		       trial_ends_at, subscription_expires_at,
		       COALESCE(active_modules, '[]'::jsonb), COALESCE(max_professionals, 5),
		       created_at, updated_at, deleted_at
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
		SELECT id, name, slug, COALESCE(segment, 'general'), active,
		       COALESCE(plan, 'trial'), COALESCE(subscription_status, 'trialing'),
		       trial_ends_at, subscription_expires_at,
		       COALESCE(active_modules, '[]'::jsonb), COALESCE(max_professionals, 5),
		       created_at, updated_at, deleted_at
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
		SET name = $2, slug = $3, segment = COALESCE(NULLIF($4, ''), segment), updated_at = NOW()
		WHERE id = $1
	`, input.ID, input.Name, strings.ToLower(strings.TrimSpace(input.Slug)), input.Segment)
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

func (r *TenantRepository) UpdateSubscription(id string, input tenant.UpdateSubscriptionInput) (tenant.Tenant, error) {
	modules := input.ActiveModules
	if len(modules) == 0 {
		modules = defaultModules()
	}
	modulesJSON, _ := json.Marshal(modules)

	var trialEnds sql.NullTime
	if input.TrialEndsAt != "" {
		if t, err := time.Parse(time.RFC3339, input.TrialEndsAt); err == nil {
			trialEnds = sql.NullTime{Time: t, Valid: true}
		}
	}
	var subExpires sql.NullTime
	if input.SubscriptionExpiresAt != "" {
		if t, err := time.Parse(time.RFC3339, input.SubscriptionExpiresAt); err == nil {
			subExpires = sql.NullTime{Time: t, Valid: true}
		}
	}

	result, err := r.db.Exec(`
		UPDATE catalog.tenants
		SET plan = $2, subscription_status = $3, trial_ends_at = $4, subscription_expires_at = $5,
		    active_modules = $6, max_professionals = $7, updated_at = NOW()
		WHERE id = $1
	`, id, input.Plan, input.SubscriptionStatus, trialEnds, subExpires, string(modulesJSON), input.MaxProfessionals)
	if err != nil {
		return tenant.Tenant{}, fmt.Errorf("update tenant subscription: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
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
	var trialEndsAt, subExpiresAt sql.NullTime
	var activeModulesRaw string

	if err := scanner.Scan(
		&item.ID, &item.Name, &item.Slug, &item.Segment, &item.Active,
		&item.Plan, &item.SubscriptionStatus,
		&trialEndsAt, &subExpiresAt,
		&activeModulesRaw, &item.MaxProfessionals,
		&createdAt, &updatedAt, &deletedAt,
	); err != nil {
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
	if trialEndsAt.Valid {
		item.TrialEndsAt = trialEndsAt.Time.Format(time.RFC3339)
	}
	if subExpiresAt.Valid {
		item.SubscriptionExpiresAt = subExpiresAt.Time.Format(time.RFC3339)
	}
	if activeModulesRaw != "" {
		_ = json.Unmarshal([]byte(activeModulesRaw), &item.ActiveModules)
	}
	if len(item.ActiveModules) == 0 {
		item.ActiveModules = defaultModules()
	}
	return item, nil
}

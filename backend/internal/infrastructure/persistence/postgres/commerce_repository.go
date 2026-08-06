package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/commerce"
)

type PackageRepository struct {
	db *sql.DB
}

func NewPackageRepository(db *sql.DB) *PackageRepository {
	return &PackageRepository{db: db}
}

func (r *PackageRepository) Create(input commerce.Package) (commerce.Package, error) {
	item := input
	item.ID = uuid.NewString()
	item.CreatedAt = time.Now().Format(time.RFC3339)
	item.UpdatedAt = item.CreatedAt
	if _, err := r.db.Exec(`
		INSERT INTO catalog.packages (id, tenant_id, name, description, price_cents, sessions_count, duration_days, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, item.ID, item.TenantID, item.Name, item.Description, item.PriceCents, item.SessionsCount, item.DurationDays, item.Active); err != nil {
		return commerce.Package{}, fmt.Errorf("insert package: %w", err)
	}
	return r.GetByID(item.TenantID, item.ID)
}

func (r *PackageRepository) Update(input commerce.Package) (commerce.Package, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.packages
		SET name = COALESCE(NULLIF($3, ''), name),
		    description = $4,
		    price_cents = $5,
		    sessions_count = $6,
		    duration_days = $7,
		    active = $8,
		    updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`, input.ID, input.TenantID, strings.TrimSpace(input.Name), strings.TrimSpace(input.Description), input.PriceCents, input.SessionsCount, input.DurationDays, input.Active)
	if err != nil {
		return commerce.Package{}, fmt.Errorf("update package: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return commerce.Package{}, fmt.Errorf("check package update: %w", err)
	}
	if affected == 0 {
		return commerce.Package{}, fmt.Errorf("package not found")
	}
	return r.GetByID(input.TenantID, input.ID)
}

func (r *PackageRepository) GetByID(tenantID string, id string) (commerce.Package, error) {
	var item commerce.Package
	var createdAt, updatedAt time.Time
	err := r.db.QueryRow(`
		SELECT id, tenant_id, name, description, price_cents, sessions_count, duration_days, active, created_at, updated_at
		FROM catalog.packages
		WHERE id = $1 AND tenant_id = $2
	`, strings.TrimSpace(id), strings.TrimSpace(tenantID)).Scan(&item.ID, &item.TenantID, &item.Name, &item.Description, &item.PriceCents, &item.SessionsCount, &item.DurationDays, &item.Active, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return commerce.Package{}, fmt.Errorf("package not found")
		}
		return commerce.Package{}, fmt.Errorf("get package: %w", err)
	}
	item.CreatedAt = createdAt.Format(time.RFC3339)
	item.UpdatedAt = updatedAt.Format(time.RFC3339)
	return item, nil
}

func (r *PackageRepository) ListByTenant(tenantID string) ([]commerce.Package, error) {
	rows, err := r.db.Query(`
		SELECT id, tenant_id, name, description, price_cents, sessions_count, duration_days, active, created_at, updated_at
		FROM catalog.packages
		WHERE tenant_id = $1
		ORDER BY name ASC
	`, strings.TrimSpace(tenantID))
	if err != nil {
		return nil, fmt.Errorf("query packages: %w", err)
	}
	defer rows.Close()
	items := []commerce.Package{}
	for rows.Next() {
		var item commerce.Package
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&item.ID, &item.TenantID, &item.Name, &item.Description, &item.PriceCents, &item.SessionsCount, &item.DurationDays, &item.Active, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan package: %w", err)
		}
		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PackageRepository) SetActive(tenantID string, id string, active bool) (commerce.Package, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.packages
		SET active = $3, updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`, id, tenantID, active)
	if err != nil {
		return commerce.Package{}, fmt.Errorf("set package active: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return commerce.Package{}, fmt.Errorf("check package active: %w", err)
	}
	if affected == 0 {
		return commerce.Package{}, fmt.Errorf("package not found")
	}
	return r.GetByID(tenantID, id)
}

func (r *PackageRepository) Delete(tenantID string, id string) error {
	result, err := r.db.Exec(`DELETE FROM catalog.packages WHERE id = $1 AND tenant_id = $2`, id, tenantID)
	if err != nil {
		return fmt.Errorf("delete package: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check package delete: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("package not found")
	}
	return nil
}

type PackageCreditRepository struct {
	db *sql.DB
}

func NewPackageCreditRepository(db *sql.DB) *PackageCreditRepository {
	return &PackageCreditRepository{db: db}
}

func (r *PackageCreditRepository) Create(input commerce.PackageCredit) (commerce.PackageCredit, error) {
	item := input
	item.ID = uuid.NewString()
	item.CreatedAt = time.Now().Format(time.RFC3339)
	_, err := r.db.Exec(`
		INSERT INTO catalog.package_credits (id, tenant_id, package_id, customer_id, sessions_total, sessions_used, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, item.ID, item.TenantID, item.PackageID, item.CustomerID, item.SessionsTotal, item.SessionsUsed, nullIfEmpty(item.ExpiresAt))
	if err != nil {
		return commerce.PackageCredit{}, fmt.Errorf("insert package credit: %w", err)
	}
	return r.GetByID(item.TenantID, item.ID)
}

func (r *PackageCreditRepository) GetByID(tenantID string, id string) (commerce.PackageCredit, error) {
	var item commerce.PackageCredit
	var createdAt, expiresAt sql.NullTime
	err := r.db.QueryRow(`
		SELECT pc.id, pc.tenant_id, pc.package_id, pc.customer_id, pc.sessions_total, pc.sessions_used, pc.expires_at, pc.created_at,
		       COALESCE(p.name, ''), COALESCE(c.name, '')
		FROM catalog.package_credits pc
		LEFT JOIN catalog.packages p ON p.id = pc.package_id
		LEFT JOIN catalog.customers c ON c.id = pc.customer_id
		WHERE pc.id = $1 AND pc.tenant_id = $2
	`, strings.TrimSpace(id), strings.TrimSpace(tenantID)).Scan(&item.ID, &item.TenantID, &item.PackageID, &item.CustomerID, &item.SessionsTotal, &item.SessionsUsed, &expiresAt, &createdAt, &item.PackageName, &item.CustomerName)
	if err != nil {
		if err == sql.ErrNoRows {
			return commerce.PackageCredit{}, fmt.Errorf("package credit not found")
		}
		return commerce.PackageCredit{}, fmt.Errorf("get package credit: %w", err)
	}
	item.CreatedAt = createdAt.Time.Format(time.RFC3339)
	if expiresAt.Valid {
		item.ExpiresAt = expiresAt.Time.UTC().Format(time.RFC3339)
	}
	return item, nil
}

func (r *PackageCreditRepository) ListByCustomer(customerID string) ([]commerce.PackageCredit, error) {
	rows, err := r.db.Query(`
		SELECT pc.id, pc.tenant_id, pc.package_id, pc.customer_id, pc.sessions_total, pc.sessions_used, pc.expires_at, pc.created_at,
		       COALESCE(p.name, ''), COALESCE(c.name, '')
		FROM catalog.package_credits pc
		LEFT JOIN catalog.packages p ON p.id = pc.package_id
		LEFT JOIN catalog.customers c ON c.id = pc.customer_id
		WHERE pc.customer_id = $1
		ORDER BY pc.created_at DESC
	`, strings.TrimSpace(customerID))
	if err != nil {
		return nil, fmt.Errorf("query package credits: %w", err)
	}
	defer rows.Close()
	return scanPackageCredits(rows)
}

func (r *PackageCreditRepository) ListByTenant(tenantID string) ([]commerce.PackageCredit, error) {
	rows, err := r.db.Query(`
		SELECT pc.id, pc.tenant_id, pc.package_id, pc.customer_id, pc.sessions_total, pc.sessions_used, pc.expires_at, pc.created_at,
		       COALESCE(p.name, ''), COALESCE(c.name, '')
		FROM catalog.package_credits pc
		LEFT JOIN catalog.packages p ON p.id = pc.package_id
		LEFT JOIN catalog.customers c ON c.id = pc.customer_id
		WHERE pc.tenant_id = $1
		ORDER BY pc.created_at DESC
	`, strings.TrimSpace(tenantID))
	if err != nil {
		return nil, fmt.Errorf("query package credits: %w", err)
	}
	defer rows.Close()
	return scanPackageCredits(rows)
}

func (r *PackageCreditRepository) IncrementUsed(tenantID string, id string) (commerce.PackageCredit, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.package_credits
		SET sessions_used = sessions_used + 1
		WHERE id = $1 AND tenant_id = $2 AND sessions_used < sessions_total
	`, id, tenantID)
	if err != nil {
		return commerce.PackageCredit{}, fmt.Errorf("increment package credit: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return commerce.PackageCredit{}, fmt.Errorf("check package credit increment: %w", err)
	}
	if affected == 0 {
		return commerce.PackageCredit{}, fmt.Errorf("package credit exhausted")
	}
	return r.GetByID(tenantID, id)
}

type GiftCardRepository struct {
	db *sql.DB
}

func NewGiftCardRepository(db *sql.DB) *GiftCardRepository {
	return &GiftCardRepository{db: db}
}

func (r *GiftCardRepository) Create(input commerce.GiftCard) (commerce.GiftCard, error) {
	item := input
	item.ID = uuid.NewString()
	item.CreatedAt = time.Now().Format(time.RFC3339)
	if item.Status == "" {
		item.Status = "active"
	}
	if _, err := r.db.Exec(`
		INSERT INTO catalog.gift_cards (id, tenant_id, code, value_cents, balance_cents, customer_id, status, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, item.ID, item.TenantID, item.Code, item.ValueCents, item.BalanceCents, item.CustomerID, item.Status, nullIfEmpty(item.ExpiresAt)); err != nil {
		return commerce.GiftCard{}, fmt.Errorf("insert gift card: %w", err)
	}
	return r.GetByID(item.TenantID, item.ID)
}

func (r *GiftCardRepository) GetByID(tenantID string, id string) (commerce.GiftCard, error) {
	var item commerce.GiftCard
	var createdAt, expiresAt sql.NullTime
	err := r.db.QueryRow(`
		SELECT g.id, g.tenant_id, g.code, g.value_cents, g.balance_cents, g.customer_id, g.status, g.expires_at, g.created_at,
		       COALESCE(c.name, '')
		FROM catalog.gift_cards g
		LEFT JOIN catalog.customers c ON c.id = g.customer_id
		WHERE g.id = $1 AND g.tenant_id = $2
	`, strings.TrimSpace(id), strings.TrimSpace(tenantID)).Scan(&item.ID, &item.TenantID, &item.Code, &item.ValueCents, &item.BalanceCents, &item.CustomerID, &item.Status, &expiresAt, &createdAt, &item.CustomerName)
	if err != nil {
		if err == sql.ErrNoRows {
			return commerce.GiftCard{}, fmt.Errorf("gift card not found")
		}
		return commerce.GiftCard{}, fmt.Errorf("get gift card: %w", err)
	}
	item.CreatedAt = createdAt.Time.Format(time.RFC3339)
	if expiresAt.Valid {
		item.ExpiresAt = expiresAt.Time.UTC().Format(time.RFC3339)
	}
	return item, nil
}

func (r *GiftCardRepository) GetByCode(tenantID string, code string) (commerce.GiftCard, error) {
	var item commerce.GiftCard
	var createdAt, expiresAt sql.NullTime
	err := r.db.QueryRow(`
		SELECT g.id, g.tenant_id, g.code, g.value_cents, g.balance_cents, g.customer_id, g.status, g.expires_at, g.created_at,
		       COALESCE(c.name, '')
		FROM catalog.gift_cards g
		LEFT JOIN catalog.customers c ON c.id = g.customer_id
		WHERE g.code = $1 AND g.tenant_id = $2
	`, strings.TrimSpace(code), strings.TrimSpace(tenantID)).Scan(&item.ID, &item.TenantID, &item.Code, &item.ValueCents, &item.BalanceCents, &item.CustomerID, &item.Status, &expiresAt, &createdAt, &item.CustomerName)
	if err != nil {
		if err == sql.ErrNoRows {
			return commerce.GiftCard{}, fmt.Errorf("gift card not found")
		}
		return commerce.GiftCard{}, fmt.Errorf("get gift card by code: %w", err)
	}
	item.CreatedAt = createdAt.Time.Format(time.RFC3339)
	if expiresAt.Valid {
		item.ExpiresAt = expiresAt.Time.UTC().Format(time.RFC3339)
	}
	return item, nil
}

func (r *GiftCardRepository) ListByTenant(tenantID string) ([]commerce.GiftCard, error) {
	rows, err := r.db.Query(`
		SELECT g.id, g.tenant_id, g.code, g.value_cents, g.balance_cents, g.customer_id, g.status, g.expires_at, g.created_at,
		       COALESCE(c.name, '')
		FROM catalog.gift_cards g
		LEFT JOIN catalog.customers c ON c.id = g.customer_id
		WHERE g.tenant_id = $1
		ORDER BY g.created_at DESC
	`, strings.TrimSpace(tenantID))
	if err != nil {
		return nil, fmt.Errorf("query gift cards: %w", err)
	}
	defer rows.Close()
	items := []commerce.GiftCard{}
	for rows.Next() {
		var item commerce.GiftCard
		var createdAt, expiresAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.TenantID, &item.Code, &item.ValueCents, &item.BalanceCents, &item.CustomerID, &item.Status, &expiresAt, &createdAt, &item.CustomerName); err != nil {
			return nil, fmt.Errorf("scan gift card: %w", err)
		}
		item.CreatedAt = createdAt.Time.Format(time.RFC3339)
		if expiresAt.Valid {
			item.ExpiresAt = expiresAt.Time.UTC().Format(time.RFC3339)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *GiftCardRepository) Update(input commerce.GiftCard) (commerce.GiftCard, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.gift_cards
		SET balance_cents = $3,
		    customer_id = $4,
		    status = $5
		WHERE id = $1 AND tenant_id = $2
	`, input.ID, input.TenantID, input.BalanceCents, input.CustomerID, input.Status)
	if err != nil {
		return commerce.GiftCard{}, fmt.Errorf("update gift card: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return commerce.GiftCard{}, fmt.Errorf("check gift card update: %w", err)
	}
	if affected == 0 {
		return commerce.GiftCard{}, fmt.Errorf("gift card not found")
	}
	return r.GetByID(input.TenantID, input.ID)
}

func (r *GiftCardRepository) Delete(tenantID string, id string) error {
	result, err := r.db.Exec(`DELETE FROM catalog.gift_cards WHERE id = $1 AND tenant_id = $2`, id, tenantID)
	if err != nil {
		return fmt.Errorf("delete gift card: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check gift card delete: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("gift card not found")
	}
	return nil
}

type MembershipRepository struct {
	db *sql.DB
}

func NewMembershipRepository(db *sql.DB) *MembershipRepository {
	return &MembershipRepository{db: db}
}

func (r *MembershipRepository) Create(input commerce.Membership) (commerce.Membership, error) {
	item := input
	item.ID = uuid.NewString()
	item.CreatedAt = time.Now().Format(time.RFC3339)
	item.UpdatedAt = item.CreatedAt
	if item.BillingPeriodDays <= 0 {
		item.BillingPeriodDays = 30
	}
	if _, err := r.db.Exec(`
		INSERT INTO catalog.memberships (id, tenant_id, name, description, price_cents, billing_period_days, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, item.ID, item.TenantID, item.Name, item.Description, item.PriceCents, item.BillingPeriodDays, item.Active); err != nil {
		return commerce.Membership{}, fmt.Errorf("insert membership: %w", err)
	}
	return r.GetByID(item.TenantID, item.ID)
}

func (r *MembershipRepository) Update(input commerce.Membership) (commerce.Membership, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.memberships
		SET name = COALESCE(NULLIF($3, ''), name),
		    description = $4,
		    price_cents = $5,
		    billing_period_days = $6,
		    active = $7,
		    updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`, input.ID, input.TenantID, strings.TrimSpace(input.Name), strings.TrimSpace(input.Description), input.PriceCents, input.BillingPeriodDays, input.Active)
	if err != nil {
		return commerce.Membership{}, fmt.Errorf("update membership: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return commerce.Membership{}, fmt.Errorf("check membership update: %w", err)
	}
	if affected == 0 {
		return commerce.Membership{}, fmt.Errorf("membership not found")
	}
	return r.GetByID(input.TenantID, input.ID)
}

func (r *MembershipRepository) GetByID(tenantID string, id string) (commerce.Membership, error) {
	var item commerce.Membership
	var createdAt, updatedAt time.Time
	err := r.db.QueryRow(`
		SELECT id, tenant_id, name, description, price_cents, billing_period_days, active, created_at, updated_at
		FROM catalog.memberships
		WHERE id = $1 AND tenant_id = $2
	`, strings.TrimSpace(id), strings.TrimSpace(tenantID)).Scan(&item.ID, &item.TenantID, &item.Name, &item.Description, &item.PriceCents, &item.BillingPeriodDays, &item.Active, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return commerce.Membership{}, fmt.Errorf("membership not found")
		}
		return commerce.Membership{}, fmt.Errorf("get membership: %w", err)
	}
	item.CreatedAt = createdAt.Format(time.RFC3339)
	item.UpdatedAt = updatedAt.Format(time.RFC3339)
	return item, nil
}

func (r *MembershipRepository) ListByTenant(tenantID string) ([]commerce.Membership, error) {
	rows, err := r.db.Query(`
		SELECT id, tenant_id, name, description, price_cents, billing_period_days, active, created_at, updated_at
		FROM catalog.memberships
		WHERE tenant_id = $1
		ORDER BY name ASC
	`, strings.TrimSpace(tenantID))
	if err != nil {
		return nil, fmt.Errorf("query memberships: %w", err)
	}
	defer rows.Close()
	items := []commerce.Membership{}
	for rows.Next() {
		var item commerce.Membership
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&item.ID, &item.TenantID, &item.Name, &item.Description, &item.PriceCents, &item.BillingPeriodDays, &item.Active, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan membership: %w", err)
		}
		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *MembershipRepository) SetActive(tenantID string, id string, active bool) (commerce.Membership, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.memberships
		SET active = $3, updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`, id, tenantID, active)
	if err != nil {
		return commerce.Membership{}, fmt.Errorf("set membership active: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return commerce.Membership{}, fmt.Errorf("check membership active: %w", err)
	}
	if affected == 0 {
		return commerce.Membership{}, fmt.Errorf("membership not found")
	}
	return r.GetByID(tenantID, id)
}

func (r *MembershipRepository) Delete(tenantID string, id string) error {
	result, err := r.db.Exec(`DELETE FROM catalog.memberships WHERE id = $1 AND tenant_id = $2`, id, tenantID)
	if err != nil {
		return fmt.Errorf("delete membership: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check membership delete: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("membership not found")
	}
	return nil
}

type MemberPlanRepository struct {
	db *sql.DB
}

func NewMemberPlanRepository(db *sql.DB) *MemberPlanRepository {
	return &MemberPlanRepository{db: db}
}

func (r *MemberPlanRepository) Create(input commerce.MemberPlan) (commerce.MemberPlan, error) {
	item := input
	item.ID = uuid.NewString()
	item.CreatedAt = time.Now().Format(time.RFC3339)
	if item.Status == "" {
		item.Status = "active"
	}
	if item.StartedAt == "" {
		item.StartedAt = item.CreatedAt
	}
	_, err := r.db.Exec(`
		INSERT INTO catalog.member_plans (id, tenant_id, membership_id, customer_id, status, started_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, item.ID, item.TenantID, item.MembershipID, item.CustomerID, item.Status, item.StartedAt, nullIfEmpty(item.ExpiresAt))
	if err != nil {
		return commerce.MemberPlan{}, fmt.Errorf("insert member plan: %w", err)
	}
	return r.GetByID(item.TenantID, item.ID)
}

func (r *MemberPlanRepository) GetByID(tenantID string, id string) (commerce.MemberPlan, error) {
	var item commerce.MemberPlan
	var createdAt, startedAt, expiresAt sql.NullTime
	err := r.db.QueryRow(`
		SELECT mp.id, mp.tenant_id, mp.membership_id, mp.customer_id, mp.status, mp.started_at, mp.expires_at, mp.created_at,
		       COALESCE(m.name, ''), COALESCE(c.name, '')
		FROM catalog.member_plans mp
		LEFT JOIN catalog.memberships m ON m.id = mp.membership_id
		LEFT JOIN catalog.customers c ON c.id = mp.customer_id
		WHERE mp.id = $1 AND mp.tenant_id = $2
	`, strings.TrimSpace(id), strings.TrimSpace(tenantID)).Scan(&item.ID, &item.TenantID, &item.MembershipID, &item.CustomerID, &item.Status, &startedAt, &expiresAt, &createdAt, &item.Membership, &item.CustomerName)
	if err != nil {
		if err == sql.ErrNoRows {
			return commerce.MemberPlan{}, fmt.Errorf("member plan not found")
		}
		return commerce.MemberPlan{}, fmt.Errorf("get member plan: %w", err)
	}
	item.CreatedAt = createdAt.Time.Format(time.RFC3339)
	if startedAt.Valid {
		item.StartedAt = startedAt.Time.UTC().Format(time.RFC3339)
	}
	if expiresAt.Valid {
		item.ExpiresAt = expiresAt.Time.UTC().Format(time.RFC3339)
	}
	return item, nil
}

func (r *MemberPlanRepository) ListByCustomer(customerID string) ([]commerce.MemberPlan, error) {
	rows, err := r.db.Query(`
		SELECT mp.id, mp.tenant_id, mp.membership_id, mp.customer_id, mp.status, mp.started_at, mp.expires_at, mp.created_at,
		       COALESCE(m.name, ''), COALESCE(c.name, '')
		FROM catalog.member_plans mp
		LEFT JOIN catalog.memberships m ON m.id = mp.membership_id
		LEFT JOIN catalog.customers c ON c.id = mp.customer_id
		WHERE mp.customer_id = $1
		ORDER BY mp.created_at DESC
	`, strings.TrimSpace(customerID))
	if err != nil {
		return nil, fmt.Errorf("query member plans: %w", err)
	}
	defer rows.Close()
	return scanMemberPlans(rows)
}

func (r *MemberPlanRepository) ListByTenant(tenantID string) ([]commerce.MemberPlan, error) {
	rows, err := r.db.Query(`
		SELECT mp.id, mp.tenant_id, mp.membership_id, mp.customer_id, mp.status, mp.started_at, mp.expires_at, mp.created_at,
		       COALESCE(m.name, ''), COALESCE(c.name, '')
		FROM catalog.member_plans mp
		LEFT JOIN catalog.memberships m ON m.id = mp.membership_id
		LEFT JOIN catalog.customers c ON c.id = mp.customer_id
		WHERE mp.tenant_id = $1
		ORDER BY mp.created_at DESC
	`, strings.TrimSpace(tenantID))
	if err != nil {
		return nil, fmt.Errorf("query member plans: %w", err)
	}
	defer rows.Close()
	return scanMemberPlans(rows)
}

func (r *MemberPlanRepository) UpdateStatus(tenantID string, id string, status string) (commerce.MemberPlan, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.member_plans
		SET status = $3
		WHERE id = $1 AND tenant_id = $2
	`, id, tenantID, status)
	if err != nil {
		return commerce.MemberPlan{}, fmt.Errorf("update member plan status: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return commerce.MemberPlan{}, fmt.Errorf("check member plan update: %w", err)
	}
	if affected == 0 {
		return commerce.MemberPlan{}, fmt.Errorf("member plan not found")
	}
	return r.GetByID(tenantID, id)
}

func scanPackageCredits(rows *sql.Rows) ([]commerce.PackageCredit, error) {
	items := []commerce.PackageCredit{}
	for rows.Next() {
		var item commerce.PackageCredit
		var createdAt, expiresAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.TenantID, &item.PackageID, &item.CustomerID, &item.SessionsTotal, &item.SessionsUsed, &expiresAt, &createdAt, &item.PackageName, &item.CustomerName); err != nil {
			return nil, fmt.Errorf("scan package credit: %w", err)
		}
		item.CreatedAt = createdAt.Time.Format(time.RFC3339)
		if expiresAt.Valid {
			item.ExpiresAt = expiresAt.Time.UTC().Format(time.RFC3339)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func scanMemberPlans(rows *sql.Rows) ([]commerce.MemberPlan, error) {
	items := []commerce.MemberPlan{}
	for rows.Next() {
		var item commerce.MemberPlan
		var createdAt, startedAt, expiresAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.TenantID, &item.MembershipID, &item.CustomerID, &item.Status, &startedAt, &expiresAt, &createdAt, &item.Membership, &item.CustomerName); err != nil {
			return nil, fmt.Errorf("scan member plan: %w", err)
		}
		item.CreatedAt = createdAt.Time.Format(time.RFC3339)
		if startedAt.Valid {
			item.StartedAt = startedAt.Time.UTC().Format(time.RFC3339)
		}
		if expiresAt.Valid {
			item.ExpiresAt = expiresAt.Time.UTC().Format(time.RFC3339)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

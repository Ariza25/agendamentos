package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"erpagendamentos/backend/internal/domain/commission"
	"github.com/google/uuid"
)

type CommissionRepository struct {
	db *sql.DB
}

func NewCommissionRepository(db *sql.DB) *CommissionRepository {
	return &CommissionRepository{db: db}
}

func (r *CommissionRepository) CreateRule(rule commission.CommissionRule) (commission.CommissionRule, error) {
	if rule.ID == "" {
		rule.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	err := r.db.QueryRow(`
		INSERT INTO catalog.commission_rules (id, tenant_id, professional_id, service_id, percentage, fixed_amount_cents, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at
	`, rule.ID, rule.TenantID, rule.ProfessionalID, rule.ServiceID, rule.Percentage, rule.FixedAmountCents, now, now).Scan(&rule.CreatedAt, &rule.UpdatedAt)
	if err != nil {
		return commission.CommissionRule{}, fmt.Errorf("create commission rule: %w", err)
	}
	return rule, nil
}

func (r *CommissionRepository) UpdateRule(rule commission.CommissionRule) (commission.CommissionRule, error) {
	now := time.Now().UTC()
	result, err := r.db.Exec(`
		UPDATE catalog.commission_rules
		SET professional_id = $1, service_id = $2, percentage = $3, fixed_amount_cents = $4, updated_at = $5
		WHERE tenant_id = $6 AND id = $7
	`, rule.ProfessionalID, rule.ServiceID, rule.Percentage, rule.FixedAmountCents, now, rule.TenantID, rule.ID)
	if err != nil {
		return commission.CommissionRule{}, fmt.Errorf("update commission rule: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return commission.CommissionRule{}, errors.New("commission rule not found")
	}
	rule.UpdatedAt = now.Format(time.RFC3339)
	return rule, nil
}

func (r *CommissionRepository) DeleteRule(tenantID string, id string) error {
	_, err := r.db.Exec("DELETE FROM catalog.commission_rules WHERE tenant_id = $1 AND id = $2", tenantID, id)
	return err
}

func (r *CommissionRepository) ListRules(tenantID string) ([]commission.CommissionRule, error) {
	rows, err := r.db.Query(`
		SELECT id, tenant_id, professional_id, service_id, percentage, fixed_amount_cents, created_at, updated_at
		FROM catalog.commission_rules
		WHERE tenant_id = $1
		ORDER BY created_at DESC
	`, tenantID)
	if err != nil {
		return nil, fmt.Errorf("list commission rules: %w", err)
	}
	defer rows.Close()

	var rules []commission.CommissionRule
	for rows.Next() {
		var c commission.CommissionRule
		var created, updated time.Time
		if err := rows.Scan(&c.ID, &c.TenantID, &c.ProfessionalID, &c.ServiceID, &c.Percentage, &c.FixedAmountCents, &created, &updated); err != nil {
			return nil, err
		}
		c.CreatedAt = created.Format(time.RFC3339)
		c.UpdatedAt = updated.Format(time.RFC3339)
		rules = append(rules, c)
	}
	return rules, nil
}

func (r *CommissionRepository) FindMatchingRule(tenantID string, professionalID string, serviceID string) (commission.CommissionRule, error) {
	var c commission.CommissionRule
	var created, updated time.Time

	// 1. Try specific prof + service
	err := r.db.QueryRow(`
		SELECT id, tenant_id, professional_id, service_id, percentage, fixed_amount_cents, created_at, updated_at
		FROM catalog.commission_rules
		WHERE tenant_id = $1 AND professional_id = $2 AND service_id = $3
		LIMIT 1
	`, tenantID, professionalID, serviceID).Scan(&c.ID, &c.TenantID, &c.ProfessionalID, &c.ServiceID, &c.Percentage, &c.FixedAmountCents, &created, &updated)
	if err == nil {
		c.CreatedAt = created.Format(time.RFC3339)
		c.UpdatedAt = updated.Format(time.RFC3339)
		return c, nil
	}

	// 2. Try professional only
	err = r.db.QueryRow(`
		SELECT id, tenant_id, professional_id, service_id, percentage, fixed_amount_cents, created_at, updated_at
		FROM catalog.commission_rules
		WHERE tenant_id = $1 AND professional_id = $2 AND service_id = ''
		LIMIT 1
	`, tenantID, professionalID).Scan(&c.ID, &c.TenantID, &c.ProfessionalID, &c.ServiceID, &c.Percentage, &c.FixedAmountCents, &created, &updated)
	if err == nil {
		c.CreatedAt = created.Format(time.RFC3339)
		c.UpdatedAt = updated.Format(time.RFC3339)
		return c, nil
	}

	// 3. Try service only
	err = r.db.QueryRow(`
		SELECT id, tenant_id, professional_id, service_id, percentage, fixed_amount_cents, created_at, updated_at
		FROM catalog.commission_rules
		WHERE tenant_id = $1 AND professional_id = '' AND service_id = $2
		LIMIT 1
	`, tenantID, serviceID).Scan(&c.ID, &c.TenantID, &c.ProfessionalID, &c.ServiceID, &c.Percentage, &c.FixedAmountCents, &created, &updated)
	if err == nil {
		c.CreatedAt = created.Format(time.RFC3339)
		c.UpdatedAt = updated.Format(time.RFC3339)
		return c, nil
	}

	// 4. Default global rule
	err = r.db.QueryRow(`
		SELECT id, tenant_id, professional_id, service_id, percentage, fixed_amount_cents, created_at, updated_at
		FROM catalog.commission_rules
		WHERE tenant_id = $1 AND professional_id = '' AND service_id = ''
		LIMIT 1
	`, tenantID).Scan(&c.ID, &c.TenantID, &c.ProfessionalID, &c.ServiceID, &c.Percentage, &c.FixedAmountCents, &created, &updated)
	if err == nil {
		c.CreatedAt = created.Format(time.RFC3339)
		c.UpdatedAt = updated.Format(time.RFC3339)
		return c, nil
	}

	return commission.CommissionRule{}, errors.New("no matching commission rule found")
}

func (r *CommissionRepository) CreateRecord(record commission.CommissionRecord) (commission.CommissionRecord, error) {
	if record.ID == "" {
		record.ID = uuid.NewString()
	}
	if record.Status == "" {
		record.Status = "pending"
	}
	now := time.Now().UTC()
	err := r.db.QueryRow(`
		INSERT INTO catalog.commission_records (id, tenant_id, appointment_id, professional_id, service_id, base_amount_cents, commission_amount_cents, status, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at
	`, record.ID, record.TenantID, record.AppointmentID, record.ProfessionalID, record.ServiceID, record.BaseAmountCents, record.CommissionAmountCents, record.Status, record.Notes, now).Scan(&record.CreatedAt)
	if err != nil {
		return commission.CommissionRecord{}, fmt.Errorf("create commission record: %w", err)
	}
	return record, nil
}

func (r *CommissionRepository) ListRecords(tenantID string, professionalID string, status string) ([]commission.CommissionRecord, error) {
	query := `
		SELECT cr.id, cr.tenant_id, cr.appointment_id, cr.professional_id, p.name, cr.service_id, COALESCE(s.name, ''), cr.base_amount_cents, cr.commission_amount_cents, cr.status, cr.settled_at, cr.notes, cr.created_at
		FROM catalog.commission_records cr
		JOIN catalog.professionals p ON cr.professional_id = p.id
		LEFT JOIN catalog.services s ON cr.service_id = s.id
		WHERE cr.tenant_id = $1
	`
	args := []any{tenantID}
	if professionalID != "" {
		args = append(args, professionalID)
		query += fmt.Sprintf(" AND cr.professional_id = $%d", len(args))
	}
	if status != "" {
		args = append(args, status)
		query += fmt.Sprintf(" AND cr.status = $%d", len(args))
	}
	query += " ORDER BY cr.created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list commission records: %w", err)
	}
	defer rows.Close()

	var records []commission.CommissionRecord
	for rows.Next() {
		var rec commission.CommissionRecord
		var settled *time.Time
		var created time.Time
		if err := rows.Scan(&rec.ID, &rec.TenantID, &rec.AppointmentID, &rec.ProfessionalID, &rec.ProfessionalName, &rec.ServiceID, &rec.ServiceName, &rec.BaseAmountCents, &rec.CommissionAmountCents, &rec.Status, &settled, &rec.Notes, &created); err != nil {
			return nil, err
		}
		if settled != nil {
			rec.SettledAt = settled.Format(time.RFC3339)
		}
		rec.CreatedAt = created.Format(time.RFC3339)
		records = append(records, rec)
	}
	return records, nil
}

func (r *CommissionRepository) SettleRecords(tenantID string, professionalID string, recordIDs []string) error {
	now := time.Now().UTC()
	if len(recordIDs) > 0 {
		query := `
			UPDATE catalog.commission_records
			SET status = 'paid', settled_at = $1
			WHERE tenant_id = $2 AND status = 'pending' AND id = ANY($3)
		`
		_, err := r.db.Exec(query, now, tenantID, recordIDs)
		return err
	}

	query := `
		UPDATE catalog.commission_records
		SET status = 'paid', settled_at = $1
		WHERE tenant_id = $2 AND status = 'pending'
	`
	args := []any{now, tenantID}
	if professionalID != "" {
		query += " AND professional_id = $3"
		args = append(args, professionalID)
	}
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *CommissionRepository) GetSummaries(tenantID string) ([]commission.Summary, error) {
	query := `
		SELECT cr.professional_id, p.name,
			COALESCE(SUM(CASE WHEN cr.status = 'pending' THEN cr.commission_amount_cents ELSE 0 END), 0) AS pending_cents,
			COALESCE(SUM(CASE WHEN cr.status = 'paid' THEN cr.commission_amount_cents ELSE 0 END), 0) AS paid_cents,
			COALESCE(SUM(cr.commission_amount_cents), 0) AS total_cents
		FROM catalog.commission_records cr
		JOIN catalog.professionals p ON cr.professional_id = p.id
		WHERE cr.tenant_id = $1
		GROUP BY cr.professional_id, p.name
		ORDER BY p.name ASC
	`
	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("get commission summaries: %w", err)
	}
	defer rows.Close()

	var summaries []commission.Summary
	for rows.Next() {
		var s commission.Summary
		if err := rows.Scan(&s.ProfessionalID, &s.ProfessionalName, &s.PendingCents, &s.PaidCents, &s.TotalCents); err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}
	return summaries, nil
}


package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	auditdomain "erpagendamentos/backend/internal/domain/audit"
)

type AuditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Record(entry auditdomain.AuditLog) error {
	detailsJSON, err := json.Marshal(entry.Details)
	if err != nil {
		detailsJSON = []byte("{}")
	}

	query := `
		INSERT INTO catalog.audit_logs (
			id, tenant_id, user_id, user_email, action, resource, resource_id,
			details, ip_address, user_agent, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err = r.db.Exec(
		query,
		entry.ID,
		entry.TenantID,
		entry.UserID,
		entry.UserEmail,
		entry.Action,
		entry.Resource,
		entry.ResourceID,
		detailsJSON,
		entry.IPAddress,
		entry.UserAgent,
		entry.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("record audit log: %w", err)
	}
	return nil
}

func (r *AuditRepository) List(tenantID string, limit int, offset int) ([]auditdomain.AuditLog, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	tenantID = strings.TrimSpace(tenantID)
	var countQuery string
	var selectQuery string
	var args []any

	if tenantID != "" && tenantID != "default" {
		countQuery = `SELECT COUNT(*) FROM catalog.audit_logs WHERE tenant_id = $1`
		selectQuery = `
			SELECT id, tenant_id, user_id, user_email, action, resource, resource_id,
			       details, ip_address, user_agent, created_at
			FROM catalog.audit_logs
			WHERE tenant_id = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []any{tenantID, limit, offset}
	} else {
		countQuery = `SELECT COUNT(*) FROM catalog.audit_logs`
		selectQuery = `
			SELECT id, tenant_id, user_id, user_email, action, resource, resource_id,
			       details, ip_address, user_agent, created_at
			FROM catalog.audit_logs
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`
		args = []any{limit, offset}
	}

	var total int
	if tenantID != "" && tenantID != "default" {
		if err := r.db.QueryRow(countQuery, tenantID).Scan(&total); err != nil {
			return nil, 0, fmt.Errorf("count audit logs: %w", err)
		}
	} else {
		if err := r.db.QueryRow(countQuery).Scan(&total); err != nil {
			return nil, 0, fmt.Errorf("count audit logs: %w", err)
		}
	}

	rows, err := r.db.Query(selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query audit logs: %w", err)
	}
	defer rows.Close()

	var logs []auditdomain.AuditLog
	for rows.Next() {
		var l auditdomain.AuditLog
		var rawDetails []byte
		if err := rows.Scan(
			&l.ID,
			&l.TenantID,
			&l.UserID,
			&l.UserEmail,
			&l.Action,
			&l.Resource,
			&l.ResourceID,
			&rawDetails,
			&l.IPAddress,
			&l.UserAgent,
			&l.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan audit log: %w", err)
		}
		if len(rawDetails) > 0 {
			_ = json.Unmarshal(rawDetails, &l.Details)
		}
		logs = append(logs, l)
	}

	return logs, total, nil
}


package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/professional"
)

type ProfessionalRepository struct {
	db *sql.DB
}

func NewProfessionalRepository(db *sql.DB) *ProfessionalRepository {
	return &ProfessionalRepository{db: db}
}

func (r *ProfessionalRepository) Create(input professional.Professional) (professional.Professional, error) {
	item := input
	item.ID = uuid.NewString()
	item.Name = strings.TrimSpace(item.Name)
	item.Email = strings.ToLower(strings.TrimSpace(item.Email))
	if _, err := r.db.Exec(`
		INSERT INTO catalog.professionals (id, tenant_id, user_id, name, email, phone, color, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, item.ID, item.TenantID, item.UserID, item.Name, item.Email, item.Phone, item.Color, item.Active); err != nil {
		return professional.Professional{}, fmt.Errorf("insert professional: %w", err)
	}
	return r.GetByID(item.TenantID, item.ID)
}

func (r *ProfessionalRepository) Update(input professional.Professional) (professional.Professional, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.professionals
		SET name = COALESCE(NULLIF($3, ''), name),
		    email = COALESCE(NULLIF(LOWER($4), ''), email),
		    phone = $5,
		    color = COALESCE(NULLIF($6, ''), color),
		    user_id = COALESCE(NULLIF($7, ''), user_id),
		    active = $8,
		    updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`, input.ID, input.TenantID, strings.TrimSpace(input.Name), strings.TrimSpace(input.Email), strings.TrimSpace(input.Phone), strings.TrimSpace(input.Color), strings.TrimSpace(input.UserID), input.Active)
	if err != nil {
		return professional.Professional{}, fmt.Errorf("update professional: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return professional.Professional{}, fmt.Errorf("check professional update: %w", err)
	}
	if affected == 0 {
		return professional.Professional{}, fmt.Errorf("professional not found")
	}
	return r.GetByID(input.TenantID, input.ID)
}

func (r *ProfessionalRepository) GetByID(tenantID string, id string) (professional.Professional, error) {
	item, err := scanProfessionalRow(r.db.QueryRow(`
		SELECT id, tenant_id, user_id, name, email, phone, color, active, created_at, updated_at
		FROM catalog.professionals
		WHERE id = $1 AND tenant_id = $2
	`, strings.TrimSpace(id), strings.TrimSpace(tenantID)))
	if err != nil {
		return professional.Professional{}, err
	}
	return item, nil
}

func (r *ProfessionalRepository) ListByTenant(tenantID string) ([]professional.Professional, error) {
	rows, err := r.db.Query(`
		SELECT id, tenant_id, user_id, name, email, phone, color, active, created_at, updated_at
		FROM catalog.professionals
		WHERE tenant_id = $1
		ORDER BY name ASC
	`, strings.TrimSpace(tenantID))
	if err != nil {
		return nil, fmt.Errorf("query professionals: %w", err)
	}
	defer rows.Close()

	items := []professional.Professional{}
	for rows.Next() {
		item, err := scanProfessionalRows(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *ProfessionalRepository) SetActive(tenantID string, id string, active bool) (professional.Professional, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.professionals
		SET active = $3, updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`, id, tenantID, active)
	if err != nil {
		return professional.Professional{}, fmt.Errorf("set professional active: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return professional.Professional{}, fmt.Errorf("check professional active: %w", err)
	}
	if affected == 0 {
		return professional.Professional{}, fmt.Errorf("professional not found")
	}
	return r.GetByID(tenantID, id)
}

func (r *ProfessionalRepository) Delete(tenantID string, id string) error {
	result, err := r.db.Exec(`
		DELETE FROM catalog.professionals
		WHERE id = $1 AND tenant_id = $2
	`, id, tenantID)
	if err != nil {
		return fmt.Errorf("delete professional: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check professional delete: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("professional not found")
	}
	return nil
}

type professionalScanner interface {
	Scan(dest ...any) error
}

func scanProfessionalRow(scanner professionalScanner) (professional.Professional, error) {
	return scanProfessionalRows(scanner)
}

func scanProfessionalRows(scanner professionalScanner) (professional.Professional, error) {
	var item professional.Professional
	var createdAt, updatedAt time.Time
	if err := scanner.Scan(&item.ID, &item.TenantID, &item.UserID, &item.Name, &item.Email, &item.Phone, &item.Color, &item.Active, &createdAt, &updatedAt); err != nil {
		if err == sql.ErrNoRows {
			return professional.Professional{}, fmt.Errorf("professional not found")
		}
		return professional.Professional{}, fmt.Errorf("scan professional: %w", err)
	}
	item.CreatedAt = createdAt.Format(time.RFC3339)
	item.UpdatedAt = updatedAt.Format(time.RFC3339)
	return item, nil
}

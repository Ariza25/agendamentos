package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/appointment"
)

type AppointmentRepository struct {
	db *sql.DB
}

func NewAppointmentRepository(db *sql.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

func (r *AppointmentRepository) Create(input appointment.Appointment) (appointment.Appointment, error) {
	item := input
	item.ID = uuid.NewString()
	if _, err := r.db.Exec(`
		INSERT INTO catalog.appointments (id, tenant_id, professional_id, service_id, customer_id, starts_at, ends_at, status, notes, created_by)
		VALUES ($1, $2, $3, $4, $5, $6::timestamptz, $7::timestamptz, $8, $9, $10)
	`, item.ID, item.TenantID, item.ProfessionalID, item.ServiceID, item.CustomerID, item.StartsAt, item.EndsAt, item.Status, item.Notes, item.CreatedBy); err != nil {
		return appointment.Appointment{}, fmt.Errorf("insert appointment: %w", err)
	}
	return r.GetByID(item.TenantID, item.ID)
}

func (r *AppointmentRepository) Update(input appointment.Appointment) (appointment.Appointment, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.appointments
		SET professional_id = $3,
		    service_id = COALESCE(NULLIF($4, ''), service_id),
		    customer_id = $5,
		    starts_at = $6::timestamptz,
		    ends_at = $7::timestamptz,
		    status = $8,
		    notes = COALESCE(NULLIF($9, ''), notes),
		    updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`, input.ID, input.TenantID, input.ProfessionalID, strings.TrimSpace(input.ServiceID), input.CustomerID, input.StartsAt, input.EndsAt, input.Status, strings.TrimSpace(input.Notes))
	if err != nil {
		return appointment.Appointment{}, fmt.Errorf("update appointment: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return appointment.Appointment{}, fmt.Errorf("check appointment update: %w", err)
	}
	if affected == 0 {
		return appointment.Appointment{}, fmt.Errorf("appointment not found")
	}
	return r.GetByID(input.TenantID, input.ID)
}

func (r *AppointmentRepository) GetByID(tenantID string, id string) (appointment.Appointment, error) {
	item, err := scanAppointmentRow(r.db.QueryRow(`
		SELECT id, tenant_id, professional_id, service_id, customer_id, starts_at, ends_at, status, notes, created_by, created_at, updated_at
		FROM catalog.appointments
		WHERE id = $1 AND tenant_id = $2
	`, strings.TrimSpace(id), strings.TrimSpace(tenantID)))
	if err != nil {
		return appointment.Appointment{}, err
	}
	return item, nil
}

func (r *AppointmentRepository) ListByRange(tenantID string, from string, to string, professionalID string, customerID string, status string) ([]appointment.Appointment, error) {
	tenantID = strings.TrimSpace(tenantID)
	query := `
		SELECT id, tenant_id, professional_id, service_id, customer_id, starts_at, ends_at, status, notes, created_by, created_at, updated_at
		FROM catalog.appointments
		WHERE tenant_id = $1
	`
	args := []any{tenantID}
	index := 2
	if from != "" {
		query += fmt.Sprintf(" AND ends_at >= $%d::timestamptz", index)
		args = append(args, from)
		index++
	}
	if to != "" {
		query += fmt.Sprintf(" AND starts_at < $%d::timestamptz", index)
		args = append(args, to)
		index++
	}
	if professionalID != "" {
		query += fmt.Sprintf(" AND professional_id = $%d", index)
		args = append(args, professionalID)
		index++
	}
	if customerID != "" {
		query += fmt.Sprintf(" AND customer_id = $%d", index)
		args = append(args, customerID)
		index++
	}
	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", index)
		args = append(args, status)
		index++
	}
	query += " ORDER BY starts_at ASC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query appointments: %w", err)
	}
	defer rows.Close()

	items := []appointment.Appointment{}
	for rows.Next() {
		item, err := scanAppointmentRows(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *AppointmentRepository) UpdateStatus(tenantID string, id string, status string, updatedBy string) (appointment.Appointment, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.appointments
		SET status = $3, updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`, id, tenantID, status)
	if err != nil {
		return appointment.Appointment{}, fmt.Errorf("update appointment status: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return appointment.Appointment{}, fmt.Errorf("check appointment status: %w", err)
	}
	if affected == 0 {
		return appointment.Appointment{}, fmt.Errorf("appointment not found")
	}
	return r.GetByID(tenantID, id)
}

func (r *AppointmentRepository) HasConflict(tenantID string, professionalID string, startsAt string, endsAt string, excludeID string) (bool, error) {
	tenantID = strings.TrimSpace(tenantID)
	professionalID = strings.TrimSpace(professionalID)
	excludeID = strings.TrimSpace(excludeID)
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM catalog.appointments
			WHERE tenant_id = $1
			  AND professional_id = $2
			  AND status NOT IN ('cancelled', 'no_show')
			  AND $3::timestamptz < ends_at
			  AND $4::timestamptz > starts_at
	`
	args := []any{tenantID, professionalID, startsAt, endsAt}
	if excludeID != "" {
		query += " AND id <> $5"
		args = append(args, excludeID)
	}
	query += ")"

	var exists bool
	if err := r.db.QueryRow(query, args...).Scan(&exists); err != nil {
		return false, fmt.Errorf("check appointment conflict: %w", err)
	}
	return exists, nil
}

func (r *AppointmentRepository) CountByStatus(tenantID string, from string, to string) ([]appointment.StatusCount, error) {
	rows, err := r.db.Query(`
		SELECT status, COUNT(*) AS count
		FROM catalog.appointments
		WHERE tenant_id = $1
		  AND ends_at >= $2::timestamptz
		  AND starts_at < $3::timestamptz
		GROUP BY status
	`, strings.TrimSpace(tenantID), from, to)
	if err != nil {
		return nil, fmt.Errorf("count appointments by status: %w", err)
	}
	defer rows.Close()

	items := []appointment.StatusCount{}
	for rows.Next() {
		var item appointment.StatusCount
		if err := rows.Scan(&item.Status, &item.Count); err != nil {
			return nil, fmt.Errorf("scan appointment count: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *AppointmentRepository) Delete(tenantID string, id string) error {
	result, err := r.db.Exec(`
		DELETE FROM catalog.appointments
		WHERE id = $1 AND tenant_id = $2
	`, id, tenantID)
	if err != nil {
		return fmt.Errorf("delete appointment: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check appointment delete: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("appointment not found")
	}
	return nil
}

type appointmentScanner interface {
	Scan(dest ...any) error
}

func scanAppointmentRow(scanner appointmentScanner) (appointment.Appointment, error) {
	return scanAppointmentRows(scanner)
}

func scanAppointmentRows(scanner appointmentScanner) (appointment.Appointment, error) {
	var item appointment.Appointment
	var startsAt, endsAt, createdAt, updatedAt time.Time
	if err := scanner.Scan(&item.ID, &item.TenantID, &item.ProfessionalID, &item.ServiceID, &item.CustomerID, &startsAt, &endsAt, &item.Status, &item.Notes, &item.CreatedBy, &createdAt, &updatedAt); err != nil {
		if err == sql.ErrNoRows {
			return appointment.Appointment{}, fmt.Errorf("appointment not found")
		}
		return appointment.Appointment{}, fmt.Errorf("scan appointment: %w", err)
	}
	item.StartsAt = startsAt.UTC().Format(time.RFC3339)
	item.EndsAt = endsAt.UTC().Format(time.RFC3339)
	item.CreatedAt = createdAt.Format(time.RFC3339)
	item.UpdatedAt = updatedAt.Format(time.RFC3339)
	return item, nil
}

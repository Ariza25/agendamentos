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
	_, err := r.db.Exec(`
		INSERT INTO catalog.appointments (id, tenant_id, professional_id, service_id, resource_id, customer_id, starts_at, ends_at, status, notes, created_by, management_code_hash, series_id, recurrence_rule, timezone, payment_status, no_show_fee_cents, source, cancellation_reason, lgpd_accepted_at, marketing_consent)
		VALUES ($1, $2, $3, $4, $5, $6, $7::timestamptz, $8::timestamptz, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20::timestamptz, $21)
	`, item.ID, item.TenantID, item.ProfessionalID, item.ServiceID, item.ResourceID, item.CustomerID, item.StartsAt, item.EndsAt, item.Status, item.Notes, item.CreatedBy, nullIfEmpty(item.ManagementCodeHash), item.SeriesID, item.RecurrenceRule, item.Timezone, item.PaymentStatus, item.NoShowFeeCents, item.Source, item.CancellationReason, nullIfEmpty(item.LGPDAcceptedAt), item.MarketingConsent)
	if err != nil {
		return appointment.Appointment{}, fmt.Errorf("insert appointment: %w", err)
	}
	return r.GetByID(item.TenantID, item.ID)
}

func (r *AppointmentRepository) Update(input appointment.Appointment) (appointment.Appointment, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.appointments
		SET professional_id = $3,
		    service_id = COALESCE(NULLIF($4, ''), service_id),
		    resource_id = COALESCE(NULLIF($5, ''), resource_id),
		    customer_id = $6,
		    starts_at = $7::timestamptz,
		    ends_at = $8::timestamptz,
		    status = $9,
		    notes = COALESCE(NULLIF($10, ''), notes),
		    management_code_hash = COALESCE(NULLIF($11, ''), management_code_hash),
		    series_id = COALESCE(NULLIF($12, ''), series_id),
		    recurrence_rule = COALESCE(NULLIF($13, ''), recurrence_rule),
		    timezone = COALESCE(NULLIF($14, ''), timezone),
		    payment_status = COALESCE(NULLIF($15, ''), payment_status),
		    no_show_fee_cents = $16,
		    source = COALESCE(NULLIF($17, ''), source),
		    cancellation_reason = COALESCE(NULLIF($18, ''), cancellation_reason),
		    check_in_at = $19,
		    no_show_at = $20,
		    marketing_consent = $21,
		    updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`, input.ID, input.TenantID, input.ProfessionalID, strings.TrimSpace(input.ServiceID), input.ResourceID, input.CustomerID, input.StartsAt, input.EndsAt, input.Status, strings.TrimSpace(input.Notes), nullIfEmpty(input.ManagementCodeHash), input.SeriesID, input.RecurrenceRule, input.Timezone, input.PaymentStatus, input.NoShowFeeCents, input.Source, input.CancellationReason, nullIfEmpty(input.CheckInAt), nullIfEmpty(input.NoShowAt), input.MarketingConsent)
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
	return scanAppointmentRow(r.db.QueryRow(appointmentSelectColumns()+`
		FROM catalog.appointments
		WHERE id = $1 AND tenant_id = $2
	`, strings.TrimSpace(id), strings.TrimSpace(tenantID)))
}

func (r *AppointmentRepository) GetByManagementCodeHash(tenantID string, hash string) (appointment.Appointment, error) {
	return scanAppointmentRow(r.db.QueryRow(appointmentSelectColumns()+`
		FROM catalog.appointments
		WHERE tenant_id = $1 AND management_code_hash = $2
	`, strings.TrimSpace(tenantID), strings.TrimSpace(hash)))
}

func (r *AppointmentRepository) ListByRange(tenantID string, from string, to string, professionalID string, customerID string, status string) ([]appointment.Appointment, error) {
	tenantID = strings.TrimSpace(tenantID)
	query := appointmentSelectColumns() + `
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

func (r *AppointmentRepository) ListDueReminders(from string, to string) ([]appointment.Appointment, error) {
	rows, err := r.db.Query(appointmentSelectColumns()+`
		FROM catalog.appointments
		WHERE starts_at >= $1::timestamptz
		  AND starts_at < $2::timestamptz
		  AND status IN ('scheduled', 'confirmed')
		  AND reminder_sent_at IS NULL
		ORDER BY starts_at ASC
	`, from, to)
	if err != nil {
		return nil, fmt.Errorf("query due reminders: %w", err)
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

func (r *AppointmentRepository) SetReminderSent(tenantID string, id string) error {
	result, err := r.db.Exec(`
		UPDATE catalog.appointments
		SET reminder_sent_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`, id, tenantID)
	if err != nil {
		return fmt.Errorf("set reminder sent: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check reminder sent: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("appointment not found")
	}
	return nil
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

func (r *AppointmentRepository) CountResourceOverlaps(tenantID string, resourceID string, startsAt string, endsAt string, excludeID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM catalog.appointments
		WHERE tenant_id = $1
		  AND resource_id = $2
		  AND status NOT IN ('cancelled', 'no_show')
		  AND $3::timestamptz < ends_at
		  AND $4::timestamptz > starts_at
	`
	args := []any{strings.TrimSpace(tenantID), strings.TrimSpace(resourceID), startsAt, endsAt}
	if excludeID = strings.TrimSpace(excludeID); excludeID != "" {
		query += " AND id <> $5"
		args = append(args, excludeID)
	}
	var count int
	if err := r.db.QueryRow(query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("count resource overlaps: %w", err)
	}
	return count, nil
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

func (r *AppointmentRepository) Analytics(tenantID string, from string, to string) (appointment.Analytics, error) {
	tenantID = strings.TrimSpace(tenantID)
	result := appointment.Analytics{}

	if err := r.db.QueryRow(`
		SELECT
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE status = 'scheduled') AS scheduled,
			COUNT(*) FILTER (WHERE status = 'confirmed') AS confirmed,
			COUNT(*) FILTER (WHERE status = 'completed') AS completed,
			COUNT(*) FILTER (WHERE status = 'cancelled') AS cancelled,
			COUNT(*) FILTER (WHERE status = 'no_show') AS no_show,
			COALESCE(SUM(s.price_cents) FILTER (WHERE a.status = 'completed'), 0) AS revenue
		FROM catalog.appointments a
		LEFT JOIN catalog.services s ON s.id = a.service_id
		WHERE a.tenant_id = $1
		  AND a.ends_at >= $2::timestamptz
		  AND a.starts_at < $3::timestamptz
	`, tenantID, from, to).Scan(&result.TotalAppointments, &result.Scheduled, &result.Confirmed, &result.Completed, &result.Cancelled, &result.NoShow, &result.RevenueCents); err != nil {
		return appointment.Analytics{}, fmt.Errorf("analytics totals: %w", err)
	}
	if result.TotalAppointments > 0 {
		result.NoShowRate = float64(result.NoShow) / float64(result.TotalAppointments) * 100
	}

	serviceRows, err := r.db.Query(`
		SELECT a.service_id, COALESCE(s.name, '') AS name, COUNT(*) AS count, COALESCE(SUM(s.price_cents) FILTER (WHERE a.status = 'completed'), 0) AS revenue
		FROM catalog.appointments a
		LEFT JOIN catalog.services s ON s.id = a.service_id
		WHERE a.tenant_id = $1
		  AND a.ends_at >= $2::timestamptz
		  AND a.starts_at < $3::timestamptz
		GROUP BY a.service_id, s.name
		ORDER BY count DESC
		LIMIT 10
	`, tenantID, from, to)
	if err != nil {
		return appointment.Analytics{}, fmt.Errorf("analytics by service: %w", err)
	}
	defer serviceRows.Close()
	for serviceRows.Next() {
		var metric appointment.ServiceMetric
		if err := serviceRows.Scan(&metric.ServiceID, &metric.ServiceName, &metric.Count, &metric.RevenueCents); err != nil {
			return appointment.Analytics{}, fmt.Errorf("scan analytics service: %w", err)
		}
		result.ByService = append(result.ByService, metric)
	}
	if err := serviceRows.Err(); err != nil {
		return appointment.Analytics{}, fmt.Errorf("analytics by service: %w", err)
	}

	professionalRows, err := r.db.Query(`
		SELECT a.professional_id, COALESCE(p.name, '') AS name, COUNT(*) AS count
		FROM catalog.appointments a
		LEFT JOIN catalog.professionals p ON p.id = a.professional_id
		WHERE a.tenant_id = $1
		  AND a.ends_at >= $2::timestamptz
		  AND a.starts_at < $3::timestamptz
		GROUP BY a.professional_id, p.name
		ORDER BY count DESC
		LIMIT 10
	`, tenantID, from, to)
	if err != nil {
		return appointment.Analytics{}, fmt.Errorf("analytics by professional: %w", err)
	}
	defer professionalRows.Close()
	for professionalRows.Next() {
		var metric appointment.ProfessionalMetric
		if err := professionalRows.Scan(&metric.ProfessionalID, &metric.ProfessionalName, &metric.Count); err != nil {
			return appointment.Analytics{}, fmt.Errorf("scan analytics professional: %w", err)
		}
		result.ByProfessional = append(result.ByProfessional, metric)
	}
	if err := professionalRows.Err(); err != nil {
		return appointment.Analytics{}, fmt.Errorf("analytics by professional: %w", err)
	}

	dayRows, err := r.db.Query(`
		SELECT (a.starts_at AT TIME ZONE 'UTC')::date AS date, COUNT(*) AS count
		FROM catalog.appointments a
		WHERE a.tenant_id = $1
		  AND a.ends_at >= $2::timestamptz
		  AND a.starts_at < $3::timestamptz
		  AND a.status IN ('scheduled', 'confirmed', 'in_progress', 'completed')
		GROUP BY date
		ORDER BY date ASC
	`, tenantID, from, to)
	if err != nil {
		return appointment.Analytics{}, fmt.Errorf("analytics by day: %w", err)
	}
	defer dayRows.Close()
	for dayRows.Next() {
		var metric appointment.DayMetric
		if err := dayRows.Scan(&metric.Date, &metric.Count); err != nil {
			return appointment.Analytics{}, fmt.Errorf("scan analytics day: %w", err)
		}
		result.ByDay = append(result.ByDay, metric)
	}
	if err := dayRows.Err(); err != nil {
		return appointment.Analytics{}, fmt.Errorf("analytics by day: %w", err)
	}

	return result, nil
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

func appointmentSelectColumns() string {
	return `
		SELECT id, tenant_id, professional_id, service_id, resource_id, customer_id, starts_at, ends_at, status, notes, created_by,
		       COALESCE(management_code_hash, ''), series_id, COALESCE(recurrence_rule, ''), COALESCE(timezone, ''), payment_status, no_show_fee_cents, source, cancellation_reason,
		       check_in_at, no_show_at, reminder_sent_at, lgpd_accepted_at, COALESCE(marketing_consent, false), created_at, updated_at`
}

func scanAppointmentRows(scanner appointmentScanner) (appointment.Appointment, error) {
	var item appointment.Appointment
	var startsAt, endsAt, createdAt, updatedAt time.Time
	var checkInAt, noShowAt, reminderSentAt, lgpdAcceptedAt sql.NullTime
	if err := scanner.Scan(&item.ID, &item.TenantID, &item.ProfessionalID, &item.ServiceID, &item.ResourceID, &item.CustomerID, &startsAt, &endsAt, &item.Status, &item.Notes, &item.CreatedBy, &item.ManagementCodeHash, &item.SeriesID, &item.RecurrenceRule, &item.Timezone, &item.PaymentStatus, &item.NoShowFeeCents, &item.Source, &item.CancellationReason, &checkInAt, &noShowAt, &reminderSentAt, &lgpdAcceptedAt, &item.MarketingConsent, &createdAt, &updatedAt); err != nil {
		if err == sql.ErrNoRows {
			return appointment.Appointment{}, fmt.Errorf("appointment not found")
		}
		return appointment.Appointment{}, fmt.Errorf("scan appointment: %w", err)
	}
	item.StartsAt = startsAt.UTC().Format(time.RFC3339)
	item.EndsAt = endsAt.UTC().Format(time.RFC3339)
	item.CreatedAt = createdAt.Format(time.RFC3339)
	item.UpdatedAt = updatedAt.Format(time.RFC3339)
	if checkInAt.Valid {
		item.CheckInAt = checkInAt.Time.UTC().Format(time.RFC3339)
	}
	if noShowAt.Valid {
		item.NoShowAt = noShowAt.Time.UTC().Format(time.RFC3339)
	}
	if reminderSentAt.Valid {
		item.ReminderSentAt = reminderSentAt.Time.UTC().Format(time.RFC3339)
	}
	if lgpdAcceptedAt.Valid {
		item.LGPDAcceptedAt = lgpdAcceptedAt.Time.UTC().Format(time.RFC3339)
	}
	return item, nil
}

func nullIfEmpty(value string) any {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return value
}

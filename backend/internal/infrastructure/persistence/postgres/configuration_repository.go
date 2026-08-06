package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"erpagendamentos/backend/internal/domain/configuration"
)

type ConfigurationRepository struct {
	db *sql.DB
}

func NewConfigurationRepository(db *sql.DB) *ConfigurationRepository {
	return &ConfigurationRepository{db: db}
}

func (r *ConfigurationRepository) Get(tenantID string) (configuration.Configuration, error) {
	tenantID = normalizePostgresTenantID(tenantID)
	var item configuration.Configuration
	var workingHoursRaw []byte
	var updatedAt time.Time
	err := r.db.QueryRow(`
		SELECT tenant_id, establishment_name, working_hours, appointment_interval_minutes, default_duration_minutes, allow_online_booking, advance_notice_hours, updated_at
		FROM catalog.configurations
		WHERE tenant_id = $1
	`, tenantID).Scan(&item.TenantID, &item.EstablishmentName, &workingHoursRaw, &item.AppointmentIntervalMinutes, &item.DefaultDurationMinutes, &item.AllowOnlineBooking, &item.AdvanceNoticeHours, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return configuration.Configuration{
				TenantID:                   tenantID,
				EstablishmentName:          "Estabelecimento",
				WorkingHours:               []configuration.WorkingHours{},
				AppointmentIntervalMinutes: 30,
				DefaultDurationMinutes:     30,
				AllowOnlineBooking:         true,
				AdvanceNoticeHours:         0,
			}, nil
		}
		return configuration.Configuration{}, fmt.Errorf("get configuration: %w", err)
	}
	item.WorkingHours = []configuration.WorkingHours{}
	if len(workingHoursRaw) > 0 {
		_ = json.Unmarshal(workingHoursRaw, &item.WorkingHours)
	}
	if item.WorkingHours == nil {
		item.WorkingHours = []configuration.WorkingHours{}
	}
	item.UpdatedAt = updatedAt.Format(time.RFC3339)
	return item, nil
}

func (r *ConfigurationRepository) Update(input configuration.Configuration) (configuration.Configuration, error) {
	tenantID := normalizePostgresTenantID(input.TenantID)
	input.TenantID = tenantID
	if input.WorkingHours == nil {
		input.WorkingHours = []configuration.WorkingHours{}
	}
	workingHoursRaw, err := json.Marshal(input.WorkingHours)
	if err != nil {
		return configuration.Configuration{}, fmt.Errorf("marshal working hours: %w", err)
	}
	if _, err := r.db.Exec(`
		INSERT INTO catalog.configurations (tenant_id, establishment_name, working_hours, appointment_interval_minutes, default_duration_minutes, allow_online_booking, advance_notice_hours, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		ON CONFLICT (tenant_id) DO UPDATE SET
			establishment_name = EXCLUDED.establishment_name,
			working_hours = EXCLUDED.working_hours,
			appointment_interval_minutes = EXCLUDED.appointment_interval_minutes,
			default_duration_minutes = EXCLUDED.default_duration_minutes,
			allow_online_booking = EXCLUDED.allow_online_booking,
			advance_notice_hours = EXCLUDED.advance_notice_hours,
			updated_at = NOW()
	`, tenantID, strings.TrimSpace(input.EstablishmentName), workingHoursRaw, input.AppointmentIntervalMinutes, input.DefaultDurationMinutes, input.AllowOnlineBooking, input.AdvanceNoticeHours); err != nil {
		return configuration.Configuration{}, fmt.Errorf("update configuration: %w", err)
	}
	return r.Get(tenantID)
}

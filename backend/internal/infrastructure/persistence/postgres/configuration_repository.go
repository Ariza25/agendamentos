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
	var closedDaysRaw []byte
	var updatedAt time.Time
	err := r.db.QueryRow(`
		SELECT tenant_id, establishment_name, working_hours, appointment_interval_minutes, default_duration_minutes,
		       allow_online_booking, advance_notice_hours, buffer_minutes, closed_days, notifications_enabled,
		       notification_reminder_hours, notification_channel, deposit_enabled, deposit_type, deposit_amount_cents,
		       deposit_percentage, no_show_fee_cents, cancellation_policy_hours, calendar_sync_enabled, updated_at
		FROM catalog.configurations
		WHERE tenant_id = $1
	`, tenantID).Scan(&item.TenantID, &item.EstablishmentName, &workingHoursRaw, &item.AppointmentIntervalMinutes, &item.DefaultDurationMinutes, &item.AllowOnlineBooking, &item.AdvanceNoticeHours, &item.BufferMinutes, &closedDaysRaw, &item.NotificationsEnabled, &item.NotificationReminderHours, &item.NotificationChannel, &item.DepositEnabled, &item.DepositType, &item.DepositAmountCents, &item.DepositPercentage, &item.NoShowFeeCents, &item.CancellationPolicyHours, &item.CalendarSyncEnabled, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return r.defaultConfiguration(tenantID), nil
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
	item.ClosedDays = []configuration.ClosedDay{}
	if len(closedDaysRaw) > 0 {
		_ = json.Unmarshal(closedDaysRaw, &item.ClosedDays)
	}
	if item.ClosedDays == nil {
		item.ClosedDays = []configuration.ClosedDay{}
	}
	item.UpdatedAt = updatedAt.Format(time.RFC3339)
	return item, nil
}

func (r *ConfigurationRepository) defaultConfiguration(tenantID string) configuration.Configuration {
	return configuration.Configuration{
		TenantID:                   tenantID,
		EstablishmentName:          "Estabelecimento",
		WorkingHours:               []configuration.WorkingHours{},
		AppointmentIntervalMinutes: 30,
		DefaultDurationMinutes:     30,
		AllowOnlineBooking:         true,
		AdvanceNoticeHours:         0,
		BufferMinutes:              0,
		ClosedDays:                 []configuration.ClosedDay{},
		NotificationsEnabled:       false,
		NotificationReminderHours:  24,
		NotificationChannel:        "whatsapp",
		DepositEnabled:             false,
		DepositType:                "fixed",
		DepositAmountCents:         0,
		DepositPercentage:          0,
		NoShowFeeCents:             0,
		CancellationPolicyHours:    0,
		CalendarSyncEnabled:        false,
	}
}

func (r *ConfigurationRepository) Update(input configuration.Configuration) (configuration.Configuration, error) {
	tenantID := normalizePostgresTenantID(input.TenantID)
	input.TenantID = tenantID
	if input.WorkingHours == nil {
		input.WorkingHours = []configuration.WorkingHours{}
	}
	if input.ClosedDays == nil {
		input.ClosedDays = []configuration.ClosedDay{}
	}
	workingHoursRaw, err := json.Marshal(input.WorkingHours)
	if err != nil {
		return configuration.Configuration{}, fmt.Errorf("marshal working hours: %w", err)
	}
	closedDaysRaw, err := json.Marshal(input.ClosedDays)
	if err != nil {
		return configuration.Configuration{}, fmt.Errorf("marshal closed days: %w", err)
	}
	if _, err := r.db.Exec(`
		INSERT INTO catalog.configurations (tenant_id, establishment_name, working_hours, appointment_interval_minutes, default_duration_minutes, allow_online_booking, advance_notice_hours, buffer_minutes, closed_days, notifications_enabled, notification_reminder_hours, notification_channel, deposit_enabled, deposit_type, deposit_amount_cents, deposit_percentage, no_show_fee_cents, cancellation_policy_hours, calendar_sync_enabled, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, NOW())
		ON CONFLICT (tenant_id) DO UPDATE SET
			establishment_name = EXCLUDED.establishment_name,
			working_hours = EXCLUDED.working_hours,
			appointment_interval_minutes = EXCLUDED.appointment_interval_minutes,
			default_duration_minutes = EXCLUDED.default_duration_minutes,
			allow_online_booking = EXCLUDED.allow_online_booking,
			advance_notice_hours = EXCLUDED.advance_notice_hours,
			buffer_minutes = EXCLUDED.buffer_minutes,
			closed_days = EXCLUDED.closed_days,
			notifications_enabled = EXCLUDED.notifications_enabled,
			notification_reminder_hours = EXCLUDED.notification_reminder_hours,
			notification_channel = EXCLUDED.notification_channel,
			deposit_enabled = EXCLUDED.deposit_enabled,
			deposit_type = EXCLUDED.deposit_type,
			deposit_amount_cents = EXCLUDED.deposit_amount_cents,
			deposit_percentage = EXCLUDED.deposit_percentage,
			no_show_fee_cents = EXCLUDED.no_show_fee_cents,
			cancellation_policy_hours = EXCLUDED.cancellation_policy_hours,
			calendar_sync_enabled = EXCLUDED.calendar_sync_enabled,
			updated_at = NOW()
	`, tenantID, strings.TrimSpace(input.EstablishmentName), workingHoursRaw, input.AppointmentIntervalMinutes, input.DefaultDurationMinutes, input.AllowOnlineBooking, input.AdvanceNoticeHours, input.BufferMinutes, closedDaysRaw, input.NotificationsEnabled, input.NotificationReminderHours, input.NotificationChannel, input.DepositEnabled, input.DepositType, input.DepositAmountCents, input.DepositPercentage, input.NoShowFeeCents, input.CancellationPolicyHours, input.CalendarSyncEnabled); err != nil {
		return configuration.Configuration{}, fmt.Errorf("update configuration: %w", err)
	}
	return r.Get(tenantID)
}

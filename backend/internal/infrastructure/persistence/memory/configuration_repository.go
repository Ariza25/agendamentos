package memory

import (
	"strings"
	"sync"
	"time"

	"erpagendamentos/backend/internal/domain/configuration"
)

type ConfigurationRepository struct {
	mu             sync.Mutex
	configurations map[string]configuration.Configuration
}

func NewConfigurationRepository() *ConfigurationRepository {
	return &ConfigurationRepository{
		configurations: make(map[string]configuration.Configuration),
	}
}

func (r *ConfigurationRepository) Get(tenantID string) (configuration.Configuration, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID = normalizeConfigurationTenantID(tenantID)
	item, exists := r.configurations[tenantID]
	if !exists {
		return defaultMemoryConfiguration(tenantID), nil
	}
	return cloneConfiguration(item), nil
}

func (r *ConfigurationRepository) Update(input configuration.Configuration) (configuration.Configuration, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID := normalizeConfigurationTenantID(input.TenantID)
	input.TenantID = tenantID
	input.UpdatedAt = time.Now().Format(time.RFC3339)
	if input.WorkingHours == nil {
		input.WorkingHours = []configuration.WorkingHours{}
	}
	if input.ClosedDays == nil {
		input.ClosedDays = []configuration.ClosedDay{}
	}
	r.configurations[tenantID] = input
	return cloneConfiguration(input), nil
}

func defaultMemoryConfiguration(tenantID string) configuration.Configuration {
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
		UpdatedAt:                  time.Now().Format(time.RFC3339),
	}
}

func normalizeConfigurationTenantID(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "default"
	}
	return value
}

func cloneConfiguration(value configuration.Configuration) configuration.Configuration {
	if value.WorkingHours == nil {
		value.WorkingHours = []configuration.WorkingHours{}
	}
	if value.ClosedDays == nil {
		value.ClosedDays = []configuration.ClosedDay{}
	}
	return value
}

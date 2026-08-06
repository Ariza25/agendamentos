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
		return configuration.Configuration{
			TenantID:                   tenantID,
			EstablishmentName:          "Estabelecimento",
			WorkingHours:               []configuration.WorkingHours{},
			AppointmentIntervalMinutes: 30,
			DefaultDurationMinutes:     30,
			AllowOnlineBooking:         true,
			AdvanceNoticeHours:         0,
			UpdatedAt:                  time.Now().Format(time.RFC3339),
		}, nil
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
	r.configurations[tenantID] = input
	return cloneConfiguration(input), nil
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
	return value
}

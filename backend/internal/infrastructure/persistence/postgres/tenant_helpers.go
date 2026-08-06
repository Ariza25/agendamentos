package postgres

import "strings"

func normalizePostgresTenantID(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "default"
	}
	return value
}

func isTenantIdentityConflict(err error) bool {
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "23505") ||
		strings.Contains(message, "idx_catalog_tenants_active_slug") ||
		strings.Contains(message, "tenant name already exists") ||
		strings.Contains(message, "tenant slug already exists")
}

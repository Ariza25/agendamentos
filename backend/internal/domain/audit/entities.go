package audit

import "time"

type AuditLog struct {
	ID         string         `json:"id"`
	TenantID   string         `json:"tenantId"`
	UserID     string         `json:"userId"`
	UserEmail  string         `json:"userEmail"`
	Action     string         `json:"action"`     // ex: "checkout.completed", "cash_register.close", "tenant.update_plan", "customer.delete"
	Resource   string         `json:"resource"`   // ex: "appointment", "cash_register", "tenant", "customer", "commission"
	ResourceID string         `json:"resourceId"`
	Details    map[string]any `json:"details"`
	IPAddress  string         `json:"ipAddress"`
	UserAgent  string         `json:"userAgent"`
	CreatedAt  time.Time      `json:"createdAt"`
}

type Repository interface {
	Record(entry AuditLog) error
	List(tenantID string, limit int, offset int) ([]AuditLog, int, error)
}


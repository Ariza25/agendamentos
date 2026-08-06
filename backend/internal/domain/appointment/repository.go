package appointment

type Appointment struct {
	ID             string `json:"id"`
	TenantID       string `json:"tenantId"`
	ProfessionalID string `json:"professionalId"`
	ServiceID      string `json:"serviceId"`
	CustomerID     string `json:"customerId"`
	StartsAt       string `json:"startsAt"`
	EndsAt         string `json:"endsAt"`
	Status         string `json:"status"`
	Notes          string `json:"notes"`
	CreatedBy      string `json:"createdBy"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type Repository interface {
	Create(input Appointment) (Appointment, error)
	Update(input Appointment) (Appointment, error)
	GetByID(tenantID string, id string) (Appointment, error)
	ListByRange(tenantID string, from string, to string, professionalID string, customerID string, status string) ([]Appointment, error)
	UpdateStatus(tenantID string, id string, status string, updatedBy string) (Appointment, error)
	HasConflict(tenantID string, professionalID string, startsAt string, endsAt string, excludeID string) (bool, error)
	CountByStatus(tenantID string, from string, to string) ([]StatusCount, error)
	Delete(tenantID string, id string) error
}

package waitinglist

const (
	StatusWaiting   = "waiting"
	StatusNotified  = "notified"
	StatusConverted = "converted"
	StatusCancelled = "cancelled"
)

type Entry struct {
	ID             string `json:"id"`
	TenantID       string `json:"tenantId"`
	ServiceID      string `json:"serviceId,omitempty"`
	ProfessionalID string `json:"professionalId,omitempty"`
	CustomerID     string `json:"customerId,omitempty"`
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	Status         string `json:"status"`
	CreatedAt      string `json:"createdAt"`
	NotifiedAt     string `json:"notifiedAt,omitempty"`
}

type Repository interface {
	Create(input Entry) (Entry, error)
	ListByTenant(tenantID string, status string) ([]Entry, error)
	GetByID(tenantID string, id string) (Entry, error)
	UpdateStatus(tenantID string, id string, status string) (Entry, error)
	ListWaitingForSlot(tenantID string, serviceID string, professionalID string, order string) ([]Entry, error)
}

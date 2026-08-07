package review

type Review struct {
	ID            string `json:"id"`
	TenantID      string `json:"tenantId"`
	TenantSlug    string `json:"tenantSlug,omitempty"`
	AppointmentID string `json:"appointmentId,omitempty"`
	CustomerID    string `json:"customerId,omitempty"`
	CustomerName  string `json:"customerName,omitempty"`
	Status        string `json:"status"`
	Rating        int    `json:"rating,omitempty"`
	Comment       string `json:"comment,omitempty"`
	InvitedAt     string `json:"invitedAt"`
	SubmittedAt   string `json:"submittedAt,omitempty"`
	CreatedAt     string `json:"createdAt"`
}

type Repository interface {
	Create(input Review) (Review, error)
	GetByID(tenantID string, id string) (Review, error)
	GetByAppointment(tenantID string, appointmentID string) (Review, error)
	ListByTenant(tenantID string, status string) ([]Review, error)
	Submit(tenantID string, id string, rating int, comment string) (Review, error)
}

package payment

type Method string

const (
	MethodPix  Method = "pix"
	MethodCard Method = "card"
)

type Status string

const (
	StatusPending  Status = "pending"
	StatusPaid     Status = "paid"
	StatusFailed   Status = "failed"
	StatusRefunded Status = "refunded"
)

type Payment struct {
	ID                string `json:"id"`
	TenantID          string `json:"tenantId"`
	AppointmentID     string `json:"appointmentId,omitempty"`
	AmountCents       int    `json:"amountCents"`
	Method            string `json:"method"`
	Provider          string `json:"provider"`
	ProviderPaymentID string `json:"providerPaymentId,omitempty"`
	Status            string `json:"status"`
	Description       string `json:"description,omitempty"`
	PaidAt            string `json:"paidAt,omitempty"`
	CreatedAt         string `json:"createdAt"`
}

type IntentInput struct {
	AmountCents int
	Description string
}

type IntentResult struct {
	ProviderPaymentID string
	Status            string
	ProviderData      string
}

type Provider interface {
	CreatePayment(input IntentInput) (IntentResult, error)
	GetPaymentStatus(providerPaymentID string) (string, error)
	Refund(providerPaymentID string) error
}

type Repository interface {
	Create(input Payment) (Payment, error)
	GetByID(tenantID string, id string) (Payment, error)
	GetByAppointment(tenantID string, appointmentID string) (Payment, error)
	ListByTenant(tenantID string, limit int) ([]Payment, error)
	ListByTenantRange(tenantID string, from string, to string) ([]Payment, error)
	UpdateStatus(tenantID string, id string, status string, providerPaymentID string) (Payment, error)
}

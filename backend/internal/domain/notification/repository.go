package notification

type Type string

const (
	TypeConfirmation Type = "confirmation"
	TypeReminder     Type = "reminder"
	TypePostService  Type = "post_service"
	TypeNoShowFee    Type = "no_show_fee"
	TypeWaitingList  Type = "waiting_list"
)

type Log struct {
	ID            string `json:"id"`
	TenantID      string `json:"tenantId"`
	AppointmentID string `json:"appointmentId,omitempty"`
	Type          string `json:"type"`
	Channel       string `json:"channel"`
	Recipient     string `json:"recipient"`
	Subject       string `json:"subject"`
	Body          string `json:"body"`
	Status        string `json:"status"`
	Error         string `json:"error,omitempty"`
	CreatedAt     string `json:"createdAt"`
	SentAt        string `json:"sentAt"`
}

type Repository interface {
	Create(input Log) (Log, error)
	ListByTenant(tenantID string, limit int) ([]Log, error)
}

type Provider interface {
	Send(channel string, recipient string, subject string, body string) error
}

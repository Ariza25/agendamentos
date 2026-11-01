package commission

type CommissionRule struct {
	ID               string  `json:"id"`
	TenantID         string  `json:"tenantId"`
	ProfessionalID   string  `json:"professionalId,omitempty"`
	ServiceID        string  `json:"serviceId,omitempty"`
	Percentage       float64 `json:"percentage"` // e.g. 50.0 for 50%
	FixedAmountCents int     `json:"fixedAmountCents"`
	CreatedAt        string  `json:"createdAt"`
	UpdatedAt        string  `json:"updatedAt"`
}

type CommissionRecord struct {
	ID                    string `json:"id"`
	TenantID              string `json:"tenantId"`
	AppointmentID         string `json:"appointmentId"`
	ProfessionalID        string `json:"professionalId"`
	ProfessionalName      string `json:"professionalName,omitempty"`
	ServiceID             string `json:"serviceId"`
	ServiceName           string `json:"serviceName,omitempty"`
	BaseAmountCents       int    `json:"baseAmountCents"`
	CommissionAmountCents int    `json:"commissionAmountCents"`
	Status                string `json:"status"` // "pending", "paid", "cancelled"
	PayoutID              string `json:"payoutId,omitempty"`
	SettledAt             string `json:"settledAt,omitempty"`
	Notes                 string `json:"notes,omitempty"`
	CreatedAt             string `json:"createdAt"`
}

type Summary struct {
	ProfessionalID   string `json:"professionalId"`
	ProfessionalName string `json:"professionalName"`
	PendingCents     int    `json:"pendingCents"`
	PaidCents        int    `json:"paidCents"`
	TotalCents       int    `json:"totalCents"`
	Count            int    `json:"count,omitempty"`
}

type Payout struct {
	ID                     string `json:"id"`
	TenantID               string `json:"tenantId"`
	ProfessionalID         string `json:"professionalId"`
	ProfessionalName       string `json:"professionalName,omitempty"`
	TotalCents             int    `json:"totalCents"`
	RecordsCount           int    `json:"recordsCount"`
	FromDate               string `json:"fromDate,omitempty"`
	ToDate                 string `json:"toDate,omitempty"`
	PaymentMethod          string `json:"paymentMethod"`
	Notes                  string `json:"notes,omitempty"`
	FinancialTransactionID string `json:"financialTransactionId,omitempty"`
	CreatedByUserID        string `json:"createdByUserId,omitempty"`
	PaidAt                 string `json:"paidAt"`
	CreatedAt              string `json:"createdAt"`
}

type Repository interface {
	CreateRule(rule CommissionRule) (CommissionRule, error)
	UpdateRule(rule CommissionRule) (CommissionRule, error)
	DeleteRule(tenantID string, id string) error
	ListRules(tenantID string) ([]CommissionRule, error)
	FindMatchingRule(tenantID string, professionalID string, serviceID string) (CommissionRule, error)

	CreateRecord(record CommissionRecord) (CommissionRecord, error)
	ListRecords(tenantID string, professionalID string, status string) ([]CommissionRecord, error)
	SettleRecords(tenantID string, professionalID string, recordIDs []string) error
	GetSummaries(tenantID string) ([]Summary, error)

	CreatePayout(payout Payout, recordIDs []string) (Payout, error)
	ListPayouts(tenantID string, professionalID string) ([]Payout, error)
}

package finance

type Transaction struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenantId"`
	Type        string `json:"type"`     // "income", "expense"
	Category    string `json:"category"` // "service_sale", "product_sale", "rent", "supplies", "commission_payment", "salary", "utility", "marketing", "other"
	AmountCents int    `json:"amountCents"`
	Description string `json:"description"`
	ReferenceID    string `json:"referenceId,omitempty"`
	CashRegisterID string `json:"cashRegisterId,omitempty"`
	PaymentMethod  string `json:"paymentMethod,omitempty"`
	DueDate        string `json:"dueDate"`
	PaidAt      string `json:"paidAt,omitempty"`
	Status      string `json:"status"` // "pending", "paid", "cancelled"
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type Summary struct {
	TotalIncomeCents   int64           `json:"totalIncomeCents"`
	TotalExpenseCents  int64           `json:"totalExpenseCents"`
	NetProfitCents     int64           `json:"netProfitCents"`
	PendingIncomeCents int64           `json:"pendingIncomeCents"`
	PendingExpenseCents int64          `json:"pendingExpenseCents"`
	ByCategory         []CategoryTotal `json:"byCategory"`
}

type CategoryTotal struct {
	Category    string `json:"category"`
	Type        string `json:"type"`
	AmountCents int64  `json:"amountCents"`
}

type Repository interface {
	Create(tx Transaction) (Transaction, error)
	Update(tx Transaction) (Transaction, error)
	GetByID(tenantID string, id string) (Transaction, error)
	List(tenantID string, txType string, status string, fromDate string, toDate string) ([]Transaction, error)
	Delete(tenantID string, id string) error
	GetSummary(tenantID string, fromDate string, toDate string) (Summary, error)
}


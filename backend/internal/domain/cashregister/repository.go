package cashregister

type CashRegister struct {
	ID                  string `json:"id"`
	TenantID            string `json:"tenantId"`
	OpenedByUserID      string `json:"openedByUserId"`
	ClosedByUserID      string `json:"closedByUserId,omitempty"`
	Status              string `json:"status"` // "open", "closed"
	InitialBalanceCents int    `json:"initialBalanceCents"`
	ClosingBalanceCents int    `json:"closingBalanceCents"`
	SystemBalanceCents  int    `json:"systemBalanceCents"`
	DifferenceCents     int    `json:"differenceCents"`
	OpenedAt            string `json:"openedAt"`
	ClosedAt            string `json:"closedAt,omitempty"`
	Notes               string `json:"notes,omitempty"`
}

type CashRegisterEntry struct {
	ID              string `json:"id"`
	TenantID        string `json:"tenantId"`
	CashRegisterID  string `json:"cashRegisterId"`
	Type            string `json:"type"` // "sangria", "suprimento"
	AmountCents     int    `json:"amountCents"`
	Reason          string `json:"reason"`
	CreatedByUserID string `json:"createdByUserId"`
	CreatedAt       string `json:"createdAt"`
}

type ShiftSummary struct {
	CashRegister          CashRegister        `json:"cashRegister"`
	Entries               []CashRegisterEntry `json:"entries"`
	TotalSangriasCents    int                 `json:"totalSangriasCents"`
	TotalSuprimentosCents int                 `json:"totalSuprimentosCents"`
	TotalCashIncomeCents  int                 `json:"totalCashIncomeCents"`
	TotalCardIncomeCents  int                 `json:"totalCardIncomeCents"`
	TotalPixIncomeCents   int                 `json:"totalPixIncomeCents"`
	TotalOtherIncomeCents int                 `json:"totalOtherIncomeCents"`
	CurrentDrawerCents    int                 `json:"currentDrawerCents"`
}

type Repository interface {
	Create(register CashRegister) (CashRegister, error)
	Update(register CashRegister) (CashRegister, error)
	GetByID(tenantID string, id string) (CashRegister, error)
	GetCurrentOpen(tenantID string) (CashRegister, error)
	List(tenantID string, limit int) ([]CashRegister, error)
	CreateEntry(entry CashRegisterEntry) (CashRegisterEntry, error)
	ListEntries(tenantID string, cashRegisterID string) ([]CashRegisterEntry, error)
}


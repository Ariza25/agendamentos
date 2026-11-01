package commission

import (
	"errors"
	"fmt"
	"time"

	"erpagendamentos/backend/internal/domain/commission"
	financedomain "erpagendamentos/backend/internal/domain/finance"
	prodomain "erpagendamentos/backend/internal/domain/professional"
)

type PayoutInput struct {
	ProfessionalID  string   `json:"professionalId"`
	RecordIDs       []string `json:"recordIds,omitempty"`
	FromDate        string   `json:"fromDate,omitempty"`
	ToDate          string   `json:"toDate,omitempty"`
	PaymentMethod   string   `json:"paymentMethod"` // "pix", "cash", "transfer", "other"
	Notes           string   `json:"notes,omitempty"`
	CreatedByUserID string   `json:"createdByUserId,omitempty"`
}

type UseCase struct {
	repo             commission.Repository
	financeRepo      financedomain.Repository
	professionalRepo prodomain.Repository
}

func NewUseCase(repo commission.Repository, optionalRepos ...any) *UseCase {
	uc := &UseCase{repo: repo}
	for _, r := range optionalRepos {
		if fr, ok := r.(financedomain.Repository); ok {
			uc.financeRepo = fr
		}
		if pr, ok := r.(prodomain.Repository); ok {
			uc.professionalRepo = pr
		}
	}
	return uc
}

func (u *UseCase) CreateRule(rule commission.CommissionRule) (commission.CommissionRule, error) {
	if rule.TenantID == "" {
		return commission.CommissionRule{}, errors.New("tenant id is required")
	}
	if rule.Percentage < 0 || rule.Percentage > 100 {
		return commission.CommissionRule{}, errors.New("percentage must be between 0 and 100")
	}
	return u.repo.CreateRule(rule)
}

func (u *UseCase) UpdateRule(rule commission.CommissionRule) (commission.CommissionRule, error) {
	if rule.Percentage < 0 || rule.Percentage > 100 {
		return commission.CommissionRule{}, errors.New("percentage must be between 0 and 100")
	}
	return u.repo.UpdateRule(rule)
}

func (u *UseCase) DeleteRule(tenantID string, id string) error {
	return u.repo.DeleteRule(tenantID, id)
}

func (u *UseCase) ListRules(tenantID string) ([]commission.CommissionRule, error) {
	return u.repo.ListRules(tenantID)
}

func (u *UseCase) CalculateAndRecord(tenantID string, appointmentID string, professionalID string, professionalName string, serviceID string, serviceName string, servicePriceCents int) (commission.CommissionRecord, error) {
	rule, err := u.repo.FindMatchingRule(tenantID, professionalID, serviceID)
	if err != nil {
		// No rule, record with 0 or skip
		return commission.CommissionRecord{}, err
	}

	commissionCents := int(float64(servicePriceCents) * (rule.Percentage / 100.0))
	if rule.FixedAmountCents > 0 {
		commissionCents += rule.FixedAmountCents
	}

	record := commission.CommissionRecord{
		TenantID:              tenantID,
		AppointmentID:         appointmentID,
		ProfessionalID:        professionalID,
		ProfessionalName:      professionalName,
		ServiceID:             serviceID,
		ServiceName:           serviceName,
		BaseAmountCents:       servicePriceCents,
		CommissionAmountCents: commissionCents,
		Status:                "pending",
	}

	return u.repo.CreateRecord(record)
}

func (u *UseCase) ListRecords(tenantID string, professionalID string, status string) ([]commission.CommissionRecord, error) {
	return u.repo.ListRecords(tenantID, professionalID, status)
}

func (u *UseCase) SettleCommissions(tenantID string, professionalID string, recordIDs []string) error {
	return u.repo.SettleRecords(tenantID, professionalID, recordIDs)
}

func (u *UseCase) GetSummaries(tenantID string) ([]commission.Summary, error) {
	return u.repo.GetSummaries(tenantID)
}

func (u *UseCase) ExecutePayout(tenantID string, input PayoutInput) (commission.Payout, error) {
	if tenantID == "" || input.ProfessionalID == "" {
		return commission.Payout{}, errors.New("tenantId and professionalId are required")
	}

	// 1. Fetch pending records for professional
	records, err := u.repo.ListRecords(tenantID, input.ProfessionalID, "pending")
	if err != nil {
		return commission.Payout{}, fmt.Errorf("list pending records: %w", err)
	}

	targetMap := make(map[string]bool)
	for _, id := range input.RecordIDs {
		targetMap[id] = true
	}

	totalCents := 0
	eligibleRecordIDs := []string{}
	for _, rec := range records {
		if len(targetMap) > 0 && !targetMap[rec.ID] {
			continue
		}
		totalCents += rec.CommissionAmountCents
		eligibleRecordIDs = append(eligibleRecordIDs, rec.ID)
	}

	if totalCents == 0 || len(eligibleRecordIDs) == 0 {
		return commission.Payout{}, errors.New("nenhuma comissão pendente encontrada para liquidar")
	}

	// 2. Resolve professional name
	profName := "Profissional"
	if u.professionalRepo != nil {
		if prof, err := u.professionalRepo.GetByID(tenantID, input.ProfessionalID); err == nil && prof.Name != "" {
			profName = prof.Name
		}
	}

	method := input.PaymentMethod
	if method == "" {
		method = "pix"
	}

	// 3. Create financial transaction (expense)
	var txID string
	if u.financeRepo != nil {
		desc := fmt.Sprintf("Repasse de Comissões - %s (%d atendimentos)", profName, len(eligibleRecordIDs))
		tx, err := u.financeRepo.Create(financedomain.Transaction{
			TenantID:      tenantID,
			Type:          "expense",
			Category:      "Comissões",
			Description:   desc,
			AmountCents:   totalCents,
			DueDate:       time.Now().Format("2006-01-02"),
			Status:        "paid",
			PaidAt:        time.Now().Format(time.RFC3339),
			PaymentMethod: method,
		})
		if err == nil {
			txID = tx.ID
		}
	}

	// 4. Record Payout
	payout := commission.Payout{
		TenantID:               tenantID,
		ProfessionalID:         input.ProfessionalID,
		ProfessionalName:       profName,
		TotalCents:             totalCents,
		RecordsCount:           len(eligibleRecordIDs),
		FromDate:               input.FromDate,
		ToDate:                 input.ToDate,
		PaymentMethod:          method,
		Notes:                  input.Notes,
		FinancialTransactionID: txID,
		CreatedByUserID:        input.CreatedByUserID,
	}

	savedPayout, err := u.repo.CreatePayout(payout, eligibleRecordIDs)
	if err != nil {
		return commission.Payout{}, fmt.Errorf("create payout: %w", err)
	}
	return savedPayout, nil
}

func (u *UseCase) ListPayouts(tenantID string, professionalID string) ([]commission.Payout, error) {
	return u.repo.ListPayouts(tenantID, professionalID)
}

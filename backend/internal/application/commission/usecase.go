package commission

import (
	"errors"

	"erpagendamentos/backend/internal/domain/commission"
)

type UseCase struct {
	repo commission.Repository
}

func NewUseCase(repo commission.Repository) *UseCase {
	return &UseCase{repo: repo}
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


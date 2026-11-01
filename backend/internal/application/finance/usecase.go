package finance

import (
	"errors"
	"strings"

	"erpagendamentos/backend/internal/domain/finance"
)

type UseCase struct {
	repo finance.Repository
}

func NewUseCase(repo finance.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) Create(input finance.Transaction) (finance.Transaction, error) {
	if input.TenantID == "" {
		return finance.Transaction{}, errors.New("tenant id is required")
	}
	input.Description = strings.TrimSpace(input.Description)
	if input.Description == "" {
		return finance.Transaction{}, errors.New("description is required")
	}
	if input.Type != "income" && input.Type != "expense" {
		return finance.Transaction{}, errors.New("type must be income or expense")
	}
	if input.Status == "" {
		input.Status = "pending"
	}
	return u.repo.Create(input)
}

func (u *UseCase) Update(input finance.Transaction) (finance.Transaction, error) {
	input.Description = strings.TrimSpace(input.Description)
	if input.Description == "" {
		return finance.Transaction{}, errors.New("description is required")
	}
	return u.repo.Update(input)
}

func (u *UseCase) GetByID(tenantID string, id string) (finance.Transaction, error) {
	return u.repo.GetByID(tenantID, id)
}

func (u *UseCase) List(tenantID string, txType string, status string, fromDate string, toDate string) ([]finance.Transaction, error) {
	return u.repo.List(tenantID, txType, status, fromDate, toDate)
}

func (u *UseCase) Delete(tenantID string, id string) error {
	return u.repo.Delete(tenantID, id)
}

func (u *UseCase) GetSummary(tenantID string, fromDate string, toDate string) (finance.Summary, error) {
	return u.repo.GetSummary(tenantID, fromDate, toDate)
}

func (u *UseCase) RecordAppointmentIncome(tenantID string, appointmentID string, amountCents int, description string) error {
	_, err := u.repo.Create(finance.Transaction{
		TenantID:    tenantID,
		Type:        "income",
		Category:    "service_sale",
		AmountCents: amountCents,
		Description: description,
		ReferenceID: appointmentID,
		Status:      "paid",
	})
	return err
}


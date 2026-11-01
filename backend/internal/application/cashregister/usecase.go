package cashregister

import (
	"errors"
	"strings"

	"erpagendamentos/backend/internal/domain/cashregister"
	"erpagendamentos/backend/internal/domain/finance"
)

var (
	ErrShiftAlreadyOpen = errors.New("a cash register shift is already open")
	ErrNoOpenShift      = errors.New("no cash register shift is currently open")
)

type UseCase struct {
	repo        cashregister.Repository
	financeRepo finance.Repository
}

func NewUseCase(repo cashregister.Repository, financeRepo finance.Repository) *UseCase {
	return &UseCase{
		repo:        repo,
		financeRepo: financeRepo,
	}
}

func (u *UseCase) OpenShift(tenantID string, userID string, initialBalanceCents int, notes string) (cashregister.CashRegister, error) {
	if tenantID == "" {
		return cashregister.CashRegister{}, errors.New("tenant id is required")
	}
	if initialBalanceCents < 0 {
		return cashregister.CashRegister{}, errors.New("initial balance cannot be negative")
	}

	// Check if already open
	if current, err := u.repo.GetCurrentOpen(tenantID); err == nil && current.ID != "" {
		return cashregister.CashRegister{}, ErrShiftAlreadyOpen
	}

	return u.repo.Create(cashregister.CashRegister{
		TenantID:            tenantID,
		OpenedByUserID:      userID,
		InitialBalanceCents: initialBalanceCents,
		SystemBalanceCents:  initialBalanceCents,
		Notes:               strings.TrimSpace(notes),
	})
}

func (u *UseCase) RecordEntry(tenantID string, userID string, entryType string, amountCents int, reason string) (cashregister.CashRegisterEntry, error) {
	if entryType != "sangria" && entryType != "suprimento" {
		return cashregister.CashRegisterEntry{}, errors.New("entry type must be 'sangria' or 'suprimento'")
	}
	if amountCents <= 0 {
		return cashregister.CashRegisterEntry{}, errors.New("amount must be greater than zero")
	}

	current, err := u.repo.GetCurrentOpen(tenantID)
	if err != nil {
		return cashregister.CashRegisterEntry{}, ErrNoOpenShift
	}

	entry, err := u.repo.CreateEntry(cashregister.CashRegisterEntry{
		TenantID:        tenantID,
		CashRegisterID:  current.ID,
		Type:            entryType,
		AmountCents:     amountCents,
		Reason:          strings.TrimSpace(reason),
		CreatedByUserID: userID,
	})
	if err != nil {
		return cashregister.CashRegisterEntry{}, err
	}

	// Update system balance
	if entryType == "sangria" {
		current.SystemBalanceCents -= amountCents
		// Record finance expense
		if u.financeRepo != nil {
			_, _ = u.financeRepo.Create(finance.Transaction{
				TenantID:       tenantID,
				Type:           "expense",
				Category:       "sangria",
				AmountCents:    amountCents,
				Description:    "Sangria de Caixa: " + reason,
				CashRegisterID: current.ID,
				PaymentMethod:  "cash",
				Status:         "paid",
			})
		}
	} else {
		current.SystemBalanceCents += amountCents
		// Record finance income
		if u.financeRepo != nil {
			_, _ = u.financeRepo.Create(finance.Transaction{
				TenantID:       tenantID,
				Type:           "income",
				Category:       "suprimento",
				AmountCents:    amountCents,
				Description:    "Suprimento de Caixa: " + reason,
				CashRegisterID: current.ID,
				PaymentMethod:  "cash",
				Status:         "paid",
			})
		}
	}

	_, _ = u.repo.Update(current)
	return entry, nil
}

func (u *UseCase) CloseShift(tenantID string, userID string, closingBalanceCents int, notes string) (cashregister.CashRegister, error) {
	current, err := u.repo.GetCurrentOpen(tenantID)
	if err != nil {
		return cashregister.CashRegister{}, ErrNoOpenShift
	}

	current.ClosedByUserID = userID
	current.Status = "closed"
	current.ClosingBalanceCents = closingBalanceCents
	current.DifferenceCents = closingBalanceCents - current.SystemBalanceCents
	if notes != "" {
		current.Notes = current.Notes + " | Fechamento: " + strings.TrimSpace(notes)
	}

	return u.repo.Update(current)
}

func (u *UseCase) GetCurrentShift(tenantID string) (cashregister.ShiftSummary, error) {
	current, err := u.repo.GetCurrentOpen(tenantID)
	if err != nil {
		return cashregister.ShiftSummary{}, ErrNoOpenShift
	}

	entries, _ := u.repo.ListEntries(tenantID, current.ID)
	var sangrias, suprimentos int
	for _, e := range entries {
		if e.Type == "sangria" {
			sangrias += e.AmountCents
		} else if e.Type == "suprimento" {
			suprimentos += e.AmountCents
		}
	}

	var cashIncome, cardIncome, pixIncome, otherIncome int
	if u.financeRepo != nil {
		txs, _ := u.financeRepo.List(tenantID, "income", "paid", "", "")
		for _, tx := range txs {
			if tx.CashRegisterID == current.ID && tx.Category != "suprimento" {
				switch tx.PaymentMethod {
				case "cash", "dinheiro":
					cashIncome += tx.AmountCents
				case "credit_card", "debit_card", "cartao":
					cardIncome += tx.AmountCents
				case "pix":
					pixIncome += tx.AmountCents
				default:
					otherIncome += tx.AmountCents
				}
			}
		}
	}

	currentDrawer := current.InitialBalanceCents + suprimentos + cashIncome - sangrias
	current.SystemBalanceCents = currentDrawer

	return cashregister.ShiftSummary{
		CashRegister:          current,
		Entries:               entries,
		TotalSangriasCents:    sangrias,
		TotalSuprimentosCents: suprimentos,
		TotalCashIncomeCents:  cashIncome,
		TotalCardIncomeCents:  cardIncome,
		TotalPixIncomeCents:   pixIncome,
		TotalOtherIncomeCents: otherIncome,
		CurrentDrawerCents:    currentDrawer,
	}, nil
}

func (u *UseCase) ListShifts(tenantID string, limit int) ([]cashregister.CashRegister, error) {
	return u.repo.List(tenantID, limit)
}

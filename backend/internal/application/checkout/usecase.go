package checkout

import (
	"context"
	"errors"
	"fmt"
	"strings"

	appointmentdomain "erpagendamentos/backend/internal/domain/appointment"
	auditdomain "erpagendamentos/backend/internal/domain/audit"
	cashregisterdomain "erpagendamentos/backend/internal/domain/cashregister"
	commissiondomain "erpagendamentos/backend/internal/domain/commission"
	financedomain "erpagendamentos/backend/internal/domain/finance"
	professionaldomain "erpagendamentos/backend/internal/domain/professional"
	servicedomain "erpagendamentos/backend/internal/domain/service"
	stockdomain "erpagendamentos/backend/internal/domain/stock"
)

type ProductItem struct {
	ProductID string  `json:"productId"`
	Quantity  float64 `json:"quantity"`
}

type CheckoutInput struct {
	AppointmentID string        `json:"appointmentId"`
	Products      []ProductItem `json:"products,omitempty"`
	DiscountCents int           `json:"discountCents,omitempty"`
	SurchargeCents int          `json:"surchargeCents,omitempty"`
	PaymentMethod string        `json:"paymentMethod"` // "pix", "cash", "credit_card", "debit_card", "other"
	Notes         string        `json:"notes,omitempty"`
}

type CheckoutResult struct {
	Appointment        appointmentdomain.Appointment `json:"appointment"`
	TotalChargedCents  int                           `json:"totalChargedCents"`
	ServiceTotalCents  int                           `json:"serviceTotalCents"`
	ProductsTotalCents int                           `json:"productsTotalCents"`
	CommissionCents    int                           `json:"commissionCents"`
}

type TxRunner interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type AuditLogger interface {
	Log(entry auditdomain.AuditLog) error
}

type OutboxPublisher interface {
	Publish(tenantID string, topic string, payload map[string]any) error
}

type UseCase struct {
	appointmentRepo  appointmentdomain.Repository
	serviceRepo      servicedomain.Repository
	professionalRepo professionaldomain.Repository
	stockRepo        stockdomain.Repository
	commissionRepo   commissiondomain.Repository
	financeRepo      financedomain.Repository
	cashRegisterRepo cashregisterdomain.Repository
	txRunner         TxRunner
	audit            AuditLogger
	outbox           OutboxPublisher
}

func NewUseCase(
	appointmentRepo appointmentdomain.Repository,
	serviceRepo servicedomain.Repository,
	professionalRepo professionaldomain.Repository,
	stockRepo stockdomain.Repository,
	commissionRepo commissiondomain.Repository,
	financeRepo financedomain.Repository,
	cashRegisterRepo cashregisterdomain.Repository,
) *UseCase {
	return &UseCase{
		appointmentRepo:  appointmentRepo,
		serviceRepo:      serviceRepo,
		professionalRepo: professionalRepo,
		stockRepo:        stockRepo,
		commissionRepo:   commissionRepo,
		financeRepo:      financeRepo,
		cashRegisterRepo: cashRegisterRepo,
	}
}

func (u *UseCase) WithTxRunner(runner TxRunner) *UseCase {
	u.txRunner = runner
	return u
}

func (u *UseCase) WithAudit(audit AuditLogger) *UseCase {
	u.audit = audit
	return u
}

func (u *UseCase) WithOutbox(outbox OutboxPublisher) *UseCase {
	u.outbox = outbox
	return u
}

func (u *UseCase) Checkout(tenantID string, input CheckoutInput) (CheckoutResult, error) {
	if tenantID == "" {
		return CheckoutResult{}, errors.New("tenant id is required")
	}
	if input.AppointmentID == "" {
		return CheckoutResult{}, errors.New("appointment id is required")
	}
	if input.PaymentMethod == "" {
		input.PaymentMethod = "cash"
	}

	appt, err := u.appointmentRepo.GetByID(tenantID, input.AppointmentID)
	if err != nil {
		return CheckoutResult{}, fmt.Errorf("appointment not found: %w", err)
	}

	// 1. Calculate services total
	serviceTotalCents := 0
	serviceName := "Serviço"
	if appt.ServiceID != "" {
		if s, err := u.serviceRepo.GetByID(tenantID, appt.ServiceID); err == nil {
			serviceTotalCents += s.PriceCents
			serviceName = s.Name
		}
	}
	for _, sid := range appt.ServiceIDs {
		if sid != appt.ServiceID {
			if s, err := u.serviceRepo.GetByID(tenantID, sid); err == nil {
				serviceTotalCents += s.PriceCents
			}
		}
	}

	// 2. Process products from stock
	productsTotalCents := 0
	for _, item := range input.Products {
		if item.ProductID == "" || item.Quantity <= 0 {
			continue
		}
		prod, err := u.stockRepo.GetProductByID(tenantID, item.ProductID)
		if err != nil {
			continue
		}
		productsTotalCents += int(float64(prod.PriceCents) * item.Quantity)

		// Create stock movement (out / sale)
		_, _ = u.stockRepo.CreateStockMovement(stockdomain.StockMovement{
			TenantID:      tenantID,
			ProductID:     item.ProductID,
			Type:          "out",
			Quantity:      item.Quantity,
			Notes:         "Venda em comanda: Agendamento #" + appt.ID[:8],
			AppointmentID: appt.ID,
		})
	}

	// 3. Final total calculation
	totalCharged := serviceTotalCents + productsTotalCents - input.DiscountCents + input.SurchargeCents
	if totalCharged < 0 {
		totalCharged = 0
	}

	// 4. Check for active cash register shift
	var cashRegisterID string
	if u.cashRegisterRepo != nil {
		if shift, err := u.cashRegisterRepo.GetCurrentOpen(tenantID); err == nil && shift.ID != "" {
			cashRegisterID = shift.ID
		}
	}

	// 5. Record Financial Income
	if u.financeRepo != nil {
		_, _ = u.financeRepo.Create(financedomain.Transaction{
			TenantID:       tenantID,
			Type:           "income",
			Category:       "service_sale",
			AmountCents:    totalCharged,
			Description:    fmt.Sprintf("Fechamento de Comanda: %s", serviceName),
			ReferenceID:    appt.ID,
			CashRegisterID: cashRegisterID,
			PaymentMethod:  input.PaymentMethod,
			Status:         "paid",
		})
	}

	// 6. Automatic Commission Calculation
	totalCommissionCents := 0
	if u.commissionRepo != nil && appt.ProfessionalID != "" {
		prof, _ := u.professionalRepo.GetByID(tenantID, appt.ProfessionalID)
		profName := prof.Name
		if profName == "" {
			profName = "Profissional"
		}

		// Check if a rule matches
		rule, err := u.commissionRepo.FindMatchingRule(tenantID, appt.ProfessionalID, appt.ServiceID)
		var ratePercentage float64
		if err == nil && rule.Percentage > 0 {
			ratePercentage = rule.Percentage
		} else if prof.CommissionPercentage > 0 {
			ratePercentage = prof.CommissionPercentage
		} else {
			ratePercentage = 0
		}

		if ratePercentage > 0 {
			totalCommissionCents = int(float64(serviceTotalCents) * (ratePercentage / 100.0))
			_, _ = u.commissionRepo.CreateRecord(commissiondomain.CommissionRecord{
				TenantID:              tenantID,
				AppointmentID:         appt.ID,
				ProfessionalID:        appt.ProfessionalID,
				ProfessionalName:      profName,
				ServiceID:             appt.ServiceID,
				ServiceName:           serviceName,
				BaseAmountCents:       serviceTotalCents,
				CommissionAmountCents: totalCommissionCents,
				Status:                "pending",
				Notes:                 "Gerado automaticamente pelo fechamento de comanda",
			})
		}
	}

	// 7. Update appointment status
	appt.Status = "completed"
	appt.PaymentStatus = "paid"
	if input.Notes != "" {
		if appt.Notes != "" {
			appt.Notes += " | " + strings.TrimSpace(input.Notes)
		} else {
			appt.Notes = strings.TrimSpace(input.Notes)
		}
	}
	updatedAppt, err := u.appointmentRepo.Update(appt)
	if err != nil {
		return CheckoutResult{}, fmt.Errorf("update appointment: %w", err)
	}

	// 8. Trilha de Auditoria (Audit Log)
	if u.audit != nil {
		_ = u.audit.Log(auditdomain.AuditLog{
			TenantID:   tenantID,
			Action:     "checkout.completed",
			Resource:   "appointment",
			ResourceID: appt.ID,
			Details: map[string]any{
				"totalChargedCents":  totalCharged,
				"serviceTotalCents":  serviceTotalCents,
				"productsTotalCents": productsTotalCents,
				"commissionCents":    totalCommissionCents,
				"paymentMethod":      input.PaymentMethod,
			},
		})
	}

	// 9. Fila de Eventos Transacionais (Outbox & Cross-Instance SSE)
	if u.outbox != nil {
		_ = u.outbox.Publish(tenantID, "checkout.completed", map[string]any{
			"appointmentId":      appt.ID,
			"totalChargedCents":  totalCharged,
			"serviceTotalCents":  serviceTotalCents,
			"productsTotalCents": productsTotalCents,
			"paymentMethod":      input.PaymentMethod,
		})
	}

	return CheckoutResult{
		Appointment:        updatedAppt,
		TotalChargedCents:  totalCharged,
		ServiceTotalCents:  serviceTotalCents,
		ProductsTotalCents: productsTotalCents,
		CommissionCents:    totalCommissionCents,
	}, nil
}


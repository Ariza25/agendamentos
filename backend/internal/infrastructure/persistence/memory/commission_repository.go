package memory

import (
	"errors"
	"sync"
	"time"

	"erpagendamentos/backend/internal/domain/commission"
	"github.com/google/uuid"
)

type CommissionRepository struct {
	mu      sync.RWMutex
	rules   map[string]commission.CommissionRule
	records map[string]commission.CommissionRecord
	payouts map[string]commission.Payout
}

func NewCommissionRepository() *CommissionRepository {
	return &CommissionRepository{
		rules:   make(map[string]commission.CommissionRule),
		records: make(map[string]commission.CommissionRecord),
		payouts: make(map[string]commission.Payout),
	}
}

func (r *CommissionRepository) CreateRule(rule commission.CommissionRule) (commission.CommissionRule, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if rule.ID == "" {
		rule.ID = uuid.NewString()
	}
	now := time.Now().UTC().Format(time.RFC3339)
	rule.CreatedAt = now
	rule.UpdatedAt = now
	r.rules[rule.TenantID+":"+rule.ID] = rule
	return rule, nil
}

func (r *CommissionRepository) UpdateRule(rule commission.CommissionRule) (commission.CommissionRule, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := rule.TenantID + ":" + rule.ID
	existing, found := r.rules[key]
	if !found {
		return commission.CommissionRule{}, errors.New("commission rule not found")
	}
	rule.CreatedAt = existing.CreatedAt
	rule.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	r.rules[key] = rule
	return rule, nil
}

func (r *CommissionRepository) DeleteRule(tenantID string, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.rules, tenantID+":"+id)
	return nil
}

func (r *CommissionRepository) ListRules(tenantID string) ([]commission.CommissionRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []commission.CommissionRule
	for _, rule := range r.rules {
		if rule.TenantID == tenantID {
			result = append(result, rule)
		}
	}
	return result, nil
}

func (r *CommissionRepository) FindMatchingRule(tenantID string, professionalID string, serviceID string) (commission.CommissionRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 1. Specific professional + service
	for _, rule := range r.rules {
		if rule.TenantID == tenantID && rule.ProfessionalID == professionalID && rule.ServiceID == serviceID && professionalID != "" && serviceID != "" {
			return rule, nil
		}
	}
	// 2. Specific professional only
	for _, rule := range r.rules {
		if rule.TenantID == tenantID && rule.ProfessionalID == professionalID && rule.ServiceID == "" && professionalID != "" {
			return rule, nil
		}
	}
	// 3. Specific service only
	for _, rule := range r.rules {
		if rule.TenantID == tenantID && rule.ProfessionalID == "" && rule.ServiceID == serviceID && serviceID != "" {
			return rule, nil
		}
	}
	// 4. Default global rule
	for _, rule := range r.rules {
		if rule.TenantID == tenantID && rule.ProfessionalID == "" && rule.ServiceID == "" {
			return rule, nil
		}
	}

	return commission.CommissionRule{}, errors.New("no matching commission rule found")
}

func (r *CommissionRepository) CreateRecord(record commission.CommissionRecord) (commission.CommissionRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if record.ID == "" {
		record.ID = uuid.NewString()
	}
	if record.Status == "" {
		record.Status = "pending"
	}
	record.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	r.records[record.TenantID+":"+record.ID] = record
	return record, nil
}

func (r *CommissionRepository) ListRecords(tenantID string, professionalID string, status string) ([]commission.CommissionRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []commission.CommissionRecord
	for _, rec := range r.records {
		if rec.TenantID == tenantID {
			if professionalID != "" && rec.ProfessionalID != professionalID {
				continue
			}
			if status != "" && rec.Status != status {
				continue
			}
			result = append(result, rec)
		}
	}
	return result, nil
}

func (r *CommissionRepository) SettleRecords(tenantID string, professionalID string, recordIDs []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC().Format(time.RFC3339)
	targetIDs := make(map[string]bool)
	for _, id := range recordIDs {
		targetIDs[id] = true
	}

	for key, rec := range r.records {
		if rec.TenantID == tenantID {
			if professionalID != "" && rec.ProfessionalID != professionalID {
				continue
			}
			if len(targetIDs) > 0 && !targetIDs[rec.ID] {
				continue
			}
			if rec.Status == "pending" {
				rec.Status = "paid"
				rec.SettledAt = now
				r.records[key] = rec
			}
		}
	}
	return nil
}

func (r *CommissionRepository) GetSummaries(tenantID string) ([]commission.Summary, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	summaries := make(map[string]*commission.Summary)
	for _, rec := range r.records {
		if rec.TenantID == tenantID {
			s, exists := summaries[rec.ProfessionalID]
			if !exists {
				s = &commission.Summary{
					ProfessionalID:   rec.ProfessionalID,
					ProfessionalName: rec.ProfessionalName,
				}
				summaries[rec.ProfessionalID] = s
			}
			if rec.Status == "pending" {
				s.PendingCents += rec.CommissionAmountCents
			} else if rec.Status == "paid" {
				s.PaidCents += rec.CommissionAmountCents
			}
			s.TotalCents += rec.CommissionAmountCents
		}
	}

	var result []commission.Summary
	for _, s := range summaries {
		result = append(result, *s)
	}
	return result, nil
}

func (r *CommissionRepository) CreatePayout(payout commission.Payout, recordIDs []string) (commission.Payout, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if payout.ID == "" {
		payout.ID = uuid.NewString()
	}
	now := time.Now().UTC().Format(time.RFC3339)
	payout.PaidAt = now
	payout.CreatedAt = now

	targetIDs := make(map[string]bool)
	for _, id := range recordIDs {
		targetIDs[id] = true
	}

	for key, rec := range r.records {
		if rec.TenantID == payout.TenantID && rec.ProfessionalID == payout.ProfessionalID {
			if len(targetIDs) > 0 && !targetIDs[rec.ID] {
				continue
			}
			if rec.Status == "pending" {
				rec.Status = "paid"
				rec.SettledAt = now
				rec.PayoutID = payout.ID
				r.records[key] = rec
			}
		}
	}

	r.payouts[payout.TenantID+":"+payout.ID] = payout
	return payout, nil
}

func (r *CommissionRepository) ListPayouts(tenantID string, professionalID string) ([]commission.Payout, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []commission.Payout
	for _, p := range r.payouts {
		if p.TenantID == tenantID {
			if professionalID != "" && p.ProfessionalID != professionalID {
				continue
			}
			result = append(result, p)
		}
	}
	return result, nil
}

package commerce

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	domain "erpagendamentos/backend/internal/domain/commerce"
)

var (
	ErrPackageNotFound       = errors.New("package not found")
	ErrInvalidPackage        = errors.New("invalid package")
	ErrGiftCardNotFound      = errors.New("gift card not found")
	ErrInvalidGiftCard       = errors.New("invalid gift card")
	ErrMembershipNotFound    = errors.New("membership not found")
	ErrInvalidMembership     = errors.New("invalid membership")
	ErrMemberPlanNotFound    = errors.New("member plan not found")
	ErrInvalidMemberPlan     = errors.New("invalid member plan")
	ErrPackageCreditNotFound = errors.New("package credit not found")
)

type UseCase struct {
	packages       domain.PackageRepository
	packageCredits domain.PackageCreditRepository
	giftCards      domain.GiftCardRepository
	memberships    domain.MembershipRepository
	memberPlans    domain.MemberPlanRepository
}

func NewUseCase(
	packages domain.PackageRepository,
	packageCredits domain.PackageCreditRepository,
	giftCards domain.GiftCardRepository,
	memberships domain.MembershipRepository,
	memberPlans domain.MemberPlanRepository,
) *UseCase {
	return &UseCase{
		packages:       packages,
		packageCredits: packageCredits,
		giftCards:      giftCards,
		memberships:    memberships,
		memberPlans:    memberPlans,
	}
}

// Packages

func (u *UseCase) ListPackages(tenantID string) ([]domain.Package, error) {
	return u.packages.ListByTenant(strings.TrimSpace(tenantID))
}

func (u *UseCase) GetPackage(tenantID string, id string) (domain.Package, error) {
	item, err := u.packages.GetByID(strings.TrimSpace(tenantID), strings.TrimSpace(id))
	if err != nil {
		return domain.Package{}, ErrPackageNotFound
	}
	return item, nil
}

func (u *UseCase) CreatePackage(input domain.Package) (domain.Package, error) {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.Name = strings.TrimSpace(input.Name)
	if input.TenantID == "" || input.Name == "" || input.PriceCents <= 0 || input.SessionsCount <= 0 {
		return domain.Package{}, ErrInvalidPackage
	}
	return u.packages.Create(input)
}

func (u *UseCase) UpdatePackage(input domain.Package) (domain.Package, error) {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.ID = strings.TrimSpace(input.ID)
	if input.TenantID == "" || input.ID == "" {
		return domain.Package{}, ErrInvalidPackage
	}
	if strings.TrimSpace(input.Name) == "" {
		return domain.Package{}, ErrInvalidPackage
	}
	item, err := u.packages.Update(input)
	if err != nil {
		return domain.Package{}, ErrPackageNotFound
	}
	return item, nil
}

func (u *UseCase) SetPackageActive(tenantID string, id string, active bool) (domain.Package, error) {
	item, err := u.packages.SetActive(strings.TrimSpace(tenantID), strings.TrimSpace(id), active)
	if err != nil {
		return domain.Package{}, ErrPackageNotFound
	}
	return item, nil
}

func (u *UseCase) DeletePackage(tenantID string, id string) error {
	if err := u.packages.Delete(strings.TrimSpace(tenantID), strings.TrimSpace(id)); err != nil {
		return ErrPackageNotFound
	}
	return nil
}

// Package credits

func (u *UseCase) ListPackageCredits(tenantID string) ([]domain.PackageCredit, error) {
	return u.packageCredits.ListByTenant(strings.TrimSpace(tenantID))
}

func (u *UseCase) ListPackageCreditsByCustomer(customerID string) ([]domain.PackageCredit, error) {
	return u.packageCredits.ListByCustomer(strings.TrimSpace(customerID))
}

func (u *UseCase) CreatePackageCredit(input domain.PackageCredit) (domain.PackageCredit, error) {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.CustomerID = strings.TrimSpace(input.CustomerID)
	if input.TenantID == "" || input.CustomerID == "" || input.SessionsTotal <= 0 {
		return domain.PackageCredit{}, ErrInvalidPackage
	}
	return u.packageCredits.Create(input)
}

func (u *UseCase) IncrementPackageCredit(tenantID string, id string) (domain.PackageCredit, error) {
	item, err := u.packageCredits.IncrementUsed(strings.TrimSpace(tenantID), strings.TrimSpace(id))
	if err != nil {
		return domain.PackageCredit{}, ErrPackageCreditNotFound
	}
	return item, nil
}

// Gift cards

func (u *UseCase) ListGiftCards(tenantID string) ([]domain.GiftCard, error) {
	return u.giftCards.ListByTenant(strings.TrimSpace(tenantID))
}

func (u *UseCase) GetGiftCard(tenantID string, id string) (domain.GiftCard, error) {
	item, err := u.giftCards.GetByID(strings.TrimSpace(tenantID), strings.TrimSpace(id))
	if err != nil {
		return domain.GiftCard{}, ErrGiftCardNotFound
	}
	return item, nil
}

func (u *UseCase) CreateGiftCard(input domain.GiftCard) (domain.GiftCard, error) {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.Code = strings.ToUpper(strings.TrimSpace(input.Code))
	if input.TenantID == "" || input.ValueCents <= 0 {
		return domain.GiftCard{}, ErrInvalidGiftCard
	}
	if input.Code == "" {
		input.Code = generateGiftCardCode()
	}
	if input.BalanceCents == 0 {
		input.BalanceCents = input.ValueCents
	}
	if input.Status == "" {
		input.Status = "active"
	}
	return u.giftCards.Create(input)
}

func (u *UseCase) UpdateGiftCard(input domain.GiftCard) (domain.GiftCard, error) {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.ID = strings.TrimSpace(input.ID)
	if input.TenantID == "" || input.ID == "" {
		return domain.GiftCard{}, ErrInvalidGiftCard
	}
	item, err := u.giftCards.Update(input)
	if err != nil {
		return domain.GiftCard{}, ErrGiftCardNotFound
	}
	return item, nil
}

func (u *UseCase) DeleteGiftCard(tenantID string, id string) error {
	if err := u.giftCards.Delete(strings.TrimSpace(tenantID), strings.TrimSpace(id)); err != nil {
		return ErrGiftCardNotFound
	}
	return nil
}

// Memberships

func (u *UseCase) ListMemberships(tenantID string) ([]domain.Membership, error) {
	return u.memberships.ListByTenant(strings.TrimSpace(tenantID))
}

func (u *UseCase) GetMembership(tenantID string, id string) (domain.Membership, error) {
	item, err := u.memberships.GetByID(strings.TrimSpace(tenantID), strings.TrimSpace(id))
	if err != nil {
		return domain.Membership{}, ErrMembershipNotFound
	}
	return item, nil
}

func (u *UseCase) CreateMembership(input domain.Membership) (domain.Membership, error) {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.Name = strings.TrimSpace(input.Name)
	if input.TenantID == "" || input.Name == "" || input.PriceCents <= 0 {
		return domain.Membership{}, ErrInvalidMembership
	}
	if input.BillingPeriodDays <= 0 {
		input.BillingPeriodDays = 30
	}
	return u.memberships.Create(input)
}

func (u *UseCase) UpdateMembership(input domain.Membership) (domain.Membership, error) {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.ID = strings.TrimSpace(input.ID)
	if input.TenantID == "" || input.ID == "" {
		return domain.Membership{}, ErrInvalidMembership
	}
	item, err := u.memberships.Update(input)
	if err != nil {
		return domain.Membership{}, ErrMembershipNotFound
	}
	return item, nil
}

func (u *UseCase) SetMembershipActive(tenantID string, id string, active bool) (domain.Membership, error) {
	item, err := u.memberships.SetActive(strings.TrimSpace(tenantID), strings.TrimSpace(id), active)
	if err != nil {
		return domain.Membership{}, ErrMembershipNotFound
	}
	return item, nil
}

func (u *UseCase) DeleteMembership(tenantID string, id string) error {
	if err := u.memberships.Delete(strings.TrimSpace(tenantID), strings.TrimSpace(id)); err != nil {
		return ErrMembershipNotFound
	}
	return nil
}

// Member plans

func (u *UseCase) ListMemberPlans(tenantID string) ([]domain.MemberPlan, error) {
	return u.memberPlans.ListByTenant(strings.TrimSpace(tenantID))
}

func (u *UseCase) ListMemberPlansByCustomer(customerID string) ([]domain.MemberPlan, error) {
	return u.memberPlans.ListByCustomer(strings.TrimSpace(customerID))
}

func (u *UseCase) CreateMemberPlan(input domain.MemberPlan) (domain.MemberPlan, error) {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.CustomerID = strings.TrimSpace(input.CustomerID)
	if input.TenantID == "" || input.CustomerID == "" || input.MembershipID == "" {
		return domain.MemberPlan{}, ErrInvalidMemberPlan
	}
	if input.Status == "" {
		input.Status = "active"
	}
	return u.memberPlans.Create(input)
}

func (u *UseCase) UpdateMemberPlanStatus(tenantID string, id string, status string) (domain.MemberPlan, error) {
	item, err := u.memberPlans.UpdateStatus(strings.TrimSpace(tenantID), strings.TrimSpace(id), strings.TrimSpace(status))
	if err != nil {
		return domain.MemberPlan{}, ErrMemberPlanNotFound
	}
	return item, nil
}

func generateGiftCardCode() string {
	return strings.ToUpper(strings.ReplaceAll(uuid.NewString()[:8], "-", ""))
}

func (u *UseCase) Now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

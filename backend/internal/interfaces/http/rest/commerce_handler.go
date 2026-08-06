package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	appcommerce "erpagendamentos/backend/internal/application/commerce"
	commercedomain "erpagendamentos/backend/internal/domain/commerce"
)

type CommerceHandler struct {
	useCase *appcommerce.UseCase
}

func NewCommerceHandler(useCase *appcommerce.UseCase) *CommerceHandler {
	return &CommerceHandler{useCase: useCase}
}

// Packages

func (h *CommerceHandler) ListPackages(w http.ResponseWriter, r *http.Request) {
	items, err := h.useCase.ListPackages(resolveRequestTenantID(r, nil))
	if err != nil {
		writeInternalError(w, "Could not load the packages.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"packages": items})
}

func (h *CommerceHandler) CreatePackage(w http.ResponseWriter, r *http.Request) {
	var input commercedomain.Package
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the package data.")
		return
	}
	input.TenantID = resolveRequestTenantID(r, nil)
	item, err := h.useCase.CreatePackage(input)
	if err != nil {
		writeCommerceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *CommerceHandler) UpdatePackage(w http.ResponseWriter, r *http.Request) {
	var input commercedomain.Package
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the package data.")
		return
	}
	input.TenantID = resolveRequestTenantID(r, nil)
	input.ID = r.PathValue("packageID")
	item, err := h.useCase.UpdatePackage(input)
	if err != nil {
		writeCommerceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *CommerceHandler) SetPackageActive(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Active bool `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Provide the package active state.")
		return
	}
	item, err := h.useCase.SetPackageActive(resolveRequestTenantID(r, nil), r.PathValue("packageID"), input.Active)
	if err != nil {
		writeCommerceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *CommerceHandler) DeletePackage(w http.ResponseWriter, r *http.Request) {
	if err := h.useCase.DeletePackage(resolveRequestTenantID(r, nil), r.PathValue("packageID")); err != nil {
		writeCommerceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Package credits

func (h *CommerceHandler) ListPackageCredits(w http.ResponseWriter, r *http.Request) {
	items, err := h.useCase.ListPackageCredits(resolveRequestTenantID(r, nil))
	if err != nil {
		writeInternalError(w, "Could not load the package credits.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"packageCredits": items})
}

func (h *CommerceHandler) CreatePackageCredit(w http.ResponseWriter, r *http.Request) {
	var input commercedomain.PackageCredit
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the package credit data.")
		return
	}
	input.TenantID = resolveRequestTenantID(r, nil)
	item, err := h.useCase.CreatePackageCredit(input)
	if err != nil {
		writeCommerceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

// Gift cards

func (h *CommerceHandler) ListGiftCards(w http.ResponseWriter, r *http.Request) {
	items, err := h.useCase.ListGiftCards(resolveRequestTenantID(r, nil))
	if err != nil {
		writeInternalError(w, "Could not load the gift cards.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"giftCards": items})
}

func (h *CommerceHandler) CreateGiftCard(w http.ResponseWriter, r *http.Request) {
	var input commercedomain.GiftCard
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the gift card data.")
		return
	}
	input.TenantID = resolveRequestTenantID(r, nil)
	item, err := h.useCase.CreateGiftCard(input)
	if err != nil {
		writeCommerceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *CommerceHandler) UpdateGiftCard(w http.ResponseWriter, r *http.Request) {
	var input commercedomain.GiftCard
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the gift card data.")
		return
	}
	input.TenantID = resolveRequestTenantID(r, nil)
	input.ID = r.PathValue("giftCardID")
	item, err := h.useCase.UpdateGiftCard(input)
	if err != nil {
		writeCommerceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *CommerceHandler) DeleteGiftCard(w http.ResponseWriter, r *http.Request) {
	if err := h.useCase.DeleteGiftCard(resolveRequestTenantID(r, nil), r.PathValue("giftCardID")); err != nil {
		writeCommerceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Memberships

func (h *CommerceHandler) ListMemberships(w http.ResponseWriter, r *http.Request) {
	items, err := h.useCase.ListMemberships(resolveRequestTenantID(r, nil))
	if err != nil {
		writeInternalError(w, "Could not load the memberships.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"memberships": items})
}

func (h *CommerceHandler) CreateMembership(w http.ResponseWriter, r *http.Request) {
	var input commercedomain.Membership
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the membership data.")
		return
	}
	input.TenantID = resolveRequestTenantID(r, nil)
	item, err := h.useCase.CreateMembership(input)
	if err != nil {
		writeCommerceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *CommerceHandler) UpdateMembership(w http.ResponseWriter, r *http.Request) {
	var input commercedomain.Membership
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the membership data.")
		return
	}
	input.TenantID = resolveRequestTenantID(r, nil)
	input.ID = r.PathValue("membershipID")
	item, err := h.useCase.UpdateMembership(input)
	if err != nil {
		writeCommerceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *CommerceHandler) SetMembershipActive(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Active bool `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Provide the membership active state.")
		return
	}
	item, err := h.useCase.SetMembershipActive(resolveRequestTenantID(r, nil), r.PathValue("membershipID"), input.Active)
	if err != nil {
		writeCommerceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *CommerceHandler) DeleteMembership(w http.ResponseWriter, r *http.Request) {
	if err := h.useCase.DeleteMembership(resolveRequestTenantID(r, nil), r.PathValue("membershipID")); err != nil {
		writeCommerceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Member plans

func (h *CommerceHandler) ListMemberPlans(w http.ResponseWriter, r *http.Request) {
	items, err := h.useCase.ListMemberPlans(resolveRequestTenantID(r, nil))
	if err != nil {
		writeInternalError(w, "Could not load the member plans.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"memberPlans": items})
}

func (h *CommerceHandler) CreateMemberPlan(w http.ResponseWriter, r *http.Request) {
	var input commercedomain.MemberPlan
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the member plan data.")
		return
	}
	input.TenantID = resolveRequestTenantID(r, nil)
	item, err := h.useCase.CreateMemberPlan(input)
	if err != nil {
		writeCommerceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *CommerceHandler) UpdateMemberPlanStatus(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Provide the new member plan status.")
		return
	}
	item, err := h.useCase.UpdateMemberPlanStatus(resolveRequestTenantID(r, nil), r.PathValue("planID"), input.Status)
	if err != nil {
		writeCommerceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func writeCommerceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, appcommerce.ErrPackageNotFound) ||
		errors.Is(err, appcommerce.ErrGiftCardNotFound) ||
		errors.Is(err, appcommerce.ErrMembershipNotFound) ||
		errors.Is(err, appcommerce.ErrMemberPlanNotFound) ||
		errors.Is(err, appcommerce.ErrPackageCreditNotFound):
		writeNotFoundError(w, "not_found", "Item not found.")
	case errors.Is(err, appcommerce.ErrInvalidPackage) ||
		errors.Is(err, appcommerce.ErrInvalidGiftCard) ||
		errors.Is(err, appcommerce.ErrInvalidMembership) ||
		errors.Is(err, appcommerce.ErrInvalidMemberPlan):
		writeBadRequestError(w, "invalid_input", "Invalid item data.")
	default:
		writeInternalError(w, "Could not process the request.")
	}
}

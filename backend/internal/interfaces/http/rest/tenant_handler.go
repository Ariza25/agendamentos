package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	apptenant "erpagendamentos/backend/internal/application/tenant"
	tenantdomain "erpagendamentos/backend/internal/domain/tenant"
)

type TenantHandler struct {
	useCase *apptenant.UseCase
}

func NewTenantHandler(useCase *apptenant.UseCase) *TenantHandler {
	return &TenantHandler{useCase: useCase}
}

func (h *TenantHandler) List(w http.ResponseWriter, r *http.Request) {
	tenants, err := h.useCase.ListTenants()
	if err != nil {
		writeInternalError(w, "Could not load the establishments.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"tenants": tenants})
}

func (h *TenantHandler) Get(w http.ResponseWriter, r *http.Request) {
	tenant, err := h.useCase.GetTenant(r.PathValue("tenantID"))
	if err != nil {
		writeNotFoundError(w, "not_found", "Establishment not found.")
		return
	}
	writeJSON(w, http.StatusOK, tenant)
}

func (h *TenantHandler) GetCurrent(w http.ResponseWriter, r *http.Request) {
	tenant, err := h.useCase.GetTenant(resolveRequestTenantID(r, nil))
	if err != nil {
		writeNotFoundError(w, "not_found", "Establishment not found.")
		return
	}
	writeJSON(w, http.StatusOK, tenant)
}

func (h *TenantHandler) UpdateCurrent(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	var input tenantdomain.Tenant
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the establishment data.")
		return
	}
	input.ID = tenantID

	current, err := h.useCase.GetTenant(tenantID)
	if err != nil {
		writeNotFoundError(w, "not_found", "Establishment not found.")
		return
	}

	input.Plan = current.Plan
	input.Active = current.Active
	input.TrialEndsAt = current.TrialEndsAt
	input.ActiveModules = current.ActiveModules
	if input.Segment == "" {
		input.Segment = current.Segment
	}

	tenant, err := h.useCase.UpdateTenant(input)
	if err != nil {
		writeTenantError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tenant)
}

func (h *TenantHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input tenantdomain.Tenant
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the establishment data.")
		return
	}
	tenant, err := h.useCase.CreateTenant(input)
	if err != nil {
		writeTenantError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, tenant)
}

func (h *TenantHandler) Update(w http.ResponseWriter, r *http.Request) {
	var input tenantdomain.Tenant
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the establishment data.")
		return
	}
	input.ID = r.PathValue("tenantID")
	tenant, err := h.useCase.UpdateTenant(input)
	if err != nil {
		writeTenantError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tenant)
}

func (h *TenantHandler) SetActive(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Active bool `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Provide whether the establishment should be active.")
		return
	}
	tenant, err := h.useCase.SetActive(r.PathValue("tenantID"), input.Active)
	if err != nil {
		writeTenantError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tenant)
}

func (h *TenantHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	var input tenantdomain.UpdateSubscriptionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read subscription update data.")
		return
	}
	tenantID := r.PathValue("tenantID")
	updated, err := h.useCase.UpdateSubscription(tenantID, input)
	if err != nil {
		writeTenantError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *TenantHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.useCase.DeleteTenant(r.PathValue("tenantID")); err != nil {
		writeTenantError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeTenantError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apptenant.ErrTenantNotFound):
		writeNotFoundError(w, "not_found", "Establishment not found.")
	case errors.Is(err, apptenant.ErrInvalidTenant):
		writeBadRequestError(w, "invalid_input", "Invalid establishment data.")
	case errors.Is(err, apptenant.ErrSlugAlreadyUsed):
		writeConflictError(w, "slug_already_used", "An establishment with this identifier already exists.")
	case errors.Is(err, apptenant.ErrLastActiveTenant):
		writeConflictError(w, "last_active_tenant", "It is not possible to remove the last active establishment.")
	default:
		writeInternalError(w, "Could not process the establishment.")
	}
}

func trimmedTenantSlug(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

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
		writeInternalError(w, "Nao foi possivel carregar os estabelecimentos.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"tenants": tenants})
}

func (h *TenantHandler) Get(w http.ResponseWriter, r *http.Request) {
	tenant, err := h.useCase.GetTenant(r.PathValue("tenantID"))
	if err != nil {
		writeNotFoundError(w, "not_found", "Estabelecimento nao encontrado.")
		return
	}
	writeJSON(w, http.StatusOK, tenant)
}

func (h *TenantHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input tenantdomain.Tenant
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Nao foi possivel ler os dados do estabelecimento.")
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
		writeInvalidPayload(w, "Nao foi possivel ler os dados do estabelecimento.")
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
		writeInvalidPayload(w, "Informe se o estabelecimento deve estar ativo.")
		return
	}
	tenant, err := h.useCase.SetActive(r.PathValue("tenantID"), input.Active)
	if err != nil {
		writeTenantError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tenant)
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
		writeNotFoundError(w, "not_found", "Estabelecimento nao encontrado.")
	case errors.Is(err, apptenant.ErrInvalidTenant):
		writeBadRequestError(w, "invalid_input", "Dados do estabelecimento invalidos.")
	case errors.Is(err, apptenant.ErrSlugAlreadyUsed):
		writeConflictError(w, "slug_already_used", "Ja existe um estabelecimento com este identificador.")
	case errors.Is(err, apptenant.ErrLastActiveTenant):
		writeConflictError(w, "last_active_tenant", "Nao e possivel remover o ultimo estabelecimento ativo.")
	default:
		writeInternalError(w, "Nao foi possivel processar o estabelecimento.")
	}
}

func trimmedTenantSlug(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	appresource "erpagendamentos/backend/internal/application/resource"
	resourcedomain "erpagendamentos/backend/internal/domain/resource"
)

type ResourceHandler struct {
	useCase *appresource.UseCase
}

func NewResourceHandler(useCase *appresource.UseCase) *ResourceHandler {
	return &ResourceHandler{useCase: useCase}
}

func (h *ResourceHandler) List(w http.ResponseWriter, r *http.Request) {
	resources, err := h.useCase.ListByTenant(resolveRequestTenantID(r, nil))
	if err != nil {
		writeInternalError(w, "Could not load the resources.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"resources": resources})
}

func (h *ResourceHandler) Get(w http.ResponseWriter, r *http.Request) {
	item, err := h.useCase.GetByID(resolveRequestTenantID(r, nil), r.PathValue("resourceID"))
	if err != nil {
		writeNotFoundError(w, "not_found", "Resource not found.")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *ResourceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input resourcedomain.Resource
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the resource data.")
		return
	}
	input.TenantID = resolveRequestTenantID(r, nil)
	item, err := h.useCase.Create(input)
	if err != nil {
		writeResourceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *ResourceHandler) Update(w http.ResponseWriter, r *http.Request) {
	var input resourcedomain.Resource
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the resource data.")
		return
	}
	input.TenantID = resolveRequestTenantID(r, nil)
	input.ID = r.PathValue("resourceID")
	item, err := h.useCase.Update(input)
	if err != nil {
		writeResourceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *ResourceHandler) SetActive(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Active bool `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Provide whether the resource should be active.")
		return
	}
	item, err := h.useCase.SetActive(resolveRequestTenantID(r, nil), r.PathValue("resourceID"), input.Active)
	if err != nil {
		writeResourceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *ResourceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.useCase.Delete(resolveRequestTenantID(r, nil), r.PathValue("resourceID")); err != nil {
		writeResourceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeResourceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, appresource.ErrResourceNotFound):
		writeNotFoundError(w, "not_found", "Resource not found.")
	case errors.Is(err, appresource.ErrInvalidResource):
		writeBadRequestError(w, "invalid_input", "Invalid resource data.")
	default:
		writeInternalError(w, "Could not process the resource.")
	}
}

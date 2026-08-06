package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	appservice "erpagendamentos/backend/internal/application/service"
	servicedomain "erpagendamentos/backend/internal/domain/service"
)

type ServiceHandler struct {
	useCase *appservice.UseCase
}

func NewServiceHandler(useCase *appservice.UseCase) *ServiceHandler {
	return &ServiceHandler{useCase: useCase}
}

func (h *ServiceHandler) List(w http.ResponseWriter, r *http.Request) {
	services, err := h.useCase.ListByTenant(resolveRequestTenantID(r, nil))
	if err != nil {
		writeInternalError(w, "Could not load the services.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"services": services})
}

func (h *ServiceHandler) Get(w http.ResponseWriter, r *http.Request) {
	service, err := h.useCase.GetByID(resolveRequestTenantID(r, nil), r.PathValue("serviceID"))
	if err != nil {
		writeNotFoundError(w, "not_found", "Service not found.")
		return
	}
	writeJSON(w, http.StatusOK, service)
}

func (h *ServiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input servicedomain.Service
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the service data.")
		return
	}
	input.TenantID = resolveRequestTenantID(r, nil)
	service, err := h.useCase.Create(input)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, service)
}

func (h *ServiceHandler) Update(w http.ResponseWriter, r *http.Request) {
	var input servicedomain.Service
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the service data.")
		return
	}
	input.TenantID = resolveRequestTenantID(r, nil)
	input.ID = r.PathValue("serviceID")
	service, err := h.useCase.Update(input)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, service)
}

func (h *ServiceHandler) SetActive(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Active bool `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Provide whether the service should be active.")
		return
	}
	service, err := h.useCase.SetActive(resolveRequestTenantID(r, nil), r.PathValue("serviceID"), input.Active)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, service)
}

func (h *ServiceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.useCase.Delete(resolveRequestTenantID(r, nil), r.PathValue("serviceID")); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, appservice.ErrServiceNotFound):
		writeNotFoundError(w, "not_found", "Service not found.")
	case errors.Is(err, appservice.ErrInvalidService):
		writeBadRequestError(w, "invalid_input", "Invalid service data.")
	default:
		writeInternalError(w, "Could not process the service.")
	}
}

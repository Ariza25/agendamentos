package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	appprofessional "erpagendamentos/backend/internal/application/professional"
	professionaldomain "erpagendamentos/backend/internal/domain/professional"
)

type ProfessionalHandler struct {
	useCase *appprofessional.UseCase
}

func NewProfessionalHandler(useCase *appprofessional.UseCase) *ProfessionalHandler {
	return &ProfessionalHandler{useCase: useCase}
}

func (h *ProfessionalHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	professionals, err := h.useCase.ListByTenant(tenantID)
	if err != nil {
		writeInternalError(w, "Could not load the professionals.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"professionals": professionals})
}

func (h *ProfessionalHandler) Get(w http.ResponseWriter, r *http.Request) {
	professional, err := h.useCase.GetByID(resolveRequestTenantID(r, nil), r.PathValue("professionalID"))
	if err != nil {
		writeNotFoundError(w, "not_found", "Professional not found.")
		return
	}
	writeJSON(w, http.StatusOK, professional)
}

func (h *ProfessionalHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input professionaldomain.Professional
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the professional data.")
		return
	}
	input.TenantID = resolveRequestTenantID(r, nil)
	professional, err := h.useCase.Create(input)
	if err != nil {
		writeProfessionalError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, professional)
}

func (h *ProfessionalHandler) Update(w http.ResponseWriter, r *http.Request) {
	var input professionaldomain.Professional
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the professional data.")
		return
	}
	input.TenantID = resolveRequestTenantID(r, nil)
	input.ID = r.PathValue("professionalID")
	professional, err := h.useCase.Update(input)
	if err != nil {
		writeProfessionalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, professional)
}

func (h *ProfessionalHandler) SetActive(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Active bool `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Provide whether the professional should be active.")
		return
	}
	professional, err := h.useCase.SetActive(resolveRequestTenantID(r, nil), r.PathValue("professionalID"), input.Active)
	if err != nil {
		writeProfessionalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, professional)
}

func (h *ProfessionalHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.useCase.Delete(resolveRequestTenantID(r, nil), r.PathValue("professionalID")); err != nil {
		writeProfessionalError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeProfessionalError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, appprofessional.ErrProfessionalNotFound):
		writeNotFoundError(w, "not_found", "Professional not found.")
	case errors.Is(err, appprofessional.ErrInvalidProfessional):
		writeBadRequestError(w, "invalid_input", "Invalid professional data.")
	default:
		writeInternalError(w, "Could not process the professional.")
	}
}

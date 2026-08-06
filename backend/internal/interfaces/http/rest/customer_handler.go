package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	appcustomer "erpagendamentos/backend/internal/application/customer"
	customerdomain "erpagendamentos/backend/internal/domain/customer"
)

type CustomerHandler struct {
	useCase *appcustomer.UseCase
}

func NewCustomerHandler(useCase *appcustomer.UseCase) *CustomerHandler {
	return &CustomerHandler{useCase: useCase}
}

func (h *CustomerHandler) List(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	customers, err := h.useCase.ListByTenant(resolveRequestTenantID(r, nil), search)
	if err != nil {
		writeInternalError(w, "Could not load the customers.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"customers": customers})
}

func (h *CustomerHandler) Get(w http.ResponseWriter, r *http.Request) {
	customer, err := h.useCase.GetByID(resolveRequestTenantID(r, nil), r.PathValue("customerID"))
	if err != nil {
		writeNotFoundError(w, "not_found", "Customer not found.")
		return
	}
	writeJSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input customerdomain.Customer
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the customer data.")
		return
	}
	input.TenantID = resolveRequestTenantID(r, nil)
	customer, err := h.useCase.Create(input)
	if err != nil {
		writeCustomerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, customer)
}

func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	var input customerdomain.Customer
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the customer data.")
		return
	}
	input.TenantID = resolveRequestTenantID(r, nil)
	input.ID = r.PathValue("customerID")
	customer, err := h.useCase.Update(input)
	if err != nil {
		writeCustomerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.useCase.Delete(resolveRequestTenantID(r, nil), r.PathValue("customerID")); err != nil {
		writeCustomerError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeCustomerError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, appcustomer.ErrCustomerNotFound):
		writeNotFoundError(w, "not_found", "Customer not found.")
	case errors.Is(err, appcustomer.ErrInvalidCustomer):
		writeBadRequestError(w, "invalid_input", "Invalid customer data.")
	default:
		writeInternalError(w, "Could not process the customer.")
	}
}

package rest

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"erpagendamentos/backend/internal/application/stock"
	stockdomain "erpagendamentos/backend/internal/domain/stock"
)

type StockHandler struct {
	stockUseCase *stock.UseCase
}

func NewStockHandler(stockUseCase *stock.UseCase) *StockHandler {
	return &StockHandler{stockUseCase: stockUseCase}
}

func (h *StockHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	activeOnly := r.URL.Query().Get("activeOnly") == "true"

	products, err := h.stockUseCase.ListProducts(tenantID, activeOnly)
	if err != nil {
		writeInternalError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"products": products})
}

func (h *StockHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	var input stockdomain.Product
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "invalid request body")
		return
	}
	input.TenantID = tenantID

	created, err := h.stockUseCase.CreateProduct(input)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_product", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *StockHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	id := r.PathValue("productID")
	if id == "" {
		id = strings.TrimPrefix(r.URL.Path, "/api/v1/admin/products/")
	}
	id = strings.TrimSpace(id)
	if id == "" {
		writeAPIError(w, http.StatusBadRequest, "invalid_id", "product id is required")
		return
	}

	var input stockdomain.Product
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "invalid request body")
		return
	}
	input.TenantID = tenantID
	input.ID = id

	updated, err := h.stockUseCase.UpdateProduct(input)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "update_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *StockHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	id := r.PathValue("productID")
	if id == "" {
		id = strings.TrimPrefix(r.URL.Path, "/api/v1/admin/products/")
	}
	id = strings.TrimSpace(id)

	if err := h.stockUseCase.DeleteProduct(tenantID, id); err != nil {
		writeInternalError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *StockHandler) CreateMovement(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	var input stockdomain.StockMovement
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "invalid request body")
		return
	}
	input.TenantID = tenantID

	created, err := h.stockUseCase.CreateMovement(input)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_movement", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *StockHandler) ListMovements(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	productID := r.URL.Query().Get("productId")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	movements, err := h.stockUseCase.ListMovements(tenantID, productID, limit)
	if err != nil {
		writeInternalError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"movements": movements})
}

func (h *StockHandler) SetServiceSupplies(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	serviceID := r.PathValue("serviceID")
	if serviceID == "" {
		serviceID = strings.TrimPrefix(r.URL.Path, "/api/v1/admin/services/")
		serviceID = strings.TrimSuffix(serviceID, "/supplies")
	}

	var supplies []stockdomain.ServiceSupply
	if err := json.NewDecoder(r.Body).Decode(&supplies); err != nil {
		writeInvalidPayload(w, "invalid request body")
		return
	}

	if err := h.stockUseCase.SetServiceSupplies(tenantID, serviceID, supplies); err != nil {
		writeAPIError(w, http.StatusBadRequest, "set_supplies_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *StockHandler) ListServiceSupplies(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	serviceID := r.PathValue("serviceID")
	if serviceID == "" {
		serviceID = strings.TrimPrefix(r.URL.Path, "/api/v1/admin/services/")
		serviceID = strings.TrimSuffix(serviceID, "/supplies")
	}

	supplies, err := h.stockUseCase.ListServiceSupplies(tenantID, serviceID)
	if err != nil {
		writeInternalError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"supplies": supplies})
}

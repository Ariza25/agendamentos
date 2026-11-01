package stock

import (
	"errors"
	"strings"

	"erpagendamentos/backend/internal/domain/stock"
)

type UseCase struct {
	repo stock.Repository
}

func NewUseCase(repo stock.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) CreateProduct(input stock.Product) (stock.Product, error) {
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		return stock.Product{}, errors.New("product name is required")
	}
	if input.TenantID == "" {
		return stock.Product{}, errors.New("tenant id is required")
	}
	if input.Unit == "" {
		input.Unit = "un"
	}
	return u.repo.CreateProduct(input)
}

func (u *UseCase) UpdateProduct(input stock.Product) (stock.Product, error) {
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		return stock.Product{}, errors.New("product name is required")
	}
	return u.repo.UpdateProduct(input)
}

func (u *UseCase) GetProduct(tenantID string, id string) (stock.Product, error) {
	return u.repo.GetProductByID(tenantID, id)
}

func (u *UseCase) ListProducts(tenantID string, activeOnly bool) ([]stock.Product, error) {
	return u.repo.ListProducts(tenantID, activeOnly)
}

func (u *UseCase) DeleteProduct(tenantID string, id string) error {
	return u.repo.DeleteProduct(tenantID, id)
}

func (u *UseCase) CreateMovement(input stock.StockMovement) (stock.StockMovement, error) {
	if input.ProductID == "" || input.Quantity == 0 {
		return stock.StockMovement{}, errors.New("product id and non-zero quantity required")
	}
	if input.Type == "" {
		input.Type = "in"
	}
	return u.repo.CreateStockMovement(input)
}

func (u *UseCase) ListMovements(tenantID string, productID string, limit int) ([]stock.StockMovement, error) {
	return u.repo.ListStockMovements(tenantID, productID, limit)
}

func (u *UseCase) SetServiceSupplies(tenantID string, serviceID string, supplies []stock.ServiceSupply) error {
	return u.repo.SetServiceSupplies(tenantID, serviceID, supplies)
}

func (u *UseCase) ListServiceSupplies(tenantID string, serviceID string) ([]stock.ServiceSupply, error) {
	return u.repo.ListServiceSupplies(tenantID, serviceID)
}

func (u *UseCase) DeductSuppliesForAppointment(tenantID string, appointmentID string, serviceID string) error {
	return u.repo.DeductSuppliesForAppointment(tenantID, appointmentID, serviceID)
}


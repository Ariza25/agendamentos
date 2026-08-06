package memory

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/customer"
)

type CustomerRepository struct {
	mu        sync.Mutex
	customers map[string]customer.Customer
}

func NewCustomerRepository() *CustomerRepository {
	return &CustomerRepository{
		customers: make(map[string]customer.Customer),
	}
}

func (r *CustomerRepository) Create(input customer.Customer) (customer.Customer, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item := input
	item.ID = uuid.NewString()
	item.CreatedAt = time.Now().Format(time.RFC3339)
	item.UpdatedAt = item.CreatedAt
	r.customers[item.ID] = item
	return item, nil
}

func (r *CustomerRepository) Update(input customer.Customer) (customer.Customer, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.customers[input.ID]
	if !exists {
		return customer.Customer{}, fmt.Errorf("customer not found")
	}
	if strings.TrimSpace(input.Name) != "" {
		current.Name = input.Name
	}
	current.Phone = strings.TrimSpace(input.Phone)
	current.Email = strings.ToLower(strings.TrimSpace(input.Email))
	current.Birthday = strings.TrimSpace(input.Birthday)
	current.Notes = strings.TrimSpace(input.Notes)
	current.UpdatedAt = time.Now().Format(time.RFC3339)
	r.customers[current.ID] = current
	return current, nil
}

func (r *CustomerRepository) GetByID(tenantID string, id string) (customer.Customer, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.customers[strings.TrimSpace(id)]
	if !exists || item.TenantID != strings.TrimSpace(tenantID) {
		return customer.Customer{}, fmt.Errorf("customer not found")
	}
	return item, nil
}

func (r *CustomerRepository) FindByPhone(tenantID string, phone string) (customer.Customer, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID = strings.TrimSpace(tenantID)
	phone = onlyDigits(phone)
	for _, item := range r.customers {
		if item.TenantID == tenantID && onlyDigits(item.Phone) == phone {
			return item, nil
		}
	}
	return customer.Customer{}, fmt.Errorf("customer not found")
}

func onlyDigits(value string) string {
	var builder strings.Builder
	for _, char := range value {
		if char >= '0' && char <= '9' {
			builder.WriteRune(char)
		}
	}
	return builder.String()
}

func (r *CustomerRepository) ListByTenant(tenantID string, search string) ([]customer.Customer, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID = strings.TrimSpace(tenantID)
	search = strings.ToLower(strings.TrimSpace(search))
	items := []customer.Customer{}
	for _, item := range r.customers {
		if item.TenantID != tenantID {
			continue
		}
		if search != "" {
			haystack := strings.ToLower(strings.Join([]string{item.Name, item.Phone, item.Email}, " "))
			if !strings.Contains(haystack, search) {
				continue
			}
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *CustomerRepository) Delete(tenantID string, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.customers[strings.TrimSpace(id)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return fmt.Errorf("customer not found")
	}
	delete(r.customers, current.ID)
	return nil
}

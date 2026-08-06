package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/customer"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Create(input customer.Customer) (customer.Customer, error) {
	item := input
	item.ID = uuid.NewString()
	if _, err := r.db.Exec(`
		INSERT INTO catalog.customers (id, tenant_id, name, phone, email, birthday, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, item.ID, item.TenantID, item.Name, item.Phone, item.Email, item.Birthday, item.Notes); err != nil {
		return customer.Customer{}, fmt.Errorf("insert customer: %w", err)
	}
	return r.GetByID(item.TenantID, item.ID)
}

func (r *CustomerRepository) Update(input customer.Customer) (customer.Customer, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.customers
		SET name = COALESCE(NULLIF($3, ''), name),
		    phone = COALESCE(NULLIF($4, ''), phone),
		    email = COALESCE(NULLIF(LOWER($5), ''), email),
		    birthday = COALESCE(NULLIF($6, ''), birthday),
		    notes = COALESCE(NULLIF($7, ''), notes),
		    updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`, input.ID, input.TenantID, strings.TrimSpace(input.Name), strings.TrimSpace(input.Phone), strings.TrimSpace(input.Email), strings.TrimSpace(input.Birthday), strings.TrimSpace(input.Notes))
	if err != nil {
		return customer.Customer{}, fmt.Errorf("update customer: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return customer.Customer{}, fmt.Errorf("check customer update: %w", err)
	}
	if affected == 0 {
		return customer.Customer{}, fmt.Errorf("customer not found")
	}
	return r.GetByID(input.TenantID, input.ID)
}

func (r *CustomerRepository) GetByID(tenantID string, id string) (customer.Customer, error) {
	item, err := scanCustomerRow(r.db.QueryRow(`
		SELECT id, tenant_id, name, phone, email, birthday, notes, created_at, updated_at
		FROM catalog.customers
		WHERE id = $1 AND tenant_id = $2
	`, strings.TrimSpace(id), strings.TrimSpace(tenantID)))
	if err != nil {
		return customer.Customer{}, err
	}
	return item, nil
}

func (r *CustomerRepository) FindByPhone(tenantID string, phone string) (customer.Customer, error) {
	item, err := scanCustomerRow(r.db.QueryRow(`
		SELECT id, tenant_id, name, phone, email, birthday, notes, created_at, updated_at
		FROM catalog.customers
		WHERE tenant_id = $1 AND regexp_replace(phone, '\D', '', 'g') = regexp_replace($2, '\D', '', 'g')
	`, strings.TrimSpace(tenantID), strings.TrimSpace(phone)))
	if err != nil {
		return customer.Customer{}, err
	}
	return item, nil
}

func (r *CustomerRepository) ListByTenant(tenantID string, search string) ([]customer.Customer, error) {
	tenantID = strings.TrimSpace(tenantID)
	search = strings.TrimSpace(search)
	query := `
		SELECT id, tenant_id, name, phone, email, birthday, notes, created_at, updated_at
		FROM catalog.customers
		WHERE tenant_id = $1
	`
	args := []any{tenantID}
	if search != "" {
		query += ` AND (LOWER(name) LIKE '%' || LOWER($2) || '%' OR LOWER(phone) LIKE '%' || LOWER($2) || '%' OR LOWER(email) LIKE '%' || LOWER($2) || '%')`
		args = append(args, search)
	}
	query += ` ORDER BY name ASC`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query customers: %w", err)
	}
	defer rows.Close()

	items := []customer.Customer{}
	for rows.Next() {
		item, err := scanCustomerRows(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *CustomerRepository) Delete(tenantID string, id string) error {
	result, err := r.db.Exec(`
		DELETE FROM catalog.customers
		WHERE id = $1 AND tenant_id = $2
	`, id, tenantID)
	if err != nil {
		return fmt.Errorf("delete customer: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check customer delete: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("customer not found")
	}
	return nil
}

type customerScanner interface {
	Scan(dest ...any) error
}

func scanCustomerRow(scanner customerScanner) (customer.Customer, error) {
	return scanCustomerRows(scanner)
}

func scanCustomerRows(scanner customerScanner) (customer.Customer, error) {
	var item customer.Customer
	var createdAt, updatedAt time.Time
	if err := scanner.Scan(&item.ID, &item.TenantID, &item.Name, &item.Phone, &item.Email, &item.Birthday, &item.Notes, &createdAt, &updatedAt); err != nil {
		if err == sql.ErrNoRows {
			return customer.Customer{}, fmt.Errorf("customer not found")
		}
		return customer.Customer{}, fmt.Errorf("scan customer: %w", err)
	}
	item.CreatedAt = createdAt.Format(time.RFC3339)
	item.UpdatedAt = updatedAt.Format(time.RFC3339)
	return item, nil
}

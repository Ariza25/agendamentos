package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/waitinglist"
)

type WaitingListRepository struct {
	db *sql.DB
}

func NewWaitingListRepository(db *sql.DB) *WaitingListRepository {
	return &WaitingListRepository{db: db}
}

func (r *WaitingListRepository) Create(input waitinglist.Entry) (waitinglist.Entry, error) {
	item := input
	item.ID = uuid.NewString()
	item.CreatedAt = time.Now().Format(time.RFC3339)
	if item.Status == "" {
		item.Status = waitinglist.StatusWaiting
	}
	if _, err := r.db.Exec(`
		INSERT INTO catalog.waiting_list (id, tenant_id, service_id, professional_id, customer_id, name, phone, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, item.ID, item.TenantID, item.ServiceID, item.ProfessionalID, item.CustomerID, item.Name, item.Phone, item.Status); err != nil {
		return waitinglist.Entry{}, fmt.Errorf("insert waiting list entry: %w", err)
	}
	return r.GetByID(item.TenantID, item.ID)
}

func (r *WaitingListRepository) GetByID(tenantID string, id string) (waitinglist.Entry, error) {
	return scanWaitingRow(r.db.QueryRow(`
		SELECT id, tenant_id, service_id, professional_id, customer_id, name, phone, status, created_at, notified_at
		FROM catalog.waiting_list
		WHERE id = $1 AND tenant_id = $2
	`, strings.TrimSpace(id), strings.TrimSpace(tenantID)))
}

func (r *WaitingListRepository) ListByTenant(tenantID string, status string) ([]waitinglist.Entry, error) {
	query := `
		SELECT id, tenant_id, service_id, professional_id, customer_id, name, phone, status, created_at, notified_at
		FROM catalog.waiting_list
		WHERE tenant_id = $1
	`
	args := []any{strings.TrimSpace(tenantID)}
	if strings.TrimSpace(status) != "" {
		query += " AND status = $2"
		args = append(args, strings.TrimSpace(status))
	}
	query += " ORDER BY created_at ASC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query waiting list: %w", err)
	}
	defer rows.Close()

	items := []waitinglist.Entry{}
	for rows.Next() {
		item, err := scanWaitingRows(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *WaitingListRepository) UpdateStatus(tenantID string, id string, status string) (waitinglist.Entry, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.waiting_list
		SET status = $3,
		    notified_at = CASE WHEN $3 = 'notified' THEN NOW() ELSE notified_at END
		WHERE id = $1 AND tenant_id = $2
	`, id, tenantID, status)
	if err != nil {
		return waitinglist.Entry{}, fmt.Errorf("update waiting list status: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return waitinglist.Entry{}, fmt.Errorf("check waiting list update: %w", err)
	}
	if affected == 0 {
		return waitinglist.Entry{}, fmt.Errorf("waiting list entry not found")
	}
	return r.GetByID(tenantID, id)
}

func (r *WaitingListRepository) ListWaitingForSlot(tenantID string, serviceID string, professionalID string, order string) ([]waitinglist.Entry, error) {
	query := `
		SELECT id, tenant_id, service_id, professional_id, customer_id, name, phone, status, created_at, notified_at
		FROM catalog.waiting_list
		WHERE tenant_id = $1 AND status = 'waiting'
	`
	args := []any{strings.TrimSpace(tenantID)}
	index := 2
	if strings.TrimSpace(serviceID) != "" {
		query += fmt.Sprintf(" AND (service_id = $%d OR service_id = '')", index)
		args = append(args, strings.TrimSpace(serviceID))
		index++
	}
	if strings.TrimSpace(professionalID) != "" {
		query += fmt.Sprintf(" AND (professional_id = $%d OR professional_id = '')", index)
		args = append(args, strings.TrimSpace(professionalID))
		index++
	}
	if strings.TrimSpace(order) == "oldest" {
		query += " ORDER BY created_at ASC"
	} else {
		query += " ORDER BY created_at DESC"
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query waiting list slot: %w", err)
	}
	defer rows.Close()

	items := []waitinglist.Entry{}
	for rows.Next() {
		item, err := scanWaitingRows(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

type waitingScanner interface {
	Scan(dest ...any) error
}

func scanWaitingRow(scanner waitingScanner) (waitinglist.Entry, error) {
	return scanWaitingRows(scanner)
}

func scanWaitingRows(scanner waitingScanner) (waitinglist.Entry, error) {
	var item waitinglist.Entry
	var createdAt, notifiedAt sql.NullTime
	if err := scanner.Scan(&item.ID, &item.TenantID, &item.ServiceID, &item.ProfessionalID, &item.CustomerID, &item.Name, &item.Phone, &item.Status, &createdAt, &notifiedAt); err != nil {
		if err == sql.ErrNoRows {
			return waitinglist.Entry{}, fmt.Errorf("waiting list entry not found")
		}
		return waitinglist.Entry{}, fmt.Errorf("scan waiting list entry: %w", err)
	}
	item.CreatedAt = createdAt.Time.Format(time.RFC3339)
	if notifiedAt.Valid {
		item.NotifiedAt = notifiedAt.Time.UTC().Format(time.RFC3339)
	}
	return item, nil
}

package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/review"
)

type ReviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(input review.Review) (review.Review, error) {
	item := input
	if item.ID == "" {
		item.ID = uuid.NewString()
	}
	item.CreatedAt = time.Now().Format(time.RFC3339)
	if _, err := r.db.Exec(`
		INSERT INTO catalog.reviews (id, tenant_id, appointment_id, customer_id, status, rating, comment, invited_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, item.ID, item.TenantID, item.AppointmentID, item.CustomerID, item.Status, item.Rating, item.Comment, item.InvitedAt); err != nil {
		return review.Review{}, fmt.Errorf("insert review: %w", err)
	}
	return r.GetByID(item.TenantID, item.ID)
}

func (r *ReviewRepository) GetByID(tenantID string, id string) (review.Review, error) {
	return scanReviewRow(r.db.QueryRow(`
		SELECT r.id, r.tenant_id, r.appointment_id, r.customer_id, r.status, r.rating, r.comment, r.invited_at, r.submitted_at, r.created_at, c.name, t.slug
		FROM catalog.reviews r
		LEFT JOIN catalog.customers c ON c.id = r.customer_id AND c.tenant_id = r.tenant_id
		LEFT JOIN catalog.tenants t ON t.id = r.tenant_id
		WHERE r.id = $1 AND r.tenant_id = $2
	`, strings.TrimSpace(id), strings.TrimSpace(tenantID)))
}

func (r *ReviewRepository) GetByAppointment(tenantID string, appointmentID string) (review.Review, error) {
	return scanReviewRow(r.db.QueryRow(`
		SELECT r.id, r.tenant_id, r.appointment_id, r.customer_id, r.status, r.rating, r.comment, r.invited_at, r.submitted_at, r.created_at, c.name, t.slug
		FROM catalog.reviews r
		LEFT JOIN catalog.customers c ON c.id = r.customer_id AND c.tenant_id = r.tenant_id
		LEFT JOIN catalog.tenants t ON t.id = r.tenant_id
		WHERE r.tenant_id = $1 AND r.appointment_id = $2
		ORDER BY r.created_at DESC
		LIMIT 1
	`, strings.TrimSpace(tenantID), strings.TrimSpace(appointmentID)))
}

func (r *ReviewRepository) ListByTenant(tenantID string, status string) ([]review.Review, error) {
	query := `
		SELECT r.id, r.tenant_id, r.appointment_id, r.customer_id, r.status, r.rating, r.comment, r.invited_at, r.submitted_at, r.created_at, c.name, t.slug
		FROM catalog.reviews r
		LEFT JOIN catalog.customers c ON c.id = r.customer_id AND c.tenant_id = r.tenant_id
		LEFT JOIN catalog.tenants t ON t.id = r.tenant_id
		WHERE r.tenant_id = $1`
	args := []any{strings.TrimSpace(tenantID)}
	if status = strings.TrimSpace(status); status != "" {
		args = append(args, status)
		query += fmt.Sprintf(" AND r.status = $%d", len(args))
	}
	query += " ORDER BY r.created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query reviews: %w", err)
	}
	defer rows.Close()

	items := []review.Review{}
	for rows.Next() {
		item, err := scanReviewRows(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *ReviewRepository) Submit(tenantID string, id string, rating int, comment string) (review.Review, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.reviews
		SET status = 'submitted', rating = $3, comment = $4, submitted_at = NOW()
		WHERE id = $1 AND tenant_id = $2 AND status = 'invited'
	`, strings.TrimSpace(id), strings.TrimSpace(tenantID), rating, comment)
	if err != nil {
		return review.Review{}, fmt.Errorf("submit review: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return review.Review{}, fmt.Errorf("check review update: %w", err)
	}
	if affected == 0 {
		return review.Review{}, fmt.Errorf("review not found")
	}
	return r.GetByID(tenantID, id)
}

type reviewScanner interface {
	Scan(dest ...any) error
}

func scanReviewRow(scanner reviewScanner) (review.Review, error) {
	return scanReviewRows(scanner)
}

func scanReviewRows(scanner reviewScanner) (review.Review, error) {
	var item review.Review
	var submittedAt sql.NullTime
	var invitedAt, createdAt time.Time
	var customerName, tenantSlug sql.NullString
	if err := scanner.Scan(&item.ID, &item.TenantID, &item.AppointmentID, &item.CustomerID, &item.Status, &item.Rating, &item.Comment, &invitedAt, &submittedAt, &createdAt, &customerName, &tenantSlug); err != nil {
		if err == sql.ErrNoRows {
			return review.Review{}, fmt.Errorf("review not found")
		}
		return review.Review{}, fmt.Errorf("scan review: %w", err)
	}
	item.InvitedAt = invitedAt.Format(time.RFC3339)
	if submittedAt.Valid {
		item.SubmittedAt = submittedAt.Time.Format(time.RFC3339)
	}
	if customerName.Valid {
		item.CustomerName = customerName.String
	}
	if tenantSlug.Valid {
		item.TenantSlug = tenantSlug.String
	}
	item.CreatedAt = createdAt.Format(time.RFC3339)
	return item, nil
}

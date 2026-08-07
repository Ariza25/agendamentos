package memory

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/review"
)

type ReviewRepository struct {
	mu      sync.Mutex
	reviews map[string]review.Review
}

func NewReviewRepository() *ReviewRepository {
	return &ReviewRepository{
		reviews: make(map[string]review.Review),
	}
}

func (r *ReviewRepository) Create(input review.Review) (review.Review, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item := input
	if item.ID == "" {
		item.ID = uuid.NewString()
	}
	if item.Status == "" {
		item.Status = "invited"
	}
	if item.InvitedAt == "" {
		item.InvitedAt = time.Now().UTC().Format(time.RFC3339)
	}
	item.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	r.reviews[item.ID] = item
	return item, nil
}

func (r *ReviewRepository) GetByID(tenantID string, id string) (review.Review, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.reviews[strings.TrimSpace(id)]
	if !exists || item.TenantID != strings.TrimSpace(tenantID) {
		return review.Review{}, fmt.Errorf("review not found")
	}
	return item, nil
}

func (r *ReviewRepository) GetByAppointment(tenantID string, appointmentID string) (review.Review, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID = strings.TrimSpace(tenantID)
	appointmentID = strings.TrimSpace(appointmentID)
	for _, item := range r.reviews {
		if item.TenantID == tenantID && item.AppointmentID == appointmentID {
			return item, nil
		}
	}
	return review.Review{}, fmt.Errorf("review not found")
}

func (r *ReviewRepository) ListByTenant(tenantID string, status string) ([]review.Review, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID = strings.TrimSpace(tenantID)
	status = strings.TrimSpace(status)
	items := []review.Review{}
	for _, item := range r.reviews {
		if item.TenantID != tenantID {
			continue
		}
		if status != "" && item.Status != status {
			continue
		}
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt > items[j].CreatedAt
	})
	return items, nil
}

func (r *ReviewRepository) Submit(tenantID string, id string, rating int, comment string) (review.Review, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.reviews[strings.TrimSpace(id)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return review.Review{}, fmt.Errorf("review not found")
	}
	current.Status = "submitted"
	current.Rating = rating
	current.Comment = comment
	current.SubmittedAt = time.Now().UTC().Format(time.RFC3339)
	r.reviews[current.ID] = current
	return current, nil
}

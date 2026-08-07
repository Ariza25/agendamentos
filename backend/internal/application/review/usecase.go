package review

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	domain "erpagendamentos/backend/internal/domain/review"
)

var (
	ErrReviewNotFound  = errors.New("review not found")
	ErrInvalidReview   = errors.New("invalid review")
	ErrReviewSubmitted = errors.New("review already submitted")
)

const (
	StatusInvited   = "invited"
	StatusSubmitted = "submitted"
)

type UseCase struct {
	repository domain.Repository
}

func NewUseCase(repository domain.Repository) *UseCase {
	return &UseCase{repository: repository}
}

func (u *UseCase) ListByTenant(tenantID string, status string) ([]domain.Review, error) {
	return u.repository.ListByTenant(strings.TrimSpace(tenantID), strings.TrimSpace(status))
}

func (u *UseCase) GetByID(tenantID string, id string) (domain.Review, error) {
	return u.repository.GetByID(strings.TrimSpace(tenantID), strings.TrimSpace(id))
}

func (u *UseCase) CreateInvitation(tenantID string, appointmentID string, customerID string) (domain.Review, error) {
	tenantID = strings.TrimSpace(tenantID)
	appointmentID = strings.TrimSpace(appointmentID)
	if tenantID == "" || appointmentID == "" {
		return domain.Review{}, ErrInvalidReview
	}
	if existing, err := u.repository.GetByAppointment(tenantID, appointmentID); err == nil {
		return existing, nil
	}
	return u.repository.Create(domain.Review{
		ID:            uuid.NewString(),
		TenantID:      tenantID,
		AppointmentID: appointmentID,
		CustomerID:    strings.TrimSpace(customerID),
		Status:        StatusInvited,
		InvitedAt:     time.Now().UTC().Format(time.RFC3339),
	})
}

func (u *UseCase) Submit(tenantID string, id string, rating int, comment string) (domain.Review, error) {
	tenantID = strings.TrimSpace(tenantID)
	id = strings.TrimSpace(id)
	comment = strings.TrimSpace(comment)
	if rating < 1 || rating > 5 {
		return domain.Review{}, ErrInvalidReview
	}
	current, err := u.repository.GetByID(tenantID, id)
	if err != nil {
		return domain.Review{}, ErrReviewNotFound
	}
	if current.Status != StatusInvited {
		return domain.Review{}, ErrReviewSubmitted
	}
	return u.repository.Submit(tenantID, id, rating, comment)
}

package appointment

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	domain "erpagendamentos/backend/internal/domain/appointment"
)

const (
	StatusScheduled  = "scheduled"
	StatusConfirmed  = "confirmed"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusCancelled  = "cancelled"
	StatusNoShow     = "no_show"
)

var (
	ErrAppointmentNotFound = errors.New("appointment not found")
	ErrInvalidAppointment  = errors.New("invalid appointment")
	ErrTimeConflict        = errors.New("time conflict")
	ErrInvalidStatusChange = errors.New("invalid status change")
)

var allowedStatuses = map[string]bool{
	StatusScheduled:  true,
	StatusConfirmed:  true,
	StatusInProgress: true,
	StatusCompleted:  true,
	StatusCancelled:  true,
	StatusNoShow:     true,
}

type UseCase struct {
	repository domain.Repository
}

type RecurrenceInput struct {
	Frequency string `json:"frequency"`
	Count     int    `json:"count"`
}

func NewUseCase(repository domain.Repository) *UseCase {
	return &UseCase{repository: repository}
}

func (u *UseCase) ListByRange(tenantID string, from string, to string, professionalID string, customerID string, status string) ([]domain.Appointment, error) {
	return u.repository.ListByRange(
		strings.TrimSpace(tenantID),
		strings.TrimSpace(from),
		strings.TrimSpace(to),
		strings.TrimSpace(professionalID),
		strings.TrimSpace(customerID),
		strings.TrimSpace(status),
	)
}

func (u *UseCase) GetByID(tenantID string, id string) (domain.Appointment, error) {
	item, err := u.repository.GetByID(strings.TrimSpace(tenantID), strings.TrimSpace(id))
	if err != nil {
		return domain.Appointment{}, ErrAppointmentNotFound
	}
	return item, nil
}

func (u *UseCase) GetByManagementCode(tenantID string, codeHash string) (domain.Appointment, error) {
	item, err := u.repository.GetByManagementCodeHash(strings.TrimSpace(tenantID), strings.TrimSpace(codeHash))
	if err != nil {
		return domain.Appointment{}, ErrAppointmentNotFound
	}
	return item, nil
}

func (u *UseCase) Create(input domain.Appointment) (domain.Appointment, error) {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.ProfessionalID = strings.TrimSpace(input.ProfessionalID)
	input.ServiceID = strings.TrimSpace(input.ServiceID)
	input.CustomerID = strings.TrimSpace(input.CustomerID)
	input.Notes = strings.TrimSpace(input.Notes)
	input.Status = strings.TrimSpace(input.Status)

	if input.TenantID == "" || input.ProfessionalID == "" || input.CustomerID == "" {
		return domain.Appointment{}, ErrInvalidAppointment
	}
	startsAt, err := parseRFC3339(input.StartsAt)
	if err != nil {
		return domain.Appointment{}, ErrInvalidAppointment
	}
	endsAt, err := parseRFC3339(input.EndsAt)
	if err != nil {
		return domain.Appointment{}, ErrInvalidAppointment
	}
	if !endsAt.After(startsAt) {
		return domain.Appointment{}, ErrInvalidAppointment
	}
	if input.Status == "" {
		input.Status = StatusScheduled
	}
	if !allowedStatuses[input.Status] {
		return domain.Appointment{}, ErrInvalidAppointment
	}

	input.StartsAt = startsAt.UTC().Format(time.RFC3339)
	input.EndsAt = endsAt.UTC().Format(time.RFC3339)

	conflict, err := u.repository.HasConflict(input.TenantID, input.ProfessionalID, input.StartsAt, input.EndsAt, "")
	if err != nil {
		return domain.Appointment{}, err
	}
	if conflict {
		return domain.Appointment{}, ErrTimeConflict
	}

	input.ID = uuid.NewString()
	return u.repository.Create(input)
}

func (u *UseCase) CreateSeries(input domain.Appointment, recurrence RecurrenceInput) ([]domain.Appointment, error) {
	if recurrence.Frequency != "daily" && recurrence.Frequency != "weekly" {
		return nil, ErrInvalidAppointment
	}
	if recurrence.Count <= 0 || recurrence.Count > 60 {
		recurrence.Count = 4
	}
	base, err := parseRFC3339(input.StartsAt)
	if err != nil {
		return nil, ErrInvalidAppointment
	}
	input.SeriesID = uuid.NewString()

	created := []domain.Appointment{}
	for index := 0; index < recurrence.Count; index++ {
		item := input
		if index > 0 {
			if recurrence.Frequency == "daily" {
				base = base.AddDate(0, 0, 1)
			} else {
				base = base.AddDate(0, 0, 7)
			}
			item.StartsAt = base.Format(time.RFC3339)
			itemEndsAt, err := parseRFC3339(input.EndsAt)
			if err != nil {
				return nil, ErrInvalidAppointment
			}
			offset := itemEndsAt.Sub(base.Add(time.Duration(-1) * time.Hour))
			item.EndsAt = base.Add(offset).Format(time.RFC3339)
		}
		appointment, err := u.Create(item)
		if err != nil {
			if index == 0 {
				return nil, err
			}
			break
		}
		created = append(created, appointment)
	}
	if len(created) == 0 {
		return nil, ErrTimeConflict
	}
	return created, nil
}

func (u *UseCase) CheckIn(tenantID string, id string, updatedBy string) (domain.Appointment, error) {
	tenantID = strings.TrimSpace(tenantID)
	id = strings.TrimSpace(id)
	current, err := u.repository.GetByID(tenantID, id)
	if err != nil {
		return domain.Appointment{}, ErrAppointmentNotFound
	}
	if current.Status != StatusScheduled && current.Status != StatusConfirmed {
		return domain.Appointment{}, ErrInvalidStatusChange
	}
	current.Status = StatusInProgress
	current.CheckInAt = time.Now().UTC().Format(time.RFC3339)
	return u.repository.Update(current)
}

func (u *UseCase) MarkNoShow(tenantID string, id string, updatedBy string, feeCents int) (domain.Appointment, error) {
	tenantID = strings.TrimSpace(tenantID)
	id = strings.TrimSpace(id)
	current, err := u.repository.GetByID(tenantID, id)
	if err != nil {
		return domain.Appointment{}, ErrAppointmentNotFound
	}
	if current.Status != StatusScheduled && current.Status != StatusConfirmed && current.Status != StatusInProgress {
		return domain.Appointment{}, ErrInvalidStatusChange
	}
	current.Status = StatusNoShow
	current.NoShowAt = time.Now().UTC().Format(time.RFC3339)
	current.NoShowFeeCents = feeCents
	current.PaymentStatus = "no_show_fee"
	return u.repository.Update(current)
}

func (u *UseCase) ListDueReminders(from string, to string) ([]domain.Appointment, error) {
	return u.repository.ListDueReminders(strings.TrimSpace(from), strings.TrimSpace(to))
}

func (u *UseCase) SetReminderSent(tenantID string, id string) error {
	return u.repository.SetReminderSent(strings.TrimSpace(tenantID), strings.TrimSpace(id))
}

func (u *UseCase) Analytics(tenantID string, from string, to string) (domain.Analytics, error) {
	return u.repository.Analytics(strings.TrimSpace(tenantID), strings.TrimSpace(from), strings.TrimSpace(to))
}

func (u *UseCase) Update(input domain.Appointment) (domain.Appointment, error) {
	tenantID := strings.TrimSpace(input.TenantID)
	id := strings.TrimSpace(input.ID)
	current, err := u.repository.GetByID(tenantID, id)
	if err != nil {
		return domain.Appointment{}, ErrAppointmentNotFound
	}

	input.TenantID = tenantID
	input.ID = id
	input.Status = strings.TrimSpace(input.Status)
	if input.Status == "" {
		input.Status = current.Status
	}
	if !allowedStatuses[input.Status] {
		return domain.Appointment{}, ErrInvalidAppointment
	}
	input.ProfessionalID = strings.TrimSpace(input.ProfessionalID)
	if input.ProfessionalID == "" {
		input.ProfessionalID = current.ProfessionalID
	}
	input.ServiceID = strings.TrimSpace(input.ServiceID)
	if input.ServiceID == "" {
		input.ServiceID = current.ServiceID
	}
	input.CustomerID = strings.TrimSpace(input.CustomerID)
	if input.CustomerID == "" {
		input.CustomerID = current.CustomerID
	}
	input.CreatedBy = current.CreatedBy

	startsAtRaw := strings.TrimSpace(input.StartsAt)
	if startsAtRaw == "" {
		startsAtRaw = current.StartsAt
	}
	endsAtRaw := strings.TrimSpace(input.EndsAt)
	if endsAtRaw == "" {
		endsAtRaw = current.EndsAt
	}
	startsAt, err := parseRFC3339(startsAtRaw)
	if err != nil {
		return domain.Appointment{}, ErrInvalidAppointment
	}
	endsAt, err := parseRFC3339(endsAtRaw)
	if err != nil {
		return domain.Appointment{}, ErrInvalidAppointment
	}
	if !endsAt.After(startsAt) {
		return domain.Appointment{}, ErrInvalidAppointment
	}
	input.StartsAt = startsAt.UTC().Format(time.RFC3339)
	input.EndsAt = endsAt.UTC().Format(time.RFC3339)

	conflict, err := u.repository.HasConflict(tenantID, input.ProfessionalID, input.StartsAt, input.EndsAt, id)
	if err != nil {
		return domain.Appointment{}, err
	}
	if conflict {
		return domain.Appointment{}, ErrTimeConflict
	}

	return u.repository.Update(input)
}

func (u *UseCase) UpdateStatus(tenantID string, id string, status string, updatedBy string) (domain.Appointment, error) {
	tenantID = strings.TrimSpace(tenantID)
	id = strings.TrimSpace(id)
	status = strings.TrimSpace(status)
	updatedBy = strings.TrimSpace(updatedBy)
	if !allowedStatuses[status] {
		return domain.Appointment{}, ErrInvalidAppointment
	}

	current, err := u.repository.GetByID(tenantID, id)
	if err != nil {
		return domain.Appointment{}, ErrAppointmentNotFound
	}
	if !validStatusChange(current.Status, status) {
		return domain.Appointment{}, ErrInvalidStatusChange
	}

	return u.repository.UpdateStatus(tenantID, id, status, updatedBy)
}

func (u *UseCase) Summary(tenantID string, date string) (map[string]any, error) {
	tenantID = strings.TrimSpace(tenantID)
	date = strings.TrimSpace(date)
	day, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, ErrInvalidAppointment
	}
	from := day.Format("2006-01-02T00:00:00Z07:00")
	to := day.AddDate(0, 0, 1).Format("2006-01-02T00:00:00Z07:00")

	counts, err := u.repository.CountByStatus(tenantID, from, to)
	if err != nil {
		return nil, err
	}

	result := map[string]any{"date": date}
	total := 0
	for _, count := range counts {
		result[count.Status] = count.Count
		total += count.Count
	}
	result["total"] = total
	return result, nil
}

func (u *UseCase) Delete(tenantID string, id string) error {
	if err := u.repository.Delete(strings.TrimSpace(tenantID), strings.TrimSpace(id)); err != nil {
		return ErrAppointmentNotFound
	}
	return nil
}

func parseRFC3339(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, fmt.Errorf("empty time")
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		parsed, err = time.ParseInLocation("2006-01-02T15:04", value, time.Local)
	}
	return parsed, err
}

func validStatusChange(from string, to string) bool {
	if from == to {
		return true
	}
	switch from {
	case StatusScheduled, StatusConfirmed:
		return true
	case StatusInProgress:
		return to == StatusCompleted || to == StatusCancelled
	case StatusCompleted, StatusCancelled, StatusNoShow:
		return false
	default:
		return true
	}
}

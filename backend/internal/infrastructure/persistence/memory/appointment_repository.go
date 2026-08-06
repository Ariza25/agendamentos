package memory

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/appointment"
)

type AppointmentRepository struct {
	mu           sync.Mutex
	appointments map[string]appointment.Appointment
}

func NewAppointmentRepository() *AppointmentRepository {
	return &AppointmentRepository{
		appointments: make(map[string]appointment.Appointment),
	}
}

func (r *AppointmentRepository) Create(input appointment.Appointment) (appointment.Appointment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item := input
	item.ID = uuid.NewString()
	item.CreatedAt = time.Now().Format(time.RFC3339)
	item.UpdatedAt = item.CreatedAt
	r.appointments[item.ID] = item
	return item, nil
}

func (r *AppointmentRepository) Update(input appointment.Appointment) (appointment.Appointment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.appointments[input.ID]
	if !exists {
		return appointment.Appointment{}, fmt.Errorf("appointment not found")
	}
	input.CreatedAt = current.CreatedAt
	input.UpdatedAt = time.Now().Format(time.RFC3339)
	r.appointments[input.ID] = input
	return input, nil
}

func (r *AppointmentRepository) GetByID(tenantID string, id string) (appointment.Appointment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.appointments[strings.TrimSpace(id)]
	if !exists || item.TenantID != strings.TrimSpace(tenantID) {
		return appointment.Appointment{}, fmt.Errorf("appointment not found")
	}
	return item, nil
}

func (r *AppointmentRepository) ListByRange(tenantID string, from string, to string, professionalID string, customerID string, status string) ([]appointment.Appointment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID = strings.TrimSpace(tenantID)
	professionalID = strings.TrimSpace(professionalID)
	customerID = strings.TrimSpace(customerID)
	status = strings.TrimSpace(status)
	fromTime, hasFrom := parseAppointmentBound(from)
	toTime, hasTo := parseAppointmentBound(to)
	items := []appointment.Appointment{}
	for _, item := range r.appointments {
		if item.TenantID != tenantID {
			continue
		}
		itemStart, itemEnd, ok := parseAppointmentInterval(item)
		if !ok {
			continue
		}
		if hasFrom && itemEnd.Before(fromTime) {
			continue
		}
		if hasTo && !itemStart.Before(toTime) {
			continue
		}
		if professionalID != "" && item.ProfessionalID != professionalID {
			continue
		}
		if customerID != "" && item.CustomerID != customerID {
			continue
		}
		if status != "" && item.Status != status {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *AppointmentRepository) UpdateStatus(tenantID string, id string, status string, updatedBy string) (appointment.Appointment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.appointments[strings.TrimSpace(id)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return appointment.Appointment{}, fmt.Errorf("appointment not found")
	}
	current.Status = strings.TrimSpace(status)
	current.UpdatedAt = time.Now().Format(time.RFC3339)
	r.appointments[current.ID] = current
	return current, nil
}

func (r *AppointmentRepository) HasConflict(tenantID string, professionalID string, startsAt string, endsAt string, excludeID string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID = strings.TrimSpace(tenantID)
	professionalID = strings.TrimSpace(professionalID)
	startsAt = strings.TrimSpace(startsAt)
	endsAt = strings.TrimSpace(endsAt)
	excludeID = strings.TrimSpace(excludeID)
	starts, err := time.Parse(time.RFC3339, startsAt)
	if err != nil {
		return false, fmt.Errorf("invalid appointment start time")
	}
	ends, err := time.Parse(time.RFC3339, endsAt)
	if err != nil {
		return false, fmt.Errorf("invalid appointment end time")
	}
	for _, item := range r.appointments {
		if item.TenantID != tenantID || item.ProfessionalID != professionalID {
			continue
		}
		if excludeID != "" && item.ID == excludeID {
			continue
		}
		if item.Status == "cancelled" || item.Status == "no_show" {
			continue
		}
		itemStart, itemEnd, ok := parseAppointmentInterval(item)
		if !ok {
			continue
		}
		if starts.Before(itemEnd) && ends.After(itemStart) {
			return true, nil
		}
	}
	return false, nil
}

func (r *AppointmentRepository) CountByStatus(tenantID string, from string, to string) ([]appointment.StatusCount, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID = strings.TrimSpace(tenantID)
	fromTime, hasFrom := parseAppointmentBound(from)
	toTime, hasTo := parseAppointmentBound(to)
	counts := map[string]int{}
	for _, item := range r.appointments {
		if item.TenantID != tenantID {
			continue
		}
		itemStart, itemEnd, ok := parseAppointmentInterval(item)
		if !ok {
			continue
		}
		if hasFrom && itemEnd.Before(fromTime) {
			continue
		}
		if hasTo && !itemStart.Before(toTime) {
			continue
		}
		counts[item.Status]++
	}
	result := []appointment.StatusCount{}
	for status, count := range counts {
		result = append(result, appointment.StatusCount{Status: status, Count: count})
	}
	return result, nil
}

func (r *AppointmentRepository) Delete(tenantID string, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.appointments[strings.TrimSpace(id)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return fmt.Errorf("appointment not found")
	}
	delete(r.appointments, current.ID)
	return nil
}

func parseAppointmentBound(value string) (time.Time, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, false
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, false
	}
	return parsed, true
}

func parseAppointmentInterval(item appointment.Appointment) (time.Time, time.Time, bool) {
	startsAt, err := time.Parse(time.RFC3339, item.StartsAt)
	if err != nil {
		return time.Time{}, time.Time{}, false
	}
	endsAt, err := time.Parse(time.RFC3339, item.EndsAt)
	if err != nil {
		return time.Time{}, time.Time{}, false
	}
	return startsAt, endsAt, true
}

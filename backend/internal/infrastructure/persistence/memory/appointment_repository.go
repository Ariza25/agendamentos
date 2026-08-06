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

func (r *AppointmentRepository) GetByManagementCodeHash(tenantID string, hash string) (appointment.Appointment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID = strings.TrimSpace(tenantID)
	hash = strings.TrimSpace(hash)
	for _, item := range r.appointments {
		if item.TenantID == tenantID && item.ManagementCodeHash == hash {
			return item, nil
		}
	}
	return appointment.Appointment{}, fmt.Errorf("appointment not found")
}

func (r *AppointmentRepository) ListDueReminders(from string, to string) ([]appointment.Appointment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	fromTime, err := time.Parse(time.RFC3339, from)
	if err != nil {
		return nil, fmt.Errorf("invalid reminder window start")
	}
	toTime, err := time.Parse(time.RFC3339, to)
	if err != nil {
		return nil, fmt.Errorf("invalid reminder window end")
	}
	items := []appointment.Appointment{}
	for _, item := range r.appointments {
		itemStart, _, ok := parseAppointmentInterval(item)
		if !ok {
			continue
		}
		if itemStart.Before(fromTime) || !itemStart.Before(toTime) {
			continue
		}
		if item.Status != "scheduled" && item.Status != "confirmed" {
			continue
		}
		if item.ReminderSentAt != "" {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *AppointmentRepository) SetReminderSent(tenantID string, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.appointments[strings.TrimSpace(id)]
	if !exists || current.TenantID != strings.TrimSpace(tenantID) {
		return fmt.Errorf("appointment not found")
	}
	current.ReminderSentAt = time.Now().Format(time.RFC3339)
	current.UpdatedAt = time.Now().Format(time.RFC3339)
	r.appointments[current.ID] = current
	return nil
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

func (r *AppointmentRepository) Analytics(tenantID string, from string, to string) (appointment.Analytics, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantID = strings.TrimSpace(tenantID)
	fromTime, hasFrom := parseAppointmentBound(from)
	toTime, hasTo := parseAppointmentBound(to)
	result := appointment.Analytics{}
	counts := map[string]int{}
	serviceCounts := map[string]*appointment.ServiceMetric{}
	professionalCounts := map[string]*appointment.ProfessionalMetric{}
	dayCounts := map[string]int{}
	byDay := map[string]*appointment.DayMetric{}

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
		result.TotalAppointments++
		counts[item.Status]++
		if serviceCounts[item.ServiceID] == nil {
			serviceCounts[item.ServiceID] = &appointment.ServiceMetric{ServiceID: item.ServiceID}
		}
		serviceCounts[item.ServiceID].Count++
		if professionalCounts[item.ProfessionalID] == nil {
			professionalCounts[item.ProfessionalID] = &appointment.ProfessionalMetric{ProfessionalID: item.ProfessionalID}
		}
		professionalCounts[item.ProfessionalID].Count++
		date := itemStart.Local().Format("2006-01-02")
		if dayCounts[date] == 0 {
			dayCounts[date] = 1
			byDay[date] = &appointment.DayMetric{Date: date}
		} else {
			dayCounts[date]++
			byDay[date].Count = dayCounts[date]
		}
	}
	result.Scheduled = counts["scheduled"]
	result.Confirmed = counts["confirmed"]
	result.Completed = counts["completed"]
	result.Cancelled = counts["cancelled"]
	result.NoShow = counts["no_show"]
	if result.TotalAppointments > 0 {
		result.NoShowRate = float64(result.NoShow) / float64(result.TotalAppointments) * 100
	}
	for _, metric := range serviceCounts {
		result.ByService = append(result.ByService, *metric)
	}
	for _, metric := range professionalCounts {
		result.ByProfessional = append(result.ByProfessional, *metric)
	}
	for _, metric := range byDay {
		result.ByDay = append(result.ByDay, *metric)
	}
	return result, nil
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

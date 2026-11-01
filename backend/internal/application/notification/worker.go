package notification

import (
	"context"
	"log"
	"strings"
	"time"

	appappointment "erpagendamentos/backend/internal/application/appointment"
	appointmentdomain "erpagendamentos/backend/internal/domain/appointment"
	configurationdomain "erpagendamentos/backend/internal/domain/configuration"
	customerdomain "erpagendamentos/backend/internal/domain/customer"
	domain "erpagendamentos/backend/internal/domain/notification"
	servicedomain "erpagendamentos/backend/internal/domain/service"
)

type ReminderWorker struct {
	appointments   *appappointment.UseCase
	notifications  *UseCase
	customers      customerdomain.Repository
	services       servicedomain.Repository
	configurations configurationdomain.Repository
	interval       time.Duration
	enabled        bool
}

func NewReminderWorker(
	appointments *appappointment.UseCase,
	notifications *UseCase,
	customers customerdomain.Repository,
	services servicedomain.Repository,
	configurations configurationdomain.Repository,
	interval time.Duration,
) *ReminderWorker {
	if interval <= 0 {
		interval = 30 * time.Minute
	}
	return &ReminderWorker{
		appointments:   appointments,
		notifications:  notifications,
		customers:      customers,
		services:       services,
		configurations: configurations,
		interval:       interval,
		enabled:        appointments != nil && notifications != nil,
	}
}

func (w *ReminderWorker) Start(ctx context.Context) {
	if w == nil || !w.enabled {
		return
	}
	go w.loop(ctx)
}

func (w *ReminderWorker) loop(ctx context.Context) {
	w.runOnce()
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.runOnce()
		}
	}
}

func (w *ReminderWorker) runOnce() {
	now := time.Now().UTC()
	from := now.Add(-24 * time.Hour).Format(time.RFC3339)
	to := now.Add(48 * time.Hour).Format(time.RFC3339)

	appointments, err := w.appointments.ListDueReminders(from, to)
	if err != nil {
		log.Printf(`{"level":"ERROR","message":"reminder worker failed to load due appointments","error":"%v"}`, err)
		return
	}

	sent := 0
	for _, appointment := range appointments {
		if w.sendReminder(appointment) {
			sent++
		}
	}

	log.Printf(`{"level":"INFO","message":"reminder worker completed","due":%d,"sent":%d}`, len(appointments), sent)
}

func (w *ReminderWorker) sendReminder(appointment appointmentdomain.Appointment) bool {
	configuration, err := w.configurations.Get(appointment.TenantID)
	if err != nil {
		return false
	}
	if !configuration.NotificationsEnabled {
		return false
	}

	customer, err := w.customers.GetByID(appointment.TenantID, appointment.CustomerID)
	if err != nil {
		return false
	}

	recipient := strings.TrimSpace(customer.Phone)
	channel := strings.TrimSpace(configuration.NotificationChannel)
	if channel == "" {
		channel = "whatsapp"
	}
	if channel == "email" {
		recipient = strings.TrimSpace(customer.Email)
	}
	if recipient == "" {
		return false
	}

	serviceName := "Agendamento"
	if appointment.ServiceID != "" {
		if service, err := w.services.GetByID(appointment.TenantID, appointment.ServiceID); err == nil {
			serviceName = service.Name
		}
	}

	hours := configuration.NotificationReminderHours
	if hours <= 0 {
		hours = 24
	}
	when := formatAppointmentLocal(appointment.StartsAt)
	subject, body := RenderReminder(TemplateData{
		EstablishmentName: configuration.EstablishmentName,
		CustomerName:      customer.Name,
		ServiceName:       serviceName,
		StartsAt:          when,
	})

	logEntry, err := w.notifications.Send(domain.Log{
		TenantID:      appointment.TenantID,
		AppointmentID: appointment.ID,
		Type:          string(domain.TypeReminder),
		Channel:       channel,
		Recipient:     recipient,
		Subject:       subject,
		Body:          body,
	})
	if err != nil {
		log.Printf(`{"level":"WARN","message":"reminder log failed","error":"%v"}`, err)
		return false
	}
	if logEntry.Status == "failed" {
		log.Printf(`{"level":"WARN","message":"reminder send failed","appointment":"%s","error":"%s"}`, appointment.ID, logEntry.Error)
		return false
	}

	if err := w.appointments.SetReminderSent(appointment.TenantID, appointment.ID); err != nil {
		log.Printf(`{"level":"WARN","message":"reminder flag failed","error":"%v"}`, err)
		return false
	}
	return true
}

func formatAppointmentLocal(value string) string {
	parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(value))
	if err != nil {
		return strings.TrimSpace(value)
	}
	return parsed.Local().Format("02/01/2006 às 15:04")
}

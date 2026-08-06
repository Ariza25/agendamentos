package calendar

import (
	"fmt"
	"strings"
	"time"

	appointmentdomain "erpagendamentos/backend/internal/domain/appointment"
)

func BuildICalCalendar(appointments []appointmentdomain.Appointment, establishmentName string) string {
	var builder strings.Builder
	now := time.Now().UTC()
	builder.WriteString("BEGIN:VCALENDAR\r\n")
	builder.WriteString("VERSION:2.0\r\n")
	builder.WriteString("PRODID:-//ERPAgendamentos//Calendar//PT\r\n")
	builder.WriteString("CALSCALE:GREGORIAN\r\n")
	builder.WriteString(fmt.Sprintf("X-WR-CALNAME:%s\r\n", sanitizeLine(establishmentName)))
	builder.WriteString(fmt.Sprintf("X-WR-TIMEZONE:America/Sao_Paulo\r\n"))
	for _, appointment := range appointments {
		if appointment.Status == "cancelled" || appointment.Status == "no_show" {
			continue
		}
		startsAt, startErr := time.Parse(time.RFC3339, appointment.StartsAt)
		endsAt, endErr := time.Parse(time.RFC3339, appointment.EndsAt)
		if startErr != nil || endErr != nil {
			continue
		}
		uid := appointment.ID
		if uid == "" {
			uid = fmt.Sprintf("%s-%s", appointment.TenantID, appointment.StartsAt)
		}
		builder.WriteString("BEGIN:VEVENT\r\n")
		builder.WriteString(fmt.Sprintf("UID:%s\r\n", uid))
		builder.WriteString(fmt.Sprintf("DTSTAMP:%s\r\n", formatUTC(now)))
		builder.WriteString(fmt.Sprintf("DTSTART:%s\r\n", formatUTC(startsAt)))
		builder.WriteString(fmt.Sprintf("DTEND:%s\r\n", formatUTC(endsAt)))
		summary := "Agendamento"
		if strings.TrimSpace(appointment.ServiceID) != "" {
			summary = appointment.ServiceID
		}
		builder.WriteString(fmt.Sprintf("SUMMARY:%s\r\n", sanitizeLine(summary)))
		builder.WriteString("STATUS:CONFIRMED\r\n")
		builder.WriteString("END:VEVENT\r\n")
	}
	builder.WriteString("END:VCALENDAR\r\n")
	return builder.String()
}

func formatUTC(value time.Time) string {
	return value.UTC().Format("20060102T150405Z")
}

func sanitizeLine(value string) string {
	value = strings.ReplaceAll(value, "\r", " ")
	value = strings.ReplaceAll(value, "\n", " ")
	return strings.TrimSpace(value)
}

package rest

import (
	"net/http"
	"time"

	appappointment "erpagendamentos/backend/internal/application/appointment"
	appconfiguration "erpagendamentos/backend/internal/application/configuration"
	"erpagendamentos/backend/internal/infrastructure/calendar"
)

type CalendarHandler struct {
	appointments  *appappointment.UseCase
	configuration *appconfiguration.UseCase
}

func NewCalendarHandler(appointments *appappointment.UseCase, configuration *appconfiguration.UseCase) *CalendarHandler {
	return &CalendarHandler{appointments: appointments, configuration: configuration}
}

func (h *CalendarHandler) ExportICal(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	query := r.URL.Query()
	from := query.Get("from")
	to := query.Get("to")
	if from == "" {
		from = time.Now().AddDate(0, 0, -30).Format("2006-01-02T00:00:00Z07:00")
	}
	if to == "" {
		to = time.Now().AddDate(0, 0, 90).Format("2006-01-02T00:00:00Z07:00")
	}

	appointments, err := h.appointments.ListByRange(tenantID, from, to, "", "", "")
	if err != nil {
		writeInternalError(w, "Could not export the calendar.")
		return
	}

	establishmentName := "Agenda"
	if configuration, err := h.configuration.Get(tenantID); err == nil && configuration.EstablishmentName != "" {
		establishmentName = configuration.EstablishmentName
	}

	content := calendar.BuildICalCalendar(appointments, establishmentName)
	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=\"agenda.ics\"")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(content))
}

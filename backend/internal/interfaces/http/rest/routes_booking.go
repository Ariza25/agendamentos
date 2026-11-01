package rest

import "net/http"

func registerBookingRoutes(mux *http.ServeMux, config RouterConfig) {
	handler := NewBookingHandler(config.BookingUseCase, config.EventBroker)

	mux.HandleFunc("GET /api/v1/public/{slug}/info", handler.Info)
	mux.HandleFunc("GET /api/v1/public/{slug}/services", handler.Services)
	mux.HandleFunc("GET /api/v1/public/{slug}/professionals", handler.Professionals)
	mux.HandleFunc("GET /api/v1/public/{slug}/availability", handler.Availability)
	mux.HandleFunc("POST /api/v1/public/{slug}/appointments", handler.Create)
	mux.HandleFunc("GET /api/v1/public/{slug}/appointments/{code}", handler.GetByManagementCode)
	mux.HandleFunc("PATCH /api/v1/public/{slug}/appointments/{code}/cancel", handler.Cancel)
	mux.HandleFunc("PATCH /api/v1/public/{slug}/appointments/{code}/reschedule", handler.Reschedule)
	mux.HandleFunc("POST /api/v1/public/{slug}/portal/lookup", handler.PortalLookup)
	mux.HandleFunc("POST /api/v1/public/{slug}/portal/appointments/{appointmentID}/cancel", handler.PortalCancel)
	mux.HandleFunc("POST /api/v1/public/{slug}/portal/appointments/{appointmentID}/reschedule", handler.PortalReschedule)
}

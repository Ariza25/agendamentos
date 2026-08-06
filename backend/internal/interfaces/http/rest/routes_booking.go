package rest

import "net/http"

func registerBookingRoutes(mux *http.ServeMux, config RouterConfig) {
	handler := NewBookingHandler(config.BookingUseCase)

	mux.HandleFunc("GET /api/v1/public/{slug}/info", handler.Info)
	mux.HandleFunc("GET /api/v1/public/{slug}/services", handler.Services)
	mux.HandleFunc("GET /api/v1/public/{slug}/professionals", handler.Professionals)
	mux.HandleFunc("GET /api/v1/public/{slug}/availability", handler.Availability)
	mux.HandleFunc("POST /api/v1/public/{slug}/appointments", handler.Create)
}

package rest

import "net/http"

func registerReviewRoutes(mux *http.ServeMux, config RouterConfig) {
	handler := NewReviewHandler(config.ReviewUseCase, config.BookingUseCase)

	mux.HandleFunc("GET /api/v1/admin/reviews", requireAdmin(config.AuthUseCase, handler.List))
	mux.HandleFunc("GET /api/v1/public/{slug}/reviews/{reviewID}", handler.GetPublic)
	mux.HandleFunc("POST /api/v1/public/{slug}/reviews/{reviewID}", handler.SubmitPublic)
}

func registerExportRoutes(mux *http.ServeMux, config RouterConfig) {
	handler := NewExportHandler(
		config.AppointmentUseCase,
		config.CustomerUseCase,
		config.PaymentUseCase,
		config.ProfessionalUseCase,
		config.ServiceUseCase,
		config.ReviewUseCase,
		config.ResourceUseCase,
	)

	mux.HandleFunc("GET /api/v1/admin/export/appointments.csv", requireAdmin(config.AuthUseCase, handler.ExportAppointments))
	mux.HandleFunc("GET /api/v1/admin/export/customers.csv", requireAdmin(config.AuthUseCase, handler.ExportCustomers))
	mux.HandleFunc("GET /api/v1/admin/export/payments.csv", requireAdmin(config.AuthUseCase, handler.ExportPayments))
	mux.HandleFunc("GET /api/v1/admin/export/reviews.csv", requireAdmin(config.AuthUseCase, handler.ExportReviews))
}

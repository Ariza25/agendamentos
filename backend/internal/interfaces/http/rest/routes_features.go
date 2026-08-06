package rest

import "net/http"

func registerNotificationRoutes(mux *http.ServeMux, config RouterConfig) {
	handler := NewNotificationHandler(config.NotificationUseCase)

	mux.HandleFunc("GET /api/v1/admin/notifications", requireAdmin(config.AuthUseCase, handler.List))
	mux.HandleFunc("POST /api/v1/admin/notifications/test", requireAdmin(config.AuthUseCase, handler.SendTest))
}

func registerPaymentRoutes(mux *http.ServeMux, config RouterConfig) {
	handler := NewPaymentHandler(config.PaymentUseCase)

	mux.HandleFunc("GET /api/v1/admin/payments", requireAdmin(config.AuthUseCase, handler.List))
	mux.HandleFunc("POST /api/v1/admin/payments", requireAdmin(config.AuthUseCase, handler.Create))
	mux.HandleFunc("GET /api/v1/admin/payments/{paymentID}", requireAdmin(config.AuthUseCase, handler.Get))
	mux.HandleFunc("POST /api/v1/admin/payments/{paymentID}/mark-paid", requireAdmin(config.AuthUseCase, handler.MarkPaid))
	mux.HandleFunc("POST /api/v1/admin/payments/{paymentID}/refund", requireAdmin(config.AuthUseCase, handler.Refund))
}

func registerWaitingListRoutes(mux *http.ServeMux, config RouterConfig) {
	handler := NewWaitingListHandler(config.WaitingListUseCase, config.BookingUseCase)

	mux.HandleFunc("GET /api/v1/admin/waiting-list", requireAdmin(config.AuthUseCase, handler.List))
	mux.HandleFunc("PATCH /api/v1/admin/waiting-list/{entryID}/status", requireAdmin(config.AuthUseCase, handler.UpdateStatus))
	mux.HandleFunc("POST /api/v1/public/{slug}/waiting-list", handler.PublicSignup)
}

func registerCommerceRoutes(mux *http.ServeMux, config RouterConfig) {
	handler := NewCommerceHandler(config.CommerceUseCase)

	mux.HandleFunc("GET /api/v1/admin/packages", requireAdmin(config.AuthUseCase, handler.ListPackages))
	mux.HandleFunc("POST /api/v1/admin/packages", requireAdmin(config.AuthUseCase, handler.CreatePackage))
	mux.HandleFunc("PUT /api/v1/admin/packages/{packageID}", requireAdmin(config.AuthUseCase, handler.UpdatePackage))
	mux.HandleFunc("PATCH /api/v1/admin/packages/{packageID}/active", requireAdmin(config.AuthUseCase, handler.SetPackageActive))
	mux.HandleFunc("DELETE /api/v1/admin/packages/{packageID}", requireAdmin(config.AuthUseCase, handler.DeletePackage))

	mux.HandleFunc("GET /api/v1/admin/package-credits", requireAdmin(config.AuthUseCase, handler.ListPackageCredits))
	mux.HandleFunc("POST /api/v1/admin/package-credits", requireAdmin(config.AuthUseCase, handler.CreatePackageCredit))

	mux.HandleFunc("GET /api/v1/admin/gift-cards", requireAdmin(config.AuthUseCase, handler.ListGiftCards))
	mux.HandleFunc("POST /api/v1/admin/gift-cards", requireAdmin(config.AuthUseCase, handler.CreateGiftCard))
	mux.HandleFunc("PUT /api/v1/admin/gift-cards/{giftCardID}", requireAdmin(config.AuthUseCase, handler.UpdateGiftCard))
	mux.HandleFunc("DELETE /api/v1/admin/gift-cards/{giftCardID}", requireAdmin(config.AuthUseCase, handler.DeleteGiftCard))

	mux.HandleFunc("GET /api/v1/admin/memberships", requireAdmin(config.AuthUseCase, handler.ListMemberships))
	mux.HandleFunc("POST /api/v1/admin/memberships", requireAdmin(config.AuthUseCase, handler.CreateMembership))
	mux.HandleFunc("PUT /api/v1/admin/memberships/{membershipID}", requireAdmin(config.AuthUseCase, handler.UpdateMembership))
	mux.HandleFunc("PATCH /api/v1/admin/memberships/{membershipID}/active", requireAdmin(config.AuthUseCase, handler.SetMembershipActive))
	mux.HandleFunc("DELETE /api/v1/admin/memberships/{membershipID}", requireAdmin(config.AuthUseCase, handler.DeleteMembership))

	mux.HandleFunc("GET /api/v1/admin/member-plans", requireAdmin(config.AuthUseCase, handler.ListMemberPlans))
	mux.HandleFunc("POST /api/v1/admin/member-plans", requireAdmin(config.AuthUseCase, handler.CreateMemberPlan))
	mux.HandleFunc("PATCH /api/v1/admin/member-plans/{planID}/status", requireAdmin(config.AuthUseCase, handler.UpdateMemberPlanStatus))
}

func registerCalendarRoutes(mux *http.ServeMux, config RouterConfig) {
	handler := NewCalendarHandler(config.AppointmentUseCase, config.ConfigurationUseCase)

	mux.HandleFunc("GET /api/v1/admin/calendar/export.ics", requireAdmin(config.AuthUseCase, handler.ExportICal))
}

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

func registerStockRoutes(mux *http.ServeMux, config RouterConfig) {
	if config.StockUseCase == nil {
		return
	}
	handler := NewStockHandler(config.StockUseCase)

	mux.HandleFunc("GET /api/v1/admin/products", requireAdmin(config.AuthUseCase, handler.ListProducts))
	mux.HandleFunc("POST /api/v1/admin/products", requireAdmin(config.AuthUseCase, handler.CreateProduct))
	mux.HandleFunc("PUT /api/v1/admin/products/{productID}", requireAdmin(config.AuthUseCase, handler.UpdateProduct))
	mux.HandleFunc("DELETE /api/v1/admin/products/{productID}", requireAdmin(config.AuthUseCase, handler.DeleteProduct))

	mux.HandleFunc("GET /api/v1/admin/stock/movements", requireAdmin(config.AuthUseCase, handler.ListMovements))
	mux.HandleFunc("POST /api/v1/admin/stock/movements", requireAdmin(config.AuthUseCase, handler.CreateMovement))

	mux.HandleFunc("GET /api/v1/admin/services/{serviceID}/supplies", requireAdmin(config.AuthUseCase, handler.ListServiceSupplies))
	mux.HandleFunc("PUT /api/v1/admin/services/{serviceID}/supplies", requireAdmin(config.AuthUseCase, handler.SetServiceSupplies))
}

func registerCommissionRoutes(mux *http.ServeMux, config RouterConfig) {
	if config.CommissionUseCase == nil {
		return
	}
	handler := NewCommissionHandler(config.CommissionUseCase)

	mux.HandleFunc("GET /api/v1/admin/commissions/rules", requireAdmin(config.AuthUseCase, handler.ListRules))
	mux.HandleFunc("POST /api/v1/admin/commissions/rules", requireAdmin(config.AuthUseCase, handler.CreateRule))
	mux.HandleFunc("PUT /api/v1/admin/commissions/rules/{ruleID}", requireAdmin(config.AuthUseCase, handler.UpdateRule))
	mux.HandleFunc("DELETE /api/v1/admin/commissions/rules/{ruleID}", requireAdmin(config.AuthUseCase, handler.DeleteRule))

	mux.HandleFunc("GET /api/v1/admin/commissions", requireAdmin(config.AuthUseCase, handler.ListRecords))
	mux.HandleFunc("POST /api/v1/admin/commissions/settle", requireAdmin(config.AuthUseCase, handler.SettleCommissions))
	mux.HandleFunc("POST /api/v1/admin/commissions/payout", requireAdmin(config.AuthUseCase, handler.Payout))
	mux.HandleFunc("GET /api/v1/admin/commissions/payouts", requireAdmin(config.AuthUseCase, handler.ListPayouts))
	mux.HandleFunc("GET /api/v1/admin/commissions/summaries", requireAdmin(config.AuthUseCase, handler.GetSummaries))
}

func registerFinanceRoutes(mux *http.ServeMux, config RouterConfig) {
	if config.FinanceUseCase == nil {
		return
	}
	handler := NewFinanceHandler(config.FinanceUseCase)

	mux.HandleFunc("GET /api/v1/admin/finance/transactions", requireAdmin(config.AuthUseCase, handler.ListTransactions))
	mux.HandleFunc("POST /api/v1/admin/finance/transactions", requireAdmin(config.AuthUseCase, handler.CreateTransaction))
	mux.HandleFunc("PUT /api/v1/admin/finance/transactions/{txID}", requireAdmin(config.AuthUseCase, handler.UpdateTransaction))
	mux.HandleFunc("DELETE /api/v1/admin/finance/transactions/{txID}", requireAdmin(config.AuthUseCase, handler.DeleteTransaction))
	mux.HandleFunc("GET /api/v1/admin/finance/summary", requireAdmin(config.AuthUseCase, handler.GetSummary))
}

func registerCalendarSyncRoutes(mux *http.ServeMux, config RouterConfig) {
	if config.CalendarSyncUseCase == nil {
		return
	}
	handler := NewCalendarSyncHandler(config.CalendarSyncUseCase)

	mux.HandleFunc("GET /api/v1/admin/calendar/connections", requireAdmin(config.AuthUseCase, handler.ListConnections))
	mux.HandleFunc("POST /api/v1/admin/calendar/connections", requireAdmin(config.AuthUseCase, handler.Connect))
	mux.HandleFunc("DELETE /api/v1/admin/calendar/connections/{connID}", requireAdmin(config.AuthUseCase, handler.Disconnect))
}

func registerWebhookRoutes(mux *http.ServeMux, config RouterConfig) {
	if config.PaymentUseCase == nil {
		return
	}
	mpHandler := NewMercadoPagoWebhookHandler(
		config.PaymentUseCase,
		config.AppointmentUseCase,
		config.FinanceUseCase,
		config.EventBroker,
	)
	mux.HandleFunc("POST /api/v1/webhooks/mercadopago", mpHandler.HandleWebhook)
	mux.HandleFunc("GET /api/v1/webhooks/mercadopago", mpHandler.HandleWebhook)
}

func registerOnboardingRoutes(mux *http.ServeMux, config RouterConfig) {
	if config.OnboardingUseCase == nil {
		return
	}
	handler := NewOnboardingHandler(config.OnboardingUseCase)
	mux.HandleFunc("POST /api/v1/auth/signup", handler.Signup)
}

func registerCashRegisterRoutes(mux *http.ServeMux, config RouterConfig) {
	if config.CashRegisterUseCase == nil {
		return
	}
	handler := NewCashRegisterHandler(config.CashRegisterUseCase)

	mux.HandleFunc("GET /api/v1/admin/cash-register/current", requireAdmin(config.AuthUseCase, handler.Current))
	mux.HandleFunc("POST /api/v1/admin/cash-register/open", requireAdmin(config.AuthUseCase, handler.Open))
	mux.HandleFunc("POST /api/v1/admin/cash-register/entry", requireAdmin(config.AuthUseCase, handler.Entry))
	mux.HandleFunc("POST /api/v1/admin/cash-register/close", requireAdmin(config.AuthUseCase, handler.Close))
	mux.HandleFunc("GET /api/v1/admin/cash-register/history", requireAdmin(config.AuthUseCase, handler.History))
}

func registerCheckoutRoutes(mux *http.ServeMux, config RouterConfig) {
	if config.CheckoutUseCase == nil {
		return
	}
	handler := NewCheckoutHandler(config.CheckoutUseCase, config.EventBroker)
	mux.HandleFunc("POST /api/v1/admin/checkout", requireAdmin(config.AuthUseCase, handler.Checkout))
}

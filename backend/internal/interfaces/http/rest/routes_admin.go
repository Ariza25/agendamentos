package rest

import "net/http"

func registerTenantRoutes(mux *http.ServeMux, config RouterConfig) {
	tenantHandler := NewTenantHandler(config.TenantUseCase)

	mux.HandleFunc("GET /api/v1/admin/tenants", requireMasterAdmin(config.AuthUseCase, tenantHandler.List))
	mux.HandleFunc("POST /api/v1/admin/tenants", requireMasterAdmin(config.AuthUseCase, tenantHandler.Create))
	mux.HandleFunc("GET /api/v1/admin/tenants/{tenantID}", requireMasterAdmin(config.AuthUseCase, tenantHandler.Get))
	mux.HandleFunc("PUT /api/v1/admin/tenants/{tenantID}", requireMasterAdmin(config.AuthUseCase, tenantHandler.Update))
	mux.HandleFunc("PATCH /api/v1/admin/tenants/{tenantID}/active", requireMasterAdmin(config.AuthUseCase, tenantHandler.SetActive))
	mux.HandleFunc("DELETE /api/v1/admin/tenants/{tenantID}", requireMasterAdmin(config.AuthUseCase, tenantHandler.Delete))
}

func registerProfessionalRoutes(mux *http.ServeMux, config RouterConfig) {
	handler := NewProfessionalHandler(config.ProfessionalUseCase)

	mux.HandleFunc("GET /api/v1/admin/professionals", requireAdmin(config.AuthUseCase, handler.List))
	mux.HandleFunc("POST /api/v1/admin/professionals", requireAdmin(config.AuthUseCase, handler.Create))
	mux.HandleFunc("GET /api/v1/admin/professionals/{professionalID}", requireAdmin(config.AuthUseCase, handler.Get))
	mux.HandleFunc("PUT /api/v1/admin/professionals/{professionalID}", requireAdmin(config.AuthUseCase, handler.Update))
	mux.HandleFunc("PATCH /api/v1/admin/professionals/{professionalID}/active", requireAdmin(config.AuthUseCase, handler.SetActive))
	mux.HandleFunc("DELETE /api/v1/admin/professionals/{professionalID}", requireAdmin(config.AuthUseCase, handler.Delete))
}

func registerServiceRoutes(mux *http.ServeMux, config RouterConfig) {
	handler := NewServiceHandler(config.ServiceUseCase)

	mux.HandleFunc("GET /api/v1/admin/services", requireAdmin(config.AuthUseCase, handler.List))
	mux.HandleFunc("POST /api/v1/admin/services", requireAdmin(config.AuthUseCase, handler.Create))
	mux.HandleFunc("GET /api/v1/admin/services/{serviceID}", requireAdmin(config.AuthUseCase, handler.Get))
	mux.HandleFunc("PUT /api/v1/admin/services/{serviceID}", requireAdmin(config.AuthUseCase, handler.Update))
	mux.HandleFunc("PATCH /api/v1/admin/services/{serviceID}/active", requireAdmin(config.AuthUseCase, handler.SetActive))
	mux.HandleFunc("DELETE /api/v1/admin/services/{serviceID}", requireAdmin(config.AuthUseCase, handler.Delete))
}

func registerCustomerRoutes(mux *http.ServeMux, config RouterConfig) {
	handler := NewCustomerHandler(config.CustomerUseCase)

	mux.HandleFunc("GET /api/v1/admin/customers", requireAdmin(config.AuthUseCase, handler.List))
	mux.HandleFunc("POST /api/v1/admin/customers", requireAdmin(config.AuthUseCase, handler.Create))
	mux.HandleFunc("GET /api/v1/admin/customers/{customerID}", requireAdmin(config.AuthUseCase, handler.Get))
	mux.HandleFunc("PUT /api/v1/admin/customers/{customerID}", requireAdmin(config.AuthUseCase, handler.Update))
	mux.HandleFunc("DELETE /api/v1/admin/customers/{customerID}", requireAdmin(config.AuthUseCase, handler.Delete))
}

func registerAppointmentRoutes(mux *http.ServeMux, config RouterConfig) {
	handler := NewAppointmentHandler(config.AppointmentUseCase)

	mux.HandleFunc("GET /api/v1/admin/appointments", requireAdmin(config.AuthUseCase, handler.List))
	mux.HandleFunc("POST /api/v1/admin/appointments", requireAdmin(config.AuthUseCase, handler.Create))
	mux.HandleFunc("GET /api/v1/admin/appointments/summary", requireAdmin(config.AuthUseCase, handler.Summary))
	mux.HandleFunc("GET /api/v1/admin/appointments/analytics", requireAdmin(config.AuthUseCase, handler.Analytics))
	mux.HandleFunc("GET /api/v1/admin/appointments/{appointmentID}", requireAdmin(config.AuthUseCase, handler.Get))
	mux.HandleFunc("PUT /api/v1/admin/appointments/{appointmentID}", requireAdmin(config.AuthUseCase, handler.Update))
	mux.HandleFunc("PATCH /api/v1/admin/appointments/{appointmentID}/status", requireAdmin(config.AuthUseCase, handler.UpdateStatus))
	mux.HandleFunc("POST /api/v1/admin/appointments/{appointmentID}/check-in", requireAdmin(config.AuthUseCase, handler.CheckIn))
	mux.HandleFunc("POST /api/v1/admin/appointments/{appointmentID}/no-show", requireAdmin(config.AuthUseCase, handler.MarkNoShow))
	mux.HandleFunc("DELETE /api/v1/admin/appointments/{appointmentID}", requireAdmin(config.AuthUseCase, handler.Delete))
}

func registerConfigurationRoutes(mux *http.ServeMux, config RouterConfig) {
	handler := NewConfigurationHandler(config.ConfigurationUseCase)

	mux.HandleFunc("GET /api/v1/admin/configuration", requireOwnerOrAdmin(config.AuthUseCase, handler.Get))
	mux.HandleFunc("PUT /api/v1/admin/configuration", requireOwnerOrAdmin(config.AuthUseCase, handler.Update))
}

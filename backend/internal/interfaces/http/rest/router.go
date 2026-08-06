package rest

import (
	"context"
	"net/http"

	appappointment "erpagendamentos/backend/internal/application/appointment"
	appauth "erpagendamentos/backend/internal/application/auth"
	appbooking "erpagendamentos/backend/internal/application/booking"
	appconfiguration "erpagendamentos/backend/internal/application/configuration"
	appcustomer "erpagendamentos/backend/internal/application/customer"
	appprofessional "erpagendamentos/backend/internal/application/professional"
	appservice "erpagendamentos/backend/internal/application/service"
	apptenant "erpagendamentos/backend/internal/application/tenant"
)

type Router struct {
	handler http.Handler
}

// RouterConfig is the HTTP composition boundary. Named dependencies prevent
// accidental positional wiring as the application grows.
type RouterConfig struct {
	AuthUseCase           *appauth.UseCase
	AuthMaintenanceConfig appauth.MaintenanceConfig
	TenantUseCase         *apptenant.UseCase
	ProfessionalUseCase   *appprofessional.UseCase
	ServiceUseCase        *appservice.UseCase
	CustomerUseCase       *appcustomer.UseCase
	AppointmentUseCase    *appappointment.UseCase
	ConfigurationUseCase  *appconfiguration.UseCase
	BookingUseCase        *appbooking.UseCase
}

func NewRouter(config RouterConfig) *Router {
	mux := http.NewServeMux()

	startBackgroundWorkers(config)
	registerAuthRoutes(mux, config)
	registerTenantRoutes(mux, config)
	registerProfessionalRoutes(mux, config)
	registerServiceRoutes(mux, config)
	registerCustomerRoutes(mux, config)
	registerAppointmentRoutes(mux, config)
	registerConfigurationRoutes(mux, config)
	registerBookingRoutes(mux, config)

	handler := withObservability(
		withSecurityHeaders(
			withCORS(
				withRateLimit(config.AuthUseCase,
					withRequestBodyLimit(
						withCSRFProtection(
							withPrivilegedStepUp(config.AuthUseCase,
								withProtectedNamespaces(config.AuthUseCase, mux),
							),
						),
					),
				),
			),
		),
	)
	return &Router{handler: handler}
}

func startBackgroundWorkers(config RouterConfig) {
	if config.AuthUseCase != nil {
		appauth.NewMaintenanceWorker(config.AuthUseCase, config.AuthMaintenanceConfig).Start(context.Background())
	}
}

func (r *Router) Handler() http.Handler {
	return r.handler
}

func (r *Router) Close() {
}

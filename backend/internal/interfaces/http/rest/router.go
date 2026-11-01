package rest

import (
	"context"
	"net/http"
	"time"

	appappointment "erpagendamentos/backend/internal/application/appointment"
	appauth "erpagendamentos/backend/internal/application/auth"
	appbooking "erpagendamentos/backend/internal/application/booking"
	appcalendar "erpagendamentos/backend/internal/application/calendar"
	appcashregister "erpagendamentos/backend/internal/application/cashregister"
	appcheckout "erpagendamentos/backend/internal/application/checkout"
	appcommerce "erpagendamentos/backend/internal/application/commerce"
	appcommission "erpagendamentos/backend/internal/application/commission"
	appconfiguration "erpagendamentos/backend/internal/application/configuration"
	appcustomer "erpagendamentos/backend/internal/application/customer"
	appfinance "erpagendamentos/backend/internal/application/finance"
	appnotification "erpagendamentos/backend/internal/application/notification"
	apponboarding "erpagendamentos/backend/internal/application/onboarding"
	apppayment "erpagendamentos/backend/internal/application/payment"
	appprofessional "erpagendamentos/backend/internal/application/professional"
	appresource "erpagendamentos/backend/internal/application/resource"
	appreview "erpagendamentos/backend/internal/application/review"
	appservice "erpagendamentos/backend/internal/application/service"
	appstock "erpagendamentos/backend/internal/application/stock"
	apptenant "erpagendamentos/backend/internal/application/tenant"
	appwaitinglist "erpagendamentos/backend/internal/application/waitinglist"
	configurationdomain "erpagendamentos/backend/internal/domain/configuration"
	customerdomain "erpagendamentos/backend/internal/domain/customer"
	servicedomain "erpagendamentos/backend/internal/domain/service"
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
	ResourceUseCase       *appresource.UseCase
	CustomerUseCase       *appcustomer.UseCase
	AppointmentUseCase    *appappointment.UseCase
	ConfigurationUseCase  *appconfiguration.UseCase
	BookingUseCase        *appbooking.UseCase
	NotificationUseCase   *appnotification.UseCase
	PaymentUseCase        *apppayment.UseCase
	WaitingListUseCase    *appwaitinglist.UseCase
	CommerceUseCase       *appcommerce.UseCase
	ReviewUseCase         *appreview.UseCase
	StockUseCase          *appstock.UseCase
	CommissionUseCase     *appcommission.UseCase
	FinanceUseCase        *appfinance.UseCase
	CalendarSyncUseCase   *appcalendar.UseCase
	OnboardingUseCase     *apponboarding.UseCase
	CashRegisterUseCase   *appcashregister.UseCase
	CheckoutUseCase       *appcheckout.UseCase
	EventBroker           *EventBroker
	CustomerRepository    customerdomain.Repository
	ServiceRepository     servicedomain.Repository
	ConfigurationRepo     configurationdomain.Repository
	ReminderInterval      time.Duration
	ReminderEnabled       bool
}

func NewRouter(config RouterConfig) *Router {
	mux := http.NewServeMux()

	startBackgroundWorkers(config)
	registerAuthRoutes(mux, config)
	registerTenantRoutes(mux, config)
	registerProfessionalRoutes(mux, config)
	registerServiceRoutes(mux, config)
	registerResourceRoutes(mux, config)
	registerCustomerRoutes(mux, config)
	registerAppointmentRoutes(mux, config)
	registerConfigurationRoutes(mux, config)
	registerBookingRoutes(mux, config)
	registerNotificationRoutes(mux, config)
	registerPaymentRoutes(mux, config)
	registerWaitingListRoutes(mux, config)
	registerCommerceRoutes(mux, config)
	registerCalendarRoutes(mux, config)
	registerReviewRoutes(mux, config)
	registerExportRoutes(mux, config)
	registerStockRoutes(mux, config)
	registerCommissionRoutes(mux, config)
	registerFinanceRoutes(mux, config)
	registerCalendarSyncRoutes(mux, config)
	registerWebhookRoutes(mux, config)
	registerOnboardingRoutes(mux, config)
	registerCashRegisterRoutes(mux, config)
	registerCheckoutRoutes(mux, config)

	// Prometheus metrics
	mux.HandleFunc("GET /metrics", HandleMetrics)

	// Real-time SSE event stream
	if config.EventBroker != nil {
		RegisterMetricsBroker(config.EventBroker)
		mux.HandleFunc("GET /api/v1/admin/events", requireAdmin(config.AuthUseCase, config.EventBroker.ServeHTTP))
	}

	widgetHandler := NewWidgetHandler()
	mux.HandleFunc("GET /widget.js", widgetHandler.ServeScript)

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
	if config.ReminderEnabled && config.NotificationUseCase != nil {
		appnotification.NewReminderWorker(
			config.AppointmentUseCase,
			config.NotificationUseCase,
			config.CustomerRepository,
			config.ServiceRepository,
			config.ConfigurationRepo,
			config.ReminderInterval,
		).Start(context.Background())
	}
}

func (r *Router) Handler() http.Handler {
	return r.handler
}

func (r *Router) Close() {
}

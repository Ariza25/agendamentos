package bootstrap

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"erpagendamentos/backend/internal/application/appointment"
	"erpagendamentos/backend/internal/application/audit"
	"erpagendamentos/backend/internal/application/auth"
	"erpagendamentos/backend/internal/application/booking"
	"erpagendamentos/backend/internal/application/calendar"
	"erpagendamentos/backend/internal/application/cashregister"
	"erpagendamentos/backend/internal/application/checkout"
	"erpagendamentos/backend/internal/application/commerce"
	"erpagendamentos/backend/internal/application/commission"
	"erpagendamentos/backend/internal/application/configuration"
	"erpagendamentos/backend/internal/application/customer"
	"erpagendamentos/backend/internal/application/finance"
	"erpagendamentos/backend/internal/application/job"
	"erpagendamentos/backend/internal/application/notification"
	"erpagendamentos/backend/internal/application/onboarding"
	"erpagendamentos/backend/internal/application/outbox"
	"erpagendamentos/backend/internal/application/payment"
	"erpagendamentos/backend/internal/application/professional"
	"erpagendamentos/backend/internal/application/resource"
	"erpagendamentos/backend/internal/application/review"
	"erpagendamentos/backend/internal/application/service"
	"erpagendamentos/backend/internal/application/stock"
	"erpagendamentos/backend/internal/application/tenant"
	"erpagendamentos/backend/internal/application/waitinglist"
	authdomain "erpagendamentos/backend/internal/domain/auth"
	configurationdomain "erpagendamentos/backend/internal/domain/configuration"
	customerdomain "erpagendamentos/backend/internal/domain/customer"
	notificationdomain "erpagendamentos/backend/internal/domain/notification"
	paymentdomain "erpagendamentos/backend/internal/domain/payment"
	servicedomain "erpagendamentos/backend/internal/domain/service"
	tenantdomain "erpagendamentos/backend/internal/domain/tenant"
	notificationinfra "erpagendamentos/backend/internal/infrastructure/notifications"
	paymentinfra "erpagendamentos/backend/internal/infrastructure/payments"
	"erpagendamentos/backend/internal/infrastructure/persistence"
	"erpagendamentos/backend/internal/interfaces/http/rest"
)

type Container struct {
	db                      *sql.DB
	databaseRequired        bool
	authUseCase             *auth.UseCase
	tenantUseCase           *tenant.UseCase
	professionalUseCase     *professional.UseCase
	serviceUseCase          *service.UseCase
	resourceUseCase         *resource.UseCase
	customerUseCase         *customer.UseCase
	appointmentUseCase      *appointment.UseCase
	configurationUseCase    *configuration.UseCase
	bookingUseCase          *booking.UseCase
	notificationUseCase     *notification.UseCase
	paymentUseCase          *payment.UseCase
	waitingListUseCase      *waitinglist.UseCase
	commerceUseCase         *commerce.UseCase
	reviewUseCase           *review.UseCase
	stockUseCase            *stock.UseCase
	commissionUseCase       *commission.UseCase
	financeUseCase          *finance.UseCase
	calendarUseCase         *calendar.UseCase
	onboardingUseCase       *onboarding.UseCase
	cashRegisterUseCase     *cashregister.UseCase
	checkoutUseCase         *checkout.UseCase
	auditUseCase            *audit.UseCase
	outboxUseCase           *outbox.UseCase
	jobQueue                *job.Queue
	eventBroker             *rest.EventBroker
	authMaintenance         auth.MaintenanceConfig
	customerRepository      customerdomain.Repository
	serviceRepository       servicedomain.Repository
	configurationRepository configurationdomain.Repository
	router                  *rest.Router
}

func NewContainer() *Container {
	loadEnvFile(".env")

	repositories, db, databaseRequired := initializeRepositories()
	authRepository := repositories.auth

	ensureDefaultTenant(repositories.tenant)
	authComponents := buildAuthComponents(authRepository)
	seedHomologationData(repositories, authComponents.hasher)

	eventBroker := rest.NewEventBroker()
	stockUseCase := stock.NewUseCase(repositories.stock)
	commissionUseCase := commission.NewUseCase(repositories.commission, repositories.finance, repositories.professional)
	financeUseCase := finance.NewUseCase(repositories.finance)
	calendarUseCase := calendar.NewUseCase(repositories.calendar)
	jobQueue := job.NewQueue(repositories.job)
	jobQueue.Start(2, 5*time.Second)

	appointmentUseCase := appointment.NewUseCase(repositories.appointment, repositories.resource)
	notificationUseCase := notification.NewUseCase(repositories.notification, buildNotificationProvider())
	paymentUseCase := payment.NewUseCase(repositories.payment, buildPaymentProvider())
	onboardingUseCase := onboarding.NewUseCase(
		repositories.tenant,
		authComponents.useCase,
		repositories.configuration,
		repositories.service,
		repositories.professional,
	)

	auditUseCase := audit.NewUseCase(repositories.audit)
	outboxUseCase := outbox.NewUseCase(repositories.outbox, func(tenantID string, topic string, payload map[string]any) error {
		if eventBroker != nil {
			eventBroker.Publish(tenantID, topic, payload)
		}
		return nil
	})
	outboxUseCase.StartWorker(context.Background(), 5*time.Second)

	container := &Container{
		db:                   db,
		databaseRequired:     databaseRequired,
		authUseCase:          authComponents.useCase,
		tenantUseCase:        tenant.NewUseCase(repositories.tenant),
		professionalUseCase:  professional.NewUseCase(repositories.professional),
		serviceUseCase:       service.NewUseCase(repositories.service),
		resourceUseCase:      resource.NewUseCase(repositories.resource),
		customerUseCase:      customer.NewUseCase(repositories.customer),
		appointmentUseCase:   appointmentUseCase,
		configurationUseCase: configuration.NewUseCase(repositories.configuration),
		notificationUseCase:  notificationUseCase,
		paymentUseCase:       paymentUseCase,
		waitingListUseCase:   waitinglist.NewUseCase(repositories.waitingList),
		reviewUseCase:        review.NewUseCase(repositories.review),
		stockUseCase:         stockUseCase,
		commissionUseCase:    commissionUseCase,
		financeUseCase:       financeUseCase,
		calendarUseCase:      calendarUseCase,
		onboardingUseCase:    onboardingUseCase,
		cashRegisterUseCase:  cashregister.NewUseCase(repositories.cashRegister, repositories.finance),
		checkoutUseCase: checkout.NewUseCase(
			repositories.appointment,
			repositories.service,
			repositories.professional,
			repositories.stock,
			repositories.commission,
			repositories.finance,
			repositories.cashRegister,
		).WithTxRunner(persistence.NewPostgresTxRunner(db)).WithAudit(auditUseCase).WithOutbox(outboxUseCase),
		auditUseCase:         auditUseCase,
		outboxUseCase:        outboxUseCase,
		jobQueue:             jobQueue,
		eventBroker:          eventBroker,
		commerceUseCase: commerce.NewUseCase(
			repositories.packages,
			repositories.packageCredits,
			repositories.giftCards,
			repositories.memberships,
			repositories.memberPlans,
		),
		authMaintenance:         authComponents.maintenance,
		customerRepository:      repositories.customer,
		serviceRepository:       repositories.service,
		configurationRepository: repositories.configuration,
	}
	container.bookingUseCase = booking.NewUseCase(
		repositories.tenant,
		repositories.configuration,
		repositories.service,
		repositories.professional,
		repositories.resource,
		repositories.customer,
		appointmentUseCase,
		notificationUseCase,
	)
	return container
}

func buildNotificationProvider() notificationdomain.Provider {
	host := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	if host == "" {
		return notificationinfra.NewLogProvider()
	}
	return notificationinfra.NewEmailProvider(
		host,
		intEnv("SMTP_PORT", 587),
		strings.TrimSpace(os.Getenv("SMTP_USER")),
		os.Getenv("SMTP_PASS"),
		strings.TrimSpace(os.Getenv("SMTP_FROM")),
	)
}

func buildPaymentProvider() paymentdomain.Provider {
	if mpToken := strings.TrimSpace(os.Getenv("MERCADOPAGO_ACCESS_TOKEN")); mpToken != "" {
		webhookURL := strings.TrimSpace(os.Getenv("MERCADOPAGO_WEBHOOK_URL"))
		return paymentinfra.NewMercadoPagoProvider(mpToken, webhookURL)
	}
	if secretKey := strings.TrimSpace(os.Getenv("STRIPE_SECRET_KEY")); secretKey != "" {
		return paymentinfra.NewStripeProvider(secretKey)
	}
	return paymentinfra.NewMockProvider()
}

func isProductionEnv() bool {
	value := deploymentEnvironment()
	return value == "production" || value == "prod"
}

func deploymentEnvironment() string {
	value := strings.ToLower(strings.TrimSpace(os.Getenv("ERPAGENDAMENTOS_ENV")))
	if value == "" {
		value = strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	}
	return value
}

func firstEnv(keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value
		}
	}
	return ""
}

func boolEnv(key string, defaultValue bool) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	if value == "" {
		return defaultValue
	}
	switch value {
	case "1", "true", "yes", "y", "on", "enabled":
		return true
	case "0", "false", "no", "n", "off", "disabled":
		return false
	default:
		log.Printf(`{"level":"WARN","message":"invalid boolean env, using default","env":"%s","value":"%s","default":%t}`, key, value, defaultValue)
		return defaultValue
	}
}

func intEnv(key string, defaultValue int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		log.Printf(`{"level":"WARN","message":"invalid integer env, using default","env":"%s","value":"%s","default":%d}`, key, value, defaultValue)
		return defaultValue
	}
	return parsed
}

func durationSecondsEnv(key string, defaultValue time.Duration) time.Duration {
	seconds := intEnv(key, int(defaultValue.Seconds()))
	return time.Duration(seconds) * time.Second
}

func (c *Container) Router() *rest.Router {
	if c.router != nil {
		return c.router
	}
	c.router = rest.NewRouter(rest.RouterConfig{
		AuthUseCase:           c.authUseCase,
		AuthMaintenanceConfig: c.authMaintenance,
		TenantUseCase:         c.tenantUseCase,
		ProfessionalUseCase:   c.professionalUseCase,
		ServiceUseCase:        c.serviceUseCase,
		ResourceUseCase:       c.resourceUseCase,
		CustomerUseCase:       c.customerUseCase,
		AppointmentUseCase:    c.appointmentUseCase,
		ConfigurationUseCase:  c.configurationUseCase,
		BookingUseCase:        c.bookingUseCase,
		NotificationUseCase:   c.notificationUseCase,
		PaymentUseCase:        c.paymentUseCase,
		WaitingListUseCase:    c.waitingListUseCase,
		CommerceUseCase:       c.commerceUseCase,
		ReviewUseCase:         c.reviewUseCase,
		StockUseCase:          c.stockUseCase,
		CommissionUseCase:     c.commissionUseCase,
		FinanceUseCase:        c.financeUseCase,
		CalendarSyncUseCase:   c.calendarUseCase,
		OnboardingUseCase:     c.onboardingUseCase,
		CashRegisterUseCase:   c.cashRegisterUseCase,
		CheckoutUseCase:       c.checkoutUseCase,
		AuditUseCase:          c.auditUseCase,
		OutboxUseCase:         c.outboxUseCase,
		ReadinessChecker:      c.ReadinessCheck,
		DBStatsGetter:         c.DBStats,
		EventBroker:           c.eventBroker,
		CustomerRepository:    c.customerRepository,
		ServiceRepository:     c.serviceRepository,
		ConfigurationRepo:     c.configurationRepository,
		ReminderInterval:      durationSecondsEnv("REMINDER_INTERVAL_SECONDS", 30*time.Minute),
		ReminderEnabled:       boolEnv("REMINDER_WORKER_ENABLED", true),
	})
	return c.router
}

func (c *Container) ReadinessCheck() error {
	if c.databaseRequired {
		if c.db == nil {
			return fmt.Errorf("database is not connected")
		}
		if err := c.db.Ping(); err != nil {
			return fmt.Errorf("database ping failed: %w", err)
		}
	}
	return nil
}

func (c *Container) DBStats() map[string]any {
	if c.db == nil {
		return map[string]any{"driver": "memory", "connected": true}
	}
	stats := c.db.Stats()
	return map[string]any{
		"driver":             "postgres",
		"connected":          true,
		"openConnections":    stats.OpenConnections,
		"inUse":              stats.InUse,
		"idle":               stats.Idle,
		"waitCount":          stats.WaitCount,
		"maxOpenConnections": stats.MaxOpenConnections,
	}
}

func (c *Container) Close() error {
	var closeErrors []error
	if c.jobQueue != nil {
		c.jobQueue.Stop()
	}
	if c.router != nil {
		c.router.Close()
	}
	if c.db != nil {
		if err := c.db.Close(); err != nil {
			closeErrors = append(closeErrors, err)
		}
	}
	return errors.Join(closeErrors...)
}

func ensureDefaultTenant(repository tenantdomain.Repository) {
	_, err := repository.GetByID("default")
	if err == nil {
		return
	}
	if _, err := repository.Create(tenantdomain.Tenant{
		ID:     "default",
		Name:   "Estabelecimento Padrão",
		Slug:   "default",
		Active: true,
	}); err != nil {
		log.Printf(`{"level":"WARN","message":"default tenant seed skipped","error":"%v"}`, err)
	}
}

func seedAdminAccount(repository authdomain.Repository, hasher auth.PasswordHasher) error {
	email := strings.ToLower(strings.TrimSpace(os.Getenv("ADMIN_EMAIL")))
	password := strings.TrimSpace(os.Getenv("ADMIN_PASSWORD"))
	if email == "" || password == "" {
		return nil
	}

	if existing, err := repository.GetStaffAccountByEmail(email); err == nil {
		if existing.Role != "admin" {
			updated, updateErr := repository.UpdateRole(existing.ID, "admin")
			if updateErr != nil {
				return fmt.Errorf("reconcile admin role: %w", updateErr)
			}
			existing = updated
		}
		if hasher.Compare(existing.PasswordHash, password) == nil {
			return nil
		}
		passwordHash, hashErr := hasher.Hash(password)
		if hashErr != nil {
			return fmt.Errorf("hash admin password: %w", hashErr)
		}
		if updateErr := repository.UpdatePassword(existing.ID, passwordHash); updateErr != nil {
			return fmt.Errorf("reconcile admin password: %w", updateErr)
		}
		return nil
	}

	name := strings.TrimSpace(os.Getenv("ADMIN_NAME"))
	if name == "" {
		name = "Administrador"
	}

	passwordHash, err := hasher.Hash(password)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}

	_, err = repository.CreateAccount(authdomain.Account{
		Name:         name,
		Email:        email,
		Phone:        strings.TrimSpace(os.Getenv("ADMIN_PHONE")),
		PasswordHash: passwordHash,
		Role:         "admin",
		TenantID:     "default",
	})
	if err != nil {
		return fmt.Errorf("create admin account: %w", err)
	}

	return nil
}

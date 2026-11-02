# Graph Report - erpagendamentos  (2026-11-02)

## Corpus Check
- 282 files · ~134,479 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2713 nodes · 8114 edges · 116 communities (103 shown, 13 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 687 edges (avg confidence: 0.79)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9de11490`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- adminFetch
- schema.sql
- AdminRoutes.tsx
- Payment
- net/http.ResponseWriter
- go_pkg_testing
- go_pkg_net_http
- PaymentsPage.tsx
- adminAuth.ts
- writeAPIError
- ref_react
- adminApi.ts
- Review
- UseCase
- RouterConfig
- Tenant
- landing/src/App.tsx
- Service
- go_pkg_time
- auth/usecase.go
- auth_cookie.go
- database/sql.DB
- resource_test.go
- Customer
- CommissionRepository
- PackagesPage.tsx
- Resource
- CashRegister
- Professional
- Queue
- Entry
- devDependencies
- devDependencies
- UseCase
- writeJSON
- gencert.go
- UseCase
- withIdempotency
- devDependencies
- ManagePage.tsx
- MemberPlan
- writeNotFoundError
- Configuration
- NewContainer
- BookingPage.tsx
- useToast
- InventoryPage
- AuthRepository
- UseCase
- net/http.Request
- publicFetch
- memoryRepositories
- Membership
- Package
- UseCase
- UseCase
- PackageCredit
- GiftCard
- UseCase
- auth_handler.go
- AuditLog
- go_pkg_errors
- EventBroker
- rate_limit.go
- Account
- Appointment
- writeInvalidPayload
- bookingApi.ts
- templates.go
- Container
- compilerOptions
- AppointmentRepository
- compilerOptions
- compilerOptions
- OutboxRepository
- UseCase
- totp.go
- writeResourceError
- testing.T
- auth/entities.go
- ERP Agendamentos
- NewBcryptHasher
- UseCase
- email_provider.go
- widget_handler.go
- booking/src/lib/apiError.ts
- time.Duration
- ExportHandler
- FinancePage.tsx
- .UpdateAccountRole
- observability.go
- admin/vite.config.ts
- appointment/repository.go
- buildAuthComponents
- booking/usecase.go
- TenantsPage.tsx
- go_pkg_erpagendamentos_backend_internal_domain_appointment
- booking/src/App.tsx
- graphify
- rules/graphify.md
- workflows/graphify.md
- generate-certs.sh
- backup_postgres.sh
- certbot_ssl.sh
- catalog.finance_transactions
- ReminderWorker
- go_pkg_path_filepath
- erpagendamentos/backend
- writeServiceError
- PortalPage.tsx
- go_pkg_context
- repositories.go
- ref_react_dom

## God Nodes (most connected - your core abstractions)
1. `writeJSON()` - 135 edges
2. `resolveRequestTenantID()` - 109 edges
3. `adminFetch()` - 78 edges
4. `useToast()` - 74 edges
5. `writeInvalidPayload()` - 72 edges
6. `UseCase` - 67 edges
7. `invalidateAdminGetCache()` - 65 edges
8. `writeInternalError()` - 58 edges
9. `RouterConfig` - 56 edges
10. `writeAPIError()` - 46 edges

## Surprising Connections (you probably didn't know these)
- `NewContainer()` --calls--> `loadEnvFile()`  [INFERRED]
  backend/internal/bootstrap/container.go → backend/internal/bootstrap/env.go
- `main()` --calls--> `newLazyAppHandler()`  [INFERRED]
  backend/cmd/api/main.go → backend/cmd/api/lazy_handler.go
- `newResourceUseCase()` --calls--> `NewUseCase()`  [INFERRED]
  backend/internal/application/appointment/resource_test.go → backend/internal/application/appointment/usecase.go
- `newTestUseCase()` --calls--> `NewUseCase()`  [INFERRED]
  backend/internal/application/booking/usecase_test.go → backend/internal/application/booking/usecase.go
- `TestManagementCodeNormalization()` --calls--> `hashManagementCode()`  [INFERRED]
  backend/internal/application/booking/usecase_test.go → backend/internal/application/booking/usecase.go

## Import Cycles
- None detected.

## Communities (116 total, 13 thin omitted)

### Community 0 - "adminFetch"
Cohesion: 0.05
Nodes (76): closeCashRegister(), connectCalendar(), createCashRegisterEntry(), createCommissionRule(), createFinancialTransaction(), createMemberPlan(), createMembership(), createProduct() (+68 more)

### Community 1 - "schema.sql"
Cohesion: 0.12
Nodes (35): catalog.appointments, catalog.audit_logs, catalog.auth_rate_limits, catalog.auth_sessions, catalog.calendar_connections, catalog.cash_register_entries, catalog.cash_registers, catalog.commission_payouts (+27 more)

### Community 2 - "AdminRoutes.tsx"
Cohesion: 0.06
Nodes (74): AdminRoutesProps, AgendaPage, AppointmentsPage, CommissionsPage, ConfigurationPage, CustomersPage, GiftCardsPage, InventoryPage (+66 more)

### Community 3 - "Payment"
Cohesion: 0.06
Nodes (26): UseCase, NewUseCase(), normalizeStatus(), buildPaymentProvider(), IntentResult, Payment, NewMercadoPagoProvider(), NewMockProvider() (+18 more)

### Community 4 - "net/http.ResponseWriter"
Cohesion: 0.09
Nodes (13): resolveRequestTenantID(), writeConfigurationError(), writeNotificationError(), writePaymentError(), writeInternalError(), writeReviewError(), net/http.ResponseWriter, CalendarSyncHandler (+5 more)

### Community 5 - "go_pkg_testing"
Cohesion: 0.16
Nodes (6): go_pkg_erpagendamentos_backend_internal_application_calendar, go_pkg_erpagendamentos_backend_internal_domain_calendar, go_pkg_erpagendamentos_backend_internal_domain_professional, go_pkg_erpagendamentos_backend_internal_domain_stock, go_pkg_erpagendamentos_backend_internal_infrastructure_persistence_memory, go_pkg_testing

### Community 6 - "go_pkg_net_http"
Cohesion: 0.10
Nodes (26): deploymentEnvironment(), isProductionEnv(), go_pkg_encoding_json, go_pkg_erpagendamentos_backend_internal_application_booking, go_pkg_erpagendamentos_backend_internal_application_cashregister, go_pkg_erpagendamentos_backend_internal_application_checkout, go_pkg_erpagendamentos_backend_internal_application_commerce, go_pkg_erpagendamentos_backend_internal_application_commission (+18 more)

### Community 7 - "PaymentsPage.tsx"
Cohesion: 0.07
Nodes (37): PaymentsPage, createPayment(), downloadAdminCsv(), fetchNotifications(), fetchPayments(), fetchReviews(), fetchWaitingList(), formatDateTime() (+29 more)

### Community 8 - "adminAuth.ts"
Cohesion: 0.06
Nodes (50): App(), handleLogout(), AppProviders(), AppProvidersProps, ResetPasswordPage(), handleSubmit(), AppRouter(), admin_src_index (+42 more)

### Community 9 - "writeAPIError"
Cohesion: 0.11
Nodes (26): bearerAuthAllowed(), hasRole(), isAdminRole(), requestWithAuthClaims(), requireMasterAdmin(), requireOwnerOrAdmin(), requirePrivilegedAdmin(), requireRoles() (+18 more)

### Community 10 - "ref_react"
Cohesion: 0.06
Nodes (59): AnalyticsPage, AppointmentFormModal(), computeEndsAt(), handleSubmit(), AppointmentFormModalProps, STATUS_OPTIONS, CheckoutComandaModal(), handleSubmit() (+51 more)

### Community 11 - "adminApi.ts"
Cohesion: 0.07
Nodes (60): AppointmentFilters, fetchCommissionPayouts(), fetchCommissionRecords(), fetchCommissionRules(), fetchCommissionSummaries(), fetchCustomers(), fetchProfessionals(), fetchResources() (+52 more)

### Community 12 - "Review"
Cohesion: 0.14
Nodes (12): UseCase, NewUseCase(), Review, NewReviewRepository(), scanReviewRow(), scanReviewRows(), NewReviewHandler(), ReviewRepository (+4 more)

### Community 13 - "UseCase"
Cohesion: 0.09
Nodes (7): UseCase, Product, ServiceSupply, StockMovement, NewStockHandler(), StockRepository, StockRepository

### Community 14 - "RouterConfig"
Cohesion: 0.21
Nodes (29): requireAdmin(), NewBookingHandler(), NewRouter(), registerAppointmentRoutes(), registerConfigurationRoutes(), registerCustomerRoutes(), registerProfessionalRoutes(), registerResourceRoutes() (+21 more)

### Community 15 - "Tenant"
Cohesion: 0.10
Nodes (13): UseCase, Tenant, UpdateSubscriptionInput, defaultModules(), isTenantIdentityConflict(), defaultModules(), scanTenantRow(), scanTenantRows() (+5 more)

### Community 16 - "landing/src/App.tsx"
Cohesion: 0.10
Nodes (22): App(), CtaSection(), FEATURES, FeaturesSection(), Footer(), Header(), NAV_ITEMS, HeroSection() (+14 more)

### Community 17 - "Service"
Cohesion: 0.11
Nodes (8): UseCase, Service, scanServiceRow(), scanServiceRows(), NewServiceHandler(), ServiceRepository, ServiceRepository, serviceScanner

### Community 18 - "go_pkg_time"
Cohesion: 0.13
Nodes (11): go_pkg_database_sql, go_pkg_erpagendamentos_backend_internal_domain_job, go_pkg_erpagendamentos_backend_internal_domain_resource, go_pkg_erpagendamentos_backend_internal_domain_review, go_pkg_erpagendamentos_backend_internal_domain_waitinglist, go_pkg_fmt, go_pkg_github_com_google_uuid, go_pkg_sort (+3 more)

### Community 19 - "auth/usecase.go"
Cohesion: 0.23
Nodes (9): PasswordResetRequestInput, ValidationError, hashResetToken(), isMFAEligibleRole(), NewValidationError(), passwordResetCodeChallengeKey(), passwordResetCodeToken(), validatePassword() (+1 more)

### Community 20 - "auth_cookie.go"
Cohesion: 0.10
Nodes (31): lazyAppHandler, newLazyAppHandler(), writeLazyJSON(), isPrivilegedAdminRole(), authCookieDomain(), authCookiePartitioned(), authCookieSameSite(), authCookieSecure() (+23 more)

### Community 21 - "database/sql.DB"
Cohesion: 0.11
Nodes (34): initializeRepositories(), postgresRepositories(), NewAppointmentRepository(), NewAuditRepository(), NewAuthRepository(), NewCalendarRepository(), NewCashRegisterRepository(), NewGiftCardRepository() (+26 more)

### Community 22 - "resource_test.go"
Cohesion: 0.54
Nodes (7): UseCase, newAppointment(), newResourceUseCase(), TestCreateAutoAssignsAvailableResource(), TestCreateRejectsBusyExplicitResource(), TestCreateRejectsUnknownResource(), TestUpdateAutoReassignsFreedResource()

### Community 23 - "Customer"
Cohesion: 0.15
Nodes (8): Customer, NewCustomerRepository(), onlyDigits(), scanCustomerRow(), scanCustomerRows(), CustomerRepository, CustomerRepository, customerScanner

### Community 24 - "CommissionRepository"
Cohesion: 0.08
Nodes (7): CommissionRecord, CommissionRule, Payout, Summary, PayoutInput, CommissionRepository, CommissionRepository

### Community 25 - "PackagesPage.tsx"
Cohesion: 0.11
Nodes (27): PackagesPage, createCommissionPayout(), createPackage(), createPackageCredit(), fetchGiftCards(), fetchPackageCredits(), fetchPackages(), updatePackage() (+19 more)

### Community 26 - "Resource"
Cohesion: 0.13
Nodes (8): UseCase, Resource, NewResourceRepository(), scanResourceRow(), NewResourceHandler(), ResourceRepository, ResourceRepository, resourceScanner

### Community 27 - "CashRegister"
Cohesion: 0.14
Nodes (5): CashRegister, CashRegisterEntry, ShiftSummary, CashRegisterRepository, CashRegisterRepository

### Community 28 - "Professional"
Cohesion: 0.18
Nodes (6): Professional, scanProfessionalRow(), scanProfessionalRows(), ProfessionalRepository, ProfessionalRepository, professionalScanner

### Community 29 - "Queue"
Cohesion: 0.10
Nodes (8): Queue, NewQueue(), Job, Repository, sync.WaitGroup, HandlerFunc, JobRepository, JobRepository

### Community 30 - "Entry"
Cohesion: 0.12
Nodes (12): UseCase, NewUseCase(), Entry, NewWaitingListRepository(), scanWaitingRow(), scanWaitingRows(), NewWaitingListHandler(), WaitingListRepository (+4 more)

### Community 31 - "devDependencies"
Cohesion: 0.06
Nodes (34): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+26 more)

### Community 32 - "devDependencies"
Cohesion: 0.06
Nodes (34): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+26 more)

### Community 33 - "UseCase"
Cohesion: 0.25
Nodes (7): UseCase, UseCase, NewCalendarHandler(), NewConfigurationHandler(), NewMercadoPagoWebhookHandler(), CalendarHandler, MercadoPagoWebhookHandler

### Community 34 - "writeJSON"
Cohesion: 0.18
Nodes (5): writeBookingError(), health(), writeJSON(), BookingHandler, CashRegisterHandler

### Community 35 - "gencert.go"
Cohesion: 0.11
Nodes (14): NewSecretCipher(), go_pkg_crypto_aes, go_pkg_crypto_cipher, go_pkg_crypto_rand, go_pkg_crypto_rsa, go_pkg_crypto_sha256, go_pkg_crypto_x509, go_pkg_crypto_x509_pkix (+6 more)

### Community 36 - "UseCase"
Cohesion: 0.12
Nodes (8): UseCase, NewUseCase(), Connection, Repository, SyncEvent, NewCalendarSyncHandler(), CalendarRepository, CalendarRepository

### Community 37 - "withIdempotency"
Cohesion: 0.16
Nodes (11): NewIdempotencyStore(), TestIdempotencyBlocksInFlightConcurrentRequests(), TestIdempotencyDoesNotCacheServerFailures(), TestIdempotencyIgnoredOnGetRequests(), TestIdempotencyReplaysCompletedResponse(), withIdempotency(), bytes.Buffer, net/http.Header (+3 more)

### Community 38 - "devDependencies"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, devDependencies, tailwindcss, @tailwindcss/vite, @types/node (+24 more)

### Community 39 - "ManagePage.tsx"
Cohesion: 0.14
Nodes (21): Service, ManageStep, Button, ButtonProps, ButtonSize, ButtonVariant, sizeClasses, variantClasses (+13 more)

### Community 40 - "MemberPlan"
Cohesion: 0.14
Nodes (7): MemberPlan, NewMemberPlanRepository(), nullIfEmpty(), scanMemberPlans(), database/sql.Rows, MemberPlanRepository, MemberPlanRepository

### Community 41 - "writeNotFoundError"
Cohesion: 0.14
Nodes (7): writeCustomerError(), writeProfessionalError(), writeNotFoundError(), writeTenantError(), CustomerHandler, ProfessionalHandler, TenantHandler

### Community 42 - "Configuration"
Cohesion: 0.18
Nodes (11): workingHoursFor(), Configuration, cloneConfiguration(), defaultMemoryConfiguration(), NewConfigurationRepository(), normalizeConfigurationTenantID(), normalizePostgresTenantID(), ClosedDay (+3 more)

### Community 43 - "NewContainer"
Cohesion: 0.10
Nodes (38): Repository, NewUseCase(), NewUseCase(), UseCase, NewUseCase(), UseCase, NewUseCase(), UseCase (+30 more)

### Community 44 - "BookingPage.tsx"
Cohesion: 0.13
Nodes (18): formatDateLong(), formatDateTime(), formatMoney(), formatTime(), fromDateKey(), SHORT_DAY_LABELS, toDateKey(), Confirmation() (+10 more)

### Community 45 - "useToast"
Cohesion: 0.07
Nodes (57): AdminRoutes(), EstablishmentPage, AdminShell(), handleCopyBookingUrl(), AdminShellProps, AdminNavigationItem, navigationGroups, navigationItems (+49 more)

### Community 46 - "InventoryPage"
Cohesion: 0.31
Nodes (9): fetchServiceSupplies(), fetchStockMovements(), setServiceSupplies(), InventoryPage(), handleAddSupply(), handleDeleteProduct(), handleRemoveSupply(), handleSaveMovement() (+1 more)

### Community 47 - "AuthRepository"
Cohesion: 0.10
Nodes (7): DeviceSession, PasswordResetToken, cloneAuthAccount(), resetTokenExpired(), sessionExpired(), AuthRepository, rateLimitEntry

### Community 48 - "UseCase"
Cohesion: 0.19
Nodes (9): appointmentDurationMinutes(), formatAppointmentForTemplate(), generateManagementCode(), UseCase, hashManagementCode(), parseRFC3339(), PublicBookingOrigin(), BookingInput (+1 more)

### Community 49 - "net/http.Request"
Cohesion: 0.15
Nodes (7): writeAppointmentError(), resolveActorTenantID(), AuthHandler, newAuthMFAHandler(), HandleMetrics(), net/http.Request, AppointmentHandler

### Community 50 - "publicFetch"
Cohesion: 0.11
Nodes (19): cancelPublicAppointment(), createPublicBooking(), fetchPublicAppointment(), fetchPublicInfo(), fetchPublicProfessionals(), fetchPublicServices(), publicFetch(), reschedulePublicAppointment() (+11 more)

### Community 51 - "memoryRepositories"
Cohesion: 0.16
Nodes (20): TestCalendarUseCase(), TestCashRegisterFlow(), TestCheckoutFlow(), TestCommissionUseCase(), TestFinanceUseCase(), TestJobQueue(), TestOnboardingSignup(), TestStockUseCase() (+12 more)

### Community 52 - "Membership"
Cohesion: 0.14
Nodes (4): Membership, NewMembershipRepository(), MembershipRepository, MembershipRepository

### Community 53 - "Package"
Cohesion: 0.14
Nodes (4): Package, NewPackageRepository(), PackageRepository, PackageRepository

### Community 54 - "UseCase"
Cohesion: 0.11
Nodes (8): UseCase, Summary, Transaction, NewFinanceHandler(), CategoryTotal, CategoryTotal, FinanceRepository, FinanceRepository

### Community 55 - "UseCase"
Cohesion: 0.17
Nodes (4): RegisterInput, SecretCipher, Session, UseCase

### Community 56 - "PackageCredit"
Cohesion: 0.18
Nodes (5): PackageCredit, NewPackageCreditRepository(), scanPackageCredits(), PackageCreditRepository, PackageCreditRepository

### Community 57 - "GiftCard"
Cohesion: 0.14
Nodes (5): generateGiftCardCode(), GiftCard, NewGiftCardRepository(), GiftCardRepository, GiftCardRepository

### Community 58 - "UseCase"
Cohesion: 0.17
Nodes (10): UseCase, NewUseCase(), Log, NewNotificationRepository(), NewNotificationHandler(), NotificationRepository, Provider, Repository (+2 more)

### Community 59 - "auth_handler.go"
Cohesion: 0.17
Nodes (11): browserSessionPayload(), buildPasswordResetURL(), exposeAuthTokenForTests(), NewAuthHandler(), secureTemporaryPassword(), writeAdminAuthError(), AuthHandler, logError() (+3 more)

### Community 60 - "AuditLog"
Cohesion: 0.15
Nodes (12): AuditLog, Repository, UseCase, NewUseCase(), TestAuditLogRecordingAndListing(), TestOutboxPublishAndDispatch(), NewAuditRepository(), NewOutboxRepository() (+4 more)

### Community 61 - "go_pkg_errors"
Cohesion: 0.16
Nodes (7): CheckoutInput, CheckoutResult, ProductItem, go_pkg_erpagendamentos_backend_internal_domain_cashregister, go_pkg_erpagendamentos_backend_internal_domain_commission, go_pkg_erpagendamentos_backend_internal_domain_finance, go_pkg_errors

### Community 62 - "EventBroker"
Cohesion: 0.23
Nodes (8): NewCheckoutHandler(), EventBroker, NewEventBroker(), RegisterMetricsBroker(), sync.RWMutex, CheckoutHandler, Event, MetricsRegistry

### Community 63 - "rate_limit.go"
Cohesion: 0.15
Nodes (21): buildCompoundIdempotencyKey(), forwardedClientIP(), isSensitiveRateLimitedRoute(), isTrustedProxy(), newMemoryRateLimiter(), normalizedClientID(), persistentRateLimitConfig(), rateLimitRouteKey() (+13 more)

### Community 64 - "Account"
Cohesion: 0.15
Nodes (5): Account, UpdateProfileInput, updateAuthPasswordAndVersion(), database/sql.Tx, AuthRepository

### Community 65 - "Appointment"
Cohesion: 0.12
Nodes (7): RecurrenceInput, parseRFC3339(), validStatusChange(), Appointment, parseAppointmentBound(), parseAppointmentInterval(), AppointmentRepository

### Community 66 - "writeInvalidPayload"
Cohesion: 0.19
Nodes (5): writeCommerceError(), writeBadRequestError(), writeInvalidPayload(), writeWaitingListError(), CommerceHandler

### Community 67 - "bookingApi.ts"
Cohesion: 0.16
Nodes (14): backendOrigin, buildApiUrl(), AvailabilityResponse, BookingResult, fetchPublicReview(), Professional, PublicAppointment, PublicConfiguration (+6 more)

### Community 68 - "templates.go"
Cohesion: 0.29
Nodes (16): greeting(), RenderCancellation(), RenderConfirmation(), RenderReminder(), RenderReschedule(), sampleData(), TestGreetingWithoutName(), TestRenderCancellation() (+8 more)

### Community 69 - "Container"
Cohesion: 0.13
Nodes (8): UseCase, NewUseCase(), UseCase, Container, NewCustomerHandler(), NewExportHandler(), NewProfessionalHandler(), Repository

### Community 70 - "compilerOptions"
Cohesion: 0.09
Nodes (21): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+13 more)

### Community 71 - "AppointmentRepository"
Cohesion: 0.20
Nodes (5): appointmentSelectColumns(), scanAppointmentRow(), scanAppointmentRows(), AppointmentRepository, appointmentScanner

### Community 72 - "compilerOptions"
Cohesion: 0.09
Nodes (21): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+13 more)

### Community 73 - "compilerOptions"
Cohesion: 0.09
Nodes (21): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+13 more)

### Community 74 - "OutboxRepository"
Cohesion: 0.16
Nodes (4): quoteLiteral(), OutboxRepository, OutboxEvent, OutboxRepository

### Community 75 - "UseCase"
Cohesion: 0.24
Nodes (7): UseCase, NewUseCase(), NewSystemHandler(), go_pkg_runtime, EventDispatcher, Repository, SystemHandler

### Community 76 - "totp.go"
Cohesion: 0.20
Nodes (8): generateTOTPSecret(), totpCode(), totpURI(), validateTOTP(), go_pkg_crypto_sha1, go_pkg_encoding_base32, go_pkg_encoding_binary, go_pkg_net_url

### Community 78 - "testing.T"
Cohesion: 0.35
Nodes (11): UseCase, newTestUseCase(), TestAvailableSlotsRespectsClosedDays(), TestCancelPublic(), TestCreateBookingGeneratesManagementCode(), TestManagementCodeNormalization(), TestPortalCancelValidatesOwnership(), TestPortalLookup() (+3 more)

### Community 79 - "auth/entities.go"
Cohesion: 0.20
Nodes (6): ChangePasswordInput, LoginInput, MFAStatus, PasswordResetCodeConfirmInput, PasswordResetConfirmInput, AuthState

### Community 80 - "ERP Agendamentos"
Cohesion: 0.14
Nodes (13): Banco de Dados & Schema, CI, Endpoints principais, ERP Agendamentos, Estrutura, Pré-requisitos, Rodando com Docker (produção local), Rodando em desenvolvimento (sem banco) (+5 more)

### Community 81 - "NewBcryptHasher"
Cohesion: 0.40
Nodes (3): NewBcryptHasher(), go_pkg_golang_org_x_crypto_bcrypt, BcryptHasher

### Community 82 - "UseCase"
Cohesion: 0.27
Nodes (8): UseCase, NewUseCase(), NewCommerceHandler(), GiftCardRepository, MemberPlanRepository, MembershipRepository, PackageCreditRepository, PackageRepository

### Community 83 - "email_provider.go"
Cohesion: 0.21
Nodes (8): buildNotificationProvider(), intEnv(), logFallback(), NewEmailProvider(), NewLogProvider(), go_pkg_net_smtp, EmailProvider, LogProvider

### Community 84 - "widget_handler.go"
Cohesion: 0.60
Nodes (3): NewWidgetHandler(), resolveWidgetBookingOrigin(), WidgetHandler

### Community 85 - "booking/src/lib/apiError.ts"
Cohesion: 0.27
Nodes (8): ApiError, ApiErrorPayload, buildApiError(), parseApiResponse(), API_ERROR_MESSAGES, FIELD_ERROR_MESSAGES, getApiErrorMessage(), translateFieldErrors()

### Community 86 - "time.Duration"
Cohesion: 0.14
Nodes (10): isProductionAuthEnvironment(), Claims, CleanupSummary, decodeSegment(), encodeSegment(), NewJWTManager(), go_pkg_crypto_hmac, time.Duration (+2 more)

### Community 88 - "FinancePage.tsx"
Cohesion: 0.15
Nodes (21): FinancePage, fetchCashRegisterHistory(), fetchCurrentCashRegister(), fetchFinancialSummary(), fetchFinancialTransactions(), ActiveTab, FinancePage(), handleCloseShift() (+13 more)

### Community 89 - ".UpdateAccountRole"
Cohesion: 0.39
Nodes (4): canManageTenantAccount(), isPrivilegedAdminRole(), isStaffRole(), normalizeTenantScope()

### Community 91 - "observability.go"
Cohesion: 0.15
Nodes (14): isHopByHopHeader(), RecordRequestMetric(), isSensitiveLogKey(), redactSensitiveText(), resolveRequestID(), sanitizeLogString(), sanitizeLogValue(), truncateLogString() (+6 more)

### Community 92 - "admin/vite.config.ts"
Cohesion: 0.60
Nodes (3): ref_tailwindcss_vite, ref_vite, ref_vitejs_plugin_react

### Community 93 - "appointment/repository.go"
Cohesion: 0.36
Nodes (5): DayMetric, ProfessionalMetric, ServiceMetric, Analytics, StatusCount

### Community 94 - "buildAuthComponents"
Cohesion: 0.29
Nodes (9): Repository, TokenManager, PasswordHasher, NewUseCase(), buildAuthComponents(), boolEnv(), durationSecondsEnv(), seedAdminAccount() (+1 more)

### Community 95 - "booking/usecase.go"
Cohesion: 0.29
Nodes (12): countOverlappingSlots(), isClosedDay(), localMidnight(), overlapsAnyBusy(), parseClock(), PortalAppointment, PortalCustomer, PortalResult (+4 more)

### Community 96 - "TenantsPage.tsx"
Cohesion: 0.09
Nodes (24): TenantsPage, checkSystemReadiness(), createTenant(), fetchAuditLogs(), fetchSystemStatus(), fetchTenants(), updateTenant(), updateTenantSubscription() (+16 more)

### Community 98 - "go_pkg_erpagendamentos_backend_internal_domain_appointment"
Cohesion: 0.60
Nodes (4): BuildICalCalendar(), formatUTC(), sanitizeLine(), go_pkg_erpagendamentos_backend_internal_domain_appointment

### Community 99 - "booking/src/App.tsx"
Cohesion: 0.32
Nodes (4): App(), booking_src_index, initWidgetBridge(), SlugEntryPage()

### Community 110 - "ReminderWorker"
Cohesion: 0.17
Nodes (11): MaintenanceWorker, MaintenanceConfig, UseCase, NewMaintenanceWorker(), formatAppointmentLocal(), UseCase, NewReminderWorker(), startBackgroundWorkers() (+3 more)

### Community 115 - "PortalPage.tsx"
Cohesion: 0.16
Nodes (15): fetchAvailability(), PortalAppointment, portalCancelAppointment(), portalLookup(), portalRescheduleAppointment(), PortalResult, Slot, ACTIONABLE_STATUSES (+7 more)

### Community 116 - "go_pkg_context"
Cohesion: 0.24
Nodes (6): NewPostgresTxRunner(), go_pkg_context, go_pkg_erpagendamentos_backend_internal_domain_payment, go_pkg_io, PostgresTxRunner, TxRunner

### Community 122 - "repositories.go"
Cohesion: 0.11
Nodes (18): main(), resolveDurationEnv(), resolvePositiveIntEnv(), applyDatabaseSQL(), resolveSchemaFile(), loadEnvFile(), resolveMaxRequestBodyBytes(), shouldLimitRequestBody() (+10 more)

## Knowledge Gaps
- **257 isolated node(s):** `graphify-mcp`, `name`, `private`, `version`, `type` (+252 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **13 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Container` connect `Container` to `Payment`, `go_pkg_net_http`, `Review`, `UseCase`, `Tenant`, `Service`, `auth_cookie.go`, `database/sql.DB`, `Resource`, `Queue`, `Entry`, `UseCase`, `UseCase`, `NewContainer`, `UseCase`, `UseCase`, `UseCase`, `UseCase`, `AuditLog`, `EventBroker`, `UseCase`, `UseCase`, `buildAuthComponents`, `ReminderWorker`?**
  _High betweenness centrality (0.045) - this node is a cross-community bridge._
- **Why does `RouterConfig` connect `RouterConfig` to `Payment`, `go_pkg_net_http`, `Review`, `UseCase`, `Tenant`, `Service`, `Resource`, `Entry`, `UseCase`, `UseCase`, `NewContainer`, `UseCase`, `UseCase`, `UseCase`, `UseCase`, `AuditLog`, `EventBroker`, `Container`, `UseCase`, `UseCase`, `time.Duration`, `ReminderWorker`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Why does `UseCase` connect `UseCase` to `net/http.ResponseWriter`, `writeAPIError`, `RouterConfig`, `auth/usecase.go`, `auth_cookie.go`, `withIdempotency`, `NewContainer`, `AuthRepository`, `net/http.Request`, `auth_handler.go`, `rate_limit.go`, `Account`, `Container`, `totp.go`, `auth/entities.go`, `time.Duration`, `.UpdateAccountRole`, `buildAuthComponents`, `booking/usecase.go`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **Are the 131 inferred relationships involving `writeJSON()` (e.g. with `health()` and `readiness()`) actually correct?**
  _`writeJSON()` has 131 INFERRED edges - model-reasoned connections that need verification._
- **Are the 104 inferred relationships involving `resolveRequestTenantID()` (e.g. with `.Analytics()` and `.CheckIn()`) actually correct?**
  _`resolveRequestTenantID()` has 104 INFERRED edges - model-reasoned connections that need verification._
- **What connects `graphify-mcp`, `name`, `private` to the rest of the system?**
  _257 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `adminFetch` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._
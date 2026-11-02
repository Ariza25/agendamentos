# Graph Report - erpagendamentos  (2026-11-02)

## Corpus Check
- 282 files · ~135,791 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2717 nodes · 8122 edges · 107 communities (93 shown, 14 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 687 edges (avg confidence: 0.79)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `df50977e`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- adminFetch
- schema.sql
- ref_react
- Payment
- net/http.Request
- go_pkg_testing
- go_pkg_net_http
- downloadAdminCsv
- adminAuth.ts
- writeAPIError
- AppointmentsPage.tsx
- adminApi.ts
- Review
- UseCase
- Repository
- Tenant
- landing/src/App.tsx
- Service
- go_pkg_time
- auth/entities.go
- RouterConfig
- database/sql.DB
- resource_test.go
- Customer
- UseCase
- formatMoney
- Professional
- CashRegister
- MembershipsPage
- Queue
- Entry
- devDependencies
- devDependencies
- UseCase
- memory/configuration_repository.go
- gencert.go
- UseCase
- withIdempotency
- devDependencies
- ManagePage.tsx
- postgres/commerce_repository.go
- net/http.ResponseWriter
- time.Time
- Container
- BookingPage.tsx
- useToast
- writeProfessionalError
- AuthRepository
- UseCase
- AppointmentHandler
- bookingApi.ts
- memoryRepositories
- Membership
- Package
- UseCase
- WaitingListPage
- PackageCredit
- GiftCard
- sync.Mutex
- AuditLog
- checkout/usecase.go
- EventBroker
- rate_limit.go
- Appointment
- writeInvalidPayload
- templates.go
- compilerOptions
- AppointmentRepository
- compilerOptions
- compilerOptions
- OutboxRepository
- UseCase
- totp.go
- writeResourceError
- testing.T
- ERP Agendamentos
- NewBcryptHasher
- UseCase
- go_pkg_context
- booking/src/lib/apiError.ts
- time.Duration
- ExportHandler
- FinancePage.tsx
- observability.go
- admin/vite.config.ts
- appointment/repository.go
- UseCase
- TenantsPage.tsx
- ical.go
- booking/src/App.tsx
- graphify
- rules/graphify.md
- workflows/graphify.md
- generate-certs.sh
- backup_postgres.sh
- certbot_ssl.sh
- catalog.finance_transactions
- context.Context
- go_pkg_path_filepath
- erpagendamentos/backend
- writeServiceError
- PortalPage.tsx
- go_pkg_os
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
- `newResourceUseCase()` --calls--> `NewUseCase()`  [INFERRED]
  backend/internal/application/appointment/resource_test.go → backend/internal/application/appointment/usecase.go
- `newTestUseCase()` --calls--> `NewUseCase()`  [INFERRED]
  backend/internal/application/booking/usecase_test.go → backend/internal/application/booking/usecase.go
- `TestManagementCodeNormalization()` --calls--> `hashManagementCode()`  [INFERRED]
  backend/internal/application/booking/usecase_test.go → backend/internal/application/booking/usecase.go
- `TestOutboxPublishAndDispatch()` --calls--> `NewUseCase()`  [INFERRED]
  backend/internal/application/outbox/usecase_test.go → backend/internal/application/outbox/usecase.go

## Import Cycles
- None detected.

## Communities (107 total, 14 thin omitted)

### Community 0 - "adminFetch"
Cohesion: 0.05
Nodes (81): closeCashRegister(), connectCalendar(), createCashRegisterEntry(), createCustomer(), createFinancialTransaction(), createGiftCard(), createMembership(), createPackage() (+73 more)

### Community 1 - "schema.sql"
Cohesion: 0.12
Nodes (35): catalog.appointments, catalog.audit_logs, catalog.auth_rate_limits, catalog.auth_sessions, catalog.calendar_connections, catalog.cash_register_entries, catalog.cash_registers, catalog.commission_payouts (+27 more)

### Community 2 - "ref_react"
Cohesion: 0.07
Nodes (89): AdminRoutesProps, AgendaPage, AnalyticsPage, AppointmentsPage, CommissionsPage, ConfigurationPage, CustomersPage, EstablishmentPage (+81 more)

### Community 3 - "Payment"
Cohesion: 0.06
Nodes (25): UseCase, NewUseCase(), normalizeStatus(), buildPaymentProvider(), IntentResult, Payment, NewMercadoPagoProvider(), NewMockProvider() (+17 more)

### Community 4 - "net/http.Request"
Cohesion: 0.10
Nodes (13): resolveRequestTenantID(), health(), writePaymentError(), writeCSV(), writeInternalError(), writeJSON(), net/http.Request, CalendarSyncHandler (+5 more)

### Community 5 - "go_pkg_testing"
Cohesion: 0.10
Nodes (12): TestJobQueue(), NewJobRepository(), go_pkg_erpagendamentos_backend_internal_application_calendar, go_pkg_erpagendamentos_backend_internal_application_cashregister, go_pkg_erpagendamentos_backend_internal_application_checkout, go_pkg_erpagendamentos_backend_internal_application_onboarding, go_pkg_erpagendamentos_backend_internal_application_stock, go_pkg_erpagendamentos_backend_internal_domain_stock (+4 more)

### Community 6 - "go_pkg_net_http"
Cohesion: 0.11
Nodes (18): deploymentEnvironment(), isProductionEnv(), secureTemporaryPassword(), go_pkg_encoding_json, go_pkg_erpagendamentos_backend_internal_application_booking, go_pkg_erpagendamentos_backend_internal_application_commerce, go_pkg_erpagendamentos_backend_internal_application_configuration, go_pkg_erpagendamentos_backend_internal_application_customer (+10 more)

### Community 7 - "downloadAdminCsv"
Cohesion: 0.11
Nodes (21): downloadAdminCsv(), fetchPayments(), fetchReviews(), handleExport(), CustomersPage(), handleDelete(), handleExport(), handleSearch() (+13 more)

### Community 8 - "adminAuth.ts"
Cohesion: 0.06
Nodes (50): App(), handleLogout(), AppProviders(), AppProvidersProps, ResetPasswordPage(), handleSubmit(), AppRouter(), admin_src_index (+42 more)

### Community 9 - "writeAPIError"
Cohesion: 0.08
Nodes (25): requireMasterAdmin(), resolveAuthClaims(), resolveTenantID(), browserSessionPayload(), buildPasswordResetURL(), exposeAuthTokenForTests(), AuthHandler, newAuthMFAHandler() (+17 more)

### Community 10 - "AppointmentsPage.tsx"
Cohesion: 0.08
Nodes (48): AppointmentFormModal(), computeEndsAt(), handleSubmit(), AppointmentFormModalProps, STATUS_OPTIONS, checkInAppointment(), createAppointment(), deleteAppointment() (+40 more)

### Community 11 - "adminApi.ts"
Cohesion: 0.06
Nodes (63): AppointmentFilters, createCommissionRule(), deleteCommissionRule(), fetchCommissionPayouts(), fetchCommissionRecords(), fetchCommissionRules(), fetchCommissionSummaries(), fetchCustomers() (+55 more)

### Community 12 - "Review"
Cohesion: 0.14
Nodes (12): UseCase, NewUseCase(), Review, NewReviewRepository(), scanReviewRow(), scanReviewRows(), NewReviewHandler(), ReviewRepository (+4 more)

### Community 13 - "UseCase"
Cohesion: 0.08
Nodes (12): NewUseCase(), UseCase, NewUseCase(), TestStockUseCase(), Product, Repository, ServiceSupply, StockMovement (+4 more)

### Community 14 - "Repository"
Cohesion: 0.18
Nodes (11): UseCase, NewUseCase(), formatAppointmentLocal(), UseCase, NewReminderWorker(), NewCalendarHandler(), NewConfigurationHandler(), Repository (+3 more)

### Community 15 - "Tenant"
Cohesion: 0.11
Nodes (12): UseCase, Tenant, UpdateSubscriptionInput, defaultModules(), isTenantIdentityConflict(), defaultModules(), scanTenantRow(), scanTenantRows() (+4 more)

### Community 16 - "landing/src/App.tsx"
Cohesion: 0.10
Nodes (22): App(), CtaSection(), FEATURES, FeaturesSection(), Footer(), Header(), NAV_ITEMS, HeroSection() (+14 more)

### Community 17 - "Service"
Cohesion: 0.12
Nodes (8): UseCase, Service, scanServiceRow(), scanServiceRows(), NewServiceHandler(), ServiceRepository, ServiceRepository, serviceScanner

### Community 18 - "go_pkg_time"
Cohesion: 0.09
Nodes (23): go_pkg_database_sql, go_pkg_erpagendamentos_backend_internal_domain_appointment, go_pkg_erpagendamentos_backend_internal_domain_auth, go_pkg_erpagendamentos_backend_internal_domain_calendar, go_pkg_erpagendamentos_backend_internal_domain_cashregister, go_pkg_erpagendamentos_backend_internal_domain_commerce, go_pkg_erpagendamentos_backend_internal_domain_configuration, go_pkg_erpagendamentos_backend_internal_domain_customer (+15 more)

### Community 19 - "auth/entities.go"
Cohesion: 0.12
Nodes (14): ChangePasswordInput, LoginInput, PasswordResetCodeConfirmInput, PasswordResetConfirmInput, PasswordResetRequestInput, RegisterInput, Session, ValidationError (+6 more)

### Community 20 - "RouterConfig"
Cohesion: 0.07
Nodes (71): bearerAuthAllowed(), hasRole(), isAdminRole(), isPrivilegedAdminRole(), requestWithAuthClaims(), requireAdmin(), requireOwnerOrAdmin(), requirePrivilegedAdmin() (+63 more)

### Community 21 - "database/sql.DB"
Cohesion: 0.12
Nodes (32): applyDatabaseSQL(), resolveSchemaFile(), initializeRepositories(), postgresRepositories(), NewAppointmentRepository(), NewAuditRepository(), NewAuthRepository(), NewCalendarRepository() (+24 more)

### Community 22 - "resource_test.go"
Cohesion: 0.54
Nodes (7): UseCase, newAppointment(), newResourceUseCase(), TestCreateAutoAssignsAvailableResource(), TestCreateRejectsBusyExplicitResource(), TestCreateRejectsUnknownResource(), TestUpdateAutoReassignsFreedResource()

### Community 23 - "Customer"
Cohesion: 0.13
Nodes (10): UseCase, Customer, onlyDigits(), scanCustomerRow(), scanCustomerRows(), NewCustomerHandler(), BookingResult, CustomerRepository (+2 more)

### Community 24 - "UseCase"
Cohesion: 0.08
Nodes (11): UseCase, NewUseCase(), CommissionRecord, CommissionRule, Payout, Repository, Summary, NewCommissionHandler() (+3 more)

### Community 25 - "formatMoney"
Cohesion: 0.06
Nodes (39): CheckoutComandaModal(), handleSubmit(), checkoutAppointment(), createCommissionPayout(), createMemberPlan(), createPackageCredit(), fetchGiftCards(), fetchPackageCredits() (+31 more)

### Community 26 - "Professional"
Cohesion: 0.06
Nodes (19): UseCase, UseCase, NewUseCase(), Professional, Resource, NewResourceRepository(), scanProfessionalRow(), scanProfessionalRows() (+11 more)

### Community 27 - "CashRegister"
Cohesion: 0.13
Nodes (9): UseCase, NewUseCase(), CashRegister, CashRegisterEntry, Repository, ShiftSummary, NewCashRegisterHandler(), CashRegisterRepository (+1 more)

### Community 28 - "MembershipsPage"
Cohesion: 0.36
Nodes (8): fetchMemberPlans(), fetchMemberships(), MembershipsPage(), handleDelete(), handlePlanStatusChange(), handleToggleActive(), loadMemberships(), planStatusTone()

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
Cohesion: 0.16
Nodes (6): RecurrenceInput, Repository, UseCase, NewUseCase(), parseRFC3339(), validStatusChange()

### Community 34 - "memory/configuration_repository.go"
Cohesion: 0.46
Nodes (5): cloneConfiguration(), defaultMemoryConfiguration(), NewConfigurationRepository(), normalizeConfigurationTenantID(), ConfigurationRepository

### Community 35 - "gencert.go"
Cohesion: 0.12
Nodes (13): NewSecretCipher(), go_pkg_crypto_aes, go_pkg_crypto_cipher, go_pkg_crypto_rand, go_pkg_crypto_rsa, go_pkg_crypto_x509, go_pkg_crypto_x509_pkix, go_pkg_encoding_base64 (+5 more)

### Community 36 - "UseCase"
Cohesion: 0.11
Nodes (10): UseCase, NewUseCase(), TestCalendarUseCase(), Connection, Repository, SyncEvent, NewCalendarRepository(), NewCalendarSyncHandler() (+2 more)

### Community 37 - "withIdempotency"
Cohesion: 0.14
Nodes (14): isHopByHopHeader(), NewIdempotencyStore(), TestIdempotencyBlocksInFlightConcurrentRequests(), TestIdempotencyDoesNotCacheServerFailures(), TestIdempotencyIgnoredOnGetRequests(), TestIdempotencyReplaysCompletedResponse(), withIdempotency(), go_pkg_bytes (+6 more)

### Community 38 - "devDependencies"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, devDependencies, tailwindcss, @tailwindcss/vite, @types/node (+24 more)

### Community 39 - "ManagePage.tsx"
Cohesion: 0.14
Nodes (24): Professional, PublicAppointment, Service, ManageStep, Button, ButtonProps, ButtonSize, ButtonVariant (+16 more)

### Community 40 - "postgres/commerce_repository.go"
Cohesion: 0.14
Nodes (7): MemberPlan, nullIfEmpty(), NewMemberPlanRepository(), scanMemberPlans(), database/sql.Rows, MemberPlanRepository, MemberPlanRepository

### Community 41 - "net/http.ResponseWriter"
Cohesion: 0.10
Nodes (15): writeLazyJSON(), writeBookingError(), writeConfigurationError(), writeCustomerError(), HandleMetrics(), writeNotificationError(), writeConflictError(), writeNotFoundError() (+7 more)

### Community 42 - "time.Time"
Cohesion: 0.18
Nodes (15): countOverlappingSlots(), isClosedDay(), localMidnight(), overlapsAnyBusy(), parseClock(), workingHoursFor(), Configuration, normalizePostgresTenantID() (+7 more)

### Community 43 - "Container"
Cohesion: 0.12
Nodes (19): NewUseCase(), NewUseCase(), UseCase, NewUseCase(), NewUseCase(), NewUseCase(), NewUseCase(), ensureDefaultTenant() (+11 more)

### Community 44 - "BookingPage.tsx"
Cohesion: 0.13
Nodes (19): formatDateLong(), formatDateTime(), formatMoney(), formatTime(), fromDateKey(), nextDateKeys(), SHORT_DAY_LABELS, toDateKey() (+11 more)

### Community 45 - "useToast"
Cohesion: 0.06
Nodes (47): AdminRoutes(), AdminShell(), handleCopyBookingUrl(), AdminShellProps, AdminNavigationItem, navigationGroups, navigationItems, EstablishmentModal() (+39 more)

### Community 47 - "AuthRepository"
Cohesion: 0.06
Nodes (12): Account, AuthState, DeviceSession, PasswordResetToken, UpdateProfileInput, cloneAuthAccount(), resetTokenExpired(), sessionExpired() (+4 more)

### Community 48 - "UseCase"
Cohesion: 0.18
Nodes (10): appointmentDurationMinutes(), formatAppointmentForTemplate(), generateManagementCode(), UseCase, hashManagementCode(), parseRFC3339(), BookingInput, PortalAppointment (+2 more)

### Community 49 - "AppointmentHandler"
Cohesion: 0.19
Nodes (5): writeAppointmentError(), resolveActorTenantID(), writeAdminAuthError(), writeBadRequestError(), AppointmentHandler

### Community 50 - "bookingApi.ts"
Cohesion: 0.10
Nodes (30): backendOrigin, buildApiUrl(), AvailabilityResponse, BookingResult, cancelPublicAppointment(), createPublicBooking(), fetchAvailability(), fetchPublicAppointment() (+22 more)

### Community 51 - "memoryRepositories"
Cohesion: 0.17
Nodes (19): TestCashRegisterFlow(), TestCheckoutFlow(), TestCommissionUseCase(), TestOnboardingSignup(), memoryRepositories(), NewAppointmentRepository(), NewAuthRepository(), NewCashRegisterRepository() (+11 more)

### Community 52 - "Membership"
Cohesion: 0.15
Nodes (3): Membership, MembershipRepository, MembershipRepository

### Community 53 - "Package"
Cohesion: 0.15
Nodes (3): Package, PackageRepository, PackageRepository

### Community 54 - "UseCase"
Cohesion: 0.10
Nodes (11): UseCase, NewUseCase(), TestFinanceUseCase(), Repository, Summary, Transaction, NewFinanceHandler(), CategoryTotal (+3 more)

### Community 55 - "WaitingListPage"
Cohesion: 0.67
Nodes (4): fetchWaitingList(), WaitingListPage(), handleStatusChange(), loadEntries()

### Community 56 - "PackageCredit"
Cohesion: 0.19
Nodes (4): PackageCredit, scanPackageCredits(), PackageCreditRepository, PackageCreditRepository

### Community 57 - "GiftCard"
Cohesion: 0.14
Nodes (5): generateGiftCardCode(), GiftCard, NewGiftCardRepository(), GiftCardRepository, GiftCardRepository

### Community 58 - "sync.Mutex"
Cohesion: 0.13
Nodes (13): lazyAppHandler, UseCase, NewUseCase(), Log, NewNotificationRepository(), NewNotificationHandler(), sync.Mutex, NotificationRepository (+5 more)

### Community 60 - "AuditLog"
Cohesion: 0.17
Nodes (11): AuditLog, Repository, UseCase, NewUseCase(), TestAuditLogRecordingAndListing(), NewAuditRepository(), NewAuditHandler(), TestAuditAndSystemHandlers() (+3 more)

### Community 61 - "checkout/usecase.go"
Cohesion: 0.16
Nodes (6): CheckoutInput, CheckoutResult, ProductItem, go_pkg_erpagendamentos_backend_internal_application_commission, go_pkg_erpagendamentos_backend_internal_domain_commission, go_pkg_erpagendamentos_backend_internal_domain_finance

### Community 62 - "EventBroker"
Cohesion: 0.14
Nodes (14): UseCase, NewBookingHandler(), NewCheckoutHandler(), EventBroker, NewMercadoPagoWebhookHandler(), RegisterMetricsBroker(), AuditLogger, OutboxPublisher (+6 more)

### Community 63 - "rate_limit.go"
Cohesion: 0.17
Nodes (19): buildCompoundIdempotencyKey(), forwardedClientIP(), isSensitiveRateLimitedRoute(), isTrustedProxy(), newMemoryRateLimiter(), normalizedClientID(), persistentRateLimitConfig(), rateLimitRouteKey() (+11 more)

### Community 65 - "Appointment"
Cohesion: 0.23
Nodes (4): Appointment, parseAppointmentBound(), parseAppointmentInterval(), AppointmentRepository

### Community 66 - "writeInvalidPayload"
Cohesion: 0.29
Nodes (3): writeCommerceError(), writeInvalidPayload(), CommerceHandler

### Community 68 - "templates.go"
Cohesion: 0.27
Nodes (17): greeting(), PublicBookingOrigin(), RenderCancellation(), RenderConfirmation(), RenderReminder(), RenderReschedule(), sampleData(), TestGreetingWithoutName() (+9 more)

### Community 70 - "compilerOptions"
Cohesion: 0.09
Nodes (21): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+13 more)

### Community 71 - "AppointmentRepository"
Cohesion: 0.19
Nodes (5): appointmentSelectColumns(), scanAppointmentRow(), scanAppointmentRows(), AppointmentRepository, appointmentScanner

### Community 72 - "compilerOptions"
Cohesion: 0.09
Nodes (21): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+13 more)

### Community 73 - "compilerOptions"
Cohesion: 0.09
Nodes (21): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+13 more)

### Community 74 - "OutboxRepository"
Cohesion: 0.14
Nodes (6): TestOutboxPublishAndDispatch(), NewOutboxRepository(), quoteLiteral(), OutboxRepository, OutboxEvent, OutboxRepository

### Community 75 - "UseCase"
Cohesion: 0.24
Nodes (7): UseCase, NewUseCase(), NewSystemHandler(), go_pkg_runtime, EventDispatcher, Repository, SystemHandler

### Community 76 - "totp.go"
Cohesion: 0.15
Nodes (11): MFAStatus, generateTOTPSecret(), totpCode(), totpURI(), validateTOTP(), isMFAEligibleRole(), go_pkg_crypto_hmac, go_pkg_crypto_sha1 (+3 more)

### Community 78 - "testing.T"
Cohesion: 0.35
Nodes (11): UseCase, newTestUseCase(), TestAvailableSlotsRespectsClosedDays(), TestCancelPublic(), TestCreateBookingGeneratesManagementCode(), TestManagementCodeNormalization(), TestPortalCancelValidatesOwnership(), TestPortalLookup() (+3 more)

### Community 80 - "ERP Agendamentos"
Cohesion: 0.14
Nodes (13): Banco de Dados & Schema, CI, Endpoints principais, ERP Agendamentos, Estrutura, Pré-requisitos, Rodando com Docker (produção local), Rodando em desenvolvimento (sem banco) (+5 more)

### Community 81 - "NewBcryptHasher"
Cohesion: 0.40
Nodes (3): NewBcryptHasher(), go_pkg_golang_org_x_crypto_bcrypt, BcryptHasher

### Community 82 - "UseCase"
Cohesion: 0.29
Nodes (9): UseCase, NewUseCase(), NewCommerceHandler(), repositorySet, GiftCardRepository, MemberPlanRepository, MembershipRepository, PackageCreditRepository (+1 more)

### Community 83 - "go_pkg_context"
Cohesion: 0.11
Nodes (12): buildNotificationProvider(), intEnv(), logFallback(), NewEmailProvider(), NewLogProvider(), go_pkg_context, go_pkg_erpagendamentos_backend_internal_domain_payment, go_pkg_io (+4 more)

### Community 85 - "booking/src/lib/apiError.ts"
Cohesion: 0.24
Nodes (10): ApiError, ApiErrorPayload, buildApiError(), parseApiResponse(), API_ERROR_MESSAGES, EXACT_MESSAGE_MAP, getApiErrorMessage(), SUBSTRING_MESSAGE_PATTERNS (+2 more)

### Community 86 - "time.Duration"
Cohesion: 0.18
Nodes (9): Claims, CleanupSummary, decodeSegment(), encodeSegment(), NewJWTManager(), go_pkg_crypto_sha256, time.Duration, JWTManager (+1 more)

### Community 88 - "FinancePage.tsx"
Cohesion: 0.15
Nodes (21): FinancePage, fetchCashRegisterHistory(), fetchCurrentCashRegister(), fetchFinancialSummary(), fetchFinancialTransactions(), ActiveTab, FinancePage(), handleCloseShift() (+13 more)

### Community 91 - "observability.go"
Cohesion: 0.27
Nodes (8): isSensitiveLogKey(), redactSensitiveText(), sanitizeLogString(), sanitizeLogValue(), truncateLogString(), go_pkg_regexp, contextKey, responseRecorder

### Community 92 - "admin/vite.config.ts"
Cohesion: 0.60
Nodes (3): ref_tailwindcss_vite, ref_vite, ref_vitejs_plugin_react

### Community 93 - "appointment/repository.go"
Cohesion: 0.36
Nodes (5): DayMetric, ProfessionalMetric, ServiceMetric, Analytics, StatusCount

### Community 94 - "UseCase"
Cohesion: 0.12
Nodes (16): Repository, SecretCipher, TokenManager, canManageTenantAccount(), PasswordHasher, UseCase, isPrivilegedAdminRole(), isProductionAuthEnvironment() (+8 more)

### Community 96 - "TenantsPage.tsx"
Cohesion: 0.09
Nodes (24): TenantsPage, checkSystemReadiness(), createTenant(), fetchAuditLogs(), fetchSystemStatus(), fetchTenants(), updateTenant(), updateTenantSubscription() (+16 more)

### Community 98 - "ical.go"
Cohesion: 0.83
Nodes (3): BuildICalCalendar(), formatUTC(), sanitizeLine()

### Community 99 - "booking/src/App.tsx"
Cohesion: 0.16
Nodes (5): App(), booking_src_index, initWidgetBridge(), BookingPage(), SlugEntryPage()

### Community 110 - "context.Context"
Cohesion: 0.20
Nodes (9): MaintenanceWorker, MaintenanceConfig, UseCase, NewMaintenanceWorker(), NewPostgresTxRunner(), context.Context, NopTxRunner, PostgresTxRunner (+1 more)

### Community 115 - "PortalPage.tsx"
Cohesion: 0.16
Nodes (14): PortalAppointment, portalCancelAppointment(), portalLookup(), portalRescheduleAppointment(), PortalResult, Slot, ACTIONABLE_STATUSES, PortalPage() (+6 more)

### Community 122 - "go_pkg_os"
Cohesion: 0.11
Nodes (17): newLazyAppHandler(), main(), resolveDurationEnv(), resolvePositiveIntEnv(), loadEnvFile(), resolveMaxRequestBodyBytes(), shouldLimitRequestBody(), withRequestBodyLimit() (+9 more)

## Knowledge Gaps
- **259 isolated node(s):** `graphify-mcp`, `name`, `private`, `version`, `type` (+254 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **14 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Container` connect `Container` to `Payment`, `go_pkg_net_http`, `Review`, `UseCase`, `Repository`, `Tenant`, `Service`, `RouterConfig`, `database/sql.DB`, `Customer`, `UseCase`, `Professional`, `CashRegister`, `Queue`, `Entry`, `UseCase`, `UseCase`, `UseCase`, `UseCase`, `sync.Mutex`, `AuditLog`, `EventBroker`, `UseCase`, `UseCase`, `UseCase`, `context.Context`?**
  _High betweenness centrality (0.046) - this node is a cross-community bridge._
- **Why does `RouterConfig` connect `RouterConfig` to `Payment`, `go_pkg_net_http`, `Review`, `UseCase`, `Repository`, `Tenant`, `Service`, `Customer`, `UseCase`, `Professional`, `CashRegister`, `Entry`, `UseCase`, `UseCase`, `Container`, `UseCase`, `UseCase`, `sync.Mutex`, `AuditLog`, `EventBroker`, `UseCase`, `UseCase`, `time.Duration`, `UseCase`, `context.Context`?**
  _High betweenness centrality (0.031) - this node is a cross-community bridge._
- **Why does `UseCase` connect `UseCase` to `net/http.Request`, `withIdempotency`, `writeAPIError`, `time.Time`, `Container`, `totp.go`, `AuthRepository`, `AppointmentHandler`, `auth/entities.go`, `RouterConfig`, `time.Duration`, `rate_limit.go`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **Are the 131 inferred relationships involving `writeJSON()` (e.g. with `health()` and `readiness()`) actually correct?**
  _`writeJSON()` has 131 INFERRED edges - model-reasoned connections that need verification._
- **Are the 104 inferred relationships involving `resolveRequestTenantID()` (e.g. with `.Analytics()` and `.CheckIn()`) actually correct?**
  _`resolveRequestTenantID()` has 104 INFERRED edges - model-reasoned connections that need verification._
- **What connects `graphify-mcp`, `name`, `private` to the rest of the system?**
  _259 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `adminFetch` be split into smaller, more focused modules?**
  _Cohesion score 0.04938271604938271 - nodes in this community are weakly interconnected._
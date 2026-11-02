# Graph Report - erpagendamentos  (2026-11-02)

## Corpus Check
- 282 files · ~134,385 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2712 nodes · 8112 edges · 118 communities (99 shown, 19 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 687 edges (avg confidence: 0.79)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `5c80d7ba`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- adminFetch
- schema.sql
- ref_react
- Payment
- writeJSON
- go_pkg_testing
- go_pkg_net_http
- downloadAdminCsv
- adminAuth.ts
- writeAPIError
- AppointmentsPage.tsx
- adminApi.ts
- Review
- UseCase
- UseCase
- Tenant
- landing/src/App.tsx
- Service
- go_pkg_time
- .Add
- RouterConfig
- database/sql.DB
- resource_test.go
- Customer
- CommissionRepository
- formatMoney
- Resource
- CashRegister
- Professional
- JobRepository
- Entry
- devDependencies
- devDependencies
- Repository
- UseCase
- auth/usecase.go
- UseCase
- Queue
- devDependencies
- booking/src/pages/WaitingListPage.tsx
- MemberPlan
- net/http.ResponseWriter
- Configuration
- NewContainer
- BookingPage.tsx
- useToast
- InventoryPage
- AuthRepository
- UseCase
- net/http.Request
- bookingApi.ts
- testing.T
- Membership
- Package
- Transaction
- UseCase
- PackageCredit
- GiftCard
- Repository
- BookingPage
- UseCase
- checkout/usecase.go
- EventBroker
- rate_limit.go
- AuthRepository
- Appointment
- writeInvalidPayload
- UseCase
- templates.go
- UseCase
- compilerOptions
- AppointmentRepository
- compilerOptions
- compilerOptions
- writeProfessionalError
- UseCase
- MembershipsPage
- writeResourceError
- time.Duration
- auth/entities.go
- ERP Agendamentos
- DeviceSession
- UseCase
- email_provider.go
- booking/src/lib/apiError.ts
- Claims
- ExportHandler
- FinancePage.tsx
- .UpdateAccountRole
- SecretCipher
- observability.go
- admin/vite.config.ts
- appointment/repository.go
- buildAuthComponents
- time.Time
- TenantsPage
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
- go_pkg_context
- writePaymentError
- TenantSubscriptionModal
- body_limit.go
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
- `handleSubmit()` --calls--> `updateCurrentTenant()`  [EXTRACTED]
  admin/src/components/EstablishmentModal.tsx → admin/src/lib/adminApi.ts
- `handleSubmit()` --calls--> `updateConfiguration()`  [EXTRACTED]
  admin/src/pages/ConfigurationPage.tsx → admin/src/lib/adminApi.ts
- `handleSubmit()` --calls--> `updateCurrentTenant()`  [EXTRACTED]
  admin/src/pages/EstablishmentPage.tsx → admin/src/lib/adminApi.ts
- `main()` --calls--> `newLazyAppHandler()`  [INFERRED]
  backend/cmd/api/main.go → backend/cmd/api/lazy_handler.go

## Import Cycles
- None detected.

## Communities (118 total, 19 thin omitted)

### Community 0 - "adminFetch"
Cohesion: 0.05
Nodes (81): closeCashRegister(), connectCalendar(), createCashRegisterEntry(), createCustomer(), createFinancialTransaction(), createGiftCard(), createMembership(), createPackage() (+73 more)

### Community 1 - "schema.sql"
Cohesion: 0.12
Nodes (35): catalog.appointments, catalog.audit_logs, catalog.auth_rate_limits, catalog.auth_sessions, catalog.calendar_connections, catalog.cash_register_entries, catalog.cash_registers, catalog.commission_payouts (+27 more)

### Community 2 - "ref_react"
Cohesion: 0.07
Nodes (88): AdminRoutesProps, AgendaPage, AppointmentsPage, CommissionsPage, ConfigurationPage, CustomersPage, EstablishmentPage, GiftCardsPage (+80 more)

### Community 3 - "Payment"
Cohesion: 0.06
Nodes (25): UseCase, NewUseCase(), normalizeStatus(), buildPaymentProvider(), IntentResult, Payment, NewMercadoPagoProvider(), NewMockProvider() (+17 more)

### Community 4 - "writeJSON"
Cohesion: 0.11
Nodes (8): resolveRequestTenantID(), writeInternalError(), writeJSON(), CalendarSyncHandler, CashRegisterHandler, CommissionHandler, FinanceHandler, StockHandler

### Community 6 - "go_pkg_net_http"
Cohesion: 0.11
Nodes (24): deploymentEnvironment(), isProductionEnv(), go_pkg_encoding_json, go_pkg_erpagendamentos_backend_internal_application_booking, go_pkg_erpagendamentos_backend_internal_application_calendar, go_pkg_erpagendamentos_backend_internal_application_cashregister, go_pkg_erpagendamentos_backend_internal_application_checkout, go_pkg_erpagendamentos_backend_internal_application_commerce (+16 more)

### Community 7 - "downloadAdminCsv"
Cohesion: 0.11
Nodes (21): downloadAdminCsv(), fetchPayments(), fetchReviews(), handleExport(), CustomersPage(), handleDelete(), handleExport(), handleSearch() (+13 more)

### Community 8 - "adminAuth.ts"
Cohesion: 0.06
Nodes (48): App(), handleLogout(), AppProviders(), AppProvidersProps, ResetPasswordPage(), handleSubmit(), AppRouter(), admin_src_index (+40 more)

### Community 9 - "writeAPIError"
Cohesion: 0.10
Nodes (21): resolveAuthClaims(), resolveTenantID(), browserSessionPayload(), exposeAuthTokenForTests(), AuthHandler, health(), readiness(), logError() (+13 more)

### Community 10 - "AppointmentsPage.tsx"
Cohesion: 0.06
Nodes (62): AnalyticsPage, AppointmentFormModal(), computeEndsAt(), handleSubmit(), AppointmentFormModalProps, STATUS_OPTIONS, CheckoutComandaModal(), handleSubmit() (+54 more)

### Community 11 - "adminApi.ts"
Cohesion: 0.06
Nodes (64): AppointmentFilters, createCommissionRule(), deleteCommissionRule(), fetchCommissionPayouts(), fetchCommissionRecords(), fetchCommissionRules(), fetchCommissionSummaries(), fetchCustomers() (+56 more)

### Community 12 - "Review"
Cohesion: 0.14
Nodes (12): UseCase, NewUseCase(), Review, NewReviewRepository(), scanReviewRow(), scanReviewRows(), NewReviewHandler(), ReviewRepository (+4 more)

### Community 13 - "UseCase"
Cohesion: 0.09
Nodes (7): UseCase, Product, ServiceSupply, StockMovement, NewStockHandler(), StockRepository, StockRepository

### Community 14 - "UseCase"
Cohesion: 0.22
Nodes (5): UseCase, NewAppointmentHandler(), NewFinanceHandler(), NewMercadoPagoWebhookHandler(), MercadoPagoWebhookHandler

### Community 15 - "Tenant"
Cohesion: 0.10
Nodes (13): UseCase, Tenant, UpdateSubscriptionInput, defaultModules(), isTenantIdentityConflict(), defaultModules(), scanTenantRow(), scanTenantRows() (+5 more)

### Community 16 - "landing/src/App.tsx"
Cohesion: 0.10
Nodes (22): App(), CtaSection(), FEATURES, FeaturesSection(), Footer(), Header(), NAV_ITEMS, HeroSection() (+14 more)

### Community 17 - "Service"
Cohesion: 0.17
Nodes (6): Service, scanServiceRow(), scanServiceRows(), ServiceRepository, ServiceRepository, serviceScanner

### Community 18 - "go_pkg_time"
Cohesion: 0.09
Nodes (27): go_pkg_database_sql, go_pkg_erpagendamentos_backend_internal_domain_appointment, go_pkg_erpagendamentos_backend_internal_domain_calendar, go_pkg_erpagendamentos_backend_internal_domain_cashregister, go_pkg_erpagendamentos_backend_internal_domain_commerce, go_pkg_erpagendamentos_backend_internal_domain_commission, go_pkg_erpagendamentos_backend_internal_domain_configuration, go_pkg_erpagendamentos_backend_internal_domain_customer (+19 more)

### Community 19 - ".Add"
Cohesion: 0.15
Nodes (11): ChangePasswordInput, PasswordResetCodeConfirmInput, PasswordResetConfirmInput, PasswordResetRequestInput, ValidationError, hashResetToken(), NewValidationError(), passwordResetCodeChallengeKey() (+3 more)

### Community 20 - "RouterConfig"
Cohesion: 0.05
Nodes (76): lazyAppHandler, newLazyAppHandler(), writeLazyJSON(), Container, bearerAuthAllowed(), hasRole(), isAdminRole(), isPrivilegedAdminRole() (+68 more)

### Community 21 - "database/sql.DB"
Cohesion: 0.10
Nodes (37): applyDatabaseSQL(), resolveSchemaFile(), initializeRepositories(), postgresRepositories(), NewAppointmentRepository(), NewAuditRepository(), NewAuthRepository(), NewCalendarRepository() (+29 more)

### Community 22 - "resource_test.go"
Cohesion: 0.54
Nodes (7): UseCase, newAppointment(), newResourceUseCase(), TestCreateAutoAssignsAvailableResource(), TestCreateRejectsBusyExplicitResource(), TestCreateRejectsUnknownResource(), TestUpdateAutoReassignsFreedResource()

### Community 23 - "Customer"
Cohesion: 0.18
Nodes (8): Customer, onlyDigits(), scanCustomerRow(), scanCustomerRows(), BookingResult, CustomerRepository, CustomerRepository, customerScanner

### Community 24 - "CommissionRepository"
Cohesion: 0.10
Nodes (6): CommissionRecord, CommissionRule, Payout, Summary, CommissionRepository, CommissionRepository

### Community 25 - "formatMoney"
Cohesion: 0.10
Nodes (26): createCommissionPayout(), createMemberPlan(), createPackageCredit(), fetchGiftCards(), fetchPackageCredits(), fetchPackages(), formatMoney(), CommissionPayoutModal() (+18 more)

### Community 26 - "Resource"
Cohesion: 0.14
Nodes (7): UseCase, Resource, scanResourceRow(), NewResourceHandler(), ResourceRepository, ResourceRepository, resourceScanner

### Community 27 - "CashRegister"
Cohesion: 0.12
Nodes (10): UseCase, NewUseCase(), CashRegister, CashRegisterEntry, Repository, ShiftSummary, Repository, NewCashRegisterHandler() (+2 more)

### Community 28 - "Professional"
Cohesion: 0.13
Nodes (9): UseCase, NewUseCase(), Professional, scanProfessionalRow(), scanProfessionalRows(), NewProfessionalHandler(), ProfessionalRepository, ProfessionalRepository (+1 more)

### Community 29 - "JobRepository"
Cohesion: 0.15
Nodes (3): Job, JobRepository, JobRepository

### Community 30 - "Entry"
Cohesion: 0.12
Nodes (12): UseCase, NewUseCase(), Entry, NewWaitingListRepository(), scanWaitingRow(), scanWaitingRows(), NewWaitingListHandler(), WaitingListRepository (+4 more)

### Community 31 - "devDependencies"
Cohesion: 0.06
Nodes (34): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+26 more)

### Community 32 - "devDependencies"
Cohesion: 0.06
Nodes (34): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+26 more)

### Community 33 - "Repository"
Cohesion: 0.29
Nodes (7): UseCase, NewUseCase(), NewCalendarHandler(), NewConfigurationHandler(), Repository, CalendarHandler, ConfigurationHandler

### Community 34 - "UseCase"
Cohesion: 0.18
Nodes (4): RecurrenceInput, UseCase, parseRFC3339(), validStatusChange()

### Community 35 - "auth/usecase.go"
Cohesion: 0.07
Nodes (24): loadEnvFile(), secureTemporaryPassword(), go_pkg_bufio, go_pkg_bytes, go_pkg_crypto_aes, go_pkg_crypto_cipher, go_pkg_crypto_hmac, go_pkg_crypto_rand (+16 more)

### Community 36 - "UseCase"
Cohesion: 0.11
Nodes (10): UseCase, NewUseCase(), TestCalendarUseCase(), Connection, Repository, SyncEvent, NewCalendarRepository(), NewCalendarSyncHandler() (+2 more)

### Community 37 - "Queue"
Cohesion: 0.25
Nodes (5): Queue, NewQueue(), Repository, sync.WaitGroup, HandlerFunc

### Community 38 - "devDependencies"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, devDependencies, tailwindcss, @tailwindcss/vite, @types/node (+24 more)

### Community 39 - "booking/src/pages/WaitingListPage.tsx"
Cohesion: 0.15
Nodes (20): Button, ButtonProps, ButtonSize, ButtonVariant, sizeClasses, variantClasses, CheckboxProps, Field() (+12 more)

### Community 40 - "MemberPlan"
Cohesion: 0.16
Nodes (5): MemberPlan, nullIfEmpty(), scanMemberPlans(), MemberPlanRepository, MemberPlanRepository

### Community 41 - "net/http.ResponseWriter"
Cohesion: 0.13
Nodes (12): writeBookingError(), writeCustomerError(), HandleMetrics(), writeNotificationError(), writeConflictError(), writeNotFoundError(), writeReviewError(), writeTenantError() (+4 more)

### Community 42 - "Configuration"
Cohesion: 0.21
Nodes (11): workingHoursFor(), Configuration, cloneConfiguration(), defaultMemoryConfiguration(), NewConfigurationRepository(), normalizeConfigurationTenantID(), normalizePostgresTenantID(), ClosedDay (+3 more)

### Community 43 - "NewContainer"
Cohesion: 0.12
Nodes (23): Repository, NewUseCase(), NewUseCase(), UseCase, NewUseCase(), NewUseCase(), UseCase, NewUseCase() (+15 more)

### Community 44 - "BookingPage.tsx"
Cohesion: 0.13
Nodes (18): formatDateLong(), formatDateTime(), formatMoney(), formatTime(), fromDateKey(), SHORT_DAY_LABELS, toDateKey(), Confirmation() (+10 more)

### Community 45 - "useToast"
Cohesion: 0.06
Nodes (47): AdminRoutes(), AdminShell(), handleCopyBookingUrl(), AdminShellProps, AdminNavigationItem, navigationGroups, navigationItems, EstablishmentModal() (+39 more)

### Community 46 - "InventoryPage"
Cohesion: 0.31
Nodes (9): fetchServiceSupplies(), fetchStockMovements(), setServiceSupplies(), InventoryPage(), handleAddSupply(), handleDeleteProduct(), handleRemoveSupply(), handleSaveMovement() (+1 more)

### Community 47 - "AuthRepository"
Cohesion: 0.14
Nodes (6): Account, PasswordResetToken, UpdateProfileInput, cloneAuthAccount(), resetTokenExpired(), AuthRepository

### Community 48 - "UseCase"
Cohesion: 0.19
Nodes (8): appointmentDurationMinutes(), generateManagementCode(), UseCase, hashManagementCode(), BookingInput, PortalAppointment, PortalCustomer, PortalResult

### Community 49 - "net/http.Request"
Cohesion: 0.15
Nodes (8): writeAppointmentError(), resolveActorTenantID(), buildPasswordResetURL(), AuthHandler, writeAdminAuthError(), newAuthMFAHandler(), net/http.Request, AppointmentHandler

### Community 50 - "bookingApi.ts"
Cohesion: 0.09
Nodes (37): backendOrigin, buildApiUrl(), AvailabilityResponse, BookingResult, cancelPublicAppointment(), fetchAvailability(), fetchPublicAppointment(), fetchPublicInfo() (+29 more)

### Community 51 - "testing.T"
Cohesion: 0.10
Nodes (38): UseCase, newTestUseCase(), TestAvailableSlotsRespectsClosedDays(), TestCancelPublic(), TestCreateBookingGeneratesManagementCode(), TestManagementCodeNormalization(), TestPortalCancelValidatesOwnership(), TestPortalLookup() (+30 more)

### Community 52 - "Membership"
Cohesion: 0.15
Nodes (3): Membership, MembershipRepository, MembershipRepository

### Community 53 - "Package"
Cohesion: 0.15
Nodes (3): Package, PackageRepository, PackageRepository

### Community 54 - "Transaction"
Cohesion: 0.15
Nodes (6): Summary, Transaction, CategoryTotal, CategoryTotal, FinanceRepository, FinanceRepository

### Community 55 - "UseCase"
Cohesion: 0.12
Nodes (11): MFAStatus, RegisterInput, SecretCipher, Session, generateTOTPSecret(), totpCode(), totpURI(), validateTOTP() (+3 more)

### Community 56 - "PackageCredit"
Cohesion: 0.18
Nodes (4): PackageCredit, scanPackageCredits(), database/sql.Rows, PackageCreditRepository

### Community 57 - "GiftCard"
Cohesion: 0.14
Nodes (5): generateGiftCardCode(), GiftCard, NewGiftCardRepository(), GiftCardRepository, GiftCardRepository

### Community 58 - "Repository"
Cohesion: 0.27
Nodes (6): formatAppointmentLocal(), UseCase, NewReminderWorker(), NewUseCase(), ReminderWorker, Repository

### Community 59 - "BookingPage"
Cohesion: 0.25
Nodes (3): createPublicBooking(), BookingPage(), handleBookingSubmit()

### Community 60 - "UseCase"
Cohesion: 0.06
Nodes (23): AuditLog, Repository, UseCase, NewUseCase(), TestAuditLogRecordingAndListing(), UseCase, NewUseCase(), TestOutboxPublishAndDispatch() (+15 more)

### Community 61 - "checkout/usecase.go"
Cohesion: 0.60
Nodes (3): CheckoutInput, CheckoutResult, ProductItem

### Community 62 - "EventBroker"
Cohesion: 0.23
Nodes (8): NewBookingHandler(), NewCheckoutHandler(), EventBroker, NewEventBroker(), sync.RWMutex, CheckoutHandler, Event, MetricsRegistry

### Community 63 - "rate_limit.go"
Cohesion: 0.05
Nodes (43): UseCase, NewUseCase(), Log, NewNotificationRepository(), buildCompoundIdempotencyKey(), isHopByHopHeader(), NewIdempotencyStore(), TestIdempotencyBlocksInFlightConcurrentRequests() (+35 more)

### Community 64 - "AuthRepository"
Cohesion: 0.15
Nodes (3): updateAuthPasswordAndVersion(), database/sql.Tx, AuthRepository

### Community 65 - "Appointment"
Cohesion: 0.23
Nodes (4): Appointment, parseAppointmentBound(), parseAppointmentInterval(), AppointmentRepository

### Community 66 - "writeInvalidPayload"
Cohesion: 0.17
Nodes (6): writeCommerceError(), writeConfigurationError(), writeBadRequestError(), writeInvalidPayload(), writeWaitingListError(), CommerceHandler

### Community 67 - "UseCase"
Cohesion: 0.22
Nodes (3): UseCase, NewCommissionHandler(), PayoutInput

### Community 68 - "templates.go"
Cohesion: 0.27
Nodes (17): greeting(), PublicBookingOrigin(), RenderCancellation(), RenderConfirmation(), RenderReminder(), RenderReschedule(), sampleData(), TestGreetingWithoutName() (+9 more)

### Community 69 - "UseCase"
Cohesion: 0.24
Nodes (4): UseCase, NewUseCase(), NewCustomerHandler(), Repository

### Community 70 - "compilerOptions"
Cohesion: 0.09
Nodes (21): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+13 more)

### Community 71 - "AppointmentRepository"
Cohesion: 0.21
Nodes (5): appointmentSelectColumns(), scanAppointmentRow(), scanAppointmentRows(), AppointmentRepository, appointmentScanner

### Community 72 - "compilerOptions"
Cohesion: 0.09
Nodes (21): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+13 more)

### Community 73 - "compilerOptions"
Cohesion: 0.09
Nodes (21): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+13 more)

### Community 75 - "UseCase"
Cohesion: 0.28
Nodes (3): UseCase, NewExportHandler(), NewServiceHandler()

### Community 76 - "MembershipsPage"
Cohesion: 0.36
Nodes (8): fetchMemberPlans(), fetchMemberships(), MembershipsPage(), handleDelete(), handlePlanStatusChange(), handleToggleActive(), loadMemberships(), planStatusTone()

### Community 80 - "ERP Agendamentos"
Cohesion: 0.14
Nodes (13): Banco de Dados & Schema, CI, Endpoints principais, ERP Agendamentos, Estrutura, Pré-requisitos, Rodando com Docker (produção local), Rodando em desenvolvimento (sem banco) (+5 more)

### Community 82 - "UseCase"
Cohesion: 0.29
Nodes (9): UseCase, NewUseCase(), NewCommerceHandler(), repositorySet, GiftCardRepository, MemberPlanRepository, MembershipRepository, PackageCreditRepository (+1 more)

### Community 83 - "email_provider.go"
Cohesion: 0.21
Nodes (8): buildNotificationProvider(), intEnv(), logFallback(), NewEmailProvider(), NewLogProvider(), go_pkg_net_smtp, EmailProvider, LogProvider

### Community 85 - "booking/src/lib/apiError.ts"
Cohesion: 0.27
Nodes (8): ApiError, ApiErrorPayload, buildApiError(), parseApiResponse(), API_ERROR_MESSAGES, FIELD_ERROR_MESSAGES, getApiErrorMessage(), translateFieldErrors()

### Community 86 - "Claims"
Cohesion: 0.33
Nodes (5): isProductionAuthEnvironment(), Claims, decodeSegment(), encodeSegment(), JWTManager

### Community 88 - "FinancePage.tsx"
Cohesion: 0.13
Nodes (24): FinancePage, fetchAppointmentSummary(), fetchCashRegisterHistory(), fetchConfiguration(), fetchCurrentCashRegister(), fetchFinancialSummary(), fetchFinancialTransactions(), cachedAdminJson() (+16 more)

### Community 89 - ".UpdateAccountRole"
Cohesion: 0.39
Nodes (4): canManageTenantAccount(), isPrivilegedAdminRole(), isStaffRole(), normalizeTenantScope()

### Community 90 - "SecretCipher"
Cohesion: 0.38
Nodes (3): NewSecretCipher(), crypto/cipher.AEAD, SecretCipher

### Community 91 - "observability.go"
Cohesion: 0.19
Nodes (12): RecordRequestMetric(), RegisterMetricsBroker(), isSensitiveLogKey(), redactSensitiveText(), resolveRequestID(), sanitizeLogString(), sanitizeLogValue(), truncateLogString() (+4 more)

### Community 92 - "admin/vite.config.ts"
Cohesion: 0.60
Nodes (3): ref_tailwindcss_vite, ref_vite, ref_vitejs_plugin_react

### Community 93 - "appointment/repository.go"
Cohesion: 0.31
Nodes (5): DayMetric, ProfessionalMetric, ServiceMetric, Analytics, StatusCount

### Community 94 - "buildAuthComponents"
Cohesion: 0.17
Nodes (12): Repository, TokenManager, PasswordHasher, NewUseCase(), buildAuthComponents(), boolEnv(), durationSecondsEnv(), seedAdminAccount() (+4 more)

### Community 95 - "time.Time"
Cohesion: 0.29
Nodes (11): countOverlappingSlots(), formatAppointmentForTemplate(), isClosedDay(), localMidnight(), overlapsAnyBusy(), parseClock(), parseRFC3339(), ProfessionalSlots (+3 more)

### Community 96 - "TenantsPage"
Cohesion: 0.18
Nodes (11): checkSystemReadiness(), fetchAuditLogs(), fetchSystemStatus(), fetchTenants(), TenantsPage(), handleCheckReadiness(), handleDelete(), handleToggleActive() (+3 more)

### Community 98 - "ical.go"
Cohesion: 0.83
Nodes (3): BuildICalCalendar(), formatUTC(), sanitizeLine()

### Community 99 - "booking/src/App.tsx"
Cohesion: 0.20
Nodes (8): App(), booking_src_index, fetchPublicReview(), submitPublicReview(), initWidgetBridge(), ReviewPage(), handleSubmit(), SlugEntryPage()

### Community 110 - "context.Context"
Cohesion: 0.31
Nodes (6): MaintenanceWorker, MaintenanceConfig, UseCase, NewMaintenanceWorker(), context.Context, NopTxRunner

### Community 115 - "PortalPage.tsx"
Cohesion: 0.25
Nodes (6): PortalAppointment, PortalResult, Slot, ACTIONABLE_STATUSES, PortalStep, STATUS_LABELS

### Community 116 - "go_pkg_context"
Cohesion: 0.13
Nodes (11): main(), resolveDurationEnv(), resolvePositiveIntEnv(), NewPostgresTxRunner(), go_pkg_context, go_pkg_erpagendamentos_backend_internal_domain_payment, go_pkg_io, go_pkg_os_signal (+3 more)

### Community 118 - "TenantSubscriptionModal"
Cohesion: 0.40
Nodes (3): updateTenantSubscription(), TenantSubscriptionModal(), handleSubmit()

### Community 122 - "body_limit.go"
Cohesion: 0.83
Nodes (3): resolveMaxRequestBodyBytes(), shouldLimitRequestBody(), withRequestBodyLimit()

## Knowledge Gaps
- **257 isolated node(s):** `graphify-mcp`, `name`, `private`, `version`, `type` (+252 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **19 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Container` connect `RouterConfig` to `Payment`, `go_pkg_net_http`, `Review`, `UseCase`, `UseCase`, `Tenant`, `database/sql.DB`, `Resource`, `CashRegister`, `Professional`, `Entry`, `Repository`, `UseCase`, `UseCase`, `Queue`, `NewContainer`, `UseCase`, `UseCase`, `Repository`, `UseCase`, `EventBroker`, `rate_limit.go`, `UseCase`, `UseCase`, `UseCase`, `UseCase`, `buildAuthComponents`, `context.Context`?**
  _High betweenness centrality (0.054) - this node is a cross-community bridge._
- **Why does `UseCase` connect `UseCase` to `auth/usecase.go`, `writeJSON`, `writeAPIError`, `NewContainer`, `time.Duration`, `AuthRepository`, `DeviceSession`, `net/http.Request`, `.Add`, `RouterConfig`, `Claims`, `.UpdateAccountRole`, `rate_limit.go`, `buildAuthComponents`, `time.Time`?**
  _High betweenness centrality (0.032) - this node is a cross-community bridge._
- **Why does `RouterConfig` connect `RouterConfig` to `Payment`, `go_pkg_net_http`, `Review`, `UseCase`, `UseCase`, `Tenant`, `Resource`, `CashRegister`, `Professional`, `Entry`, `Repository`, `UseCase`, `UseCase`, `NewContainer`, `UseCase`, `UseCase`, `Repository`, `UseCase`, `EventBroker`, `rate_limit.go`, `UseCase`, `UseCase`, `UseCase`, `time.Duration`, `UseCase`, `context.Context`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Are the 131 inferred relationships involving `writeJSON()` (e.g. with `health()` and `readiness()`) actually correct?**
  _`writeJSON()` has 131 INFERRED edges - model-reasoned connections that need verification._
- **Are the 104 inferred relationships involving `resolveRequestTenantID()` (e.g. with `.Analytics()` and `.CheckIn()`) actually correct?**
  _`resolveRequestTenantID()` has 104 INFERRED edges - model-reasoned connections that need verification._
- **What connects `graphify-mcp`, `name`, `private` to the rest of the system?**
  _257 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `adminFetch` be split into smaller, more focused modules?**
  _Cohesion score 0.05030864197530864 - nodes in this community are weakly interconnected._
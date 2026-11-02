# Graph Report - erpagendamentos  (2026-11-02)

## Corpus Check
- 285 files · ~138,044 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2748 nodes · 8250 edges · 123 communities (107 shown, 16 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 696 edges (avg confidence: 0.79)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `5a44f6b7`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- adminFetch
- schema.sql
- ref_react
- Payment
- writeJSON
- go_pkg_net_http
- container.go
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
- CommissionRepository
- InventoryPage
- Resource
- CashRegister
- MembershipsPage
- Queue
- Entry
- devDependencies
- devDependencies
- UseCase
- Configuration
- auth_handler.go
- UseCase
- rate_limit.go
- devDependencies
- ManagePage.tsx
- MemberPlan
- net/http.ResponseWriter
- time.Time
- Container
- BookingPage.tsx
- useToast
- writeProfessionalError
- Account
- UseCase
- AppointmentHandler
- bookingApi.ts
- memoryRepositories
- Membership
- Package
- Transaction
- net/http.Request
- PackageCredit
- GiftCard
- Professional
- AuthRepository
- AuditLog
- go_pkg_errors
- EventBroker
- UseCase
- NewContainer
- Appointment
- writeInvalidPayload
- UseCase
- templates.go
- hashResetToken
- compilerOptions
- AppointmentRepository
- compilerOptions
- compilerOptions
- OutboxRepository
- UseCase
- Session
- writeResourceError
- testing.T
- DeviceSession
- ERP Agendamentos
- NewBcryptHasher
- UseCase
- email_provider.go
- UseCase
- booking/src/lib/apiError.ts
- time.Duration
- ExportHandler
- FinancePage.tsx
- Claims
- SecretCipher
- observability.go
- admin/vite.config.ts
- appointment/repository.go
- UseCase
- PackagesPage
- adminHttp.ts
- TenantSubscriptionModal
- ical.go
- booking/src/App.tsx
- graphify
- rules/graphify.md
- workflows/graphify.md
- generate-certs.sh
- backup_postgres.sh
- certbot_ssl.sh
- catalog.finance_transactions
- GiftCardsPage
- MaintenanceWorker
- go_pkg_path_filepath
- erpagendamentos/backend
- writeServiceError
- PortalPage.tsx
- calendar_sync_handler.go
- widget_handler.go
- createPackageCredit
- ProfessionalsPage
- ResourcesPage
- main.go
- ref_react_dom

## God Nodes (most connected - your core abstractions)
1. `writeJSON()` - 137 edges
2. `resolveRequestTenantID()` - 109 edges
3. `adminFetch()` - 78 edges
4. `useToast()` - 77 edges
5. `writeInvalidPayload()` - 73 edges
6. `UseCase` - 72 edges
7. `invalidateAdminGetCache()` - 65 edges
8. `writeInternalError()` - 58 edges
9. `RouterConfig` - 56 edges
10. `writeAPIError()` - 48 edges

## Surprising Connections (you probably didn't know these)
- `NewContainer()` --calls--> `loadEnvFile()`  [INFERRED]
  backend/internal/bootstrap/container.go → backend/internal/bootstrap/env.go
- `handleSubmit()` --calls--> `updateCurrentTenant()`  [EXTRACTED]
  admin/src/components/EstablishmentModal.tsx → admin/src/lib/adminApi.ts
- `handleSubmit()` --calls--> `updateConfiguration()`  [EXTRACTED]
  admin/src/pages/ConfigurationPage.tsx → admin/src/lib/adminApi.ts
- `handleSubmit()` --calls--> `updateCurrentTenant()`  [EXTRACTED]
  admin/src/pages/EstablishmentPage.tsx → admin/src/lib/adminApi.ts
- `newResourceUseCase()` --calls--> `NewUseCase()`  [INFERRED]
  backend/internal/application/appointment/resource_test.go → backend/internal/application/appointment/usecase.go

## Import Cycles
- None detected.

## Communities (123 total, 16 thin omitted)

### Community 0 - "adminFetch"
Cohesion: 0.05
Nodes (86): closeCashRegister(), connectCalendar(), createCashRegisterEntry(), createCommissionRule(), createCustomer(), createFinancialTransaction(), createGiftCard(), createMemberPlan() (+78 more)

### Community 1 - "schema.sql"
Cohesion: 0.12
Nodes (35): catalog.appointments, catalog.audit_logs, catalog.auth_rate_limits, catalog.auth_sessions, catalog.calendar_connections, catalog.cash_register_entries, catalog.cash_registers, catalog.commission_payouts (+27 more)

### Community 2 - "ref_react"
Cohesion: 0.08
Nodes (85): AdminRoutes(), AdminRoutesProps, AgendaPage, AnalyticsPage, AppointmentsPage, CommissionsPage, ConfigurationPage, CustomersPage (+77 more)

### Community 3 - "Payment"
Cohesion: 0.06
Nodes (26): UseCase, NewUseCase(), normalizeStatus(), buildPaymentProvider(), IntentResult, Payment, NewMercadoPagoProvider(), NewMockProvider() (+18 more)

### Community 4 - "writeJSON"
Cohesion: 0.10
Nodes (11): resolveRequestTenantID(), health(), writePaymentError(), writeInternalError(), writeJSON(), CalendarSyncHandler, CashRegisterHandler, CommissionHandler (+3 more)

### Community 5 - "go_pkg_net_http"
Cohesion: 0.07
Nodes (26): loadEnvFile(), appendQueryParam(), fetchGoogleUserInfo(), go_pkg_bufio, go_pkg_encoding_json, go_pkg_erpagendamentos_backend_internal_application_booking, go_pkg_erpagendamentos_backend_internal_application_cashregister, go_pkg_erpagendamentos_backend_internal_application_commission (+18 more)

### Community 6 - "container.go"
Cohesion: 0.11
Nodes (22): deploymentEnvironment(), isProductionEnv(), go_pkg_context, go_pkg_erpagendamentos_backend_internal_application_checkout, go_pkg_erpagendamentos_backend_internal_application_commerce, go_pkg_erpagendamentos_backend_internal_application_configuration, go_pkg_erpagendamentos_backend_internal_application_customer, go_pkg_erpagendamentos_backend_internal_application_professional (+14 more)

### Community 7 - "downloadAdminCsv"
Cohesion: 0.11
Nodes (21): downloadAdminCsv(), fetchPayments(), fetchReviews(), handleExport(), CustomersPage(), handleDelete(), handleExport(), handleSearch() (+13 more)

### Community 8 - "adminAuth.ts"
Cohesion: 0.06
Nodes (52): App(), handleLogout(), AppProviders(), AppProvidersProps, ResetPasswordPage(), handleSubmit(), AppRouter(), admin_src_index (+44 more)

### Community 9 - "writeAPIError"
Cohesion: 0.11
Nodes (21): resolveAuthClaims(), browserSessionPayload(), exposeAuthTokenForTests(), AuthHandler, writeAdminAuthError(), AuthHandler, readiness(), logError() (+13 more)

### Community 10 - "AppointmentsPage.tsx"
Cohesion: 0.07
Nodes (54): AppointmentFormModal(), computeEndsAt(), handleSubmit(), AppointmentFormModalProps, STATUS_OPTIONS, CheckoutComandaModal(), handleSubmit(), checkInAppointment() (+46 more)

### Community 11 - "adminApi.ts"
Cohesion: 0.05
Nodes (71): CheckoutComandaModalProps, SelectedProduct, AppointmentFilters, checkoutAppointment(), fetchCommissionPayouts(), fetchCommissionRecords(), fetchCommissionRules(), fetchCommissionSummaries() (+63 more)

### Community 12 - "Review"
Cohesion: 0.14
Nodes (12): UseCase, NewUseCase(), Review, NewReviewRepository(), scanReviewRow(), scanReviewRows(), NewReviewHandler(), ReviewRepository (+4 more)

### Community 13 - "UseCase"
Cohesion: 0.08
Nodes (10): UseCase, NewUseCase(), TestStockUseCase(), Product, Repository, ServiceSupply, StockMovement, NewStockHandler() (+2 more)

### Community 14 - "Repository"
Cohesion: 0.16
Nodes (10): formatAppointmentLocal(), UseCase, NewReminderWorker(), NewPostgresTxRunner(), context.Context, ReminderWorker, NopTxRunner, PostgresTxRunner (+2 more)

### Community 15 - "Tenant"
Cohesion: 0.11
Nodes (12): UseCase, Tenant, UpdateSubscriptionInput, defaultModules(), isTenantIdentityConflict(), defaultModules(), scanTenantRow(), scanTenantRows() (+4 more)

### Community 16 - "landing/src/App.tsx"
Cohesion: 0.10
Nodes (22): App(), CtaSection(), FEATURES, FeaturesSection(), Footer(), Header(), NAV_ITEMS, HeroSection() (+14 more)

### Community 17 - "Service"
Cohesion: 0.11
Nodes (9): UseCase, NewUseCase(), Service, scanServiceRow(), scanServiceRows(), NewServiceHandler(), ServiceRepository, ServiceRepository (+1 more)

### Community 18 - "go_pkg_time"
Cohesion: 0.12
Nodes (11): go_pkg_database_sql, go_pkg_erpagendamentos_backend_internal_domain_job, go_pkg_erpagendamentos_backend_internal_domain_review, go_pkg_erpagendamentos_backend_internal_domain_waitinglist, go_pkg_fmt, go_pkg_github_com_google_uuid, go_pkg_github_com_jackc_pgx_v5_stdlib, go_pkg_sort (+3 more)

### Community 19 - "auth/entities.go"
Cohesion: 0.18
Nodes (10): ChangePasswordInput, LoginInput, PasswordResetCodeConfirmInput, PasswordResetConfirmInput, RegisterInput, ValidationError, NewValidationError(), passwordResetCodeToken() (+2 more)

### Community 20 - "RouterConfig"
Cohesion: 0.06
Nodes (72): lazyAppHandler, writeLazyJSON(), bearerAuthAllowed(), hasRole(), isAdminRole(), isPrivilegedAdminRole(), requestWithAuthClaims(), requireAdmin() (+64 more)

### Community 21 - "database/sql.DB"
Cohesion: 0.12
Nodes (33): applyDatabaseSQL(), resolveSchemaFile(), initializeRepositories(), postgresRepositories(), NewAppointmentRepository(), NewAuditRepository(), NewAuthRepository(), NewCalendarRepository() (+25 more)

### Community 22 - "resource_test.go"
Cohesion: 0.54
Nodes (7): UseCase, newAppointment(), newResourceUseCase(), TestCreateAutoAssignsAvailableResource(), TestCreateRejectsBusyExplicitResource(), TestCreateRejectsUnknownResource(), TestUpdateAutoReassignsFreedResource()

### Community 23 - "Customer"
Cohesion: 0.11
Nodes (13): UseCase, NewUseCase(), Customer, NewCustomerRepository(), onlyDigits(), scanCustomerRow(), scanCustomerRows(), NewCustomerHandler() (+5 more)

### Community 24 - "CommissionRepository"
Cohesion: 0.12
Nodes (5): CommissionRecord, CommissionRule, Payout, CommissionRepository, CommissionRepository

### Community 25 - "InventoryPage"
Cohesion: 0.31
Nodes (9): fetchServiceSupplies(), fetchStockMovements(), setServiceSupplies(), InventoryPage(), handleAddSupply(), handleDeleteProduct(), handleRemoveSupply(), handleSaveMovement() (+1 more)

### Community 26 - "Resource"
Cohesion: 0.13
Nodes (9): UseCase, Resource, NewResourceRepository(), scanResourceRow(), NewExportHandler(), NewResourceHandler(), ResourceRepository, ResourceRepository (+1 more)

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
Nodes (6): RecurrenceInput, UseCase, parseRFC3339(), validStatusChange(), NewCalendarHandler(), CalendarHandler

### Community 34 - "Configuration"
Cohesion: 0.24
Nodes (9): Configuration, cloneConfiguration(), defaultMemoryConfiguration(), normalizeConfigurationTenantID(), normalizePostgresTenantID(), ClosedDay, WorkingHours, ConfigurationRepository (+1 more)

### Community 35 - "auth_handler.go"
Cohesion: 0.09
Nodes (18): secureTemporaryPassword(), go_pkg_crypto_aes, go_pkg_crypto_cipher, go_pkg_crypto_hmac, go_pkg_crypto_rand, go_pkg_crypto_rsa, go_pkg_crypto_sha1, go_pkg_crypto_sha256 (+10 more)

### Community 36 - "UseCase"
Cohesion: 0.12
Nodes (8): UseCase, NewUseCase(), Connection, Repository, SyncEvent, NewCalendarSyncHandler(), CalendarRepository, CalendarRepository

### Community 37 - "rate_limit.go"
Cohesion: 0.05
Nodes (45): UseCase, NewUseCase(), Log, NewNotificationRepository(), buildCompoundIdempotencyKey(), isHopByHopHeader(), NewIdempotencyStore(), TestIdempotencyBlocksInFlightConcurrentRequests() (+37 more)

### Community 38 - "devDependencies"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, devDependencies, tailwindcss, @tailwindcss/vite, @types/node (+24 more)

### Community 39 - "ManagePage.tsx"
Cohesion: 0.14
Nodes (24): Professional, PublicAppointment, Service, ManageStep, Button, ButtonProps, ButtonSize, ButtonVariant (+16 more)

### Community 40 - "MemberPlan"
Cohesion: 0.14
Nodes (7): MemberPlan, NewMemberPlanRepository(), nullIfEmpty(), scanMemberPlans(), database/sql.Rows, MemberPlanRepository, MemberPlanRepository

### Community 41 - "net/http.ResponseWriter"
Cohesion: 0.12
Nodes (11): writeCustomerError(), HandleMetrics(), writeNotificationError(), writeNotFoundError(), writeRawJSON(), writeReviewError(), writeTenantError(), writeWaitingListError() (+3 more)

### Community 42 - "time.Time"
Cohesion: 0.27
Nodes (12): countOverlappingSlots(), formatAppointmentForTemplate(), isClosedDay(), localMidnight(), overlapsAnyBusy(), parseClock(), parseRFC3339(), workingHoursFor() (+4 more)

### Community 43 - "Container"
Cohesion: 0.15
Nodes (11): UseCase, NewUseCase(), UseCase, Container, NewAuthGoogleHandler(), NewConfigurationHandler(), NewOnboardingHandler(), Repository (+3 more)

### Community 44 - "BookingPage.tsx"
Cohesion: 0.13
Nodes (19): formatDateLong(), formatDateTime(), formatMoney(), formatTime(), fromDateKey(), nextDateKeys(), SHORT_DAY_LABELS, toDateKey() (+11 more)

### Community 45 - "useToast"
Cohesion: 0.05
Nodes (53): AdminShell(), handleCopyBookingUrl(), AdminShellProps, AdminNavigationItem, navigationGroups, navigationItems, EstablishmentModal(), handleCopyLink() (+45 more)

### Community 47 - "Account"
Cohesion: 0.15
Nodes (6): Account, PasswordResetToken, UpdateProfileInput, cloneAuthAccount(), resetTokenExpired(), AuthRepository

### Community 48 - "UseCase"
Cohesion: 0.19
Nodes (8): appointmentDurationMinutes(), generateManagementCode(), UseCase, hashManagementCode(), BookingInput, PortalAppointment, PortalCustomer, PortalResult

### Community 49 - "AppointmentHandler"
Cohesion: 0.26
Nodes (3): writeAppointmentError(), resolveActorTenantID(), AppointmentHandler

### Community 50 - "bookingApi.ts"
Cohesion: 0.10
Nodes (30): backendOrigin, buildApiUrl(), AvailabilityResponse, BookingResult, cancelPublicAppointment(), createPublicBooking(), fetchAvailability(), fetchPublicAppointment() (+22 more)

### Community 51 - "memoryRepositories"
Cohesion: 0.25
Nodes (17): TestCashRegisterFlow(), TestCheckoutFlow(), TestCommissionUseCase(), TestOnboardingSignup(), memoryRepositories(), NewAppointmentRepository(), NewAuthRepository(), NewCashRegisterRepository() (+9 more)

### Community 52 - "Membership"
Cohesion: 0.16
Nodes (4): Membership, NewMembershipRepository(), MembershipRepository, MembershipRepository

### Community 53 - "Package"
Cohesion: 0.15
Nodes (4): Package, NewPackageRepository(), PackageRepository, PackageRepository

### Community 54 - "Transaction"
Cohesion: 0.15
Nodes (6): Summary, Transaction, CategoryTotal, CategoryTotal, FinanceRepository, FinanceRepository

### Community 55 - "net/http.Request"
Cohesion: 0.14
Nodes (9): resolveTenantID(), resolveFrontendURL(), resolveGoogleRedirectURI(), buildPasswordResetURL(), newAuthMFAHandler(), writeBookingError(), net/http.Request, AuthGoogleHandler (+1 more)

### Community 56 - "PackageCredit"
Cohesion: 0.18
Nodes (5): PackageCredit, NewPackageCreditRepository(), scanPackageCredits(), PackageCreditRepository, PackageCreditRepository

### Community 57 - "GiftCard"
Cohesion: 0.14
Nodes (5): generateGiftCardCode(), GiftCard, NewGiftCardRepository(), GiftCardRepository, GiftCardRepository

### Community 58 - "Professional"
Cohesion: 0.14
Nodes (8): UseCase, Professional, scanProfessionalRow(), scanProfessionalRows(), NewProfessionalHandler(), ProfessionalRepository, ProfessionalRepository, professionalScanner

### Community 59 - "AuthRepository"
Cohesion: 0.11
Nodes (4): AuthState, updateAuthPasswordAndVersion(), database/sql.Tx, AuthRepository

### Community 60 - "AuditLog"
Cohesion: 0.16
Nodes (10): AuditLog, Repository, UseCase, NewUseCase(), TestAuditLogRecordingAndListing(), NewAuditRepository(), NewAuditHandler(), AuditRepository (+2 more)

### Community 61 - "go_pkg_errors"
Cohesion: 0.15
Nodes (6): CheckoutInput, CheckoutResult, ProductItem, go_pkg_erpagendamentos_backend_internal_domain_cashregister, go_pkg_erpagendamentos_backend_internal_domain_finance, go_pkg_errors

### Community 62 - "EventBroker"
Cohesion: 0.22
Nodes (8): NewAppointmentHandler(), NewBookingHandler(), EventBroker, NewEventBroker(), RegisterMetricsBroker(), sync.RWMutex, Event, MetricsRegistry

### Community 63 - "UseCase"
Cohesion: 0.14
Nodes (6): UseCase, NewUseCase(), Repository, Summary, NewCommissionHandler(), PayoutInput

### Community 64 - "NewContainer"
Cohesion: 0.22
Nodes (13): Repository, NewUseCase(), NewUseCase(), NewUseCase(), NewUseCase(), NewUseCase(), NewUseCase(), NewUseCase() (+5 more)

### Community 65 - "Appointment"
Cohesion: 0.19
Nodes (5): Appointment, StatusCount, parseAppointmentBound(), parseAppointmentInterval(), AppointmentRepository

### Community 66 - "writeInvalidPayload"
Cohesion: 0.20
Nodes (5): writeCommerceError(), writeConfigurationError(), writeBadRequestError(), writeInvalidPayload(), CommerceHandler

### Community 67 - "UseCase"
Cohesion: 0.19
Nodes (7): UseCase, NewUseCase(), TestFinanceUseCase(), Repository, NewFinanceHandler(), NewMercadoPagoWebhookHandler(), MercadoPagoWebhookHandler

### Community 68 - "templates.go"
Cohesion: 0.27
Nodes (17): greeting(), PublicBookingOrigin(), RenderCancellation(), RenderConfirmation(), RenderReminder(), RenderReschedule(), sampleData(), TestGreetingWithoutName() (+9 more)

### Community 69 - "hashResetToken"
Cohesion: 0.29
Nodes (3): PasswordResetRequestInput, hashResetToken(), passwordResetCodeChallengeKey()

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
Cohesion: 0.20
Nodes (10): UseCase, NewUseCase(), TestOutboxPublishAndDispatch(), NewOutboxRepository(), NewSystemHandler(), TestAuditAndSystemHandlers(), go_pkg_runtime, EventDispatcher (+2 more)

### Community 76 - "Session"
Cohesion: 0.17
Nodes (7): MFAStatus, Session, generateTOTPSecret(), totpCode(), totpURI(), validateTOTP(), isMFAEligibleRole()

### Community 78 - "testing.T"
Cohesion: 0.23
Nodes (15): UseCase, newTestUseCase(), TestAvailableSlotsRespectsClosedDays(), TestCancelPublic(), TestCreateBookingGeneratesManagementCode(), TestManagementCodeNormalization(), TestPortalCancelValidatesOwnership(), TestPortalLookup() (+7 more)

### Community 80 - "ERP Agendamentos"
Cohesion: 0.14
Nodes (13): Banco de Dados & Schema, CI, Endpoints principais, ERP Agendamentos, Estrutura, Pré-requisitos, Rodando com Docker (produção local), Rodando em desenvolvimento (sem banco) (+5 more)

### Community 81 - "NewBcryptHasher"
Cohesion: 0.40
Nodes (3): NewBcryptHasher(), go_pkg_golang_org_x_crypto_bcrypt, BcryptHasher

### Community 82 - "UseCase"
Cohesion: 0.21
Nodes (9): UseCase, NewUseCase(), NewCommerceHandler(), repositorySet, GiftCardRepository, MemberPlanRepository, MembershipRepository, PackageCreditRepository (+1 more)

### Community 83 - "email_provider.go"
Cohesion: 0.21
Nodes (8): buildNotificationProvider(), intEnv(), logFallback(), NewEmailProvider(), NewLogProvider(), go_pkg_net_smtp, EmailProvider, LogProvider

### Community 84 - "UseCase"
Cohesion: 0.33
Nodes (6): UseCase, NewCheckoutHandler(), AuditLogger, OutboxPublisher, TxRunner, CheckoutHandler

### Community 85 - "booking/src/lib/apiError.ts"
Cohesion: 0.24
Nodes (10): ApiError, ApiErrorPayload, buildApiError(), parseApiResponse(), API_ERROR_MESSAGES, EXACT_MESSAGE_MAP, getApiErrorMessage(), SUBSTRING_MESSAGE_PATTERNS (+2 more)

### Community 86 - "time.Duration"
Cohesion: 0.24
Nodes (4): boolEnv(), durationSecondsEnv(), CleanupSummary, time.Duration

### Community 88 - "FinancePage.tsx"
Cohesion: 0.14
Nodes (23): FinancePage, fetchCashRegisterHistory(), fetchConfiguration(), fetchCurrentCashRegister(), fetchFinancialSummary(), fetchFinancialTransactions(), cachedAdminJson(), ActiveTab (+15 more)

### Community 89 - "Claims"
Cohesion: 0.48
Nodes (4): Claims, decodeSegment(), encodeSegment(), JWTManager

### Community 90 - "SecretCipher"
Cohesion: 0.38
Nodes (3): NewSecretCipher(), crypto/cipher.AEAD, SecretCipher

### Community 91 - "observability.go"
Cohesion: 0.19
Nodes (11): RecordRequestMetric(), isSensitiveLogKey(), redactSensitiveText(), resolveRequestID(), sanitizeLogString(), sanitizeLogValue(), truncateLogString(), withObservability() (+3 more)

### Community 92 - "admin/vite.config.ts"
Cohesion: 0.60
Nodes (3): ref_tailwindcss_vite, ref_vite, ref_vitejs_plugin_react

### Community 93 - "appointment/repository.go"
Cohesion: 0.53
Nodes (4): DayMetric, ProfessionalMetric, ServiceMetric, Analytics

### Community 94 - "UseCase"
Cohesion: 0.15
Nodes (13): Repository, SecretCipher, TokenManager, canManageTenantAccount(), PasswordHasher, UseCase, isPrivilegedAdminRole(), isProductionAuthEnvironment() (+5 more)

### Community 95 - "PackagesPage"
Cohesion: 0.47
Nodes (6): fetchPackageCredits(), fetchPackages(), PackagesPage(), handleDelete(), handleToggleActive(), loadPackages()

### Community 96 - "adminHttp.ts"
Cohesion: 0.13
Nodes (15): checkSystemReadiness(), fetchAuditLogs(), fetchSystemStatus(), fetchTenants(), adminGetCache, AdminMeResponse, parseJson(), backendOrigin (+7 more)

### Community 97 - "TenantSubscriptionModal"
Cohesion: 0.40
Nodes (3): updateTenantSubscription(), TenantSubscriptionModal(), handleSubmit()

### Community 98 - "ical.go"
Cohesion: 0.83
Nodes (3): BuildICalCalendar(), formatUTC(), sanitizeLine()

### Community 99 - "booking/src/App.tsx"
Cohesion: 0.16
Nodes (5): App(), booking_src_index, initWidgetBridge(), BookingPage(), SlugEntryPage()

### Community 109 - "GiftCardsPage"
Cohesion: 0.50
Nodes (5): fetchGiftCards(), GiftCardsPage(), handleDelete(), loadGiftCards(), statusTone()

### Community 110 - "MaintenanceWorker"
Cohesion: 0.39
Nodes (5): MaintenanceWorker, MaintenanceConfig, UseCase, NewMaintenanceWorker(), authComponents

### Community 115 - "PortalPage.tsx"
Cohesion: 0.16
Nodes (14): PortalAppointment, portalCancelAppointment(), portalLookup(), portalRescheduleAppointment(), PortalResult, Slot, ACTIONABLE_STATUSES, PortalPage() (+6 more)

### Community 117 - "widget_handler.go"
Cohesion: 0.60
Nodes (3): NewWidgetHandler(), resolveWidgetBookingOrigin(), WidgetHandler

### Community 118 - "createPackageCredit"
Cohesion: 0.67
Nodes (3): createPackageCredit(), AssignCreditModal(), handleSubmit()

### Community 119 - "ProfessionalsPage"
Cohesion: 0.83
Nodes (4): ProfessionalsPage(), handleDelete(), handleToggleActive(), loadProfessionals()

### Community 120 - "ResourcesPage"
Cohesion: 0.83
Nodes (4): ResourcesPage(), handleDelete(), handleToggleActive(), loadResources()

### Community 122 - "main.go"
Cohesion: 0.13
Nodes (12): newLazyAppHandler(), main(), resolveDurationEnv(), resolvePositiveIntEnv(), resolveMaxRequestBodyBytes(), shouldLimitRequestBody(), withRequestBodyLimit(), go_pkg_erpagendamentos_backend_internal_application_stock (+4 more)

## Knowledge Gaps
- **263 isolated node(s):** `graphify-mcp`, `name`, `private`, `version`, `type` (+258 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **16 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Container` connect `Container` to `Payment`, `container.go`, `Review`, `UseCase`, `Repository`, `Tenant`, `Service`, `RouterConfig`, `database/sql.DB`, `Customer`, `Resource`, `CashRegister`, `Queue`, `Entry`, `UseCase`, `UseCase`, `rate_limit.go`, `UseCase`, `Professional`, `AuditLog`, `EventBroker`, `UseCase`, `NewContainer`, `UseCase`, `UseCase`, `UseCase`, `UseCase`, `time.Duration`, `UseCase`, `MaintenanceWorker`?**
  _High betweenness centrality (0.047) - this node is a cross-community bridge._
- **Why does `UseCase` connect `UseCase` to `NewContainer`, `writeJSON`, `hashResetToken`, `rate_limit.go`, `writeAPIError`, `time.Time`, `Container`, `Session`, `MaintenanceWorker`, `DeviceSession`, `Account`, `AppointmentHandler`, `auth/entities.go`, `RouterConfig`, `time.Duration`, `net/http.Request`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Why does `RouterConfig` connect `RouterConfig` to `Payment`, `container.go`, `Review`, `UseCase`, `Repository`, `Tenant`, `Service`, `Customer`, `Resource`, `CashRegister`, `Entry`, `UseCase`, `UseCase`, `rate_limit.go`, `Container`, `UseCase`, `Professional`, `AuditLog`, `EventBroker`, `UseCase`, `UseCase`, `UseCase`, `UseCase`, `UseCase`, `time.Duration`, `UseCase`, `MaintenanceWorker`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Are the 133 inferred relationships involving `writeJSON()` (e.g. with `health()` and `readiness()`) actually correct?**
  _`writeJSON()` has 133 INFERRED edges - model-reasoned connections that need verification._
- **Are the 104 inferred relationships involving `resolveRequestTenantID()` (e.g. with `.Analytics()` and `.CheckIn()`) actually correct?**
  _`resolveRequestTenantID()` has 104 INFERRED edges - model-reasoned connections that need verification._
- **What connects `graphify-mcp`, `name`, `private` to the rest of the system?**
  _263 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `adminFetch` be split into smaller, more focused modules?**
  _Cohesion score 0.047058823529411764 - nodes in this community are weakly interconnected._
# Graph Report - erpagendamentos  (2026-11-02)

## Corpus Check
- 285 files · ~138,113 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2749 nodes · 8249 edges · 118 communities (101 shown, 17 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 695 edges (avg confidence: 0.79)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `3e8af476`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- adminApi.ts
- schema.sql
- ref_react
- Payment
- net/http.Request
- go_pkg_net_http
- withIdempotency
- downloadAdminCsv
- adminAuth.ts
- writeAPIError
- AppointmentsPage.tsx
- types.ts
- Review
- UseCase
- context.Context
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
- InventoryPage
- UseCase
- CashRegister
- MembershipsPage
- JobRepository
- Entry
- devDependencies
- devDependencies
- UseCase
- memory/configuration_repository.go
- booking/usecase.go
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
- writeJSON
- PackageCredit
- GiftCard
- Professional
- AuthRepository
- AuditLog
- checkout/usecase.go
- EventBroker
- auth_google_handler.go
- Queue
- Appointment
- writeInvalidPayload
- UseCase
- templates.go
- .UpdateAccountRole
- compilerOptions
- AppointmentRepository
- compilerOptions
- compilerOptions
- OutboxRepository
- UseCase
- totp.go
- writeResourceError
- testing.T
- lazyAppHandler
- ERP Agendamentos
- TestAuthGoogleHandlerStatusAndLogin
- UseCase
- main.go
- UseCase
- booking/src/lib/apiError.ts
- time.Duration
- ExportHandler
- writePaymentError
- Claims
- buildAuthComponents
- UseCase
- admin/vite.config.ts
- appointment/repository.go
- UseCase
- PackagesPage
- TenantsPage
- connection.go
- ical.go
- booking/src/App.tsx
- graphify
- rules/graphify.md
- workflows/graphify.md
- generate-certs.sh
- backup_postgres.sh
- certbot_ssl.sh
- catalog.finance_transactions
- WaitingListPage
- go_pkg_path_filepath
- erpagendamentos/backend
- writeServiceError
- PortalPage.tsx
- CommissionsPage.tsx
- go_pkg_os
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
- `handleSubmit()` --calls--> `updateTenantSubscription()`  [EXTRACTED]
  admin/src/pages/TenantsPage.tsx → admin/src/lib/adminApi.ts

## Import Cycles
- None detected.

## Communities (118 total, 17 thin omitted)

### Community 0 - "adminApi.ts"
Cohesion: 0.05
Nodes (106): AppointmentFilters, checkInAppointment(), checkoutAppointment(), checkSystemReadiness(), closeCashRegister(), connectCalendar(), createAppointment(), createCashRegisterEntry() (+98 more)

### Community 1 - "schema.sql"
Cohesion: 0.12
Nodes (35): catalog.appointments, catalog.audit_logs, catalog.auth_rate_limits, catalog.auth_sessions, catalog.calendar_connections, catalog.cash_register_entries, catalog.cash_registers, catalog.commission_payouts (+27 more)

### Community 2 - "ref_react"
Cohesion: 0.07
Nodes (89): AdminRoutesProps, AgendaPage, AnalyticsPage, AppointmentsPage, ConfigurationPage, CustomersPage, EstablishmentPage, GiftCardsPage (+81 more)

### Community 3 - "Payment"
Cohesion: 0.06
Nodes (25): UseCase, NewUseCase(), normalizeStatus(), buildPaymentProvider(), IntentResult, Payment, NewMercadoPagoProvider(), NewMockProvider() (+17 more)

### Community 4 - "net/http.Request"
Cohesion: 0.10
Nodes (10): resolveRequestTenantID(), resolveTenantID(), buildPasswordResetURL(), HandleMetrics(), writeInternalError(), net/http.Request, CalendarSyncHandler, CommissionHandler (+2 more)

### Community 5 - "go_pkg_net_http"
Cohesion: 0.10
Nodes (26): deploymentEnvironment(), isProductionEnv(), go_pkg_context, go_pkg_encoding_json, go_pkg_erpagendamentos_backend_internal_application_booking, go_pkg_erpagendamentos_backend_internal_application_calendar, go_pkg_erpagendamentos_backend_internal_application_cashregister, go_pkg_erpagendamentos_backend_internal_application_checkout (+18 more)

### Community 6 - "withIdempotency"
Cohesion: 0.15
Nodes (12): isHopByHopHeader(), NewIdempotencyStore(), TestIdempotencyBlocksInFlightConcurrentRequests(), TestIdempotencyDoesNotCacheServerFailures(), TestIdempotencyIgnoredOnGetRequests(), TestIdempotencyReplaysCompletedResponse(), withIdempotency(), bytes.Buffer (+4 more)

### Community 7 - "downloadAdminCsv"
Cohesion: 0.11
Nodes (22): downloadAdminCsv(), fetchCustomers(), fetchPayments(), fetchReviews(), handleExport(), CustomersPage(), handleDelete(), handleExport() (+14 more)

### Community 8 - "adminAuth.ts"
Cohesion: 0.05
Nodes (55): App(), handleLogout(), AppProviders(), AppProvidersProps, ResetPasswordPage(), handleSubmit(), AppRouter(), admin_src_index (+47 more)

### Community 9 - "writeAPIError"
Cohesion: 0.10
Nodes (18): resolveAuthClaims(), browserSessionPayload(), exposeAuthTokenForTests(), AuthHandler, newAuthMFAHandler(), AuthHandler, readiness(), logError() (+10 more)

### Community 10 - "AppointmentsPage.tsx"
Cohesion: 0.06
Nodes (56): AppointmentFormModal(), computeEndsAt(), handleSubmit(), AppointmentFormModalProps, STATUS_OPTIONS, CheckoutComandaModal(), handleSubmit(), CheckoutComandaModalProps (+48 more)

### Community 11 - "types.ts"
Cohesion: 0.05
Nodes (54): FinancePage, fetchCashRegisterHistory(), fetchFinancialTransactions(), ActiveTab, FinancePage(), handleCloseShift(), handleCreate(), handleCreateEntry() (+46 more)

### Community 12 - "Review"
Cohesion: 0.14
Nodes (12): UseCase, NewUseCase(), Review, NewReviewRepository(), scanReviewRow(), scanReviewRows(), NewReviewHandler(), ReviewRepository (+4 more)

### Community 13 - "UseCase"
Cohesion: 0.09
Nodes (7): UseCase, Product, ServiceSupply, StockMovement, NewStockHandler(), StockRepository, StockRepository

### Community 14 - "context.Context"
Cohesion: 0.19
Nodes (7): MaintenanceWorker, MaintenanceConfig, UseCase, NewMaintenanceWorker(), formatAppointmentLocal(), context.Context, NopTxRunner

### Community 15 - "Tenant"
Cohesion: 0.11
Nodes (12): UseCase, Tenant, UpdateSubscriptionInput, defaultModules(), isTenantIdentityConflict(), defaultModules(), scanTenantRow(), scanTenantRows() (+4 more)

### Community 16 - "landing/src/App.tsx"
Cohesion: 0.10
Nodes (22): App(), CtaSection(), FEATURES, FeaturesSection(), Footer(), Header(), NAV_ITEMS, HeroSection() (+14 more)

### Community 17 - "Service"
Cohesion: 0.11
Nodes (8): UseCase, Service, scanServiceRow(), scanServiceRows(), NewServiceHandler(), ServiceRepository, ServiceRepository, serviceScanner

### Community 18 - "go_pkg_time"
Cohesion: 0.08
Nodes (26): go_pkg_database_sql, go_pkg_erpagendamentos_backend_internal_domain_appointment, go_pkg_erpagendamentos_backend_internal_domain_auth, go_pkg_erpagendamentos_backend_internal_domain_calendar, go_pkg_erpagendamentos_backend_internal_domain_cashregister, go_pkg_erpagendamentos_backend_internal_domain_commerce, go_pkg_erpagendamentos_backend_internal_domain_commission, go_pkg_erpagendamentos_backend_internal_domain_configuration (+18 more)

### Community 19 - "auth/entities.go"
Cohesion: 0.12
Nodes (13): ChangePasswordInput, LoginInput, PasswordResetCodeConfirmInput, PasswordResetConfirmInput, PasswordResetRequestInput, ValidationError, hashResetToken(), NewValidationError() (+5 more)

### Community 20 - "RouterConfig"
Cohesion: 0.06
Nodes (74): hasRole(), isAdminRole(), isPrivilegedAdminRole(), requestWithAuthClaims(), requireAdmin(), requireMasterAdmin(), requireOwnerOrAdmin(), requirePrivilegedAdmin() (+66 more)

### Community 21 - "database/sql.DB"
Cohesion: 0.08
Nodes (35): applyDatabaseSQL(), resolveSchemaFile(), initializeRepositories(), postgresRepositories(), NewAppointmentRepository(), NewAuditRepository(), NewAuthRepository(), NewCalendarRepository() (+27 more)

### Community 22 - "resource_test.go"
Cohesion: 0.54
Nodes (7): UseCase, newAppointment(), newResourceUseCase(), TestCreateAutoAssignsAvailableResource(), TestCreateRejectsBusyExplicitResource(), TestCreateRejectsUnknownResource(), TestUpdateAutoReassignsFreedResource()

### Community 23 - "Customer"
Cohesion: 0.06
Nodes (20): UseCase, UseCase, Customer, Resource, NewCustomerRepository(), onlyDigits(), NewResourceRepository(), scanCustomerRow() (+12 more)

### Community 24 - "UseCase"
Cohesion: 0.08
Nodes (9): UseCase, CommissionRecord, CommissionRule, Payout, Summary, NewCommissionHandler(), PayoutInput, CommissionRepository (+1 more)

### Community 25 - "InventoryPage"
Cohesion: 0.21
Nodes (13): fetchServices(), fetchServiceSupplies(), fetchStockMovements(), InventoryPage(), handleAddSupply(), handleDeleteProduct(), handleRemoveSupply(), handleSaveMovement() (+5 more)

### Community 26 - "UseCase"
Cohesion: 0.16
Nodes (11): UseCase, NewUseCase(), Log, NewNotificationRepository(), NewNotificationHandler(), NotificationRepository, Provider, Repository (+3 more)

### Community 27 - "CashRegister"
Cohesion: 0.13
Nodes (7): UseCase, CashRegister, CashRegisterEntry, ShiftSummary, NewCashRegisterHandler(), CashRegisterRepository, CashRegisterRepository

### Community 28 - "MembershipsPage"
Cohesion: 0.19
Nodes (13): fetchGiftCards(), fetchMemberPlans(), fetchMemberships(), GiftCardsPage(), handleDelete(), loadGiftCards(), statusTone(), MembershipsPage() (+5 more)

### Community 29 - "JobRepository"
Cohesion: 0.14
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

### Community 33 - "UseCase"
Cohesion: 0.16
Nodes (6): RecurrenceInput, UseCase, parseRFC3339(), validStatusChange(), NewCalendarHandler(), CalendarHandler

### Community 34 - "memory/configuration_repository.go"
Cohesion: 0.46
Nodes (5): cloneConfiguration(), defaultMemoryConfiguration(), NewConfigurationRepository(), normalizeConfigurationTenantID(), ConfigurationRepository

### Community 35 - "booking/usecase.go"
Cohesion: 0.09
Nodes (22): isSensitiveLogKey(), redactSensitiveText(), sanitizeLogString(), sanitizeLogValue(), truncateLogString(), go_pkg_bytes, go_pkg_crypto_aes, go_pkg_crypto_cipher (+14 more)

### Community 36 - "UseCase"
Cohesion: 0.11
Nodes (10): UseCase, NewUseCase(), TestCalendarUseCase(), Connection, Repository, SyncEvent, NewCalendarRepository(), NewCalendarSyncHandler() (+2 more)

### Community 37 - "rate_limit.go"
Cohesion: 0.16
Nodes (20): buildCompoundIdempotencyKey(), forwardedClientIP(), isSensitiveRateLimitedRoute(), isTrustedProxy(), newMemoryRateLimiter(), normalizedClientID(), persistentRateLimitConfig(), rateLimitRouteKey() (+12 more)

### Community 38 - "devDependencies"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, devDependencies, tailwindcss, @tailwindcss/vite, @types/node (+24 more)

### Community 39 - "ManagePage.tsx"
Cohesion: 0.14
Nodes (24): Professional, PublicAppointment, Service, ManageStep, Button, ButtonProps, ButtonSize, ButtonVariant (+16 more)

### Community 40 - "MemberPlan"
Cohesion: 0.18
Nodes (5): MemberPlan, scanMemberPlans(), database/sql.Rows, MemberPlanRepository, MemberPlanRepository

### Community 41 - "net/http.ResponseWriter"
Cohesion: 0.11
Nodes (16): writeAdminAuthError(), writeConfigurationError(), writeCustomerError(), writeNotificationError(), writeBadRequestError(), writeConflictError(), writeNotFoundError(), writeRawJSON() (+8 more)

### Community 42 - "time.Time"
Cohesion: 0.18
Nodes (15): countOverlappingSlots(), isClosedDay(), localMidnight(), overlapsAnyBusy(), parseClock(), workingHoursFor(), Configuration, normalizePostgresTenantID() (+7 more)

### Community 43 - "Container"
Cohesion: 0.14
Nodes (24): NewUseCase(), NewUseCase(), UseCase, NewReminderWorker(), UseCase, NewUseCase(), NewUseCase(), NewUseCase() (+16 more)

### Community 44 - "BookingPage.tsx"
Cohesion: 0.13
Nodes (19): formatDateLong(), formatDateTime(), formatMoney(), formatTime(), fromDateKey(), nextDateKeys(), SHORT_DAY_LABELS, toDateKey() (+11 more)

### Community 45 - "useToast"
Cohesion: 0.05
Nodes (54): AdminRoutes(), AdminShell(), handleCopyBookingUrl(), AdminShellProps, AdminNavigationItem, navigationGroups, navigationItems, EstablishmentModal() (+46 more)

### Community 47 - "Account"
Cohesion: 0.14
Nodes (6): Account, PasswordResetToken, UpdateProfileInput, cloneAuthAccount(), resetTokenExpired(), AuthRepository

### Community 48 - "UseCase"
Cohesion: 0.18
Nodes (10): appointmentDurationMinutes(), formatAppointmentForTemplate(), generateManagementCode(), UseCase, hashManagementCode(), parseRFC3339(), BookingInput, PortalAppointment (+2 more)

### Community 49 - "AppointmentHandler"
Cohesion: 0.24
Nodes (3): writeAppointmentError(), resolveActorTenantID(), AppointmentHandler

### Community 50 - "bookingApi.ts"
Cohesion: 0.10
Nodes (30): backendOrigin, buildApiUrl(), AvailabilityResponse, BookingResult, cancelPublicAppointment(), createPublicBooking(), fetchAvailability(), fetchPublicAppointment() (+22 more)

### Community 51 - "memoryRepositories"
Cohesion: 0.13
Nodes (22): TestCashRegisterFlow(), TestCheckoutFlow(), TestCommissionUseCase(), TestFinanceUseCase(), TestJobQueue(), TestOutboxPublishAndDispatch(), TestStockUseCase(), memoryRepositories() (+14 more)

### Community 52 - "Membership"
Cohesion: 0.16
Nodes (3): Membership, MembershipRepository, MembershipRepository

### Community 53 - "Package"
Cohesion: 0.17
Nodes (3): Package, PackageRepository, PackageRepository

### Community 54 - "Transaction"
Cohesion: 0.15
Nodes (6): Summary, Transaction, CategoryTotal, CategoryTotal, FinanceRepository, FinanceRepository

### Community 55 - "writeJSON"
Cohesion: 0.17
Nodes (5): writeBookingError(), health(), writeJSON(), BookingHandler, CashRegisterHandler

### Community 56 - "PackageCredit"
Cohesion: 0.18
Nodes (4): PackageCredit, scanPackageCredits(), PackageCreditRepository, PackageCreditRepository

### Community 57 - "GiftCard"
Cohesion: 0.15
Nodes (4): generateGiftCardCode(), GiftCard, GiftCardRepository, GiftCardRepository

### Community 58 - "Professional"
Cohesion: 0.13
Nodes (8): UseCase, Professional, scanProfessionalRow(), scanProfessionalRows(), NewProfessionalHandler(), ProfessionalRepository, ProfessionalRepository, professionalScanner

### Community 59 - "AuthRepository"
Cohesion: 0.10
Nodes (5): DeviceSession, sessionExpired(), updateAuthPasswordAndVersion(), database/sql.Tx, AuthRepository

### Community 60 - "AuditLog"
Cohesion: 0.16
Nodes (11): AuditLog, Repository, UseCase, NewUseCase(), TestAuditLogRecordingAndListing(), NewAuditRepository(), NewAuditHandler(), TestAuditAndSystemHandlers() (+3 more)

### Community 61 - "checkout/usecase.go"
Cohesion: 0.28
Nodes (4): CheckoutInput, CheckoutResult, ProductItem, go_pkg_erpagendamentos_backend_internal_domain_stock

### Community 62 - "EventBroker"
Cohesion: 0.18
Nodes (10): NewBookingHandler(), NewCheckoutHandler(), EventBroker, NewEventBroker(), RegisterMetricsBroker(), go_pkg_sync_atomic, sync.RWMutex, CheckoutHandler (+2 more)

### Community 63 - "auth_google_handler.go"
Cohesion: 0.25
Nodes (10): appendQueryParam(), fetchGoogleUserInfo(), NewAuthGoogleHandler(), resolveFrontendURL(), resolveGoogleRedirectURI(), go_pkg_net_url, AuthGoogleHandler, googleStateData (+2 more)

### Community 64 - "Queue"
Cohesion: 0.29
Nodes (5): Queue, NewQueue(), Repository, sync.WaitGroup, HandlerFunc

### Community 65 - "Appointment"
Cohesion: 0.23
Nodes (4): Appointment, parseAppointmentBound(), parseAppointmentInterval(), AppointmentRepository

### Community 66 - "writeInvalidPayload"
Cohesion: 0.26
Nodes (3): writeCommerceError(), writeInvalidPayload(), CommerceHandler

### Community 67 - "UseCase"
Cohesion: 0.22
Nodes (5): UseCase, NewAppointmentHandler(), NewFinanceHandler(), NewMercadoPagoWebhookHandler(), MercadoPagoWebhookHandler

### Community 68 - "templates.go"
Cohesion: 0.27
Nodes (17): greeting(), PublicBookingOrigin(), RenderCancellation(), RenderConfirmation(), RenderReminder(), RenderReschedule(), sampleData(), TestGreetingWithoutName() (+9 more)

### Community 69 - ".UpdateAccountRole"
Cohesion: 0.39
Nodes (4): canManageTenantAccount(), isPrivilegedAdminRole(), isStaffRole(), normalizeTenantScope()

### Community 70 - "compilerOptions"
Cohesion: 0.09
Nodes (21): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+13 more)

### Community 71 - "AppointmentRepository"
Cohesion: 0.17
Nodes (6): appointmentSelectColumns(), nullIfEmpty(), scanAppointmentRow(), scanAppointmentRows(), AppointmentRepository, appointmentScanner

### Community 72 - "compilerOptions"
Cohesion: 0.09
Nodes (21): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+13 more)

### Community 73 - "compilerOptions"
Cohesion: 0.09
Nodes (21): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+13 more)

### Community 75 - "UseCase"
Cohesion: 0.22
Nodes (7): UseCase, NewUseCase(), NewSystemHandler(), go_pkg_runtime, EventDispatcher, Repository, SystemHandler

### Community 76 - "totp.go"
Cohesion: 0.16
Nodes (10): MFAStatus, generateTOTPSecret(), totpCode(), totpURI(), validateTOTP(), isMFAEligibleRole(), go_pkg_crypto_hmac, go_pkg_crypto_sha1 (+2 more)

### Community 78 - "testing.T"
Cohesion: 0.35
Nodes (11): UseCase, newTestUseCase(), TestAvailableSlotsRespectsClosedDays(), TestCancelPublic(), TestCreateBookingGeneratesManagementCode(), TestManagementCodeNormalization(), TestPortalCancelValidatesOwnership(), TestPortalLookup() (+3 more)

### Community 79 - "lazyAppHandler"
Cohesion: 0.43
Nodes (3): lazyAppHandler, newLazyAppHandler(), writeLazyJSON()

### Community 80 - "ERP Agendamentos"
Cohesion: 0.14
Nodes (13): Banco de Dados & Schema, CI, Endpoints principais, ERP Agendamentos, Estrutura, Pré-requisitos, Rodando com Docker (produção local), Rodando em desenvolvimento (sem banco) (+5 more)

### Community 81 - "TestAuthGoogleHandlerStatusAndLogin"
Cohesion: 0.24
Nodes (9): TestOnboardingSignup(), NewAuthRepository(), NewServiceRepository(), NewTenantRepository(), NewBcryptHasher(), NewJWTManager(), TestAuthGoogleHandlerStatusAndLogin(), go_pkg_golang_org_x_crypto_bcrypt (+1 more)

### Community 82 - "UseCase"
Cohesion: 0.20
Nodes (8): UseCase, NewUseCase(), NewCommerceHandler(), GiftCardRepository, MemberPlanRepository, MembershipRepository, PackageCreditRepository, PackageRepository

### Community 83 - "main.go"
Cohesion: 0.12
Nodes (14): main(), resolveDurationEnv(), resolvePositiveIntEnv(), buildNotificationProvider(), intEnv(), logFallback(), NewEmailProvider(), NewLogProvider() (+6 more)

### Community 84 - "UseCase"
Cohesion: 0.17
Nodes (15): Repository, NewUseCase(), NewUseCase(), UseCase, NewUseCase(), NewUseCase(), NewUseCase(), NewUseCase() (+7 more)

### Community 85 - "booking/src/lib/apiError.ts"
Cohesion: 0.24
Nodes (10): ApiError, ApiErrorPayload, buildApiError(), parseApiResponse(), API_ERROR_MESSAGES, EXACT_MESSAGE_MAP, getApiErrorMessage(), SUBSTRING_MESSAGE_PATTERNS (+2 more)

### Community 89 - "Claims"
Cohesion: 0.33
Nodes (5): isProductionAuthEnvironment(), Claims, decodeSegment(), encodeSegment(), JWTManager

### Community 90 - "buildAuthComponents"
Cohesion: 0.16
Nodes (12): Repository, TokenManager, PasswordHasher, NewUseCase(), buildAuthComponents(), boolEnv(), durationSecondsEnv(), seedAdminAccount() (+4 more)

### Community 91 - "UseCase"
Cohesion: 0.40
Nodes (4): UseCase, NewUseCase(), NewConfigurationHandler(), ConfigurationHandler

### Community 92 - "admin/vite.config.ts"
Cohesion: 0.60
Nodes (3): ref_tailwindcss_vite, ref_vite, ref_vitejs_plugin_react

### Community 93 - "appointment/repository.go"
Cohesion: 0.36
Nodes (5): DayMetric, ProfessionalMetric, ServiceMetric, Analytics, StatusCount

### Community 94 - "UseCase"
Cohesion: 0.19
Nodes (5): RegisterInput, SecretCipher, Session, UseCase, NewAuthHandler()

### Community 95 - "PackagesPage"
Cohesion: 0.47
Nodes (6): fetchPackageCredits(), fetchPackages(), PackagesPage(), handleDelete(), handleToggleActive(), loadPackages()

### Community 96 - "TenantsPage"
Cohesion: 0.24
Nodes (8): fetchTenants(), TenantsPage(), handleCheckReadiness(), handleDelete(), handleToggleActive(), loadAuditLogs(), loadSystemStatus(), loadTenants()

### Community 97 - "connection.go"
Cohesion: 0.60
Nodes (4): Open(), resolveDurationEnv(), resolveIntEnv(), go_pkg_github_com_jackc_pgx_v5_stdlib

### Community 98 - "ical.go"
Cohesion: 0.83
Nodes (3): BuildICalCalendar(), formatUTC(), sanitizeLine()

### Community 99 - "booking/src/App.tsx"
Cohesion: 0.16
Nodes (5): App(), booking_src_index, initWidgetBridge(), BookingPage(), SlugEntryPage()

### Community 109 - "WaitingListPage"
Cohesion: 0.67
Nodes (4): fetchWaitingList(), WaitingListPage(), handleStatusChange(), loadEntries()

### Community 115 - "PortalPage.tsx"
Cohesion: 0.16
Nodes (14): PortalAppointment, portalCancelAppointment(), portalLookup(), portalRescheduleAppointment(), PortalResult, Slot, ACTIONABLE_STATUSES, PortalPage() (+6 more)

### Community 119 - "CommissionsPage.tsx"
Cohesion: 0.12
Nodes (23): CommissionsPage, createCommissionPayout(), fetchCommissionPayouts(), fetchCommissionRecords(), fetchCommissionRules(), fetchCommissionSummaries(), fetchProfessionals(), ActiveTab (+15 more)

### Community 122 - "go_pkg_os"
Cohesion: 0.13
Nodes (10): loadEnvFile(), secureTemporaryPassword(), resolveMaxRequestBodyBytes(), shouldLimitRequestBody(), withRequestBodyLimit(), go_pkg_bufio, go_pkg_erpagendamentos_backend_internal_infrastructure_persistence_memory, go_pkg_net_http_httptest (+2 more)

## Knowledge Gaps
- **263 isolated node(s):** `graphify-mcp`, `name`, `private`, `version`, `type` (+258 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **17 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Container` connect `Container` to `Payment`, `go_pkg_net_http`, `Review`, `UseCase`, `context.Context`, `Tenant`, `Service`, `RouterConfig`, `database/sql.DB`, `Customer`, `UseCase`, `UseCase`, `CashRegister`, `Entry`, `UseCase`, `UseCase`, `UseCase`, `Professional`, `AuditLog`, `EventBroker`, `Queue`, `UseCase`, `UseCase`, `lazyAppHandler`, `UseCase`, `UseCase`, `buildAuthComponents`, `UseCase`, `UseCase`?**
  _High betweenness centrality (0.051) - this node is a cross-community bridge._
- **Why does `UseCase` connect `UseCase` to `booking/usecase.go`, `net/http.Request`, `.UpdateAccountRole`, `rate_limit.go`, `withIdempotency`, `writeAPIError`, `time.Time`, `Container`, `totp.go`, `Account`, `AppointmentHandler`, `auth/entities.go`, `RouterConfig`, `time.Duration`, `Claims`, `buildAuthComponents`, `AuthRepository`, `auth_google_handler.go`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **Why does `UseCase` connect `UseCase` to `writeInvalidPayload`, `MemberPlan`, `Container`, `go_pkg_time`, `Membership`, `Package`, `RouterConfig`, `PackageCredit`, `GiftCard`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **Are the 133 inferred relationships involving `writeJSON()` (e.g. with `health()` and `readiness()`) actually correct?**
  _`writeJSON()` has 133 INFERRED edges - model-reasoned connections that need verification._
- **Are the 104 inferred relationships involving `resolveRequestTenantID()` (e.g. with `.Analytics()` and `.CheckIn()`) actually correct?**
  _`resolveRequestTenantID()` has 104 INFERRED edges - model-reasoned connections that need verification._
- **What connects `graphify-mcp`, `name`, `private` to the rest of the system?**
  _263 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `adminApi.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.04804003336113428 - nodes in this community are weakly interconnected._
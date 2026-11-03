# Graph Report - erpagendamentos  (2026-11-02)

## Corpus Check
- 287 files · ~142,146 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 18 file(s) not represented in the graph (top: (none) 7, .conf 5, .css 3)

## Summary
- 2984 nodes · 9915 edges · 156 communities (91 shown, 65 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 698 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `5b4b314a`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- InventoryPage.tsx
- schema.sql
- ref_react
- Payment
- writeInvalidPayload
- container.go
- withIdempotency
- CustomersPage
- adminAuth.ts
- writeAPIError
- AgendaPage.tsx
- adminApi.ts
- AdminShell.tsx
- UseCase
- AdminRoutes.tsx
- Tenant
- landing/src/App.tsx
- Service
- Session
- RouterConfig
- database/sql.DB
- Resource
- Customer
- UseCase
- net/http.Request
- UseCase
- CashRegister
- MembershipsPage.tsx
- Queue
- Entry
- admin/package.json
- booking/package.json
- UseCase
- CommissionsPage.tsx
- auth/usecase.go
- Connection
- time.Time
- landing/package.json
- booking/src/pages/WaitingListPage.tsx
- UseCase
- net/http.ResponseWriter
- sync.Mutex
- Repository
- PortalPage.tsx
- useToast
- TenantsPage.tsx
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
- repositorySet
- Professional
- AuthRepository
- AuditLog
- EventBroker
- AuthGoogleHandler
- .UpdateAccountRole
- Appointment
- auth_cookie.go
- UseCase
- templates.go
- ReminderWorker
- compilerOptions
- AppointmentRepository
- compilerOptions
- compilerOptions
- OutboxRepository
- UseCase
- UseCase
- writeResourceError
- testing.T
- Container
- ERP Agendamentos
- router.go
- .GetByID
- main.go
- UseCase
- booking/src/lib/apiError.ts
- time.Duration
- ExportHandler
- writeServiceError
- Claims
- NewContainer
- PackagesPage.tsx
- CashRegisterRepository
- UseCase
- observability.go
- email_provider.go
- Repository
- Summary
- BookingPage.tsx
- graphify
- rules/graphify.md
- workflows/graphify.md
- generate-certs.sh
- backup_postgres.sh
- certbot_ssl.sh
- requireRoles
- booking/src/shared/components/Button.tsx
- erpagendamentos/backend
- writeWaitingListError
- PortalPage
- withCSRFProtection
- booking/src/main.tsx

## God Nodes (most connected - your core abstractions)
1. `writeJSON()` - 137 edges
2. `resolveRequestTenantID()` - 109 edges
3. `adminFetch()` - 78 edges
4. `useToast()` - 77 edges
5. `Button` - 76 edges
6. `writeInvalidPayload()` - 73 edges
7. `UseCase` - 72 edges
8. `invalidateAdminGetCache()` - 65 edges
9. `writeInternalError()` - 58 edges
10. `Input` - 56 edges

## Surprising Connections (you probably didn't know these)
- `getPlanBadge()` --calls--> `Badge()`  [EXTRACTED]
  admin/src/pages/TenantsPage.tsx → admin/src/shared/components/Badge.tsx
- `getStatusBadge()` --calls--> `Badge()`  [EXTRACTED]
  admin/src/pages/TenantsPage.tsx → admin/src/shared/components/Badge.tsx
- `getAuditActionBadge()` --calls--> `Badge()`  [EXTRACTED]
  admin/src/pages/TenantsPage.tsx → admin/src/shared/components/Badge.tsx
- `newResourceUseCase()` --calls--> `NewUseCase()`  [INFERRED]
  backend/internal/application/appointment/resource_test.go → backend/internal/application/appointment/usecase.go
- `newTestUseCase()` --calls--> `NewUseCase()`  [INFERRED]
  backend/internal/application/booking/usecase_test.go → backend/internal/application/booking/usecase.go

## Import Cycles
- None detected.

## Communities (156 total, 65 thin omitted)

### Community 0 - "InventoryPage.tsx"
Cohesion: 0.16
Nodes (24): createProduct(), createStockMovement(), deleteProduct(), deleteService(), fetchProducts(), fetchServices(), fetchServiceSupplies(), fetchStockMovements() (+16 more)

### Community 1 - "schema.sql"
Cohesion: 0.05
Nodes (75): catalog.appointments, catalog.audit_logs, catalog.auth_rate_limits, catalog.auth_sessions, catalog.calendar_connections, catalog.cash_register_entries, catalog.cash_registers, catalog.commission_payouts (+67 more)

### Community 2 - "ref_react"
Cohesion: 0.06
Nodes (82): createAdminUser(), deleteAdminUser(), deleteGiftCard(), fetchAdminUsers(), fetchAppointmentAnalytics(), fetchCurrentTenant(), fetchGiftCards(), fetchNotifications() (+74 more)

### Community 3 - "Payment"
Cohesion: 0.05
Nodes (24): UseCase, NewUseCase(), normalizeStatus(), buildPaymentProvider(), IntentResult, Payment, NewMercadoPagoProvider(), NewMockProvider() (+16 more)

### Community 4 - "writeInvalidPayload"
Cohesion: 0.26
Nodes (3): writeCommerceError(), writeInvalidPayload(), CommerceHandler

### Community 5 - "container.go"
Cohesion: 0.10
Nodes (9): deploymentEnvironment(), isProductionEnv(), resolveMaxRequestBodyBytes(), shouldLimitRequestBody(), withRequestBodyLimit(), NewWidgetHandler(), Method, Status (+1 more)

### Community 6 - "withIdempotency"
Cohesion: 0.15
Nodes (10): isHopByHopHeader(), NewIdempotencyStore(), TestIdempotencyBlocksInFlightConcurrentRequests(), TestIdempotencyDoesNotCacheServerFailures(), TestIdempotencyIgnoredOnGetRequests(), TestIdempotencyReplaysCompletedResponse(), withIdempotency(), idempotencyEntry (+2 more)

### Community 7 - "CustomersPage"
Cohesion: 0.27
Nodes (11): deleteCustomer(), downloadAdminCsv(), fetchCustomers(), handleExport(), CustomersPage(), handleDelete(), handleExport(), handleSearch() (+3 more)

### Community 8 - "adminAuth.ts"
Cohesion: 0.06
Nodes (61): App(), handleLogout(), AppProviders(), AppProvidersProps, ResetPasswordPage(), handleSubmit(), AppRouter(), ADMIN_ROLES (+53 more)

### Community 9 - "writeAPIError"
Cohesion: 0.06
Nodes (41): requestWithAuthClaims(), requireMasterAdmin(), resolveAuthClaims(), resolveTenantID(), withProtectedNamespaces(), browserSessionPayload(), buildPasswordResetURL(), exposeAuthTokenForTests() (+33 more)

### Community 10 - "AgendaPage.tsx"
Cohesion: 0.08
Nodes (58): AppointmentFormModal(), computeEndsAt(), handleSubmit(), AppointmentFormModalProps, STATUS_OPTIONS, CheckoutComandaModal(), handleSubmit(), CheckoutComandaModalProps (+50 more)

### Community 11 - "adminApi.ts"
Cohesion: 0.04
Nodes (104): AppointmentFilters, closeCashRegister(), connectCalendar(), createCashRegisterEntry(), createCommissionPayout(), createFinancialTransaction(), createGiftCard(), createMemberPlan() (+96 more)

### Community 12 - "AdminShell.tsx"
Cohesion: 0.11
Nodes (26): AdminShell(), handleCopyBookingUrl(), AdminShellProps, SaaSPlansUpgradeModal(), AdminNavigationItem, navigationGroups, navigationItems, EstablishmentModal() (+18 more)

### Community 13 - "UseCase"
Cohesion: 0.07
Nodes (8): UseCase, Product, Repository, ServiceSupply, StockMovement, NewStockHandler(), StockRepository, StockRepository

### Community 14 - "AdminRoutes.tsx"
Cohesion: 0.14
Nodes (26): AccessDeniedPage(), AdminRouteFallback(), AdminRoutes(), AdminRoutesProps, AgendaPage, AnalyticsPage, AppointmentsPage, CommissionsPage (+18 more)

### Community 15 - "Tenant"
Cohesion: 0.07
Nodes (19): UseCase, UseCase, NewUseCase(), ensureDefaultTenant(), Tenant, UpdateSubscriptionInput, defaultModules(), isTenantIdentityConflict() (+11 more)

### Community 16 - "landing/src/App.tsx"
Cohesion: 0.11
Nodes (21): App(), CtaSection(), FEATURES, FeaturesSection(), Footer(), Header(), NAV_ITEMS, HeroSection() (+13 more)

### Community 17 - "Service"
Cohesion: 0.09
Nodes (10): UseCase, NewUseCase(), Service, scanServiceRow(), scanServiceRows(), NewServiceHandler(), ServiceRepository, ServiceRepository (+2 more)

### Community 19 - "Session"
Cohesion: 0.12
Nodes (14): ChangePasswordInput, LoginInput, PasswordResetCodeConfirmInput, PasswordResetConfirmInput, RegisterInput, Session, ValidationError, totpCode() (+6 more)

### Community 20 - "RouterConfig"
Cohesion: 0.27
Nodes (26): requireAdmin(), NewRouter(), registerAppointmentRoutes(), registerConfigurationRoutes(), registerCustomerRoutes(), registerProfessionalRoutes(), registerResourceRoutes(), registerServiceRoutes() (+18 more)

### Community 21 - "database/sql.DB"
Cohesion: 0.11
Nodes (33): applyDatabaseSQL(), resolveSchemaFile(), initializeRepositories(), postgresRepositories(), NewAppointmentRepository(), NewAuditRepository(), NewAuthRepository(), NewCalendarRepository() (+25 more)

### Community 22 - "Resource"
Cohesion: 0.10
Nodes (11): NewUseCase(), UseCase, NewUseCase(), Resource, scanResourceRow(), NewExportHandler(), NewResourceHandler(), ResourceRepository (+3 more)

### Community 23 - "Customer"
Cohesion: 0.10
Nodes (12): UseCase, NewUseCase(), Customer, NewCustomerRepository(), onlyDigits(), scanCustomerRow(), scanCustomerRows(), NewCustomerHandler() (+4 more)

### Community 24 - "UseCase"
Cohesion: 0.06
Nodes (11): UseCase, NewUseCase(), CommissionRecord, CommissionRule, Payout, Repository, Summary, NewCommissionHandler() (+3 more)

### Community 25 - "net/http.Request"
Cohesion: 0.10
Nodes (10): isAdminRole(), resolveRequestTenantID(), writeProfessionalError(), writeInternalError(), CalendarSyncHandler, CashRegisterHandler, CommissionHandler, FinanceHandler (+2 more)

### Community 26 - "UseCase"
Cohesion: 0.15
Nodes (10): UseCase, NewUseCase(), Log, NewNotificationRepository(), NewNotificationHandler(), NotificationRepository, Provider, Repository (+2 more)

### Community 27 - "CashRegister"
Cohesion: 0.18
Nodes (4): CashRegister, CashRegisterEntry, Repository, CashRegisterRepository

### Community 28 - "MembershipsPage.tsx"
Cohesion: 0.33
Nodes (12): deleteMembership(), fetchMemberPlans(), fetchMemberships(), setMembershipActive(), updateMemberPlanStatus(), MembershipsPage(), handleDelete(), handlePlanStatusChange() (+4 more)

### Community 29 - "Queue"
Cohesion: 0.07
Nodes (10): Queue, NewQueue(), TestJobQueue(), Job, Repository, NewJobRepository(), HandlerFunc, JobRepository (+2 more)

### Community 30 - "Entry"
Cohesion: 0.11
Nodes (11): UseCase, NewUseCase(), Entry, NewWaitingListRepository(), scanWaitingRow(), scanWaitingRows(), NewWaitingListHandler(), WaitingListRepository (+3 more)

### Community 31 - "admin/package.json"
Cohesion: 0.06
Nodes (34): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+26 more)

### Community 32 - "booking/package.json"
Cohesion: 0.06
Nodes (34): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+26 more)

### Community 33 - "UseCase"
Cohesion: 0.12
Nodes (8): DayMetric, ProfessionalMetric, RecurrenceInput, ServiceMetric, UseCase, parseRFC3339(), validStatusChange(), Analytics

### Community 34 - "CommissionsPage.tsx"
Cohesion: 0.19
Nodes (20): createCommissionRule(), deleteCommissionRule(), deleteProfessional(), fetchCommissionPayouts(), fetchCommissionRecords(), fetchCommissionRules(), fetchCommissionSummaries(), fetchProfessionals() (+12 more)

### Community 35 - "auth/usecase.go"
Cohesion: 0.08
Nodes (7): persistentRateLimitConfig(), TestPersistentRateLimitConfigMatch(), authClaimsContextKey, googleStateData, googleTokenResponse, persistentRateLimit, tokenPayload

### Community 36 - "Connection"
Cohesion: 0.09
Nodes (10): UseCase, NewUseCase(), TestCalendarUseCase(), Connection, Repository, SyncEvent, NewCalendarRepository(), NewCalendarSyncHandler() (+2 more)

### Community 37 - "time.Time"
Cohesion: 0.17
Nodes (15): countOverlappingSlots(), formatAppointmentForTemplate(), isClosedDay(), localMidnight(), overlapsAnyBusy(), parseClock(), parseRFC3339(), BuildICalCalendar() (+7 more)

### Community 38 - "landing/package.json"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, devDependencies, tailwindcss, @tailwindcss/vite, @types/node (+24 more)

### Community 39 - "booking/src/pages/WaitingListPage.tsx"
Cohesion: 0.16
Nodes (23): App(), fetchPublicReview(), signUpWaitingList(), submitPublicReview(), DetailsStep(), ReviewPage(), handleSubmit(), SlugEntryPage() (+15 more)

### Community 40 - "UseCase"
Cohesion: 0.10
Nodes (7): UseCase, NewUseCase(), MemberPlan, NewMemberPlanRepository(), NewCommerceHandler(), MemberPlanRepository, MemberPlanRepository

### Community 41 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (18): writeAdminAuthError(), newAuthMFAHandler(), writeConfigurationError(), writeCustomerError(), HandleMetrics(), writeNotificationError(), writeBadRequestError(), writeConflictError() (+10 more)

### Community 42 - "sync.Mutex"
Cohesion: 0.12
Nodes (18): NewUseCase(), workingHoursFor(), UseCase, NewUseCase(), Configuration, cloneConfiguration(), defaultMemoryConfiguration(), normalizeConfigurationTenantID() (+10 more)

### Community 44 - "PortalPage.tsx"
Cohesion: 0.18
Nodes (21): PortalAppointment, PortalResult, Slot, formatDateLong(), formatDateTime(), formatTime(), fromDateKey(), nextDateKeys() (+13 more)

### Community 45 - "useToast"
Cohesion: 0.09
Nodes (55): createResource(), fetchConfiguration(), updateConfiguration(), updateResource(), AdminAccount, AdminAuthError, DAY_LABELS, CommissionPayoutModal() (+47 more)

### Community 46 - "TenantsPage.tsx"
Cohesion: 0.15
Nodes (22): checkSystemReadiness(), deleteTenant(), fetchAuditLogs(), fetchSystemStatus(), fetchTenants(), setTenantActive(), parseJson(), AVAILABLE_MODULES (+14 more)

### Community 47 - "Account"
Cohesion: 0.14
Nodes (6): Account, PasswordResetToken, UpdateProfileInput, cloneAuthAccount(), resetTokenExpired(), AuthRepository

### Community 48 - "UseCase"
Cohesion: 0.19
Nodes (8): appointmentDurationMinutes(), generateManagementCode(), UseCase, hashManagementCode(), BookingInput, PortalAppointment, PortalCustomer, PortalResult

### Community 49 - "AppointmentHandler"
Cohesion: 0.26
Nodes (3): writeAppointmentError(), resolveActorTenantID(), AppointmentHandler

### Community 50 - "bookingApi.ts"
Cohesion: 0.14
Nodes (25): backendOrigin, buildApiUrl(), AvailabilityResponse, BookingResult, cancelPublicAppointment(), fetchAvailability(), fetchPublicAppointment(), fetchPublicInfo() (+17 more)

### Community 51 - "memoryRepositories"
Cohesion: 0.17
Nodes (20): TestCashRegisterFlow(), TestCheckoutFlow(), TestCommissionUseCase(), NewUseCase(), TestOnboardingSignup(), memoryRepositories(), NewAppointmentRepository(), NewAuthRepository() (+12 more)

### Community 52 - "Membership"
Cohesion: 0.12
Nodes (5): Membership, NewMembershipRepository(), MembershipRepository, MembershipRepository, MembershipRepository

### Community 53 - "Package"
Cohesion: 0.12
Nodes (5): Package, NewPackageRepository(), PackageRepository, PackageRepository, PackageRepository

### Community 54 - "Transaction"
Cohesion: 0.24
Nodes (3): Transaction, FinanceRepository, FinanceRepository

### Community 55 - "writeJSON"
Cohesion: 0.21
Nodes (6): writeBookingError(), health(), writePaymentError(), writeJSON(), BookingHandler, PaymentHandler

### Community 56 - "PackageCredit"
Cohesion: 0.12
Nodes (6): PackageCredit, NewPackageCreditRepository(), scanPackageCredits(), PackageCreditRepository, PackageCreditRepository, PackageCreditRepository

### Community 57 - "repositorySet"
Cohesion: 0.10
Nodes (7): generateGiftCardCode(), GiftCard, NewGiftCardRepository(), repositorySet, GiftCardRepository, GiftCardRepository, GiftCardRepository

### Community 58 - "Professional"
Cohesion: 0.05
Nodes (22): UseCase, NewUseCase(), UseCase, NewUseCase(), Professional, Review, NewReviewRepository(), scanProfessionalRow() (+14 more)

### Community 60 - "AuditLog"
Cohesion: 0.15
Nodes (13): AuditLog, Repository, UseCase, NewUseCase(), TestAuditLogRecordingAndListing(), TestOutboxPublishAndDispatch(), NewAuditRepository(), NewOutboxRepository() (+5 more)

### Community 62 - "EventBroker"
Cohesion: 0.23
Nodes (8): NewBookingHandler(), NewCheckoutHandler(), EventBroker, NewMercadoPagoWebhookHandler(), RegisterMetricsBroker(), CheckoutHandler, Event, MercadoPagoWebhookHandler

### Community 63 - "AuthGoogleHandler"
Cohesion: 0.27
Nodes (6): appendQueryParam(), fetchGoogleUserInfo(), resolveFrontendURL(), resolveGoogleRedirectURI(), AuthGoogleHandler, googleUserInfo

### Community 64 - ".UpdateAccountRole"
Cohesion: 0.39
Nodes (4): canManageTenantAccount(), isPrivilegedAdminRole(), isStaffRole(), normalizeTenantScope()

### Community 65 - "Appointment"
Cohesion: 0.10
Nodes (8): Repository, Appointment, StatusCount, parseAppointmentBound(), parseAppointmentInterval(), parseTimeFlexible(), BookingResult, AppointmentRepository

### Community 66 - "auth_cookie.go"
Cohesion: 0.41
Nodes (13): authCookieDomain(), authCookiePartitioned(), authCookieSameSite(), authCookieSecure(), authTokenFromCookie(), clearAuthCookie(), csrfTokenFromCookie(), isAuthCookieProductionEnv() (+5 more)

### Community 67 - "UseCase"
Cohesion: 0.25
Nodes (4): UseCase, NewUseCase(), ShiftSummary, NewCashRegisterHandler()

### Community 68 - "templates.go"
Cohesion: 0.49
Nodes (9): greeting(), PublicBookingOrigin(), RenderCancellation(), RenderConfirmation(), RenderReminder(), RenderReschedule(), writeFooter(), writeLine() (+1 more)

### Community 69 - "ReminderWorker"
Cohesion: 0.14
Nodes (9): MaintenanceWorker, MaintenanceConfig, NewMaintenanceWorker(), formatAppointmentLocal(), NewReminderWorker(), startBackgroundWorkers(), ReminderWorker, NopTxRunner (+1 more)

### Community 70 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 71 - "AppointmentRepository"
Cohesion: 0.16
Nodes (6): appointmentSelectColumns(), nullIfEmpty(), scanAppointmentRow(), scanAppointmentRows(), AppointmentRepository, appointmentScanner

### Community 72 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 73 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 74 - "OutboxRepository"
Cohesion: 0.18
Nodes (4): quoteLiteral(), OutboxRepository, OutboxEvent, OutboxRepository

### Community 75 - "UseCase"
Cohesion: 0.17
Nodes (6): UseCase, NewUseCase(), NewSystemHandler(), EventDispatcher, Repository, SystemHandler

### Community 78 - "testing.T"
Cohesion: 0.14
Nodes (23): newAppointment(), newResourceUseCase(), TestCreateAutoAssignsAvailableResource(), TestCreateRejectsBusyExplicitResource(), TestCreateRejectsUnknownResource(), TestUpdateAutoReassignsFreedResource(), newTestUseCase(), TestAvailableSlotsRespectsClosedDays() (+15 more)

### Community 79 - "Container"
Cohesion: 0.12
Nodes (11): lazyAppHandler, writeLazyJSON(), boolEnv(), durationSecondsEnv(), Container, isPrivilegedAdminRole(), isSensitiveAdminMutation(), mfaEnforcementEnabled() (+3 more)

### Community 80 - "ERP Agendamentos"
Cohesion: 0.14
Nodes (13): Banco de Dados & Schema, CI, Endpoints principais, ERP Agendamentos, Estrutura, Pré-requisitos, Rodando com Docker (produção local), Rodando em desenvolvimento (sem banco) (+5 more)

### Community 81 - "router.go"
Cohesion: 0.14
Nodes (4): CheckoutInput, ProductItem, CategoryTotal, Type

### Community 83 - "main.go"
Cohesion: 0.38
Nodes (4): newLazyAppHandler(), main(), resolveDurationEnv(), resolvePositiveIntEnv()

### Community 84 - "UseCase"
Cohesion: 0.27
Nodes (5): UseCase, AuditLogger, CheckoutResult, OutboxPublisher, TxRunner

### Community 85 - "booking/src/lib/apiError.ts"
Cohesion: 0.26
Nodes (10): ApiError, ApiErrorPayload, buildApiError(), parseApiResponse(), API_ERROR_MESSAGES, EXACT_MESSAGE_MAP, getApiErrorMessage(), SUBSTRING_MESSAGE_PATTERNS (+2 more)

### Community 86 - "time.Duration"
Cohesion: 0.12
Nodes (3): CleanupSummary, DeviceSession, sessionExpired()

### Community 89 - "Claims"
Cohesion: 0.31
Nodes (5): isProductionAuthEnvironment(), Claims, decodeSegment(), encodeSegment(), JWTManager

### Community 90 - "NewContainer"
Cohesion: 0.12
Nodes (13): NewUseCase(), TestFinanceUseCase(), NewUseCase(), TestStockUseCase(), buildAuthComponents(), NewContainer(), seedAdminAccount(), loadEnvFile() (+5 more)

### Community 91 - "PackagesPage.tsx"
Cohesion: 0.24
Nodes (16): deletePackage(), deleteResource(), fetchPackageCredits(), fetchPackages(), fetchResources(), setPackageActive(), setResourceActive(), PackagesPage() (+8 more)

### Community 94 - "UseCase"
Cohesion: 0.09
Nodes (16): MFAStatus, PasswordResetRequestInput, SecretCipher, TokenManager, generateTOTPSecret(), totpURI(), PasswordHasher, UseCase (+8 more)

### Community 95 - "observability.go"
Cohesion: 0.52
Nodes (6): isSensitiveLogKey(), redactSensitiveText(), sanitizeLogString(), sanitizeLogValue(), truncateLogString(), contextKey

### Community 96 - "email_provider.go"
Cohesion: 0.21
Nodes (7): buildNotificationProvider(), intEnv(), logFallback(), NewEmailProvider(), NewLogProvider(), EmailProvider, LogProvider

### Community 99 - "BookingPage.tsx"
Cohesion: 0.15
Nodes (14): createPublicBooking(), formatMoney(), BookingPage(), handleBookingSubmit(), CenteredLayout(), DateTimeStep(), EmptyNote(), ErrorCard() (+6 more)

### Community 109 - "requireRoles"
Cohesion: 0.53
Nodes (5): hasRole(), requireOwnerOrAdmin(), requirePrivilegedAdmin(), requireRoles(), registerAuthRoutes()

### Community 110 - "booking/src/shared/components/Button.tsx"
Cohesion: 0.33
Nodes (5): ButtonProps, ButtonSize, ButtonVariant, sizeClasses, variantClasses

### Community 115 - "PortalPage"
Cohesion: 0.39
Nodes (8): portalCancelAppointment(), portalLookup(), portalRescheduleAppointment(), PortalPage(), handleCancel(), handleConfirmReschedule(), handleLookup(), refresh()

### Community 117 - "withCSRFProtection"
Cohesion: 0.32
Nodes (6): isLocalDevelopmentOrigin(), isOriginAllowed(), resolveAllowedOrigins(), withCORS(), requiresCSRFProtection(), withCSRFProtection()

## Knowledge Gaps
- **286 isolated node(s):** `graphify-mcp`, `name`, `private`, `version`, `type` (+281 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 546 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **65 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Container` connect `Container` to `Payment`, `container.go`, `UseCase`, `Tenant`, `Service`, `database/sql.DB`, `Resource`, `Customer`, `UseCase`, `UseCase`, `Queue`, `Entry`, `UseCase`, `Connection`, `UseCase`, `sync.Mutex`, `UseCase`, `Professional`, `AuditLog`, `EventBroker`, `UseCase`, `ReminderWorker`, `UseCase`, `UseCase`, `UseCase`, `NewContainer`, `UseCase`?**
  _High betweenness centrality (0.038) - this node is a cross-community bridge._
- **Why does `RouterConfig` connect `RouterConfig` to `Payment`, `UseCase`, `Tenant`, `Service`, `Resource`, `Customer`, `UseCase`, `UseCase`, `Entry`, `UseCase`, `Connection`, `UseCase`, `sync.Mutex`, `UseCase`, `Professional`, `AuditLog`, `EventBroker`, `UseCase`, `ReminderWorker`, `UseCase`, `UseCase`, `router.go`, `UseCase`, `time.Duration`, `UseCase`, `requireRoles`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **Why does `repositorySet` connect `repositorySet` to `Payment`, `UseCase`, `Tenant`, `Service`, `go_pkg_time`, `database/sql.DB`, `Resource`, `Customer`, `UseCase`, `UseCase`, `CashRegister`, `Queue`, `Entry`, `Connection`, `UseCase`, `sync.Mutex`, `Repository`, `memoryRepositories`, `Membership`, `Package`, `PackageCredit`, `Professional`, `AuditLog`, `Appointment`, `UseCase`, `UseCase`, `Repository`?**
  _High betweenness centrality (0.022) - this node is a cross-community bridge._
- **Are the 133 inferred relationships involving `writeJSON()` (e.g. with `health()` and `readiness()`) actually correct?**
  _`writeJSON()` has 133 INFERRED edges - model-reasoned connections that need verification._
- **Are the 104 inferred relationships involving `resolveRequestTenantID()` (e.g. with `.Analytics()` and `.CheckIn()`) actually correct?**
  _`resolveRequestTenantID()` has 104 INFERRED edges - model-reasoned connections that need verification._
- **What connects `graphify-mcp`, `name`, `private` to the rest of the system?**
  _286 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `schema.sql` be split into smaller, more focused modules?**
  _Cohesion score 0.05297334244702666 - nodes in this community are weakly interconnected._
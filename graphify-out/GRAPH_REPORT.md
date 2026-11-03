# Graph Report - erpagendamentos  (2026-11-02)

## Corpus Check
- 287 files · ~141,341 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 18 file(s) not represented in the graph (top: (none) 7, .conf 5, .css 3)

## Summary
- 2981 nodes · 9891 edges · 159 communities (97 shown, 62 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 698 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `d74f6c1e`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- adminFetch
- schema.sql
- ref_react
- Payment
- resolveRequestTenantID
- container.go
- withIdempotency
- PaymentsPage.tsx
- adminAuth.ts
- writeAPIError
- AppointmentsPage.tsx
- adminApi.ts
- Review
- UseCase
- AdminRoutes.tsx
- Tenant
- landing/src/App.tsx
- Service
- go_pkg_time
- auth/entities.go
- RouterConfig
- database/sql.DB
- Resource
- Customer
- UseCase
- writeJSON
- UseCase
- CashRegister
- MembershipsPage.tsx
- Job
- Entry
- admin/package.json
- booking/package.json
- Appointment
- FinancePage.tsx
- auth/usecase.go
- Connection
- rate_limit.go
- landing/package.json
- booking/src/pages/WaitingListPage.tsx
- MemberPlan
- writeNotFoundError
- time.Time
- UseCase
- PortalPage.tsx
- useToast
- TenantsPage.tsx
- AuthRepository
- UseCase
- AppointmentHandler
- bookingApi.ts
- NewContainer
- UseCase
- Package
- Transaction
- net/http.Request
- PackageCredit
- GiftCard
- Professional
- Account
- AuditLog
- EventBroker
- AuthGoogleHandler
- ServiceRepository
- AppointmentRepository
- auth_cookie.go
- Repository
- templates.go
- context.Context
- compilerOptions
- AppointmentRepository
- compilerOptions
- compilerOptions
- sync.RWMutex
- UseCase
- CustomerRepository
- net/http.ResponseWriter
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
- Repository
- Claims
- SecretCipher
- PackagesPage
- Analytics
- UseCase
- observability.go
- email_provider.go
- Queue
- BuildICalCalendar
- BookingPage.tsx
- graphify
- rules/graphify.md
- workflows/graphify.md
- generate-certs.sh
- backup_postgres.sh
- certbot_ssl.sh
- auth_context.go
- booking/src/shared/components/Button.tsx
- erpagendamentos/backend
- idempotencyRecorder
- PortalPage
- ReminderWorker
- withCSRFProtection
- booking/src/main.tsx
- OutboxRepository
- JobRepository

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

## Communities (159 total, 62 thin omitted)

### Community 0 - "adminFetch"
Cohesion: 0.06
Nodes (61): connectCalendar(), createAdminUser(), createGiftCard(), createMembership(), createPackage(), createPackageCredit(), createPayment(), createProfessional() (+53 more)

### Community 1 - "schema.sql"
Cohesion: 0.05
Nodes (75): catalog.appointments, catalog.audit_logs, catalog.auth_rate_limits, catalog.auth_sessions, catalog.calendar_connections, catalog.cash_register_entries, catalog.cash_registers, catalog.commission_payouts (+67 more)

### Community 2 - "ref_react"
Cohesion: 0.08
Nodes (66): deleteGiftCard(), fetchAppointmentAnalytics(), fetchGiftCards(), fetchNotifications(), fetchReviews(), fetchWaitingList(), updateWaitingListStatus(), AnalyticsPage() (+58 more)

### Community 3 - "Payment"
Cohesion: 0.05
Nodes (23): UseCase, NewUseCase(), normalizeStatus(), buildPaymentProvider(), IntentResult, Payment, NewMercadoPagoProvider(), NewMockProvider() (+15 more)

### Community 4 - "resolveRequestTenantID"
Cohesion: 0.12
Nodes (10): resolveRequestTenantID(), writeCommerceError(), writeConfigurationError(), writeProfessionalError(), writeBadRequestError(), writeInvalidPayload(), writeServiceError(), CommerceHandler (+2 more)

### Community 5 - "container.go"
Cohesion: 0.08
Nodes (10): resolveMaxRequestBodyBytes(), shouldLimitRequestBody(), withRequestBodyLimit(), NewWidgetHandler(), Method, Status, googleStateData, googleTokenResponse (+2 more)

### Community 6 - "withIdempotency"
Cohesion: 0.20
Nodes (9): isHopByHopHeader(), NewIdempotencyStore(), TestIdempotencyBlocksInFlightConcurrentRequests(), TestIdempotencyDoesNotCacheServerFailures(), TestIdempotencyIgnoredOnGetRequests(), TestIdempotencyReplaysCompletedResponse(), withIdempotency(), idempotencyEntry (+1 more)

### Community 7 - "PaymentsPage.tsx"
Cohesion: 0.16
Nodes (21): deleteCustomer(), downloadAdminCsv(), fetchCustomers(), fetchPayments(), markPaymentPaid(), refundPayment(), handleExport(), CustomersPage() (+13 more)

### Community 8 - "adminAuth.ts"
Cohesion: 0.06
Nodes (69): App(), handleLogout(), AppProviders(), AppProvidersProps, ResetPasswordPage(), handleSubmit(), AppRouter(), ADMIN_ROLES (+61 more)

### Community 9 - "writeAPIError"
Cohesion: 0.07
Nodes (26): resolveAuthClaims(), resolveTenantID(), buildPasswordResetURL(), exposeAuthTokenForTests(), AuthHandler, secureTemporaryPassword(), newAuthMFAHandler(), AuthHandler (+18 more)

### Community 10 - "AppointmentsPage.tsx"
Cohesion: 0.12
Nodes (43): AppointmentFormModal(), computeEndsAt(), handleSubmit(), AppointmentFormModalProps, STATUS_OPTIONS, checkInAppointment(), createAppointment(), createCustomer() (+35 more)

### Community 11 - "adminApi.ts"
Cohesion: 0.05
Nodes (84): handleSubmit(), CheckoutComandaModalProps, SelectedProduct, AppointmentFilters, checkoutAppointment(), createCommissionPayout(), createCommissionRule(), createProduct() (+76 more)

### Community 12 - "Review"
Cohesion: 0.12
Nodes (11): UseCase, NewUseCase(), Review, NewReviewRepository(), scanReviewRow(), scanReviewRows(), NewReviewHandler(), ReviewRepository (+3 more)

### Community 13 - "UseCase"
Cohesion: 0.07
Nodes (8): UseCase, Product, Repository, ServiceSupply, StockMovement, NewStockHandler(), StockRepository, StockRepository

### Community 14 - "AdminRoutes.tsx"
Cohesion: 0.09
Nodes (33): AccessDeniedPage(), AdminRouteFallback(), AdminRoutes(), AdminRoutesProps, AgendaPage, AnalyticsPage, AppointmentsPage, CommissionsPage (+25 more)

### Community 15 - "Tenant"
Cohesion: 0.08
Nodes (15): UseCase, NewUseCase(), ensureDefaultTenant(), Tenant, UpdateSubscriptionInput, defaultModules(), isTenantIdentityConflict(), defaultModules() (+7 more)

### Community 16 - "landing/src/App.tsx"
Cohesion: 0.11
Nodes (21): App(), CtaSection(), FEATURES, FeaturesSection(), Footer(), Header(), NAV_ITEMS, HeroSection() (+13 more)

### Community 17 - "Service"
Cohesion: 0.12
Nodes (6): UseCase, NewUseCase(), Service, NewServiceHandler(), ServiceRepository, Repository

### Community 18 - "go_pkg_time"
Cohesion: 0.16
Nodes (3): Open(), resolveDurationEnv(), resolveIntEnv()

### Community 19 - "auth/entities.go"
Cohesion: 0.14
Nodes (11): ChangePasswordInput, LoginInput, PasswordResetConfirmInput, RegisterInput, Session, ValidationError, totpCode(), validateTOTP() (+3 more)

### Community 20 - "RouterConfig"
Cohesion: 0.26
Nodes (27): requireAdmin(), NewRouter(), startBackgroundWorkers(), registerAppointmentRoutes(), registerConfigurationRoutes(), registerCustomerRoutes(), registerProfessionalRoutes(), registerResourceRoutes() (+19 more)

### Community 21 - "database/sql.DB"
Cohesion: 0.13
Nodes (28): postgresRepositories(), NewAppointmentRepository(), NewAuditRepository(), NewAuthRepository(), NewCalendarRepository(), NewCashRegisterRepository(), NewGiftCardRepository(), NewMemberPlanRepository() (+20 more)

### Community 22 - "Resource"
Cohesion: 0.10
Nodes (10): UseCase, NewUseCase(), Resource, scanResourceRow(), NewExportHandler(), NewResourceHandler(), ResourceRepository, ResourceRepository (+2 more)

### Community 23 - "Customer"
Cohesion: 0.13
Nodes (8): UseCase, NewUseCase(), Customer, NewCustomerRepository(), onlyDigits(), NewCustomerHandler(), Repository, CustomerRepository

### Community 24 - "UseCase"
Cohesion: 0.06
Nodes (10): UseCase, CommissionRecord, CommissionRule, Payout, Repository, Summary, NewCommissionHandler(), PayoutInput (+2 more)

### Community 25 - "writeJSON"
Cohesion: 0.11
Nodes (11): health(), writeConflictError(), writeInternalError(), writeJSON(), writeReviewError(), writeWaitingListError(), CalendarSyncHandler, CommissionHandler (+3 more)

### Community 26 - "UseCase"
Cohesion: 0.15
Nodes (10): UseCase, NewUseCase(), Log, NewNotificationRepository(), NewNotificationHandler(), NotificationRepository, Provider, Repository (+2 more)

### Community 27 - "CashRegister"
Cohesion: 0.11
Nodes (8): UseCase, CashRegister, CashRegisterEntry, Repository, ShiftSummary, NewCashRegisterHandler(), CashRegisterRepository, CashRegisterRepository

### Community 28 - "MembershipsPage.tsx"
Cohesion: 0.25
Nodes (15): createMemberPlan(), deleteMembership(), fetchMemberPlans(), fetchMemberships(), setMembershipActive(), updateMemberPlanStatus(), handleSubmit(), MembershipsPage() (+7 more)

### Community 30 - "Entry"
Cohesion: 0.11
Nodes (11): UseCase, NewUseCase(), Entry, NewWaitingListRepository(), scanWaitingRow(), scanWaitingRows(), NewWaitingListHandler(), WaitingListRepository (+3 more)

### Community 31 - "admin/package.json"
Cohesion: 0.06
Nodes (34): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+26 more)

### Community 32 - "booking/package.json"
Cohesion: 0.06
Nodes (34): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+26 more)

### Community 33 - "Appointment"
Cohesion: 0.11
Nodes (8): RecurrenceInput, Repository, UseCase, parseRFC3339(), validStatusChange(), Appointment, NewCalendarHandler(), CalendarHandler

### Community 34 - "FinancePage.tsx"
Cohesion: 0.15
Nodes (27): closeCashRegister(), createCashRegisterEntry(), createFinancialTransaction(), deleteFinancialTransaction(), fetchCashRegisterHistory(), fetchCurrentCashRegister(), fetchFinancialSummary(), fetchFinancialTransactions() (+19 more)

### Community 35 - "auth/usecase.go"
Cohesion: 0.11
Nodes (12): TokenManager, canManageTenantAccount(), PasswordHasher, isPrivilegedAdminRole(), isStaffRole(), NewUseCase(), normalizeTenantScope(), seedAdminAccount() (+4 more)

### Community 36 - "Connection"
Cohesion: 0.09
Nodes (10): UseCase, NewUseCase(), TestCalendarUseCase(), Connection, Repository, SyncEvent, NewCalendarRepository(), NewCalendarSyncHandler() (+2 more)

### Community 37 - "rate_limit.go"
Cohesion: 0.19
Nodes (17): buildCompoundIdempotencyKey(), forwardedClientIP(), isSensitiveRateLimitedRoute(), isTrustedProxy(), newMemoryRateLimiter(), normalizedClientID(), persistentRateLimitConfig(), rateLimitRouteKey() (+9 more)

### Community 38 - "landing/package.json"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, devDependencies, tailwindcss, @tailwindcss/vite, @types/node (+24 more)

### Community 39 - "booking/src/pages/WaitingListPage.tsx"
Cohesion: 0.16
Nodes (23): App(), fetchPublicReview(), signUpWaitingList(), submitPublicReview(), DetailsStep(), ReviewPage(), handleSubmit(), SlugEntryPage() (+15 more)

### Community 40 - "MemberPlan"
Cohesion: 0.13
Nodes (6): MemberPlan, NewMemberPlanRepository(), scanMemberPlans(), MemberPlanRepository, MemberPlanRepository, MemberPlanRepository

### Community 41 - "writeNotFoundError"
Cohesion: 0.16
Nodes (6): writeCustomerError(), writeNotificationError(), writeNotFoundError(), writeTenantError(), CustomerHandler, TenantHandler

### Community 42 - "time.Time"
Cohesion: 0.16
Nodes (17): countOverlappingSlots(), isClosedDay(), localMidnight(), overlapsAnyBusy(), parseClock(), workingHoursFor(), Configuration, cloneConfiguration() (+9 more)

### Community 43 - "UseCase"
Cohesion: 0.36
Nodes (5): UseCase, NewAuthGoogleHandler(), NewOnboardingHandler(), SignupInput, OnboardingHandler

### Community 44 - "PortalPage.tsx"
Cohesion: 0.18
Nodes (21): PortalAppointment, PortalResult, Slot, formatDateLong(), formatDateTime(), formatTime(), fromDateKey(), nextDateKeys() (+13 more)

### Community 45 - "useToast"
Cohesion: 0.08
Nodes (62): AdminShell(), handleCopyBookingUrl(), AdminShellProps, SaaSPlansUpgradeModal(), EstablishmentModal(), handleCopyLink(), handleNameChange(), handleSlugChange() (+54 more)

### Community 46 - "TenantsPage.tsx"
Cohesion: 0.13
Nodes (25): checkSystemReadiness(), createTenant(), deleteTenant(), fetchAuditLogs(), fetchSystemStatus(), fetchTenants(), setTenantActive(), updateTenant() (+17 more)

### Community 47 - "AuthRepository"
Cohesion: 0.09
Nodes (7): DeviceSession, PasswordResetToken, cloneAuthAccount(), resetTokenExpired(), sessionExpired(), AuthRepository, rateLimitEntry

### Community 48 - "UseCase"
Cohesion: 0.17
Nodes (11): appointmentDurationMinutes(), formatAppointmentForTemplate(), generateManagementCode(), UseCase, hashManagementCode(), parseRFC3339(), BookingInput, BookingResult (+3 more)

### Community 49 - "AppointmentHandler"
Cohesion: 0.21
Nodes (4): writeAppointmentError(), resolveActorTenantID(), writeAdminAuthError(), AppointmentHandler

### Community 50 - "bookingApi.ts"
Cohesion: 0.14
Nodes (25): backendOrigin, buildApiUrl(), AvailabilityResponse, BookingResult, cancelPublicAppointment(), fetchAvailability(), fetchPublicAppointment(), fetchPublicInfo() (+17 more)

### Community 51 - "NewContainer"
Cohesion: 0.10
Nodes (36): NewUseCase(), NewUseCase(), TestCashRegisterFlow(), NewUseCase(), TestCheckoutFlow(), NewUseCase(), TestCommissionUseCase(), NewUseCase() (+28 more)

### Community 52 - "UseCase"
Cohesion: 0.10
Nodes (6): UseCase, Membership, NewMembershipRepository(), NewCommerceHandler(), MembershipRepository, MembershipRepository

### Community 53 - "Package"
Cohesion: 0.12
Nodes (5): Package, NewPackageRepository(), PackageRepository, PackageRepository, PackageRepository

### Community 54 - "Transaction"
Cohesion: 0.12
Nodes (5): Repository, Summary, Transaction, FinanceRepository, FinanceRepository

### Community 55 - "net/http.Request"
Cohesion: 0.19
Nodes (4): writeBookingError(), resolveWidgetBookingOrigin(), BookingHandler, CashRegisterHandler

### Community 56 - "PackageCredit"
Cohesion: 0.13
Nodes (7): NewUseCase(), PackageCredit, NewPackageCreditRepository(), scanPackageCredits(), PackageCreditRepository, PackageCreditRepository, PackageCreditRepository

### Community 57 - "GiftCard"
Cohesion: 0.11
Nodes (6): generateGiftCardCode(), GiftCard, NewGiftCardRepository(), GiftCardRepository, GiftCardRepository, GiftCardRepository

### Community 58 - "Professional"
Cohesion: 0.10
Nodes (10): UseCase, NewUseCase(), Professional, scanProfessionalRow(), scanProfessionalRows(), NewProfessionalHandler(), ProfessionalRepository, ProfessionalRepository (+2 more)

### Community 59 - "Account"
Cohesion: 0.07
Nodes (6): Repository, Account, AuthState, UpdateProfileInput, updateAuthPasswordAndVersion(), AuthRepository

### Community 60 - "AuditLog"
Cohesion: 0.17
Nodes (11): AuditLog, Repository, UseCase, NewUseCase(), TestAuditLogRecordingAndListing(), NewAuditRepository(), NewAuditHandler(), TestAuditAndSystemHandlers() (+3 more)

### Community 62 - "EventBroker"
Cohesion: 0.13
Nodes (10): UseCase, NewAppointmentHandler(), NewBookingHandler(), EventBroker, NewEventBroker(), NewFinanceHandler(), NewMercadoPagoWebhookHandler(), RegisterMetricsBroker() (+2 more)

### Community 63 - "AuthGoogleHandler"
Cohesion: 0.24
Nodes (7): appendQueryParam(), fetchGoogleUserInfo(), resolveFrontendURL(), resolveGoogleRedirectURI(), browserSessionPayload(), AuthGoogleHandler, googleUserInfo

### Community 64 - "ServiceRepository"
Cohesion: 0.25
Nodes (4): scanServiceRow(), scanServiceRows(), ServiceRepository, serviceScanner

### Community 65 - "AppointmentRepository"
Cohesion: 0.14
Nodes (5): StatusCount, parseAppointmentBound(), parseAppointmentInterval(), parseTimeFlexible(), AppointmentRepository

### Community 66 - "auth_cookie.go"
Cohesion: 0.41
Nodes (13): authCookieDomain(), authCookiePartitioned(), authCookieSameSite(), authCookieSecure(), authTokenFromCookie(), clearAuthCookie(), csrfTokenFromCookie(), isAuthCookieProductionEnv() (+5 more)

### Community 67 - "Repository"
Cohesion: 0.24
Nodes (6): NewUseCase(), UseCase, NewUseCase(), NewConfigurationHandler(), Repository, ConfigurationHandler

### Community 68 - "templates.go"
Cohesion: 0.27
Nodes (16): greeting(), PublicBookingOrigin(), RenderCancellation(), RenderConfirmation(), RenderReminder(), RenderReschedule(), sampleData(), TestGreetingWithoutName() (+8 more)

### Community 69 - "context.Context"
Cohesion: 0.22
Nodes (5): MaintenanceWorker, MaintenanceConfig, NewMaintenanceWorker(), NopTxRunner, TxRunner

### Community 70 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 71 - "AppointmentRepository"
Cohesion: 0.21
Nodes (5): appointmentSelectColumns(), scanAppointmentRow(), scanAppointmentRows(), AppointmentRepository, appointmentScanner

### Community 72 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 73 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 74 - "sync.RWMutex"
Cohesion: 0.24
Nodes (3): OutboxRepository, OutboxEvent, MetricsRegistry

### Community 75 - "UseCase"
Cohesion: 0.16
Nodes (8): UseCase, NewUseCase(), TestOutboxPublishAndDispatch(), NewOutboxRepository(), NewSystemHandler(), EventDispatcher, Repository, SystemHandler

### Community 76 - "CustomerRepository"
Cohesion: 0.29
Nodes (4): scanCustomerRow(), scanCustomerRows(), CustomerRepository, customerScanner

### Community 77 - "net/http.ResponseWriter"
Cohesion: 0.14
Nodes (7): HandleMetrics(), writePaymentError(), writeResourceError(), PaymentHandler, ResourceHandler, responseRecorder, StockHandler

### Community 78 - "testing.T"
Cohesion: 0.15
Nodes (18): newAppointment(), newResourceUseCase(), TestCreateAutoAssignsAvailableResource(), TestCreateRejectsBusyExplicitResource(), TestCreateRejectsUnknownResource(), TestUpdateAutoReassignsFreedResource(), newTestUseCase(), TestAvailableSlotsRespectsClosedDays() (+10 more)

### Community 79 - "Container"
Cohesion: 0.13
Nodes (10): lazyAppHandler, writeLazyJSON(), boolEnv(), Container, isPrivilegedAdminRole(), isSensitiveAdminMutation(), mfaEnforcementEnabled(), withPrivilegedStepUp() (+2 more)

### Community 80 - "ERP Agendamentos"
Cohesion: 0.14
Nodes (13): Banco de Dados & Schema, CI, Endpoints principais, ERP Agendamentos, Estrutura, Pré-requisitos, Rodando com Docker (produção local), Rodando em desenvolvimento (sem banco) (+5 more)

### Community 81 - "router.go"
Cohesion: 0.09
Nodes (6): applyDatabaseSQL(), resolveSchemaFile(), CheckoutInput, ProductItem, CategoryTotal, Type

### Community 83 - "main.go"
Cohesion: 0.38
Nodes (4): newLazyAppHandler(), main(), resolveDurationEnv(), resolvePositiveIntEnv()

### Community 84 - "UseCase"
Cohesion: 0.23
Nodes (7): UseCase, NewCheckoutHandler(), AuditLogger, CheckoutResult, OutboxPublisher, TxRunner, CheckoutHandler

### Community 85 - "booking/src/lib/apiError.ts"
Cohesion: 0.26
Nodes (10): ApiError, ApiErrorPayload, buildApiError(), parseApiResponse(), API_ERROR_MESSAGES, EXACT_MESSAGE_MAP, getApiErrorMessage(), SUBSTRING_MESSAGE_PATTERNS (+2 more)

### Community 88 - "Repository"
Cohesion: 0.22
Nodes (4): NewQueue(), TestJobQueue(), Repository, NewJobRepository()

### Community 89 - "Claims"
Cohesion: 0.24
Nodes (5): isProductionAuthEnvironment(), Claims, decodeSegment(), encodeSegment(), JWTManager

### Community 91 - "PackagesPage"
Cohesion: 0.46
Nodes (8): deletePackage(), fetchPackageCredits(), fetchPackages(), setPackageActive(), PackagesPage(), handleDelete(), handleToggleActive(), loadPackages()

### Community 93 - "Analytics"
Cohesion: 0.29
Nodes (4): DayMetric, ProfessionalMetric, ServiceMetric, Analytics

### Community 94 - "UseCase"
Cohesion: 0.11
Nodes (11): MFAStatus, PasswordResetCodeConfirmInput, PasswordResetRequestInput, SecretCipher, generateTOTPSecret(), totpURI(), UseCase, hashResetToken() (+3 more)

### Community 95 - "observability.go"
Cohesion: 0.43
Nodes (6): isSensitiveLogKey(), redactSensitiveText(), sanitizeLogString(), sanitizeLogValue(), truncateLogString(), contextKey

### Community 96 - "email_provider.go"
Cohesion: 0.19
Nodes (8): buildNotificationProvider(), durationSecondsEnv(), intEnv(), logFallback(), NewEmailProvider(), NewLogProvider(), EmailProvider, LogProvider

### Community 98 - "BuildICalCalendar"
Cohesion: 0.67
Nodes (3): BuildICalCalendar(), formatUTC(), sanitizeLine()

### Community 99 - "BookingPage.tsx"
Cohesion: 0.15
Nodes (14): createPublicBooking(), formatMoney(), BookingPage(), handleBookingSubmit(), CenteredLayout(), DateTimeStep(), EmptyNote(), ErrorCard() (+6 more)

### Community 109 - "auth_context.go"
Cohesion: 0.25
Nodes (11): hasRole(), isAdminRole(), requestWithAuthClaims(), requireMasterAdmin(), requireOwnerOrAdmin(), requirePrivilegedAdmin(), requireRoles(), withProtectedNamespaces() (+3 more)

### Community 110 - "booking/src/shared/components/Button.tsx"
Cohesion: 0.33
Nodes (5): ButtonProps, ButtonSize, ButtonVariant, sizeClasses, variantClasses

### Community 115 - "PortalPage"
Cohesion: 0.39
Nodes (8): portalCancelAppointment(), portalLookup(), portalRescheduleAppointment(), PortalPage(), handleCancel(), handleConfirmReschedule(), handleLookup(), refresh()

### Community 116 - "ReminderWorker"
Cohesion: 0.39
Nodes (3): formatAppointmentLocal(), NewReminderWorker(), ReminderWorker

### Community 117 - "withCSRFProtection"
Cohesion: 0.32
Nodes (6): isLocalDevelopmentOrigin(), isOriginAllowed(), resolveAllowedOrigins(), withCORS(), requiresCSRFProtection(), withCSRFProtection()

## Knowledge Gaps
- **287 isolated node(s):** `graphify-mcp`, `name`, `private`, `version`, `type` (+282 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 547 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **62 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Container` connect `Container` to `Payment`, `container.go`, `Review`, `UseCase`, `Tenant`, `Service`, `database/sql.DB`, `Resource`, `Customer`, `UseCase`, `UseCase`, `CashRegister`, `Entry`, `Appointment`, `Connection`, `UseCase`, `UseCase`, `NewContainer`, `UseCase`, `Professional`, `AuditLog`, `EventBroker`, `Repository`, `context.Context`, `UseCase`, `UseCase`, `UseCase`, `Queue`?**
  _High betweenness centrality (0.042) - this node is a cross-community bridge._
- **Why does `RouterConfig` connect `RouterConfig` to `Payment`, `Review`, `UseCase`, `Tenant`, `Service`, `Resource`, `Customer`, `UseCase`, `UseCase`, `CashRegister`, `Entry`, `Appointment`, `Connection`, `UseCase`, `UseCase`, `UseCase`, `Professional`, `AuditLog`, `EventBroker`, `Repository`, `context.Context`, `UseCase`, `router.go`, `UseCase`, `time.Duration`, `UseCase`, `auth_context.go`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **Why does `repositorySet` connect `NewContainer` to `Payment`, `Review`, `UseCase`, `Tenant`, `Service`, `go_pkg_time`, `database/sql.DB`, `Resource`, `Customer`, `UseCase`, `UseCase`, `CashRegister`, `Entry`, `Appointment`, `auth/usecase.go`, `Connection`, `MemberPlan`, `UseCase`, `Package`, `Transaction`, `PackageCredit`, `GiftCard`, `Professional`, `Account`, `AuditLog`, `Repository`, `UseCase`, `Repository`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Are the 133 inferred relationships involving `writeJSON()` (e.g. with `health()` and `readiness()`) actually correct?**
  _`writeJSON()` has 133 INFERRED edges - model-reasoned connections that need verification._
- **Are the 104 inferred relationships involving `resolveRequestTenantID()` (e.g. with `.Analytics()` and `.CheckIn()`) actually correct?**
  _`resolveRequestTenantID()` has 104 INFERRED edges - model-reasoned connections that need verification._
- **What connects `graphify-mcp`, `name`, `private` to the rest of the system?**
  _287 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `adminFetch` be split into smaller, more focused modules?**
  _Cohesion score 0.06393442622950819 - nodes in this community are weakly interconnected._
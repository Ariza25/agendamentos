# Graph Report - erpagendamentos  (2026-11-02)

## Corpus Check
- 287 files · ~140,998 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 18 file(s) not represented in the graph (top: (none) 7, .conf 5, .css 3)

## Summary
- 2979 nodes · 9873 edges · 158 communities (98 shown, 60 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 697 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `5e80b524`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- adminFetch
- schema.sql
- ref_react
- Payment
- resolveRequestTenantID
- container.go
- IdempotencyStore
- PaymentsPage.tsx
- adminAuth.ts
- writeAPIError
- AppointmentsPage
- adminApi.ts
- Review
- UseCase
- AdminRoutes.tsx
- Tenant
- landing/src/App.tsx
- Service
- go_pkg_time
- .registerWithRole
- RouterConfig
- database/sql.DB
- Resource
- Customer
- UseCase
- InventoryPage.tsx
- UseCase
- CashRegister
- MembershipsPage.tsx
- Queue
- Entry
- admin/package.json
- booking/package.json
- UseCase
- FinancePage.tsx
- go_pkg_strings
- Connection
- rate_limit.go
- landing/package.json
- booking/src/pages/WaitingListPage.tsx
- UseCase
- net/http.ResponseWriter
- Configuration
- UseCase
- PortalPage.tsx
- useToast
- Repository
- AuthRepository
- UseCase
- net/http.Request
- bookingApi.ts
- memoryRepositories
- Membership
- Package
- UseCase
- writeJSON
- PackageCredit
- GiftCard
- Professional
- Account
- AuditLog
- EventBroker
- AuthGoogleHandler
- ConfigurationPage.tsx
- Appointment
- auth_cookie.go
- NewContainer
- templates.go
- .UpdateAccountRole
- compilerOptions
- AppointmentRepository
- compilerOptions
- compilerOptions
- OutboxRepository
- UseCase
- time.Time
- writeResourceError
- testing.T
- Container
- ERP Agendamentos
- bcrypt_hasher.go
- .GetByID
- main.go
- UseCase
- booking/src/lib/apiError.ts
- time.Duration
- ExportHandler
- writeCustomerError
- Claims
- SecretCipher
- TeamPage
- Analytics
- UseCase
- observability.go
- email_provider.go
- connection.go
- ical.go
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
- idempotencyRecorder
- PortalPage
- body_limit.go
- .ServeScript
- booking/src/main.tsx
- CommissionsPage.tsx

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

## Communities (158 total, 60 thin omitted)

### Community 0 - "adminFetch"
Cohesion: 0.06
Nodes (63): connectCalendar(), createAdminUser(), createCustomer(), createGiftCard(), createMemberPlan(), createMembership(), createPackage(), createPackageCredit() (+55 more)

### Community 1 - "schema.sql"
Cohesion: 0.05
Nodes (75): catalog.appointments, catalog.audit_logs, catalog.auth_rate_limits, catalog.auth_sessions, catalog.calendar_connections, catalog.cash_register_entries, catalog.cash_registers, catalog.commission_payouts (+67 more)

### Community 2 - "ref_react"
Cohesion: 0.09
Nodes (60): deleteGiftCard(), fetchGiftCards(), fetchNotifications(), fetchReviews(), fetchWaitingList(), updateWaitingListStatus(), APPOINTMENT_STATUS_LABELS, APPOINTMENT_STATUS_TONES (+52 more)

### Community 3 - "Payment"
Cohesion: 0.05
Nodes (23): UseCase, NewUseCase(), normalizeStatus(), buildPaymentProvider(), IntentResult, Payment, NewMercadoPagoProvider(), NewMockProvider() (+15 more)

### Community 4 - "resolveRequestTenantID"
Cohesion: 0.09
Nodes (9): resolveRequestTenantID(), writeCommerceError(), writeInternalError(), writeInvalidPayload(), CalendarSyncHandler, CommerceHandler, CommissionHandler, FinanceHandler (+1 more)

### Community 5 - "container.go"
Cohesion: 0.11
Nodes (3): deploymentEnvironment(), isProductionEnv(), CategoryTotal

### Community 6 - "IdempotencyStore"
Cohesion: 0.32
Nodes (3): NewIdempotencyStore(), idempotencyEntry, IdempotencyStore

### Community 7 - "PaymentsPage.tsx"
Cohesion: 0.16
Nodes (21): deleteCustomer(), downloadAdminCsv(), fetchCustomers(), fetchPayments(), markPaymentPaid(), refundPayment(), handleExport(), CustomersPage() (+13 more)

### Community 8 - "adminAuth.ts"
Cohesion: 0.07
Nodes (60): App(), handleLogout(), AppProviders(), AppProvidersProps, ResetPasswordPage(), handleSubmit(), AppRouter(), ADMIN_ROLES (+52 more)

### Community 9 - "writeAPIError"
Cohesion: 0.08
Nodes (27): isAdminRole(), requestWithAuthClaims(), requireMasterAdmin(), resolveAuthClaims(), resolveTenantID(), withProtectedNamespaces(), browserSessionPayload(), exposeAuthTokenForTests() (+19 more)

### Community 10 - "AppointmentsPage"
Cohesion: 0.15
Nodes (23): checkInAppointment(), deleteAppointment(), fetchAppointments(), fetchAppointmentSummary(), invalidateMutations(), markAppointmentNoShow(), updateAppointment(), updateAppointmentStatus() (+15 more)

### Community 11 - "adminApi.ts"
Cohesion: 0.05
Nodes (72): AppointmentFilters, checkSystemReadiness(), deleteTenant(), fetchAuditLogs(), fetchConfiguration(), fetchSystemStatus(), fetchTenants(), setTenantActive() (+64 more)

### Community 12 - "Review"
Cohesion: 0.12
Nodes (12): UseCase, NewUseCase(), Review, NewReviewRepository(), scanReviewRow(), scanReviewRows(), NewAppointmentHandler(), NewReviewHandler() (+4 more)

### Community 13 - "UseCase"
Cohesion: 0.07
Nodes (11): UseCase, NewUseCase(), TestStockUseCase(), Product, Repository, ServiceSupply, StockMovement, NewStockRepository() (+3 more)

### Community 14 - "AdminRoutes.tsx"
Cohesion: 0.08
Nodes (40): AccessDeniedPage(), AdminRouteFallback(), AdminRoutes(), AdminRoutesProps, AgendaPage, AnalyticsPage, AppointmentsPage, CommissionsPage (+32 more)

### Community 15 - "Tenant"
Cohesion: 0.08
Nodes (15): UseCase, NewUseCase(), ensureDefaultTenant(), Tenant, UpdateSubscriptionInput, defaultModules(), isTenantIdentityConflict(), defaultModules() (+7 more)

### Community 16 - "landing/src/App.tsx"
Cohesion: 0.11
Nodes (21): App(), CtaSection(), FEATURES, FeaturesSection(), Footer(), Header(), NAV_ITEMS, HeroSection() (+13 more)

### Community 17 - "Service"
Cohesion: 0.05
Nodes (20): MaintenanceWorker, MaintenanceConfig, NewMaintenanceWorker(), formatAppointmentLocal(), NewReminderWorker(), UseCase, NewUseCase(), Service (+12 more)

### Community 18 - "go_pkg_time"
Cohesion: 0.13
Nodes (3): CheckoutInput, ProductItem, Type

### Community 19 - ".registerWithRole"
Cohesion: 0.14
Nodes (12): ChangePasswordInput, PasswordResetCodeConfirmInput, PasswordResetConfirmInput, PasswordResetRequestInput, RegisterInput, ValidationError, hashResetToken(), NewValidationError() (+4 more)

### Community 20 - "RouterConfig"
Cohesion: 0.26
Nodes (27): requireAdmin(), NewRouter(), startBackgroundWorkers(), registerAppointmentRoutes(), registerConfigurationRoutes(), registerCustomerRoutes(), registerProfessionalRoutes(), registerResourceRoutes() (+19 more)

### Community 21 - "database/sql.DB"
Cohesion: 0.13
Nodes (30): applyDatabaseSQL(), resolveSchemaFile(), initializeRepositories(), postgresRepositories(), NewAppointmentRepository(), NewAuditRepository(), NewAuthRepository(), NewCalendarRepository() (+22 more)

### Community 22 - "Resource"
Cohesion: 0.11
Nodes (9): UseCase, NewUseCase(), Resource, scanResourceRow(), NewResourceHandler(), ResourceRepository, ResourceRepository, resourceScanner (+1 more)

### Community 23 - "Customer"
Cohesion: 0.09
Nodes (13): UseCase, Customer, NewCustomerRepository(), onlyDigits(), scanCustomerRow(), scanCustomerRows(), NewCustomerHandler(), NewExportHandler() (+5 more)

### Community 24 - "UseCase"
Cohesion: 0.06
Nodes (10): UseCase, CommissionRecord, CommissionRule, Payout, Repository, Summary, NewCommissionHandler(), PayoutInput (+2 more)

### Community 25 - "InventoryPage.tsx"
Cohesion: 0.16
Nodes (24): createProduct(), createStockMovement(), deleteProduct(), deleteService(), fetchProducts(), fetchServices(), fetchServiceSupplies(), fetchStockMovements() (+16 more)

### Community 26 - "UseCase"
Cohesion: 0.15
Nodes (10): UseCase, NewUseCase(), Log, NewNotificationRepository(), NewNotificationHandler(), NotificationRepository, Provider, Repository (+2 more)

### Community 27 - "CashRegister"
Cohesion: 0.10
Nodes (9): UseCase, NewUseCase(), CashRegister, CashRegisterEntry, Repository, ShiftSummary, NewCashRegisterHandler(), CashRegisterRepository (+1 more)

### Community 28 - "MembershipsPage.tsx"
Cohesion: 0.31
Nodes (13): deleteMembership(), fetchMemberPlans(), fetchMemberships(), setMembershipActive(), updateMemberPlanStatus(), MembershipsPage(), handleDelete(), handlePlanStatusChange() (+5 more)

### Community 29 - "Queue"
Cohesion: 0.08
Nodes (9): Queue, NewQueue(), TestJobQueue(), Job, Repository, NewJobRepository(), HandlerFunc, JobRepository (+1 more)

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
Cohesion: 0.18
Nodes (4): RecurrenceInput, UseCase, parseRFC3339(), validStatusChange()

### Community 34 - "FinancePage.tsx"
Cohesion: 0.16
Nodes (26): closeCashRegister(), createCashRegisterEntry(), createFinancialTransaction(), deleteFinancialTransaction(), fetchCashRegisterHistory(), fetchCurrentCashRegister(), fetchFinancialSummary(), fetchFinancialTransactions() (+18 more)

### Community 35 - "go_pkg_strings"
Cohesion: 0.08
Nodes (8): appendQueryParam(), secureTemporaryPassword(), Method, Status, authClaimsContextKey, googleStateData, googleTokenResponse, tokenPayload

### Community 36 - "Connection"
Cohesion: 0.09
Nodes (10): UseCase, NewUseCase(), TestCalendarUseCase(), Connection, Repository, SyncEvent, NewCalendarRepository(), NewCalendarSyncHandler() (+2 more)

### Community 37 - "rate_limit.go"
Cohesion: 0.17
Nodes (19): buildCompoundIdempotencyKey(), forwardedClientIP(), isSensitiveRateLimitedRoute(), isTrustedProxy(), newMemoryRateLimiter(), normalizedClientID(), persistentRateLimitConfig(), rateLimitRouteKey() (+11 more)

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
Cohesion: 0.10
Nodes (16): writeConfigurationError(), writeNotificationError(), writeProfessionalError(), writeBadRequestError(), writeConflictError(), writeNotFoundError(), writeReviewError(), writeServiceError() (+8 more)

### Community 42 - "Configuration"
Cohesion: 0.13
Nodes (17): workingHoursFor(), UseCase, NewUseCase(), Configuration, cloneConfiguration(), defaultMemoryConfiguration(), normalizeConfigurationTenantID(), normalizePostgresTenantID() (+9 more)

### Community 43 - "UseCase"
Cohesion: 0.36
Nodes (5): UseCase, NewAuthGoogleHandler(), NewOnboardingHandler(), SignupInput, OnboardingHandler

### Community 44 - "PortalPage.tsx"
Cohesion: 0.18
Nodes (21): PortalAppointment, PortalResult, Slot, formatDateLong(), formatDateTime(), formatTime(), fromDateKey(), nextDateKeys() (+13 more)

### Community 45 - "useToast"
Cohesion: 0.09
Nodes (62): handleCopyBookingUrl(), EstablishmentModal(), handleCopyLink(), handleNameChange(), handleSlugChange(), handleSubmit(), EstablishmentModalProps, getPublicBookingUrl() (+54 more)

### Community 46 - "Repository"
Cohesion: 0.09
Nodes (8): Repository, PasswordHasher, seedAdminAccount(), seedAccount(), seedAppointments(), seedHomologationData(), AuthState, authComponents

### Community 47 - "AuthRepository"
Cohesion: 0.08
Nodes (7): DeviceSession, PasswordResetToken, cloneAuthAccount(), resetTokenExpired(), sessionExpired(), AuthRepository, rateLimitEntry

### Community 48 - "UseCase"
Cohesion: 0.19
Nodes (8): appointmentDurationMinutes(), generateManagementCode(), UseCase, hashManagementCode(), BookingInput, PortalAppointment, PortalCustomer, PortalResult

### Community 49 - "net/http.Request"
Cohesion: 0.14
Nodes (8): writeAppointmentError(), resolveActorTenantID(), buildPasswordResetURL(), AuthHandler, writeAdminAuthError(), newAuthMFAHandler(), HandleMetrics(), AppointmentHandler

### Community 50 - "bookingApi.ts"
Cohesion: 0.14
Nodes (25): backendOrigin, buildApiUrl(), AvailabilityResponse, BookingResult, cancelPublicAppointment(), fetchAvailability(), fetchPublicAppointment(), fetchPublicInfo() (+17 more)

### Community 51 - "memoryRepositories"
Cohesion: 0.16
Nodes (23): NewUseCase(), TestCashRegisterFlow(), TestCheckoutFlow(), NewUseCase(), TestCommissionUseCase(), NewUseCase(), TestOnboardingSignup(), buildAuthComponents() (+15 more)

### Community 52 - "Membership"
Cohesion: 0.12
Nodes (5): Membership, NewMembershipRepository(), MembershipRepository, MembershipRepository, MembershipRepository

### Community 53 - "Package"
Cohesion: 0.12
Nodes (5): Package, NewPackageRepository(), PackageRepository, PackageRepository, PackageRepository

### Community 54 - "UseCase"
Cohesion: 0.08
Nodes (11): UseCase, NewUseCase(), TestFinanceUseCase(), Repository, Summary, Transaction, NewFinanceHandler(), NewMercadoPagoWebhookHandler() (+3 more)

### Community 55 - "writeJSON"
Cohesion: 0.15
Nodes (7): writeBookingError(), health(), writePaymentError(), writeJSON(), BookingHandler, CashRegisterHandler, PaymentHandler

### Community 56 - "PackageCredit"
Cohesion: 0.12
Nodes (6): PackageCredit, NewPackageCreditRepository(), scanPackageCredits(), PackageCreditRepository, PackageCreditRepository, PackageCreditRepository

### Community 57 - "GiftCard"
Cohesion: 0.11
Nodes (6): generateGiftCardCode(), GiftCard, NewGiftCardRepository(), GiftCardRepository, GiftCardRepository, GiftCardRepository

### Community 58 - "Professional"
Cohesion: 0.10
Nodes (10): UseCase, NewUseCase(), Professional, scanProfessionalRow(), scanProfessionalRows(), NewProfessionalHandler(), ProfessionalRepository, ProfessionalRepository (+2 more)

### Community 59 - "Account"
Cohesion: 0.14
Nodes (4): Account, UpdateProfileInput, updateAuthPasswordAndVersion(), AuthRepository

### Community 60 - "AuditLog"
Cohesion: 0.20
Nodes (7): AuditLog, Repository, UseCase, NewAuditHandler(), AuditRepository, AuditRepository, AuditHandler

### Community 62 - "EventBroker"
Cohesion: 0.21
Nodes (8): NewBookingHandler(), NewCheckoutHandler(), EventBroker, NewEventBroker(), RegisterMetricsBroker(), CheckoutHandler, Event, MetricsRegistry

### Community 63 - "AuthGoogleHandler"
Cohesion: 0.31
Nodes (5): fetchGoogleUserInfo(), resolveFrontendURL(), resolveGoogleRedirectURI(), AuthGoogleHandler, googleUserInfo

### Community 64 - "ConfigurationPage.tsx"
Cohesion: 0.12
Nodes (11): updateConfiguration(), backendOrigin, DAY_LABELS, ActiveTab, CHANNEL_OPTIONS, ConfigurationPage(), handleSubmit(), EMPTY_DAY (+3 more)

### Community 65 - "Appointment"
Cohesion: 0.11
Nodes (7): Repository, Appointment, StatusCount, parseAppointmentBound(), parseAppointmentInterval(), parseTimeFlexible(), AppointmentRepository

### Community 66 - "auth_cookie.go"
Cohesion: 0.23
Nodes (19): authCookieDomain(), authCookiePartitioned(), authCookieSameSite(), authCookieSecure(), authTokenFromCookie(), clearAuthCookie(), csrfTokenFromCookie(), isAuthCookieProductionEnv() (+11 more)

### Community 67 - "NewContainer"
Cohesion: 0.18
Nodes (9): NewUseCase(), NewUseCase(), NewUseCase(), buildNotificationProvider(), NewContainer(), loadEnvFile(), TestSeedHomologationData(), NewLogProvider() (+1 more)

### Community 68 - "templates.go"
Cohesion: 0.27
Nodes (16): greeting(), PublicBookingOrigin(), RenderCancellation(), RenderConfirmation(), RenderReminder(), RenderReschedule(), sampleData(), TestGreetingWithoutName() (+8 more)

### Community 69 - ".UpdateAccountRole"
Cohesion: 0.39
Nodes (4): canManageTenantAccount(), isPrivilegedAdminRole(), isStaffRole(), normalizeTenantScope()

### Community 70 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 71 - "AppointmentRepository"
Cohesion: 0.19
Nodes (5): appointmentSelectColumns(), scanAppointmentRow(), scanAppointmentRows(), AppointmentRepository, appointmentScanner

### Community 72 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 73 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 74 - "OutboxRepository"
Cohesion: 0.12
Nodes (5): quoteLiteral(), OutboxRepository, OutboxEvent, Repository, OutboxRepository

### Community 75 - "UseCase"
Cohesion: 0.20
Nodes (11): NewUseCase(), TestAuditLogRecordingAndListing(), UseCase, NewUseCase(), TestOutboxPublishAndDispatch(), NewAuditRepository(), NewOutboxRepository(), NewSystemHandler() (+3 more)

### Community 76 - "time.Time"
Cohesion: 0.35
Nodes (9): countOverlappingSlots(), formatAppointmentForTemplate(), isClosedDay(), localMidnight(), overlapsAnyBusy(), parseClock(), parseRFC3339(), ProfessionalSlots (+1 more)

### Community 78 - "testing.T"
Cohesion: 0.13
Nodes (22): newAppointment(), newResourceUseCase(), TestCreateAutoAssignsAvailableResource(), TestCreateRejectsBusyExplicitResource(), TestCreateRejectsUnknownResource(), TestUpdateAutoReassignsFreedResource(), newTestUseCase(), TestAvailableSlotsRespectsClosedDays() (+14 more)

### Community 79 - "Container"
Cohesion: 0.13
Nodes (9): lazyAppHandler, writeLazyJSON(), Container, isPrivilegedAdminRole(), isSensitiveAdminMutation(), mfaEnforcementEnabled(), withPrivilegedStepUp(), Router (+1 more)

### Community 80 - "ERP Agendamentos"
Cohesion: 0.14
Nodes (13): Banco de Dados & Schema, CI, Endpoints principais, ERP Agendamentos, Estrutura, Pré-requisitos, Rodando com Docker (produção local), Rodando em desenvolvimento (sem banco) (+5 more)

### Community 82 - ".GetByID"
Cohesion: 0.31
Nodes (3): nullIfEmpty(), scanMemberPlans(), MemberPlanRepository

### Community 83 - "main.go"
Cohesion: 0.38
Nodes (4): newLazyAppHandler(), main(), resolveDurationEnv(), resolvePositiveIntEnv()

### Community 84 - "UseCase"
Cohesion: 0.22
Nodes (6): UseCase, NewUseCase(), AuditLogger, CheckoutResult, OutboxPublisher, TxRunner

### Community 85 - "booking/src/lib/apiError.ts"
Cohesion: 0.26
Nodes (10): ApiError, ApiErrorPayload, buildApiError(), parseApiResponse(), API_ERROR_MESSAGES, EXACT_MESSAGE_MAP, getApiErrorMessage(), SUBSTRING_MESSAGE_PATTERNS (+2 more)

### Community 89 - "Claims"
Cohesion: 0.23
Nodes (6): TokenManager, isProductionAuthEnvironment(), Claims, decodeSegment(), encodeSegment(), JWTManager

### Community 91 - "TeamPage"
Cohesion: 0.52
Nodes (7): deleteAdminUser(), fetchAdminUsers(), updateAdminUserRole(), TeamPage(), handleDelete(), handleRoleChange(), loadUsers()

### Community 93 - "Analytics"
Cohesion: 0.33
Nodes (4): DayMetric, ProfessionalMetric, ServiceMetric, Analytics

### Community 94 - "UseCase"
Cohesion: 0.12
Nodes (11): LoginInput, MFAStatus, SecretCipher, Session, generateTOTPSecret(), totpCode(), totpURI(), validateTOTP() (+3 more)

### Community 95 - "observability.go"
Cohesion: 0.52
Nodes (6): isSensitiveLogKey(), redactSensitiveText(), sanitizeLogString(), sanitizeLogValue(), truncateLogString(), contextKey

### Community 96 - "email_provider.go"
Cohesion: 0.47
Nodes (3): logFallback(), NewEmailProvider(), EmailProvider

### Community 97 - "connection.go"
Cohesion: 0.60
Nodes (3): Open(), resolveDurationEnv(), resolveIntEnv()

### Community 98 - "ical.go"
Cohesion: 0.83
Nodes (3): BuildICalCalendar(), formatUTC(), sanitizeLine()

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

### Community 116 - "body_limit.go"
Cohesion: 0.83
Nodes (3): resolveMaxRequestBodyBytes(), shouldLimitRequestBody(), withRequestBodyLimit()

### Community 117 - ".ServeScript"
Cohesion: 0.50
Nodes (3): NewWidgetHandler(), resolveWidgetBookingOrigin(), WidgetHandler

### Community 119 - "CommissionsPage.tsx"
Cohesion: 0.14
Nodes (29): createCommissionPayout(), createCommissionRule(), deleteCommissionRule(), fetchAppointmentAnalytics(), fetchCommissionPayouts(), fetchCommissionRecords(), fetchCommissionRules(), fetchCommissionSummaries() (+21 more)

## Knowledge Gaps
- **287 isolated node(s):** `graphify-mcp`, `name`, `private`, `version`, `type` (+282 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 548 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **60 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Container` connect `Container` to `Payment`, `container.go`, `Review`, `UseCase`, `Tenant`, `Service`, `database/sql.DB`, `Resource`, `Customer`, `UseCase`, `UseCase`, `CashRegister`, `Queue`, `Entry`, `UseCase`, `Connection`, `UseCase`, `Configuration`, `UseCase`, `UseCase`, `memoryRepositories`, `UseCase`, `Professional`, `AuditLog`, `EventBroker`, `NewContainer`, `UseCase`, `UseCase`, `UseCase`?**
  _High betweenness centrality (0.043) - this node is a cross-community bridge._
- **Why does `RouterConfig` connect `RouterConfig` to `Payment`, `container.go`, `Review`, `UseCase`, `Tenant`, `Service`, `Resource`, `Customer`, `UseCase`, `UseCase`, `CashRegister`, `Entry`, `UseCase`, `Connection`, `UseCase`, `Configuration`, `UseCase`, `UseCase`, `UseCase`, `Professional`, `AuditLog`, `EventBroker`, `UseCase`, `UseCase`, `time.Duration`, `UseCase`, `requireRoles`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **Why does `UseCase` connect `UseCase` to `resolveRequestTenantID`, `container.go`, `Review`, `Container`, `net/http.Request`, `memoryRepositories`, `RouterConfig`, `UseCase`, `Professional`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Are the 133 inferred relationships involving `writeJSON()` (e.g. with `health()` and `readiness()`) actually correct?**
  _`writeJSON()` has 133 INFERRED edges - model-reasoned connections that need verification._
- **Are the 104 inferred relationships involving `resolveRequestTenantID()` (e.g. with `.Analytics()` and `.CheckIn()`) actually correct?**
  _`resolveRequestTenantID()` has 104 INFERRED edges - model-reasoned connections that need verification._
- **What connects `graphify-mcp`, `name`, `private` to the rest of the system?**
  _287 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `adminFetch` be split into smaller, more focused modules?**
  _Cohesion score 0.06093189964157706 - nodes in this community are weakly interconnected._
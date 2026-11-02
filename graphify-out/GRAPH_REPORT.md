# Graph Report - erpagendamentos  (2026-11-01)

## Corpus Check
- 266 files · ~127,642 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 16 file(s) not represented in the graph (top: (none) 7, .conf 5, .css 3)

## Summary
- 2807 nodes · 9128 edges · 149 communities (90 shown, 59 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 673 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `4e939256`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- AppointmentsPage.tsx
- auth_cookie.go
- container.go
- Payment
- writeAPIError
- schema.sql
- Button
- adminFetch
- adminAuth.ts
- ref_react
- resolveRequestTenantID
- UseCase
- ProfessionalsPage.tsx
- auth/usecase.go
- UseCase
- Tenant
- PageContainer
- writeJSON
- AdminRoutes.tsx
- UseCase
- Service
- Customer
- Queue
- Professional
- Entry
- admin/package.json
- database/sql.DB
- booking/package.json
- Appointment
- Review
- Connection
- landing/package.json
- booking/src/pages/WaitingListPage.tsx
- landing/src/App.tsx
- MemberPlan
- AuthRepository
- CashRegister
- PortalPage.tsx
- UseCase
- UseCase
- Package
- net/http.ResponseWriter
- bookingApi.ts
- Resource
- GiftCard
- PackageCredit
- Account
- AuthRepository
- time.Duration
- Container
- BookingPage.tsx
- ServicesPage.tsx
- .registerWithRole
- adminApi.ts
- UseCase
- repositorySet
- time.Time
- templates.go
- compilerOptions
- testing.T
- AppointmentRepository
- compilerOptions
- compilerOptions
- ReminderWorker
- Configuration
- memoryRepositories
- observability.go
- net/http.Request
- ExportHandler
- ERP Agendamentos
- MembershipsPage.tsx
- UseCase
- Claims
- RouterConfig
- NewContainer
- booking/src/lib/apiError.ts
- TenantsPage.tsx
- rate_limit.go
- EventBroker
- withIdempotency
- PortalPage
- checkout/usecase.go
- writeProfessionalError
- writeServiceError
- PaymentsPage.tsx
- postgres/commerce_repository.go
- writePaymentError
- Analytics
- ServiceRepository
- bcrypt_hasher.go
- booking/src/shared/components/Button.tsx
- SecretCipher
- body_limit.go
- booking/src/main.tsx
- widget_handler.go
- graphify
- rules/graphify.md
- workflows/graphify.md
- generate-certs.sh
- backup_postgres.sh
- certbot_ssl.sh
- erpagendamentos/backend
- CustomerRepository
- writeNotFoundError
- .UpdateAccountRole
- writeResourceError
- requireRoles

## God Nodes (most connected - your core abstractions)
1. `writeJSON()` - 132 edges
2. `resolveRequestTenantID()` - 108 edges
3. `adminFetch()` - 74 edges
4. `Button` - 72 edges
5. `writeInvalidPayload()` - 71 edges
6. `useToast()` - 68 edges
7. `UseCase` - 67 edges
8. `invalidateAdminGetCache()` - 64 edges
9. `writeInternalError()` - 58 edges
10. `RouterConfig` - 54 edges

## Surprising Connections (you probably didn't know these)
- `handleExport()` --calls--> `downloadAdminCsv()`  [EXTRACTED]
  admin/src/pages/ReviewsPage.tsx → admin/src/lib/adminApi.ts
- `newResourceUseCase()` --calls--> `NewUseCase()`  [INFERRED]
  backend/internal/application/appointment/resource_test.go → backend/internal/application/appointment/usecase.go
- `newTestUseCase()` --calls--> `NewUseCase()`  [INFERRED]
  backend/internal/application/booking/usecase_test.go → backend/internal/application/booking/usecase.go
- `TestManagementCodeNormalization()` --calls--> `hashManagementCode()`  [INFERRED]
  backend/internal/application/booking/usecase_test.go → backend/internal/application/booking/usecase.go
- `buildAuthComponents()` --calls--> `isProductionEnv()`  [INFERRED]
  backend/internal/bootstrap/auth_components.go → backend/internal/bootstrap/container.go

## Import Cycles
- None detected.

## Communities (149 total, 59 thin omitted)

### Community 0 - "AppointmentsPage.tsx"
Cohesion: 0.09
Nodes (48): AppointmentFormModal(), computeEndsAt(), handleSubmit(), AppointmentFormModalProps, STATUS_OPTIONS, CheckoutComandaModal(), handleSubmit(), checkInAppointment() (+40 more)

### Community 1 - "auth_cookie.go"
Cohesion: 0.11
Nodes (28): lazyAppHandler, writeLazyJSON(), bearerAuthAllowed(), isPrivilegedAdminRole(), authCookieDomain(), authCookiePartitioned(), authCookieSameSite(), authCookieSecure() (+20 more)

### Community 2 - "container.go"
Cohesion: 0.10
Nodes (4): deploymentEnvironment(), isProductionEnv(), CategoryTotal, Type

### Community 3 - "Payment"
Cohesion: 0.05
Nodes (23): UseCase, NewUseCase(), normalizeStatus(), buildPaymentProvider(), IntentResult, Payment, NewMercadoPagoProvider(), NewMockProvider() (+15 more)

### Community 4 - "writeAPIError"
Cohesion: 0.09
Nodes (26): requestWithAuthClaims(), requireMasterAdmin(), resolveAuthClaims(), resolveTenantID(), withProtectedNamespaces(), browserSessionPayload(), exposeAuthTokenForTests(), AuthHandler (+18 more)

### Community 5 - "schema.sql"
Cohesion: 0.06
Nodes (70): catalog.appointments, catalog.auth_rate_limits, catalog.auth_sessions, catalog.calendar_connections, catalog.cash_register_entries, catalog.cash_registers, catalog.commission_payouts, catalog.commission_records (+62 more)

### Community 6 - "Button"
Cohesion: 0.07
Nodes (50): AdminShell(), AdminShellProps, SaaSPlansUpgradeModal(), AdminNavigationItem, navigationGroups, navigationItems, createAdminUser(), createPackageCredit() (+42 more)

### Community 7 - "adminFetch"
Cohesion: 0.09
Nodes (46): closeCashRegister(), connectCalendar(), createCashRegisterEntry(), createFinancialTransaction(), createPackage(), createPayment(), deleteFinancialTransaction(), disconnectCalendar() (+38 more)

### Community 8 - "adminAuth.ts"
Cohesion: 0.08
Nodes (47): App(), handleLogout(), AppProviders(), AppProvidersProps, ResetPasswordPage(), handleSubmit(), AppRouter(), ADMIN_ROLES (+39 more)

### Community 9 - "ref_react"
Cohesion: 0.10
Nodes (42): createCustomer(), createGiftCard(), deleteCustomer(), deleteGiftCard(), deletePackage(), fetchCustomers(), fetchGiftCards(), fetchPackageCredits() (+34 more)

### Community 10 - "resolveRequestTenantID"
Cohesion: 0.10
Nodes (9): isAdminRole(), resolveRequestTenantID(), writeInternalError(), writeWaitingListError(), CalendarSyncHandler, CommissionHandler, FinanceHandler, StockHandler (+1 more)

### Community 11 - "UseCase"
Cohesion: 0.06
Nodes (10): UseCase, CommissionRecord, CommissionRule, Payout, Repository, Summary, NewCommissionHandler(), PayoutInput (+2 more)

### Community 12 - "ProfessionalsPage.tsx"
Cohesion: 0.33
Nodes (11): createProfessional(), deleteProfessional(), setProfessionalActive(), updateProfessional(), COLOR_OPTIONS, ProfessionalFormModal(), handleSubmit(), ProfessionalsPage() (+3 more)

### Community 13 - "auth/usecase.go"
Cohesion: 0.09
Nodes (4): Method, Status, authClaimsContextKey, tokenPayload

### Community 14 - "UseCase"
Cohesion: 0.07
Nodes (9): UseCase, NewUseCase(), Product, Repository, ServiceSupply, StockMovement, NewStockHandler(), StockRepository (+1 more)

### Community 15 - "Tenant"
Cohesion: 0.08
Nodes (16): UseCase, NewUseCase(), ensureDefaultTenant(), Tenant, UpdateSubscriptionInput, defaultModules(), isTenantIdentityConflict(), defaultModules() (+8 more)

### Community 16 - "PageContainer"
Cohesion: 0.08
Nodes (43): fetchAppointmentAnalytics(), fetchCurrentTenant(), fetchNotifications(), fetchReviews(), fetchWaitingList(), updateWaitingListStatus(), formatDate(), fromDateKey() (+35 more)

### Community 17 - "writeJSON"
Cohesion: 0.14
Nodes (9): writeBookingError(), health(), writeConflictError(), writeJSON(), writeReviewError(), writeTenantError(), BookingHandler, ReviewHandler (+1 more)

### Community 18 - "AdminRoutes.tsx"
Cohesion: 0.15
Nodes (25): AccessDeniedPage(), AdminRouteFallback(), AdminRoutes(), AdminRoutesProps, AgendaPage, AnalyticsPage, AppointmentsPage, CommissionsPage (+17 more)

### Community 19 - "UseCase"
Cohesion: 0.11
Nodes (12): MFAStatus, PasswordResetCodeConfirmInput, PasswordResetRequestInput, SecretCipher, generateTOTPSecret(), totpURI(), UseCase, hashResetToken() (+4 more)

### Community 20 - "Service"
Cohesion: 0.12
Nodes (6): UseCase, NewUseCase(), Service, NewServiceHandler(), ServiceRepository, Repository

### Community 22 - "Customer"
Cohesion: 0.12
Nodes (9): UseCase, NewUseCase(), Customer, NewCustomerRepository(), onlyDigits(), NewCustomerHandler(), NewExportHandler(), Repository (+1 more)

### Community 23 - "Queue"
Cohesion: 0.08
Nodes (6): Queue, Job, Repository, HandlerFunc, JobRepository, JobRepository

### Community 24 - "Professional"
Cohesion: 0.10
Nodes (10): UseCase, NewUseCase(), Professional, scanProfessionalRow(), scanProfessionalRows(), NewProfessionalHandler(), ProfessionalRepository, ProfessionalRepository (+2 more)

### Community 25 - "Entry"
Cohesion: 0.11
Nodes (11): UseCase, NewUseCase(), Entry, NewWaitingListRepository(), scanWaitingRow(), scanWaitingRows(), NewWaitingListHandler(), WaitingListRepository (+3 more)

### Community 26 - "admin/package.json"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+24 more)

### Community 28 - "database/sql.DB"
Cohesion: 0.12
Nodes (29): applyDatabaseSQL(), resolveSchemaFile(), initializeRepositories(), postgresRepositories(), NewAppointmentRepository(), NewAuthRepository(), NewCalendarRepository(), NewCashRegisterRepository() (+21 more)

### Community 29 - "booking/package.json"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+24 more)

### Community 30 - "Appointment"
Cohesion: 0.11
Nodes (7): Repository, Appointment, StatusCount, parseAppointmentBound(), parseAppointmentInterval(), BookingResult, AppointmentRepository

### Community 31 - "Review"
Cohesion: 0.12
Nodes (11): UseCase, NewUseCase(), Review, NewReviewRepository(), scanReviewRow(), scanReviewRows(), NewAppointmentHandler(), ReviewRepository (+3 more)

### Community 32 - "Connection"
Cohesion: 0.09
Nodes (10): UseCase, NewUseCase(), TestCalendarUseCase(), Connection, Repository, SyncEvent, NewCalendarRepository(), NewCalendarSyncHandler() (+2 more)

### Community 33 - "landing/package.json"
Cohesion: 0.06
Nodes (30): dependencies, lucide-react, react, react-dom, devDependencies, tailwindcss, @tailwindcss/vite, @types/react (+22 more)

### Community 34 - "booking/src/pages/WaitingListPage.tsx"
Cohesion: 0.17
Nodes (22): App(), fetchPublicReview(), signUpWaitingList(), submitPublicReview(), DetailsStep(), ReviewPage(), handleSubmit(), SlugEntryPage() (+14 more)

### Community 35 - "landing/src/App.tsx"
Cohesion: 0.11
Nodes (21): App(), CtaSection(), FEATURES, FeaturesSection(), Footer(), Header(), NAV_ITEMS, HeroSection() (+13 more)

### Community 36 - "MemberPlan"
Cohesion: 0.16
Nodes (5): NewUseCase(), MemberPlan, NewMemberPlanRepository(), MemberPlanRepository, MemberPlanRepository

### Community 37 - "AuthRepository"
Cohesion: 0.08
Nodes (6): AuthState, DeviceSession, PasswordResetToken, resetTokenExpired(), sessionExpired(), AuthRepository

### Community 38 - "CashRegister"
Cohesion: 0.10
Nodes (9): UseCase, NewUseCase(), CashRegister, CashRegisterEntry, Repository, ShiftSummary, NewCashRegisterHandler(), CashRegisterRepository (+1 more)

### Community 39 - "PortalPage.tsx"
Cohesion: 0.19
Nodes (23): PortalAppointment, PortalResult, Slot, formatDateLong(), formatDateTime(), formatTime(), fromDateKey(), nextDateKeys() (+15 more)

### Community 40 - "UseCase"
Cohesion: 0.18
Nodes (10): appointmentDurationMinutes(), formatAppointmentForTemplate(), generateManagementCode(), UseCase, hashManagementCode(), parseRFC3339(), BookingInput, PortalAppointment (+2 more)

### Community 41 - "UseCase"
Cohesion: 0.10
Nodes (6): UseCase, Membership, NewMembershipRepository(), NewCommerceHandler(), MembershipRepository, MembershipRepository

### Community 42 - "Package"
Cohesion: 0.12
Nodes (5): Package, NewPackageRepository(), PackageRepository, PackageRepository, PackageRepository

### Community 43 - "net/http.ResponseWriter"
Cohesion: 0.15
Nodes (6): writeCommerceError(), writeConfigurationError(), writeInvalidPayload(), CashRegisterHandler, CommerceHandler, responseRecorder

### Community 44 - "bookingApi.ts"
Cohesion: 0.14
Nodes (23): backendOrigin, buildApiUrl(), AvailabilityResponse, cancelPublicAppointment(), fetchAvailability(), fetchPublicAppointment(), fetchPublicInfo(), fetchPublicProfessionals() (+15 more)

### Community 45 - "Resource"
Cohesion: 0.11
Nodes (9): UseCase, NewUseCase(), Resource, scanResourceRow(), NewResourceHandler(), ResourceRepository, ResourceRepository, resourceScanner (+1 more)

### Community 46 - "GiftCard"
Cohesion: 0.11
Nodes (6): generateGiftCardCode(), GiftCard, NewGiftCardRepository(), GiftCardRepository, GiftCardRepository, GiftCardRepository

### Community 47 - "PackageCredit"
Cohesion: 0.13
Nodes (6): PackageCredit, NewPackageCreditRepository(), scanPackageCredits(), PackageCreditRepository, PackageCreditRepository, PackageCreditRepository

### Community 48 - "Account"
Cohesion: 0.11
Nodes (4): Repository, Account, UpdateProfileInput, cloneAuthAccount()

### Community 50 - "time.Duration"
Cohesion: 0.17
Nodes (8): MaintenanceConfig, buildAuthComponents(), boolEnv(), durationSecondsEnv(), intEnv(), CleanupSummary, NewJWTManager(), authComponents

### Community 51 - "Container"
Cohesion: 0.16
Nodes (12): UseCase, NewUseCase(), UseCase, NewUseCase(), Container, NewCalendarHandler(), NewConfigurationHandler(), NewOnboardingHandler() (+4 more)

### Community 52 - "BookingPage.tsx"
Cohesion: 0.13
Nodes (15): BookingResult, createPublicBooking(), formatMoney(), BookingPage(), handleBookingSubmit(), CenteredLayout(), EmptyNote(), ErrorCard() (+7 more)

### Community 53 - "ServicesPage.tsx"
Cohesion: 0.11
Nodes (34): createResource(), createService(), deleteResource(), deleteService(), fetchResources(), setResourceActive(), setServiceActive(), updateResource() (+26 more)

### Community 54 - ".registerWithRole"
Cohesion: 0.14
Nodes (13): ChangePasswordInput, LoginInput, PasswordResetConfirmInput, RegisterInput, Session, ValidationError, totpCode(), validateTOTP() (+5 more)

### Community 55 - "adminApi.ts"
Cohesion: 0.06
Nodes (71): CheckoutComandaModalProps, SelectedProduct, AppointmentFilters, createCommissionPayout(), createCommissionRule(), createProduct(), createStockMovement(), deleteCommissionRule() (+63 more)

### Community 56 - "UseCase"
Cohesion: 0.18
Nodes (4): RecurrenceInput, UseCase, parseRFC3339(), validStatusChange()

### Community 57 - "repositorySet"
Cohesion: 0.13
Nodes (11): UseCase, NewUseCase(), Log, NewNotificationRepository(), NewNotificationHandler(), repositorySet, NotificationRepository, Provider (+3 more)

### Community 58 - "time.Time"
Cohesion: 0.23
Nodes (11): countOverlappingSlots(), isClosedDay(), localMidnight(), overlapsAnyBusy(), parseClock(), BuildICalCalendar(), formatUTC(), sanitizeLine() (+3 more)

### Community 59 - "templates.go"
Cohesion: 0.27
Nodes (16): greeting(), PublicBookingOrigin(), RenderCancellation(), RenderConfirmation(), RenderReminder(), RenderReschedule(), sampleData(), TestGreetingWithoutName() (+8 more)

### Community 60 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 61 - "testing.T"
Cohesion: 0.18
Nodes (16): newAppointment(), newResourceUseCase(), TestCreateAutoAssignsAvailableResource(), TestCreateRejectsBusyExplicitResource(), TestCreateRejectsUnknownResource(), TestUpdateAutoReassignsFreedResource(), newTestUseCase(), TestAvailableSlotsRespectsClosedDays() (+8 more)

### Community 62 - "AppointmentRepository"
Cohesion: 0.17
Nodes (6): appointmentSelectColumns(), nullIfEmpty(), scanAppointmentRow(), scanAppointmentRows(), AppointmentRepository, appointmentScanner

### Community 63 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 64 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 65 - "ReminderWorker"
Cohesion: 0.22
Nodes (5): MaintenanceWorker, NewMaintenanceWorker(), formatAppointmentLocal(), NewReminderWorker(), ReminderWorker

### Community 66 - "Configuration"
Cohesion: 0.25
Nodes (10): workingHoursFor(), Configuration, cloneConfiguration(), defaultMemoryConfiguration(), normalizeConfigurationTenantID(), normalizePostgresTenantID(), ClosedDay, WorkingHours (+2 more)

### Community 67 - "memoryRepositories"
Cohesion: 0.26
Nodes (15): TestCashRegisterFlow(), TestCheckoutFlow(), TestCommissionUseCase(), TestOnboardingSignup(), memoryRepositories(), NewAppointmentRepository(), NewAuthRepository(), NewCashRegisterRepository() (+7 more)

### Community 68 - "observability.go"
Cohesion: 0.10
Nodes (16): newLazyAppHandler(), main(), resolveDurationEnv(), resolvePositiveIntEnv(), buildNotificationProvider(), logFallback(), NewEmailProvider(), NewLogProvider() (+8 more)

### Community 69 - "net/http.Request"
Cohesion: 0.12
Nodes (11): writeAppointmentError(), resolveActorTenantID(), buildPasswordResetURL(), AuthHandler, secureTemporaryPassword(), writeAdminAuthError(), newAuthMFAHandler(), HandleMetrics() (+3 more)

### Community 71 - "ERP Agendamentos"
Cohesion: 0.14
Nodes (13): Banco de Dados & Schema, CI, Endpoints principais, ERP Agendamentos, Estrutura, Pré-requisitos, Rodando com Docker (produção local), Rodando em desenvolvimento (sem banco) (+5 more)

### Community 72 - "MembershipsPage.tsx"
Cohesion: 0.19
Nodes (21): createMemberPlan(), createMembership(), deleteMembership(), fetchMemberPlans(), fetchMemberships(), setMembershipActive(), updateMemberPlanStatus(), updateMembership() (+13 more)

### Community 73 - "UseCase"
Cohesion: 0.09
Nodes (9): UseCase, Repository, Summary, Transaction, NewFinanceHandler(), NewMercadoPagoWebhookHandler(), FinanceRepository, FinanceRepository (+1 more)

### Community 74 - "Claims"
Cohesion: 0.21
Nodes (7): TokenManager, isProductionAuthEnvironment(), NewUseCase(), Claims, decodeSegment(), encodeSegment(), JWTManager

### Community 75 - "RouterConfig"
Cohesion: 0.23
Nodes (28): requireAdmin(), NewReviewHandler(), NewRouter(), startBackgroundWorkers(), registerAppointmentRoutes(), registerConfigurationRoutes(), registerCustomerRoutes(), registerProfessionalRoutes() (+20 more)

### Community 76 - "NewContainer"
Cohesion: 0.18
Nodes (11): NewUseCase(), NewUseCase(), NewUseCase(), NewUseCase(), NewUseCase(), TestFinanceUseCase(), NewQueue(), TestJobQueue() (+3 more)

### Community 77 - "booking/src/lib/apiError.ts"
Cohesion: 0.29
Nodes (8): ApiError, ApiErrorPayload, buildApiError(), parseApiResponse(), API_ERROR_MESSAGES, FIELD_ERROR_MESSAGES, getApiErrorMessage(), translateFieldErrors()

### Community 78 - "TenantsPage.tsx"
Cohesion: 0.13
Nodes (23): createTenant(), deleteTenant(), fetchTenants(), setTenantActive(), updateTenant(), updateTenantSubscription(), AVAILABLE_MODULES, PLAN_PRESETS (+15 more)

### Community 79 - "rate_limit.go"
Cohesion: 0.16
Nodes (19): buildCompoundIdempotencyKey(), forwardedClientIP(), isSensitiveRateLimitedRoute(), isTrustedProxy(), newMemoryRateLimiter(), normalizedClientID(), persistentRateLimitConfig(), rateLimitRouteKey() (+11 more)

### Community 80 - "EventBroker"
Cohesion: 0.20
Nodes (9): UseCase, NewBookingHandler(), NewCheckoutHandler(), EventBroker, NewEventBroker(), RegisterMetricsBroker(), CheckoutHandler, Event (+1 more)

### Community 81 - "withIdempotency"
Cohesion: 0.15
Nodes (10): isHopByHopHeader(), NewIdempotencyStore(), TestIdempotencyBlocksInFlightConcurrentRequests(), TestIdempotencyDoesNotCacheServerFailures(), TestIdempotencyIgnoredOnGetRequests(), TestIdempotencyReplaysCompletedResponse(), withIdempotency(), idempotencyEntry (+2 more)

### Community 82 - "PortalPage"
Cohesion: 0.39
Nodes (8): portalCancelAppointment(), portalLookup(), portalRescheduleAppointment(), PortalPage(), handleCancel(), handleConfirmReschedule(), handleLookup(), refresh()

### Community 83 - "checkout/usecase.go"
Cohesion: 0.38
Nodes (3): CheckoutInput, CheckoutResult, ProductItem

### Community 86 - "PaymentsPage.tsx"
Cohesion: 0.22
Nodes (15): downloadAdminCsv(), fetchPayments(), markPaymentPaid(), refundPayment(), handleExport(), handleExport(), PAYMENT_METHOD_LABELS, PAYMENT_STATUS_LABELS (+7 more)

### Community 87 - "postgres/commerce_repository.go"
Cohesion: 0.19
Nodes (4): NewMemberPlanRepository(), scanMemberPlans(), MemberPlanRepository, MembershipRepository

### Community 90 - "Analytics"
Cohesion: 0.33
Nodes (4): DayMetric, ProfessionalMetric, ServiceMetric, Analytics

### Community 91 - "ServiceRepository"
Cohesion: 0.25
Nodes (4): scanServiceRow(), scanServiceRows(), ServiceRepository, serviceScanner

### Community 93 - "booking/src/shared/components/Button.tsx"
Cohesion: 0.33
Nodes (5): ButtonProps, ButtonSize, ButtonVariant, sizeClasses, variantClasses

### Community 95 - "body_limit.go"
Cohesion: 0.83
Nodes (3): resolveMaxRequestBodyBytes(), shouldLimitRequestBody(), withRequestBodyLimit()

### Community 144 - "CustomerRepository"
Cohesion: 0.29
Nodes (4): scanCustomerRow(), scanCustomerRows(), CustomerRepository, customerScanner

### Community 145 - "writeNotFoundError"
Cohesion: 0.31
Nodes (4): writeCustomerError(), writeNotificationError(), writeNotFoundError(), CustomerHandler

### Community 146 - ".UpdateAccountRole"
Cohesion: 0.39
Nodes (4): canManageTenantAccount(), isPrivilegedAdminRole(), isStaffRole(), normalizeTenantScope()

### Community 148 - "requireRoles"
Cohesion: 0.53
Nodes (5): hasRole(), requireOwnerOrAdmin(), requirePrivilegedAdmin(), requireRoles(), registerAuthRoutes()

## Knowledge Gaps
- **274 isolated node(s):** `graphify-mcp`, `name`, `private`, `version`, `type` (+269 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 517 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **59 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Container` connect `Container` to `auth_cookie.go`, `container.go`, `Payment`, `UseCase`, `UseCase`, `Tenant`, `UseCase`, `Service`, `Customer`, `Queue`, `Professional`, `Entry`, `database/sql.DB`, `Review`, `Connection`, `CashRegister`, `UseCase`, `UseCase`, `Resource`, `time.Duration`, `UseCase`, `repositorySet`, `UseCase`, `NewContainer`, `EventBroker`?**
  _High betweenness centrality (0.035) - this node is a cross-community bridge._
- **Why does `UseCase` connect `UseCase` to `auth_cookie.go`, `writeAPIError`, `AuthRepository`, `net/http.Request`, `Claims`, `RouterConfig`, `resolveRequestTenantID`, `auth/usecase.go`, `rate_limit.go`, `Account`, `withIdempotency`, `time.Duration`, `.UpdateAccountRole`, `Container`, `requireRoles`, `.registerWithRole`, `time.Time`?**
  _High betweenness centrality (0.029) - this node is a cross-community bridge._
- **Why does `repositorySet` connect `repositorySet` to `Payment`, `UseCase`, `UseCase`, `Tenant`, `Service`, `go_pkg_strings`, `Customer`, `Queue`, `Professional`, `Entry`, `database/sql.DB`, `Appointment`, `Review`, `Connection`, `MemberPlan`, `CashRegister`, `UseCase`, `Package`, `Resource`, `GiftCard`, `PackageCredit`, `Account`, `Container`, `memoryRepositories`, `UseCase`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Are the 128 inferred relationships involving `writeJSON()` (e.g. with `health()` and `readiness()`) actually correct?**
  _`writeJSON()` has 128 INFERRED edges - model-reasoned connections that need verification._
- **Are the 103 inferred relationships involving `resolveRequestTenantID()` (e.g. with `.Analytics()` and `.CheckIn()`) actually correct?**
  _`resolveRequestTenantID()` has 103 INFERRED edges - model-reasoned connections that need verification._
- **What connects `graphify-mcp`, `name`, `private` to the rest of the system?**
  _274 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `AppointmentsPage.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.08521303258145363 - nodes in this community are weakly interconnected._
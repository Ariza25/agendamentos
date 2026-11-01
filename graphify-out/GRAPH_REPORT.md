# Graph Report - erpagendamentos  (2026-11-01)

## Corpus Check
- 260 files · ~124,478 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 15 file(s) not represented in the graph (top: (none) 6, .conf 5, .css 3)

## Summary
- 2768 nodes · 9023 edges · 143 communities (85 shown, 58 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 661 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `0a5040c3`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- AppointmentsPage.tsx
- Payment
- Button
- adminApi.ts
- Professional
- adminAuth.ts
- auth/usecase.go
- writeAPIError
- ref_react
- container.go
- auth_cookie.go
- UseCase
- adminFetch
- UseCase
- UseCase
- Tenant
- PaymentsPage
- observability.go
- AdminRoutes.tsx
- Service
- net/http.ResponseWriter
- go_pkg_strings
- Customer
- resolveRequestTenantID
- Entry
- admin/package.json
- booking/package.json
- time.Duration
- Connection
- RouterConfig
- landing/package.json
- booking/src/pages/WaitingListPage.tsx
- Membership
- database/sql.DB
- Account
- UseCase
- CashRegister
- GiftCard
- PackageCredit
- Queue
- PortalPage.tsx
- bookingApi.ts
- landing/src/App.tsx
- Resource
- Package
- writeInvalidPayload
- Review
- MemberPlan
- writeJSON
- writePaymentError
- ConfigurationPage.tsx
- BookingPage.tsx
- ServiceRepository
- AppointmentHandler
- rate_limit.go
- UseCase
- AppointmentRepository
- memoryRepositories
- compilerOptions
- Appointment
- TenantsPage
- compilerOptions
- compilerOptions
- Repository
- Container
- Configuration
- sync.Mutex
- schema.sql
- EventBroker
- writeResourceError
- AppointmentRepository
- time.Time
- Transaction
- auth_context.go
- ERP Agendamentos
- TeamPage
- ical.go
- metrics.go
- writeCustomerError
- booking/src/lib/apiError.ts
- Claims
- ExportHandler
- PackagesPage.tsx
- Repository
- CashRegisterRepository
- PortalPage
- .UpdateAccountRole
- UseCase
- writeProfessionalError
- writeServiceError
- SecretCipher
- Analytics
- checkout/usecase.go
- Summary
- booking/src/shared/components/Button.tsx
- widget_handler.go
- booking/src/main.tsx
- graphify
- rules/graphify.md
- workflows/graphify.md
- generate-certs.sh
- erpagendamentos/backend

## God Nodes (most connected - your core abstractions)
1. `writeJSON()` - 132 edges
2. `resolveRequestTenantID()` - 108 edges
3. `adminFetch()` - 74 edges
4. `Button` - 72 edges
5. `writeInvalidPayload()` - 71 edges
6. `useToast()` - 68 edges
7. `UseCase` - 65 edges
8. `invalidateAdminGetCache()` - 64 edges
9. `writeInternalError()` - 58 edges
10. `RouterConfig` - 54 edges

## Surprising Connections (you probably didn't know these)
- `getPlanBadge()` --calls--> `Badge()`  [EXTRACTED]
  admin/src/pages/TenantsPage.tsx → admin/src/shared/components/Badge.tsx
- `getStatusBadge()` --calls--> `Badge()`  [EXTRACTED]
  admin/src/pages/TenantsPage.tsx → admin/src/shared/components/Badge.tsx
- `newResourceUseCase()` --calls--> `NewUseCase()`  [INFERRED]
  backend/internal/application/appointment/resource_test.go → backend/internal/application/appointment/usecase.go
- `newTestUseCase()` --calls--> `NewUseCase()`  [INFERRED]
  backend/internal/application/booking/usecase_test.go → backend/internal/application/booking/usecase.go
- `TestManagementCodeNormalization()` --calls--> `hashManagementCode()`  [INFERRED]
  backend/internal/application/booking/usecase_test.go → backend/internal/application/booking/usecase.go

## Import Cycles
- None detected.

## Communities (143 total, 58 thin omitted)

### Community 0 - "AppointmentsPage.tsx"
Cohesion: 0.09
Nodes (48): computeEndsAt(), handleSubmit(), AppointmentFormModalProps, STATUS_OPTIONS, checkInAppointment(), createAppointment(), deleteAppointment(), fetchAppointmentAnalytics() (+40 more)

### Community 1 - "Payment"
Cohesion: 0.05
Nodes (24): UseCase, NewUseCase(), normalizeStatus(), buildPaymentProvider(), IntentResult, Payment, NewMercadoPagoProvider(), NewMockProvider() (+16 more)

### Community 2 - "Button"
Cohesion: 0.12
Nodes (51): AppointmentFormModal(), CheckoutComandaModal(), createMemberPlan(), createMembership(), deleteMembership(), fetchMemberPlans(), fetchMemberships(), setMembershipActive() (+43 more)

### Community 3 - "adminApi.ts"
Cohesion: 0.05
Nodes (84): handleSubmit(), CheckoutComandaModalProps, SelectedProduct, AppointmentFilters, checkoutAppointment(), createCommissionPayout(), createCommissionRule(), createProduct() (+76 more)

### Community 4 - "Professional"
Cohesion: 0.10
Nodes (10): UseCase, NewUseCase(), Professional, scanProfessionalRow(), scanProfessionalRows(), NewProfessionalHandler(), ProfessionalRepository, ProfessionalRepository (+2 more)

### Community 5 - "adminAuth.ts"
Cohesion: 0.07
Nodes (48): App(), handleLogout(), AppProviders(), AppProvidersProps, ResetPasswordPage(), handleSubmit(), AppRouter(), ADMIN_ROLES (+40 more)

### Community 7 - "writeAPIError"
Cohesion: 0.11
Nodes (21): resolveAuthClaims(), browserSessionPayload(), exposeAuthTokenForTests(), AuthHandler, AuthHandler, readiness(), RecordRequestMetric(), logError() (+13 more)

### Community 8 - "ref_react"
Cohesion: 0.07
Nodes (72): deleteGiftCard(), deleteProfessional(), fetchGiftCards(), fetchNotifications(), fetchReviews(), fetchWaitingList(), setProfessionalActive(), updateWaitingListStatus() (+64 more)

### Community 9 - "container.go"
Cohesion: 0.11
Nodes (4): deploymentEnvironment(), isProductionEnv(), CategoryTotal, Type

### Community 10 - "auth_cookie.go"
Cohesion: 0.23
Nodes (20): bearerAuthAllowed(), authCookieDomain(), authCookiePartitioned(), authCookieSameSite(), authCookieSecure(), authTokenFromCookie(), clearAuthCookie(), csrfTokenFromCookie() (+12 more)

### Community 11 - "UseCase"
Cohesion: 0.06
Nodes (10): UseCase, CommissionRecord, CommissionRule, Payout, Repository, Summary, NewCommissionHandler(), PayoutInput (+2 more)

### Community 12 - "adminFetch"
Cohesion: 0.06
Nodes (70): closeCashRegister(), connectCalendar(), createAdminUser(), createCashRegisterEntry(), createCustomer(), createFinancialTransaction(), createGiftCard(), createPackage() (+62 more)

### Community 13 - "UseCase"
Cohesion: 0.07
Nodes (25): ChangePasswordInput, LoginInput, MFAStatus, PasswordResetCodeConfirmInput, PasswordResetConfirmInput, PasswordResetRequestInput, RegisterInput, SecretCipher (+17 more)

### Community 14 - "UseCase"
Cohesion: 0.07
Nodes (8): UseCase, Product, Repository, ServiceSupply, StockMovement, NewStockHandler(), StockRepository, StockRepository

### Community 15 - "Tenant"
Cohesion: 0.11
Nodes (12): UseCase, Tenant, UpdateSubscriptionInput, defaultModules(), isTenantIdentityConflict(), defaultModules(), scanTenantRow(), scanTenantRows() (+4 more)

### Community 16 - "PaymentsPage"
Cohesion: 0.17
Nodes (19): deleteCustomer(), downloadAdminCsv(), fetchCustomers(), fetchPayments(), markPaymentPaid(), refundPayment(), handleExport(), CustomersPage() (+11 more)

### Community 17 - "observability.go"
Cohesion: 0.09
Nodes (17): newLazyAppHandler(), main(), resolveDurationEnv(), resolvePositiveIntEnv(), buildNotificationProvider(), intEnv(), logFallback(), NewEmailProvider() (+9 more)

### Community 18 - "AdminRoutes.tsx"
Cohesion: 0.08
Nodes (41): AccessDeniedPage(), AdminRouteFallback(), AdminRoutes(), AdminRoutesProps, AgendaPage, AnalyticsPage, AppointmentsPage, CommissionsPage (+33 more)

### Community 19 - "Service"
Cohesion: 0.12
Nodes (6): UseCase, NewUseCase(), Service, NewServiceHandler(), ServiceRepository, Repository

### Community 20 - "net/http.ResponseWriter"
Cohesion: 0.15
Nodes (6): buildPasswordResetURL(), newAuthMFAHandler(), writeCommerceError(), HandleMetrics(), CommerceHandler, responseRecorder

### Community 21 - "go_pkg_strings"
Cohesion: 0.15
Nodes (5): Open(), resolveDurationEnv(), resolveIntEnv(), Method, Status

### Community 22 - "Customer"
Cohesion: 0.09
Nodes (14): UseCase, NewUseCase(), Customer, NewCustomerRepository(), onlyDigits(), scanCustomerRow(), scanCustomerRows(), NewCustomerHandler() (+6 more)

### Community 23 - "resolveRequestTenantID"
Cohesion: 0.12
Nodes (6): resolveRequestTenantID(), writeInternalError(), CalendarSyncHandler, CommissionHandler, FinanceHandler, StockHandler

### Community 24 - "Entry"
Cohesion: 0.11
Nodes (11): UseCase, NewUseCase(), Entry, NewWaitingListRepository(), scanWaitingRow(), scanWaitingRows(), NewWaitingListHandler(), WaitingListRepository (+3 more)

### Community 25 - "admin/package.json"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+24 more)

### Community 26 - "booking/package.json"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+24 more)

### Community 28 - "time.Duration"
Cohesion: 0.13
Nodes (10): MaintenanceWorker, MaintenanceConfig, NewMaintenanceWorker(), buildAuthComponents(), boolEnv(), durationSecondsEnv(), seedAdminAccount(), CleanupSummary (+2 more)

### Community 29 - "Connection"
Cohesion: 0.09
Nodes (10): UseCase, NewUseCase(), TestCalendarUseCase(), Connection, Repository, SyncEvent, NewCalendarRepository(), NewCalendarSyncHandler() (+2 more)

### Community 30 - "RouterConfig"
Cohesion: 0.25
Nodes (27): requireAdmin(), NewRouter(), startBackgroundWorkers(), registerAppointmentRoutes(), registerConfigurationRoutes(), registerCustomerRoutes(), registerProfessionalRoutes(), registerResourceRoutes() (+19 more)

### Community 31 - "landing/package.json"
Cohesion: 0.06
Nodes (30): dependencies, lucide-react, react, react-dom, devDependencies, tailwindcss, @tailwindcss/vite, @types/react (+22 more)

### Community 32 - "booking/src/pages/WaitingListPage.tsx"
Cohesion: 0.17
Nodes (22): App(), fetchPublicReview(), signUpWaitingList(), submitPublicReview(), DetailsStep(), ReviewPage(), handleSubmit(), SlugEntryPage() (+14 more)

### Community 33 - "Membership"
Cohesion: 0.12
Nodes (5): Membership, NewMembershipRepository(), MembershipRepository, MembershipRepository, MembershipRepository

### Community 34 - "database/sql.DB"
Cohesion: 0.14
Nodes (28): applyDatabaseSQL(), resolveSchemaFile(), initializeRepositories(), postgresRepositories(), NewAppointmentRepository(), NewAuthRepository(), NewCalendarRepository(), NewCashRegisterRepository() (+20 more)

### Community 35 - "Account"
Cohesion: 0.04
Nodes (12): Repository, Account, AuthState, DeviceSession, PasswordResetToken, UpdateProfileInput, cloneAuthAccount(), resetTokenExpired() (+4 more)

### Community 36 - "UseCase"
Cohesion: 0.18
Nodes (10): appointmentDurationMinutes(), formatAppointmentForTemplate(), generateManagementCode(), UseCase, hashManagementCode(), parseRFC3339(), BookingInput, PortalAppointment (+2 more)

### Community 37 - "CashRegister"
Cohesion: 0.14
Nodes (7): UseCase, CashRegister, CashRegisterEntry, Repository, ShiftSummary, NewCashRegisterHandler(), CashRegisterRepository

### Community 38 - "GiftCard"
Cohesion: 0.12
Nodes (5): GiftCard, NewGiftCardRepository(), GiftCardRepository, GiftCardRepository, GiftCardRepository

### Community 39 - "PackageCredit"
Cohesion: 0.13
Nodes (6): PackageCredit, NewPackageCreditRepository(), scanPackageCredits(), PackageCreditRepository, PackageCreditRepository, PackageCreditRepository

### Community 40 - "Queue"
Cohesion: 0.08
Nodes (9): Queue, NewQueue(), TestJobQueue(), Job, Repository, NewJobRepository(), HandlerFunc, JobRepository (+1 more)

### Community 41 - "PortalPage.tsx"
Cohesion: 0.19
Nodes (23): PortalAppointment, PortalResult, Slot, formatDateLong(), formatDateTime(), formatTime(), fromDateKey(), nextDateKeys() (+15 more)

### Community 42 - "bookingApi.ts"
Cohesion: 0.14
Nodes (23): backendOrigin, buildApiUrl(), AvailabilityResponse, cancelPublicAppointment(), fetchAvailability(), fetchPublicAppointment(), fetchPublicInfo(), fetchPublicProfessionals() (+15 more)

### Community 43 - "landing/src/App.tsx"
Cohesion: 0.16
Nodes (14): App(), CtaSection(), FEATURES, FeaturesSection(), Footer(), Header(), NAV_ITEMS, HeroSection() (+6 more)

### Community 44 - "Resource"
Cohesion: 0.11
Nodes (9): UseCase, NewUseCase(), Resource, scanResourceRow(), NewResourceHandler(), ResourceRepository, ResourceRepository, resourceScanner (+1 more)

### Community 45 - "Package"
Cohesion: 0.12
Nodes (5): Package, NewPackageRepository(), PackageRepository, PackageRepository, PackageRepository

### Community 46 - "writeInvalidPayload"
Cohesion: 0.10
Nodes (14): writeConfigurationError(), writeNotificationError(), writeBadRequestError(), writeConflictError(), writeInvalidPayload(), writeNotFoundError(), writeRawJSON(), writeReviewError() (+6 more)

### Community 47 - "Review"
Cohesion: 0.12
Nodes (12): UseCase, NewUseCase(), Review, NewReviewRepository(), scanReviewRow(), scanReviewRows(), NewAppointmentHandler(), NewReviewHandler() (+4 more)

### Community 48 - "MemberPlan"
Cohesion: 0.11
Nodes (8): NewUseCase(), MemberPlan, NewMemberPlanRepository(), nullIfEmpty(), scanMemberPlans(), MemberPlanRepository, MemberPlanRepository, MemberPlanRepository

### Community 49 - "writeJSON"
Cohesion: 0.21
Nodes (5): writeBookingError(), health(), writeJSON(), BookingHandler, CashRegisterHandler

### Community 51 - "ConfigurationPage.tsx"
Cohesion: 0.14
Nodes (11): fetchConfiguration(), updateConfiguration(), DAY_LABELS, ActiveTab, CHANNEL_OPTIONS, ConfigurationPage(), handleSubmit(), EMPTY_DAY (+3 more)

### Community 52 - "BookingPage.tsx"
Cohesion: 0.13
Nodes (15): BookingResult, createPublicBooking(), formatMoney(), BookingPage(), handleBookingSubmit(), CenteredLayout(), EmptyNote(), ErrorCard() (+7 more)

### Community 53 - "ServiceRepository"
Cohesion: 0.25
Nodes (4): scanServiceRow(), scanServiceRows(), ServiceRepository, serviceScanner

### Community 54 - "AppointmentHandler"
Cohesion: 0.21
Nodes (4): writeAppointmentError(), resolveActorTenantID(), writeAdminAuthError(), AppointmentHandler

### Community 55 - "rate_limit.go"
Cohesion: 0.30
Nodes (13): forwardedClientIP(), isSensitiveRateLimitedRoute(), isTrustedProxy(), newMemoryRateLimiter(), normalizedClientID(), persistentRateLimitConfig(), rateLimitRouteKey(), remoteClientIP() (+5 more)

### Community 56 - "UseCase"
Cohesion: 0.15
Nodes (3): generateGiftCardCode(), UseCase, NewCommerceHandler()

### Community 57 - "AppointmentRepository"
Cohesion: 0.15
Nodes (4): StatusCount, parseAppointmentBound(), parseAppointmentInterval(), AppointmentRepository

### Community 58 - "memoryRepositories"
Cohesion: 0.07
Nodes (54): newAppointment(), newResourceUseCase(), TestCreateAutoAssignsAvailableResource(), TestCreateRejectsBusyExplicitResource(), TestCreateRejectsUnknownResource(), TestUpdateAutoReassignsFreedResource(), newTestUseCase(), TestAvailableSlotsRespectsClosedDays() (+46 more)

### Community 59 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 60 - "Appointment"
Cohesion: 0.12
Nodes (7): RecurrenceInput, Repository, UseCase, NewUseCase(), parseRFC3339(), validStatusChange(), Appointment

### Community 61 - "TenantsPage"
Cohesion: 0.36
Nodes (9): deleteTenant(), fetchTenants(), setTenantActive(), TenantsPage(), getPlanBadge(), getStatusBadge(), handleDelete(), handleToggleActive() (+1 more)

### Community 62 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 63 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 64 - "Repository"
Cohesion: 0.15
Nodes (10): UseCase, NewUseCase(), formatAppointmentLocal(), NewReminderWorker(), NewCalendarHandler(), NewConfigurationHandler(), Repository, ReminderWorker (+2 more)

### Community 65 - "Container"
Cohesion: 0.10
Nodes (14): lazyAppHandler, writeLazyJSON(), NewUseCase(), UseCase, NewUseCase(), NewUseCase(), ensureDefaultTenant(), Container (+6 more)

### Community 66 - "Configuration"
Cohesion: 0.23
Nodes (10): workingHoursFor(), Configuration, cloneConfiguration(), defaultMemoryConfiguration(), normalizeConfigurationTenantID(), normalizePostgresTenantID(), ClosedDay, WorkingHours (+2 more)

### Community 67 - "sync.Mutex"
Cohesion: 0.12
Nodes (12): UseCase, NewUseCase(), Log, NewNotificationRepository(), NewNotificationHandler(), NotificationRepository, Provider, Repository (+4 more)

### Community 68 - "schema.sql"
Cohesion: 0.06
Nodes (70): catalog.appointments, catalog.auth_rate_limits, catalog.auth_sessions, catalog.calendar_connections, catalog.cash_register_entries, catalog.cash_registers, catalog.commission_payouts, catalog.commission_records (+62 more)

### Community 69 - "EventBroker"
Cohesion: 0.27
Nodes (5): NewBookingHandler(), EventBroker, NewEventBroker(), Event, MetricsRegistry

### Community 71 - "AppointmentRepository"
Cohesion: 0.21
Nodes (5): appointmentSelectColumns(), scanAppointmentRow(), scanAppointmentRows(), AppointmentRepository, appointmentScanner

### Community 72 - "time.Time"
Cohesion: 0.33
Nodes (8): countOverlappingSlots(), isClosedDay(), localMidnight(), overlapsAnyBusy(), parseClock(), ProfessionalSlots, Slot, rateLimitEntry

### Community 73 - "Transaction"
Cohesion: 0.24
Nodes (3): Transaction, FinanceRepository, FinanceRepository

### Community 74 - "auth_context.go"
Cohesion: 0.11
Nodes (21): hasRole(), isAdminRole(), isPrivilegedAdminRole(), requestWithAuthClaims(), requireMasterAdmin(), requireOwnerOrAdmin(), requirePrivilegedAdmin(), requireRoles() (+13 more)

### Community 75 - "ERP Agendamentos"
Cohesion: 0.14
Nodes (13): Banco de Dados & Schema, CI, Endpoints principais, ERP Agendamentos, Estrutura, Pré-requisitos, Rodando com Docker (produção local), Rodando em desenvolvimento (sem banco) (+5 more)

### Community 76 - "TeamPage"
Cohesion: 0.52
Nodes (7): deleteAdminUser(), fetchAdminUsers(), updateAdminUserRole(), TeamPage(), handleDelete(), handleRoleChange(), loadUsers()

### Community 77 - "ical.go"
Cohesion: 0.83
Nodes (3): BuildICalCalendar(), formatUTC(), sanitizeLine()

### Community 80 - "booking/src/lib/apiError.ts"
Cohesion: 0.29
Nodes (8): ApiError, ApiErrorPayload, buildApiError(), parseApiResponse(), API_ERROR_MESSAGES, FIELD_ERROR_MESSAGES, getApiErrorMessage(), translateFieldErrors()

### Community 82 - "Claims"
Cohesion: 0.31
Nodes (5): isProductionAuthEnvironment(), Claims, decodeSegment(), encodeSegment(), JWTManager

### Community 84 - "PackagesPage.tsx"
Cohesion: 0.24
Nodes (16): deletePackage(), deleteResource(), fetchPackageCredits(), fetchPackages(), fetchResources(), setPackageActive(), setResourceActive(), PackagesPage() (+8 more)

### Community 85 - "Repository"
Cohesion: 0.24
Nodes (5): UseCase, NewUseCase(), Repository, NewCheckoutHandler(), CheckoutHandler

### Community 87 - "PortalPage"
Cohesion: 0.39
Nodes (8): portalCancelAppointment(), portalLookup(), portalRescheduleAppointment(), PortalPage(), handleCancel(), handleConfirmReschedule(), handleLookup(), refresh()

### Community 89 - ".UpdateAccountRole"
Cohesion: 0.39
Nodes (4): canManageTenantAccount(), isPrivilegedAdminRole(), isStaffRole(), normalizeTenantScope()

### Community 90 - "UseCase"
Cohesion: 0.25
Nodes (3): UseCase, NewUseCase(), NewFinanceHandler()

### Community 97 - "Analytics"
Cohesion: 0.33
Nodes (4): DayMetric, ProfessionalMetric, ServiceMetric, Analytics

### Community 98 - "checkout/usecase.go"
Cohesion: 0.38
Nodes (3): CheckoutInput, CheckoutResult, ProductItem

### Community 104 - "booking/src/shared/components/Button.tsx"
Cohesion: 0.33
Nodes (5): ButtonProps, ButtonSize, ButtonVariant, sizeClasses, variantClasses

### Community 105 - "widget_handler.go"
Cohesion: 0.67
Nodes (3): NewWidgetHandler(), resolveWidgetBookingOrigin(), WidgetHandler

## Knowledge Gaps
- **266 isolated node(s):** `graphify-mcp`, `name`, `private`, `version`, `type` (+261 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 506 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **58 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Container` connect `Container` to `Payment`, `Professional`, `container.go`, `UseCase`, `UseCase`, `UseCase`, `Tenant`, `Service`, `Customer`, `Entry`, `time.Duration`, `Connection`, `database/sql.DB`, `UseCase`, `CashRegister`, `Queue`, `Resource`, `Review`, `UseCase`, `Appointment`, `Repository`, `sync.Mutex`, `EventBroker`, `auth_context.go`, `Repository`, `UseCase`?**
  _High betweenness centrality (0.038) - this node is a cross-community bridge._
- **Why does `UseCase` connect `UseCase` to `Container`, `Account`, `auth/usecase.go`, `writeAPIError`, `time.Time`, `auth_context.go`, `Claims`, `net/http.ResponseWriter`, `rate_limit.go`, `AppointmentHandler`, `resolveRequestTenantID`, `.UpdateAccountRole`, `time.Duration`, `RouterConfig`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Why does `RouterConfig` connect `RouterConfig` to `Payment`, `Professional`, `container.go`, `UseCase`, `UseCase`, `UseCase`, `Tenant`, `Service`, `Customer`, `Entry`, `time.Duration`, `Connection`, `UseCase`, `CashRegister`, `Resource`, `Review`, `UseCase`, `Appointment`, `Repository`, `Container`, `sync.Mutex`, `EventBroker`, `auth_context.go`, `Repository`, `UseCase`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Are the 128 inferred relationships involving `writeJSON()` (e.g. with `health()` and `readiness()`) actually correct?**
  _`writeJSON()` has 128 INFERRED edges - model-reasoned connections that need verification._
- **Are the 103 inferred relationships involving `resolveRequestTenantID()` (e.g. with `.Analytics()` and `.CheckIn()`) actually correct?**
  _`resolveRequestTenantID()` has 103 INFERRED edges - model-reasoned connections that need verification._
- **What connects `graphify-mcp`, `name`, `private` to the rest of the system?**
  _266 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `AppointmentsPage.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.08766233766233766 - nodes in this community are weakly interconnected._
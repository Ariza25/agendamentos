# Graph Report - erpagendamentos  (2026-11-02)

## Corpus Check
- 280 files · ~132,775 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 16 file(s) not represented in the graph (top: (none) 7, .conf 5, .css 3)

## Summary
- 2912 nodes · 9498 edges · 150 communities (90 shown, 60 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 683 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `51cf66d6`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- adminApi.ts
- Payment
- schema.sql
- Button
- container.go
- adminAuth.ts
- AppointmentsPage.tsx
- Entry
- net/http.ResponseWriter
- Tenant
- TenantsPage.tsx
- CommissionRule
- auth/usecase.go
- PackagesPage
- UseCase
- FinancePage.tsx
- net/http.Request
- writeJSON
- UseCase
- adminFetch
- Customer
- Professional
- UseCase
- CashRegister
- auth_cookie.go
- landing/src/App.tsx
- admin/package.json
- UseCase
- Service
- RouterConfig
- booking/package.json
- Review
- Appointment
- database/sql.DB
- landing/package.json
- booking/src/pages/WaitingListPage.tsx
- EventBroker
- Repository
- UseCase
- TenantsPage
- Account
- time.Time
- PackageCredit
- PortalPage.tsx
- AdminRoutes.tsx
- MemberPlan
- UseCase
- Package
- bookingApi.ts
- GiftCard
- writeInvalidPayload
- AuditLog
- .registerWithRole
- repositorySet
- Resource
- MembershipsPage.tsx
- BookingPage.tsx
- AppointmentRepository
- AuthRepository
- templates.go
- compilerOptions
- UseCase
- testing.T
- NewContainer
- compilerOptions
- compilerOptions
- PaymentsPage.tsx
- Connection
- AppointmentHandler
- ERP Agendamentos
- time.Duration
- Repository
- requireRoles
- CommissionsPage.tsx
- Claims
- .GetByID
- observability.go
- UseCase
- booking/src/lib/apiError.ts
- ProfessionalRepository
- UseCase
- Summary
- Queue
- PortalPage
- context.Context
- .UpdateAccountRole
- Container
- SecretCipher
- writeNotFoundError
- writeResourceError
- Analytics
- bcrypt_hasher.go
- writeServiceError
- booking/src/shared/components/Button.tsx
- ExportHandler
- widget_handler.go
- body_limit.go
- booking/src/main.tsx
- graphify
- rules/graphify.md
- workflows/graphify.md
- generate-certs.sh
- backup_postgres.sh
- certbot_ssl.sh
- erpagendamentos/backend

## God Nodes (most connected - your core abstractions)
1. `writeJSON()` - 134 edges
2. `resolveRequestTenantID()` - 108 edges
3. `adminFetch()` - 77 edges
4. `Button` - 72 edges
5. `writeInvalidPayload()` - 71 edges
6. `useToast()` - 68 edges
7. `UseCase` - 67 edges
8. `invalidateAdminGetCache()` - 64 edges
9. `writeInternalError()` - 58 edges
10. `RouterConfig` - 56 edges

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

## Communities (150 total, 60 thin omitted)

### Community 0 - "adminApi.ts"
Cohesion: 0.06
Nodes (68): CheckoutComandaModalProps, SelectedProduct, AppointmentFilters, createProduct(), createStockMovement(), deleteProduct(), fetchProducts(), fetchServices() (+60 more)

### Community 1 - "Payment"
Cohesion: 0.05
Nodes (22): UseCase, NewUseCase(), normalizeStatus(), buildPaymentProvider(), IntentResult, Payment, NewMercadoPagoProvider(), NewMockProvider() (+14 more)

### Community 2 - "schema.sql"
Cohesion: 0.05
Nodes (75): catalog.appointments, catalog.audit_logs, catalog.auth_rate_limits, catalog.auth_sessions, catalog.calendar_connections, catalog.cash_register_entries, catalog.cash_registers, catalog.commission_payouts (+67 more)

### Community 3 - "Button"
Cohesion: 0.13
Nodes (39): fetchConfiguration(), updateConfiguration(), DAY_LABELS, CommissionPayoutModal(), ActiveTab, CHANNEL_OPTIONS, ConfigurationPage(), handleSubmit() (+31 more)

### Community 4 - "container.go"
Cohesion: 0.08
Nodes (8): deploymentEnvironment(), isProductionEnv(), CheckoutInput, ProductItem, CategoryTotal, Type, Method, Status

### Community 5 - "adminAuth.ts"
Cohesion: 0.07
Nodes (53): App(), handleLogout(), AppProviders(), AppProvidersProps, ResetPasswordPage(), handleSubmit(), AppRouter(), ADMIN_ROLES (+45 more)

### Community 6 - "AppointmentsPage.tsx"
Cohesion: 0.09
Nodes (46): AppointmentFormModal(), computeEndsAt(), handleSubmit(), AppointmentFormModalProps, STATUS_OPTIONS, CheckoutComandaModal(), handleSubmit(), checkInAppointment() (+38 more)

### Community 7 - "Entry"
Cohesion: 0.11
Nodes (11): UseCase, NewUseCase(), Entry, NewWaitingListRepository(), scanWaitingRow(), scanWaitingRows(), NewWaitingListHandler(), WaitingListRepository (+3 more)

### Community 8 - "net/http.ResponseWriter"
Cohesion: 0.09
Nodes (13): isAdminRole(), resolveRequestTenantID(), writeConfigurationError(), writePaymentError(), writeInternalError(), writeWaitingListError(), CalendarSyncHandler, CommissionHandler (+5 more)

### Community 9 - "Tenant"
Cohesion: 0.07
Nodes (20): UseCase, NewUseCase(), UseCase, NewUseCase(), ensureDefaultTenant(), Tenant, UpdateSubscriptionInput, defaultModules() (+12 more)

### Community 10 - "TenantsPage.tsx"
Cohesion: 0.08
Nodes (69): deleteGiftCard(), fetchAppointmentAnalytics(), fetchGiftCards(), fetchNotifications(), fetchReviews(), fetchWaitingList(), updateWaitingListStatus(), fromDateKey() (+61 more)

### Community 11 - "CommissionRule"
Cohesion: 0.13
Nodes (4): CommissionRule, Payout, CommissionRepository, CommissionRepository

### Community 12 - "auth/usecase.go"
Cohesion: 0.08
Nodes (7): newMemoryRateLimiter(), persistentRateLimitConfig(), resolvePositiveEnvInt(), TestPersistentRateLimitConfigMatch(), authClaimsContextKey, persistentRateLimit, tokenPayload

### Community 13 - "PackagesPage"
Cohesion: 0.46
Nodes (8): deletePackage(), fetchPackageCredits(), fetchPackages(), setPackageActive(), PackagesPage(), handleDelete(), handleToggleActive(), loadPackages()

### Community 14 - "UseCase"
Cohesion: 0.07
Nodes (8): UseCase, Product, Repository, ServiceSupply, StockMovement, NewStockHandler(), StockRepository, StockRepository

### Community 15 - "FinancePage.tsx"
Cohesion: 0.15
Nodes (27): closeCashRegister(), createCashRegisterEntry(), createFinancialTransaction(), deleteFinancialTransaction(), fetchCashRegisterHistory(), fetchCurrentCashRegister(), fetchFinancialSummary(), fetchFinancialTransactions() (+19 more)

### Community 16 - "net/http.Request"
Cohesion: 0.07
Nodes (40): requestWithAuthClaims(), requireMasterAdmin(), resolveAuthClaims(), resolveTenantID(), withProtectedNamespaces(), browserSessionPayload(), buildPasswordResetURL(), exposeAuthTokenForTests() (+32 more)

### Community 17 - "writeJSON"
Cohesion: 0.14
Nodes (6): writeBookingError(), writeJSON(), writeTenantError(), BookingHandler, CashRegisterHandler, TenantHandler

### Community 19 - "UseCase"
Cohesion: 0.17
Nodes (11): appointmentDurationMinutes(), formatAppointmentForTemplate(), generateManagementCode(), UseCase, hashManagementCode(), parseRFC3339(), BookingInput, BookingResult (+3 more)

### Community 20 - "adminFetch"
Cohesion: 0.06
Nodes (63): connectCalendar(), createAdminUser(), createCustomer(), createGiftCard(), createMemberPlan(), createMembership(), createPackage(), createPackageCredit() (+55 more)

### Community 21 - "Customer"
Cohesion: 0.10
Nodes (11): UseCase, NewUseCase(), Customer, onlyDigits(), scanCustomerRow(), scanCustomerRows(), NewCustomerHandler(), Repository (+3 more)

### Community 23 - "Professional"
Cohesion: 0.13
Nodes (7): UseCase, NewUseCase(), Professional, NewAppointmentHandler(), NewProfessionalHandler(), ProfessionalRepository, Repository

### Community 24 - "UseCase"
Cohesion: 0.10
Nodes (14): LoginInput, MFAStatus, SecretCipher, Session, TokenManager, generateTOTPSecret(), totpCode(), totpURI() (+6 more)

### Community 26 - "CashRegister"
Cohesion: 0.11
Nodes (8): UseCase, CashRegister, CashRegisterEntry, Repository, ShiftSummary, NewCashRegisterHandler(), CashRegisterRepository, CashRegisterRepository

### Community 27 - "auth_cookie.go"
Cohesion: 0.06
Nodes (38): lazyAppHandler, writeLazyJSON(), bearerAuthAllowed(), isPrivilegedAdminRole(), authCookieDomain(), authCookiePartitioned(), authCookieSameSite(), authCookieSecure() (+30 more)

### Community 28 - "landing/src/App.tsx"
Cohesion: 0.11
Nodes (21): App(), CtaSection(), FEATURES, FeaturesSection(), Footer(), Header(), NAV_ITEMS, HeroSection() (+13 more)

### Community 29 - "admin/package.json"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+24 more)

### Community 30 - "UseCase"
Cohesion: 0.09
Nodes (9): UseCase, Repository, Summary, Transaction, NewFinanceHandler(), NewMercadoPagoWebhookHandler(), FinanceRepository, FinanceRepository (+1 more)

### Community 31 - "Service"
Cohesion: 0.09
Nodes (10): UseCase, NewUseCase(), Service, scanServiceRow(), scanServiceRows(), NewServiceHandler(), ServiceRepository, ServiceRepository (+2 more)

### Community 32 - "RouterConfig"
Cohesion: 0.26
Nodes (26): requireAdmin(), NewRouter(), registerAppointmentRoutes(), registerConfigurationRoutes(), registerCustomerRoutes(), registerProfessionalRoutes(), registerResourceRoutes(), registerServiceRoutes() (+18 more)

### Community 33 - "booking/package.json"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+24 more)

### Community 34 - "Review"
Cohesion: 0.12
Nodes (11): UseCase, NewUseCase(), Review, NewReviewRepository(), scanReviewRow(), scanReviewRows(), NewReviewHandler(), ReviewRepository (+3 more)

### Community 35 - "Appointment"
Cohesion: 0.11
Nodes (6): Repository, Appointment, StatusCount, parseAppointmentBound(), parseAppointmentInterval(), AppointmentRepository

### Community 36 - "database/sql.DB"
Cohesion: 0.10
Nodes (34): applyDatabaseSQL(), resolveSchemaFile(), initializeRepositories(), postgresRepositories(), NewAppointmentRepository(), NewAuditRepository(), NewAuthRepository(), NewCalendarRepository() (+26 more)

### Community 37 - "landing/package.json"
Cohesion: 0.06
Nodes (30): dependencies, lucide-react, react, react-dom, devDependencies, tailwindcss, @tailwindcss/vite, @types/react (+22 more)

### Community 38 - "booking/src/pages/WaitingListPage.tsx"
Cohesion: 0.17
Nodes (22): App(), fetchPublicReview(), signUpWaitingList(), submitPublicReview(), DetailsStep(), ReviewPage(), handleSubmit(), SlugEntryPage() (+14 more)

### Community 39 - "EventBroker"
Cohesion: 0.27
Nodes (5): NewBookingHandler(), EventBroker, RegisterMetricsBroker(), Event, MetricsRegistry

### Community 40 - "Repository"
Cohesion: 0.09
Nodes (3): Repository, DeviceSession, sessionExpired()

### Community 41 - "UseCase"
Cohesion: 0.08
Nodes (12): UseCase, NewUseCase(), TestOutboxPublishAndDispatch(), NewOutboxRepository(), quoteLiteral(), NewSystemHandler(), OutboxRepository, EventDispatcher (+4 more)

### Community 42 - "TenantsPage"
Cohesion: 0.20
Nodes (17): checkSystemReadiness(), deleteTenant(), fetchAuditLogs(), fetchSystemStatus(), fetchTenants(), setTenantActive(), parseJson(), TenantsPage() (+9 more)

### Community 43 - "Account"
Cohesion: 0.12
Nodes (7): Account, PasswordResetToken, UpdateProfileInput, cloneAuthAccount(), resetTokenExpired(), AuthRepository, rateLimitEntry

### Community 44 - "time.Time"
Cohesion: 0.13
Nodes (20): countOverlappingSlots(), isClosedDay(), localMidnight(), overlapsAnyBusy(), parseClock(), workingHoursFor(), Configuration, BuildICalCalendar() (+12 more)

### Community 45 - "PackageCredit"
Cohesion: 0.13
Nodes (5): PackageCredit, scanPackageCredits(), PackageCreditRepository, PackageCreditRepository, PackageCreditRepository

### Community 46 - "PortalPage.tsx"
Cohesion: 0.19
Nodes (23): PortalAppointment, PortalResult, Slot, formatDateLong(), formatDateTime(), formatTime(), fromDateKey(), nextDateKeys() (+15 more)

### Community 47 - "AdminRoutes.tsx"
Cohesion: 0.09
Nodes (39): AccessDeniedPage(), AdminRouteFallback(), AdminRoutes(), AdminRoutesProps, AgendaPage, AnalyticsPage, AppointmentsPage, CommissionsPage (+31 more)

### Community 48 - "MemberPlan"
Cohesion: 0.16
Nodes (5): NewUseCase(), MemberPlan, scanMemberPlans(), MemberPlanRepository, MemberPlanRepository

### Community 49 - "UseCase"
Cohesion: 0.11
Nodes (5): UseCase, Membership, NewCommerceHandler(), MembershipRepository, MembershipRepository

### Community 50 - "Package"
Cohesion: 0.12
Nodes (4): Package, PackageRepository, PackageRepository, PackageRepository

### Community 51 - "bookingApi.ts"
Cohesion: 0.14
Nodes (23): backendOrigin, buildApiUrl(), AvailabilityResponse, cancelPublicAppointment(), fetchAvailability(), fetchPublicAppointment(), fetchPublicInfo(), fetchPublicProfessionals() (+15 more)

### Community 52 - "GiftCard"
Cohesion: 0.11
Nodes (5): generateGiftCardCode(), GiftCard, GiftCardRepository, GiftCardRepository, GiftCardRepository

### Community 53 - "writeInvalidPayload"
Cohesion: 0.29
Nodes (3): writeCommerceError(), writeInvalidPayload(), CommerceHandler

### Community 54 - "AuditLog"
Cohesion: 0.17
Nodes (11): AuditLog, Repository, UseCase, NewUseCase(), TestAuditLogRecordingAndListing(), NewAuditRepository(), NewAuditHandler(), TestAuditAndSystemHandlers() (+3 more)

### Community 55 - ".registerWithRole"
Cohesion: 0.15
Nodes (12): ChangePasswordInput, PasswordResetCodeConfirmInput, PasswordResetConfirmInput, PasswordResetRequestInput, RegisterInput, ValidationError, hashResetToken(), NewValidationError() (+4 more)

### Community 56 - "repositorySet"
Cohesion: 0.13
Nodes (11): UseCase, NewUseCase(), Log, NewNotificationRepository(), NewNotificationHandler(), repositorySet, NotificationRepository, Provider (+3 more)

### Community 57 - "Resource"
Cohesion: 0.10
Nodes (10): NewUseCase(), UseCase, Resource, scanResourceRow(), NewExportHandler(), NewResourceHandler(), ResourceRepository, ResourceRepository (+2 more)

### Community 58 - "MembershipsPage.tsx"
Cohesion: 0.33
Nodes (12): deleteMembership(), fetchMemberPlans(), fetchMemberships(), setMembershipActive(), updateMemberPlanStatus(), MembershipsPage(), handleDelete(), handlePlanStatusChange() (+4 more)

### Community 59 - "BookingPage.tsx"
Cohesion: 0.13
Nodes (15): BookingResult, createPublicBooking(), formatMoney(), BookingPage(), handleBookingSubmit(), CenteredLayout(), EmptyNote(), ErrorCard() (+7 more)

### Community 60 - "AppointmentRepository"
Cohesion: 0.19
Nodes (5): appointmentSelectColumns(), scanAppointmentRow(), scanAppointmentRows(), AppointmentRepository, appointmentScanner

### Community 61 - "AuthRepository"
Cohesion: 0.11
Nodes (3): AuthState, updateAuthPasswordAndVersion(), AuthRepository

### Community 62 - "templates.go"
Cohesion: 0.27
Nodes (16): greeting(), PublicBookingOrigin(), RenderCancellation(), RenderConfirmation(), RenderReminder(), RenderReschedule(), sampleData(), TestGreetingWithoutName() (+8 more)

### Community 63 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 64 - "UseCase"
Cohesion: 0.18
Nodes (4): RecurrenceInput, UseCase, parseRFC3339(), validStatusChange()

### Community 65 - "testing.T"
Cohesion: 0.15
Nodes (18): newAppointment(), newResourceUseCase(), TestCreateAutoAssignsAvailableResource(), TestCreateRejectsBusyExplicitResource(), TestCreateRejectsUnknownResource(), TestUpdateAutoReassignsFreedResource(), newTestUseCase(), TestAvailableSlotsRespectsClosedDays() (+10 more)

### Community 66 - "NewContainer"
Cohesion: 0.08
Nodes (36): NewUseCase(), TestCashRegisterFlow(), NewUseCase(), TestCheckoutFlow(), NewUseCase(), TestCommissionUseCase(), NewUseCase(), TestFinanceUseCase() (+28 more)

### Community 67 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 68 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 69 - "PaymentsPage.tsx"
Cohesion: 0.16
Nodes (21): deleteCustomer(), downloadAdminCsv(), fetchCustomers(), fetchPayments(), markPaymentPaid(), refundPayment(), handleExport(), CustomersPage() (+13 more)

### Community 70 - "Connection"
Cohesion: 0.09
Nodes (10): UseCase, NewUseCase(), TestCalendarUseCase(), Connection, Repository, SyncEvent, NewCalendarRepository(), NewCalendarSyncHandler() (+2 more)

### Community 71 - "AppointmentHandler"
Cohesion: 0.11
Nodes (10): writeAppointmentError(), resolveActorTenantID(), writeAdminAuthError(), writeBadRequestError(), writeConflictError(), writeRawJSON(), writeReviewError(), AppointmentHandler (+2 more)

### Community 72 - "ERP Agendamentos"
Cohesion: 0.14
Nodes (13): Banco de Dados & Schema, CI, Endpoints principais, ERP Agendamentos, Estrutura, Pré-requisitos, Rodando com Docker (produção local), Rodando em desenvolvimento (sem banco) (+5 more)

### Community 73 - "time.Duration"
Cohesion: 0.15
Nodes (9): buildAuthComponents(), boolEnv(), durationSecondsEnv(), intEnv(), CleanupSummary, NewJWTManager(), NewSecretCipher(), memoryRateLimiter (+1 more)

### Community 75 - "requireRoles"
Cohesion: 0.43
Nodes (6): hasRole(), requireOwnerOrAdmin(), requirePrivilegedAdmin(), requireRoles(), NewAuthHandler(), registerAuthRoutes()

### Community 76 - "CommissionsPage.tsx"
Cohesion: 0.16
Nodes (24): createCommissionPayout(), createCommissionRule(), deleteCommissionRule(), deleteProfessional(), fetchCommissionPayouts(), fetchCommissionRecords(), fetchCommissionRules(), fetchCommissionSummaries() (+16 more)

### Community 77 - "Claims"
Cohesion: 0.31
Nodes (5): isProductionAuthEnvironment(), Claims, decodeSegment(), encodeSegment(), JWTManager

### Community 78 - ".GetByID"
Cohesion: 0.21
Nodes (3): nullIfEmpty(), MemberPlanRepository, MembershipRepository

### Community 79 - "observability.go"
Cohesion: 0.10
Nodes (16): newLazyAppHandler(), main(), resolveDurationEnv(), resolvePositiveIntEnv(), buildNotificationProvider(), logFallback(), NewEmailProvider(), NewLogProvider() (+8 more)

### Community 80 - "UseCase"
Cohesion: 0.21
Nodes (7): UseCase, NewCheckoutHandler(), AuditLogger, CheckoutResult, OutboxPublisher, TxRunner, CheckoutHandler

### Community 81 - "booking/src/lib/apiError.ts"
Cohesion: 0.29
Nodes (8): ApiError, ApiErrorPayload, buildApiError(), parseApiResponse(), API_ERROR_MESSAGES, FIELD_ERROR_MESSAGES, getApiErrorMessage(), translateFieldErrors()

### Community 82 - "ProfessionalRepository"
Cohesion: 0.29
Nodes (4): scanProfessionalRow(), scanProfessionalRows(), ProfessionalRepository, professionalScanner

### Community 83 - "UseCase"
Cohesion: 0.22
Nodes (3): UseCase, NewCommissionHandler(), PayoutInput

### Community 87 - "Queue"
Cohesion: 0.08
Nodes (6): Queue, Job, Repository, HandlerFunc, JobRepository, JobRepository

### Community 88 - "PortalPage"
Cohesion: 0.39
Nodes (8): portalCancelAppointment(), portalLookup(), portalRescheduleAppointment(), PortalPage(), handleCancel(), handleConfirmReschedule(), handleLookup(), refresh()

### Community 89 - "context.Context"
Cohesion: 0.20
Nodes (7): MaintenanceWorker, MaintenanceConfig, NewMaintenanceWorker(), startBackgroundWorkers(), authComponents, NopTxRunner, TxRunner

### Community 90 - ".UpdateAccountRole"
Cohesion: 0.33
Nodes (4): canManageTenantAccount(), isPrivilegedAdminRole(), isStaffRole(), normalizeTenantScope()

### Community 92 - "Container"
Cohesion: 0.12
Nodes (12): NewUseCase(), UseCase, NewUseCase(), formatAppointmentLocal(), NewReminderWorker(), Container, NewCalendarHandler(), NewConfigurationHandler() (+4 more)

### Community 95 - "writeNotFoundError"
Cohesion: 0.18
Nodes (6): writeCustomerError(), writeNotificationError(), writeProfessionalError(), writeNotFoundError(), CustomerHandler, ProfessionalHandler

### Community 99 - "Analytics"
Cohesion: 0.33
Nodes (4): DayMetric, ProfessionalMetric, ServiceMetric, Analytics

### Community 102 - "booking/src/shared/components/Button.tsx"
Cohesion: 0.33
Nodes (5): ButtonProps, ButtonSize, ButtonVariant, sizeClasses, variantClasses

### Community 104 - "widget_handler.go"
Cohesion: 0.60
Nodes (3): NewWidgetHandler(), resolveWidgetBookingOrigin(), WidgetHandler

### Community 105 - "body_limit.go"
Cohesion: 0.83
Nodes (3): resolveMaxRequestBodyBytes(), shouldLimitRequestBody(), withRequestBodyLimit()

## Knowledge Gaps
- **274 isolated node(s):** `graphify-mcp`, `name`, `private`, `version`, `type` (+269 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 533 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **60 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Container` connect `Container` to `Payment`, `container.go`, `Entry`, `Tenant`, `UseCase`, `UseCase`, `Customer`, `Professional`, `UseCase`, `CashRegister`, `auth_cookie.go`, `UseCase`, `Service`, `Review`, `database/sql.DB`, `EventBroker`, `UseCase`, `UseCase`, `AuditLog`, `repositorySet`, `Resource`, `UseCase`, `NewContainer`, `Connection`, `time.Duration`, `UseCase`, `UseCase`, `Queue`, `context.Context`?**
  _High betweenness centrality (0.043) - this node is a cross-community bridge._
- **Why does `UseCase` connect `UseCase` to `RouterConfig`, `AppointmentHandler`, `Repository`, `time.Duration`, `Tenant`, `Account`, `auth/usecase.go`, `Claims`, `time.Time`, `requireRoles`, `net/http.Request`, `net/http.ResponseWriter`, `.registerWithRole`, `context.Context`, `.UpdateAccountRole`, `auth_cookie.go`, `Container`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **Why does `RouterConfig` connect `RouterConfig` to `Payment`, `container.go`, `Entry`, `Tenant`, `UseCase`, `UseCase`, `Customer`, `Professional`, `UseCase`, `CashRegister`, `UseCase`, `Service`, `Review`, `EventBroker`, `UseCase`, `UseCase`, `AuditLog`, `repositorySet`, `Resource`, `UseCase`, `Connection`, `time.Duration`, `requireRoles`, `UseCase`, `UseCase`, `context.Context`, `Container`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **Are the 130 inferred relationships involving `writeJSON()` (e.g. with `health()` and `readiness()`) actually correct?**
  _`writeJSON()` has 130 INFERRED edges - model-reasoned connections that need verification._
- **Are the 103 inferred relationships involving `resolveRequestTenantID()` (e.g. with `.Analytics()` and `.CheckIn()`) actually correct?**
  _`resolveRequestTenantID()` has 103 INFERRED edges - model-reasoned connections that need verification._
- **What connects `graphify-mcp`, `name`, `private` to the rest of the system?**
  _274 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `adminApi.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.058738738738738736 - nodes in this community are weakly interconnected._
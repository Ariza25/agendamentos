# Graph Report - erpagendamentos  (2026-11-01)

## Corpus Check
- 270 files · ~119,712 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 15 file(s) not represented in the graph (top: (none) 6, .conf 5, .css 3)

## Summary
- 2743 nodes · 8863 edges · 155 communities (92 shown, 63 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 651 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `340fbbea`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- adminAuth.ts
- auth_cookie.go
- Payment
- adminApi.ts
- ref_react
- go_pkg_errors
- Button
- UseCase
- PageContainer
- email_provider.go
- UseCase
- adminFetch
- container.go
- Service
- net/http.Request
- Tenant
- UseCase
- Customer
- UseCase
- Professional
- net/http.ResponseWriter
- Entry
- admin/package.json
- database/sql.DB
- AppointmentHandler
- booking/package.json
- TeamPage.tsx
- Appointment
- UseCase
- Account
- Connection
- Job
- rate_limit.go
- RouterConfig
- landing/package.json
- .Add
- booking/src/pages/WaitingListPage.tsx
- MemberPlan
- CashRegister
- Repository
- UseCase
- Configuration
- GiftCard
- PortalPage.tsx
- AdminRoutes.tsx
- PackageCredit
- Package
- bookingApi.ts
- landing/src/App.tsx
- Resource
- Container
- 002_features.sql
- writeJSON
- writeInvalidPayload
- writeNotFoundError
- 001_baseline.sql
- BookingPage.tsx
- UseCase
- memoryRepositories
- AuthRepository
- PackagesPage.tsx
- templates.go
- AppointmentRepository
- compilerOptions
- compilerOptions
- compilerOptions
- ReminderWorker
- 006_erp_finance_stock_commissions.sql
- AppointmentsPage
- time.Duration
- ERP Agendamentos
- MembershipsPage.tsx
- time.Time
- booking/src/lib/apiError.ts
- Claims
- CommissionsPage.tsx
- ExportHandler
- observability.go
- PortalPage
- .UpdateAccountRole
- UseCase
- EventBroker
- writeResourceError
- writeServiceError
- Analytics
- SecretCipher
- DeviceSession
- bcrypt_hasher.go
- 005_resources.sql
- 007_jobs_and_calendar_sync.sql
- booking/src/shared/components/Button.tsx
- main.go
- widget_handler.go
- catalog.reviews
- body_limit.go
- booking/src/main.tsx
- graphify
- rules/graphify.md
- workflows/graphify.md
- generate-certs.sh
- CustomersPage
- .GetByID
- Queue
- erpagendamentos/backend
- Repository
- Repository
- UseCase
- JobRepository
- 010_branding_pos_and_cash_registers.sql

## God Nodes (most connected - your core abstractions)
1. `writeJSON()` - 129 edges
2. `resolveRequestTenantID()` - 106 edges
3. `adminFetch()` - 72 edges
4. `writeInvalidPayload()` - 69 edges
5. `Button` - 68 edges
6. `useToast()` - 66 edges
7. `UseCase` - 65 edges
8. `invalidateAdminGetCache()` - 62 edges
9. `writeInternalError()` - 57 edges
10. `RouterConfig` - 54 edges

## Surprising Connections (you probably didn't know these)
- `newResourceUseCase()` --calls--> `NewUseCase()`  [INFERRED]
  backend/internal/application/appointment/resource_test.go → backend/internal/application/appointment/usecase.go
- `newTestUseCase()` --calls--> `NewUseCase()`  [INFERRED]
  backend/internal/application/booking/usecase_test.go → backend/internal/application/booking/usecase.go
- `TestManagementCodeNormalization()` --calls--> `hashManagementCode()`  [INFERRED]
  backend/internal/application/booking/usecase_test.go → backend/internal/application/booking/usecase.go
- `buildAuthComponents()` --calls--> `isProductionEnv()`  [INFERRED]
  backend/internal/bootstrap/auth_components.go → backend/internal/bootstrap/container.go
- `buildAuthComponents()` --calls--> `seedAdminAccount()`  [INFERRED]
  backend/internal/bootstrap/auth_components.go → backend/internal/bootstrap/container.go

## Import Cycles
- None detected.

## Communities (155 total, 63 thin omitted)

### Community 0 - "adminAuth.ts"
Cohesion: 0.06
Nodes (65): AdminShell(), AdminShellProps, App(), handleLogout(), AdminNavigationItem, navigationGroups, navigationItems, AppProviders() (+57 more)

### Community 1 - "auth_cookie.go"
Cohesion: 0.13
Nodes (26): bearerAuthAllowed(), isPrivilegedAdminRole(), authCookieDomain(), authCookiePartitioned(), authCookieSameSite(), authCookieSecure(), authTokenFromCookie(), clearAuthCookie() (+18 more)

### Community 2 - "Payment"
Cohesion: 0.05
Nodes (23): UseCase, NewUseCase(), normalizeStatus(), buildPaymentProvider(), IntentResult, Payment, NewMercadoPagoProvider(), NewMockProvider() (+15 more)

### Community 3 - "adminApi.ts"
Cohesion: 0.07
Nodes (57): AppointmentFilters, createProduct(), createStockMovement(), deleteProduct(), fetchProducts(), fetchServices(), fetchServiceSupplies(), fetchStockMovements() (+49 more)

### Community 4 - "ref_react"
Cohesion: 0.08
Nodes (38): AppointmentFormModal(), computeEndsAt(), handleSubmit(), AppointmentFormModalProps, STATUS_OPTIONS, CheckoutComandaModalProps, SelectedProduct, createAppointment() (+30 more)

### Community 5 - "go_pkg_errors"
Cohesion: 0.10
Nodes (7): CheckoutInput, CheckoutResult, ProductItem, CategoryTotal, Type, Method, Status

### Community 7 - "Button"
Cohesion: 0.10
Nodes (50): CheckoutComandaModal(), createResource(), createService(), fetchConfiguration(), updateConfiguration(), updateResource(), updateService(), DAY_LABELS (+42 more)

### Community 8 - "UseCase"
Cohesion: 0.07
Nodes (11): UseCase, NewUseCase(), TestStockUseCase(), Product, Repository, ServiceSupply, StockMovement, NewStockRepository() (+3 more)

### Community 9 - "PageContainer"
Cohesion: 0.07
Nodes (63): createPayment(), deleteTenant(), fetchAppointmentAnalytics(), fetchCurrentTenant(), fetchNotifications(), fetchPayments(), fetchReviews(), fetchTenants() (+55 more)

### Community 10 - "email_provider.go"
Cohesion: 0.21
Nodes (7): buildNotificationProvider(), intEnv(), logFallback(), NewEmailProvider(), NewLogProvider(), EmailProvider, LogProvider

### Community 11 - "UseCase"
Cohesion: 0.07
Nodes (11): UseCase, NewUseCase(), TestCommissionUseCase(), CommissionRecord, CommissionRule, Repository, Summary, NewCommissionRepository() (+3 more)

### Community 12 - "adminFetch"
Cohesion: 0.07
Nodes (63): closeCashRegister(), connectCalendar(), createAdminUser(), createCashRegisterEntry(), createFinancialTransaction(), createGiftCard(), createMemberPlan(), createMembership() (+55 more)

### Community 13 - "container.go"
Cohesion: 0.09
Nodes (4): deploymentEnvironment(), isProductionEnv(), authClaimsContextKey, tokenPayload

### Community 14 - "Service"
Cohesion: 0.09
Nodes (10): UseCase, NewUseCase(), Service, scanServiceRow(), scanServiceRows(), NewServiceHandler(), ServiceRepository, ServiceRepository (+2 more)

### Community 15 - "net/http.Request"
Cohesion: 0.07
Nodes (33): requestWithAuthClaims(), requireMasterAdmin(), resolveAuthClaims(), resolveTenantID(), withProtectedNamespaces(), browserSessionPayload(), buildPasswordResetURL(), exposeAuthTokenForTests() (+25 more)

### Community 16 - "Tenant"
Cohesion: 0.08
Nodes (13): UseCase, NewUseCase(), ensureDefaultTenant(), Tenant, isTenantIdentityConflict(), scanTenantRow(), scanTenantRows(), NewTenantHandler() (+5 more)

### Community 17 - "UseCase"
Cohesion: 0.10
Nodes (6): UseCase, Membership, NewMembershipRepository(), NewCommerceHandler(), MembershipRepository, MembershipRepository

### Community 18 - "Customer"
Cohesion: 0.09
Nodes (13): UseCase, Customer, NewCustomerRepository(), onlyDigits(), scanCustomerRow(), scanCustomerRows(), NewCustomerHandler(), NewExportHandler() (+5 more)

### Community 19 - "UseCase"
Cohesion: 0.09
Nodes (7): UseCase, Repository, Summary, Transaction, NewFinanceHandler(), FinanceRepository, FinanceRepository

### Community 20 - "Professional"
Cohesion: 0.11
Nodes (9): UseCase, Professional, scanProfessionalRow(), scanProfessionalRows(), NewProfessionalHandler(), ProfessionalRepository, ProfessionalRepository, professionalScanner (+1 more)

### Community 21 - "net/http.ResponseWriter"
Cohesion: 0.09
Nodes (13): isAdminRole(), resolveRequestTenantID(), writeConfigurationError(), writeProfessionalError(), writeInternalError(), writeWaitingListError(), CalendarSyncHandler, CommissionHandler (+5 more)

### Community 22 - "Entry"
Cohesion: 0.11
Nodes (11): UseCase, NewUseCase(), Entry, NewWaitingListRepository(), scanWaitingRow(), scanWaitingRows(), NewWaitingListHandler(), WaitingListRepository (+3 more)

### Community 23 - "admin/package.json"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+24 more)

### Community 24 - "database/sql.DB"
Cohesion: 0.12
Nodes (31): applyDatabaseSQL(), applyMigrationFile(), applySQLDirectory(), applySQLFile(), initializeRepositories(), postgresRepositories(), NewAppointmentRepository(), NewAuthRepository() (+23 more)

### Community 25 - "AppointmentHandler"
Cohesion: 0.14
Nodes (8): writeAppointmentError(), resolveActorTenantID(), writeAdminAuthError(), writeBadRequestError(), writeConflictError(), writeReviewError(), AppointmentHandler, ReviewHandler

### Community 26 - "booking/package.json"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+24 more)

### Community 27 - "TeamPage.tsx"
Cohesion: 0.40
Nodes (9): deleteAdminUser(), fetchAdminUsers(), updateAdminUserRole(), roleLabel(), ROLE_OPTIONS, TeamPage(), handleDelete(), handleRoleChange() (+1 more)

### Community 28 - "Appointment"
Cohesion: 0.11
Nodes (6): Repository, Appointment, StatusCount, parseAppointmentBound(), parseAppointmentInterval(), AppointmentRepository

### Community 29 - "UseCase"
Cohesion: 0.10
Nodes (14): LoginInput, MFAStatus, RegisterInput, SecretCipher, Session, TokenManager, generateTOTPSecret(), totpCode() (+6 more)

### Community 30 - "Account"
Cohesion: 0.13
Nodes (6): Account, PasswordResetToken, UpdateProfileInput, cloneAuthAccount(), resetTokenExpired(), AuthRepository

### Community 31 - "Connection"
Cohesion: 0.22
Nodes (3): Connection, CalendarRepository, CalendarRepository

### Community 33 - "rate_limit.go"
Cohesion: 0.08
Nodes (26): UseCase, NewUseCase(), Review, NewReviewRepository(), scanReviewRow(), scanReviewRows(), forwardedClientIP(), isSensitiveRateLimitedRoute() (+18 more)

### Community 34 - "RouterConfig"
Cohesion: 0.20
Nodes (31): hasRole(), requireAdmin(), requireOwnerOrAdmin(), requirePrivilegedAdmin(), requireRoles(), NewRouter(), registerAppointmentRoutes(), registerConfigurationRoutes() (+23 more)

### Community 35 - "landing/package.json"
Cohesion: 0.06
Nodes (30): dependencies, lucide-react, react, react-dom, devDependencies, tailwindcss, @tailwindcss/vite, @types/react (+22 more)

### Community 36 - ".Add"
Cohesion: 0.15
Nodes (13): ChangePasswordInput, PasswordResetCodeConfirmInput, PasswordResetConfirmInput, PasswordResetRequestInput, ValidationError, PasswordHasher, hashResetToken(), NewValidationError() (+5 more)

### Community 37 - "booking/src/pages/WaitingListPage.tsx"
Cohesion: 0.17
Nodes (22): App(), fetchPublicReview(), signUpWaitingList(), submitPublicReview(), DetailsStep(), ReviewPage(), handleSubmit(), SlugEntryPage() (+14 more)

### Community 38 - "MemberPlan"
Cohesion: 0.14
Nodes (6): NewUseCase(), MemberPlan, NewMemberPlanRepository(), scanMemberPlans(), MemberPlanRepository, MemberPlanRepository

### Community 39 - "CashRegister"
Cohesion: 0.11
Nodes (8): UseCase, CashRegister, CashRegisterEntry, Repository, ShiftSummary, NewCashRegisterHandler(), CashRegisterRepository, CashRegisterRepository

### Community 41 - "UseCase"
Cohesion: 0.18
Nodes (10): appointmentDurationMinutes(), formatAppointmentForTemplate(), generateManagementCode(), UseCase, hashManagementCode(), parseRFC3339(), BookingInput, PortalAppointment (+2 more)

### Community 42 - "Configuration"
Cohesion: 0.25
Nodes (10): workingHoursFor(), Configuration, cloneConfiguration(), defaultMemoryConfiguration(), normalizeConfigurationTenantID(), normalizePostgresTenantID(), ClosedDay, WorkingHours (+2 more)

### Community 43 - "GiftCard"
Cohesion: 0.11
Nodes (6): generateGiftCardCode(), GiftCard, NewGiftCardRepository(), GiftCardRepository, GiftCardRepository, GiftCardRepository

### Community 44 - "PortalPage.tsx"
Cohesion: 0.19
Nodes (23): PortalAppointment, PortalResult, Slot, formatDateLong(), formatDateTime(), formatTime(), fromDateKey(), nextDateKeys() (+15 more)

### Community 45 - "AdminRoutes.tsx"
Cohesion: 0.15
Nodes (25): AccessDeniedPage(), AdminRouteFallback(), AdminRoutes(), AdminRoutesProps, AgendaPage, AnalyticsPage, AppointmentsPage, CommissionsPage (+17 more)

### Community 46 - "PackageCredit"
Cohesion: 0.12
Nodes (6): PackageCredit, NewPackageCreditRepository(), scanPackageCredits(), PackageCreditRepository, PackageCreditRepository, PackageCreditRepository

### Community 47 - "Package"
Cohesion: 0.12
Nodes (5): Package, NewPackageRepository(), PackageRepository, PackageRepository, PackageRepository

### Community 48 - "bookingApi.ts"
Cohesion: 0.14
Nodes (23): backendOrigin, buildApiUrl(), AvailabilityResponse, cancelPublicAppointment(), fetchAvailability(), fetchPublicAppointment(), fetchPublicInfo(), fetchPublicProfessionals() (+15 more)

### Community 49 - "landing/src/App.tsx"
Cohesion: 0.16
Nodes (14): App(), CtaSection(), FEATURES, FeaturesSection(), Footer(), Header(), NAV_ITEMS, HeroSection() (+6 more)

### Community 50 - "Resource"
Cohesion: 0.10
Nodes (10): NewUseCase(), UseCase, NewUseCase(), Resource, scanResourceRow(), NewResourceHandler(), ResourceRepository, ResourceRepository (+2 more)

### Community 51 - "Container"
Cohesion: 0.11
Nodes (18): lazyAppHandler, writeLazyJSON(), NewUseCase(), UseCase, NewUseCase(), NewUseCase(), UseCase, NewUseCase() (+10 more)

### Community 52 - "002_features.sql"
Cohesion: 0.14
Nodes (20): catalog.gift_cards, catalog.member_plans, catalog.memberships, catalog.notification_logs, catalog.package_credits, catalog.packages, catalog.payments, catalog.service_professionals (+12 more)

### Community 53 - "writeJSON"
Cohesion: 0.23
Nodes (5): writeBookingError(), writePaymentError(), writeJSON(), BookingHandler, PaymentHandler

### Community 54 - "writeInvalidPayload"
Cohesion: 0.29
Nodes (3): writeCommerceError(), writeInvalidPayload(), CommerceHandler

### Community 55 - "writeNotFoundError"
Cohesion: 0.18
Nodes (6): writeCustomerError(), writeNotificationError(), writeNotFoundError(), writeTenantError(), CustomerHandler, TenantHandler

### Community 56 - "001_baseline.sql"
Cohesion: 0.17
Nodes (21): catalog.appointments, catalog.auth_rate_limits, catalog.auth_sessions, catalog.configurations, catalog.customers, catalog.password_reset_tokens, catalog.professionals, catalog.services (+13 more)

### Community 57 - "BookingPage.tsx"
Cohesion: 0.13
Nodes (15): BookingResult, createPublicBooking(), formatMoney(), BookingPage(), handleBookingSubmit(), CenteredLayout(), EmptyNote(), ErrorCard() (+7 more)

### Community 58 - "UseCase"
Cohesion: 0.16
Nodes (6): RecurrenceInput, UseCase, parseRFC3339(), validStatusChange(), NewCalendarHandler(), CalendarHandler

### Community 59 - "memoryRepositories"
Cohesion: 0.11
Nodes (31): newAppointment(), newResourceUseCase(), TestCreateAutoAssignsAvailableResource(), TestCreateRejectsBusyExplicitResource(), TestCreateRejectsUnknownResource(), TestUpdateAutoReassignsFreedResource(), newTestUseCase(), TestAvailableSlotsRespectsClosedDays() (+23 more)

### Community 61 - "PackagesPage.tsx"
Cohesion: 0.17
Nodes (23): deleteGiftCard(), deletePackage(), deleteResource(), fetchGiftCards(), fetchPackageCredits(), fetchPackages(), fetchResources(), setPackageActive() (+15 more)

### Community 62 - "templates.go"
Cohesion: 0.27
Nodes (16): greeting(), PublicBookingOrigin(), RenderCancellation(), RenderConfirmation(), RenderReminder(), RenderReschedule(), sampleData(), TestGreetingWithoutName() (+8 more)

### Community 63 - "AppointmentRepository"
Cohesion: 0.18
Nodes (6): appointmentSelectColumns(), nullIfEmpty(), scanAppointmentRow(), scanAppointmentRows(), AppointmentRepository, appointmentScanner

### Community 65 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 66 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 67 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 68 - "ReminderWorker"
Cohesion: 0.18
Nodes (8): MaintenanceWorker, MaintenanceConfig, NewMaintenanceWorker(), formatAppointmentLocal(), NewReminderWorker(), startBackgroundWorkers(), authComponents, ReminderWorker

### Community 69 - "006_erp_finance_stock_commissions.sql"
Cohesion: 0.21
Nodes (12): catalog.commission_records, catalog.commission_rules, catalog.financial_transactions, catalog.products, catalog.service_supplies, catalog.stock_movements, idx_catalog_commission_records_prof, idx_catalog_commission_rules_tenant (+4 more)

### Community 70 - "AppointmentsPage"
Cohesion: 0.18
Nodes (22): handleSubmit(), checkInAppointment(), checkoutAppointment(), deleteAppointment(), fetchAppointments(), fetchAppointmentSummary(), invalidateMutations(), markAppointmentNoShow() (+14 more)

### Community 71 - "time.Duration"
Cohesion: 0.19
Nodes (5): buildAuthComponents(), boolEnv(), durationSecondsEnv(), CleanupSummary, NewJWTManager()

### Community 72 - "ERP Agendamentos"
Cohesion: 0.14
Nodes (13): CI, Endpoints principais, ERP Agendamentos, Estrutura, Migrations, Pré-requisitos, Rodando com Docker (produção local), Rodando em desenvolvimento (sem banco) (+5 more)

### Community 73 - "MembershipsPage.tsx"
Cohesion: 0.33
Nodes (12): deleteMembership(), fetchMemberPlans(), fetchMemberships(), setMembershipActive(), updateMemberPlanStatus(), MembershipsPage(), handleDelete(), handlePlanStatusChange() (+4 more)

### Community 74 - "time.Time"
Cohesion: 0.23
Nodes (11): countOverlappingSlots(), isClosedDay(), localMidnight(), overlapsAnyBusy(), parseClock(), BuildICalCalendar(), formatUTC(), sanitizeLine() (+3 more)

### Community 75 - "booking/src/lib/apiError.ts"
Cohesion: 0.29
Nodes (8): ApiError, ApiErrorPayload, buildApiError(), parseApiResponse(), API_ERROR_MESSAGES, FIELD_ERROR_MESSAGES, getApiErrorMessage(), translateFieldErrors()

### Community 76 - "Claims"
Cohesion: 0.31
Nodes (5): isProductionAuthEnvironment(), Claims, decodeSegment(), encodeSegment(), JWTManager

### Community 77 - "CommissionsPage.tsx"
Cohesion: 0.21
Nodes (20): createCommissionRule(), deleteCommissionRule(), deleteProfessional(), fetchCommissionRecords(), fetchCommissionRules(), fetchCommissionSummaries(), fetchProfessionals(), setProfessionalActive() (+12 more)

### Community 79 - "observability.go"
Cohesion: 0.43
Nodes (6): isSensitiveLogKey(), redactSensitiveText(), sanitizeLogString(), sanitizeLogValue(), truncateLogString(), contextKey

### Community 80 - "PortalPage"
Cohesion: 0.39
Nodes (8): portalCancelAppointment(), portalLookup(), portalRescheduleAppointment(), PortalPage(), handleCancel(), handleConfirmReschedule(), handleLookup(), refresh()

### Community 81 - ".UpdateAccountRole"
Cohesion: 0.39
Nodes (4): canManageTenantAccount(), isPrivilegedAdminRole(), isStaffRole(), normalizeTenantScope()

### Community 82 - "UseCase"
Cohesion: 0.15
Nodes (10): UseCase, NewUseCase(), Log, NewNotificationRepository(), NewNotificationHandler(), NotificationRepository, Provider, Repository (+2 more)

### Community 83 - "EventBroker"
Cohesion: 0.20
Nodes (10): UseCase, NewAppointmentHandler(), NewBookingHandler(), NewCheckoutHandler(), EventBroker, NewMercadoPagoWebhookHandler(), RegisterMetricsBroker(), CheckoutHandler (+2 more)

### Community 86 - "Analytics"
Cohesion: 0.33
Nodes (4): DayMetric, ProfessionalMetric, ServiceMetric, Analytics

### Community 91 - "005_resources.sql"
Cohesion: 0.40
Nodes (3): catalog.resources, idx_appointments_resource, idx_resources_tenant

### Community 92 - "007_jobs_and_calendar_sync.sql"
Cohesion: 0.47
Nodes (4): catalog.calendar_connections, catalog.jobs, idx_calendar_connections_tenant, idx_catalog_jobs_poll

### Community 93 - "booking/src/shared/components/Button.tsx"
Cohesion: 0.33
Nodes (5): ButtonProps, ButtonSize, ButtonVariant, sizeClasses, variantClasses

### Community 94 - "main.go"
Cohesion: 0.60
Nodes (4): newLazyAppHandler(), main(), resolveDurationEnv(), resolvePositiveIntEnv()

### Community 95 - "widget_handler.go"
Cohesion: 0.60
Nodes (3): NewWidgetHandler(), resolveWidgetBookingOrigin(), WidgetHandler

### Community 96 - "catalog.reviews"
Cohesion: 0.60
Nodes (3): catalog.reviews, idx_reviews_tenant_created, idx_reviews_tenant_status

### Community 97 - "body_limit.go"
Cohesion: 0.83
Nodes (3): resolveMaxRequestBodyBytes(), shouldLimitRequestBody(), withRequestBodyLimit()

### Community 113 - "CustomersPage"
Cohesion: 0.27
Nodes (11): deleteCustomer(), downloadAdminCsv(), fetchCustomers(), handleExport(), CustomersPage(), handleDelete(), handleExport(), handleSearch() (+3 more)

### Community 116 - "Queue"
Cohesion: 0.27
Nodes (3): Queue, HandlerFunc, MetricsRegistry

### Community 125 - "Repository"
Cohesion: 0.20
Nodes (4): NewQueue(), TestJobQueue(), Repository, NewJobRepository()

### Community 126 - "Repository"
Cohesion: 0.22
Nodes (4): NewUseCase(), TestCalendarUseCase(), Repository, NewCalendarRepository()

### Community 127 - "UseCase"
Cohesion: 0.25
Nodes (3): UseCase, SyncEvent, NewCalendarSyncHandler()

### Community 130 - "010_branding_pos_and_cash_registers.sql"
Cohesion: 0.60
Nodes (4): catalog.cash_register_entries, catalog.cash_registers, idx_catalog_cash_reg_entries_reg, idx_catalog_cash_registers_tenant

## Knowledge Gaps
- **263 isolated node(s):** `graphify-mcp`, `name`, `private`, `version`, `type` (+258 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 510 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **63 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Container` connect `Container` to `auth_cookie.go`, `Payment`, `UseCase`, `UseCase`, `container.go`, `Service`, `Tenant`, `UseCase`, `Customer`, `UseCase`, `Professional`, `Entry`, `database/sql.DB`, `UseCase`, `rate_limit.go`, `CashRegister`, `UseCase`, `Resource`, `UseCase`, `ReminderWorker`, `time.Duration`, `UseCase`, `EventBroker`, `Queue`, `UseCase`?**
  _High betweenness centrality (0.037) - this node is a cross-community bridge._
- **Why does `RouterConfig` connect `RouterConfig` to `Payment`, `go_pkg_errors`, `UseCase`, `UseCase`, `Service`, `Tenant`, `UseCase`, `Customer`, `UseCase`, `Professional`, `Entry`, `UseCase`, `rate_limit.go`, `CashRegister`, `UseCase`, `Resource`, `Container`, `UseCase`, `ReminderWorker`, `time.Duration`, `UseCase`, `EventBroker`, `UseCase`?**
  _High betweenness centrality (0.031) - this node is a cross-community bridge._
- **Why does `UseCase` connect `UseCase` to `auth_cookie.go`, `RouterConfig`, `rate_limit.go`, `.Add`, `ReminderWorker`, `time.Duration`, `Repository`, `time.Time`, `Claims`, `container.go`, `net/http.Request`, `.UpdateAccountRole`, `Container`, `net/http.ResponseWriter`, `DeviceSession`, `AppointmentHandler`, `Account`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Are the 125 inferred relationships involving `writeJSON()` (e.g. with `health()` and `readiness()`) actually correct?**
  _`writeJSON()` has 125 INFERRED edges - model-reasoned connections that need verification._
- **Are the 101 inferred relationships involving `resolveRequestTenantID()` (e.g. with `.Analytics()` and `.CheckIn()`) actually correct?**
  _`resolveRequestTenantID()` has 101 INFERRED edges - model-reasoned connections that need verification._
- **What connects `graphify-mcp`, `name`, `private` to the rest of the system?**
  _263 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `adminAuth.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.05555555555555555 - nodes in this community are weakly interconnected._
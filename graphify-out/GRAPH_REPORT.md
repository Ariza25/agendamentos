# Graph Report - erpagendamentos  (2026-11-01)

## Corpus Check
- 259 files · ~110,656 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 15 file(s) not represented in the graph (top: (none) 6, .conf 5, .css 3)

## Summary
- 2634 nodes · 8454 edges · 131 communities (81 shown, 50 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 627 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `1d364b09`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- adminAuth.ts
- writeAPIError
- adminApi.ts
- Resource
- RouterConfig
- AppointmentsPage.tsx
- net/http.Request
- Professional
- Button
- go_pkg_time
- UseCase
- ref_react
- UseCase
- adminFetch
- Payment
- auth/usecase.go
- writeJSON
- Tenant
- rate_limit.go
- container.go
- Customer
- Membership
- booking/src/pages/WaitingListPage.tsx
- Entry
- database/sql.DB
- admin/package.json
- booking/package.json
- TeamPage.tsx
- Appointment
- Connection
- Queue
- landing/package.json
- GiftCard
- AuthRepository
- PortalPage.tsx
- net/http.ResponseWriter
- UseCase
- Account
- Service
- AdminRoutes.tsx
- .registerWithRole
- UseCase
- Review
- landing/src/App.tsx
- UseCase
- PackageCredit
- 002_features.sql
- MemberPlan
- bookingApi.ts
- AuthRepository
- Configuration
- payment/repository.go
- CommissionsPage.tsx
- 001_baseline.sql
- BookingPage.tsx
- AppointmentRepository
- AppointmentHandler
- WaitingListPage
- compilerOptions
- UseCase
- ReminderWorker
- compilerOptions
- compilerOptions
- 006_erp_finance_stock_commissions.sql
- FinancePage.tsx
- buildAuthComponents
- PaymentsPage.tsx
- ExportHandler
- ERP Agendamentos
- MembershipsPage.tsx
- UseCase
- email_provider.go
- writeInvalidPayload
- booking/src/lib/apiError.ts
- .GetByID
- .UpdateAccountRole
- writeNotFoundError
- Container
- time.Duration
- PortalPage
- writeProfessionalError
- writeResourceError
- writeServiceError
- SecretCipher
- ical.go
- body_limit.go
- booking/src/main.tsx
- generate-certs.sh
- 005_resources.sql
- 007_jobs_and_calendar_sync.sql
- booking/src/shared/components/Button.tsx
- widget_handler.go
- catalog.reviews
- graphify
- rules/graphify.md
- workflows/graphify.md
- erpagendamentos/backend

## God Nodes (most connected - your core abstractions)
1. `writeJSON()` - 123 edges
2. `resolveRequestTenantID()` - 100 edges
3. `adminFetch()` - 68 edges
4. `Button` - 66 edges
5. `UseCase` - 65 edges
6. `writeInvalidPayload()` - 65 edges
7. `useToast()` - 64 edges
8. `invalidateAdminGetCache()` - 59 edges
9. `writeInternalError()` - 57 edges
10. `RouterConfig` - 50 edges

## Surprising Connections (you probably didn't know these)
- `NewContainer()` --calls--> `loadEnvFile()`  [INFERRED]
  backend/internal/bootstrap/container.go → backend/internal/bootstrap/env.go
- `newTestUseCase()` --calls--> `NewUseCase()`  [INFERRED]
  backend/internal/application/booking/usecase_test.go → backend/internal/application/booking/usecase.go
- `TestManagementCodeNormalization()` --calls--> `hashManagementCode()`  [INFERRED]
  backend/internal/application/booking/usecase_test.go → backend/internal/application/booking/usecase.go
- `buildAuthComponents()` --calls--> `boolEnv()`  [INFERRED]
  backend/internal/bootstrap/auth_components.go → backend/internal/bootstrap/container.go
- `buildAuthComponents()` --calls--> `durationSecondsEnv()`  [INFERRED]
  backend/internal/bootstrap/auth_components.go → backend/internal/bootstrap/container.go

## Import Cycles
- None detected.

## Communities (131 total, 50 thin omitted)

### Community 0 - "adminAuth.ts"
Cohesion: 0.05
Nodes (66): AdminShell(), AdminShellProps, App(), handleLogout(), AdminNavigationItem, navigationGroups, navigationItems, AppProviders() (+58 more)

### Community 1 - "writeAPIError"
Cohesion: 0.10
Nodes (25): requestWithAuthClaims(), requireMasterAdmin(), resolveAuthClaims(), resolveTenantID(), withProtectedNamespaces(), browserSessionPayload(), exposeAuthTokenForTests(), AuthHandler (+17 more)

### Community 2 - "adminApi.ts"
Cohesion: 0.08
Nodes (50): AppointmentFilters, createProduct(), createStockMovement(), deleteProduct(), fetchProducts(), fetchServices(), fetchServiceSupplies(), fetchStockMovements() (+42 more)

### Community 3 - "Resource"
Cohesion: 0.06
Nodes (40): newAppointment(), newResourceUseCase(), TestCreateAutoAssignsAvailableResource(), TestCreateRejectsBusyExplicitResource(), TestCreateRejectsUnknownResource(), TestUpdateAutoReassignsFreedResource(), NewUseCase(), newTestUseCase() (+32 more)

### Community 4 - "RouterConfig"
Cohesion: 0.05
Nodes (49): UseCase, UseCase, Repository, Summary, Transaction, hasRole(), requireAdmin(), requireOwnerOrAdmin() (+41 more)

### Community 5 - "AppointmentsPage.tsx"
Cohesion: 0.09
Nodes (50): computeEndsAt(), handleSubmit(), AppointmentFormModalProps, STATUS_OPTIONS, checkInAppointment(), createAppointment(), deleteAppointment(), fetchAppointmentAnalytics() (+42 more)

### Community 6 - "net/http.Request"
Cohesion: 0.10
Nodes (11): isAdminRole(), resolveRequestTenantID(), authTokenFromCookie(), writeCommerceError(), requiresCSRFProtection(), writeInternalError(), CalendarSyncHandler, CommerceHandler (+3 more)

### Community 8 - "Professional"
Cohesion: 0.11
Nodes (9): UseCase, Professional, scanProfessionalRow(), scanProfessionalRows(), NewProfessionalHandler(), ProfessionalRepository, ProfessionalRepository, professionalScanner (+1 more)

### Community 9 - "Button"
Cohesion: 0.11
Nodes (52): AppointmentFormModal(), createCustomer(), createService(), fetchConfiguration(), fetchNotifications(), updateConfiguration(), updateCustomer(), updateService() (+44 more)

### Community 11 - "UseCase"
Cohesion: 0.07
Nodes (8): UseCase, Product, Repository, ServiceSupply, StockMovement, NewStockHandler(), StockRepository, StockRepository

### Community 12 - "ref_react"
Cohesion: 0.08
Nodes (56): deleteGiftCard(), deletePackage(), deleteTenant(), fetchCurrentTenant(), fetchGiftCards(), fetchPackageCredits(), fetchPackages(), fetchReviews() (+48 more)

### Community 13 - "UseCase"
Cohesion: 0.07
Nodes (8): UseCase, CommissionRecord, CommissionRule, Repository, Summary, NewCommissionHandler(), CommissionRepository, CommissionRepository

### Community 14 - "adminFetch"
Cohesion: 0.08
Nodes (47): connectCalendar(), createGiftCard(), createMemberPlan(), createMembership(), createPackage(), createPackageCredit(), createPayment(), createProfessional() (+39 more)

### Community 15 - "Payment"
Cohesion: 0.05
Nodes (21): UseCase, normalizeStatus(), buildPaymentProvider(), IntentResult, Payment, NewMercadoPagoProvider(), NewMockProvider(), mapStripeStatus() (+13 more)

### Community 16 - "auth/usecase.go"
Cohesion: 0.07
Nodes (13): newLazyAppHandler(), main(), resolveDurationEnv(), resolvePositiveIntEnv(), loadEnvFile(), isSensitiveLogKey(), redactSensitiveText(), sanitizeLogString() (+5 more)

### Community 17 - "writeJSON"
Cohesion: 0.20
Nodes (6): writeBookingError(), health(), writePaymentError(), writeJSON(), BookingHandler, PaymentHandler

### Community 18 - "Tenant"
Cohesion: 0.09
Nodes (11): UseCase, Tenant, isTenantIdentityConflict(), scanTenantRow(), scanTenantRows(), NewTenantHandler(), TenantRepository, SignupInput (+3 more)

### Community 19 - "rate_limit.go"
Cohesion: 0.08
Nodes (41): lazyAppHandler, writeLazyJSON(), bearerAuthAllowed(), isPrivilegedAdminRole(), authCookieDomain(), authCookiePartitioned(), authCookieSameSite(), authCookieSecure() (+33 more)

### Community 20 - "container.go"
Cohesion: 0.11
Nodes (4): NewUseCase(), deploymentEnvironment(), isProductionEnv(), Type

### Community 21 - "Customer"
Cohesion: 0.10
Nodes (12): UseCase, Customer, onlyDigits(), scanCustomerRow(), scanCustomerRows(), NewCustomerHandler(), NewExportHandler(), BookingResult (+4 more)

### Community 22 - "Membership"
Cohesion: 0.12
Nodes (4): Membership, MembershipRepository, MembershipRepository, MembershipRepository

### Community 23 - "booking/src/pages/WaitingListPage.tsx"
Cohesion: 0.17
Nodes (22): App(), fetchPublicReview(), signUpWaitingList(), submitPublicReview(), DetailsStep(), ReviewPage(), handleSubmit(), SlugEntryPage() (+14 more)

### Community 24 - "Entry"
Cohesion: 0.11
Nodes (10): UseCase, Entry, NewWaitingListRepository(), scanWaitingRow(), scanWaitingRows(), NewWaitingListHandler(), WaitingListRepository, WaitingListRepository (+2 more)

### Community 25 - "database/sql.DB"
Cohesion: 0.12
Nodes (30): applyDatabaseSQL(), applyMigrationFile(), applySQLDirectory(), applySQLFile(), initializeRepositories(), postgresRepositories(), NewAppointmentRepository(), NewAuthRepository() (+22 more)

### Community 26 - "admin/package.json"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+24 more)

### Community 27 - "booking/package.json"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+24 more)

### Community 28 - "TeamPage.tsx"
Cohesion: 0.30
Nodes (11): createAdminUser(), deleteAdminUser(), fetchAdminUsers(), updateAdminUserRole(), roleLabel(), ROLE_OPTIONS, TeamPage(), handleDelete() (+3 more)

### Community 29 - "Appointment"
Cohesion: 0.11
Nodes (6): Repository, Appointment, NewAppointmentRepository(), parseAppointmentBound(), parseAppointmentInterval(), AppointmentRepository

### Community 30 - "Connection"
Cohesion: 0.10
Nodes (7): UseCase, Connection, Repository, SyncEvent, NewCalendarSyncHandler(), CalendarRepository, CalendarRepository

### Community 31 - "Queue"
Cohesion: 0.08
Nodes (6): Queue, Job, Repository, HandlerFunc, JobRepository, JobRepository

### Community 32 - "landing/package.json"
Cohesion: 0.06
Nodes (30): dependencies, lucide-react, react, react-dom, devDependencies, tailwindcss, @tailwindcss/vite, @types/react (+22 more)

### Community 33 - "GiftCard"
Cohesion: 0.11
Nodes (5): generateGiftCardCode(), GiftCard, GiftCardRepository, GiftCardRepository, GiftCardRepository

### Community 34 - "AuthRepository"
Cohesion: 0.08
Nodes (6): AuthState, DeviceSession, PasswordResetToken, resetTokenExpired(), sessionExpired(), AuthRepository

### Community 35 - "PortalPage.tsx"
Cohesion: 0.19
Nodes (23): PortalAppointment, PortalResult, Slot, formatDateLong(), formatDateTime(), formatTime(), fromDateKey(), nextDateKeys() (+15 more)

### Community 36 - "net/http.ResponseWriter"
Cohesion: 0.13
Nodes (11): buildPasswordResetURL(), AuthHandler, secureTemporaryPassword(), writeAdminAuthError(), newAuthMFAHandler(), HandleMetrics(), writeConflictError(), writeRawJSON() (+3 more)

### Community 37 - "UseCase"
Cohesion: 0.11
Nodes (11): LoginInput, MFAStatus, SecretCipher, Session, generateTOTPSecret(), totpCode(), totpURI(), validateTOTP() (+3 more)

### Community 38 - "Account"
Cohesion: 0.12
Nodes (4): Repository, Account, UpdateProfileInput, cloneAuthAccount()

### Community 39 - "Service"
Cohesion: 0.10
Nodes (8): UseCase, Service, scanServiceRow(), scanServiceRows(), NewServiceHandler(), ServiceRepository, ServiceRepository, serviceScanner

### Community 40 - "AdminRoutes.tsx"
Cohesion: 0.15
Nodes (25): AccessDeniedPage(), AdminRouteFallback(), AdminRoutes(), AdminRoutesProps, AgendaPage, AnalyticsPage, AppointmentsPage, CommissionsPage (+17 more)

### Community 41 - ".registerWithRole"
Cohesion: 0.13
Nodes (13): ChangePasswordInput, PasswordResetCodeConfirmInput, PasswordResetConfirmInput, PasswordResetRequestInput, RegisterInput, ValidationError, hashResetToken(), isStaffRole() (+5 more)

### Community 42 - "UseCase"
Cohesion: 0.11
Nodes (5): UseCase, Package, NewCommerceHandler(), PackageRepository, PackageRepository

### Community 43 - "Review"
Cohesion: 0.13
Nodes (8): Review, NewReviewRepository(), scanReviewRow(), scanReviewRows(), ReviewRepository, ReviewRepository, reviewScanner, Repository

### Community 44 - "landing/src/App.tsx"
Cohesion: 0.16
Nodes (14): App(), CtaSection(), FEATURES, FeaturesSection(), Footer(), Header(), NAV_ITEMS, HeroSection() (+6 more)

### Community 45 - "UseCase"
Cohesion: 0.18
Nodes (10): appointmentDurationMinutes(), formatAppointmentForTemplate(), generateManagementCode(), UseCase, hashManagementCode(), parseRFC3339(), BookingInput, PortalAppointment (+2 more)

### Community 46 - "PackageCredit"
Cohesion: 0.13
Nodes (5): PackageCredit, scanPackageCredits(), PackageCreditRepository, PackageCreditRepository, PackageCreditRepository

### Community 47 - "002_features.sql"
Cohesion: 0.14
Nodes (20): catalog.gift_cards, catalog.member_plans, catalog.memberships, catalog.notification_logs, catalog.package_credits, catalog.packages, catalog.payments, catalog.service_professionals (+12 more)

### Community 48 - "MemberPlan"
Cohesion: 0.14
Nodes (6): NewUseCase(), MemberPlan, scanMemberPlans(), MemberPlanRepository, MemberPlanRepository, MemberPlanRepository

### Community 49 - "bookingApi.ts"
Cohesion: 0.14
Nodes (23): backendOrigin, buildApiUrl(), AvailabilityResponse, cancelPublicAppointment(), fetchAvailability(), fetchPublicAppointment(), fetchPublicInfo(), fetchPublicProfessionals() (+15 more)

### Community 51 - "Configuration"
Cohesion: 0.13
Nodes (19): countOverlappingSlots(), isClosedDay(), localMidnight(), overlapsAnyBusy(), parseClock(), workingHoursFor(), Configuration, cloneConfiguration() (+11 more)

### Community 53 - "CommissionsPage.tsx"
Cohesion: 0.18
Nodes (22): createCommissionRule(), deleteCommissionRule(), deleteProfessional(), fetchCommissionRecords(), fetchCommissionRules(), fetchCommissionSummaries(), fetchProfessionals(), setProfessionalActive() (+14 more)

### Community 54 - "001_baseline.sql"
Cohesion: 0.17
Nodes (21): catalog.appointments, catalog.auth_rate_limits, catalog.auth_sessions, catalog.configurations, catalog.customers, catalog.password_reset_tokens, catalog.professionals, catalog.services (+13 more)

### Community 55 - "BookingPage.tsx"
Cohesion: 0.13
Nodes (15): BookingResult, createPublicBooking(), formatMoney(), BookingPage(), handleBookingSubmit(), CenteredLayout(), EmptyNote(), ErrorCard() (+7 more)

### Community 56 - "AppointmentRepository"
Cohesion: 0.17
Nodes (6): StatusCount, appointmentSelectColumns(), scanAppointmentRow(), scanAppointmentRows(), AppointmentRepository, appointmentScanner

### Community 57 - "AppointmentHandler"
Cohesion: 0.26
Nodes (3): writeAppointmentError(), resolveActorTenantID(), AppointmentHandler

### Community 58 - "WaitingListPage"
Cohesion: 0.70
Nodes (5): fetchWaitingList(), updateWaitingListStatus(), WaitingListPage(), handleStatusChange(), loadEntries()

### Community 59 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 60 - "UseCase"
Cohesion: 0.12
Nodes (8): DayMetric, ProfessionalMetric, RecurrenceInput, ServiceMetric, UseCase, parseRFC3339(), validStatusChange(), Analytics

### Community 61 - "ReminderWorker"
Cohesion: 0.21
Nodes (6): MaintenanceWorker, MaintenanceConfig, NewMaintenanceWorker(), formatAppointmentLocal(), NewReminderWorker(), ReminderWorker

### Community 62 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 63 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 64 - "006_erp_finance_stock_commissions.sql"
Cohesion: 0.21
Nodes (12): catalog.commission_records, catalog.commission_rules, catalog.financial_transactions, catalog.products, catalog.service_supplies, catalog.stock_movements, idx_catalog_commission_records_prof, idx_catalog_commission_rules_tenant (+4 more)

### Community 65 - "FinancePage.tsx"
Cohesion: 0.25
Nodes (15): createFinancialTransaction(), deleteFinancialTransaction(), fetchFinancialSummary(), fetchFinancialTransactions(), updateFinancialTransaction(), FinancePage(), handleCreate(), handleDelete() (+7 more)

### Community 66 - "buildAuthComponents"
Cohesion: 0.12
Nodes (14): TokenManager, PasswordHasher, isProductionAuthEnvironment(), NewUseCase(), buildAuthComponents(), seedAdminAccount(), Claims, NewBcryptHasher() (+6 more)

### Community 67 - "PaymentsPage.tsx"
Cohesion: 0.16
Nodes (21): deleteCustomer(), downloadAdminCsv(), fetchCustomers(), fetchPayments(), markPaymentPaid(), refundPayment(), handleExport(), CustomersPage() (+13 more)

### Community 69 - "ERP Agendamentos"
Cohesion: 0.14
Nodes (13): CI, Endpoints principais, ERP Agendamentos, Estrutura, Migrations, Pré-requisitos, Rodando com Docker (produção local), Rodando em desenvolvimento (sem banco) (+5 more)

### Community 70 - "MembershipsPage.tsx"
Cohesion: 0.33
Nodes (12): deleteMembership(), fetchMemberPlans(), fetchMemberships(), setMembershipActive(), updateMemberPlanStatus(), MembershipsPage(), handleDelete(), handlePlanStatusChange() (+4 more)

### Community 71 - "UseCase"
Cohesion: 0.15
Nodes (10): UseCase, NewUseCase(), Log, NewNotificationRepository(), NewNotificationHandler(), NotificationRepository, Provider, Repository (+2 more)

### Community 72 - "email_provider.go"
Cohesion: 0.16
Nodes (9): boolEnv(), buildNotificationProvider(), durationSecondsEnv(), intEnv(), logFallback(), NewEmailProvider(), NewLogProvider(), EmailProvider (+1 more)

### Community 73 - "writeInvalidPayload"
Cohesion: 0.20
Nodes (7): writeConfigurationError(), writeBadRequestError(), writeInvalidPayload(), writeTenantError(), writeWaitingListError(), TenantHandler, WaitingListHandler

### Community 74 - "booking/src/lib/apiError.ts"
Cohesion: 0.29
Nodes (8): ApiError, ApiErrorPayload, buildApiError(), parseApiResponse(), API_ERROR_MESSAGES, FIELD_ERROR_MESSAGES, getApiErrorMessage(), translateFieldErrors()

### Community 76 - ".UpdateAccountRole"
Cohesion: 0.48
Nodes (3): canManageTenantAccount(), isPrivilegedAdminRole(), normalizeTenantScope()

### Community 77 - "writeNotFoundError"
Cohesion: 0.27
Nodes (4): writeCustomerError(), writeNotificationError(), writeNotFoundError(), CustomerHandler

### Community 78 - "Container"
Cohesion: 0.05
Nodes (50): NewUseCase(), NewUseCase(), TestCalendarUseCase(), NewUseCase(), TestCommissionUseCase(), NewUseCase(), NewUseCase(), NewUseCase() (+42 more)

### Community 80 - "PortalPage"
Cohesion: 0.39
Nodes (8): portalCancelAppointment(), portalLookup(), portalRescheduleAppointment(), PortalPage(), handleCancel(), handleConfirmReschedule(), handleLookup(), refresh()

### Community 85 - "ical.go"
Cohesion: 0.83
Nodes (3): BuildICalCalendar(), formatUTC(), sanitizeLine()

### Community 87 - "body_limit.go"
Cohesion: 0.83
Nodes (3): resolveMaxRequestBodyBytes(), shouldLimitRequestBody(), withRequestBodyLimit()

### Community 90 - "005_resources.sql"
Cohesion: 0.40
Nodes (3): catalog.resources, idx_appointments_resource, idx_resources_tenant

### Community 91 - "007_jobs_and_calendar_sync.sql"
Cohesion: 0.47
Nodes (4): catalog.calendar_connections, catalog.jobs, idx_calendar_connections_tenant, idx_catalog_jobs_poll

### Community 92 - "booking/src/shared/components/Button.tsx"
Cohesion: 0.33
Nodes (5): ButtonProps, ButtonSize, ButtonVariant, sizeClasses, variantClasses

### Community 93 - "widget_handler.go"
Cohesion: 0.60
Nodes (3): NewWidgetHandler(), resolveWidgetBookingOrigin(), WidgetHandler

### Community 94 - "catalog.reviews"
Cohesion: 0.60
Nodes (3): catalog.reviews, idx_reviews_tenant_created, idx_reviews_tenant_status

## Knowledge Gaps
- **256 isolated node(s):** `graphify-mcp`, `name`, `private`, `version`, `type` (+251 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 492 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **50 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Container` connect `Container` to `Resource`, `RouterConfig`, `Professional`, `UseCase`, `UseCase`, `Payment`, `Tenant`, `rate_limit.go`, `container.go`, `Customer`, `Entry`, `database/sql.DB`, `Connection`, `Queue`, `UseCase`, `Service`, `UseCase`, `UseCase`, `UseCase`, `ReminderWorker`, `UseCase`, `email_provider.go`?**
  _High betweenness centrality (0.050) - this node is a cross-community bridge._
- **Why does `UseCase` connect `UseCase` to `writeAPIError`, `buildAuthComponents`, `AuthRepository`, `RouterConfig`, `net/http.ResponseWriter`, `Account`, `net/http.Request`, `.registerWithRole`, `.UpdateAccountRole`, `Container`, `time.Duration`, `auth/usecase.go`, `Configuration`, `rate_limit.go`, `AppointmentHandler`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **Why does `repositorySet` connect `database/sql.DB` to `Resource`, `RouterConfig`, `Professional`, `go_pkg_time`, `UseCase`, `UseCase`, `Payment`, `Tenant`, `Customer`, `Membership`, `Entry`, `Appointment`, `Connection`, `Queue`, `GiftCard`, `Account`, `UseCase`, `Review`, `PackageCredit`, `MemberPlan`, `UseCase`, `Container`?**
  _High betweenness centrality (0.022) - this node is a cross-community bridge._
- **Are the 119 inferred relationships involving `writeJSON()` (e.g. with `health()` and `readiness()`) actually correct?**
  _`writeJSON()` has 119 INFERRED edges - model-reasoned connections that need verification._
- **Are the 95 inferred relationships involving `resolveRequestTenantID()` (e.g. with `.Analytics()` and `.CheckIn()`) actually correct?**
  _`resolveRequestTenantID()` has 95 INFERRED edges - model-reasoned connections that need verification._
- **What connects `graphify-mcp`, `name`, `private` to the rest of the system?**
  _256 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `adminAuth.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.05480277024992472 - nodes in this community are weakly interconnected._
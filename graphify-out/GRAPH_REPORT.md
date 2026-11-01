# Graph Report - erpagendamentos  (2026-10-31)

## Corpus Check
- 250 files · ~106,161 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 7 file(s) not represented in the graph (top: (none) 3, .css 3, .mdc 1)

## Summary
- 2559 nodes · 8219 edges · 121 communities (90 shown, 31 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 614 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `320efef2`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Section
- adminAuth.ts
- booking/usecase.go
- ref_react
- Payment
- adminApi.ts
- Button
- Tenant
- RouterConfig
- UseCase
- response.go
- Customer
- net/http.ResponseWriter
- Entry
- adminFetch
- PackagesPage.tsx
- UseCase
- Repository
- testing.T
- Review
- UseCase
- booking/src/pages/WaitingListPage.tsx
- MemberPlan
- BookingPage
- Account
- container.go
- Professional
- GiftCard
- landing/src/App.tsx
- UseCase
- Membership
- UseCase
- ServiceRepository
- 002_features.sql
- entities.go
- AuthRepository
- Resource
- database/sql.DB
- PackageCredit
- PortalPage.tsx
- AdminRoutes.tsx
- UseCase
- writePaymentError
- 001_baseline.sql
- observability.go
- net/http.Request
- admin/package.json
- Appointment
- templates.go
- UseCase
- BookingPage.tsx
- compilerOptions
- AppointmentRepository
- Transaction
- compilerOptions
- compilerOptions
- 006_erp_finance_stock_commissions.sql
- Service
- writeJSON
- email_provider.go
- sync.Mutex
- Configuration
- requestClientIP
- go_pkg_strings
- DeviceSession
- AppointmentHandler
- Connection
- Analytics
- ResourceRepository
- ERP Agendamentos
- booking/src/lib/apiError.ts
- ProfessionalsPage.tsx
- ReminderWorker
- Claims
- Container
- .GetByID
- PortalPage
- ServicesPage.tsx
- UseCase
- booking/package.json
- landing/package.json
- .UpdateAccountRole
- SecretCipher
- time.Duration
- writeResourceError
- 005_resources.sql
- booking/src/shared/components/Button.tsx
- ReviewRepository
- Job
- widget_handler.go
- catalog.reviews
- EventBroker
- Repository
- Queue
- Repository
- memoryRepositories
- UseCase
- UseCase
- erpagendamentos/backend
- writeProfessionalError
- writeServiceError
- ResourcesPage
- writeCustomerError
- Summary
- 007_jobs_and_calendar_sync.sql
- graphify
- rules/graphify.md
- workflows/graphify.md

## God Nodes (most connected - your core abstractions)
1. `writeJSON()` - 120 edges
2. `resolveRequestTenantID()` - 99 edges
3. `adminFetch()` - 68 edges
4. `Button` - 64 edges
5. `writeInvalidPayload()` - 64 edges
6. `useToast()` - 62 edges
7. `UseCase` - 62 edges
8. `invalidateAdminGetCache()` - 59 edges
9. `writeInternalError()` - 56 edges
10. `Appointment` - 49 edges

## Surprising Connections (you probably didn't know these)
- `newTestUseCase()` --calls--> `NewUseCase()`  [INFERRED]
  backend/internal/application/booking/usecase_test.go → backend/internal/application/booking/usecase.go
- `buildAuthComponents()` --calls--> `isProductionEnv()`  [INFERRED]
  backend/internal/bootstrap/auth_components.go → backend/internal/bootstrap/container.go
- `buildAuthComponents()` --calls--> `seedAdminAccount()`  [INFERRED]
  backend/internal/bootstrap/auth_components.go → backend/internal/bootstrap/container.go
- `NewContainer()` --calls--> `buildAuthComponents()`  [INFERRED]
  backend/internal/bootstrap/container.go → backend/internal/bootstrap/auth_components.go
- `NewContainer()` --calls--> `loadEnvFile()`  [INFERRED]
  backend/internal/bootstrap/container.go → backend/internal/bootstrap/env.go

## Import Cycles
- None detected.

## Communities (121 total, 31 thin omitted)

### Community 0 - "Section"
Cohesion: 0.07
Nodes (51): deleteCustomer(), downloadAdminCsv(), fetchCurrentTenant(), fetchCustomers(), fetchPayments(), fetchReviews(), fetchWaitingList(), markPaymentPaid() (+43 more)

### Community 1 - "adminAuth.ts"
Cohesion: 0.06
Nodes (57): AdminShell(), AdminShellProps, App(), handleLogout(), AdminNavigationItem, navigationGroups, navigationItems, AppProviders() (+49 more)

### Community 3 - "ref_react"
Cohesion: 0.08
Nodes (52): checkInAppointment(), deleteAppointment(), deleteTenant(), fetchAppointmentAnalytics(), fetchAppointments(), fetchAppointmentSummary(), fetchTenants(), invalidateMutations() (+44 more)

### Community 4 - "Payment"
Cohesion: 0.09
Nodes (13): UseCase, NewUseCase(), normalizeStatus(), Payment, NewPaymentRepository(), scanPaymentRow(), scanPaymentRows(), NewPaymentHandler() (+5 more)

### Community 5 - "adminApi.ts"
Cohesion: 0.06
Nodes (71): AppointmentFilters, createCommissionRule(), createProduct(), createStockMovement(), deleteCommissionRule(), deleteProduct(), fetchCommissionRecords(), fetchCommissionRules() (+63 more)

### Community 6 - "Button"
Cohesion: 0.09
Nodes (59): AppointmentFormModal(), computeEndsAt(), handleSubmit(), AppointmentFormModalProps, STATUS_OPTIONS, createAppointment(), createCustomer(), createResource() (+51 more)

### Community 8 - "Tenant"
Cohesion: 0.09
Nodes (13): UseCase, NewUseCase(), ensureDefaultTenant(), Tenant, isTenantIdentityConflict(), NewTenantRepository(), scanTenantRow(), scanTenantRows() (+5 more)

### Community 9 - "RouterConfig"
Cohesion: 0.07
Nodes (65): bearerAuthAllowed(), hasRole(), isAdminRole(), isPrivilegedAdminRole(), requestWithAuthClaims(), requireAdmin(), requireMasterAdmin(), requireOwnerOrAdmin() (+57 more)

### Community 10 - "UseCase"
Cohesion: 0.06
Nodes (12): UseCase, NewUseCase(), TestStockUseCase(), Product, Repository, ServiceSupply, StockMovement, NewStockRepository() (+4 more)

### Community 11 - "response.go"
Cohesion: 0.22
Nodes (4): writeCSV(), writeRawJSON(), errorResponse, ExportHandler

### Community 12 - "Customer"
Cohesion: 0.12
Nodes (10): Customer, NewCustomerRepository(), onlyDigits(), NewCustomerRepository(), scanCustomerRow(), scanCustomerRows(), Repository, CustomerRepository (+2 more)

### Community 13 - "net/http.ResponseWriter"
Cohesion: 0.14
Nodes (10): writeBookingError(), writeNotificationError(), writeConflictError(), writeNotFoundError(), writeReviewError(), writeTenantError(), BookingHandler, responseRecorder (+2 more)

### Community 14 - "Entry"
Cohesion: 0.11
Nodes (12): UseCase, NewUseCase(), Entry, NewWaitingListRepository(), NewWaitingListRepository(), scanWaitingRow(), scanWaitingRows(), NewWaitingListHandler() (+4 more)

### Community 15 - "adminFetch"
Cohesion: 0.09
Nodes (48): connectCalendar(), createFinancialTransaction(), createGiftCard(), createMemberPlan(), createMembership(), createPackage(), createPackageCredit(), createPayment() (+40 more)

### Community 16 - "PackagesPage.tsx"
Cohesion: 0.11
Nodes (34): createAdminUser(), deleteAdminUser(), deleteGiftCard(), deletePackage(), fetchAdminUsers(), fetchGiftCards(), fetchPackageCredits(), fetchPackages() (+26 more)

### Community 17 - "UseCase"
Cohesion: 0.10
Nodes (13): MFAStatus, SecretCipher, Session, TokenManager, generateTOTPSecret(), totpCode(), totpURI(), validateTOTP() (+5 more)

### Community 19 - "testing.T"
Cohesion: 0.19
Nodes (16): newAppointment(), newResourceUseCase(), TestCreateAutoAssignsAvailableResource(), TestCreateRejectsBusyExplicitResource(), TestCreateRejectsUnknownResource(), TestUpdateAutoReassignsFreedResource(), NewUseCase(), newTestUseCase() (+8 more)

### Community 20 - "Review"
Cohesion: 0.16
Nodes (8): UseCase, NewUseCase(), Review, NewReviewRepository(), NewAppointmentHandler(), NewReviewHandler(), ReviewRepository, Repository

### Community 21 - "UseCase"
Cohesion: 0.16
Nodes (6): RecurrenceInput, UseCase, parseRFC3339(), validStatusChange(), NewCalendarHandler(), CalendarHandler

### Community 22 - "booking/src/pages/WaitingListPage.tsx"
Cohesion: 0.14
Nodes (24): App(), fetchPublicReview(), signUpWaitingList(), submitPublicReview(), formatMoney(), initWidgetBridge(), DetailsStep(), ReviewPage() (+16 more)

### Community 23 - "MemberPlan"
Cohesion: 0.14
Nodes (5): NewUseCase(), MemberPlan, scanMemberPlans(), MemberPlanRepository, MemberPlanRepository

### Community 24 - "BookingPage"
Cohesion: 0.14
Nodes (18): cancelPublicAppointment(), createPublicBooking(), fetchAvailability(), fetchPublicAppointment(), fetchPublicInfo(), fetchPublicProfessionals(), fetchPublicServices(), publicFetch() (+10 more)

### Community 25 - "Account"
Cohesion: 0.14
Nodes (6): Account, PasswordResetToken, UpdateProfileInput, cloneAuthAccount(), resetTokenExpired(), AuthRepository

### Community 26 - "container.go"
Cohesion: 0.08
Nodes (9): newLazyAppHandler(), main(), resolveDurationEnv(), resolvePositiveIntEnv(), deploymentEnvironment(), isProductionEnv(), newMemoryRateLimiter(), resolvePositiveEnvInt() (+1 more)

### Community 27 - "Professional"
Cohesion: 0.10
Nodes (12): UseCase, NewUseCase(), Professional, NewProfessionalRepository(), NewProfessionalRepository(), scanProfessionalRow(), scanProfessionalRows(), NewProfessionalHandler() (+4 more)

### Community 28 - "GiftCard"
Cohesion: 0.11
Nodes (5): generateGiftCardCode(), GiftCard, GiftCardRepository, GiftCardRepository, GiftCardRepository

### Community 29 - "landing/src/App.tsx"
Cohesion: 0.16
Nodes (14): App(), CtaSection(), FEATURES, FeaturesSection(), Footer(), Header(), NAV_ITEMS, HeroSection() (+6 more)

### Community 30 - "UseCase"
Cohesion: 0.07
Nodes (12): UseCase, NewUseCase(), TestCommissionUseCase(), CommissionRecord, CommissionRule, Repository, Summary, NewCommissionRepository() (+4 more)

### Community 31 - "Membership"
Cohesion: 0.12
Nodes (4): Membership, MembershipRepository, MembershipRepository, MembershipRepository

### Community 32 - "UseCase"
Cohesion: 0.11
Nodes (6): UseCase, Package, NewPackageRepository(), NewCommerceHandler(), PackageRepository, PackageRepository

### Community 33 - "ServiceRepository"
Cohesion: 0.23
Nodes (5): NewServiceRepository(), scanServiceRow(), scanServiceRows(), ServiceRepository, serviceScanner

### Community 34 - "002_features.sql"
Cohesion: 0.14
Nodes (20): catalog.gift_cards, catalog.member_plans, catalog.memberships, catalog.notification_logs, catalog.package_credits, catalog.packages, catalog.payments, catalog.service_professionals (+12 more)

### Community 35 - "entities.go"
Cohesion: 0.15
Nodes (14): ChangePasswordInput, LoginInput, PasswordResetCodeConfirmInput, PasswordResetConfirmInput, PasswordResetRequestInput, RegisterInput, ValidationError, hashResetToken() (+6 more)

### Community 37 - "Resource"
Cohesion: 0.14
Nodes (6): UseCase, NewUseCase(), Resource, NewResourceHandler(), ResourceRepository, Repository

### Community 38 - "database/sql.DB"
Cohesion: 0.14
Nodes (17): applyDatabaseSQL(), applyMigrationFile(), applySQLDirectory(), applySQLFile(), initializeRepositories(), postgresRepositories(), NewAuthRepository(), NewGiftCardRepository() (+9 more)

### Community 39 - "PackageCredit"
Cohesion: 0.14
Nodes (5): PackageCredit, scanPackageCredits(), PackageCreditRepository, PackageCreditRepository, PackageCreditRepository

### Community 40 - "PortalPage.tsx"
Cohesion: 0.17
Nodes (24): PortalAppointment, PortalResult, Slot, formatDateLong(), formatDateTime(), formatTime(), fromDateKey(), nextDateKeys() (+16 more)

### Community 41 - "AdminRoutes.tsx"
Cohesion: 0.15
Nodes (25): AccessDeniedPage(), AdminRouteFallback(), AdminRoutes(), AdminRoutesProps, AgendaPage, AnalyticsPage, AppointmentsPage, CommissionsPage (+17 more)

### Community 42 - "UseCase"
Cohesion: 0.16
Nodes (12): appointmentDurationMinutes(), formatAppointmentForTemplate(), generateManagementCode(), UseCase, hashManagementCode(), parseRFC3339(), TestManagementCodeNormalization(), PublicBookingOrigin() (+4 more)

### Community 45 - "001_baseline.sql"
Cohesion: 0.17
Nodes (21): catalog.appointments, catalog.auth_rate_limits, catalog.auth_sessions, catalog.configurations, catalog.customers, catalog.password_reset_tokens, catalog.professionals, catalog.services (+13 more)

### Community 46 - "observability.go"
Cohesion: 0.10
Nodes (16): buildPaymentProvider(), NewMockProvider(), mapStripeStatus(), NewStripeProvider(), isSensitiveLogKey(), redactSensitiveText(), sanitizeLogString(), sanitizeLogValue() (+8 more)

### Community 47 - "net/http.Request"
Cohesion: 0.09
Nodes (26): resolveActorTenantID(), resolveAuthClaims(), resolveTenantID(), authTokenFromCookie(), buildPasswordResetURL(), AuthHandler, secureTemporaryPassword(), writeAdminAuthError() (+18 more)

### Community 48 - "admin/package.json"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+24 more)

### Community 49 - "Appointment"
Cohesion: 0.10
Nodes (7): Repository, Appointment, StatusCount, parseAppointmentBound(), parseAppointmentInterval(), BookingResult, AppointmentRepository

### Community 50 - "templates.go"
Cohesion: 0.29
Nodes (15): greeting(), RenderCancellation(), RenderConfirmation(), RenderReminder(), RenderReschedule(), sampleData(), TestGreetingWithoutName(), TestRenderCancellation() (+7 more)

### Community 51 - "UseCase"
Cohesion: 0.15
Nodes (10): UseCase, NewUseCase(), Log, NewNotificationRepository(), NewNotificationHandler(), NotificationRepository, Provider, Repository (+2 more)

### Community 52 - "BookingPage.tsx"
Cohesion: 0.13
Nodes (17): backendOrigin, buildApiUrl(), AvailabilityResponse, BookingResult, Professional, PublicAppointment, PublicConfiguration, PublicInfo (+9 more)

### Community 53 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 54 - "AppointmentRepository"
Cohesion: 0.16
Nodes (7): appointmentSelectColumns(), NewAppointmentRepository(), nullIfEmpty(), scanAppointmentRow(), scanAppointmentRows(), AppointmentRepository, appointmentScanner

### Community 55 - "Transaction"
Cohesion: 0.15
Nodes (5): Repository, Transaction, NewFinanceRepository(), FinanceRepository, FinanceRepository

### Community 56 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 57 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 58 - "006_erp_finance_stock_commissions.sql"
Cohesion: 0.21
Nodes (12): catalog.commission_records, catalog.commission_rules, catalog.financial_transactions, catalog.products, catalog.service_supplies, catalog.stock_movements, idx_catalog_commission_records_prof, idx_catalog_commission_rules_tenant (+4 more)

### Community 59 - "Service"
Cohesion: 0.12
Nodes (5): UseCase, Service, NewServiceHandler(), ServiceRepository, Repository

### Community 60 - "writeJSON"
Cohesion: 0.11
Nodes (12): resolveRequestTenantID(), writeCommerceError(), writeInternalError(), writeInvalidPayload(), writeJSON(), writeWaitingListError(), CalendarSyncHandler, CommerceHandler (+4 more)

### Community 61 - "email_provider.go"
Cohesion: 0.21
Nodes (7): buildNotificationProvider(), intEnv(), logFallback(), NewEmailProvider(), NewLogProvider(), EmailProvider, LogProvider

### Community 62 - "sync.Mutex"
Cohesion: 0.21
Nodes (10): countOverlappingSlots(), isClosedDay(), localMidnight(), overlapsAnyBusy(), parseClock(), ProfessionalSlots, Slot, rateLimitEntry (+2 more)

### Community 63 - "Configuration"
Cohesion: 0.21
Nodes (12): workingHoursFor(), Configuration, cloneConfiguration(), defaultMemoryConfiguration(), NewConfigurationRepository(), normalizeConfigurationTenantID(), NewConfigurationRepository(), normalizePostgresTenantID() (+4 more)

### Community 64 - "requestClientIP"
Cohesion: 0.16
Nodes (12): browserSessionPayload(), exposeAuthTokenForTests(), AuthHandler, health(), readiness(), logError(), forwardedClientIP(), isTrustedProxy() (+4 more)

### Community 65 - "go_pkg_strings"
Cohesion: 0.19
Nodes (4): BuildICalCalendar(), formatUTC(), sanitizeLine(), CategoryTotal

### Community 67 - "AppointmentHandler"
Cohesion: 0.23
Nodes (4): writeAppointmentError(), writeConfigurationError(), writeBadRequestError(), AppointmentHandler

### Community 68 - "Connection"
Cohesion: 0.20
Nodes (4): Connection, NewCalendarRepository(), CalendarRepository, CalendarRepository

### Community 69 - "Analytics"
Cohesion: 0.33
Nodes (4): DayMetric, ProfessionalMetric, ServiceMetric, Analytics

### Community 70 - "ResourceRepository"
Cohesion: 0.27
Nodes (4): NewResourceRepository(), scanResourceRow(), ResourceRepository, resourceScanner

### Community 71 - "ERP Agendamentos"
Cohesion: 0.14
Nodes (13): CI, Endpoints principais, ERP Agendamentos, Estrutura, Migrations, Pré-requisitos, Rodando com Docker (produção local), Rodando em desenvolvimento (sem banco) (+5 more)

### Community 72 - "booking/src/lib/apiError.ts"
Cohesion: 0.29
Nodes (8): ApiError, ApiErrorPayload, buildApiError(), parseApiResponse(), API_ERROR_MESSAGES, FIELD_ERROR_MESSAGES, getApiErrorMessage(), translateFieldErrors()

### Community 73 - "ProfessionalsPage.tsx"
Cohesion: 0.31
Nodes (10): createProfessional(), deleteProfessional(), setProfessionalActive(), updateProfessional(), COLOR_OPTIONS, handleSubmit(), ProfessionalsPage(), handleDelete() (+2 more)

### Community 74 - "ReminderWorker"
Cohesion: 0.21
Nodes (6): MaintenanceWorker, MaintenanceConfig, NewMaintenanceWorker(), formatAppointmentLocal(), NewReminderWorker(), ReminderWorker

### Community 75 - "Claims"
Cohesion: 0.31
Nodes (5): isProductionAuthEnvironment(), Claims, decodeSegment(), encodeSegment(), JWTManager

### Community 76 - "Container"
Cohesion: 0.11
Nodes (14): lazyAppHandler, writeLazyJSON(), NewUseCase(), UseCase, NewUseCase(), NewUseCase(), NewUseCase(), Container (+6 more)

### Community 78 - "PortalPage"
Cohesion: 0.39
Nodes (8): portalCancelAppointment(), portalLookup(), portalRescheduleAppointment(), PortalPage(), handleCancel(), handleConfirmReschedule(), handleLookup(), refresh()

### Community 79 - "ServicesPage.tsx"
Cohesion: 0.31
Nodes (10): createService(), deleteService(), setServiceActive(), updateService(), COLOR_OPTIONS, handleSubmit(), ServicesPage(), handleDelete() (+2 more)

### Community 80 - "UseCase"
Cohesion: 0.20
Nodes (5): UseCase, NewUseCase(), TestFinanceUseCase(), NewFinanceRepository(), NewFinanceHandler()

### Community 81 - "booking/package.json"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+24 more)

### Community 82 - "landing/package.json"
Cohesion: 0.06
Nodes (30): dependencies, lucide-react, react, react-dom, devDependencies, tailwindcss, @tailwindcss/vite, @types/react (+22 more)

### Community 83 - ".UpdateAccountRole"
Cohesion: 0.48
Nodes (3): canManageTenantAccount(), isPrivilegedAdminRole(), normalizeTenantScope()

### Community 86 - "time.Duration"
Cohesion: 0.12
Nodes (8): buildAuthComponents(), boolEnv(), durationSecondsEnv(), CleanupSummary, NewBcryptHasher(), NewJWTManager(), authComponents, BcryptHasher

### Community 88 - "005_resources.sql"
Cohesion: 0.40
Nodes (3): catalog.resources, idx_appointments_resource, idx_resources_tenant

### Community 89 - "booking/src/shared/components/Button.tsx"
Cohesion: 0.33
Nodes (5): ButtonProps, ButtonSize, ButtonVariant, sizeClasses, variantClasses

### Community 90 - "ReviewRepository"
Cohesion: 0.29
Nodes (5): NewReviewRepository(), scanReviewRow(), scanReviewRows(), ReviewRepository, reviewScanner

### Community 92 - "widget_handler.go"
Cohesion: 0.60
Nodes (3): NewWidgetHandler(), resolveWidgetBookingOrigin(), WidgetHandler

### Community 93 - "catalog.reviews"
Cohesion: 0.60
Nodes (3): catalog.reviews, idx_reviews_tenant_created, idx_reviews_tenant_status

### Community 94 - "EventBroker"
Cohesion: 0.27
Nodes (5): NewBookingHandler(), EventBroker, RegisterMetricsBroker(), Event, MetricsRegistry

### Community 95 - "Repository"
Cohesion: 0.22
Nodes (4): NewUseCase(), TestCalendarUseCase(), Repository, NewCalendarRepository()

### Community 97 - "Repository"
Cohesion: 0.22
Nodes (4): NewQueue(), TestJobQueue(), Repository, NewJobRepository()

### Community 98 - "memoryRepositories"
Cohesion: 0.22
Nodes (9): memoryRepositories(), NewAuthRepository(), NewGiftCardRepository(), NewMemberPlanRepository(), NewMembershipRepository(), NewPackageCreditRepository(), NewPaymentRepository(), NewServiceRepository() (+1 more)

### Community 99 - "UseCase"
Cohesion: 0.25
Nodes (3): UseCase, SyncEvent, NewCalendarSyncHandler()

### Community 100 - "UseCase"
Cohesion: 0.29
Nodes (3): UseCase, NewCustomerHandler(), NewExportHandler()

### Community 110 - "ResourcesPage"
Cohesion: 0.52
Nodes (7): deleteResource(), fetchResources(), setResourceActive(), ResourcesPage(), handleDelete(), handleToggleActive(), loadResources()

### Community 113 - "007_jobs_and_calendar_sync.sql"
Cohesion: 0.47
Nodes (4): catalog.calendar_connections, catalog.jobs, idx_calendar_connections_tenant, idx_catalog_jobs_poll

## Knowledge Gaps
- **253 isolated node(s):** `graphify-mcp`, `name`, `private`, `version`, `type` (+248 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 459 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **31 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Container` connect `Container` to `Payment`, `Tenant`, `RouterConfig`, `UseCase`, `Customer`, `Entry`, `UseCase`, `Review`, `UseCase`, `container.go`, `Professional`, `UseCase`, `UseCase`, `Resource`, `database/sql.DB`, `UseCase`, `UseCase`, `Service`, `ReminderWorker`, `UseCase`, `time.Duration`, `EventBroker`, `Queue`, `UseCase`, `UseCase`?**
  _High betweenness centrality (0.049) - this node is a cross-community bridge._
- **Why does `UseCase` connect `UseCase` to `requestClientIP`, `DeviceSession`, `entities.go`, `RouterConfig`, `Claims`, `Container`, `net/http.Request`, `Repository`, `.UpdateAccountRole`, `time.Duration`, `Account`, `container.go`, `writeJSON`, `sync.Mutex`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **Why does `repositorySet` connect `database/sql.DB` to `Payment`, `go_pkg_time`, `Tenant`, `UseCase`, `Customer`, `Entry`, `Repository`, `Review`, `MemberPlan`, `Professional`, `GiftCard`, `UseCase`, `Membership`, `UseCase`, `Resource`, `PackageCredit`, `Appointment`, `UseCase`, `Transaction`, `Service`, `Container`, `Repository`, `Repository`, `memoryRepositories`?**
  _High betweenness centrality (0.032) - this node is a cross-community bridge._
- **Are the 116 inferred relationships involving `writeJSON()` (e.g. with `health()` and `readiness()`) actually correct?**
  _`writeJSON()` has 116 INFERRED edges - model-reasoned connections that need verification._
- **Are the 94 inferred relationships involving `resolveRequestTenantID()` (e.g. with `.Analytics()` and `.CheckIn()`) actually correct?**
  _`resolveRequestTenantID()` has 94 INFERRED edges - model-reasoned connections that need verification._
- **What connects `graphify-mcp`, `name`, `private` to the rest of the system?**
  _253 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Section` be split into smaller, more focused modules?**
  _Cohesion score 0.07175141242937853 - nodes in this community are weakly interconnected._
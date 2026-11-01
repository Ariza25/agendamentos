# Graph Report - erpagendamentos  (2026-11-01)

## Corpus Check
- 271 files · ~125,290 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 15 file(s) not represented in the graph (top: (none) 6, .conf 5, .css 3)

## Summary
- 2800 nodes · 9019 edges · 159 communities (92 shown, 67 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 661 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9f15177b`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- adminAuth.ts
- TenantsPage.tsx
- Payment
- writeInvalidPayload
- adminFetch
- adminApi.ts
- Button
- AppointmentsPage.tsx
- writeAPIError
- UseCase
- checkout/usecase.go
- container.go
- Tenant
- net/http.Request
- auth_cookie.go
- RouterConfig
- auth/usecase.go
- Container
- CashRegister
- UseCase
- Customer
- database/sql.DB
- Queue
- Professional
- Entry
- admin/package.json
- booking/package.json
- go_pkg_strings
- Service
- UseCase
- .registerWithRole
- Connection
- FinancePage.tsx
- landing/package.json
- booking/src/pages/WaitingListPage.tsx
- UseCase
- GiftCard
- Repository
- Account
- PortalPage.tsx
- AdminRoutes.tsx
- UseCase
- Package
- PackageCredit
- bookingApi.ts
- landing/src/App.tsx
- PackagesPage
- Resource
- repositorySet
- 002_features.sql
- AuthRepository
- UseCase
- ConfigurationPage.tsx
- .GetByID
- UseCase
- 001_baseline.sql
- BookingPage.tsx
- CommissionsPage.tsx
- testing.T
- UseCase
- writeJSON
- writeNotFoundError
- Appointment
- templates.go
- AppointmentRepository
- compilerOptions
- Configuration
- net/http.ResponseWriter
- compilerOptions
- compilerOptions
- ReminderWorker
- 006_erp_finance_stock_commissions.sql
- AppointmentHandler
- time.Duration
- Claims
- net/http.Handler
- ERP Agendamentos
- MembershipsPage.tsx
- DeviceSession
- EventBroker
- PaymentsPage.tsx
- Analytics
- sync.Mutex
- memoryRepositories
- observability.go
- booking/src/lib/apiError.ts
- catalog.commission_payouts
- writeProfessionalError
- ResourceRepository
- writeResourceError
- PortalPage
- .UpdateAccountRole
- writeServiceError
- ical.go
- bcrypt_hasher.go
- 005_resources.sql
- 007_jobs_and_calendar_sync.sql
- 010_branding_pos_and_cash_registers.sql
- booking/src/shared/components/Button.tsx
- SecretCipher
- widget_handler.go
- catalog.reviews
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
- `newResourceUseCase()` --calls--> `NewUseCase()`  [INFERRED]
  backend/internal/application/appointment/resource_test.go → backend/internal/application/appointment/usecase.go
- `newTestUseCase()` --calls--> `NewUseCase()`  [INFERRED]
  backend/internal/application/booking/usecase_test.go → backend/internal/application/booking/usecase.go
- `TestManagementCodeNormalization()` --calls--> `hashManagementCode()`  [INFERRED]
  backend/internal/application/booking/usecase_test.go → backend/internal/application/booking/usecase.go
- `buildAuthComponents()` --calls--> `boolEnv()`  [INFERRED]
  backend/internal/bootstrap/auth_components.go → backend/internal/bootstrap/container.go
- `buildAuthComponents()` --calls--> `isProductionEnv()`  [INFERRED]
  backend/internal/bootstrap/auth_components.go → backend/internal/bootstrap/container.go

## Import Cycles
- None detected.

## Communities (159 total, 67 thin omitted)

### Community 0 - "adminAuth.ts"
Cohesion: 0.07
Nodes (52): App(), handleLogout(), AppProviders(), AppProvidersProps, ResetPasswordPage(), handleSubmit(), AppRouter(), ADMIN_ROLES (+44 more)

### Community 1 - "TenantsPage.tsx"
Cohesion: 0.08
Nodes (48): deleteGiftCard(), deleteTenant(), fetchCurrentTenant(), fetchGiftCards(), fetchNotifications(), fetchReviews(), fetchTenants(), fetchWaitingList() (+40 more)

### Community 2 - "Payment"
Cohesion: 0.05
Nodes (23): UseCase, NewUseCase(), normalizeStatus(), buildPaymentProvider(), IntentResult, Payment, NewMercadoPagoProvider(), NewMockProvider() (+15 more)

### Community 3 - "writeInvalidPayload"
Cohesion: 0.18
Nodes (7): writeCommerceError(), writeConfigurationError(), writeBadRequestError(), writeInvalidPayload(), writeWaitingListError(), CommerceHandler, WaitingListHandler

### Community 4 - "adminFetch"
Cohesion: 0.06
Nodes (68): connectCalendar(), createAdminUser(), createCommissionPayout(), createGiftCard(), createMemberPlan(), createMembership(), createPackage(), createPackageCredit() (+60 more)

### Community 5 - "adminApi.ts"
Cohesion: 0.06
Nodes (65): CheckoutComandaModalProps, SelectedProduct, AppointmentFilters, checkoutAppointment(), createProduct(), createStockMovement(), deleteProduct(), fetchProducts() (+57 more)

### Community 6 - "Button"
Cohesion: 0.11
Nodes (57): CheckoutComandaModal(), AdminAccount, CommissionPayoutModal(), CustomerFormModal(), GiftCardFormModal(), LoginAdminPageProps, CreatePlanModal(), MembershipFormModal() (+49 more)

### Community 7 - "AppointmentsPage.tsx"
Cohesion: 0.11
Nodes (41): AppointmentFormModal(), computeEndsAt(), handleSubmit(), AppointmentFormModalProps, STATUS_OPTIONS, checkInAppointment(), createAppointment(), createCustomer() (+33 more)

### Community 8 - "writeAPIError"
Cohesion: 0.08
Nodes (34): bearerAuthAllowed(), hasRole(), isAdminRole(), isPrivilegedAdminRole(), requestWithAuthClaims(), requireMasterAdmin(), requireOwnerOrAdmin(), requirePrivilegedAdmin() (+26 more)

### Community 9 - "UseCase"
Cohesion: 0.07
Nodes (10): UseCase, NewUseCase(), TestStockUseCase(), Product, Repository, ServiceSupply, StockMovement, NewStockHandler() (+2 more)

### Community 10 - "checkout/usecase.go"
Cohesion: 0.29
Nodes (3): CheckoutInput, CheckoutResult, ProductItem

### Community 11 - "container.go"
Cohesion: 0.11
Nodes (4): deploymentEnvironment(), isProductionEnv(), CategoryTotal, Type

### Community 12 - "Tenant"
Cohesion: 0.08
Nodes (16): UseCase, NewUseCase(), ensureDefaultTenant(), Tenant, UpdateSubscriptionInput, defaultModules(), isTenantIdentityConflict(), defaultModules() (+8 more)

### Community 13 - "net/http.Request"
Cohesion: 0.10
Nodes (8): resolveRequestTenantID(), writeCSV(), writeInternalError(), CalendarSyncHandler, CommissionHandler, ExportHandler, FinanceHandler, StockHandler

### Community 14 - "auth_cookie.go"
Cohesion: 0.26
Nodes (17): authCookieDomain(), authCookiePartitioned(), authCookieSameSite(), authCookieSecure(), authTokenFromCookie(), clearAuthCookie(), csrfTokenFromCookie(), isAuthCookieProductionEnv() (+9 more)

### Community 16 - "RouterConfig"
Cohesion: 0.25
Nodes (27): requireAdmin(), NewRouter(), registerAppointmentRoutes(), registerConfigurationRoutes(), registerCustomerRoutes(), registerProfessionalRoutes(), registerResourceRoutes(), registerServiceRoutes() (+19 more)

### Community 17 - "auth/usecase.go"
Cohesion: 0.11
Nodes (3): newMemoryRateLimiter(), resolvePositiveEnvInt(), tokenPayload

### Community 18 - "Container"
Cohesion: 0.13
Nodes (16): UseCase, NewUseCase(), UseCase, NewUseCase(), NewUseCase(), Container, NewContainer(), loadEnvFile() (+8 more)

### Community 19 - "CashRegister"
Cohesion: 0.10
Nodes (9): UseCase, NewUseCase(), CashRegister, CashRegisterEntry, Repository, ShiftSummary, NewCashRegisterHandler(), CashRegisterRepository (+1 more)

### Community 20 - "UseCase"
Cohesion: 0.10
Nodes (6): UseCase, Membership, NewMembershipRepository(), NewCommerceHandler(), MembershipRepository, MembershipRepository

### Community 21 - "Customer"
Cohesion: 0.10
Nodes (12): UseCase, NewUseCase(), Customer, NewCustomerRepository(), onlyDigits(), scanCustomerRow(), scanCustomerRows(), NewCustomerHandler() (+4 more)

### Community 22 - "database/sql.DB"
Cohesion: 0.12
Nodes (30): applyDatabaseSQL(), applyMigrationFile(), applySQLDirectory(), applySQLFile(), initializeRepositories(), postgresRepositories(), NewAppointmentRepository(), NewAuthRepository() (+22 more)

### Community 23 - "Queue"
Cohesion: 0.08
Nodes (9): Queue, NewQueue(), TestJobQueue(), Job, Repository, NewJobRepository(), HandlerFunc, JobRepository (+1 more)

### Community 24 - "Professional"
Cohesion: 0.05
Nodes (22): UseCase, NewUseCase(), UseCase, NewUseCase(), Professional, Review, NewReviewRepository(), scanProfessionalRow() (+14 more)

### Community 25 - "Entry"
Cohesion: 0.11
Nodes (11): UseCase, NewUseCase(), Entry, NewWaitingListRepository(), scanWaitingRow(), scanWaitingRows(), NewWaitingListHandler(), WaitingListRepository (+3 more)

### Community 26 - "admin/package.json"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+24 more)

### Community 27 - "booking/package.json"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+24 more)

### Community 29 - "Service"
Cohesion: 0.09
Nodes (10): UseCase, Service, scanServiceRow(), scanServiceRows(), NewExportHandler(), NewServiceHandler(), ServiceRepository, ServiceRepository (+2 more)

### Community 30 - "UseCase"
Cohesion: 0.18
Nodes (4): RecurrenceInput, UseCase, parseRFC3339(), validStatusChange()

### Community 31 - ".registerWithRole"
Cohesion: 0.13
Nodes (14): ChangePasswordInput, PasswordResetCodeConfirmInput, PasswordResetConfirmInput, PasswordResetRequestInput, RegisterInput, ValidationError, PasswordHasher, hashResetToken() (+6 more)

### Community 32 - "Connection"
Cohesion: 0.09
Nodes (10): UseCase, NewUseCase(), TestCalendarUseCase(), Connection, Repository, SyncEvent, NewCalendarRepository(), NewCalendarSyncHandler() (+2 more)

### Community 33 - "FinancePage.tsx"
Cohesion: 0.15
Nodes (27): closeCashRegister(), createCashRegisterEntry(), createFinancialTransaction(), deleteFinancialTransaction(), fetchCashRegisterHistory(), fetchCurrentCashRegister(), fetchFinancialSummary(), fetchFinancialTransactions() (+19 more)

### Community 34 - "landing/package.json"
Cohesion: 0.06
Nodes (30): dependencies, lucide-react, react, react-dom, devDependencies, tailwindcss, @tailwindcss/vite, @types/react (+22 more)

### Community 35 - "booking/src/pages/WaitingListPage.tsx"
Cohesion: 0.17
Nodes (22): App(), fetchPublicReview(), signUpWaitingList(), submitPublicReview(), DetailsStep(), ReviewPage(), handleSubmit(), SlugEntryPage() (+14 more)

### Community 36 - "UseCase"
Cohesion: 0.10
Nodes (13): LoginInput, MFAStatus, SecretCipher, Session, TokenManager, generateTOTPSecret(), totpCode(), totpURI() (+5 more)

### Community 37 - "GiftCard"
Cohesion: 0.11
Nodes (6): generateGiftCardCode(), GiftCard, NewGiftCardRepository(), GiftCardRepository, GiftCardRepository, GiftCardRepository

### Community 39 - "Account"
Cohesion: 0.14
Nodes (6): Account, PasswordResetToken, UpdateProfileInput, cloneAuthAccount(), resetTokenExpired(), AuthRepository

### Community 40 - "PortalPage.tsx"
Cohesion: 0.19
Nodes (23): PortalAppointment, PortalResult, Slot, formatDateLong(), formatDateTime(), formatTime(), fromDateKey(), nextDateKeys() (+15 more)

### Community 41 - "AdminRoutes.tsx"
Cohesion: 0.09
Nodes (38): AccessDeniedPage(), AdminRouteFallback(), AdminRoutes(), AdminRoutesProps, AgendaPage, AnalyticsPage, AppointmentsPage, CommissionsPage (+30 more)

### Community 42 - "UseCase"
Cohesion: 0.17
Nodes (11): appointmentDurationMinutes(), formatAppointmentForTemplate(), generateManagementCode(), UseCase, hashManagementCode(), parseRFC3339(), PublicBookingOrigin(), BookingInput (+3 more)

### Community 43 - "Package"
Cohesion: 0.12
Nodes (5): Package, NewPackageRepository(), PackageRepository, PackageRepository, PackageRepository

### Community 44 - "PackageCredit"
Cohesion: 0.12
Nodes (7): NewUseCase(), PackageCredit, NewPackageCreditRepository(), scanPackageCredits(), PackageCreditRepository, PackageCreditRepository, PackageCreditRepository

### Community 45 - "bookingApi.ts"
Cohesion: 0.14
Nodes (23): backendOrigin, buildApiUrl(), AvailabilityResponse, cancelPublicAppointment(), fetchAvailability(), fetchPublicAppointment(), fetchPublicInfo(), fetchPublicProfessionals() (+15 more)

### Community 46 - "landing/src/App.tsx"
Cohesion: 0.16
Nodes (14): App(), CtaSection(), FEATURES, FeaturesSection(), Footer(), Header(), NAV_ITEMS, HeroSection() (+6 more)

### Community 47 - "PackagesPage"
Cohesion: 0.46
Nodes (8): deletePackage(), fetchPackageCredits(), fetchPackages(), setPackageActive(), PackagesPage(), handleDelete(), handleToggleActive(), loadPackages()

### Community 48 - "Resource"
Cohesion: 0.13
Nodes (7): NewUseCase(), UseCase, NewUseCase(), Resource, NewResourceHandler(), ResourceRepository, Repository

### Community 49 - "repositorySet"
Cohesion: 0.14
Nodes (6): MemberPlan, NewMemberPlanRepository(), scanMemberPlans(), repositorySet, MemberPlanRepository, MemberPlanRepository

### Community 50 - "002_features.sql"
Cohesion: 0.14
Nodes (20): catalog.gift_cards, catalog.member_plans, catalog.memberships, catalog.notification_logs, catalog.package_credits, catalog.packages, catalog.payments, catalog.service_professionals (+12 more)

### Community 52 - "UseCase"
Cohesion: 0.06
Nodes (11): UseCase, NewUseCase(), CommissionRecord, CommissionRule, Payout, Repository, Summary, NewCommissionHandler() (+3 more)

### Community 53 - "ConfigurationPage.tsx"
Cohesion: 0.14
Nodes (11): fetchConfiguration(), updateConfiguration(), DAY_LABELS, ActiveTab, CHANNEL_OPTIONS, ConfigurationPage(), handleSubmit(), EMPTY_DAY (+3 more)

### Community 54 - ".GetByID"
Cohesion: 0.19
Nodes (3): nullIfEmpty(), MemberPlanRepository, MembershipRepository

### Community 55 - "UseCase"
Cohesion: 0.08
Nodes (11): UseCase, NewUseCase(), TestFinanceUseCase(), Repository, Summary, Transaction, NewFinanceHandler(), NewMercadoPagoWebhookHandler() (+3 more)

### Community 56 - "001_baseline.sql"
Cohesion: 0.17
Nodes (21): catalog.appointments, catalog.auth_rate_limits, catalog.auth_sessions, catalog.configurations, catalog.customers, catalog.password_reset_tokens, catalog.professionals, catalog.services (+13 more)

### Community 57 - "BookingPage.tsx"
Cohesion: 0.13
Nodes (15): BookingResult, createPublicBooking(), formatMoney(), BookingPage(), handleBookingSubmit(), CenteredLayout(), EmptyNote(), ErrorCard() (+7 more)

### Community 58 - "CommissionsPage.tsx"
Cohesion: 0.14
Nodes (28): handleSubmit(), createCommissionRule(), deleteCommissionRule(), fetchAppointmentAnalytics(), fetchCommissionPayouts(), fetchCommissionRecords(), fetchCommissionRules(), fetchCommissionSummaries() (+20 more)

### Community 59 - "testing.T"
Cohesion: 0.19
Nodes (15): newAppointment(), newResourceUseCase(), TestCreateAutoAssignsAvailableResource(), TestCreateRejectsBusyExplicitResource(), TestCreateRejectsUnknownResource(), TestUpdateAutoReassignsFreedResource(), newTestUseCase(), TestAvailableSlotsRespectsClosedDays() (+7 more)

### Community 60 - "UseCase"
Cohesion: 0.14
Nodes (11): NewUseCase(), UseCase, NewUseCase(), Log, NewNotificationRepository(), NewNotificationHandler(), NotificationRepository, Provider (+3 more)

### Community 61 - "writeJSON"
Cohesion: 0.15
Nodes (7): writeBookingError(), health(), writePaymentError(), writeJSON(), BookingHandler, CashRegisterHandler, PaymentHandler

### Community 62 - "writeNotFoundError"
Cohesion: 0.17
Nodes (6): writeCustomerError(), writeNotificationError(), writeNotFoundError(), writeTenantError(), CustomerHandler, TenantHandler

### Community 63 - "Appointment"
Cohesion: 0.11
Nodes (6): Repository, Appointment, parseAppointmentBound(), parseAppointmentInterval(), BookingResult, AppointmentRepository

### Community 64 - "templates.go"
Cohesion: 0.29
Nodes (15): greeting(), RenderCancellation(), RenderConfirmation(), RenderReminder(), RenderReschedule(), sampleData(), TestGreetingWithoutName(), TestRenderCancellation() (+7 more)

### Community 65 - "AppointmentRepository"
Cohesion: 0.18
Nodes (6): StatusCount, appointmentSelectColumns(), scanAppointmentRow(), scanAppointmentRows(), AppointmentRepository, appointmentScanner

### Community 66 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 67 - "Configuration"
Cohesion: 0.22
Nodes (11): isClosedDay(), workingHoursFor(), Configuration, cloneConfiguration(), defaultMemoryConfiguration(), normalizeConfigurationTenantID(), normalizePostgresTenantID(), ClosedDay (+3 more)

### Community 68 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (22): browserSessionPayload(), buildPasswordResetURL(), exposeAuthTokenForTests(), AuthHandler, secureTemporaryPassword(), writeAdminAuthError(), newAuthMFAHandler(), AuthHandler (+14 more)

### Community 69 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 70 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 71 - "ReminderWorker"
Cohesion: 0.18
Nodes (8): MaintenanceWorker, MaintenanceConfig, NewMaintenanceWorker(), formatAppointmentLocal(), NewReminderWorker(), startBackgroundWorkers(), authComponents, ReminderWorker

### Community 72 - "006_erp_finance_stock_commissions.sql"
Cohesion: 0.21
Nodes (12): catalog.commission_records, catalog.commission_rules, catalog.financial_transactions, catalog.products, catalog.service_supplies, catalog.stock_movements, idx_catalog_commission_records_prof, idx_catalog_commission_rules_tenant (+4 more)

### Community 73 - "AppointmentHandler"
Cohesion: 0.26
Nodes (3): writeAppointmentError(), resolveActorTenantID(), AppointmentHandler

### Community 74 - "time.Duration"
Cohesion: 0.18
Nodes (6): buildAuthComponents(), durationSecondsEnv(), intEnv(), CleanupSummary, NewJWTManager(), NewSecretCipher()

### Community 75 - "Claims"
Cohesion: 0.31
Nodes (5): isProductionAuthEnvironment(), Claims, decodeSegment(), encodeSegment(), JWTManager

### Community 76 - "net/http.Handler"
Cohesion: 0.12
Nodes (13): lazyAppHandler, writeLazyJSON(), boolEnv(), resolveMaxRequestBodyBytes(), shouldLimitRequestBody(), withRequestBodyLimit(), resolveAllowedOrigins(), withCORS() (+5 more)

### Community 77 - "ERP Agendamentos"
Cohesion: 0.14
Nodes (13): CI, Endpoints principais, ERP Agendamentos, Estrutura, Migrations, Pré-requisitos, Rodando com Docker (produção local), Rodando em desenvolvimento (sem banco) (+5 more)

### Community 78 - "MembershipsPage.tsx"
Cohesion: 0.33
Nodes (12): deleteMembership(), fetchMemberPlans(), fetchMemberships(), setMembershipActive(), updateMemberPlanStatus(), MembershipsPage(), handleDelete(), handlePlanStatusChange() (+4 more)

### Community 80 - "EventBroker"
Cohesion: 0.20
Nodes (9): UseCase, NewUseCase(), NewBookingHandler(), NewCheckoutHandler(), EventBroker, RegisterMetricsBroker(), CheckoutHandler, Event (+1 more)

### Community 81 - "PaymentsPage.tsx"
Cohesion: 0.16
Nodes (21): deleteCustomer(), downloadAdminCsv(), fetchCustomers(), fetchPayments(), markPaymentPaid(), refundPayment(), handleExport(), CustomersPage() (+13 more)

### Community 82 - "Analytics"
Cohesion: 0.29
Nodes (4): DayMetric, ProfessionalMetric, ServiceMetric, Analytics

### Community 83 - "sync.Mutex"
Cohesion: 0.22
Nodes (9): countOverlappingSlots(), localMidnight(), overlapsAnyBusy(), parseClock(), ProfessionalSlots, Slot, rateLimitEntry, memoryRateLimiter (+1 more)

### Community 84 - "memoryRepositories"
Cohesion: 0.26
Nodes (15): TestCashRegisterFlow(), TestCheckoutFlow(), TestCommissionUseCase(), TestOnboardingSignup(), memoryRepositories(), NewAppointmentRepository(), NewAuthRepository(), NewCashRegisterRepository() (+7 more)

### Community 85 - "observability.go"
Cohesion: 0.10
Nodes (16): newLazyAppHandler(), main(), resolveDurationEnv(), resolvePositiveIntEnv(), buildNotificationProvider(), logFallback(), NewEmailProvider(), NewLogProvider() (+8 more)

### Community 86 - "booking/src/lib/apiError.ts"
Cohesion: 0.29
Nodes (8): ApiError, ApiErrorPayload, buildApiError(), parseApiResponse(), API_ERROR_MESSAGES, FIELD_ERROR_MESSAGES, getApiErrorMessage(), translateFieldErrors()

### Community 87 - "catalog.commission_payouts"
Cohesion: 0.29
Nodes (3): catalog.commission_payouts, idx_commission_payouts_tenant_prof, idx_commission_records_payout

### Community 89 - "ResourceRepository"
Cohesion: 0.31
Nodes (3): scanResourceRow(), ResourceRepository, resourceScanner

### Community 91 - "PortalPage"
Cohesion: 0.39
Nodes (8): portalCancelAppointment(), portalLookup(), portalRescheduleAppointment(), PortalPage(), handleCancel(), handleConfirmReschedule(), handleLookup(), refresh()

### Community 92 - ".UpdateAccountRole"
Cohesion: 0.39
Nodes (4): canManageTenantAccount(), isPrivilegedAdminRole(), isStaffRole(), normalizeTenantScope()

### Community 94 - "ical.go"
Cohesion: 0.83
Nodes (3): BuildICalCalendar(), formatUTC(), sanitizeLine()

### Community 98 - "005_resources.sql"
Cohesion: 0.40
Nodes (3): catalog.resources, idx_appointments_resource, idx_resources_tenant

### Community 99 - "007_jobs_and_calendar_sync.sql"
Cohesion: 0.47
Nodes (4): catalog.calendar_connections, catalog.jobs, idx_calendar_connections_tenant, idx_catalog_jobs_poll

### Community 100 - "010_branding_pos_and_cash_registers.sql"
Cohesion: 0.60
Nodes (4): catalog.cash_register_entries, catalog.cash_registers, idx_catalog_cash_reg_entries_reg, idx_catalog_cash_registers_tenant

### Community 101 - "booking/src/shared/components/Button.tsx"
Cohesion: 0.33
Nodes (5): ButtonProps, ButtonSize, ButtonVariant, sizeClasses, variantClasses

### Community 104 - "catalog.reviews"
Cohesion: 0.60
Nodes (3): catalog.reviews, idx_reviews_tenant_created, idx_reviews_tenant_status

## Knowledge Gaps
- **266 isolated node(s):** `graphify-mcp`, `name`, `private`, `version`, `type` (+261 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 528 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **67 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Container` connect `Container` to `Payment`, `UseCase`, `container.go`, `Tenant`, `CashRegister`, `UseCase`, `Customer`, `database/sql.DB`, `Queue`, `Professional`, `Entry`, `Service`, `UseCase`, `Connection`, `UseCase`, `UseCase`, `Resource`, `UseCase`, `UseCase`, `UseCase`, `ReminderWorker`, `net/http.Handler`, `EventBroker`?**
  _High betweenness centrality (0.040) - this node is a cross-community bridge._
- **Why does `UseCase` connect `UseCase` to `net/http.ResponseWriter`, `Repository`, `Account`, `writeAPIError`, `AppointmentHandler`, `time.Duration`, `Claims`, `ReminderWorker`, `net/http.Request`, `DeviceSession`, `RouterConfig`, `auth/usecase.go`, `Container`, `sync.Mutex`, `.UpdateAccountRole`, `.registerWithRole`?**
  _High betweenness centrality (0.029) - this node is a cross-community bridge._
- **Why does `RouterConfig` connect `RouterConfig` to `Payment`, `UseCase`, `container.go`, `Tenant`, `Container`, `CashRegister`, `UseCase`, `Customer`, `Professional`, `Entry`, `Service`, `UseCase`, `Connection`, `UseCase`, `UseCase`, `Resource`, `UseCase`, `UseCase`, `UseCase`, `ReminderWorker`, `time.Duration`, `EventBroker`?**
  _High betweenness centrality (0.029) - this node is a cross-community bridge._
- **Are the 128 inferred relationships involving `writeJSON()` (e.g. with `health()` and `readiness()`) actually correct?**
  _`writeJSON()` has 128 INFERRED edges - model-reasoned connections that need verification._
- **Are the 103 inferred relationships involving `resolveRequestTenantID()` (e.g. with `.Analytics()` and `.CheckIn()`) actually correct?**
  _`resolveRequestTenantID()` has 103 INFERRED edges - model-reasoned connections that need verification._
- **What connects `graphify-mcp`, `name`, `private` to the rest of the system?**
  _266 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `adminAuth.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.07168458781362007 - nodes in this community are weakly interconnected._
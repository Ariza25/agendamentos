# Graph Report - erpagendamentos  (2026-11-01)

## Corpus Check
- 263 files · ~125,703 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 16 file(s) not represented in the graph (top: (none) 7, .conf 5, .css 3)

## Summary
- 2776 nodes · 9034 edges · 144 communities (86 shown, 58 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 661 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9a990ce7`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- TenantsPage.tsx
- adminApi.ts
- Account
- Payment
- schema.sql
- adminFetch
- container.go
- templates.go
- UseCase
- Button
- adminAuth.ts
- ref_react
- UseCase
- AdminRoutes.tsx
- UseCase
- net/http.Request
- Customer
- Repository
- auth/usecase.go
- Tenant
- Professional
- writeJSON
- Queue
- Resource
- Entry
- net/http.ResponseWriter
- admin/package.json
- booking/package.json
- UseCase
- email_provider.go
- Review
- Connection
- RouterConfig
- landing/package.json
- MemberPlan
- database/sql.DB
- booking/src/pages/WaitingListPage.tsx
- Service
- Container
- PortalPage.tsx
- UseCase
- CashRegister
- Membership
- Package
- body_limit.go
- bookingApi.ts
- landing/src/App.tsx
- UseCase
- PackageCredit
- GiftCard
- time.Duration
- AuthRepository
- BookingPage.tsx
- CommissionsPage.tsx
- ReminderWorker
- PageContainer
- Appointment
- compilerOptions
- compilerOptions
- compilerOptions
- PackagesPage.tsx
- .Add
- AppointmentRepository
- AppointmentHandler
- Configuration
- InventoryPage.tsx
- testing.T
- ERP Agendamentos
- UseCase
- Transaction
- sync.Mutex
- EventBroker
- booking/src/lib/apiError.ts
- Claims
- Repository
- response.go
- memoryRepositories
- UseCase
- PortalPage
- MembershipsPage.tsx
- ResourceRepository
- writeProfessionalError
- writeInvalidPayload
- writeServiceError
- CashRegisterRepository
- main.go
- checkout/usecase.go
- SecretCipher
- writePaymentError
- Analytics
- Summary
- booking/src/shared/components/Button.tsx
- bcrypt_hasher.go
- widget_handler.go
- booking/src/main.tsx
- backup_postgres.sh
- graphify
- rules/graphify.md
- workflows/graphify.md
- generate-certs.sh
- erpagendamentos/backend
- certbot_ssl.sh

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
- `buildAuthComponents()` --calls--> `durationSecondsEnv()`  [INFERRED]
  backend/internal/bootstrap/auth_components.go → backend/internal/bootstrap/container.go

## Import Cycles
- None detected.

## Communities (144 total, 58 thin omitted)

### Community 0 - "TenantsPage.tsx"
Cohesion: 0.07
Nodes (52): createAdminUser(), deleteAdminUser(), deleteGiftCard(), deleteResource(), deleteTenant(), fetchAdminUsers(), fetchGiftCards(), fetchResources() (+44 more)

### Community 1 - "adminApi.ts"
Cohesion: 0.06
Nodes (65): AppointmentFilters, closeCashRegister(), createCashRegisterEntry(), createFinancialTransaction(), deleteFinancialTransaction(), fetchCashRegisterHistory(), fetchCurrentCashRegister(), fetchFinancialSummary() (+57 more)

### Community 2 - "Account"
Cohesion: 0.13
Nodes (6): Account, PasswordResetToken, UpdateProfileInput, cloneAuthAccount(), resetTokenExpired(), AuthRepository

### Community 3 - "Payment"
Cohesion: 0.05
Nodes (23): UseCase, NewUseCase(), normalizeStatus(), buildPaymentProvider(), IntentResult, Payment, NewMercadoPagoProvider(), NewMockProvider() (+15 more)

### Community 4 - "schema.sql"
Cohesion: 0.06
Nodes (70): catalog.appointments, catalog.auth_rate_limits, catalog.auth_sessions, catalog.calendar_connections, catalog.cash_register_entries, catalog.cash_registers, catalog.commission_payouts, catalog.commission_records (+62 more)

### Community 5 - "adminFetch"
Cohesion: 0.08
Nodes (50): connectCalendar(), createCustomer(), createGiftCard(), createMemberPlan(), createMembership(), createPackage(), createPackageCredit(), createPayment() (+42 more)

### Community 6 - "container.go"
Cohesion: 0.09
Nodes (6): deploymentEnvironment(), isProductionEnv(), CategoryTotal, Type, Method, Status

### Community 7 - "templates.go"
Cohesion: 0.27
Nodes (16): greeting(), PublicBookingOrigin(), RenderCancellation(), RenderConfirmation(), RenderReminder(), RenderReschedule(), sampleData(), TestGreetingWithoutName() (+8 more)

### Community 8 - "UseCase"
Cohesion: 0.09
Nodes (16): LoginInput, MFAStatus, RegisterInput, SecretCipher, Session, generateTOTPSecret(), totpCode(), totpURI() (+8 more)

### Community 9 - "Button"
Cohesion: 0.10
Nodes (54): AppointmentFormModal(), CheckoutComandaModal(), createService(), fetchConfiguration(), fetchNotifications(), updateConfiguration(), updateService(), DAY_LABELS (+46 more)

### Community 10 - "adminAuth.ts"
Cohesion: 0.07
Nodes (54): App(), handleLogout(), AppProviders(), AppProvidersProps, ResetPasswordPage(), handleSubmit(), AppRouter(), ADMIN_ROLES (+46 more)

### Community 11 - "ref_react"
Cohesion: 0.06
Nodes (64): computeEndsAt(), handleSubmit(), AppointmentFormModalProps, STATUS_OPTIONS, CheckoutComandaModalProps, SelectedProduct, checkInAppointment(), createAppointment() (+56 more)

### Community 12 - "UseCase"
Cohesion: 0.06
Nodes (11): UseCase, NewUseCase(), CommissionRecord, CommissionRule, Payout, Repository, Summary, NewCommissionHandler() (+3 more)

### Community 13 - "AdminRoutes.tsx"
Cohesion: 0.09
Nodes (39): AccessDeniedPage(), AdminRouteFallback(), AdminRoutes(), AdminRoutesProps, AgendaPage, AnalyticsPage, AppointmentsPage, CommissionsPage (+31 more)

### Community 14 - "UseCase"
Cohesion: 0.07
Nodes (10): UseCase, NewUseCase(), TestStockUseCase(), Product, Repository, ServiceSupply, StockMovement, NewStockHandler() (+2 more)

### Community 16 - "net/http.Request"
Cohesion: 0.07
Nodes (41): requireMasterAdmin(), resolveAuthClaims(), resolveTenantID(), authTokenFromCookie(), browserSessionPayload(), buildPasswordResetURL(), exposeAuthTokenForTests(), AuthHandler (+33 more)

### Community 17 - "Customer"
Cohesion: 0.09
Nodes (13): UseCase, Customer, NewCustomerRepository(), onlyDigits(), scanCustomerRow(), scanCustomerRows(), NewCustomerHandler(), NewExportHandler() (+5 more)

### Community 18 - "Repository"
Cohesion: 0.10
Nodes (9): Repository, TokenManager, PasswordHasher, NewUseCase(), buildAuthComponents(), seedAdminAccount(), NewJWTManager(), NewSecretCipher() (+1 more)

### Community 19 - "auth/usecase.go"
Cohesion: 0.08
Nodes (9): isSensitiveLogKey(), redactSensitiveText(), sanitizeLogString(), sanitizeLogValue(), truncateLogString(), newMemoryRateLimiter(), resolvePositiveEnvInt(), contextKey (+1 more)

### Community 20 - "Tenant"
Cohesion: 0.08
Nodes (16): UseCase, NewUseCase(), ensureDefaultTenant(), Tenant, UpdateSubscriptionInput, defaultModules(), isTenantIdentityConflict(), defaultModules() (+8 more)

### Community 21 - "Professional"
Cohesion: 0.11
Nodes (9): UseCase, Professional, scanProfessionalRow(), scanProfessionalRows(), NewProfessionalHandler(), ProfessionalRepository, ProfessionalRepository, professionalScanner (+1 more)

### Community 22 - "writeJSON"
Cohesion: 0.11
Nodes (9): resolveRequestTenantID(), writeCommerceError(), writeInternalError(), writeJSON(), CalendarSyncHandler, CommerceHandler, CommissionHandler, FinanceHandler (+1 more)

### Community 24 - "Queue"
Cohesion: 0.08
Nodes (9): Queue, NewQueue(), TestJobQueue(), Job, Repository, NewJobRepository(), HandlerFunc, JobRepository (+1 more)

### Community 25 - "Resource"
Cohesion: 0.13
Nodes (7): NewUseCase(), UseCase, NewUseCase(), Resource, NewResourceHandler(), ResourceRepository, Repository

### Community 26 - "Entry"
Cohesion: 0.11
Nodes (11): UseCase, NewUseCase(), Entry, NewWaitingListRepository(), scanWaitingRow(), scanWaitingRows(), NewWaitingListHandler(), WaitingListRepository (+3 more)

### Community 27 - "net/http.ResponseWriter"
Cohesion: 0.10
Nodes (14): writeAdminAuthError(), writeBookingError(), writeCustomerError(), writeNotificationError(), writeConflictError(), writeNotFoundError(), writeReviewError(), writeTenantError() (+6 more)

### Community 28 - "admin/package.json"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+24 more)

### Community 29 - "booking/package.json"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+24 more)

### Community 30 - "UseCase"
Cohesion: 0.16
Nodes (6): RecurrenceInput, UseCase, parseRFC3339(), validStatusChange(), NewCalendarHandler(), CalendarHandler

### Community 31 - "email_provider.go"
Cohesion: 0.16
Nodes (9): boolEnv(), buildNotificationProvider(), durationSecondsEnv(), intEnv(), logFallback(), NewEmailProvider(), NewLogProvider(), EmailProvider (+1 more)

### Community 32 - "Review"
Cohesion: 0.12
Nodes (12): UseCase, NewUseCase(), Review, NewReviewRepository(), scanReviewRow(), scanReviewRows(), NewAppointmentHandler(), NewReviewHandler() (+4 more)

### Community 33 - "Connection"
Cohesion: 0.09
Nodes (10): UseCase, NewUseCase(), TestCalendarUseCase(), Connection, Repository, SyncEvent, NewCalendarRepository(), NewCalendarSyncHandler() (+2 more)

### Community 34 - "RouterConfig"
Cohesion: 0.08
Nodes (62): lazyAppHandler, writeLazyJSON(), bearerAuthAllowed(), hasRole(), isAdminRole(), isPrivilegedAdminRole(), requestWithAuthClaims(), requireAdmin() (+54 more)

### Community 35 - "landing/package.json"
Cohesion: 0.06
Nodes (30): dependencies, lucide-react, react, react-dom, devDependencies, tailwindcss, @tailwindcss/vite, @types/react (+22 more)

### Community 36 - "MemberPlan"
Cohesion: 0.12
Nodes (8): NewUseCase(), MemberPlan, NewMemberPlanRepository(), nullIfEmpty(), scanMemberPlans(), MemberPlanRepository, MemberPlanRepository, MemberPlanRepository

### Community 37 - "database/sql.DB"
Cohesion: 0.12
Nodes (31): applyDatabaseSQL(), resolveSchemaFile(), initializeRepositories(), postgresRepositories(), NewAppointmentRepository(), NewAuthRepository(), NewCalendarRepository(), NewCashRegisterRepository() (+23 more)

### Community 38 - "booking/src/pages/WaitingListPage.tsx"
Cohesion: 0.17
Nodes (22): App(), fetchPublicReview(), signUpWaitingList(), submitPublicReview(), DetailsStep(), ReviewPage(), handleSubmit(), SlugEntryPage() (+14 more)

### Community 39 - "Service"
Cohesion: 0.09
Nodes (10): UseCase, NewUseCase(), Service, scanServiceRow(), scanServiceRows(), NewServiceHandler(), ServiceRepository, ServiceRepository (+2 more)

### Community 40 - "Container"
Cohesion: 0.13
Nodes (16): NewUseCase(), UseCase, NewUseCase(), NewUseCase(), UseCase, NewUseCase(), NewUseCase(), Container (+8 more)

### Community 41 - "PortalPage.tsx"
Cohesion: 0.19
Nodes (23): PortalAppointment, PortalResult, Slot, formatDateLong(), formatDateTime(), formatTime(), fromDateKey(), nextDateKeys() (+15 more)

### Community 42 - "UseCase"
Cohesion: 0.18
Nodes (10): appointmentDurationMinutes(), formatAppointmentForTemplate(), generateManagementCode(), UseCase, hashManagementCode(), parseRFC3339(), BookingInput, PortalAppointment (+2 more)

### Community 43 - "CashRegister"
Cohesion: 0.13
Nodes (8): UseCase, NewUseCase(), CashRegister, CashRegisterEntry, Repository, ShiftSummary, NewCashRegisterHandler(), CashRegisterRepository

### Community 44 - "Membership"
Cohesion: 0.12
Nodes (5): Membership, NewMembershipRepository(), MembershipRepository, MembershipRepository, MembershipRepository

### Community 45 - "Package"
Cohesion: 0.12
Nodes (5): Package, NewPackageRepository(), PackageRepository, PackageRepository, PackageRepository

### Community 46 - "body_limit.go"
Cohesion: 0.83
Nodes (3): resolveMaxRequestBodyBytes(), shouldLimitRequestBody(), withRequestBodyLimit()

### Community 47 - "bookingApi.ts"
Cohesion: 0.14
Nodes (23): backendOrigin, buildApiUrl(), AvailabilityResponse, cancelPublicAppointment(), fetchAvailability(), fetchPublicAppointment(), fetchPublicInfo(), fetchPublicProfessionals() (+15 more)

### Community 48 - "landing/src/App.tsx"
Cohesion: 0.14
Nodes (17): App(), CtaSection(), FEATURES, FeaturesSection(), Footer(), Header(), NAV_ITEMS, HeroSection() (+9 more)

### Community 49 - "UseCase"
Cohesion: 0.15
Nodes (10): UseCase, NewUseCase(), Log, NewNotificationRepository(), NewNotificationHandler(), NotificationRepository, Provider, Repository (+2 more)

### Community 50 - "PackageCredit"
Cohesion: 0.13
Nodes (6): PackageCredit, NewPackageCreditRepository(), scanPackageCredits(), PackageCreditRepository, PackageCreditRepository, PackageCreditRepository

### Community 51 - "GiftCard"
Cohesion: 0.12
Nodes (5): GiftCard, NewGiftCardRepository(), GiftCardRepository, GiftCardRepository, GiftCardRepository

### Community 52 - "time.Duration"
Cohesion: 0.11
Nodes (3): CleanupSummary, DeviceSession, sessionExpired()

### Community 53 - "AuthRepository"
Cohesion: 0.11
Nodes (3): AuthState, updateAuthPasswordAndVersion(), AuthRepository

### Community 54 - "BookingPage.tsx"
Cohesion: 0.13
Nodes (15): BookingResult, createPublicBooking(), formatMoney(), BookingPage(), handleBookingSubmit(), CenteredLayout(), EmptyNote(), ErrorCard() (+7 more)

### Community 55 - "CommissionsPage.tsx"
Cohesion: 0.19
Nodes (19): handleSubmit(), checkoutAppointment(), createCommissionPayout(), createCommissionRule(), deleteCommissionRule(), fetchCommissionPayouts(), fetchCommissionRecords(), fetchCommissionRules() (+11 more)

### Community 56 - "ReminderWorker"
Cohesion: 0.21
Nodes (6): MaintenanceWorker, MaintenanceConfig, NewMaintenanceWorker(), formatAppointmentLocal(), NewReminderWorker(), ReminderWorker

### Community 57 - "PageContainer"
Cohesion: 0.09
Nodes (41): deleteCustomer(), downloadAdminCsv(), fetchAppointmentAnalytics(), fetchCustomers(), fetchPayments(), fetchReviews(), markPaymentPaid(), refundPayment() (+33 more)

### Community 58 - "Appointment"
Cohesion: 0.11
Nodes (6): Repository, Appointment, StatusCount, parseAppointmentBound(), parseAppointmentInterval(), AppointmentRepository

### Community 59 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 60 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 61 - "compilerOptions"
Cohesion: 0.11
Nodes (17): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+9 more)

### Community 62 - "PackagesPage.tsx"
Cohesion: 0.50
Nodes (8): deletePackage(), fetchPackageCredits(), fetchPackages(), setPackageActive(), PackagesPage(), handleDelete(), handleToggleActive(), loadPackages()

### Community 63 - ".Add"
Cohesion: 0.20
Nodes (11): ChangePasswordInput, PasswordResetCodeConfirmInput, PasswordResetConfirmInput, PasswordResetRequestInput, ValidationError, hashResetToken(), NewValidationError(), passwordResetCodeChallengeKey() (+3 more)

### Community 64 - "AppointmentRepository"
Cohesion: 0.20
Nodes (5): appointmentSelectColumns(), scanAppointmentRow(), scanAppointmentRows(), AppointmentRepository, appointmentScanner

### Community 65 - "AppointmentHandler"
Cohesion: 0.26
Nodes (3): writeAppointmentError(), resolveActorTenantID(), AppointmentHandler

### Community 66 - "Configuration"
Cohesion: 0.25
Nodes (10): workingHoursFor(), Configuration, cloneConfiguration(), defaultMemoryConfiguration(), normalizeConfigurationTenantID(), normalizePostgresTenantID(), ClosedDay, WorkingHours (+2 more)

### Community 67 - "InventoryPage.tsx"
Cohesion: 0.22
Nodes (18): createProduct(), createStockMovement(), deleteProduct(), fetchProducts(), fetchServices(), fetchServiceSupplies(), fetchStockMovements(), setServiceSupplies() (+10 more)

### Community 68 - "testing.T"
Cohesion: 0.19
Nodes (15): newAppointment(), newResourceUseCase(), TestCreateAutoAssignsAvailableResource(), TestCreateRejectsBusyExplicitResource(), TestCreateRejectsUnknownResource(), TestUpdateAutoReassignsFreedResource(), newTestUseCase(), TestAvailableSlotsRespectsClosedDays() (+7 more)

### Community 69 - "ERP Agendamentos"
Cohesion: 0.14
Nodes (13): Banco de Dados & Schema, CI, Endpoints principais, ERP Agendamentos, Estrutura, Pré-requisitos, Rodando com Docker (produção local), Rodando em desenvolvimento (sem banco) (+5 more)

### Community 70 - "UseCase"
Cohesion: 0.17
Nodes (3): generateGiftCardCode(), UseCase, NewCommerceHandler()

### Community 71 - "Transaction"
Cohesion: 0.24
Nodes (3): Transaction, FinanceRepository, FinanceRepository

### Community 72 - "sync.Mutex"
Cohesion: 0.16
Nodes (13): countOverlappingSlots(), isClosedDay(), localMidnight(), overlapsAnyBusy(), parseClock(), BuildICalCalendar(), formatUTC(), sanitizeLine() (+5 more)

### Community 73 - "EventBroker"
Cohesion: 0.27
Nodes (5): NewBookingHandler(), EventBroker, RegisterMetricsBroker(), Event, MetricsRegistry

### Community 74 - "booking/src/lib/apiError.ts"
Cohesion: 0.29
Nodes (8): ApiError, ApiErrorPayload, buildApiError(), parseApiResponse(), API_ERROR_MESSAGES, FIELD_ERROR_MESSAGES, getApiErrorMessage(), translateFieldErrors()

### Community 75 - "Claims"
Cohesion: 0.24
Nodes (5): isProductionAuthEnvironment(), Claims, decodeSegment(), encodeSegment(), JWTManager

### Community 76 - "Repository"
Cohesion: 0.24
Nodes (5): UseCase, NewUseCase(), Repository, NewCheckoutHandler(), CheckoutHandler

### Community 77 - "response.go"
Cohesion: 0.22
Nodes (4): writeCSV(), writeRawJSON(), errorResponse, ExportHandler

### Community 78 - "memoryRepositories"
Cohesion: 0.26
Nodes (15): TestCashRegisterFlow(), TestCheckoutFlow(), TestCommissionUseCase(), TestOnboardingSignup(), memoryRepositories(), NewAppointmentRepository(), NewAuthRepository(), NewCashRegisterRepository() (+7 more)

### Community 79 - "UseCase"
Cohesion: 0.20
Nodes (6): UseCase, NewUseCase(), TestFinanceUseCase(), NewFinanceHandler(), NewMercadoPagoWebhookHandler(), MercadoPagoWebhookHandler

### Community 80 - "PortalPage"
Cohesion: 0.39
Nodes (8): portalCancelAppointment(), portalLookup(), portalRescheduleAppointment(), PortalPage(), handleCancel(), handleConfirmReschedule(), handleLookup(), refresh()

### Community 81 - "MembershipsPage.tsx"
Cohesion: 0.33
Nodes (12): deleteMembership(), fetchMemberPlans(), fetchMemberships(), setMembershipActive(), updateMemberPlanStatus(), MembershipsPage(), handleDelete(), handlePlanStatusChange() (+4 more)

### Community 82 - "ResourceRepository"
Cohesion: 0.31
Nodes (3): scanResourceRow(), ResourceRepository, resourceScanner

### Community 84 - "writeInvalidPayload"
Cohesion: 0.12
Nodes (8): writeConfigurationError(), writeResourceError(), writeBadRequestError(), writeInvalidPayload(), writeWaitingListError(), CashRegisterHandler, ResourceHandler, WaitingListHandler

### Community 87 - "main.go"
Cohesion: 0.38
Nodes (4): newLazyAppHandler(), main(), resolveDurationEnv(), resolvePositiveIntEnv()

### Community 88 - "checkout/usecase.go"
Cohesion: 0.32
Nodes (3): CheckoutInput, CheckoutResult, ProductItem

### Community 92 - "Analytics"
Cohesion: 0.33
Nodes (4): DayMetric, ProfessionalMetric, ServiceMetric, Analytics

### Community 94 - "booking/src/shared/components/Button.tsx"
Cohesion: 0.33
Nodes (5): ButtonProps, ButtonSize, ButtonVariant, sizeClasses, variantClasses

## Knowledge Gaps
- **270 isolated node(s):** `graphify-mcp`, `name`, `private`, `version`, `type` (+265 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 512 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **58 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Container` connect `Container` to `Payment`, `container.go`, `UseCase`, `UseCase`, `UseCase`, `Customer`, `Tenant`, `Professional`, `Queue`, `Resource`, `Entry`, `UseCase`, `email_provider.go`, `Review`, `Connection`, `RouterConfig`, `database/sql.DB`, `Service`, `UseCase`, `CashRegister`, `UseCase`, `ReminderWorker`, `UseCase`, `EventBroker`, `Repository`, `UseCase`?**
  _High betweenness centrality (0.038) - this node is a cross-community bridge._
- **Why does `UseCase` connect `UseCase` to `AppointmentHandler`, `Account`, `RouterConfig`, `sync.Mutex`, `Container`, `Claims`, `net/http.Request`, `Repository`, `auth/usecase.go`, `time.Duration`, `writeJSON`, `.Add`?**
  _High betweenness centrality (0.029) - this node is a cross-community bridge._
- **Why does `RouterConfig` connect `RouterConfig` to `Payment`, `container.go`, `UseCase`, `UseCase`, `UseCase`, `Customer`, `Tenant`, `Professional`, `Resource`, `Entry`, `UseCase`, `Review`, `Connection`, `Service`, `Container`, `UseCase`, `CashRegister`, `UseCase`, `time.Duration`, `ReminderWorker`, `UseCase`, `EventBroker`, `Repository`, `UseCase`?**
  _High betweenness centrality (0.029) - this node is a cross-community bridge._
- **Are the 128 inferred relationships involving `writeJSON()` (e.g. with `health()` and `readiness()`) actually correct?**
  _`writeJSON()` has 128 INFERRED edges - model-reasoned connections that need verification._
- **Are the 103 inferred relationships involving `resolveRequestTenantID()` (e.g. with `.Analytics()` and `.CheckIn()`) actually correct?**
  _`resolveRequestTenantID()` has 103 INFERRED edges - model-reasoned connections that need verification._
- **What connects `graphify-mcp`, `name`, `private` to the rest of the system?**
  _270 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `TenantsPage.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.0715846994535519 - nodes in this community are weakly interconnected._
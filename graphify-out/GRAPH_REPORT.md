# Graph Report - erpagendamentos  (2026-11-02)

## Corpus Check
- 280 files · ~132,793 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2696 nodes · 8031 edges · 124 communities (110 shown, 14 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 683 edges (avg confidence: 0.79)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `74d08265`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- adminApi.ts
- schema.sql
- TenantsPage.tsx
- Payment
- net/http.ResponseWriter
- ConfigurationPage
- go_pkg_net_http
- downloadAdminCsv
- adminAuth.ts
- writeAPIError
- useToast
- types.ts
- Review
- UseCase
- net/http.Request
- Tenant
- landing/src/App.tsx
- Service
- go_pkg_time
- UseCase
- auth_cookie.go
- database/sql.DB
- resource_test.go
- Customer
- UseCase
- CommissionsPage.tsx
- Resource
- CashRegister
- Professional
- Queue
- Entry
- devDependencies
- devDependencies
- Container
- .Update
- auth/usecase.go
- UseCase
- RouterConfig
- devDependencies
- booking/src/pages/WaitingListPage.tsx
- MemberPlan
- writeJSON
- Configuration
- NewContainer
- BookingPage.tsx
- AdminShell.tsx
- InventoryPage
- AuthRepository
- UseCase
- AppointmentHandler
- bookingApi.ts
- memoryRepositories
- Membership
- Package
- UseCase
- auth/entities.go
- PackageCredit
- GiftCard
- UseCase
- BookingPage
- AuditLog
- go_pkg_errors
- EventBroker
- rate_limit.go
- AuthRepository
- Appointment
- writeInvalidPayload
- testing.T
- templates.go
- withIdempotency
- compilerOptions
- AppointmentRepository
- compilerOptions
- compilerOptions
- writeNotFoundError
- OutboxRepository
- MembershipsPage
- totp.go
- admin/src/App.tsx
- UseCase
- ERP Agendamentos
- time.Duration
- UseCase
- email_provider.go
- admin/src/lib/apiError.ts
- booking/src/lib/apiError.ts
- Claims
- ExportHandler
- FinancePage
- .UpdateAccountRole
- SecretCipher
- observability.go
- admin/vite.config.ts
- appointment/repository.go
- NewBcryptHasher
- time.Time
- TenantsPage
- adminHttp.ts
- ical.go
- booking/src/App.tsx
- graphify
- rules/graphify.md
- workflows/graphify.md
- generate-certs.sh
- backup_postgres.sh
- certbot_ssl.sh
- catalog.finance_transactions
- AppProviders.tsx
- MaintenanceWorker
- go_pkg_path_filepath
- erpagendamentos/backend
- writeServiceError
- PortalPage.tsx
- main.go
- writePaymentError
- TenantSubscriptionModal
- tx.go
- TeamPage
- widget_handler.go
- body_limit.go
- ref_react_dom

## God Nodes (most connected - your core abstractions)
1. `writeJSON()` - 134 edges
2. `resolveRequestTenantID()` - 108 edges
3. `adminFetch()` - 77 edges
4. `writeInvalidPayload()` - 71 edges
5. `useToast()` - 68 edges
6. `UseCase` - 67 edges
7. `invalidateAdminGetCache()` - 64 edges
8. `writeInternalError()` - 58 edges
9. `RouterConfig` - 56 edges
10. `writeAPIError()` - 46 edges

## Surprising Connections (you probably didn't know these)
- `NewContainer()` --calls--> `loadEnvFile()`  [INFERRED]
  backend/internal/bootstrap/container.go → backend/internal/bootstrap/env.go
- `handleSubmit()` --calls--> `updateConfiguration()`  [EXTRACTED]
  admin/src/pages/ConfigurationPage.tsx → admin/src/lib/adminApi.ts
- `loadAuditLogs()` --calls--> `fetchAuditLogs()`  [EXTRACTED]
  admin/src/pages/TenantsPage.tsx → admin/src/lib/adminApi.ts
- `handleSubmit()` --calls--> `loginAdmin()`  [EXTRACTED]
  admin/src/pages/LoginAdminPage.tsx → admin/src/lib/adminAuth.ts
- `handleSubmit()` --calls--> `signupAdmin()`  [EXTRACTED]
  admin/src/pages/SignupPage.tsx → admin/src/lib/adminAuth.ts

## Import Cycles
- None detected.

## Communities (124 total, 14 thin omitted)

### Community 0 - "adminApi.ts"
Cohesion: 0.05
Nodes (98): AppointmentFilters, checkInAppointment(), checkoutAppointment(), checkSystemReadiness(), closeCashRegister(), connectCalendar(), createAdminUser(), createAppointment() (+90 more)

### Community 1 - "schema.sql"
Cohesion: 0.12
Nodes (35): catalog.appointments, catalog.audit_logs, catalog.auth_rate_limits, catalog.auth_sessions, catalog.calendar_connections, catalog.cash_register_entries, catalog.cash_registers, catalog.commission_payouts (+27 more)

### Community 2 - "TenantsPage.tsx"
Cohesion: 0.07
Nodes (85): AdminRoutesProps, AgendaPage, AnalyticsPage, AppointmentsPage, CustomersPage, GiftCardsPage, InventoryPage, MembershipsPage (+77 more)

### Community 3 - "Payment"
Cohesion: 0.06
Nodes (26): UseCase, NewUseCase(), normalizeStatus(), buildPaymentProvider(), IntentResult, Payment, NewMercadoPagoProvider(), NewMockProvider() (+18 more)

### Community 4 - "net/http.ResponseWriter"
Cohesion: 0.09
Nodes (13): resolveRequestTenantID(), writeConfigurationError(), writeResourceError(), writeInternalError(), writeReviewError(), writeWaitingListError(), net/http.ResponseWriter, CalendarSyncHandler (+5 more)

### Community 6 - "go_pkg_net_http"
Cohesion: 0.07
Nodes (27): deploymentEnvironment(), isProductionEnv(), go_pkg_encoding_json, go_pkg_erpagendamentos_backend_internal_application_booking, go_pkg_erpagendamentos_backend_internal_application_calendar, go_pkg_erpagendamentos_backend_internal_application_cashregister, go_pkg_erpagendamentos_backend_internal_application_checkout, go_pkg_erpagendamentos_backend_internal_application_commerce (+19 more)

### Community 7 - "downloadAdminCsv"
Cohesion: 0.11
Nodes (22): downloadAdminCsv(), fetchCustomers(), fetchPayments(), fetchReviews(), handleExport(), CustomersPage(), handleDelete(), handleExport() (+14 more)

### Community 8 - "adminAuth.ts"
Cohesion: 0.23
Nodes (17): ADMIN_ROLES, AdminLoginInput, AdminSession, buildAdminAuthHeaders(), dispatchAdminAuthChange(), getAdminAccount(), getAdminCSRFToken(), getAdminSession() (+9 more)

### Community 9 - "writeAPIError"
Cohesion: 0.08
Nodes (36): hasRole(), isAdminRole(), isPrivilegedAdminRole(), requestWithAuthClaims(), requireMasterAdmin(), requireOwnerOrAdmin(), requirePrivilegedAdmin(), requireRoles() (+28 more)

### Community 10 - "useToast"
Cohesion: 0.05
Nodes (71): AppointmentFormModal(), computeEndsAt(), handleSubmit(), AppointmentFormModalProps, STATUS_OPTIONS, CheckoutComandaModal(), handleSubmit(), CheckoutComandaModalProps (+63 more)

### Community 11 - "types.ts"
Cohesion: 0.05
Nodes (44): ConfigurationPage, FinancePage, DAY_LABELS, ActiveTab, CHANNEL_OPTIONS, EMPTY_DAY, PRESET_COLORS, ActiveTab (+36 more)

### Community 12 - "Review"
Cohesion: 0.14
Nodes (12): UseCase, NewUseCase(), Review, NewReviewRepository(), scanReviewRow(), scanReviewRows(), NewReviewHandler(), ReviewRepository (+4 more)

### Community 13 - "UseCase"
Cohesion: 0.09
Nodes (7): UseCase, Product, ServiceSupply, StockMovement, NewStockHandler(), StockRepository, StockRepository

### Community 14 - "net/http.Request"
Cohesion: 0.13
Nodes (7): buildPasswordResetURL(), newAuthMFAHandler(), health(), HandleMetrics(), writeTenantError(), net/http.Request, TenantHandler

### Community 15 - "Tenant"
Cohesion: 0.10
Nodes (13): UseCase, Tenant, UpdateSubscriptionInput, defaultModules(), isTenantIdentityConflict(), defaultModules(), scanTenantRow(), scanTenantRows() (+5 more)

### Community 16 - "landing/src/App.tsx"
Cohesion: 0.10
Nodes (22): App(), CtaSection(), FEATURES, FeaturesSection(), Footer(), Header(), NAV_ITEMS, HeroSection() (+14 more)

### Community 17 - "Service"
Cohesion: 0.11
Nodes (8): UseCase, Service, scanServiceRow(), scanServiceRows(), NewServiceHandler(), ServiceRepository, ServiceRepository, serviceScanner

### Community 18 - "go_pkg_time"
Cohesion: 0.10
Nodes (24): go_pkg_context, go_pkg_database_sql, go_pkg_erpagendamentos_backend_internal_domain_appointment, go_pkg_erpagendamentos_backend_internal_domain_commerce, go_pkg_erpagendamentos_backend_internal_domain_configuration, go_pkg_erpagendamentos_backend_internal_domain_customer, go_pkg_erpagendamentos_backend_internal_domain_job, go_pkg_erpagendamentos_backend_internal_domain_notification (+16 more)

### Community 19 - "UseCase"
Cohesion: 0.11
Nodes (16): PasswordResetCodeConfirmInput, PasswordResetRequestInput, Repository, SecretCipher, TokenManager, PasswordHasher, UseCase, hashResetToken() (+8 more)

### Community 20 - "auth_cookie.go"
Cohesion: 0.11
Nodes (28): lazyAppHandler, writeLazyJSON(), bearerAuthAllowed(), authCookieDomain(), authCookiePartitioned(), authCookieSameSite(), authCookieSecure(), authTokenFromCookie() (+20 more)

### Community 21 - "database/sql.DB"
Cohesion: 0.12
Nodes (33): applyDatabaseSQL(), resolveSchemaFile(), initializeRepositories(), postgresRepositories(), NewAppointmentRepository(), NewAuditRepository(), NewAuthRepository(), NewCalendarRepository() (+25 more)

### Community 22 - "resource_test.go"
Cohesion: 0.54
Nodes (7): UseCase, newAppointment(), newResourceUseCase(), TestCreateAutoAssignsAvailableResource(), TestCreateRejectsBusyExplicitResource(), TestCreateRejectsUnknownResource(), TestUpdateAutoReassignsFreedResource()

### Community 23 - "Customer"
Cohesion: 0.15
Nodes (8): Customer, onlyDigits(), scanCustomerRow(), scanCustomerRows(), BookingResult, CustomerRepository, CustomerRepository, customerScanner

### Community 24 - "UseCase"
Cohesion: 0.08
Nodes (11): UseCase, NewUseCase(), CommissionRecord, CommissionRule, Payout, Repository, Summary, NewCommissionHandler() (+3 more)

### Community 25 - "CommissionsPage.tsx"
Cohesion: 0.06
Nodes (41): CommissionsPage, createMemberPlan(), createPackageCredit(), fetchCommissionPayouts(), fetchCommissionRecords(), fetchCommissionRules(), fetchCommissionSummaries(), fetchGiftCards() (+33 more)

### Community 26 - "Resource"
Cohesion: 0.13
Nodes (8): UseCase, Resource, NewResourceRepository(), scanResourceRow(), NewResourceHandler(), ResourceRepository, ResourceRepository, resourceScanner

### Community 27 - "CashRegister"
Cohesion: 0.13
Nodes (7): UseCase, CashRegister, CashRegisterEntry, ShiftSummary, NewCashRegisterHandler(), CashRegisterRepository, CashRegisterRepository

### Community 28 - "Professional"
Cohesion: 0.18
Nodes (6): Professional, scanProfessionalRow(), scanProfessionalRows(), ProfessionalRepository, ProfessionalRepository, professionalScanner

### Community 29 - "Queue"
Cohesion: 0.10
Nodes (8): Queue, NewQueue(), Job, Repository, sync.WaitGroup, HandlerFunc, JobRepository, JobRepository

### Community 30 - "Entry"
Cohesion: 0.12
Nodes (12): UseCase, NewUseCase(), Entry, NewWaitingListRepository(), scanWaitingRow(), scanWaitingRows(), NewWaitingListHandler(), WaitingListRepository (+4 more)

### Community 31 - "devDependencies"
Cohesion: 0.06
Nodes (34): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+26 more)

### Community 32 - "devDependencies"
Cohesion: 0.06
Nodes (34): dependencies, lucide-react, react, react-dom, react-router-dom, devDependencies, tailwindcss, @tailwindcss/vite (+26 more)

### Community 33 - "Container"
Cohesion: 0.13
Nodes (10): UseCase, UseCase, UseCase, Container, NewAppointmentHandler(), NewCalendarHandler(), NewConfigurationHandler(), NewExportHandler() (+2 more)

### Community 34 - ".Update"
Cohesion: 0.25
Nodes (3): RecurrenceInput, parseRFC3339(), validStatusChange()

### Community 35 - "auth/usecase.go"
Cohesion: 0.08
Nodes (19): loadEnvFile(), secureTemporaryPassword(), go_pkg_bufio, go_pkg_bytes, go_pkg_crypto_aes, go_pkg_crypto_cipher, go_pkg_crypto_rand, go_pkg_crypto_rsa (+11 more)

### Community 36 - "UseCase"
Cohesion: 0.12
Nodes (8): UseCase, NewUseCase(), Connection, Repository, SyncEvent, NewCalendarSyncHandler(), CalendarRepository, CalendarRepository

### Community 37 - "RouterConfig"
Cohesion: 0.20
Nodes (30): requireAdmin(), NewBookingHandler(), NewRouter(), startBackgroundWorkers(), registerAppointmentRoutes(), registerConfigurationRoutes(), registerCustomerRoutes(), registerProfessionalRoutes() (+22 more)

### Community 38 - "devDependencies"
Cohesion: 0.06
Nodes (32): dependencies, lucide-react, react, react-dom, devDependencies, tailwindcss, @tailwindcss/vite, @types/node (+24 more)

### Community 39 - "booking/src/pages/WaitingListPage.tsx"
Cohesion: 0.15
Nodes (20): Button, ButtonProps, ButtonSize, ButtonVariant, sizeClasses, variantClasses, CheckboxProps, Field() (+12 more)

### Community 40 - "MemberPlan"
Cohesion: 0.17
Nodes (5): MemberPlan, nullIfEmpty(), scanMemberPlans(), MemberPlanRepository, MemberPlanRepository

### Community 41 - "writeJSON"
Cohesion: 0.23
Nodes (4): writeBookingError(), writeJSON(), BookingHandler, CashRegisterHandler

### Community 42 - "Configuration"
Cohesion: 0.18
Nodes (11): workingHoursFor(), Configuration, cloneConfiguration(), defaultMemoryConfiguration(), NewConfigurationRepository(), normalizeConfigurationTenantID(), normalizePostgresTenantID(), ClosedDay (+3 more)

### Community 43 - "NewContainer"
Cohesion: 0.13
Nodes (31): Repository, NewUseCase(), NewUseCase(), NewUseCase(), UseCase, NewUseCase(), NewUseCase(), NewUseCase() (+23 more)

### Community 44 - "BookingPage.tsx"
Cohesion: 0.13
Nodes (18): formatDateLong(), formatDateTime(), formatMoney(), formatTime(), fromDateKey(), SHORT_DAY_LABELS, toDateKey(), Confirmation() (+10 more)

### Community 45 - "AdminShell.tsx"
Cohesion: 0.18
Nodes (13): AdminRoutes(), AdminShell(), AdminShellProps, AdminNavigationItem, navigationGroups, navigationItems, AdminRole, roleLabel() (+5 more)

### Community 46 - "InventoryPage"
Cohesion: 0.20
Nodes (14): fetchServices(), fetchServiceSupplies(), fetchStockMovements(), setServiceSupplies(), InventoryPage(), handleAddSupply(), handleDeleteProduct(), handleRemoveSupply() (+6 more)

### Community 47 - "AuthRepository"
Cohesion: 0.14
Nodes (6): Account, PasswordResetToken, UpdateProfileInput, cloneAuthAccount(), resetTokenExpired(), AuthRepository

### Community 48 - "UseCase"
Cohesion: 0.19
Nodes (8): appointmentDurationMinutes(), generateManagementCode(), UseCase, hashManagementCode(), BookingInput, PortalAppointment, PortalCustomer, PortalResult

### Community 49 - "AppointmentHandler"
Cohesion: 0.18
Nodes (6): writeAppointmentError(), resolveActorTenantID(), writeAdminAuthError(), writeBadRequestError(), writeConflictError(), AppointmentHandler

### Community 50 - "bookingApi.ts"
Cohesion: 0.09
Nodes (37): backendOrigin, buildApiUrl(), AvailabilityResponse, BookingResult, cancelPublicAppointment(), fetchAvailability(), fetchPublicAppointment(), fetchPublicInfo() (+29 more)

### Community 51 - "memoryRepositories"
Cohesion: 0.14
Nodes (22): TestCashRegisterFlow(), TestCheckoutFlow(), TestCommissionUseCase(), TestFinanceUseCase(), TestOnboardingSignup(), TestStockUseCase(), memoryRepositories(), NewAppointmentRepository() (+14 more)

### Community 52 - "Membership"
Cohesion: 0.16
Nodes (3): Membership, MembershipRepository, MembershipRepository

### Community 53 - "Package"
Cohesion: 0.16
Nodes (3): Package, PackageRepository, PackageRepository

### Community 54 - "UseCase"
Cohesion: 0.10
Nodes (10): UseCase, Summary, Transaction, NewFinanceHandler(), NewMercadoPagoWebhookHandler(), CategoryTotal, CategoryTotal, FinanceRepository (+2 more)

### Community 55 - "auth/entities.go"
Cohesion: 0.19
Nodes (9): ChangePasswordInput, LoginInput, PasswordResetConfirmInput, RegisterInput, Session, ValidationError, NewValidationError(), validatePassword() (+1 more)

### Community 56 - "PackageCredit"
Cohesion: 0.18
Nodes (5): PackageCredit, scanPackageCredits(), database/sql.Rows, PackageCreditRepository, PackageCreditRepository

### Community 57 - "GiftCard"
Cohesion: 0.14
Nodes (5): generateGiftCardCode(), GiftCard, NewGiftCardRepository(), GiftCardRepository, GiftCardRepository

### Community 58 - "UseCase"
Cohesion: 0.17
Nodes (9): UseCase, NewUseCase(), formatAppointmentLocal(), UseCase, NewReminderWorker(), NewCustomerHandler(), Repository, context.Context (+1 more)

### Community 59 - "BookingPage"
Cohesion: 0.25
Nodes (3): createPublicBooking(), BookingPage(), handleBookingSubmit()

### Community 60 - "AuditLog"
Cohesion: 0.14
Nodes (13): AuditLog, Repository, UseCase, NewUseCase(), TestAuditLogRecordingAndListing(), TestOutboxPublishAndDispatch(), NewAuditRepository(), NewOutboxRepository() (+5 more)

### Community 61 - "go_pkg_errors"
Cohesion: 0.12
Nodes (9): CheckoutInput, CheckoutResult, ProductItem, go_pkg_erpagendamentos_backend_internal_domain_calendar, go_pkg_erpagendamentos_backend_internal_domain_cashregister, go_pkg_erpagendamentos_backend_internal_domain_commission, go_pkg_erpagendamentos_backend_internal_domain_finance, go_pkg_erpagendamentos_backend_internal_domain_stock (+1 more)

### Community 62 - "EventBroker"
Cohesion: 0.23
Nodes (8): NewCheckoutHandler(), EventBroker, NewEventBroker(), RegisterMetricsBroker(), sync.RWMutex, CheckoutHandler, Event, MetricsRegistry

### Community 63 - "rate_limit.go"
Cohesion: 0.08
Nodes (31): UseCase, NewUseCase(), Log, NewNotificationRepository(), buildCompoundIdempotencyKey(), NewNotificationHandler(), forwardedClientIP(), isSensitiveRateLimitedRoute() (+23 more)

### Community 64 - "AuthRepository"
Cohesion: 0.12
Nodes (4): AuthState, updateAuthPasswordAndVersion(), database/sql.Tx, AuthRepository

### Community 65 - "Appointment"
Cohesion: 0.16
Nodes (5): Appointment, StatusCount, parseAppointmentBound(), parseAppointmentInterval(), AppointmentRepository

### Community 66 - "writeInvalidPayload"
Cohesion: 0.26
Nodes (3): writeCommerceError(), writeInvalidPayload(), CommerceHandler

### Community 67 - "testing.T"
Cohesion: 0.23
Nodes (15): UseCase, newTestUseCase(), TestAvailableSlotsRespectsClosedDays(), TestCancelPublic(), TestCreateBookingGeneratesManagementCode(), TestManagementCodeNormalization(), TestPortalCancelValidatesOwnership(), TestPortalLookup() (+7 more)

### Community 68 - "templates.go"
Cohesion: 0.27
Nodes (17): greeting(), PublicBookingOrigin(), RenderCancellation(), RenderConfirmation(), RenderReminder(), RenderReschedule(), sampleData(), TestGreetingWithoutName() (+9 more)

### Community 69 - "withIdempotency"
Cohesion: 0.15
Nodes (12): isHopByHopHeader(), NewIdempotencyStore(), TestIdempotencyBlocksInFlightConcurrentRequests(), TestIdempotencyDoesNotCacheServerFailures(), TestIdempotencyIgnoredOnGetRequests(), TestIdempotencyReplaysCompletedResponse(), withIdempotency(), bytes.Buffer (+4 more)

### Community 70 - "compilerOptions"
Cohesion: 0.09
Nodes (21): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+13 more)

### Community 71 - "AppointmentRepository"
Cohesion: 0.19
Nodes (5): appointmentSelectColumns(), scanAppointmentRow(), scanAppointmentRows(), AppointmentRepository, appointmentScanner

### Community 72 - "compilerOptions"
Cohesion: 0.09
Nodes (21): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+13 more)

### Community 73 - "compilerOptions"
Cohesion: 0.09
Nodes (21): compilerOptions, allowJs, allowSyntheticDefaultImports, esModuleInterop, forceConsistentCasingInFileNames, isolatedModules, jsx, lib (+13 more)

### Community 74 - "writeNotFoundError"
Cohesion: 0.18
Nodes (6): writeCustomerError(), writeNotificationError(), writeProfessionalError(), writeNotFoundError(), CustomerHandler, ProfessionalHandler

### Community 75 - "OutboxRepository"
Cohesion: 0.18
Nodes (4): quoteLiteral(), OutboxRepository, OutboxEvent, OutboxRepository

### Community 76 - "MembershipsPage"
Cohesion: 0.36
Nodes (8): fetchMemberPlans(), fetchMemberships(), MembershipsPage(), handleDelete(), handlePlanStatusChange(), handleToggleActive(), loadMemberships(), planStatusTone()

### Community 77 - "totp.go"
Cohesion: 0.16
Nodes (10): MFAStatus, generateTOTPSecret(), totpCode(), totpURI(), validateTOTP(), isMFAEligibleRole(), go_pkg_crypto_hmac, go_pkg_crypto_sha1 (+2 more)

### Community 78 - "admin/src/App.tsx"
Cohesion: 0.22
Nodes (13): App(), handleLogout(), ResetPasswordPage(), handleSubmit(), AdminAuthError, confirmAdminPasswordReset(), logoutAdmin(), subscribeToAdminAuthChanges() (+5 more)

### Community 79 - "UseCase"
Cohesion: 0.22
Nodes (7): UseCase, NewUseCase(), NewSystemHandler(), go_pkg_runtime, EventDispatcher, Repository, SystemHandler

### Community 80 - "ERP Agendamentos"
Cohesion: 0.14
Nodes (13): Banco de Dados & Schema, CI, Endpoints principais, ERP Agendamentos, Estrutura, Pré-requisitos, Rodando com Docker (produção local), Rodando em desenvolvimento (sem banco) (+5 more)

### Community 81 - "time.Duration"
Cohesion: 0.16
Nodes (4): CleanupSummary, DeviceSession, sessionExpired(), time.Duration

### Community 82 - "UseCase"
Cohesion: 0.22
Nodes (8): UseCase, NewUseCase(), NewCommerceHandler(), GiftCardRepository, MemberPlanRepository, MembershipRepository, PackageCreditRepository, PackageRepository

### Community 83 - "email_provider.go"
Cohesion: 0.21
Nodes (8): buildNotificationProvider(), intEnv(), logFallback(), NewEmailProvider(), NewLogProvider(), go_pkg_net_smtp, EmailProvider, LogProvider

### Community 84 - "admin/src/lib/apiError.ts"
Cohesion: 0.27
Nodes (8): ApiError, ApiErrorPayload, buildApiError(), parseApiResponse(), API_ERROR_MESSAGES, FIELD_ERROR_MESSAGES, getApiErrorMessage(), translateFieldErrors()

### Community 85 - "booking/src/lib/apiError.ts"
Cohesion: 0.27
Nodes (8): ApiError, ApiErrorPayload, buildApiError(), parseApiResponse(), API_ERROR_MESSAGES, FIELD_ERROR_MESSAGES, getApiErrorMessage(), translateFieldErrors()

### Community 86 - "Claims"
Cohesion: 0.33
Nodes (5): isProductionAuthEnvironment(), Claims, decodeSegment(), encodeSegment(), JWTManager

### Community 88 - "FinancePage"
Cohesion: 0.27
Nodes (13): fetchCashRegisterHistory(), fetchFinancialTransactions(), FinancePage(), handleCloseShift(), handleCreate(), handleCreateEntry(), handleDelete(), handleMarkPaid() (+5 more)

### Community 89 - ".UpdateAccountRole"
Cohesion: 0.39
Nodes (4): canManageTenantAccount(), isPrivilegedAdminRole(), isStaffRole(), normalizeTenantScope()

### Community 90 - "SecretCipher"
Cohesion: 0.38
Nodes (3): NewSecretCipher(), crypto/cipher.AEAD, SecretCipher

### Community 91 - "observability.go"
Cohesion: 0.17
Nodes (12): RecordRequestMetric(), isSensitiveLogKey(), redactSensitiveText(), resolveRequestID(), sanitizeLogString(), sanitizeLogValue(), truncateLogString(), withObservability() (+4 more)

### Community 92 - "admin/vite.config.ts"
Cohesion: 0.60
Nodes (3): ref_tailwindcss_vite, ref_vite, ref_vitejs_plugin_react

### Community 93 - "appointment/repository.go"
Cohesion: 0.53
Nodes (4): DayMetric, ProfessionalMetric, ServiceMetric, Analytics

### Community 94 - "NewBcryptHasher"
Cohesion: 0.40
Nodes (3): NewBcryptHasher(), go_pkg_golang_org_x_crypto_bcrypt, BcryptHasher

### Community 95 - "time.Time"
Cohesion: 0.29
Nodes (11): countOverlappingSlots(), formatAppointmentForTemplate(), isClosedDay(), localMidnight(), overlapsAnyBusy(), parseClock(), parseRFC3339(), ProfessionalSlots (+3 more)

### Community 96 - "TenantsPage"
Cohesion: 0.24
Nodes (8): fetchTenants(), TenantsPage(), handleCheckReadiness(), handleDelete(), handleToggleActive(), loadAuditLogs(), loadSystemStatus(), loadTenants()

### Community 97 - "adminHttp.ts"
Cohesion: 0.20
Nodes (9): AdminAccount, requestAdminPasswordReset(), adminGetCache, AdminMeResponse, backendOrigin, buildApiUrl(), AdminLoginPage(), handleRecovery() (+1 more)

### Community 98 - "ical.go"
Cohesion: 0.83
Nodes (3): BuildICalCalendar(), formatUTC(), sanitizeLine()

### Community 99 - "booking/src/App.tsx"
Cohesion: 0.20
Nodes (8): App(), booking_src_index, fetchPublicReview(), submitPublicReview(), initWidgetBridge(), ReviewPage(), handleSubmit(), SlugEntryPage()

### Community 109 - "AppProviders.tsx"
Cohesion: 0.32
Nodes (5): AppProviders(), AppProvidersProps, AppRouter(), admin_src_index, ToastProvider()

### Community 110 - "MaintenanceWorker"
Cohesion: 0.39
Nodes (5): MaintenanceWorker, MaintenanceConfig, UseCase, NewMaintenanceWorker(), authComponents

### Community 115 - "PortalPage.tsx"
Cohesion: 0.25
Nodes (6): PortalAppointment, PortalResult, Slot, ACTIONABLE_STATUSES, PortalStep, STATUS_LABELS

### Community 116 - "main.go"
Cohesion: 0.38
Nodes (6): newLazyAppHandler(), main(), resolveDurationEnv(), resolvePositiveIntEnv(), go_pkg_os_signal, go_pkg_syscall

### Community 118 - "TenantSubscriptionModal"
Cohesion: 0.40
Nodes (3): updateTenantSubscription(), TenantSubscriptionModal(), handleSubmit()

### Community 119 - "tx.go"
Cohesion: 0.40
Nodes (4): NewPostgresTxRunner(), NopTxRunner, PostgresTxRunner, TxRunner

### Community 120 - "TeamPage"
Cohesion: 0.60
Nodes (5): fetchAdminUsers(), TeamPage(), handleDelete(), handleRoleChange(), loadUsers()

### Community 121 - "widget_handler.go"
Cohesion: 0.60
Nodes (3): NewWidgetHandler(), resolveWidgetBookingOrigin(), WidgetHandler

### Community 122 - "body_limit.go"
Cohesion: 0.83
Nodes (3): resolveMaxRequestBodyBytes(), shouldLimitRequestBody(), withRequestBodyLimit()

## Knowledge Gaps
- **256 isolated node(s):** `graphify-mcp`, `name`, `private`, `version`, `type` (+251 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **14 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Container` connect `Container` to `Payment`, `go_pkg_net_http`, `Review`, `UseCase`, `Tenant`, `Service`, `UseCase`, `auth_cookie.go`, `database/sql.DB`, `UseCase`, `Resource`, `CashRegister`, `Queue`, `Entry`, `UseCase`, `NewContainer`, `UseCase`, `UseCase`, `UseCase`, `AuditLog`, `EventBroker`, `rate_limit.go`, `UseCase`, `UseCase`, `MaintenanceWorker`?**
  _High betweenness centrality (0.054) - this node is a cross-community bridge._
- **Why does `UseCase` connect `UseCase` to `Container`, `auth/usecase.go`, `net/http.ResponseWriter`, `RouterConfig`, `withIdempotency`, `writeAPIError`, `NewContainer`, `totp.go`, `net/http.Request`, `AuthRepository`, `MaintenanceWorker`, `time.Duration`, `AppointmentHandler`, `Claims`, `auth/entities.go`, `.UpdateAccountRole`, `rate_limit.go`, `time.Time`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **Why does `RouterConfig` connect `RouterConfig` to `Payment`, `go_pkg_net_http`, `Review`, `UseCase`, `Tenant`, `Service`, `UseCase`, `UseCase`, `Resource`, `CashRegister`, `Entry`, `Container`, `UseCase`, `NewContainer`, `UseCase`, `UseCase`, `UseCase`, `AuditLog`, `EventBroker`, `rate_limit.go`, `UseCase`, `time.Duration`, `UseCase`, `MaintenanceWorker`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **Are the 130 inferred relationships involving `writeJSON()` (e.g. with `health()` and `readiness()`) actually correct?**
  _`writeJSON()` has 130 INFERRED edges - model-reasoned connections that need verification._
- **Are the 103 inferred relationships involving `resolveRequestTenantID()` (e.g. with `.Analytics()` and `.CheckIn()`) actually correct?**
  _`resolveRequestTenantID()` has 103 INFERRED edges - model-reasoned connections that need verification._
- **What connects `graphify-mcp`, `name`, `private` to the rest of the system?**
  _256 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `adminApi.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.05353535353535353 - nodes in this community are weakly interconnected._
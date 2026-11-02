package bootstrap

import (
	"database/sql"
	"log"
	"os"
	"strings"

	appointmentdomain "erpagendamentos/backend/internal/domain/appointment"
	authdomain "erpagendamentos/backend/internal/domain/auth"
	calendardomain "erpagendamentos/backend/internal/domain/calendar"
	cashregisterdomain "erpagendamentos/backend/internal/domain/cashregister"
	commercedomain "erpagendamentos/backend/internal/domain/commerce"
	commissiondomain "erpagendamentos/backend/internal/domain/commission"
	configurationdomain "erpagendamentos/backend/internal/domain/configuration"
	customerdomain "erpagendamentos/backend/internal/domain/customer"
	financedomain "erpagendamentos/backend/internal/domain/finance"
	jobdomain "erpagendamentos/backend/internal/domain/job"
	notificationdomain "erpagendamentos/backend/internal/domain/notification"
	paymentdomain "erpagendamentos/backend/internal/domain/payment"
	professionaldomain "erpagendamentos/backend/internal/domain/professional"
	resourcedomain "erpagendamentos/backend/internal/domain/resource"
	reviewdomain "erpagendamentos/backend/internal/domain/review"
	servicedomain "erpagendamentos/backend/internal/domain/service"
	stockdomain "erpagendamentos/backend/internal/domain/stock"
	tenantdomain "erpagendamentos/backend/internal/domain/tenant"
	waitinglistdomain "erpagendamentos/backend/internal/domain/waitinglist"
	"erpagendamentos/backend/internal/infrastructure/persistence/memory"
	"erpagendamentos/backend/internal/infrastructure/persistence/postgres"
)

type repositorySet struct {
	auth           authdomain.Repository
	tenant         tenantdomain.Repository
	professional   professionaldomain.Repository
	service        servicedomain.Repository
	resource       resourcedomain.Repository
	customer       customerdomain.Repository
	appointment    appointmentdomain.Repository
	configuration  configurationdomain.Repository
	notification   notificationdomain.Repository
	payment        paymentdomain.Repository
	waitingList    waitinglistdomain.Repository
	review         reviewdomain.Repository
	packages       commercedomain.PackageRepository
	packageCredits commercedomain.PackageCreditRepository
	giftCards      commercedomain.GiftCardRepository
	memberships    commercedomain.MembershipRepository
	memberPlans    commercedomain.MemberPlanRepository
	stock          stockdomain.Repository
	commission     commissiondomain.Repository
	finance        financedomain.Repository
	calendar       calendardomain.Repository
	job            jobdomain.Repository
	cashRegister   cashregisterdomain.Repository
}

func initializeRepositories() (repositorySet, *sql.DB, bool) {
	repositories := memoryRepositories()
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	databaseRequired := databaseURL != ""
	if databaseURL == "" {
		if isProductionEnv() {
			log.Println(`{"level":"WARN","message":"DATABASE_URL is not configured; running with in-memory persistence until database URL is provided."}`)
		}
		return repositories, nil, false
	}

	db, err := postgres.Open(databaseURL)
	if err != nil {
		if isProductionEnv() {
			log.Fatalf("postgres is required in production: %v", err)
		}
		log.Printf("postgres disabled: %v", err)
		log.Printf(`{"level":"ERROR","message":"database is unavailable","alert":"database_unavailable","error":"%v"}`, err)
		return repositories, nil, databaseRequired
	}
	if err := applyDatabaseSQL(db); err != nil {
		_ = db.Close()
		if isProductionEnv() {
			log.Fatalf("database migrations failed in production: %v", err)
		}
		log.Printf("postgres disabled after migration error: %v", err)
		log.Printf(`{"level":"ERROR","message":"database is unavailable","alert":"database_unavailable","error":"%v"}`, err)
		return repositories, nil, databaseRequired
	}

	log.Println("postgres repositories enabled")
	return postgresRepositories(db), db, databaseRequired
}

func memoryRepositories() repositorySet {
	return repositorySet{
		auth:           memory.NewAuthRepository(),
		tenant:         memory.NewTenantRepository(),
		professional:   memory.NewProfessionalRepository(),
		service:        memory.NewServiceRepository(),
		resource:       memory.NewResourceRepository(),
		customer:       memory.NewCustomerRepository(),
		appointment:    memory.NewAppointmentRepository(),
		configuration:  memory.NewConfigurationRepository(),
		notification:   memory.NewNotificationRepository(),
		payment:        memory.NewPaymentRepository(),
		waitingList:    memory.NewWaitingListRepository(),
		review:         memory.NewReviewRepository(),
		packages:       memory.NewPackageRepository(),
		packageCredits: memory.NewPackageCreditRepository(),
		giftCards:      memory.NewGiftCardRepository(),
		memberships:    memory.NewMembershipRepository(),
		memberPlans:    memory.NewMemberPlanRepository(),
		stock:          memory.NewStockRepository(),
		commission:     memory.NewCommissionRepository(),
		finance:        memory.NewFinanceRepository(),
		calendar:       memory.NewCalendarRepository(),
		job:            memory.NewJobRepository(),
		cashRegister:   memory.NewCashRegisterRepository(),
	}
}

func postgresRepositories(db *sql.DB) repositorySet {
	return repositorySet{
		auth:           postgres.NewAuthRepository(db),
		tenant:         postgres.NewTenantRepository(db),
		professional:   postgres.NewProfessionalRepository(db),
		service:        postgres.NewServiceRepository(db),
		resource:       postgres.NewResourceRepository(db),
		customer:       postgres.NewCustomerRepository(db),
		appointment:    postgres.NewAppointmentRepository(db),
		configuration:  postgres.NewConfigurationRepository(db),
		notification:   postgres.NewNotificationRepository(db),
		payment:        postgres.NewPaymentRepository(db),
		waitingList:    postgres.NewWaitingListRepository(db),
		review:         postgres.NewReviewRepository(db),
		packages:       postgres.NewPackageRepository(db),
		packageCredits: postgres.NewPackageCreditRepository(db),
		giftCards:      postgres.NewGiftCardRepository(db),
		memberships:    postgres.NewMembershipRepository(db),
		memberPlans:    postgres.NewMemberPlanRepository(db),
		stock:          postgres.NewStockRepository(db),
		commission:     postgres.NewCommissionRepository(db),
		finance:        postgres.NewFinanceRepository(db),
		calendar:       postgres.NewCalendarRepository(db),
		job:            postgres.NewJobRepository(db),
		cashRegister:   postgres.NewCashRegisterRepository(db),
	}
}

package bootstrap

import (
	"database/sql"
	"log"
	"os"
	"strings"

	appointmentdomain "erpagendamentos/backend/internal/domain/appointment"
	authdomain "erpagendamentos/backend/internal/domain/auth"
	commercedomain "erpagendamentos/backend/internal/domain/commerce"
	configurationdomain "erpagendamentos/backend/internal/domain/configuration"
	customerdomain "erpagendamentos/backend/internal/domain/customer"
	notificationdomain "erpagendamentos/backend/internal/domain/notification"
	paymentdomain "erpagendamentos/backend/internal/domain/payment"
	professionaldomain "erpagendamentos/backend/internal/domain/professional"
	servicedomain "erpagendamentos/backend/internal/domain/service"
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
	customer       customerdomain.Repository
	appointment    appointmentdomain.Repository
	configuration  configurationdomain.Repository
	notification   notificationdomain.Repository
	payment        paymentdomain.Repository
	waitingList    waitinglistdomain.Repository
	packages       commercedomain.PackageRepository
	packageCredits commercedomain.PackageCreditRepository
	giftCards      commercedomain.GiftCardRepository
	memberships    commercedomain.MembershipRepository
	memberPlans    commercedomain.MemberPlanRepository
}

func initializeRepositories() (repositorySet, *sql.DB, bool) {
	repositories := memoryRepositories()
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	databaseRequired := databaseURL != ""
	if isProductionEnv() && databaseURL == "" {
		log.Fatal("DATABASE_URL is required when ERPAGENDAMENTOS_ENV=production")
	}
	if databaseURL == "" {
		return repositories, nil, databaseRequired
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
		customer:       memory.NewCustomerRepository(),
		appointment:    memory.NewAppointmentRepository(),
		configuration:  memory.NewConfigurationRepository(),
		notification:   memory.NewNotificationRepository(),
		payment:        memory.NewPaymentRepository(),
		waitingList:    memory.NewWaitingListRepository(),
		packages:       memory.NewPackageRepository(),
		packageCredits: memory.NewPackageCreditRepository(),
		giftCards:      memory.NewGiftCardRepository(),
		memberships:    memory.NewMembershipRepository(),
		memberPlans:    memory.NewMemberPlanRepository(),
	}
}

func postgresRepositories(db *sql.DB) repositorySet {
	return repositorySet{
		auth:           postgres.NewAuthRepository(db),
		tenant:         postgres.NewTenantRepository(db),
		professional:   postgres.NewProfessionalRepository(db),
		service:        postgres.NewServiceRepository(db),
		customer:       postgres.NewCustomerRepository(db),
		appointment:    postgres.NewAppointmentRepository(db),
		configuration:  postgres.NewConfigurationRepository(db),
		notification:   postgres.NewNotificationRepository(db),
		payment:        postgres.NewPaymentRepository(db),
		waitingList:    postgres.NewWaitingListRepository(db),
		packages:       postgres.NewPackageRepository(db),
		packageCredits: postgres.NewPackageCreditRepository(db),
		giftCards:      postgres.NewGiftCardRepository(db),
		memberships:    postgres.NewMembershipRepository(db),
		memberPlans:    postgres.NewMemberPlanRepository(db),
	}
}

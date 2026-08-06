package bootstrap

import (
	"database/sql"
	"log"
	"os"
	"strings"

	appointmentdomain "erpagendamentos/backend/internal/domain/appointment"
	authdomain "erpagendamentos/backend/internal/domain/auth"
	configurationdomain "erpagendamentos/backend/internal/domain/configuration"
	customerdomain "erpagendamentos/backend/internal/domain/customer"
	professionaldomain "erpagendamentos/backend/internal/domain/professional"
	servicedomain "erpagendamentos/backend/internal/domain/service"
	tenantdomain "erpagendamentos/backend/internal/domain/tenant"
	"erpagendamentos/backend/internal/infrastructure/persistence/memory"
	"erpagendamentos/backend/internal/infrastructure/persistence/postgres"
)

type repositorySet struct {
	auth          authdomain.Repository
	tenant        tenantdomain.Repository
	professional  professionaldomain.Repository
	service       servicedomain.Repository
	customer      customerdomain.Repository
	appointment   appointmentdomain.Repository
	configuration configurationdomain.Repository
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
		auth:          memory.NewAuthRepository(),
		tenant:        memory.NewTenantRepository(),
		professional:  memory.NewProfessionalRepository(),
		service:       memory.NewServiceRepository(),
		customer:      memory.NewCustomerRepository(),
		appointment:   memory.NewAppointmentRepository(),
		configuration: memory.NewConfigurationRepository(),
	}
}

func postgresRepositories(db *sql.DB) repositorySet {
	return repositorySet{
		auth:          postgres.NewAuthRepository(db),
		tenant:        postgres.NewTenantRepository(db),
		professional:  postgres.NewProfessionalRepository(db),
		service:       postgres.NewServiceRepository(db),
		customer:      postgres.NewCustomerRepository(db),
		appointment:   postgres.NewAppointmentRepository(db),
		configuration: postgres.NewConfigurationRepository(db),
	}
}

package onboarding_test

import (
	"testing"
	"time"

	appauth "erpagendamentos/backend/internal/application/auth"
	"erpagendamentos/backend/internal/application/onboarding"
	"erpagendamentos/backend/internal/infrastructure/persistence/memory"
	"erpagendamentos/backend/internal/infrastructure/security"
)

func TestOnboardingSignup(t *testing.T) {
	tenantRepo := memory.NewTenantRepository()
	authRepo := memory.NewAuthRepository()
	configRepo := memory.NewConfigurationRepository()
	serviceRepo := memory.NewServiceRepository()
	proRepo := memory.NewProfessionalRepository()

	hasher := security.NewBcryptHasher(4)
	tokens := security.NewJWTManager("test-jwt-secret-key-32-chars-long", time.Hour)
	authUseCase := appauth.NewUseCase(authRepo, hasher, tokens)

	onboardingUseCase := onboarding.NewUseCase(
		tenantRepo,
		authUseCase,
		configRepo,
		serviceRepo,
		proRepo,
	)

	// 1. Success Signup
	session, tenant, err := onboardingUseCase.Signup(onboarding.SignupInput{
		BusinessName: "Barbearia Don Corleone",
		Slug:         "don-corleone",
		AdminName:    "Vito Corleone",
		Email:        "vito@corleone.com",
		Password:     "Senha@1234",
		Phone:        "11999998888",
		Segment:      "barbearia",
	})
	if err != nil {
		t.Fatalf("expected successful signup, got: %v", err)
	}

	if tenant.Slug != "don-corleone" {
		t.Errorf("expected slug 'don-corleone', got '%s'", tenant.Slug)
	}
	if session.User.Email != "vito@corleone.com" {
		t.Errorf("expected email 'vito@corleone.com', got '%s'", session.User.Email)
	}
	if session.User.Role != "owner" {
		t.Errorf("expected role 'owner', got '%s'", session.User.Role)
	}

	// Verify services were seeded
	services, err := serviceRepo.ListByTenant(tenant.ID)
	if err != nil || len(services) < 2 {
		t.Errorf("expected seeded services for barbearia, got %d", len(services))
	}

	// 2. Duplicate slug should fail
	_, _, err = onboardingUseCase.Signup(onboarding.SignupInput{
		BusinessName: "Outra Barbearia",
		Slug:         "don-corleone",
		AdminName:    "Outro Don",
		Email:        "outro@corleone.com",
		Password:     "Senha@1234",
	})
	if err == nil {
		t.Errorf("expected error for duplicate slug, got nil")
	}

	// 3. Short password should fail
	_, _, err = onboardingUseCase.Signup(onboarding.SignupInput{
		BusinessName: "Salão Beleza Pura",
		Slug:         "beleza-pura",
		AdminName:    "Maria Silva",
		Email:        "maria@beleza.com",
		Password:     "curta",
	})
	if err == nil {
		t.Errorf("expected error for short password, got nil")
	}

	// 4. Signup without businessName and slug should auto-generate establishment
	sessionMinimal, tenantMinimal, err := onboardingUseCase.Signup(onboarding.SignupInput{
		AdminName: "Carlos Barbeiro",
		Email:     "carlos@barbeiro.com",
		Password:  "Senha@1234",
		Phone:     "11988887777",
	})
	if err != nil {
		t.Fatalf("expected successful minimal signup, got: %v", err)
	}
	if tenantMinimal.Name != "Estabelecimento de Carlos Barbeiro" {
		t.Errorf("expected auto generated name, got: %s", tenantMinimal.Name)
	}
	if tenantMinimal.Slug == "" {
		t.Errorf("expected non-empty auto generated slug")
	}
	if sessionMinimal.User.Email != "carlos@barbeiro.com" {
		t.Errorf("expected email carlos@barbeiro.com, got: %s", sessionMinimal.User.Email)
	}

	// 5. Signup with OwnerName instead of AdminName should succeed
	sessionOwner, tenantOwner, err := onboardingUseCase.Signup(onboarding.SignupInput{
		OwnerName: "Juliana Silva",
		Email:     "juliana@silva.com",
		Password:  "Senha@1234",
		Phone:     "11988887777",
	})
	if err != nil {
		t.Fatalf("expected successful signup with OwnerName, got: %v", err)
	}
	if tenantOwner.Name != "Estabelecimento de Juliana Silva" {
		t.Errorf("expected name 'Estabelecimento de Juliana Silva', got: %s", tenantOwner.Name)
	}
	if sessionOwner.User.Name != "Juliana Silva" {
		t.Errorf("expected user name 'Juliana Silva', got: %s", sessionOwner.User.Name)
	}

	// 6. Duplicate email should fail with translated error message
	_, _, err = onboardingUseCase.Signup(onboarding.SignupInput{
		OwnerName: "Outro Don",
		Email:     "vito@corleone.com",
		Password:  "Senha@1234",
	})
	if err == nil {
		t.Errorf("expected error for duplicate email, got nil")
	}
}


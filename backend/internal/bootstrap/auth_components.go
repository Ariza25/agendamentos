package bootstrap

import (
	"log"
	"os"
	"strings"
	"time"

	"erpagendamentos/backend/internal/application/auth"
	authdomain "erpagendamentos/backend/internal/domain/auth"
	"erpagendamentos/backend/internal/infrastructure/security"
)

type authComponents struct {
	useCase     *auth.UseCase
	maintenance auth.MaintenanceConfig
	hasher      auth.PasswordHasher
}

func buildAuthComponents(authRepository authdomain.Repository) authComponents {
	jwtSecret := os.Getenv("JWT_SECRET")
	if isProductionEnv() && strings.TrimSpace(jwtSecret) == "" {
		log.Fatal("JWT_SECRET is required when ERPAGENDAMENTOS_ENV=production")
	}
	if jwtSecret == "" {
		jwtSecret = "erpagendamentos-dev-secret"
	}
	if isProductionEnv() && jwtSecret == "erpagendamentos-dev-secret" {
		log.Fatal("JWT_SECRET must not use the development default in production")
	}

	passwordHasher := security.NewBcryptHasher(12)
	tokenManager := security.NewJWTManager(jwtSecret, 7*24*time.Hour)
	mfaEncryptionKey := strings.TrimSpace(os.Getenv("MFA_ENCRYPTION_KEY"))
	if isProductionEnv() && mfaEncryptionKey == "" {
		log.Fatal("MFA_ENCRYPTION_KEY is required when ERPAGENDAMENTOS_ENV=production")
	}
	if mfaEncryptionKey == "" {
		mfaEncryptionKey = "erpagendamentos-dev-mfa-key"
	}
	mfaCipher, err := security.NewSecretCipher(mfaEncryptionKey)
	if err != nil {
		log.Fatalf("initialize mfa encryption: %v", err)
	}
	if err := seedAdminAccount(authRepository, passwordHasher); err != nil {
		log.Printf("admin seed skipped: %v", err)
	}

	return authComponents{
		useCase: auth.NewUseCase(
			authRepository,
			passwordHasher,
			tokenManager,
		).WithMFA(mfaCipher),
		maintenance: auth.MaintenanceConfig{
			Enabled:   boolEnv("FEATURE_AUTH_MAINTENANCE_ENABLED", true),
			Interval:  durationSecondsEnv("AUTH_MAINTENANCE_INTERVAL_SECONDS", time.Hour),
			Retention: durationSecondsEnv("AUTH_RESET_TOKEN_RETENTION_SECONDS", 7*24*time.Hour),
		},
		hasher: passwordHasher,
	}
}

package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	domain "erpagendamentos/backend/internal/domain/auth"
)

var (
	ErrEmailAlreadyRegistered = fmt.Errorf("email already registered")
	ErrInvalidCredentials     = fmt.Errorf("invalid credentials")
	ErrInvalidProfile         = fmt.Errorf("invalid profile")
	ErrInvalidRole            = fmt.Errorf("invalid role")
	ErrAccountNotFound        = fmt.Errorf("account not found")
	ErrLastAdmin              = fmt.Errorf("cannot remove the last admin")
	ErrInvalidResetToken      = fmt.Errorf("invalid password reset token")
	ErrMFARequired            = fmt.Errorf("mfa code required")
)

var passwordSpecialCharPattern = regexp.MustCompile(`[^A-Za-z0-9]`)

const (
	passwordResetCodeMaxAttempts       = 5
	passwordResetCodeWindow            = 15 * time.Minute
	passwordResetCodeRequestIPLimit    = 10
	passwordResetCodeRequestEmailLimit = 3
	passwordResetCodeConfirmIPLimit    = 30
)

type ValidationError struct {
	FieldErrors map[string]string
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

func NewValidationError() *ValidationError {
	return &ValidationError{FieldErrors: make(map[string]string)}
}

func (e *ValidationError) Add(field string, message string) {
	if _, exists := e.FieldErrors[field]; !exists {
		e.FieldErrors[field] = message
	}
}

func (e *ValidationError) HasErrors() bool {
	return len(e.FieldErrors) > 0
}

type PasswordHasher interface {
	Hash(value string) (string, error)
	Compare(hash string, value string) error
}

type TokenManager interface {
	Sign(claims domain.Claims) (string, error)
	Parse(token string) (domain.Claims, error)
}

type SecretCipher interface {
	Encrypt(value string) (string, error)
	Decrypt(value string) (string, error)
}

type UseCase struct {
	repository domain.Repository
	hasher     PasswordHasher
	tokens     TokenManager
	now        func() time.Time
	mfaCipher  SecretCipher
}

func (u *UseCase) WithMFA(cipher SecretCipher) *UseCase {
	u.mfaCipher = cipher
	return u
}

func NewUseCase(repository domain.Repository, hasher PasswordHasher, tokens TokenManager) *UseCase {
	return &UseCase{repository: repository, hasher: hasher, tokens: tokens, now: time.Now}
}

func (u *UseCase) ListAccounts() ([]domain.Account, error) {
	return u.repository.ListAccounts()
}

func (u *UseCase) CreateStaffAccount(input domain.RegisterInput) (domain.Account, error) {
	role := strings.TrimSpace(input.Role)
	if role == "" {
		role = "receptionist"
	}
	if !isStaffRole(role) {
		return domain.Account{}, ErrInvalidRole
	}
	session, err := u.registerWithRole(input, role)
	if err != nil {
		return domain.Account{}, err
	}
	return session.User, nil
}

func (u *UseCase) RegisterOwner(input domain.RegisterInput) (domain.Session, error) {
	input.Role = "owner"
	return u.registerWithRole(input, "owner")
}

func (u *UseCase) CreateOrUpdateTenantOwner(input domain.RegisterInput, tenantID string) (domain.Account, bool, error) {
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		return domain.Account{}, false, ErrInvalidProfile
	}
	if existing, err := u.repository.GetStaffAccountByTenant(tenantID); err == nil {
		account, err := u.repository.UpdateProfile(domain.UpdateProfileInput{
			UserID: existing.ID,
			Name:   strings.TrimSpace(input.Name),
			Phone:  strings.TrimSpace(input.Phone),
		})
		return account, false, err
	}
	input.Role = "owner"
	input.TenantID = tenantID
	session, err := u.registerWithRole(input, "owner")
	if err != nil {
		return domain.Account{}, false, err
	}
	return session.User, true, nil
}

func (u *UseCase) UpdateAccountRole(actorTenantID string, userID string, role string) (domain.Account, error) {
	actorTenantID = normalizeTenantScope(actorTenantID)
	userID = strings.TrimSpace(userID)
	role = strings.TrimSpace(role)
	if userID == "" {
		return domain.Account{}, ErrAccountNotFound
	}
	if !isStaffRole(role) {
		return domain.Account{}, ErrInvalidRole
	}

	current, err := u.repository.GetAccountByID(userID)
	if err != nil {
		return domain.Account{}, ErrAccountNotFound
	}
	if !canManageTenantAccount(actorTenantID, current.TenantID) {
		return domain.Account{}, ErrAccountNotFound
	}
	if isPrivilegedAdminRole(current.Role) && !isPrivilegedAdminRole(role) {
		if count, err := u.countAdminAccounts(current.TenantID); err != nil {
			return domain.Account{}, err
		} else if count <= 1 {
			return domain.Account{}, ErrLastAdmin
		}
	}

	account, err := u.repository.UpdateRole(userID, role)
	if err != nil {
		return domain.Account{}, ErrAccountNotFound
	}
	return account, nil
}

func (u *UseCase) DeleteStaffAccount(actorTenantID string, userID string) error {
	actorTenantID = normalizeTenantScope(actorTenantID)
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return ErrAccountNotFound
	}

	account, err := u.repository.GetAccountByID(userID)
	if err != nil {
		return ErrAccountNotFound
	}
	if !canManageTenantAccount(actorTenantID, account.TenantID) {
		return ErrAccountNotFound
	}
	if account.Role != "receptionist" && account.Role != "professional" {
		return ErrInvalidRole
	}
	if err := u.repository.DeleteAccount(userID); err != nil {
		return ErrAccountNotFound
	}
	return nil
}

func (u *UseCase) DeleteTenantOwner(userID string) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return ErrAccountNotFound
	}

	account, err := u.repository.GetAccountByID(userID)
	if err != nil {
		return ErrAccountNotFound
	}
	if account.Role != "owner" {
		return ErrInvalidRole
	}
	if err := u.repository.DeleteAccount(userID); err != nil {
		return ErrAccountNotFound
	}
	return nil
}

func (u *UseCase) ChangePassword(input domain.ChangePasswordInput) error {
	input.UserID = strings.TrimSpace(input.UserID)
	input.CurrentPassword = strings.TrimSpace(input.CurrentPassword)
	input.NewPassword = strings.TrimSpace(input.NewPassword)
	validation := NewValidationError()

	if input.CurrentPassword == "" {
		validation.Add("currentPassword", "Informe a senha atual.")
	}
	if input.NewPassword == "" {
		validation.Add("newPassword", "Informe a nova senha.")
	} else {
		validatePasswordField(input.NewPassword, "newPassword", validation)
	}
	if validation.HasErrors() {
		return validation
	}

	account, err := u.repository.GetAccountByID(input.UserID)
	if err != nil {
		return ErrAccountNotFound
	}
	if err := u.hasher.Compare(account.PasswordHash, input.CurrentPassword); err != nil {
		validation.Add("currentPassword", "Senha atual incorreta.")
		return validation
	}

	passwordHash, err := u.hasher.Hash(input.NewPassword)
	if err != nil {
		return err
	}
	return u.repository.UpdatePassword(input.UserID, passwordHash)
}

func (u *UseCase) RequestPasswordReset(input domain.PasswordResetRequestInput) (string, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if email == "" {
		return "", nil
	}
	account, err := u.repository.GetStaffAccountByEmail(email)
	if err != nil {
		return "", nil
	}
	return u.CreatePasswordResetToken(account.ID, 24*time.Hour)
}

func (u *UseCase) RequestPasswordResetCode(input domain.PasswordResetRequestInput) (string, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if email == "" {
		return "", nil
	}
	account, err := u.repository.GetStaffAccountByEmail(email)
	if err != nil {
		return "", nil
	}

	number, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	code := fmt.Sprintf("%06d", number.Int64())
	token := passwordResetCodeToken(email, code)
	if err := u.repository.CreatePasswordResetToken(domain.PasswordResetToken{
		TokenHash:    hashResetToken(token),
		ChallengeKey: passwordResetCodeChallengeKey(email),
		UserID:       account.ID,
		ExpiresAt:    u.now().Add(passwordResetCodeWindow).Format(time.RFC3339),
		MaxAttempts:  passwordResetCodeMaxAttempts,
	}); err != nil {
		return "", err
	}
	return code, nil
}

func (u *UseCase) CreatePasswordResetToken(userID string, ttl time.Duration) (string, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return "", ErrAccountNotFound
	}
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	token := base64.RawURLEncoding.EncodeToString(tokenBytes)
	if err := u.repository.CreatePasswordResetToken(domain.PasswordResetToken{
		TokenHash: hashResetToken(token),
		UserID:    userID,
		ExpiresAt: time.Now().Add(ttl).Format(time.RFC3339),
	}); err != nil {
		return "", err
	}
	return token, nil
}

func (u *UseCase) ConfirmPasswordReset(input domain.PasswordResetConfirmInput) error {
	token := strings.TrimSpace(input.Token)
	newPassword := strings.TrimSpace(input.NewPassword)
	validation := NewValidationError()
	if token == "" {
		validation.Add("token", "Link de redefinição inválido.")
	}
	if newPassword == "" {
		validation.Add("newPassword", "Informe a nova senha.")
	} else {
		validatePasswordField(newPassword, "newPassword", validation)
	}
	if validation.HasErrors() {
		return validation
	}

	passwordHash, err := u.hasher.Hash(newPassword)
	if err != nil {
		return err
	}
	if err := u.repository.ConsumePasswordResetToken(hashResetToken(token), passwordHash); err != nil {
		if errors.Is(err, domain.ErrPasswordResetTokenInvalid) {
			return ErrInvalidResetToken
		}
		return err
	}
	return nil
}

func (u *UseCase) ConfirmPasswordResetCode(input domain.PasswordResetCodeConfirmInput) error {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	code := strings.TrimSpace(input.Code)
	validation := NewValidationError()
	if email == "" {
		validation.Add("email", "Informe seu e-mail.")
	}
	if len(code) != 6 {
		validation.Add("code", "Informe o código de 6 dígitos.")
	} else if _, ok := new(big.Int).SetString(code, 10); !ok {
		validation.Add("code", "Informe apenas os números do código.")
	}
	if validation.HasErrors() {
		return validation
	}
	newPassword := strings.TrimSpace(input.NewPassword)
	passwordValidation := NewValidationError()
	if newPassword == "" {
		passwordValidation.Add("newPassword", "Informe a nova senha.")
	} else {
		validatePasswordField(newPassword, "newPassword", passwordValidation)
	}
	if passwordValidation.HasErrors() {
		return passwordValidation
	}
	passwordHash, err := u.hasher.Hash(newPassword)
	if err != nil {
		return err
	}
	err = u.repository.ConsumePasswordResetCode(
		passwordResetCodeChallengeKey(email),
		hashResetToken(passwordResetCodeToken(email, code)),
		passwordHash,
	)
	if errors.Is(err, domain.ErrPasswordResetTokenInvalid) || errors.Is(err, domain.ErrPasswordResetCodeLocked) {
		return ErrInvalidResetToken
	}
	return err
}

func passwordResetCodeToken(email string, code string) string {
	return "password-reset-code:staff:" + email + ":" + code
}

func passwordResetCodeChallengeKey(email string) string {
	return hashResetToken("password-reset-code:staff:" + email)
}

func (u *UseCase) AllowPasswordResetCodeRequest(email string, clientIP string) (bool, error) {
	checks := []struct {
		key   string
		limit int
	}{
		{key: "password-reset-code-request:ip:" + strings.TrimSpace(clientIP), limit: passwordResetCodeRequestIPLimit},
		{key: "password-reset-code-request:email:" + strings.ToLower(strings.TrimSpace(email)), limit: passwordResetCodeRequestEmailLimit},
	}
	for _, check := range checks {
		allowed, err := u.repository.AllowRateLimit(hashResetToken(check.key), check.limit, passwordResetCodeWindow)
		if err != nil || !allowed {
			return allowed, err
		}
	}
	return true, nil
}

func (u *UseCase) AllowPasswordResetCodeAttempt(clientIP string) (bool, error) {
	return u.repository.AllowRateLimit(
		hashResetToken("password-reset-code-confirm:ip:"+strings.TrimSpace(clientIP)),
		passwordResetCodeConfirmIPLimit,
		passwordResetCodeWindow,
	)
}

func (u *UseCase) AllowSecurityRequest(scope string, key string, limit int, window time.Duration) (bool, error) {
	if limit <= 0 || window <= 0 {
		return false, errors.New("invalid security rate limit")
	}
	normalized := strings.ToLower(strings.TrimSpace(scope)) + ":" + strings.ToLower(strings.TrimSpace(key))
	return u.repository.AllowRateLimit(hashResetToken(normalized), limit, window)
}

func (u *UseCase) CleanupExpiredAuthData(retention time.Duration) (domain.CleanupSummary, error) {
	if retention <= 0 {
		retention = 7 * 24 * time.Hour
	}
	return u.repository.CleanupExpiredAuthData(retention)
}

func (u *UseCase) Login(input domain.LoginInput) (domain.Session, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	password := strings.TrimSpace(input.Password)
	validation := NewValidationError()
	if email == "" {
		validation.Add("email", "Informe seu e-mail.")
	}
	if password == "" {
		validation.Add("password", "Informe sua senha.")
	}
	if validation.HasErrors() {
		return domain.Session{}, validation
	}

	account, err := u.repository.GetStaffAccountByEmail(email)
	if err != nil {
		validation.Add("email", "E-mail ou senha incorretos.")
		validation.Add("password", "E-mail ou senha incorretos.")
		return domain.Session{}, validation
	}

	if err := u.hasher.Compare(account.PasswordHash, password); err != nil {
		validation.Add("email", "E-mail ou senha incorretos.")
		validation.Add("password", "E-mail ou senha incorretos.")
		return domain.Session{}, validation
	}
	if isMFAEligibleRole(account.Role) && account.MFAEnabled {
		if strings.TrimSpace(input.MFACode) == "" {
			return domain.Session{}, ErrMFARequired
		}
		if u.mfaCipher == nil {
			return domain.Session{}, errors.New("mfa unavailable")
		}
		secret, decryptErr := u.mfaCipher.Decrypt(account.MFASecret)
		if decryptErr != nil || !validateTOTP(secret, input.MFACode, u.now()) {
			validation.Add("mfaCode", "Código de verificação inválido.")
			return domain.Session{}, validation
		}
	}

	return u.buildSession(account)
}

func (u *UseCase) MFAStatus(userID string) (domain.MFAStatus, error) {
	account, err := u.repository.GetAccountByID(userID)
	if err != nil || !isMFAEligibleRole(account.Role) {
		return domain.MFAStatus{}, ErrInvalidRole
	}
	return domain.MFAStatus{Enabled: account.MFAEnabled}, nil
}

func (u *UseCase) BeginMFASetup(userID string) (string, string, error) {
	account, err := u.repository.GetAccountByID(userID)
	if err != nil || !isMFAEligibleRole(account.Role) {
		return "", "", ErrInvalidRole
	}
	secret, err := generateTOTPSecret()
	if err != nil {
		return "", "", err
	}
	return secret, totpURI(secret, account.Email), nil
}

func (u *UseCase) EnableMFA(userID string, secret string, code string) error {
	account, err := u.repository.GetAccountByID(userID)
	if err != nil || !isMFAEligibleRole(account.Role) || u.mfaCipher == nil {
		return ErrInvalidRole
	}
	secret = strings.ToUpper(strings.TrimSpace(secret))
	if !validateTOTP(secret, code, u.now()) {
		return ErrInvalidCredentials
	}
	encrypted, err := u.mfaCipher.Encrypt(secret)
	if err != nil {
		return err
	}
	return u.repository.UpdateMFA(userID, encrypted, true)
}

func (u *UseCase) DisableMFA(userID string, password string, code string) error {
	account, err := u.repository.GetAccountByID(userID)
	if err != nil || !account.MFAEnabled || u.mfaCipher == nil {
		return ErrInvalidCredentials
	}
	if u.hasher.Compare(account.PasswordHash, strings.TrimSpace(password)) != nil {
		return ErrInvalidCredentials
	}
	secret, err := u.mfaCipher.Decrypt(account.MFASecret)
	if err != nil || !validateTOTP(secret, code, u.now()) {
		return ErrInvalidCredentials
	}
	return u.repository.UpdateMFA(userID, "", false)
}

func (u *UseCase) StepUpMFA(userID string, code string) (domain.Session, error) {
	account, err := u.repository.GetAccountByID(userID)
	if err != nil || !account.MFAEnabled || u.mfaCipher == nil {
		return domain.Session{}, ErrInvalidCredentials
	}
	secret, err := u.mfaCipher.Decrypt(account.MFASecret)
	if err != nil || !validateTOTP(secret, code, u.now()) {
		return domain.Session{}, ErrInvalidCredentials
	}
	return u.buildSessionWithMFA(account, true)
}

func isMFAEligibleRole(role string) bool {
	return role == "owner" || role == "admin"
}

func (u *UseCase) registerWithRole(input domain.RegisterInput, role string) (domain.Session, error) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(input.Email))
	password := strings.TrimSpace(input.Password)
	name := strings.TrimSpace(input.Name)
	validation := NewValidationError()

	if name == "" {
		validation.Add("name", "Informe seu nome completo.")
	}
	if normalizedEmail == "" {
		validation.Add("email", "Informe seu e-mail.")
	} else if !strings.Contains(normalizedEmail, "@") {
		validation.Add("email", "Informe um e-mail válido.")
	}
	if password == "" {
		validation.Add("password", "Informe sua senha.")
	} else {
		validatePassword(password, validation)
	}
	if validation.HasErrors() {
		return domain.Session{}, validation
	}

	passwordHash, err := u.hasher.Hash(password)
	if err != nil {
		return domain.Session{}, err
	}

	if _, err := u.repository.GetStaffAccountByEmail(normalizedEmail); err == nil {
		validation.Add("email", "Já existe uma conta cadastrada com este e-mail.")
		return domain.Session{}, validation
	}

	account, err := u.repository.CreateAccount(domain.Account{
		Name:         name,
		Email:        normalizedEmail,
		Phone:        strings.TrimSpace(input.Phone),
		PasswordHash: passwordHash,
		Role:         role,
		TenantID:     strings.TrimSpace(input.TenantID),
	})
	if err != nil {
		return domain.Session{}, err
	}

	return u.buildSession(account)
}

func validatePassword(password string, validation *ValidationError) {
	validatePasswordField(password, "password", validation)
}

func validatePasswordField(password string, field string, validation *ValidationError) {
	if len(password) < 8 {
		validation.Add(field, "A senha deve ter no mínimo 8 caracteres.")
		return
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !passwordSpecialCharPattern.MatchString(password) {
		validation.Add(field, "A senha deve conter letras maiúsculas, minúsculas, números e caractere especial.")
	}
}

func (u *UseCase) Me(userID string) (domain.Account, error) {
	return u.repository.GetAccountByID(userID)
}

func (u *UseCase) UpdateProfile(input domain.UpdateProfileInput) (domain.Account, error) {
	if strings.TrimSpace(input.UserID) == "" {
		return domain.Account{}, ErrInvalidProfile
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Phone = strings.TrimSpace(input.Phone)

	return u.repository.UpdateProfile(input)
}

func (u *UseCase) ParseToken(token string) (domain.Claims, error) {
	claims, err := u.tokens.Parse(token)
	if err != nil {
		return domain.Claims{}, err
	}
	if claims.Role == "onboarding" {
		return claims, nil
	}
	if strings.TrimSpace(claims.SessionID) == "" {
		if isProductionAuthEnvironment() {
			return domain.Claims{}, ErrInvalidCredentials
		}
	} else {
		session, sessionErr := u.repository.GetSession(claims.SessionID)
		if sessionErr != nil || session.UserID != claims.UserID {
			return domain.Claims{}, ErrInvalidCredentials
		}
		_ = u.repository.TouchSession(claims.SessionID, "", "", "")
	}
	state, err := u.repository.GetAuthState(claims.UserID)
	if err != nil || claims.AuthVersion != state.AuthVersion {
		return domain.Claims{}, ErrInvalidCredentials
	}
	claims.Role = state.Role
	claims.TenantID = state.TenantID
	return claims, nil
}

func isProductionAuthEnvironment() bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv("ERPAGENDAMENTOS_ENV")))
	return value == "production" || value == "prod"
}

func (u *UseCase) buildSession(account domain.Account) (domain.Session, error) {
	return u.buildSessionWithMFA(account, false)
}

func (u *UseCase) buildSessionWithMFA(account domain.Account, mfaVerified bool) (domain.Session, error) {
	mfaAt := int64(0)
	if mfaVerified {
		mfaAt = u.now().Unix()
	}
	sessionID := uuid.NewString()
	expiresAt := u.now().Add(7 * 24 * time.Hour)
	if err := u.repository.CreateSession(domain.DeviceSession{
		ID: sessionID, UserID: account.ID, ExpiresAt: expiresAt.Format(time.RFC3339),
	}); err != nil {
		return domain.Session{}, err
	}
	token, err := u.tokens.Sign(domain.Claims{
		UserID:      account.ID,
		Role:        account.Role,
		TenantID:    account.TenantID,
		AuthVersion: account.AuthVersion,
		MFAAt:       mfaAt,
		SessionID:   sessionID,
		ExpiresAt:   expiresAt.Unix(),
	})
	if err != nil {
		return domain.Session{}, err
	}

	return domain.Session{Token: token, User: account}, nil
}

func (u *UseCase) TrackSession(sessionID string, clientID string, deviceName string, ipAddress string) {
	if strings.TrimSpace(sessionID) == "" {
		return
	}
	_ = u.repository.TouchSession(sessionID, strings.TrimSpace(clientID), strings.TrimSpace(deviceName), strings.TrimSpace(ipAddress))
}

func (u *UseCase) ListSessions(userID string, currentSessionID string) ([]domain.DeviceSession, error) {
	sessions, err := u.repository.ListSessions(strings.TrimSpace(userID))
	if err != nil {
		return nil, err
	}
	for index := range sessions {
		sessions[index].Current = sessions[index].ID == currentSessionID
	}
	return sessions, nil
}

func (u *UseCase) RevokeSession(userID string, sessionID string) error {
	return u.repository.RevokeSession(strings.TrimSpace(userID), strings.TrimSpace(sessionID))
}

func (u *UseCase) RevokeOtherSessions(userID string, currentSessionID string) error {
	return u.repository.RevokeOtherSessions(strings.TrimSpace(userID), strings.TrimSpace(currentSessionID))
}

func (u *UseCase) CreateOnboardingToken(tenantID string) (string, error) {
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		return "", ErrInvalidProfile
	}
	return u.tokens.Sign(domain.Claims{
		UserID: "trial:" + tenantID, Role: "onboarding", TenantID: tenantID,
		ExpiresAt: u.now().Add(30 * time.Minute).Unix(),
	})
}

func (u *UseCase) countAdminAccounts(tenantID string) (int, error) {
	tenantID = normalizeTenantScope(tenantID)
	accounts, err := u.repository.ListAccounts()
	if err != nil {
		return 0, err
	}
	count := 0
	for _, account := range accounts {
		if isPrivilegedAdminRole(account.Role) && normalizeTenantScope(account.TenantID) == tenantID {
			count++
		}
	}
	return count, nil
}

func normalizeTenantScope(tenantID string) string {
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "default" {
		return ""
	}
	return tenantID
}

func canManageTenantAccount(actorTenantID string, targetTenantID string) bool {
	actorTenantID = normalizeTenantScope(actorTenantID)
	targetTenantID = normalizeTenantScope(targetTenantID)
	return actorTenantID == "" || actorTenantID == targetTenantID
}

func hashResetToken(token string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(token)))
	return hex.EncodeToString(sum[:])
}

func isStaffRole(role string) bool {
	switch role {
	case "owner", "admin", "receptionist", "professional":
		return true
	default:
		return false
	}
}

func isAdminRole(role string) bool {
	switch role {
	case "owner", "admin", "receptionist", "professional":
		return true
	default:
		return false
	}
}

func isPrivilegedAdminRole(role string) bool {
	switch role {
	case "admin", "owner":
		return true
	default:
		return false
	}
}

package memory

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/auth"
)

type AuthRepository struct {
	mu          sync.Mutex
	accounts    map[string]auth.Account
	byEmail     map[string]string
	resetTokens map[string]auth.PasswordResetToken
	rateLimits  map[string]rateLimitEntry
	sessions    map[string]auth.DeviceSession
}

type rateLimitEntry struct {
	WindowStartedAt time.Time
	Count           int
	ExpiresAt       time.Time
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{
		accounts:    make(map[string]auth.Account),
		byEmail:     make(map[string]string),
		resetTokens: make(map[string]auth.PasswordResetToken),
		rateLimits:  make(map[string]rateLimitEntry),
		sessions:    make(map[string]auth.DeviceSession),
	}
}

func (r *AuthRepository) CreateAccount(input auth.Account) (auth.Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	email := strings.ToLower(strings.TrimSpace(input.Email))
	if _, exists := r.byEmail[email]; exists {
		return auth.Account{}, fmt.Errorf("account already exists")
	}

	account := input
	account.ID = uuid.NewString()
	account.AuthVersion = 1
	account.Email = email
	account.CreatedAt = time.Now().Format(time.RFC3339)

	r.accounts[account.ID] = account
	r.byEmail[email] = account.ID
	return cloneAuthAccount(account), nil
}

func (r *AuthRepository) ListAccounts() ([]auth.Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	accounts := make([]auth.Account, 0, len(r.accounts))
	for _, account := range r.accounts {
		accounts = append(accounts, cloneAuthAccount(account))
	}
	return accounts, nil
}

func (r *AuthRepository) GetStaffAccountByEmail(email string) (auth.Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id, exists := r.byEmail[strings.ToLower(strings.TrimSpace(email))]
	if !exists {
		return auth.Account{}, fmt.Errorf("account not found")
	}
	return cloneAuthAccount(r.accounts[id]), nil
}

func (r *AuthRepository) GetAccountByID(userID string) (auth.Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	account, exists := r.accounts[userID]
	if !exists {
		return auth.Account{}, fmt.Errorf("account not found")
	}
	return cloneAuthAccount(account), nil
}

func (r *AuthRepository) GetAuthState(userID string) (auth.AuthState, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	account, exists := r.accounts[userID]
	if !exists {
		return auth.AuthState{}, fmt.Errorf("account not found")
	}
	return auth.AuthState{
		UserID:      account.ID,
		Role:        account.Role,
		TenantID:    account.TenantID,
		AuthVersion: account.AuthVersion,
	}, nil
}

func (r *AuthRepository) GetStaffAccountByTenant(tenantID string) (auth.Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, account := range r.accounts {
		if account.TenantID == tenantID && account.Role == "owner" {
			return cloneAuthAccount(account), nil
		}
	}
	return auth.Account{}, fmt.Errorf("account not found")
}

func (r *AuthRepository) UpdateProfile(input auth.UpdateProfileInput) (auth.Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	account, exists := r.accounts[input.UserID]
	if !exists {
		return auth.Account{}, fmt.Errorf("account not found")
	}
	if strings.TrimSpace(input.Name) != "" {
		account.Name = input.Name
	}
	if strings.TrimSpace(input.Phone) != "" {
		account.Phone = input.Phone
	}
	r.accounts[input.UserID] = account
	return cloneAuthAccount(account), nil
}

func (r *AuthRepository) UpdateRole(userID string, role string) (auth.Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	account, exists := r.accounts[userID]
	if !exists {
		return auth.Account{}, fmt.Errorf("account not found")
	}

	delete(r.byEmail, strings.ToLower(strings.TrimSpace(account.Email)))
	account.Role = strings.TrimSpace(role)
	account.AuthVersion++
	r.accounts[userID] = account
	r.byEmail[strings.ToLower(strings.TrimSpace(account.Email))] = account.ID
	return cloneAuthAccount(account), nil
}

func (r *AuthRepository) UpdatePassword(userID string, passwordHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	account, exists := r.accounts[userID]
	if !exists {
		return fmt.Errorf("account not found")
	}
	account.PasswordHash = passwordHash
	account.AuthVersion++
	r.accounts[userID] = account
	return nil
}

func (r *AuthRepository) UpdateMFA(userID string, encryptedSecret string, enabled bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	account, exists := r.accounts[userID]
	if !exists {
		return fmt.Errorf("account not found")
	}
	account.MFASecret = encryptedSecret
	account.MFAEnabled = enabled
	account.AuthVersion++
	r.accounts[userID] = account
	return nil
}

func (r *AuthRepository) CreatePasswordResetToken(input auth.PasswordResetToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for tokenHash, token := range r.resetTokens {
		if token.UserID == input.UserID && token.UsedAt == "" {
			token.UsedAt = time.Now().Format(time.RFC3339)
			r.resetTokens[tokenHash] = token
		}
	}
	r.resetTokens[input.TokenHash] = input
	return nil
}

func (r *AuthRepository) ConsumePasswordResetToken(tokenHash string, passwordHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	token, exists := r.resetTokens[tokenHash]
	if !exists || token.UsedAt != "" || resetTokenExpired(token) {
		return auth.ErrPasswordResetTokenInvalid
	}
	account, exists := r.accounts[token.UserID]
	if !exists {
		return fmt.Errorf("account not found")
	}
	account.PasswordHash = passwordHash
	account.AuthVersion++
	r.accounts[token.UserID] = account
	token.UsedAt = time.Now().Format(time.RFC3339)
	r.resetTokens[tokenHash] = token
	return nil
}

func (r *AuthRepository) ConsumePasswordResetCode(challengeKey string, tokenHash string, passwordHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	storedHash := ""
	var token auth.PasswordResetToken
	for candidateHash, candidate := range r.resetTokens {
		if candidate.ChallengeKey == challengeKey && candidate.UsedAt == "" {
			storedHash = candidateHash
			token = candidate
			break
		}
	}
	if storedHash == "" || resetTokenExpired(token) {
		return auth.ErrPasswordResetTokenInvalid
	}

	token.AttemptCount++
	if storedHash != tokenHash {
		if token.MaxAttempts > 0 && token.AttemptCount >= token.MaxAttempts {
			token.UsedAt = time.Now().Format(time.RFC3339)
			r.resetTokens[storedHash] = token
			return auth.ErrPasswordResetCodeLocked
		}
		r.resetTokens[storedHash] = token
		return auth.ErrPasswordResetTokenInvalid
	}

	account, exists := r.accounts[token.UserID]
	if !exists {
		return fmt.Errorf("account not found")
	}
	account.PasswordHash = passwordHash
	account.AuthVersion++
	r.accounts[token.UserID] = account
	token.UsedAt = time.Now().Format(time.RFC3339)
	r.resetTokens[storedHash] = token
	return nil
}

func (r *AuthRepository) AllowRateLimit(keyHash string, limit int, window time.Duration) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	entry := r.rateLimits[keyHash]
	if entry.WindowStartedAt.IsZero() || now.Sub(entry.WindowStartedAt) >= window {
		entry = rateLimitEntry{WindowStartedAt: now}
	}
	entry.Count++
	entry.ExpiresAt = now.Add(window)
	r.rateLimits[keyHash] = entry
	return entry.Count <= limit, nil
}

func (r *AuthRepository) CreateSession(input auth.DeviceSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().Format(time.RFC3339)
	input.CreatedAt = now
	input.LastSeenAt = now
	r.sessions[input.ID] = input
	return nil
}

func (r *AuthRepository) GetSession(sessionID string) (auth.DeviceSession, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	session, ok := r.sessions[sessionID]
	if !ok || sessionExpired(session) {
		return auth.DeviceSession{}, fmt.Errorf("session not found")
	}
	return session, nil
}

func (r *AuthRepository) ListSessions(userID string) ([]auth.DeviceSession, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	values := []auth.DeviceSession{}
	for _, session := range r.sessions {
		if session.UserID == userID && !sessionExpired(session) {
			values = append(values, session)
		}
	}
	return values, nil
}

func (r *AuthRepository) TouchSession(sessionID string, clientID string, deviceName string, ipAddress string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	session, ok := r.sessions[sessionID]
	if !ok {
		return fmt.Errorf("session not found")
	}
	if clientID != "" {
		session.ClientID = clientID
	}
	if deviceName != "" {
		session.DeviceName = deviceName
	}
	if ipAddress != "" {
		session.IPAddress = ipAddress
	}
	session.LastSeenAt = time.Now().Format(time.RFC3339)
	r.sessions[sessionID] = session
	return nil
}

func (r *AuthRepository) RevokeSession(userID string, sessionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if session, ok := r.sessions[sessionID]; ok && session.UserID == userID {
		delete(r.sessions, sessionID)
	}
	return nil
}

func (r *AuthRepository) RevokeOtherSessions(userID string, currentSessionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, session := range r.sessions {
		if session.UserID == userID && id != currentSessionID {
			delete(r.sessions, id)
		}
	}
	return nil
}

func (r *AuthRepository) CleanupExpiredAuthData(retention time.Duration) (auth.CleanupSummary, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	cutoff := time.Now().Add(-retention)
	summary := auth.CleanupSummary{}
	for tokenHash, token := range r.resetTokens {
		expiresAt, expiresErr := time.Parse(time.RFC3339, token.ExpiresAt)
		usedAt, usedErr := time.Parse(time.RFC3339, token.UsedAt)
		if (expiresErr == nil && expiresAt.Before(cutoff)) || (token.UsedAt != "" && usedErr == nil && usedAt.Before(cutoff)) {
			delete(r.resetTokens, tokenHash)
			summary.PasswordResetTokens++
		}
	}
	for key, entry := range r.rateLimits {
		if entry.ExpiresAt.Before(time.Now()) {
			delete(r.rateLimits, key)
			summary.RateLimits++
		}
	}
	for id, session := range r.sessions {
		if sessionExpired(session) {
			delete(r.sessions, id)
			summary.Sessions++
		}
	}
	return summary, nil
}

func (r *AuthRepository) DeleteAccount(userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	account, exists := r.accounts[userID]
	if !exists {
		return fmt.Errorf("account not found")
	}
	delete(r.accounts, userID)
	delete(r.byEmail, strings.ToLower(strings.TrimSpace(account.Email)))
	return nil
}

func sessionExpired(session auth.DeviceSession) bool {
	expiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
	return err != nil || time.Now().After(expiresAt)
}

func resetTokenExpired(token auth.PasswordResetToken) bool {
	expiresAt, err := time.Parse(time.RFC3339, token.ExpiresAt)
	return err != nil || time.Now().After(expiresAt)
}

func cloneAuthAccount(value auth.Account) auth.Account {
	return value
}

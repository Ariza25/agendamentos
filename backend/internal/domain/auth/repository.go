package auth

import (
	"errors"
	"time"
)

var (
	ErrPasswordResetTokenInvalid = errors.New("password reset token invalid")
	ErrPasswordResetCodeLocked   = errors.New("password reset code locked")
)

type Repository interface {
	CreateAccount(input Account) (Account, error)
	ListAccounts() ([]Account, error)
	GetStaffAccountByEmail(email string) (Account, error)
	GetAccountByID(userID string) (Account, error)
	GetAuthState(userID string) (AuthState, error)
	GetStaffAccountByTenant(tenantID string) (Account, error)
	UpdateProfile(input UpdateProfileInput) (Account, error)
	UpdateRole(userID string, role string) (Account, error)
	UpdatePassword(userID string, passwordHash string) error
	UpdateMFA(userID string, encryptedSecret string, enabled bool) error
	CreatePasswordResetToken(input PasswordResetToken) error
	ConsumePasswordResetToken(tokenHash string, passwordHash string) error
	ConsumePasswordResetCode(challengeKey string, tokenHash string, passwordHash string) error
	AllowRateLimit(keyHash string, limit int, window time.Duration) (bool, error)
	CreateSession(input DeviceSession) error
	GetSession(sessionID string) (DeviceSession, error)
	ListSessions(userID string) ([]DeviceSession, error)
	TouchSession(sessionID string, clientID string, deviceName string, ipAddress string) error
	RevokeSession(userID string, sessionID string) error
	RevokeOtherSessions(userID string, currentSessionID string) error
	CleanupExpiredAuthData(retention time.Duration) (CleanupSummary, error)
	DeleteAccount(userID string) error
}

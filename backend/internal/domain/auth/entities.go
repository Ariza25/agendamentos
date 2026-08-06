package auth

type Account struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
	TenantID     string `json:"tenantId"`
	AuthVersion  int    `json:"-"`
	MFAEnabled   bool   `json:"mfaEnabled"`
	MFASecret    string `json:"-"`
	CreatedAt    string `json:"createdAt"`
}

type Claims struct {
	UserID      string
	CustomerID  string
	Role        string
	TenantID    string
	AuthVersion int
	MFAAt       int64
	SessionID   string
	ExpiresAt   int64
}

type DeviceSession struct {
	ID         string `json:"id"`
	UserID     string `json:"-"`
	ClientID   string `json:"clientId"`
	DeviceName string `json:"deviceName"`
	IPAddress  string `json:"ipAddress"`
	LastSeenAt string `json:"lastSeenAt"`
	CreatedAt  string `json:"createdAt"`
	ExpiresAt  string `json:"expiresAt"`
	Current    bool   `json:"current"`
}

type AuthState struct {
	UserID      string
	Role        string
	TenantID    string
	AuthVersion int
}

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Role     string `json:"role"`
	TenantID string `json:"tenantId"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Audience string `json:"audience"`
	MFACode  string `json:"mfaCode"`
}

type MFAStatus struct {
	Enabled bool `json:"enabled"`
}

type UpdateProfileInput struct {
	UserID string `json:"-"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
}

type ChangePasswordInput struct {
	UserID          string `json:"-"`
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

type PasswordResetRequestInput struct {
	Email    string `json:"email"`
	Audience string `json:"audience"`
}

type PasswordResetConfirmInput struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

type PasswordResetCodeConfirmInput struct {
	Email       string `json:"email"`
	Code        string `json:"code"`
	NewPassword string `json:"newPassword"`
}

type PasswordResetToken struct {
	TokenHash    string
	ChallengeKey string
	UserID       string
	ExpiresAt    string
	UsedAt       string
	AttemptCount int
	MaxAttempts  int
}

type CleanupSummary struct {
	PasswordResetTokens int
	RateLimits          int
	Sessions            int
}

type Session struct {
	Token string  `json:"token"`
	User  Account `json:"user"`
}

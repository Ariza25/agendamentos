package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"erpagendamentos/backend/internal/domain/auth"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

const authAccountColumns = `id, name, email, phone, password_hash, role, tenant_id, auth_version, mfa_enabled, mfa_secret_encrypted, created_at`

func (r *AuthRepository) CreateAccount(input auth.Account) (auth.Account, error) {
	account := input
	account.ID = uuid.NewString()
	account.Email = strings.ToLower(strings.TrimSpace(account.Email))

	const query = `
		INSERT INTO catalog.users (id, name, email, phone, password_hash, role, tenant_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at
	`
	var createdAt time.Time
	if err := r.db.QueryRow(query, account.ID, account.Name, account.Email, account.Phone, account.PasswordHash, account.Role, account.TenantID).Scan(&createdAt); err != nil {
		return auth.Account{}, fmt.Errorf("insert user: %w", err)
	}
	account.CreatedAt = createdAt.Format(time.RFC3339)
	return r.GetAccountByID(account.ID)
}

func (r *AuthRepository) ListAccounts() ([]auth.Account, error) {
	rows, err := r.db.Query(`
		SELECT ` + authAccountColumns + `
		FROM catalog.users
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	var accounts []auth.Account
	for rows.Next() {
		var account auth.Account
		var createdAt time.Time
		if err := rows.Scan(
			&account.ID,
			&account.Name,
			&account.Email,
			&account.Phone,
			&account.PasswordHash,
			&account.Role,
			&account.TenantID,
			&account.AuthVersion,
			&account.MFAEnabled,
			&account.MFASecret,
			&createdAt,
		); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		account.CreatedAt = createdAt.Format(time.RFC3339)
		accounts = append(accounts, account)
	}

	return accounts, rows.Err()
}

func (r *AuthRepository) GetStaffAccountByEmail(email string) (auth.Account, error) {
	const query = `
		SELECT ` + authAccountColumns + `
		FROM catalog.users
		WHERE LOWER(email) = LOWER($1)
		  AND role IN ('owner', 'admin', 'receptionist', 'professional')
	`
	return r.scanAccount(query, email)
}

func (r *AuthRepository) GetAccountByID(userID string) (auth.Account, error) {
	const query = `
		SELECT ` + authAccountColumns + `
		FROM catalog.users
		WHERE id = $1
	`
	return r.scanAccount(query, userID)
}

func (r *AuthRepository) GetAuthState(userID string) (auth.AuthState, error) {
	var state auth.AuthState
	err := r.db.QueryRow(`
		SELECT u.id, u.role, u.tenant_id, u.auth_version
		FROM catalog.users u
		LEFT JOIN catalog.tenants t ON t.id = u.tenant_id
		WHERE u.id = $1
		  AND u.tenant_id <> ''
		  AND t.active = TRUE
		  AND t.deleted_at IS NULL
	`, userID).Scan(&state.UserID, &state.Role, &state.TenantID, &state.AuthVersion)
	if err != nil {
		if err == sql.ErrNoRows {
			return auth.AuthState{}, fmt.Errorf("account not found")
		}
		return auth.AuthState{}, fmt.Errorf("query auth state: %w", err)
	}
	return state, nil
}

func (r *AuthRepository) GetStaffAccountByTenant(tenantID string) (auth.Account, error) {
	const query = `
		SELECT ` + authAccountColumns + `
		FROM catalog.users
		WHERE tenant_id = $1 AND role = 'owner'
		ORDER BY created_at ASC
		LIMIT 1
	`
	return r.scanAccount(query, tenantID)
}

func (r *AuthRepository) UpdateProfile(input auth.UpdateProfileInput) (auth.Account, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.users
		SET name = COALESCE(NULLIF($2, ''), name),
		    phone = COALESCE(NULLIF($3, ''), phone),
		    updated_at = NOW()
		WHERE id = $1
	`, input.UserID, input.Name, input.Phone)
	if err != nil {
		return auth.Account{}, fmt.Errorf("update user: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return auth.Account{}, fmt.Errorf("check profile update: %w", err)
	}
	if affected == 0 {
		return auth.Account{}, fmt.Errorf("account not found")
	}
	return r.GetAccountByID(input.UserID)
}

func (r *AuthRepository) UpdateRole(userID string, role string) (auth.Account, error) {
	result, err := r.db.Exec(`
		UPDATE catalog.users
		SET role = $2, auth_version = auth_version + 1, updated_at = NOW()
		WHERE id = $1
	`, userID, strings.TrimSpace(role))
	if err != nil {
		return auth.Account{}, fmt.Errorf("update user role: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return auth.Account{}, fmt.Errorf("check role update: %w", err)
	}
	if affected == 0 {
		return auth.Account{}, fmt.Errorf("account not found")
	}
	return r.GetAccountByID(userID)
}

func (r *AuthRepository) UpdatePassword(userID string, passwordHash string) error {
	result, err := r.db.Exec(`
		UPDATE catalog.users
		SET password_hash = $2, auth_version = auth_version + 1, updated_at = NOW()
		WHERE id = $1
	`, userID, passwordHash)
	if err != nil {
		return fmt.Errorf("update user password: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check password update: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("account not found")
	}
	return nil
}

func (r *AuthRepository) UpdateMFA(userID string, encryptedSecret string, enabled bool) error {
	result, err := r.db.Exec(`
		UPDATE catalog.users
		SET mfa_secret_encrypted = $2, mfa_enabled = $3,
			auth_version = auth_version + 1, updated_at = NOW()
		WHERE id = $1
	`, userID, encryptedSecret, enabled)
	if err != nil {
		return fmt.Errorf("update user mfa: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		return fmt.Errorf("account not found")
	}
	return nil
}

func (r *AuthRepository) CreatePasswordResetToken(input auth.PasswordResetToken) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin password reset token: %w", err)
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`
		UPDATE catalog.password_reset_tokens
		SET used_at = NOW()
		WHERE user_id = $1 AND used_at IS NULL
	`, input.UserID); err != nil {
		return fmt.Errorf("invalidate password reset tokens: %w", err)
	}
	if _, err := tx.Exec(`
		INSERT INTO catalog.password_reset_tokens (
			token_hash, challenge_key, user_id, expires_at, attempt_count, max_attempts
		) VALUES ($1, $2, $3, $4::timestamptz, $5, $6)
		ON CONFLICT (token_hash) DO UPDATE SET
			challenge_key = EXCLUDED.challenge_key,
			user_id = EXCLUDED.user_id,
			expires_at = EXCLUDED.expires_at,
			used_at = NULL,
			attempt_count = EXCLUDED.attempt_count,
			max_attempts = EXCLUDED.max_attempts,
			created_at = NOW()
	`, input.TokenHash, input.ChallengeKey, input.UserID, input.ExpiresAt, input.AttemptCount, input.MaxAttempts); err != nil {
		return fmt.Errorf("insert password reset token: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit password reset token: %w", err)
	}
	return nil
}

func (r *AuthRepository) ConsumePasswordResetToken(tokenHash string, passwordHash string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin consume password reset token: %w", err)
	}
	defer tx.Rollback()
	var userID string
	var expiresAt time.Time
	var usedAt sql.NullTime
	if err := tx.QueryRow(`
		SELECT user_id, expires_at, used_at
		FROM catalog.password_reset_tokens
		WHERE token_hash = $1
		FOR UPDATE
	`, tokenHash).Scan(&userID, &expiresAt, &usedAt); err != nil {
		if err == sql.ErrNoRows {
			return auth.ErrPasswordResetTokenInvalid
		}
		return fmt.Errorf("lock password reset token: %w", err)
	}
	if usedAt.Valid || time.Now().After(expiresAt) {
		return auth.ErrPasswordResetTokenInvalid
	}
	if err := updateAuthPasswordAndVersion(tx, userID, passwordHash); err != nil {
		return err
	}
	if _, err := tx.Exec(`UPDATE catalog.password_reset_tokens SET used_at = NOW() WHERE token_hash = $1`, tokenHash); err != nil {
		return fmt.Errorf("consume password reset token: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit password reset token consumption: %w", err)
	}
	return nil
}

func (r *AuthRepository) ConsumePasswordResetCode(challengeKey string, tokenHash string, passwordHash string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin consume password reset code: %w", err)
	}
	defer tx.Rollback()
	var storedHash string
	var userID string
	var expiresAt time.Time
	var attemptCount int
	var maxAttempts int
	if err := tx.QueryRow(`
		SELECT token_hash, user_id, expires_at, attempt_count, max_attempts
		FROM catalog.password_reset_tokens
		WHERE challenge_key = $1 AND used_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
		FOR UPDATE
	`, challengeKey).Scan(&storedHash, &userID, &expiresAt, &attemptCount, &maxAttempts); err != nil {
		if err == sql.ErrNoRows {
			return auth.ErrPasswordResetTokenInvalid
		}
		return fmt.Errorf("lock password reset code: %w", err)
	}
	if time.Now().After(expiresAt) {
		_, _ = tx.Exec(`UPDATE catalog.password_reset_tokens SET used_at = NOW() WHERE token_hash = $1`, storedHash)
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit expired password reset code: %w", err)
		}
		return auth.ErrPasswordResetTokenInvalid
	}

	nextAttempt := attemptCount + 1
	if storedHash != tokenHash {
		locked := maxAttempts > 0 && nextAttempt >= maxAttempts
		if _, err := tx.Exec(`
			UPDATE catalog.password_reset_tokens
			SET attempt_count = $2,
			    used_at = CASE WHEN $3 THEN NOW() ELSE used_at END
			WHERE token_hash = $1
		`, storedHash, nextAttempt, locked); err != nil {
			return fmt.Errorf("record password reset attempt: %w", err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit password reset attempt: %w", err)
		}
		if locked {
			return auth.ErrPasswordResetCodeLocked
		}
		return auth.ErrPasswordResetTokenInvalid
	}

	if err := updateAuthPasswordAndVersion(tx, userID, passwordHash); err != nil {
		return err
	}
	if _, err := tx.Exec(`
		UPDATE catalog.password_reset_tokens
		SET attempt_count = $2, used_at = NOW()
		WHERE token_hash = $1
	`, storedHash, nextAttempt); err != nil {
		return fmt.Errorf("consume password reset code: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit password reset code consumption: %w", err)
	}
	return nil
}

func (r *AuthRepository) AllowRateLimit(keyHash string, limit int, window time.Duration) (bool, error) {
	var count int
	err := r.db.QueryRow(`
		INSERT INTO catalog.auth_rate_limits (key_hash, window_started_at, request_count, expires_at)
		VALUES ($1, NOW(), 1, NOW() + ($2 * INTERVAL '1 second'))
		ON CONFLICT (key_hash) DO UPDATE SET
			window_started_at = CASE
				WHEN catalog.auth_rate_limits.window_started_at <= NOW() - ($2 * INTERVAL '1 second') THEN NOW()
				ELSE catalog.auth_rate_limits.window_started_at
			END,
			request_count = CASE
				WHEN catalog.auth_rate_limits.window_started_at <= NOW() - ($2 * INTERVAL '1 second') THEN 1
				ELSE catalog.auth_rate_limits.request_count + 1
			END,
			expires_at = NOW() + ($2 * INTERVAL '1 second')
		RETURNING request_count
	`, keyHash, int(window.Seconds())).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("apply auth rate limit: %w", err)
	}
	return count <= limit, nil
}

func (r *AuthRepository) CreateSession(input auth.DeviceSession) error {
	_, err := r.db.Exec(`
		INSERT INTO catalog.auth_sessions (id, user_id, client_id, device_name, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6::timestamptz)
	`, input.ID, input.UserID, input.ClientID, input.DeviceName, input.IPAddress, input.ExpiresAt)
	if err != nil {
		return fmt.Errorf("create auth session: %w", err)
	}
	return nil
}

func (r *AuthRepository) GetSession(sessionID string) (auth.DeviceSession, error) {
	var session auth.DeviceSession
	var createdAt, lastSeenAt, expiresAt time.Time
	err := r.db.QueryRow(`
		SELECT id, user_id, client_id, device_name, ip_address, created_at, last_seen_at, expires_at
		FROM catalog.auth_sessions
		WHERE id = $1 AND revoked_at IS NULL AND expires_at > NOW()
	`, sessionID).Scan(&session.ID, &session.UserID, &session.ClientID, &session.DeviceName, &session.IPAddress, &createdAt, &lastSeenAt, &expiresAt)
	if err != nil {
		return auth.DeviceSession{}, fmt.Errorf("session not found")
	}
	session.CreatedAt = createdAt.Format(time.RFC3339)
	session.LastSeenAt = lastSeenAt.Format(time.RFC3339)
	session.ExpiresAt = expiresAt.Format(time.RFC3339)
	return session, nil
}

func (r *AuthRepository) ListSessions(userID string) ([]auth.DeviceSession, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, client_id, device_name, ip_address, created_at, last_seen_at, expires_at
		FROM catalog.auth_sessions
		WHERE user_id = $1 AND revoked_at IS NULL AND expires_at > NOW()
		ORDER BY last_seen_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list auth sessions: %w", err)
	}
	defer rows.Close()
	sessions := []auth.DeviceSession{}
	for rows.Next() {
		var session auth.DeviceSession
		var createdAt, lastSeenAt, expiresAt time.Time
		if err := rows.Scan(&session.ID, &session.UserID, &session.ClientID, &session.DeviceName, &session.IPAddress, &createdAt, &lastSeenAt, &expiresAt); err != nil {
			return nil, err
		}
		session.CreatedAt = createdAt.Format(time.RFC3339)
		session.LastSeenAt = lastSeenAt.Format(time.RFC3339)
		session.ExpiresAt = expiresAt.Format(time.RFC3339)
		sessions = append(sessions, session)
	}
	return sessions, rows.Err()
}

func (r *AuthRepository) TouchSession(sessionID string, clientID string, deviceName string, ipAddress string) error {
	_, err := r.db.Exec(`
		UPDATE catalog.auth_sessions
		SET client_id = COALESCE(NULLIF($2, ''), client_id),
			device_name = COALESCE(NULLIF($3, ''), device_name),
			ip_address = COALESCE(NULLIF($4, ''), ip_address), last_seen_at = NOW()
		WHERE id = $1 AND revoked_at IS NULL
		  AND (last_seen_at < NOW() - INTERVAL '5 minutes' OR client_id = '' OR device_name = '' OR ip_address = '')
	`, sessionID, clientID, deviceName, ipAddress)
	return err
}

func (r *AuthRepository) RevokeSession(userID string, sessionID string) error {
	_, err := r.db.Exec(`UPDATE catalog.auth_sessions SET revoked_at = NOW() WHERE id = $1 AND user_id = $2 AND revoked_at IS NULL`, sessionID, userID)
	return err
}

func (r *AuthRepository) RevokeOtherSessions(userID string, currentSessionID string) error {
	_, err := r.db.Exec(`UPDATE catalog.auth_sessions SET revoked_at = NOW() WHERE user_id = $1 AND id <> $2 AND revoked_at IS NULL`, userID, currentSessionID)
	return err
}

func (r *AuthRepository) CleanupExpiredAuthData(retention time.Duration) (auth.CleanupSummary, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return auth.CleanupSummary{}, fmt.Errorf("begin auth cleanup: %w", err)
	}
	defer tx.Rollback()
	seconds := int(retention.Seconds())
	tokensResult, err := tx.Exec(`
		DELETE FROM catalog.password_reset_tokens
		WHERE expires_at < NOW() - ($1 * INTERVAL '1 second')
		   OR (used_at IS NOT NULL AND used_at < NOW() - ($1 * INTERVAL '1 second'))
	`, seconds)
	if err != nil {
		return auth.CleanupSummary{}, fmt.Errorf("cleanup password reset tokens: %w", err)
	}
	rateResult, err := tx.Exec(`DELETE FROM catalog.auth_rate_limits WHERE expires_at < NOW()`)
	if err != nil {
		return auth.CleanupSummary{}, fmt.Errorf("cleanup auth rate limits: %w", err)
	}
	sessionResult, err := tx.Exec(`DELETE FROM catalog.auth_sessions WHERE expires_at < NOW() OR revoked_at < NOW() - ($1 * INTERVAL '1 second')`, seconds)
	if err != nil {
		return auth.CleanupSummary{}, fmt.Errorf("cleanup auth sessions: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return auth.CleanupSummary{}, fmt.Errorf("commit auth cleanup: %w", err)
	}
	tokens, _ := tokensResult.RowsAffected()
	rateLimits, _ := rateResult.RowsAffected()
	sessions, _ := sessionResult.RowsAffected()
	return auth.CleanupSummary{PasswordResetTokens: int(tokens), RateLimits: int(rateLimits), Sessions: int(sessions)}, nil
}

func (r *AuthRepository) DeleteAccount(userID string) error {
	result, err := r.db.Exec(`DELETE FROM catalog.users WHERE id = $1`, userID)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check user delete: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("account not found")
	}
	return nil
}

func (r *AuthRepository) scanAccount(query string, arg string) (auth.Account, error) {
	var account auth.Account
	var createdAt time.Time
	if err := r.db.QueryRow(query, arg).Scan(
		&account.ID,
		&account.Name,
		&account.Email,
		&account.Phone,
		&account.PasswordHash,
		&account.Role,
		&account.TenantID,
		&account.AuthVersion,
		&account.MFAEnabled,
		&account.MFASecret,
		&createdAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return auth.Account{}, fmt.Errorf("account not found")
		}
		return auth.Account{}, fmt.Errorf("query user: %w", err)
	}
	account.CreatedAt = createdAt.Format(time.RFC3339)
	if err := r.validateStaffTenant(account); err != nil {
		return auth.Account{}, err
	}
	return account, nil
}

func (r *AuthRepository) validateStaffTenant(account auth.Account) error {
	if strings.TrimSpace(account.TenantID) == "" {
		return fmt.Errorf("staff account without tenant")
	}
	var exists bool
	if err := r.db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM catalog.tenants
			WHERE id = $1
			  AND active = TRUE
			  AND deleted_at IS NULL
		)
	`, account.TenantID).Scan(&exists); err != nil {
		return fmt.Errorf("query account tenant: %w", err)
	}
	if !exists {
		return fmt.Errorf("staff account tenant not found")
	}
	return nil
}

func updateAuthPasswordAndVersion(tx *sql.Tx, userID string, passwordHash string) error {
	result, err := tx.Exec(`
		UPDATE catalog.users
		SET password_hash = $2, auth_version = auth_version + 1, updated_at = NOW()
		WHERE id = $1
	`, userID, passwordHash)
	if err != nil {
		return fmt.Errorf("update reset password: %w", err)
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return fmt.Errorf("account not found")
	}
	return nil
}

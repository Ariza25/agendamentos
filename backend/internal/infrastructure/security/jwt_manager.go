package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"erpagendamentos/backend/internal/domain/auth"
)

type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

type tokenPayload struct {
	Sub         string `json:"sub"`
	CustomerID  string `json:"customer_id"`
	Role        string `json:"role"`
	TenantID    string `json:"tenant_id"`
	AuthVersion int    `json:"auth_version"`
	MFAAt       int64  `json:"mfa_at,omitempty"`
	SessionID   string `json:"sid,omitempty"`
	Exp         int64  `json:"exp"`
	Iat         int64  `json:"iat"`
}

func NewJWTManager(secret string, ttl time.Duration) *JWTManager {
	return &JWTManager{secret: []byte(secret), ttl: ttl}
}

func (m *JWTManager) Sign(claims auth.Claims) (string, error) {
	headerJSON, err := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	if err != nil {
		return "", err
	}

	now := time.Now()
	expiresAt := now.Add(m.ttl).Unix()
	if claims.ExpiresAt > now.Unix() && claims.ExpiresAt < expiresAt {
		expiresAt = claims.ExpiresAt
	}
	payloadJSON, err := json.Marshal(tokenPayload{
		Sub:         claims.UserID,
		CustomerID:  claims.CustomerID,
		Role:        claims.Role,
		TenantID:    claims.TenantID,
		AuthVersion: claims.AuthVersion,
		MFAAt:       claims.MFAAt,
		SessionID:   claims.SessionID,
		Exp:         expiresAt,
		Iat:         now.Unix(),
	})
	if err != nil {
		return "", err
	}

	unsigned := encodeSegment(headerJSON) + "." + encodeSegment(payloadJSON)
	return unsigned + "." + m.sign(unsigned), nil
}

func (m *JWTManager) Parse(token string) (auth.Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return auth.Claims{}, fmt.Errorf("invalid token")
	}

	unsigned := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(parts[2]), []byte(m.sign(unsigned))) {
		return auth.Claims{}, fmt.Errorf("invalid token signature")
	}

	payloadBytes, err := decodeSegment(parts[1])
	if err != nil {
		return auth.Claims{}, fmt.Errorf("decode token payload: %w", err)
	}

	var payload tokenPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return auth.Claims{}, fmt.Errorf("invalid token payload: %w", err)
	}
	if payload.Exp < time.Now().Unix() {
		return auth.Claims{}, fmt.Errorf("token expired")
	}

	return auth.Claims{
		UserID:      payload.Sub,
		CustomerID:  payload.CustomerID,
		Role:        payload.Role,
		TenantID:    payload.TenantID,
		AuthVersion: payload.AuthVersion,
		MFAAt:       payload.MFAAt,
		SessionID:   payload.SessionID,
		ExpiresAt:   payload.Exp,
	}, nil
}

func (m *JWTManager) sign(value string) string {
	mac := hmac.New(sha256.New, m.secret)
	_, _ = mac.Write([]byte(value))
	return encodeSegment(mac.Sum(nil))
}

func encodeSegment(value []byte) string {
	return base64.RawURLEncoding.EncodeToString(value)
}

func decodeSegment(value string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(value)
}

package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func generateTOTPSecret() (string, error) {
	value := make([]byte, 20)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(value), nil
}

func validateTOTP(secret string, code string, now time.Time) bool {
	code = strings.TrimSpace(code)
	if len(code) != 6 {
		return false
	}
	for offset := int64(-1); offset <= 1; offset++ {
		if hmac.Equal([]byte(totpCode(secret, now.Unix()/30+offset)), []byte(code)) {
			return true
		}
	}
	return false
}

func totpCode(secret string, counter int64) string {
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(strings.TrimSpace(secret)))
	if err != nil {
		return ""
	}
	value := make([]byte, 8)
	binary.BigEndian.PutUint64(value, uint64(counter))
	mac := hmac.New(sha1.New, key)
	_, _ = mac.Write(value)
	sum := mac.Sum(nil)
	offset := sum[len(sum)-1] & 0x0f
	number := (binary.BigEndian.Uint32(sum[offset:offset+4]) & 0x7fffffff) % 1000000
	return fmt.Sprintf("%06d", number)
}

func totpURI(secret string, email string) string {
	label := "ERP Agendamentos:" + strings.TrimSpace(email)
	query := url.Values{"secret": {secret}, "issuer": {"ERP Agendamentos"}, "algorithm": {"SHA1"}, "digits": {strconv.Itoa(6)}, "period": {strconv.Itoa(30)}}
	return "otpauth://totp/" + url.PathEscape(label) + "?" + query.Encode()
}

package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

type SecretCipher struct {
	aead cipher.AEAD
}

func NewSecretCipher(keyMaterial string) (*SecretCipher, error) {
	key := sha256.Sum256([]byte(keyMaterial))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &SecretCipher{aead: aead}, nil
}

func (c *SecretCipher) Encrypt(value string) (string, error) {
	return c.EncryptWithContext(value, "")
}

func (c *SecretCipher) EncryptWithContext(value string, context string) (string, error) {
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := c.aead.Seal(nonce, nonce, []byte(value), []byte(context))
	return base64.RawURLEncoding.EncodeToString(sealed), nil
}

func (c *SecretCipher) Decrypt(value string) (string, error) {
	return c.DecryptWithContext(value, "")
}

func (c *SecretCipher) DecryptWithContext(value string, context string) (string, error) {
	encoded, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil || len(encoded) < c.aead.NonceSize() {
		return "", fmt.Errorf("invalid encrypted secret")
	}
	nonce := encoded[:c.aead.NonceSize()]
	plain, err := c.aead.Open(nil, nonce, encoded[c.aead.NonceSize():], []byte(context))
	if err != nil {
		return "", fmt.Errorf("decrypt secret: %w", err)
	}
	return string(plain), nil
}

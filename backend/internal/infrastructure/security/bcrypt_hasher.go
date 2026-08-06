package security

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost int) *BcryptHasher {
	return &BcryptHasher{cost: cost}
}

func (h *BcryptHasher) Hash(value string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(value), h.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (h *BcryptHasher) Compare(hash string, value string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(value))
}

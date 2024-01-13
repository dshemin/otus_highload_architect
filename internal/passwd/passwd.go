package passwd

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost int) *BcryptHasher {
	const defaultCost = 12

	if cost <= 0 {
		cost = defaultCost
	}

	return &BcryptHasher{
		cost: cost,
	}
}

func (h *BcryptHasher) Hash(p string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p), h.cost)
	if err != nil {
		return "", errors.Wrap(err, "bcrypt generate from password")
	}
	return string(hash), nil
}

var ErrPasswordMismatch = errors.New("password mismatch")

func (h *BcryptHasher) Check(given, hashed string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(given))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrPasswordMismatch
	}
	return errors.Wrap(err, "compare passwords")
}

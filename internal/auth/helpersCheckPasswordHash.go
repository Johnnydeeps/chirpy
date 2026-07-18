package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

func CheckPasswordHash(password, hash string) (bool, error) {
	isValid, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("password and hash comparison failed: %w", err)
	}
	return isValid, nil
}

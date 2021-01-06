package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// GetHash - todo
func GetHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", errors.New("Internal server error")
	}

	return string(hash), nil
}

// CheckHash - todo
func CheckHash(hash string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("Wrong username of password")
	}

	return nil
}

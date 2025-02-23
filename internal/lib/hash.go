package lib

import (
	"golang.org/x/crypto/bcrypt"
)

// Hash a password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 15)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// Compare a password and a hashed password and return a boolean.
func ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

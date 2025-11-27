package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}
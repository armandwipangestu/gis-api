package helpers

import (
	"armandwipangestu/gis-api/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Value secret get from environment variable `JWT_SECRET`
var jwtKey = []byte(config.GetEnv("JWT_SECRET", "secret_key"))

func GenerateToken(username string) string {
	// Configure expired time token, in this case we set to 60 minute from now
	expirationTime := time.Now().Add(60 * time.Minute)

	// Create claims JWT
	// Subject has username, and ExpiredAt that determine expired time token
	claims := &jwt.RegisteredClaims{
		Subject:	username,
		ExpiresAt:	jwt.NewNumericDate(expirationTime),
	}

	// Create new token with claim that already created
	// Using algorithm HS256 for signing the token
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)

	// Return token in form string
	return token
}
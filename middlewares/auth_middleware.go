package middlewares

import (
	"armandwipangestu/gis-api/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Get secret key from environment variable
// If not exist, use default "secret_key"
var jwtKey = []byte(config.GetEnv("JWT_SECRET", "secret_key"))

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get header Authorization from request
		tokenString := c.GetHeader("Authorization")

		// If token is empty, return with response 401 Unauthorized
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token is required",
			})
			c.Abort() // Stop incoming request
			return
		}

		// Remove prefix "Bearer " from token
		// Common Header value is: "Bearer <token>"
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Create struct to accomodate token claim
		claims := &jwt.RegisteredClaims{}

		// Parse token and verify signature with jwtKey
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Return secret key to verify token
			return jwtKey, nil
		})

		// If token is not valid or has an error when parsing
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthenticated",
			})
			c.Abort() // Stop incoming request
			return
		}

		// Store claim "sub" (username) to the context
		c.Set("username", claims.Subject)

		// Next to incoming handler
		c.Next()
	}
}
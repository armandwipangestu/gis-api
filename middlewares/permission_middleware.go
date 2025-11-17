package middlewares

import (
	"armandwipangestu/gis-api/database"
	"armandwipangestu/gis-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Middleware to check is user has some permission
func Permission(permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get username from context (stored from middleware Auth)
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		var user models.User
		// Load full user data with relationship roles and permissions
		err := database.DB.
				Preload("Roles.Permissions").
				Where("username = ?", username).
				First(&user).Error

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Check if user has which permission requested
		for _, role := range user.Roles {
			for _, perm := range role.Permissions {
				if perm.Name == permissionName {
					c.Next() // User has access, go to next request
					return
				}
			}
		}

		// If not found, deny access
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden - permission denied"})
		c.Abort()
	}
}
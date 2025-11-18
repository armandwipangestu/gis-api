package routes

import (
	authController "armandwipangestu/gis-api/controllers/auth"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// Initialize gin
	router := gin.Default()

	// Auth routes (no auth required)
	auth := router.Group("/api")
	{
		auth.POST("/login", authController.Login)
	}

	return router
}
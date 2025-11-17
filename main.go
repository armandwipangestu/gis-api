package main

import (
	"armandwipangestu/gis-api/config"
	"armandwipangestu/gis-api/database"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	database.InitDB()

	// Initialize Gin
	router := gin.Default()

	// Create route with method GET on root or `/` endpoint
	router.GET("/", func(c *gin.Context) {
		// Return response JSON
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	// Serve app with port 3000
	router.Run(":" + config.GetEnv("APP_PORT", "3000"))
}
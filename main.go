package main

import (
	"armandwipangestu/gis-api/config"
	"armandwipangestu/gis-api/database"
	"armandwipangestu/gis-api/database/seeders"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Initialize database (connection, migration stuff)
	database.InitDB()

	// Running the seeders
	seeders.Seed()

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
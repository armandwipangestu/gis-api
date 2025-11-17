package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
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
	router.Run(":3000")
}
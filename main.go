package main

import (
	"armandwipangestu/gis-api/config"
	"armandwipangestu/gis-api/database"
	"armandwipangestu/gis-api/database/seeders"
	"armandwipangestu/gis-api/routes"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Initialize database (connection, migration stuff)
	database.InitDB()

	// Running the seeders
	seeders.Seed()

	// Setup router
	r := routes.SetupRouter()

	// Serve app
	r.Run(":" + config.GetEnv("APP_PORT", "3000"))
}
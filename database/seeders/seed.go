package seeders

import (
	"armandwipangestu/gis-api/database"
	"log"
)

func Seed() {
	db := database.DB
	log.Println("Running database seeders...")

	// Running seeder sequential
	SeedPermissions(db)
	SeedRoles(db)
	SeedUsers(db)
	SeedSetting(db)

	log.Println("Database seeding completed!")
}
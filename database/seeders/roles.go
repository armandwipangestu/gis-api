package seeders

import (
	"armandwipangestu/gis-api/models"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	roles := []models.Role{
		{Name: "admin"},
		{Name: "user"},
	}

	// Loop and assign permissions based on role
	for _, role := range roles {
		// Cek/insert role
		db.FirstOrCreate(&role, models.Role{Name: role.Name})

		// Get all permission on database
		var allPermissions []models.Permission
		db.Find(&allPermissions)

		switch role.Name {

		case "admin":
			// Admin: assign all permission
			db.Model(&role).Association("Permissions").Replace(allPermissions)

		case "user":
			// User: assign some permission
			var viewOnly []models.Permission
			db.Where("name IN ?", []string{
				"categories-index",
				"categories-create",
				"categories-show",
				"categories-edit",
				"maps-index",
				"maps-create",
				"maps-show",
				"maps-edit",
			}).Find(&viewOnly)
			db.Model(&role).Association("Permissions").Replace(viewOnly)
			
		}
	}
}
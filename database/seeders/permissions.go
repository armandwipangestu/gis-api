package seeders

import (
	"armandwipangestu/gis-api/models"

	"gorm.io/gorm"
)

func SeedPermissions(db *gorm.DB) {
	permissions := []models.Permission{
		{Name: "dashboard-index"},

		{Name: "users-index"},
		{Name: "users-create"},
		{Name: "users-show"},
		{Name: "users-edit"},
		{Name: "users-update"},
		{Name: "users-delete"},

		{Name: "permissions-index"},
		{Name: "permissions-create"},
		{Name: "permissions-show"},
		{Name: "permissions-edit"},
		{Name: "permissions-update"},
		{Name: "permissions-delete"},

		{Name: "roles-index"},
		{Name: "roles-create"},
		{Name: "roles-show"},
		{Name: "roles-edit"},
		{Name: "roles-update"},
		{Name: "roles-delete"},

		{Name: "categories-index"},
		{Name: "categories-create"},
		{Name: "categories-show"},
		{Name: "categories-edit"},
		{Name: "categories-update"},
		{Name: "categories-delete"},

		{Name: "maps-index"},
		{Name: "maps-create"},
		{Name: "maps-show"},
		{Name: "maps-edit"},
		{Name: "maps-update"},
		{Name: "maps-delete"},

		{Name: "settings-show"},
		{Name: "settings-update"},
	}

	for _, p := range permissions {
		db.FirstOrCreate(&p, models.Permission{Name: p.Name})
	}
}
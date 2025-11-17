package seeders

import (
	"armandwipangestu/gis-api/models"
	"log"

	"gorm.io/gorm"
)

func SeedSetting(db *gorm.DB) {
	// Default values
	defaults := models.Setting{
		Title:				"GIS Desa",
		Description:		"Eksplorasi Desa secara interaktif melalui peta GIS.",
		MapCenterLat: 		"7.592589928951457",
		MapCenterLng: 		"112.26113954274147",
		MapZoom: 			16,
		VillageBoundary: 	"[]",
	}

	var out models.Setting
	if err := db.
			Where(&models.Setting{Id: 1}).
			Attrs(defaults).
			FirstOrCreate(&out).Error; err != nil {
				log.Fatalf("Failed seeding setting: %v", err)
			}
}
package models

import (
	"time"
)

type Setting struct {
	Id					uint		`json:"id" gorm:"primaryKey"`
	Title				string		`json:"title"`
	Description			string		`json:"description" gorm:"type:text"`
	MapCenterLat		string		`json:"map_center_lat" gorm:"type:double"`
	MapCenterLng		string		`json:"map_center_lang" gorm:"type:double"`
	MapZoom				int			`json:"map_zoom"`
	VillageBoundary		string		`json:"village_boundary" gorm:"type:json"`
	CreatedAt			time.Time	`json:"created_at"`
	UpdatedAt			time.Time	`json:"updated_at"`
}
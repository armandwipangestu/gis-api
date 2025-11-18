package models

import (
	"time"
)

type Map struct {
	Id				uint			`json:"id" gorm:"primaryKey"`
	Image			string			`json:"image"`
	Name			string			`json:"name"`
	Slug			string			`json:"slug" gorm:"unique"`
	Description		string			`json:"description" gorm:"type:text"`
	Address			string			`json:"address" gorm:"type:text"`
	Latitude		string			`json:"latitude" gorm:"type:double"`
	Longitude		string			`json:"longitude" gorm:"type:double"`
	Geometry		string			`json:"geometry" gorm:"type:json"`
	Status			string			`json:"status" gorm:"type:enum('active', 'inactive');default:'active'"`
	CategoryID		uint			`json:"category_id"`
	Category		Category		`json:"category" gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
	CreateadAt		time.Time		`json:"created_at"`
	UpdatedAt		time.Time		`json:"updated_at"`
}
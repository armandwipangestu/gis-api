package models

import (
	"time"
)

type Category struct {
	Id				uint		`json:"id" gorm:"primaryKey"`
	Image			string		`json:"image"`
	Name			string		`json:"name"`
	Slug			string		`json:"slug" gorm:"unique"`
	Color			string		`json:"color"`
	Description		string		`json:"description" gorm:"type:text"`
	Maps			[]Map		`json:"maps" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
}
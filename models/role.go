package models

import (
	"time"
)

type Role struct {
	Id				uint			`json:"id" gorm:"primaryKey"`
	Name			string			`json:"name"`
	Permissions		[]Permission	`json:"permissions" gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE"`
	CreatedAt		time.Time		`json:"created_at"`
	UpdatedAt		time.Time		`json:"updated_at"`
}
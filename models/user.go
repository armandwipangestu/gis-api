package models

import (
	"time"
)

type User struct {
	Id				uint			`json:"id" gorm:"primaryKey"`
	Name			string			`json:"name"`
	Username		string			`json:"username" gorm:"unique;not null"`
	Email			string			`json:"email" gorm:"unique;not null"`
	Password		string			`json:"-"` // not showing on JSON response
	Roles			[]Role			`json:"roles" gorm:"many2many:user_roles"`
	CreatedAt		time.Time		`json:"created_at"`
	UpdatedAt		time.Time		`json:"updated_at"`
}
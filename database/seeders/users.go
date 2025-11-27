package seeders

import (
	"armandwipangestu/gis-api/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	// Hash default password
	passowrd, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	// Getting role
	var adminRole, userRole models.Role
	db.Where("name = ?", "admin").First(&adminRole)
	db.Where("name = ?", "user").First(&userRole)

	// Initialize default user
	users := []models.User{
		{
			Name:		"Admin",
			Username:	"admin",
			Email:		"admin@gis-api.internal",
			Password:	string(passowrd),
			Roles:		[]models.Role{adminRole},
		},
		{
			Name:		"User",
			Username:	"user",
			Email:		"user@gis-api.internal",
			Password: 	string(passowrd),
			Roles:		[]models.Role{userRole},
		},
	}

	for _, u := range users {
		var user models.User
		err := db.Where("username = ?", u.Username).First(&user).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// Data not exist, create one
				db.Create(&u)
			} else {
				// Unexpected error, like connection or query problem
				panic(err)
			}
		} else {
			// If user exist, update data and roles
			db.Model(&user).Updates(models.User{
				Email:		u.Email,
				Password:	u.Password,
			})
			db.Model(&user).Association("Roles").Replace(u.Roles)
		}
	}
}
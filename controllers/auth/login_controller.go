package auth

import (
	"armandwipangestu/gis-api/database"
	"armandwipangestu/gis-api/helpers"
	"armandwipangestu/gis-api/models"
	"armandwipangestu/gis-api/structs"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	// Initialize struct for accomodate data login receive from request
	var req = structs.UserLoginRequest{}
	var user = models.User{}

	// Validation input from request body using ShouldBindJSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors: helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Get User with preload relationship Role and Permission
	if err := database.DB.Preload("Roles").Preload("Roles.Permissions").Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, structs.ErrorResponse{
			Success: false,
			Message: "User Not Found",
			Errors: helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Compare password from request with password in database, which has been hash using bcrpyt
	// If not equal, send response error Unauthorized
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, structs.ErrorResponse{
			Success: false,
			Message: "Invalid Password",
			Errors: helpers.TranslateErrorMessage(err),
		})
		return
	}

	// If login success, generate new token for the user
	token := helpers.GenerateToken(user.Username)

	// Mapping permissions to map[string]bool
	permissionMap := helpers.GetPermissionMap(user.Roles)

	// Send response success with status OK and data user with token
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Login Success",
		Data: structs.UserResponse{
			Id:				user.Id,
			Name:			user.Name,
			Username:		user.Username,
			Email:			user.Email,
			Permissions:	permissionMap,
			Token:			&token,
		},
	})
}
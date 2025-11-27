package admin

import (
	"armandwipangestu/gis-api/database"
	"armandwipangestu/gis-api/helpers"
	"armandwipangestu/gis-api/models"
	"armandwipangestu/gis-api/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FindUsers(c *gin.Context) {
	// Initialize slice to accomodate data
	var users []models.User
	var usersResponse []structs.UserResponse

	// Initialize total
	var total int64

	// Get parameter search, page, limit, offset
	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseUrl := helpers.BuildBaseUrl(c)

	// Prepare the query
	query := database.DB.Preload("Roles").Model(&models.User{})
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Count total data and take data based on pagination
	query.Count(&total)
	err := query.Order("id desc").Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch users",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Mapping every user to UserResponse
	for _, user := range users {
		// Initialize empty slice, default is []
		roleResponses := []structs.RoleResponse{}

		// Mapping roles if exist
		for _, role := range user.Roles{
			roleResponses = append(roleResponses, structs.RoleResponse{
				Id:			role.Id,
				Name:		role.Name,
				CreatedAt:	role.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: 	role.UpdatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		// Append to list response
		usersResponse = append(usersResponse, structs.UserResponse{
			Id:			user.Id,
			Username:	user.Username,
			Name:		user.Name,
			Email:		user.Email,
			Roles:		roleResponses, // it will empty or `[]` if role not exist
			CreatedAt: 	user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: 	user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// Send response with pagination structure
	helpers.PaginateResponse(c, usersResponse, total, page, limit, baseUrl, search, "List Data Users")
}
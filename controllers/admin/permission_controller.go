package admin

import (
	"net/http"
	"armandwipangestu/gis-api/database"
	"armandwipangestu/gis-api/helpers"
	"armandwipangestu/gis-api/models"
	"armandwipangestu/gis-api/structs"

	"github.com/gin-gonic/gin"
)

// Get all permission with search and pagination
func FindPermissions(c *gin.Context) {
	var permissions []models.Permission
	var total int64

	// Get parameter search, page, limit, and offset from helper
	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseUrl := helpers.BuildBaseUrl(c)

	// Query first from table permissions
	query := database.DB.Model(&models.Permission{})
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Count total data
	query.Count(&total)

	// Get data based on limit and offset
	err := query.Order("id desc").Limit(limit).
		Offset(offset).Find(&permissions).Error; if err != nil {
			c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
				Success: false,
				Message: "Faield to fetch permissions",
				Errors: helpers.TranslateErrorMessage(err),
			})

			return
		}

	// Response JSON with pagination
	helpers.PaginateResponse(c, permissions, total, page, limit, baseUrl, search, "List Data Permissions")
}
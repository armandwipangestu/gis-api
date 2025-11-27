package admin

import (
	"armandwipangestu/gis-api/database"
	"armandwipangestu/gis-api/helpers"
	"armandwipangestu/gis-api/models"
	"armandwipangestu/gis-api/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Show all category with feature search and pagination
func FindCategories(c *gin.Context) {
	// Initialize slice for accomodate data
	var categories []models.Category
	var total int64

	// Get parameter search, page, limit, offset
	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseUrl := helpers.BuildBaseUrl(c)

	// Prepare query
	query := database.DB.Model(&models.Category{})
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Count total data and get data based on pagination
	query.Count(&total)
	err := query.Order("id desc").Limit(limit).Offset(offset).Find(&categories).Error; if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success:	false,
			Message:	"Failed to fetch categories",
			Errors:		helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Send response with pagination structure
	helpers.PaginateResponse(c, categories, total, page, limit, baseUrl, search, "List Data Categories")
}
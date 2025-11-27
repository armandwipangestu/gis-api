package admin

import (
	"armandwipangestu/gis-api/database"
	"armandwipangestu/gis-api/helpers"
	"armandwipangestu/gis-api/models"
	"armandwipangestu/gis-api/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Show all map with feature search and pagination
func FindMaps(c *gin.Context) {
	// Initialize slice to accomodate data
	var maps []models.Map
	var total int64

	// Get parameter search, page, limit, offset
	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseUrl := helpers.BuildBaseUrl(c)

	// Prepare query
	query := database.DB.Preload("Category").Model(&models.Map{})
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Count total data and take data based on pagination
	query.Count(&total)
	err := query.Order("id desc").Limit(limit).Offset(offset).Find(&maps).Error; if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch maps",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Send response with structure pagination
	helpers.PaginateResponse(c, maps, total, page, limit, baseUrl, search, "List Data Maps")
}
package admin

import (
	"armandwipangestu/gis-api/database"
	"armandwipangestu/gis-api/helpers"
	"armandwipangestu/gis-api/models"
	"armandwipangestu/gis-api/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Dashboard function get statsitic for dashboard GIS
func Dashboard(c *gin.Context) {
	var (
		categoriesCount			int64
		mapsCount				int64
		activeMapsCount			int64
		inactiveMapsCount		int64
	)

	// Count total categories
	if err := database.DB.Model(&models.Category{}).
			Count(&categoriesCount).Error; err != nil {
				c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
					Success: false,
					Message: "Failed to get categories count",
					Errors: helpers.TranslateErrorMessage(err),
				})

				return
			}

	// Count total maps (all)
	if err := database.DB.Model(&models.Map{}).
			Count(&mapsCount).Error; err != nil {
				c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
					Success: false,
					Message: "Failed to get maps count",
					Errors: helpers.TranslateErrorMessage(err),
				})

				return
			}

	// Count total maps active
	if err := database.DB.Model(&models.Map{}).
			Where("status = ?", "active").
			Count(&activeMapsCount).Error; err != nil {
				c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
					Success: false,
					Message: "Failed to get active maps count",
					Errors: helpers.TranslateErrorMessage(err),
				})

				return
			}

	// Count total maps inactive
	if err := database.DB.Model(&models.Map{}).
			Where("status = ?", "inactive").
			Count(&inactiveMapsCount).Error; err != nil {
				c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
					Success: false,
					Message: "Failed to get inactive maps count",
					Errors:	 helpers.TranslateErrorMessage(err),
				})

				return
			}

	// Response OK
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Dashboard stats retrieved successfully",
		Data: structs.DashboardResponse{
			CategoriesCount: categoriesCount,
			MapsCount: mapsCount,
			ActiveMapsCount: activeMapsCount,
			InactiveMapsCount: inactiveMapsCount,
		},
	})
}
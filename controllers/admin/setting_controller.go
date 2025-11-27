package admin

import (
	"armandwipangestu/gis-api/database"
	"armandwipangestu/gis-api/helpers"
	"armandwipangestu/gis-api/models"
	"armandwipangestu/gis-api/structs"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Get one setting based on id
func GetSetting(c *gin.Context) {
	// Initialize variable
	var setting models.Setting

	if err := database.DB.First(&setting, 1).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to load settngis",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Setting found",
		Data:	 setting,
	})
}

// Update data setting
func UpdateSetting(c *gin.Context) {
	// Get data existing
	var setting models.Setting
	if err := database.DB.First(&setting, 1).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Setting not found",
			Errors:	 helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Bind request (JSON)
	var req structs.SettingUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Map field to model
	setting.Title 			= req.Title
	setting.Description		= req.Description
	setting.MapCenterLat	= req.MapCenterLat
	setting.MapCenterLng	= req.MapCenterLng
	setting.MapZoom			= req.MapZoom
	setting.VillageBoundary	= req.VillageBoundary

	// Store changed
	if err := database.DB.Save(&setting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update setting",
			Errors:	 helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Respon success
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Setting updated successfully",
		Data:	 setting,
	})
}
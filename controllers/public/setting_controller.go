package public

import (
	"armandwipangestu/gis-api/database"
	"armandwipangestu/gis-api/helpers"
	"armandwipangestu/gis-api/models"
	"armandwipangestu/gis-api/structs"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetSetting(c *gin.Context) {
	// Get settings singleton (assume id=1)
	var setting models.Setting

	if err := database.DB.First(&setting, 1).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to load settings",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Send success response
	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Settings retrieved successfully",
		Data:	 setting,
	})
}
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
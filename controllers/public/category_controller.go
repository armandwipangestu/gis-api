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

func PublicCategoriesWithMaps(c *gin.Context) {
	// Initialize variable
	var categories []models.Category

	// baseUrl := helpers.BuildBaseUrl(c)
	includeInactive := c.Query("include_inactive") == "1"

	// Preload relationship Maps. Default just map active.
	err := database.DB.
		Preload("Maps", func(db *gorm.DB) *gorm.DB {
			db = db.Order("id DESC")
			if !includeInactive {
				db = db.Where("status = ?", "active")
			}

			return db
		}).
		Order("name ASC").
		Find(&categories).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch public categories",
			Errors:	 helpers.TranslateErrorMessage(err),
		})

		return
	}

	out := make([]structs.CategoryResponse, 0, len(categories))
	for _, cat := range categories {
		item := structs.CategoryResponse{
			Id:					cat.Id,
			Name:				cat.Name,
			Slug:				cat.Slug,
			Description:		cat.Description,
			Image:				cat.Image,
			Maps:				make([]structs.PublicMap, 0, len(cat.Maps)),
			CreatedAt:			cat.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:			cat.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		for _, m := range cat.Maps {
			item.Maps = append(item.Maps, structs.PublicMap{
				Id:				m.Id,
				Name:			m.Name,
				Slug:			m.Slug,
				Description:	m.Description,
				Address:		m.Address,
				Latitude:		m.Latitude,
				Longitude:		m.Longitude,
				Geometry:		m.Geometry, // stay string so FE free to parse
				Status:			m.Status,
				Image:			m.Image,
				CategoryID:		m.CategoryID,
			})
		}
		out = append(out, item)
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Public categories with maps retrieved successfully",
		Data: 	 out,
	})
}
package admin

import (
	"armandwipangestu/gis-api/database"
	"armandwipangestu/gis-api/helpers"
	"armandwipangestu/gis-api/models"
	"armandwipangestu/gis-api/structs"
	"net/http"

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

// Create new permission
func CreatePermission(c *gin.Context) {
	var req structs.PermissionCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors: helpers.TranslateErrorMessage(err),
		})

		return
	}

	permission := models.Permission{
		Name: req.Name,
	}

	if err := database.DB.Create(&permission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create permission",
			Errors: helpers.TranslateErrorMessage(err),
		})

		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Permission created successfully",
		Data: permission,
	})
}


// Get 1 permission based on id
func FindPermissionById(c *gin.Context) {
	// Get id from parameter
	id := c.Param("id")

	// Variable permission
	var permission models.Permission

	// Get permission based on id
	if err := database.DB.First(&permission, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Permission not found",
			Errors: helpers.TranslateErrorMessage(err),
		})

		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Permission Found",
		Data:	permission,
	})
}

// Update permission
func UpdatePermission(c *gin.Context) {
	id := c.Param("id")
	var permission models.Permission

	if err := database.DB.First(&permission, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Permission not found",
			Errors: helpers.TranslateErrorMessage(err),
		})

		return
	}

	var req structs.PermissionUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors: helpers.TranslateErrorMessage(err),
		})

		return
	}

	permission.Name = req.Name

	if err := database.DB.Save(&permission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update permission",
			Errors: helpers.TranslateErrorMessage(err),
		})

		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Permission updated successfully",
		Data: permission,
	})
}

// Delete permission
func DeletePermission(c *gin.Context) {
	// Get id from parameter
	id := c.Param("id")

	// Variable permission
	var permission models.Permission

	// Check if permission exists or not
	if err := database.DB.First(&permission, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Permission not found",
			Errors:	helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Delete all relationship role<->permission at pivot table
	if err := database.DB.Table("role_permissions").
			Where("permission_id = ?", id).
			Delete(nil).Error; err != nil {
				c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
					Success: false,
					Message: "Failed to detach permission from roles",
					Errors: helpers.TranslateErrorMessage(err),
				})

				return
			}

	// Delete permission
	if err := database.DB.Delete(&permission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete permission",
			Errors: helpers.TranslateErrorMessage(err),
		})

		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Permission deleted successfully",
	})
}

// Get all permissions
func FindAllPermissions(c *gin.Context) {
	// Initialize slice for accomodate data permissions
	var permissions []models.Permission

	// Get data permissions from database
	database.DB.Find(&permissions)

	// Send response success with data user
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Lists Data Permissions",
		Data: permissions,
	})
}
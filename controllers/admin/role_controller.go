package admin

import (
	"armandwipangestu/gis-api/database"
	"armandwipangestu/gis-api/helpers"
	"armandwipangestu/gis-api/models"
	"armandwipangestu/gis-api/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FindRoles(c *gin.Context) {
	var roles []models.Role
	var rolesResponse []structs.RoleResponse
	var total int64

	// Get parameter search, page, limit, and offset
	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseUrl := helpers.BuildBaseUrl(c)

	// Prepare the query
	query := database.DB.Preload("Permissions").Model(&models.Role{})
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Count total data and get data based on pagination
	query.Count(&total)
	err := query.Order("id desc").
		Limit(limit).
		Offset(offset).
		Find(&roles).Error; if err != nil {
			c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
				Success: false,
				Message: "Failed to fetch roles",
				Errors: helpers.TranslateErrorMessage(err),
			})

			return
		}

	// Mapping each role to RoleResponse
	for _, role := range roles {
		permissionResponses := []structs.PermissionResponse{} // always fill with empty slice

		for _, permission := range role.Permissions {
			permissionResponses = append(permissionResponses, structs.PermissionResponse{
				Id:			permission.Id,
				Name:		permission.Name,
				CreatedAt:	permission.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: 	permission.UpdatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		rolesResponse = append(rolesResponse, structs.RoleResponse{
			Id:				role.Id,
			Name:			role.Name,
			Permissions:	permissionResponses,
			CreatedAt: 		role.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: 		role.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// Send response with structure pagination
	helpers.PaginateResponse(c, rolesResponse, total, page, limit, baseUrl, search, "List Data Roles")
}

func CreateRole(c *gin.Context) {
	var req structs.RoleCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors: helpers.TranslateErrorMessage(err),
		})

		return
	}

	var permissions []models.Permission
	if len(req.PermissionIDs) > 0 {
		database.DB.Where("id IN ?", req.PermissionIDs).Find(&permissions)
	}

	role := models.Role{
		Name:			req.Name,
		Permissions:	permissions,
	}

	if err := database.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create role",
		})

		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Role created successfully",
		Data: role,
	})
}

func FindByRoleId(c *gin.Context) {
	// Get parameter ID from URL
	id := c.Param("id")

	// Initialize variable
	var role models.Role

	// Get data role with relationship permissions
	if err := database.DB.Preload("Permissions").First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Role not found",
			Errors:	 helpers.TranslateErrorMessage(err),
		})
		
		return
	}

	// Mapping relationship permissions to response structure
	permissionResponses := []structs.PermissionResponse{}
	for _, permission := range role.Permissions{
		permissionResponses = append(permissionResponses, structs.PermissionResponse{
			Id:				permission.Id,
			Name:			permission.Name,
			CreatedAt:		permission.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: 		permission.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// Prepare end response
	roleResponse := structs.RoleResponse{
		Id:				role.Id,
		Name:			role.Name,
		Permissions:	permissionResponses,
		CreatedAt: 		role.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:		role.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// Send response as JSON
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success:	true,
		Message:	"Role found",
		Data:		roleResponse,
	})
}

func UpdateRole(c *gin.Context) {
	// Get ID from parameter
	id := c.Param("id")

	// Initialize variable
	var role models.Role

	if err := database.DB.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Role not found",
			Errors: helpers.TranslateErrorMessage(err),
		})

		return
	}

	var req structs.RoleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:	 helpers.TranslateErrorMessage(err),
		})

		return
	}

	role.Name = req.Name

	var permission []models.Permission
	if len(req.PermissionIDs) > 0 {
		database.DB.Where("id IN ?", req.PermissionIDs).Find(&permission)
	}
	database.DB.Model(&role).Association("Permissions").Replace(&permission)

	if err := database.DB.Save(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update role",
			Errors: helpers.TranslateErrorMessage(err),
		})

		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Role updated successfully",
		Data:	 role,
	})
}

func DeleteRole(c *gin.Context) {
	// Get ID from parameter
	id := c.Param("id")

	// Initialize variable
	var role models.Role

	// Get data role
	if err := database.DB.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Role not found",
			Errors:	 helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Delete all relationship role<->permission at pivot table
	if err := database.DB.Model(&role).Association("Permissions").Clear(); err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to detach role from permissions",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Delete role
	if err := database.DB.Delete(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete role",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Send response
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Role deleted successfully",
	})
}
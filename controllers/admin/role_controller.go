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


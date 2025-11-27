package admin

import (
	"armandwipangestu/gis-api/database"
	"armandwipangestu/gis-api/helpers"
	"armandwipangestu/gis-api/models"
	"armandwipangestu/gis-api/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FindUsers(c *gin.Context) {
	// Initialize slice to accomodate data
	var users []models.User
	var usersResponse []structs.UserResponse

	// Initialize total
	var total int64

	// Get parameter search, page, limit, offset
	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseUrl := helpers.BuildBaseUrl(c)

	// Prepare the query
	query := database.DB.Preload("Roles").Model(&models.User{})
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Count total data and take data based on pagination
	query.Count(&total)
	err := query.Order("id desc").Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch users",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Mapping every user to UserResponse
	for _, user := range users {
		// Initialize empty slice, default is []
		roleResponses := []structs.RoleResponse{}

		// Mapping roles if exist
		for _, role := range user.Roles{
			roleResponses = append(roleResponses, structs.RoleResponse{
				Id:			role.Id,
				Name:		role.Name,
				CreatedAt:	role.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: 	role.UpdatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		// Append to list response
		usersResponse = append(usersResponse, structs.UserResponse{
			Id:			user.Id,
			Username:	user.Username,
			Name:		user.Name,
			Email:		user.Email,
			Roles:		roleResponses, // it will empty or `[]` if role not exist
			CreatedAt: 	user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: 	user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// Send response with pagination structure
	helpers.PaginateResponse(c, usersResponse, total, page, limit, baseUrl, search, "List Data Users")
}

func CreateUser(c *gin.Context) {
	// Struct user request
	var req = structs.UserCreateRequest{}

	// Bind JSON request to struct UserCreateRequest + validation
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:	 helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Get list role based on role_ids (if passing)
	var roles []models.Role
	if len(req.RoleIDs) > 0 {
		database.DB.Where("id IN ?", req.RoleIDs).Find(&roles)
	}

	// Initialize new user
	user := models.User{
		Name:		req.Name,
		Username:	req.Username,
		Email:		req.Email,
		Password:	helpers.HashPassword(req.Password),
		Roles:		roles,
	}

	// Save user to database
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create user",
			Errors:	 helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Send response success (mapping to UserResponse so that consistent)
	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "User created successfully",
		Data: 	 structs.UserResponse{
			Id:			user.Id,
			Name:		user.Name,
			Username: 	user.Username,
			Email:		user.Email,
			CreatedAt: 	user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: 	user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func FindUserById(c *gin.Context) {
	// Get ID user from parameter url
	id := c.Param("id")

	// Initialize user
	var user models.User

	// Find user based on id and preload relationship Roles
	if err := database.DB.Preload("Roles").First(&user, id).Error; err != nil {
		// If user not found, send response 404
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:	 helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Convert roles from model to struct RoleResponse
	var roleResponses []structs.RoleResponse
	for _, role := range user.Roles{
		roleResponses = append(roleResponses, structs.RoleResponse{
			Id: 		role.Id,
			Name:		role.Name,
			CreatedAt: 	role.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: 	role.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// Send response success with UserResponse
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "User Found",
		Data:	structs.UserResponse{
			Id:			user.Id,
			Name:		user.Name,
			Username:	user.Username,
			Email:		user.Email,
			Roles:		roleResponses,
			CreatedAt: 	user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: 	user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}
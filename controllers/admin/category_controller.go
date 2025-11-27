package admin

import (
	"armandwipangestu/gis-api/database"
	"armandwipangestu/gis-api/helpers"
	"armandwipangestu/gis-api/models"
	"armandwipangestu/gis-api/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Show all category with feature search and pagination
func FindCategories(c *gin.Context) {
	// Initialize slice for accomodate data
	var categories []models.Category
	var total int64

	// Get parameter search, page, limit, offset
	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseUrl := helpers.BuildBaseUrl(c)

	// Prepare query
	query := database.DB.Model(&models.Category{})
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Count total data and get data based on pagination
	query.Count(&total)
	err := query.Order("id desc").Limit(limit).Offset(offset).Find(&categories).Error; if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success:	false,
			Message:	"Failed to fetch categories",
			Errors:		helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Send response with pagination structure
	helpers.PaginateResponse(c, categories, total, page, limit, baseUrl, search, "List Data Categories")
}

// Add new category
func CreateCategory(c *gin.Context) {
	var req structs.CategoryCreateRequest

	// Input validation
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:	 helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Get file image from form
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:	 map[string]string{"Image": "Image is required"},
		})

		return
	}

	// Upload file using helper
	uploadResult := helpers.UploadFile(c, helpers.UploadConfig{
		File:			file,
		AllowedTypes: 	[]string{".jpg", "jpeg", ".png", ".gif"},
		MaxSize: 		10 << 20, // Maximum size is 10MB
		DestinationDir: "public/uploads/categories",
	})
	if uploadResult.Response != nil {
		c.JSON(http.StatusBadRequest, uploadResult.Response)
		return
	}

	// Create slug based on category name
	slug := helpers.Slugify(req.Name)

	// Check if slug is used or not
	var existing models.Category
	if err := database.DB.Where("slug = ?", slug).First(&existing).Error; err == nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:	 map[string]string{
				"slug": "Slug already exists",
			},
		})

		return
	}

	// Create object category (based on new field model)
	category := models.Category{
		Name:			req.Name,
		Slug:			slug,
		Image:			uploadResult.FileName,
		Color:			req.Color,
		Description: 	req.Description,
	}

	// Save category
	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create category",
		})

		return
	}

	// Send success response
	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Category created successfully",
		Data:	 category,
	})
}

// Get 1 category based on id
func FindCategoryById(c *gin.Context) {
	// Get parameter id
	id := c.Param("id")

	// Initialize variable
	var category models.Category

	// Find category based on id
	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Category not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Send data category
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Category found",
		Data:	 category,
	})
}
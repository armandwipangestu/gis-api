package admin

import (
	"armandwipangestu/gis-api/database"
	"armandwipangestu/gis-api/helpers"
	"armandwipangestu/gis-api/models"
	"armandwipangestu/gis-api/structs"
	"net/http"
	"os"
	"path/filepath"

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

// Update data category
func UpdateCategory(c *gin.Context) {
	// Get parameter id
	id := c.Param("id")

	// Initialize category
	var category models.Category

	// Check if category with that id is exist or not
	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Category not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Struct category update request
	var req structs.CategoryUpdateRequest

	// Input validation
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Save path old image (if exist)
	oldImagePath := ""
	if category.Image != "" {
		oldImagePath = filepath.Join("public", "uploads", "categories", category.Image)
	}

	// If user upload new image
	file, err := c.FormFile("image")
	if err == nil {
		// Upload new image
		uploadResult := helpers.UploadFile(c, helpers.UploadConfig{
			File:			file,
			AllowedTypes: 	[]string{".jpg", ".jpeg", ".png", ".gif"},
			MaxSize: 		10 << 20, // Maximum file upload is 10MB
			DestinationDir: "public/uploads/categories",
		})

		// If upload is fail, return error
		if uploadResult.Response != nil {
			c.JSON(http.StatusBadRequest, uploadResult.Response)
			
			return
		}

		// Set new filename
		category.Image = uploadResult.FileName
	}

	// Update data
	category.Name = req.Name
	category.Slug = helpers.Slugify(req.Name)
	category.Color = req.Color
	category.Description = req.Description

	// Save changes to database
	if err := database.DB.Save(&category).Error; err != nil {
		// If save fail and has new file, remove new file, to make not orphan
		if file != nil && category.Image != "" {
			newImagePath := filepath.Join("public", "uploads", "categories", category.Image)
			os.Remove(newImagePath)
		}
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "failed to update category",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	// If has new file and old file still exist, remove old image
	if file != nil && oldImagePath != "" {
		_ = os.Remove(oldImagePath)
	}

	// Send success response
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Category updated successfully",
		Data:	 category,
	})
}

// Remove category based on id
func DeleteCategory(c *gin.Context) {
	// Get parameter id
	id := c.Param("id")

	// Initialize category
	var category models.Category

	// Check if data category found
	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Category not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Save path image to delete
	imagePath := ""
	if category.Image != "" {
		imagePath = filepath.Join("public", "uploads", "categories", category.Image)
	}

	// Delete data category
	if err := database.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete category",
			Errors:	 helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Remove file image if exist
	if imagePath != "" {
		if err := os.Remove(imagePath); err != nil && !os.IsNotExist(err) {
			// Failed to delete image, but delete category still considered successfully deleted
			c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
				Success: false,
				Message: "Category deleted but failed to remove image",
				Errors:	 map[string]string{"image": "Failed to remove image file: " + err.Error()},
			})

			return
		}
	}

	// Response success
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Category deleted successfully",
	})
}
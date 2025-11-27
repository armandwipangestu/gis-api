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

// Show all map with feature search and pagination
func FindMaps(c *gin.Context) {
	// Initialize slice to accomodate data
	var maps []models.Map
	var total int64

	// Get parameter search, page, limit, offset
	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseUrl := helpers.BuildBaseUrl(c)

	// Prepare query
	query := database.DB.Preload("Category").Model(&models.Map{})
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Count total data and take data based on pagination
	query.Count(&total)
	err := query.Order("id desc").Limit(limit).Offset(offset).Find(&maps).Error; if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch maps",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Send response with structure pagination
	helpers.PaginateResponse(c, maps, total, page, limit, baseUrl, search, "List Data Maps")
}

// Create new map
func CreateMap(c *gin.Context) {
	// Initialize variable
	var req structs.MapCreateRequest

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
			Errors:  map[string]string{"image": "Image is required"},
		})

		return
	}

	// Upload file using helper
	uploadResult := helpers.UploadFile(c, helpers.UploadConfig{
		File:			file,
		AllowedTypes: 	[]string{".jpg", ".jpeg", ".png", ".gif"},
		MaxSize:		10 << 20, // Maximum file upload is 10MB
		DestinationDir: "public/uploads/maps",
	})
	if uploadResult.Response != nil {
		c.JSON(http.StatusBadRequest, uploadResult.Response)

		return
	}

	// Create slug based on name
	slug := helpers.Slugify(req.Name)

	// Check if slug already used
	var existing models.Map
	if err := database.DB.Where("slug = ?", slug).First(&existing).Error; err == nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  map[string]string{
				"slug": "Slug already exists",
			},
		})

		return
	}

	// Create object map (based on field model)
	mp := models.Map{
		Image:			uploadResult.FileName,
		Name:			req.Name,
		Slug:			slug,
		Description: 	req.Description,
		Address: 		req.Address,
		Latitude: 		req.Latitude,
		Longitude: 		req.Longitude,
		Geometry: 		req.Geometry,
		Status: 		req.Status, // make default DB if empty
		CategoryID: 	req.CategoryID,
	}

	// Save map
	if err := database.DB.Create(&mp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create map",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Send response success
	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Map created successfully",
		Data: 	 mp,
	})
}

// Get one map based on id
func FindMapById(c *gin.Context) {
	// Get parameter id
	id := c.Param("id")

	// Initialize variable
	var mp models.Map

	// Find map based on id
	if err := database.DB.First(&mp, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Map not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Send data map
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Map Found",
		Data: 	 structs.PublicMap{
			Id: 			mp.Id,
			Image:			mp.Image,
			Name:			mp.Name,
			Slug:			mp.Slug,
			Description:	mp.Description,
			Address:		mp.Address,
			Latitude:		mp.Latitude,
			Longitude:		mp.Longitude,
			Geometry:		mp.Geometry,
			Status:			mp.Status,
			CategoryID:		mp.CategoryID,
		},
	})
}

// Update data map
func UpdateMap(c *gin.Context) {
	// Get parameter id
	id := c.Param("id")

	// Initialize map
	var mp models.Map

	// Check if map with given id is exist or not
	if err := database.DB.First(&mp, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Map not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Struct map update request
	var req structs.MapUpdateRequest

	// Input validation
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:	 helpers.TranslateErrorMessage(err),
		})

		return
	}

	// Save path old image (if exist)
	oldImagePath := ""
	if mp.Image != "" {
		oldImagePath = filepath.Join("public", "uploads", "maps", mp.Image)
	}

	// If user upload new image
	file, err := c.FormFile("image")
	if err == nil {
		// Upload new image
		uploadResult := helpers.UploadFile(c, helpers.UploadConfig{
			File:			file,
			AllowedTypes: 	[]string{".jpg", ".jpeg", ".png", ".gif"},
			MaxSize: 		10 << 20, // Maximum upload file size is 10MB
			DestinationDir: "public/uploads/maps",
		})

		// If upload failed, return error
		if uploadResult.Response != nil {
			c.JSON(http.StatusBadRequest, uploadResult.Response)
			
			return
		}

		// Set new file name
		mp.Image = uploadResult.FileName
	}

	// Update data
	mp.Name = req.Name
	mp.Slug = helpers.Slugify(req.Name)
	mp.Description = req.Description
	mp.Address = req.Address
	mp.Latitude = req.Latitude
	mp.Longitude = req.Longitude
	mp.Geometry = req.Geometry
	mp.Status = req.Status
	mp.CategoryID = req.CategoryID

	// Save changes to database
	if err := database.DB.Save(&mp).Error; err != nil {
		// If save failed and has new file, remove new file, to make not orphan
		if file != nil && mp.Image != "" {
			newImagePath := filepath.Join("public", "uploads", "maps", mp.Image)
			_ = os.Remove(newImagePath)
		}
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update map",
			Errors:	 helpers.TranslateErrorMessage(err),
		})

		return
	}

	// If has new file and old file still exist, remove the old file
	if file != nil && oldImagePath != "" {
		_ = os.Remove(oldImagePath)
	}

	// Send success response
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Map updated successfully",
		Data:	 structs.PublicMap{
			Id: 					mp.Id,
			Image: 					mp.Image,
			Name: 					mp.Name,
			Slug: 					mp.Slug,
			Description: 			mp.Description,
			Address: 				mp.Address,
			Latitude: 				mp.Latitude,
			Longitude: 				mp.Longitude,
			Geometry: 				mp.Geometry,
			Status: 				mp.Status,
			CategoryID: 			mp.CategoryID,
		},
	})
}
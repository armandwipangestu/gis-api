package helpers

import (
	"armandwipangestu/gis-api/structs"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadConfig struct {
	File				*multipart.FileHeader	// File to be uploaded
	AllowedTypes		[]string				// Extension file whitelist
	MaxSize				int64					// Maximum file size (in byte)
	DestinationDir		string					// Path folder saved
}

type UploadResult struct {
	FileName			string
	FilePath			string
	Error				error
	Response			*structs.ErrorResponse
}

// For filename has been slugify
func SlugifyFilename(filename string) string {
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)

	// Just slugify base name
	slugBase := Slugify(base)

	// Merge with real extension
	return slugBase + ext
}

func UploadFile(c *gin.Context, config UploadConfig) UploadResult {
	// Check if file exist
	if config.File == nil {
		return UploadResult{
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "File is required",
				Errors: map[string]string{"file": "No file was uploaded"},
			},
		}
	}

	// Check file size
	if config.File.Size > config.MaxSize {
		return UploadResult{
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "File to large",
				Errors: map[string]string{"file": fmt.Sprintf("Maximum file size is: %dMB", config.MaxSize/(1<<20))},
			},
		}
	}

	// Check file type
	ext := strings.ToLower(filepath.Ext(config.File.Filename))
	allowed := false
	for _, t := range config.AllowedTypes {
		if ext == t {
			allowed = true
			break
		}
	}
	if !allowed {
		return UploadResult{
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "Invalid file type",
				Errors: map[string]string{"file": fmt.Sprintf("Allowed file types: %v", config.AllowedTypes)},
			},
		}
	}

	// Generate UUID as filename
	uuidName := uuid.New().String()

	// Saved real extension name
	filename := uuidName + ext
	filePath := filepath.Join(config.DestinationDir, filename)

	// Create destination folder if not exist
	if err := os.MkdirAll(config.DestinationDir, 0755); err != nil {
		return UploadResult{
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "Failed to create upload directory",
				Errors: map[string]string{"system": err.Error()},
			},
		}
	}

	// Save file to destination folder
	if err := c.SaveUploadedFile(config.File, filePath); err != nil {
		return UploadResult{
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "Failed to save file",
				Errors:	map[string]string{"file": err.Error()},
			},
		}
	}

	// Return the result upload
	return UploadResult{
		FileName: filename,
		FilePath: filePath,
	}
}
package helpers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// TranslateErrorMessage handle validation error and database to easy of map
func TranslateErrorMessage(err error) map[string]string {
	errorsMap := make(map[string]string)

	// Handle validation for validator.v10
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := fieldError.Field()
			switch fieldError.Tag() {
			case "required":
				errorsMap[field] = fmt.Sprintf("%s is required", field)
			case "email":
				errorsMap[field] = "Invalid email format"
			case "unique":
				errorsMap[field] = fmt.Sprintf("%s already exists", field)
			case "min":
				errorsMap[field] = fmt.Sprintf("%s must be at least %s characters", field, fieldError.Param())
			case "max":
				errorsMap[field] = fmt.Sprintf("%s must be at most %s characters", field, fieldError.Param())
			case "numeric":
				errorsMap[field] = fmt.Sprintf("%s must be a number", field)
			default:
				errorsMap[field] = "Invalid value"
			}
		}
	}

	// Handle GORM error: duplicate entry
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			field := extractDuplicateField(err.Error())
			if field != "" {
				errorsMap[field] = fmt.Sprintf("%s already exists", field)
			} else {
				errorsMap[field] = "Duplicate entry"
			}
		} else if err == gorm.ErrRecordNotFound {
			errorsMap["Error"] = "Record not found"
		}
	}

	return errorsMap
}

// IsDuplicateEntryError check if error is duplicate entry
func IsDuplicateEntryError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Duplicate entry")
}

// extractDuplicateField try to extract column name from error "Duplicate Entry"
func extractDuplicateField(errMsg string) string {
	// Example error MySQL: Error 1002: Duplicate entry 'test@example.com' for key 'users.email'
	// we try to get after 'for key' then extract the field name
	re := regexp.MustCompile(`for key '(\w+\.)?(\w+)'`)
	matches := re.FindStringSubmatch(errMsg)
	if len(matches) == 3 {
		// Return capitalize field name
		return strings.Title(matches[2])
	}
	return ""
}
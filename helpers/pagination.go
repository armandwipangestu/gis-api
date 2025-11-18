package helpers

import (
	"armandwipangestu/gis-api/structs"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Structure for link pagination (like in Laravel)
type PaginationLink struct {
	Url			string			`json:"url"`
	Label		string			`json:"label"`
	Active		bool			`json:"active"`
}

// Convertion string to integer with fallback to 1 if failed or < 1
func StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil || i < 1 {
		return 1
	}
	return i
}

// Count total page based on total data and ofset per page
func TotalPage(total int64, perPage int) int {
	if perPage == 0 {
		return 1
	}

	pages := int(total) / perPage
	
	if int(total)%perPage != 0 {
		pages++
	}

	return pages
}

// Create pagination link like in Laravel
func BuildPaginationLinks(currentPage, lastPage int, baseUrl, search string) []PaginationLink {
	links := []PaginationLink{}

	// Link to previous page
	links = append(links, PaginationLink{
		Url: PageUrl(baseUrl, currentPage-1, lastPage, search),
		Label: "&laquo; Previous",
		Active: false,
	})

	// Link to all page
	for i := 1; i <= lastPage; i++ {
		links = append(links, PaginationLink{
			Url: baseUrl + "?page=" + strconv.Itoa(i) + QueryString(search),
			Label: strconv.Itoa(i),
			Active: i == currentPage,
		})
	}

	// Link to next page
	links = append(links, PaginationLink{
		Url: PageUrl(baseUrl, currentPage+1, lastPage, search),
		Label: "Next &raquo;",
		Active: false,
	})

	return links
}

// Create Url for specific page (previous/next)
func PageUrl(baseUrl string, page, lastPage int, search string) string {
	if page < 1 || page > lastPage {
		return ""
	}
	return baseUrl + "?page=" + strconv.Itoa(page) + QueryString(search)
}

// Create query string if has search parameter
func QueryString(search string) string {
	if search == "" {
		return ""
	}
	return "&search=" + search
}

// Get search parameter and pagination from query string
func GetPaginationParams(c *gin.Context) (search string, page, limit, offset int) {
	search = c.Query("search")
	page = StringToInt(c.DefaultQuery("page", "1"))
	limit = StringToInt(c.DefaultQuery("limit", "5"))
	offset = (page - 1) * limit
	return
}

// Build foundation of basic url for pagination (considering HTTP/HTTPS and proxy)
func BuildBaseUrl(c *gin.Context) string {
	// Check X-Forwarded-Proto from reverse proxy if exist
	scheme := c.Request.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		if c.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}

	return scheme + "://" + c.Request.Host + c.Request.URL.Path
}

// Send response full JSON pagination like in laravel
func PaginateResponse(c *gin.Context, data interface{}, total int64, page, limit int, baseUrl, search, message string) {
	lastPage := TotalPage(total, limit)
	from := (page-1)*limit + 1
	to := from + reflect.ValueOf(data).Len() - 1

	links := BuildPaginationLinks(page, lastPage, baseUrl, search)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: message,
		Data: gin.H{
			"current_page":		page,
			"data": 			data,
			"first_page_url": 	baseUrl + "?page=1" + QueryString(search),
			"from":				from,
			"last_page": 		lastPage,
			"last_page_url":	baseUrl + "?page=" + strconv.Itoa(lastPage) + QueryString(search),
			"links": 			links,
			"next_page_url": 	PageUrl(baseUrl, page+1, lastPage, search),
			"path":				baseUrl,
			"per_page":			limit,
			"prev_page_url": 	PageUrl(baseUrl, page-1, lastPage, search),
			"to":				to,
			"total":			total,
		},
	})
}
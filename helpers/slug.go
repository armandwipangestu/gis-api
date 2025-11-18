package helpers

import (
	"regexp"
	"strings"
)

func Slugify(text string) string {
	// Change to lowercase
	slug := strings.ToLower(text)

	// Delete character non-alfanumeric except space
	re := regexp.MustCompile(`[^a-z0-9\s-]`)
	slug = re.ReplaceAllString(slug, "")

	// Change space and strip `-` double with single strip
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")

	return slug
}
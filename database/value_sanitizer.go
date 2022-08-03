package database

import (
	"regexp"
	"strings"
)

func SanitizeValue(value string) string {
	result := strings.ReplaceAll(value, " ", "_")

	reg := regexp.MustCompile("[^a-zA-Z0-9,_-]+")
	result = reg.ReplaceAllString(result, "")

	return result
}

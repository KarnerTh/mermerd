package database

import (
	"regexp"
	"strings"
)

func SanitizeColumnType(columnType string) string {
	switch {
	case columnType == "character varying", columnType == "varchar":
		return "string"
	case strings.Contains(columnType, "timestamp"):
		return "date"
	default:
		return sanitize(strings.ReplaceAll(columnType, " ", "_"))
	}
}

func SanitizeTableName(tableName string) string {
	return strings.ReplaceAll(sanitize(tableName), " ", "_")
}

func SanitizeColumnName(columnName string) string {
	return strings.ReplaceAll(sanitize(columnName), " ", "_")
}

func sanitize(value string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9_-]+")
	return reg.ReplaceAllString(value, "")
}

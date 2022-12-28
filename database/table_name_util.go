package database

import (
	"errors"
	"strings"
)

func ParseTableName(value string, selectedSchemas []string) (TableDetail, error) {
	parts := strings.Split(value, ".")

	if len(parts) == 1 {
		if len(selectedSchemas) != 1 {
			return TableDetail{}, errors.New("If table names do not specify the schema, exactly one selected schema should be present")
		}

		return TableDetail{Schema: selectedSchemas[0], Name: parts[0]}, nil
	}

	if len(parts) == 2 {
		return TableDetail{Schema: parts[0], Name: parts[1]}, nil
	}

	return TableDetail{}, errors.New("Could not parse table name")
}

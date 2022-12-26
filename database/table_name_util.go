package database

import (
	"errors"
	"strings"
)

func ParseTableName(value string, selectedSchemas []string) (TableNameResult, error) {
	parts := strings.Split(value, ".")

	if len(parts) == 1 {
		if len(selectedSchemas) != 1 {
			return TableNameResult{}, errors.New("Could not parse table name")
		}

		return TableNameResult{Schema: selectedSchemas[0], Name: parts[0]}, nil
	}

	if len(parts) == 2 {
		return TableNameResult{Schema: parts[0], Name: parts[1]}, nil
	}

	return TableNameResult{}, errors.New("Could not parse table name")
}

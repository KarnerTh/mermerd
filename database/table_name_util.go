package database

import (
	"errors"
	"strings"
)

func ParseTableName(value string) (TableNameResult, error) {
	parts := strings.Split(value, ".")

	if len(parts) == 1 {
		return TableNameResult{Schema: "", Name: parts[0]}, nil
	}

	if len(parts) == 2 {
		return TableNameResult{Schema: parts[0], Name: parts[1]}, nil
	}

	return TableNameResult{}, errors.New("Could not parse table name")
}

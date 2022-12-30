package database

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTableName(t *testing.T) {
	testCases := []struct {
		value           string
		selectedSchemas []string
		expectedResult  TableDetail
		shouldFail      bool
	}{
		{
			value:           "public.tableA",
			selectedSchemas: []string{},
			expectedResult:  TableDetail{Schema: "public", Name: "tableA"},
			shouldFail:      false,
		},
		{
			value:           "tableA",
			selectedSchemas: []string{"public"},
			expectedResult:  TableDetail{Schema: "public", Name: "tableA"},
			shouldFail:      false,
		},
		{
			value:           "tableA",
			selectedSchemas: []string{},
			expectedResult:  TableDetail{},
			shouldFail:      true,
		},
		{
			value:           "tableA",
			selectedSchemas: []string{"public", "other_db"},
			expectedResult:  TableDetail{},
			shouldFail:      true,
		},
		{
			value:           "public.other_db.tableA",
			selectedSchemas: []string{"public", "other_db"},
			expectedResult:  TableDetail{},
			shouldFail:      true,
		},
	}

	for index, testCase := range testCases {
		t.Run(fmt.Sprintf("run #%d", index), func(t *testing.T) {
			// Arrange
			value := testCase.value
			selectedSchemas := testCase.selectedSchemas

			// Act
			result, err := ParseTableName(value, selectedSchemas)

			// Assert
			if testCase.shouldFail {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, testCase.expectedResult, result)
		})
	}
}

package database

import (
	"fmt"
	"testing"
)

func TestDbType_String(t *testing.T) {
	testCases := []struct {
		dbType         DbType
		expectedResult string
	}{
		{
			dbType:         Postgres,
			expectedResult: "pgx",
		},
		{
			dbType:         MySql,
			expectedResult: "mysql",
		},
	}

	for index, testCase := range testCases {
		t.Run(fmt.Sprintf("run #%d", index), func(t *testing.T) {
			// Arrange
			dbType := testCase.dbType

			// Act
			result := dbType.String()

			// Assert
			if result != testCase.expectedResult {
				t.Errorf("Expected %s, got %s", testCase.expectedResult, result)
			}
		})
	}
}

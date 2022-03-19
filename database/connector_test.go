package database

import (
	"fmt"
	"testing"
)

func TestNewConnector(t *testing.T) {
	testCases := []struct {
		connectionString string
		expectedDbType   DbType
	}{
		{
			connectionString: "postgresql://user:password@localhost:5432/yourDb",
			expectedDbType:   Postgres,
		},
		{
			connectionString: "postgres://user:password@localhost:5432/yourDb",
			expectedDbType:   Postgres,
		},
		{
			connectionString: "mysql://root:password@tcp(127.0.0.1:3306)/yourDb",
			expectedDbType:   MySql,
		},
	}

	for index, testCase := range testCases {
		t.Run(fmt.Sprintf("run #%d", index), func(t *testing.T) {
			// Arrange
			connectionString := testCase.connectionString

			// Act
			resultConnector, resultError := NewConnector(connectionString)

			// Assert
			if resultError != nil {
				t.Errorf("Should not throw error")
			}
			if resultConnector.GetDbType() != testCase.expectedDbType {
				t.Errorf("Expected %s, got %s", testCase.expectedDbType, resultConnector)
			}
		})
	}
}

func TestUnsupportedConnector(t *testing.T) {
	// Arrange
	connectionString := "notSupported://user:password@localhost:5432/yourDb"

	// Act
	resultConnector, resultError := NewConnector(connectionString)

	// Assert
	if resultError == nil {
		t.Errorf("Should throw error")
	}
	if resultConnector != nil {
		t.Errorf("Connector should be nil")
	}
}

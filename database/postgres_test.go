package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const connectionString = "postgresql://user:password@localhost:5432/mermerd_test"

func getConnection() postgresConnector {
	return postgresConnector{
		dbType:           Postgres,
		connectionString: connectionString,
		db:               nil,
	}
}

func getConnectionAndConnect() postgresConnector {
	connector := getConnection()
	_ = connector.Connect()
	return connector
}

func TestPostgresConnector_Connect(t *testing.T) {
	// Arrange
	connector := getConnection()

	// Act
	err := connector.Connect()

	// Assert
	assert.Nil(t, err)
}

func TestPostgresConnector_GetSchemas(t *testing.T) {
	// Arrange
	connector := getConnectionAndConnect()

	// Act
	schemas, err := connector.GetSchemas()

	// Assert
	assert.Nil(t, err)
	assert.NotEmpty(t, schemas)
}

func TestPostgresConnector_GetTables(t *testing.T) {
	// Arrange
	connector := getConnectionAndConnect()
	schema := "public"

	// Act
	tables, err := connector.GetTables(schema)

	// Assert
	expectedResult := []string{"article", "article_detail", "article_comment", "label", "article_label"}
	assert.Nil(t, err)
	assert.ElementsMatch(t, expectedResult, tables)
}

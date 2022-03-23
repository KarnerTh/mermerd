package database

import (
	"fmt"
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

func TestPostgresConnector_GetColumns(t *testing.T) {
	connector := getConnectionAndConnect()
	testCases := []struct {
		tableName       string
		expectedColumns []string
	}{
		{tableName: "article", expectedColumns: []string{"id", "title"}},
		{tableName: "article_detail", expectedColumns: []string{"id", "created_at"}},
		{tableName: "article_comment", expectedColumns: []string{"id", "article_id", "comment"}},
		{tableName: "label", expectedColumns: []string{"id", "label"}},
		{tableName: "article_label", expectedColumns: []string{"article_id", "label_id"}},
	}

	for index, testCase := range testCases {
		t.Run(fmt.Sprintf("run #%d", index), func(t *testing.T) {
			// Arrange
			tableName := testCase.tableName
			var columnNames []string

			// Act
			columns, err := connector.GetColumns(tableName)

			// Assert
			for _, column := range columns {
				columnNames = append(columnNames, column.Name)
			}

			assert.Nil(t, err)
			assert.ElementsMatch(t, testCase.expectedColumns, columnNames)
		})
	}
}

func TestPostgresConnector_GetConstraints(t *testing.T) {
	connector := getConnectionAndConnect()

	t.Run("One-to-one relation", func(t *testing.T) {
		// Arrange
		tableName := "article_detail"

		// Act
		constraintResults, err := connector.GetConstraints(tableName)

		// Assert
		assert.Nil(t, err)
		assert.Len(t, constraintResults, 1)
		constraint := constraintResults[0]
		assert.True(t, constraint.IsPrimary)
		assert.False(t, constraint.HasMultiplePK)
	})

	t.Run("Many-to-one relation #1", func(t *testing.T) {
		// Arrange
		tableName := "article_comment"

		// Act
		constraintResults, err := connector.GetConstraints(tableName)

		// Assert
		assert.Nil(t, err)
		assert.Len(t, constraintResults, 1)
		constraint := constraintResults[0]
		assert.False(t, constraint.IsPrimary)
		assert.False(t, constraint.HasMultiplePK)
	})

	t.Run("Many-to-one relation #2", func(t *testing.T) {
		// Arrange
		pkTableName := "article"
		fkTableName := "article_label"

		// Act
		constraintResults, err := connector.GetConstraints(pkTableName)

		// Assert
		assert.Nil(t, err)
		var constraint *ConstraintResult
		for _, item := range constraintResults {
			if item.FkTable == fkTableName {
				constraint = &item
			}
		}
		assert.NotNil(t, constraint)
		assert.True(t, constraint.IsPrimary)
		assert.True(t, constraint.HasMultiplePK)
	})
}

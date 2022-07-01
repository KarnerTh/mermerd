package database

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseIntegrations(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	connectorFactory := NewConnectorFactory()
	testCases := []struct {
		dbType           DbType
		connectionString string
		schema           string
	}{
		{
			dbType:           Postgres,
			connectionString: "postgresql://user:password@localhost:5432/mermerd_test",
			schema:           "public",
		},
		{
			dbType:           MySql,
			connectionString: "mysql://user:password@tcp(127.0.0.1:3306)/mermerd_test",
			schema:           "mermerd_test",
		},
	}

	for _, testCase := range testCases {
		connector, _ := connectorFactory.NewConnector(testCase.connectionString)

		getConnectionAndConnect := func() Connector {
			_ = connector.Connect()
			return connector
		}

		t.Run(testCase.dbType.String(), func(t *testing.T) {
			t.Run("Connect", func(t *testing.T) {
				// Arrange
				connector := connector

				// Act
				err := connector.Connect()

				// Assert
				assert.Nil(t, err)
			})

			t.Run("GetSchemas", func(t *testing.T) {
				// Arrange
				connector := getConnectionAndConnect()

				// Act
				schemas, err := connector.GetSchemas()

				// Assert
				assert.Nil(t, err)
				assert.NotEmpty(t, schemas)
			})

			t.Run("GetTables", func(t *testing.T) {
				// Arrange
				connector := getConnectionAndConnect()
				schema := testCase.schema

				// Act
				tables, err := connector.GetTables(schema)

				// Assert
				expectedResult := []string{
					"article",
					"article_detail",
					"article_comment",
					"label",
					"article_label",
					"test_1_a",
					"test_1_b",
				}
				assert.Nil(t, err)
				assert.ElementsMatch(t, expectedResult, tables)
			})

			t.Run("GetColumns", func(t *testing.T) {
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
					{tableName: "test_1_a", expectedColumns: []string{"id", "xid"}},
					{tableName: "test_1_b", expectedColumns: []string{"aid", "bid"}},
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
			})

			t.Run("GetConstraints", func(t *testing.T) {
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
							break
						}
					}
					assert.NotNil(t, constraint)
					assert.True(t, constraint.IsPrimary)
					assert.True(t, constraint.HasMultiplePK)
				})

				// Multiple primary keys (https://github.com/KarnerTh/mermerd/issues/8)
				t.Run("Test 1 (Issue #8)", func(t *testing.T) {
					// Arrange
					pkTableName := "test_1_b"

					// Act
					constraintResults, err := connector.GetConstraints(pkTableName)

					// Assert
					assert.Nil(t, err)
					assert.NotNil(t, constraintResults)
					assert.Len(t, constraintResults, 2)
					assert.True(t, constraintResults[0].IsPrimary)
					assert.True(t, constraintResults[0].HasMultiplePK)
					assert.Equal(t, constraintResults[0].ColumnName, "aid")
					assert.True(t, constraintResults[1].IsPrimary)
					assert.True(t, constraintResults[1].HasMultiplePK)
					assert.Equal(t, constraintResults[1].ColumnName, "bid")
				})
			})
		})
	}
}

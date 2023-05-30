package database

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestPostgresEnums(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Arrange
	var enumValues string

	// Act
	connector, _ := NewConnectorFactory().NewConnector(testConnectionPostgres.connectionString)
	if err := connector.Connect(); err != nil {
		logrus.Error(err)
		t.FailNow()
	}
	columns, err := connector.GetColumns(TableDetail{Schema: "public", Name: "test_2_enum"})

	// Assert
	for _, column := range columns {
		if column.Name == "fruit" {
			enumValues = column.EnumValues
		}
	}

	assert.Nil(t, err)
	assert.Equal(t, "apple,banana", enumValues)
}

func TestPostgresNotUniqueConstraintName(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Arrange
	tableName := TableDetail{Schema: "public", Name: "test_not_unique_constraint_name_b"}
	connector, _ := NewConnectorFactory().NewConnector(testConnectionPostgres.connectionString)
	if err := connector.Connect(); err != nil {
		logrus.Error(err)
		t.FailNow()
	}

	// Act
	_, err := connector.GetConstraints(tableName)

	// Assert
	// we only need to check if an error is thrown
	assert.Nil(t, err)
}

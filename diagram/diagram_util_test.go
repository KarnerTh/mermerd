package diagram

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/KarnerTh/mermerd/database"
	"github.com/KarnerTh/mermerd/mocks"
)

func TestGetRelation(t *testing.T) {
	testCases := []struct {
		isPrimary        bool
		hasMultiplePK    bool
		expectedRelation ErdRelationType
	}{
		{true, true, relationManyToOne},
		{false, true, relationManyToOne},
		{false, false, relationManyToOne},
		{true, false, relationOneToOne},
	}

	for index, testCase := range testCases {
		t.Run(fmt.Sprintf("run #%d", index), func(t *testing.T) {
			// Arrange
			constraint := database.ConstraintResult{
				FkTable:        "tableA",
				PkTable:        "tableB",
				ConstraintName: "constraintXY",
				IsPrimary:      testCase.isPrimary,
				HasMultiplePK:  testCase.hasMultiplePK,
			}

			// Act
			result := getRelation(constraint)

			// Assert
			assert.Equal(t, testCase.expectedRelation, result)
		})
	}
}

func TestGetAttributeKey(t *testing.T) {
	testCases := []struct {
		column                  database.ColumnResult
		expectedAttributeResult ErdAttributeKey
	}{
		{
			column: database.ColumnResult{
				Name:      "",
				DataType:  "",
				IsPrimary: true,
				IsForeign: false,
			},
			expectedAttributeResult: primaryKey,
		},
		{
			column: database.ColumnResult{
				Name:      "",
				DataType:  "",
				IsPrimary: false,
				IsForeign: true,
			},
			expectedAttributeResult: foreignKey,
		},
		{
			column: database.ColumnResult{
				Name:      "",
				DataType:  "",
				IsPrimary: true,
				IsForeign: true,
			},
			expectedAttributeResult: primaryKey,
		},
		{
			column: database.ColumnResult{
				Name:      "",
				DataType:  "",
				IsPrimary: false,
				IsForeign: false,
			},
			expectedAttributeResult: none,
		},
	}

	for index, testCase := range testCases {
		t.Run(fmt.Sprintf("run #%d", index), func(t *testing.T) {
			// Arrange
			column := testCase.column

			// Act
			result := getAttributeKey(column)

			// Assert
			assert.Equal(t, testCase.expectedAttributeResult, result)
		})
	}
}

func TestTableNameInSlice(t *testing.T) {
	t.Run("Existing item should be found", func(t *testing.T) {
		// Arrange
		tableName := "testTable"
		slice := []ErdTableData{{Name: tableName}}

		// Act
		result := tableNameInSlice(slice, tableName)

		// Assert
		assert.True(t, result)
	})

	t.Run("Missing item should not be found", func(t *testing.T) {
		// Arrange
		tableName := "testTable"
		slice := []ErdTableData{{Name: "notTheTableName"}}

		// Act
		result := tableNameInSlice(slice, tableName)

		// Assert
		assert.False(t, result)
	})
}

func TestGetColumnData(t *testing.T) {
	columnName := "testColumn"
	enumValues := "a,b"
	column := database.ColumnResult{
		Name:       columnName,
		IsPrimary:  true,
		EnumValues: enumValues,
	}

	t.Run("Get all fields", func(t *testing.T) {
		// Arrange
		configMock := mocks.MermerdConfig{}
		configMock.On("OmitAttributeKeys").Return(false).Once()
		configMock.On("ShowEnumValues").Return(true).Once()

		// Act
		result := getColumnData(&configMock, column)

		// Assert
		configMock.AssertExpectations(t)
		assert.Equal(t, columnName, result.Name)
		assert.Equal(t, enumValues, result.EnumValues)
		assert.Equal(t, primaryKey, result.AttributeKey)
	})

	t.Run("Get all fields except enum values", func(t *testing.T) {
		// Arrange
		configMock := mocks.MermerdConfig{}
		configMock.On("OmitAttributeKeys").Return(false).Once()
		configMock.On("ShowEnumValues").Return(false).Once()

		// Act
		result := getColumnData(&configMock, column)

		// Assert
		configMock.AssertExpectations(t)
		assert.Equal(t, columnName, result.Name)
		assert.Equal(t, "", result.EnumValues)
		assert.Equal(t, primaryKey, result.AttributeKey)
	})

	t.Run("Get all fields except attribute key", func(t *testing.T) {
		// Arrange
		configMock := mocks.MermerdConfig{}
		configMock.On("OmitAttributeKeys").Return(true).Once()
		configMock.On("ShowEnumValues").Return(true).Once()

		// Act
		result := getColumnData(&configMock, column)

		// Assert
		configMock.AssertExpectations(t)
		assert.Equal(t, columnName, result.Name)
		assert.Equal(t, enumValues, result.EnumValues)
		assert.Equal(t, none, result.AttributeKey)
	})

	t.Run("Get only minimal fields", func(t *testing.T) {
		// Arrange
		configMock := mocks.MermerdConfig{}
		configMock.On("OmitAttributeKeys").Return(true).Once()
		configMock.On("ShowEnumValues").Return(false).Once()

		// Act
		result := getColumnData(&configMock, column)

		// Assert
		configMock.AssertExpectations(t)
		assert.Equal(t, columnName, result.Name)
		assert.Equal(t, "", result.EnumValues)
		assert.Equal(t, none, result.AttributeKey)
	})
}

func TestShouldSkipConstraint(t *testing.T) {
	tableName1 := "Table1"
	tableName2 := "Table2"
	tables := []ErdTableData{{Name: tableName1}, {Name: tableName2}}

	t.Run("ShowAllConstraints config should never skip", func(t *testing.T) {
		// Arrange
		configMock := mocks.MermerdConfig{}
		configMock.On("ShowAllConstraints").Return(true).Once()
		constraint := database.ConstraintResult{PkTable: tableName1}

		// Act
		result := shouldSkipConstraint(&configMock, tables, constraint)

		// Assert
		configMock.AssertExpectations(t)
		assert.False(t, result)
	})

	t.Run("Skip constraint if both tables are not present", func(t *testing.T) {
		// Arrange
		configMock := mocks.MermerdConfig{}
		configMock.On("ShowAllConstraints").Return(false).Once()
		constraint := database.ConstraintResult{PkTable: tableName1, FkTable: "UnknownTable"}

		// Act
		result := shouldSkipConstraint(&configMock, tables, constraint)

		// Assert
		configMock.AssertExpectations(t)
		assert.True(t, result)
	})

	t.Run("Do not skip constraint if both tables are present", func(t *testing.T) {
		// Arrange
		configMock := mocks.MermerdConfig{}
		configMock.On("ShowAllConstraints").Return(false).Once()
		constraint := database.ConstraintResult{PkTable: tableName1, FkTable: tableName2}

		// Act
		result := shouldSkipConstraint(&configMock, tables, constraint)

		// Assert
		configMock.AssertExpectations(t)
		assert.False(t, result)
	})
}

func TestGetConstraintData(t *testing.T) {
	t.Run("OmitConstraintLabels should remove the constraint label", func(t *testing.T) {
		// Arrange
		configMock := mocks.MermerdConfig{}
		configMock.On("OmitConstraintLabels").Return(true).Once()
		constraint := database.ConstraintResult{ColumnName: "Column1"}

		// Act
		result := getConstraintData(&configMock, constraint)

		// Assert
		configMock.AssertExpectations(t)
		assert.Equal(t, result.ConstraintLabel, "")
	})
}

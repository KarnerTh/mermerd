package diagram

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/KarnerTh/mermerd/database"
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

package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstraintResultList_AppendIfNotExists(t *testing.T) {
	constraintItem := ConstraintResult{
		FkTable:        "tableA",
		PkTable:        "tableB",
		ConstraintName: "constraintXY",
		IsPrimary:      false,
		HasMultiplePK:  false,
	}
	constraintList := ConstraintResultList{constraintItem}

	t.Run("Same item should not be appended", func(t *testing.T) {
		// Arrange
		testItem := constraintItem

		// Act
		result := constraintList.AppendIfNotExists(testItem)

		// Assert
		expectedCount := len(constraintList)
		assert.Len(t, result, expectedCount)
	})

	t.Run("Different item should be appended", func(t *testing.T) {
		// Arrange
		testItem := constraintItem
		testItem.FkTable = "tableC"

		// Act
		result := constraintList.AppendIfNotExists(testItem)

		// Assert
		expectedCount := len(constraintList) + 1
		assert.Len(t, result, expectedCount)
	})
}

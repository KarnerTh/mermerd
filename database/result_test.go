package database

import "testing"

func TestConstraintResultList_AppendIfNotExists(t *testing.T) {
	constraintItem := ConstraintResult{
		FkTable:        "tableA",
		PKTable:        "tableB",
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
		if len(result) != expectedCount {
			t.Errorf("Expected %d items, but got %d", expectedCount, len(result))
		}
	})

	t.Run("Different item should be appended", func(t *testing.T) {
		// Arrange
		testItem := constraintItem
		testItem.FkTable = "tableC"

		// Act
		result := constraintList.AppendIfNotExists(testItem)

		// Assert
		expectedCount := len(constraintList) + 1
		if len(result) != expectedCount {
			t.Errorf("Expected %d items, but got %d", expectedCount, len(result))
		}
	})
}

package diagram

import (
	"fmt"
	"mermerd/database"
	"testing"
)

func TestGetRelation(t *testing.T) {
	testCases := []struct {
		isPrimary        bool
		hasMultiplePK    bool
		expectedRelation string
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
				PKTable:        "tableB",
				ConstraintName: "constraintXY",
				IsPrimary:      testCase.isPrimary,
				HasMultiplePK:  testCase.hasMultiplePK,
			}

			// Act
			result := getRelation(constraint)

			// Assert
			if result != testCase.expectedRelation {
				t.Errorf("Expected %s, got %s", testCase.expectedRelation, result)
			}
		})
	}
}

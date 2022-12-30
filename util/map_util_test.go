package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap2(t *testing.T) {
	// Arrange
	input := []int{1, 2, 3}

	// Act
	result := Map2(input, func(value int) string {
		return fmt.Sprintf("value: %d", value)
	})

	// Assert
	expectedResult := []string{"value: 1", "value: 2", "value: 3"}
	assert.ElementsMatch(t, expectedResult, result)
}

package database

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeValue(t *testing.T) {
	testCases := []struct {
		inputValue     string
		expectedResult string
	}{
		{inputValue: "timestamp without time zone", expectedResult: "timestamp_without_time_zone"},
		{inputValue: "numbers are allowed 7", expectedResult: "numbers_are_allowed_7"},
		{inputValue: "valid_stays_valid", expectedResult: "valid_stays_valid"},
		{inputValue: "valid-stays-valid", expectedResult: "valid-stays-valid"},
		{inputValue: "symbols_$_are_&_not_§_allowed", expectedResult: "symbols__are__not__allowed"},
		{inputValue: "also: not allowed", expectedResult: "also_not_allowed"},
	}

	for index, testCase := range testCases {
		t.Run(fmt.Sprintf("run #%d", index), func(t *testing.T) {
			// Arrange
			value := testCase.inputValue

			// Act
			result := SanitizeValue(value)

			// Assert
			assert.Equal(t, testCase.expectedResult, result)
		})
	}
}

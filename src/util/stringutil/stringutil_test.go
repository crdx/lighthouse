package stringutil_test

import (
	"fmt"
	"testing"

	"crdx.org/lighthouse/util/stringutil"
	"github.com/go-playground/assert/v2"
)

func TestPluralise(t *testing.T) {
	testCases := []struct {
		inputCount int
		inputUnit  string
		expected   string
	}{
		{1, "apple", "apple"},
		{2, "apple", "apples"},
		{0, "apple", "apples"},
		{1, "item", "item"},
		{2, "item", "items"},
		{-1, "apple", "apples"},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%d,%s", testCase.inputCount, testCase.inputUnit), func(t *testing.T) {
			actual := stringutil.Pluralise(testCase.inputCount, testCase.inputUnit)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

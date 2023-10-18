package timeutil_test

import (
	"fmt"
	"testing"

	"crdx.org/lighthouse/util/timeutil"
	"github.com/go-playground/assert/v2"
)

func TestTimeAgo(t *testing.T) {
	testCases := []struct {
		inputN         int
		inputLong      bool
		inputPrecision int
		expected       string
	}{
		{0, false, 0, "now"},
		{0, true, 0, "just now"},
		{60, false, 0, "1m ago"},
		{60, true, 0, "1 min ago"},
		{3600, false, 0, "1h ago"},
		{3660, false, 1, "1h ago"},
		{3660, true, 1, "1 hour ago"},
		{3660, true, 2, "1 hour 1 min ago"},
		{90000, false, 2, "1d 1h ago"},
		{90000, true, 2, "1 day 1 hour ago"},
		{1234567, true, 2, "2 weeks 6 hours ago"},
		{12345678, true, 1, "20 weeks ago"},
		{123456789, true, 2, "3 years 48 weeks ago"},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%d,%v,%d", testCase.inputN, testCase.inputLong, testCase.inputPrecision), func(t *testing.T) {
			actual := timeutil.TimeAgo(testCase.inputN, testCase.inputLong, testCase.inputPrecision)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

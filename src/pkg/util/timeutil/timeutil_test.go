package timeutil_test

import (
	"fmt"
	"testing"
	"time"

	"crdx.org/lighthouse/pkg/util/timeutil"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestToLocal(t *testing.T) {
	testTime := time.Date(2023, 10, 20, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		inputConfig *timeutil.Config
		inputTime   time.Time
		expected    time.Time
		shouldPanic bool
	}{
		{
			&timeutil.Config{Timezone: func() string { return "America/New_York" }},
			testTime,
			testTime.In(lo.Must(time.LoadLocation("America/New_York"))),
			false,
		},
		{
			&timeutil.Config{Timezone: func() string { return "Europe/London" }},
			testTime,
			testTime.In(lo.Must(time.LoadLocation("Europe/London"))),
			false,
		},
		{
			&timeutil.Config{Timezone: func() string { return "" }},
			testTime,
			time.Time{},
			true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.inputConfig.Timezone(), func(t *testing.T) {
			var actual time.Time
			var panicValue any

			timeutil.Init(testCase.inputConfig)

			func() {
				defer func() {
					if r := recover(); r != nil {
						panicValue = r
					}
				}()
				actual = timeutil.ToLocal(testCase.inputTime)
			}()

			if testCase.shouldPanic {
				assert.NotNil(t, panicValue)
			} else {
				assert.Equal(t, testCase.expected, actual)
			}
		})
	}
}

func TestTimeAgo(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			actual := timeutil.TimeAgo(testCase.inputN, testCase.inputLong, testCase.inputPrecision)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

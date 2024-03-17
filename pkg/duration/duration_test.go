package duration_test

import (
	"testing"
	"time"

	"crdx.org/lighthouse/pkg/duration"
	"github.com/stretchr/testify/assert"
)

func TestValid(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"0", true},
		{"", true},
		{"1 min", true},
		{"2 hours", true},
		{"3 days", true},
		{"4 weeks 1 day", true},
		{"1 min", true},
		{"1 hour 20 secs", true},
		{"1 day", true},
		{"1 week", true},
		{"60 secs", true},

		{"invalid", false},
		{"1 fortnight", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			actual := duration.Valid(testCase.input)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestParse(t *testing.T) {
	testCases := []struct {
		input    string
		expected time.Duration
		ok       bool
	}{
		{"", 0, true},
		{"60 secs", 60 * time.Second, true},
		{"1 min", 1 * time.Minute, true},
		{"2 hours", 2 * time.Hour, true},
		{"3 days 2 hours", 3*24*time.Hour + 2*time.Hour, true},
		{"1 min", 1 * time.Minute, true},
		{"1 hour", 1 * time.Hour, true},
		{"1 day", 24 * time.Hour, true},
		{"1 week", 7 * 24 * time.Hour, true},

		{"no weeks", 0, false},
		{"1 fortnight", 0, false},
		{"invalid", 0, false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			actual, ok := duration.Parse(testCase.input)
			assert.Equal(t, testCase.ok, ok)
			if testCase.ok {
				assert.Equal(t, testCase.expected, actual)
			}
		})
	}
}

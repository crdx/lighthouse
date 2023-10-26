package cache_test

import (
	"testing"
	"time"

	"crdx.org/lighthouse/pkg/cache"
	"github.com/stretchr/testify/assert"
)

func TestTemporalCache(t *testing.T) {
	t.Parallel()

	cache := cache.NewTemporal[string]()
	duration := 50 * time.Millisecond

	testCases := []struct {
		key               string
		expectedFirstRun  bool
		sleepDuration     time.Duration
		expectedSecondRun bool
	}{
		{"key1", false, 25 * time.Millisecond, true},
		{"key2", false, 75 * time.Millisecond, false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.key, func(t *testing.T) {
			t.Parallel()

			result := cache.SeenWithinLast(testCase.key, duration)
			assert.Equal(t, testCase.expectedFirstRun, result)

			time.Sleep(testCase.sleepDuration)

			result = cache.SeenWithinLast(testCase.key, duration)
			assert.Equal(t, testCase.expectedSecondRun, result)
		})
	}
}

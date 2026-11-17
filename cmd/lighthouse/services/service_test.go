package services

import (
	"testing"
	"time"
)

func TestExponentialBackoffDoublesUpToTheMaximum(t *testing.T) {
	for _, test := range []struct {
		previous time.Duration
		want     time.Duration
	}{
		{time.Second, 2 * time.Second},
		{20 * time.Minute, 40 * time.Minute},
		{40 * time.Minute, maximumRestartInterval},
		{maximumRestartInterval, maximumRestartInterval},
	} {
		if got := ExponentialBackoff(test.previous); got != test.want {
			t.Errorf("ExponentialBackoff(%s) = %s, want %s", test.previous, got, test.want)
		}
	}
}

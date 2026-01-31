package duration

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

// Valid returns true if s is a valid duration.
func Valid(s string) bool {
	_, ok := newDuration(s).seconds()
	return ok
}

// MustParse converts s into a time.Duration, and panics if it fails.
func MustParse(s string) time.Duration {
	return lo.Must(Parse(s))
}

// Parse converts s into a time.Duration.
func Parse(s string) (time.Duration, bool) {
	seconds, ok := newDuration(s).seconds()
	return time.Duration(seconds) * time.Second, ok //nolint:gosec // Human-readable durations will not overflow.
}

type duration struct {
	s string
}

func newDuration(s string) *duration {
	return &duration{s: s}
}

func (self *duration) tokens() []string {
	re := regexp.MustCompile(`[ ,]+`)
	return re.Split(strings.TrimSpace(self.s), -1)
}

func (self *duration) seconds() (uint, bool) {
	if self.s == "0" || self.s == "" {
		return 0, true
	}

	tokens := self.tokens()
	var total uint

	for i := 0; i < len(tokens); i += 2 {
		if i+1 >= len(tokens) {
			return 0, false
		}

		number, err := strconv.Atoi(tokens[i])
		if err != nil || number <= 0 {
			return 0, false
		}

		unit := tokens[i+1]
		seconds, ok := valueToSeconds(uint(number), unit)
		if !ok {
			return 0, false
		}

		total += seconds
	}

	return total, true
}

func valueToSeconds(value uint, unit string) (uint, bool) {
	const (
		secondMultiplier = 1
		minuteMultiplier = 60 * secondMultiplier
		hourMultiplier   = 60 * minuteMultiplier
		dayMultiplier    = 24 * hourMultiplier
		weekMultiplier   = 7 * dayMultiplier
	)

	switch unit {
	case "sec", "secs":
		return value * secondMultiplier, true
	case "min", "mins":
		return value * minuteMultiplier, true
	case "hour", "hours":
		return value * hourMultiplier, true
	case "day", "days":
		return value * dayMultiplier, true
	case "week", "weeks":
		return value * weekMultiplier, true
	default:
		return 0, false
	}
}

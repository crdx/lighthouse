package timeutil

import (
	"fmt"
	"strings"
	"time"

	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/util/stringutil"
	"github.com/samber/lo"
)

// ToLocal formats a time in the local timezone.
func ToLocal(t time.Time) time.Time {
	return t.In(lo.Must(time.LoadLocation(env.LocalTimeZone)))
}

// TimeAgo formats a number of seconds as a relative time in the past.
//
// Precision is the number of units to include. For example, for 65 seconds a precision of 1 would
// return "1 min" and a precision of 2 would return "1 min 5 secs".
func TimeAgo(seconds int, long bool, precision int) string {
	if seconds == 0 {
		if long {
			return "just now"
		} else {
			return "now"
		}
	}

	type Unit struct {
		longName  string
		shortName string
		seconds   int
	}

	units := []Unit{
		{"year", "y", 60 * 60 * 24 * 7 * 52},
		{"week", "w", 60 * 60 * 24 * 7},
		{"day", "d", 60 * 60 * 24},
		{"hour", "h", 60 * 60},
		{"min", "m", 60},
		{"sec", "s", 0},
	}

	var parts []string

	for _, unit := range units {
		if seconds < unit.seconds {
			continue
		}

		var partSeconds int
		if unit.seconds > 0 {
			partSeconds = seconds / unit.seconds
			seconds %= unit.seconds
		} else {
			partSeconds = seconds
		}

		if partSeconds > 0 {
			if long {
				parts = append(parts, fmt.Sprintf(
					"%d %s",
					partSeconds,
					stringutil.Pluralise(partSeconds, unit.longName),
				))
			} else {
				parts = append(parts, fmt.Sprintf(
					"%d%s",
					partSeconds,
					unit.shortName,
				))
			}

			if precision > 0 {
				precision--
				if precision == 0 {
					break
				}
			}
		}
	}

	return strings.Join(parts, " ") + " ago"
}

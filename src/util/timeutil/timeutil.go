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
func TimeAgo(n int, verbose bool, precision int) string {
	if n == 0 {
		if verbose {
			return "just now"
		} else {
			return "now"
		}
	}

	units := []struct {
		name  string
		value int
	}{
		{"year", 60 * 60 * 24 * 7 * 52},
		{"week", 60 * 60 * 24 * 7},
		{"day", 60 * 60 * 24},
		{"hour", 60 * 60},
		{"min", 60},
		{"sec", 0},
	}

	var a []string

	for _, unit := range units {
		if n < unit.value {
			continue
		}

		var x int
		if unit.value > 0 {
			x = n / unit.value
			n %= unit.value
		} else {
			x = n
		}

		if x > 0 {
			if verbose {
				a = append(a, fmt.Sprintf("%d %s", x, stringutil.Pluralise(x, unit.name)))
			} else {
				a = append(a, fmt.Sprintf("%d%s", x, string(unit.name[0])))
			}

			if precision > 0 {
				precision--
				if precision == 0 {
					break
				}
			}
		}
	}

	return strings.Join(a, " ") + " ago"
}

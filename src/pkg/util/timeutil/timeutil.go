package timeutil

import (
	"fmt"
	"strings"
	"time"

	"crdx.org/lighthouse/pkg/util/stringutil"
	"github.com/samber/lo"
)

type Config struct {
	Timezone func() string
}

var pkgConfig *Config

func Init(config *Config) {
	pkgConfig = config
}

// ToLocal converts a time to the local timezone.
func ToLocal(t time.Time) time.Time {
	if pkgConfig.Timezone() == "" {
		panic("no local timezone")
	}

	return t.In(lo.Must(time.LoadLocation(pkgConfig.Timezone())))
}

// FormatDuration formats a duration as a relative time.
//
// Precision is the number of units to include. For example, for 65 seconds a precision of 1 would
// return "1 min" and a precision of 2 would return "1 min 5 secs".
func FormatDuration(duration time.Duration, long bool, precision int, suffix string) string {
	seconds := int(duration.Seconds())

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
					stringutil.Pluralise(partSeconds, unit.longName, unit.longName+"s"),
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

	s := strings.Join(parts, " ")

	if suffix != "" {
		return s + " " + suffix
	}

	return s
}

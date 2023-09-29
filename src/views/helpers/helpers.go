package helpers

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"crdx.org/lighthouse/util"
)

func GetFuncMap() template.FuncMap {
	return template.FuncMap{
		"timeAgoVerbose": func(t time.Time) string {
			return timeAgo(int(time.Since(t).Seconds()), true, 1)
		},
		"timeAgo": func(t time.Time) string {
			return timeAgo(int(time.Since(t).Seconds()), false, 1)
		},
		"escape": escape,
	}
}

func escape(s string) string {
	return template.HTMLEscapeString(s)
}

func timeAgo(n int, verbose bool, precision int) string {
	if n == 0 {
		return "just now"
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
				a = append(a, fmt.Sprintf("%d %s", x, util.Pluralise(x, unit.name)))
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

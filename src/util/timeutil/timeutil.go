package timeutil

import (
	"time"

	"crdx.org/lighthouse/env"
	"github.com/samber/lo"
)

// ToLocal formats a time in the local timezone.
func ToLocal(t time.Time) time.Time {
	return t.In(lo.Must(time.LoadLocation(env.LocalTimeZone)))
}

package config

import (
	"html/template"
	"time"

	"crdx.org/lighthouse/pkg/constants"
	"crdx.org/lighthouse/pkg/env"
	"crdx.org/lighthouse/pkg/util/reflectutil"
	"crdx.org/lighthouse/pkg/util/stringutil"
	"crdx.org/lighthouse/pkg/util/timeutil"
)

func GetViewFuncMap() template.FuncMap {
	return template.FuncMap{
		"timeAgoLong": func(v any) string {
			if t, found := reflectutil.GetTime(v); found {
				return timeutil.FormatDuration(time.Since(t), true, 1, "ago")
			}
			return ""
		},
		"timeAgoShort": func(v any) string {
			if t, found := reflectutil.GetTime(v); found {
				return timeutil.FormatDuration(time.Since(t), false, 1, "ago")
			}
			return ""
		},

		"formatDurationLong": func(duration time.Duration) string {
			return timeutil.FormatDuration(duration, true, 1, "")
		},
		"formatDurationShort": func(duration time.Duration) string {
			return timeutil.FormatDuration(duration, false, 1, "")
		},

		"formatDateTimeSystem": func(v any) string {
			if t, found := reflectutil.GetTime(v); found {
				return timeutil.ToLocal(t).Format(constants.TimeFormatSystem)
			}
			return ""
		},
		"formatDateTimeReadable": func(v any) string {
			if t, found := reflectutil.GetTime(v); found {
				return timeutil.ToLocal(t).Format(constants.TimeFormatReadable)
			}
			return ""
		},

		"renderMarkdown":   func(s string) template.HTML { return template.HTML(stringutil.RenderMarkdown(s)) },
		"enableLiveReload": env.LiveReload,
		"isProduction":     env.Production,
		"disableAuth":      env.DisableAuth,
	}
}

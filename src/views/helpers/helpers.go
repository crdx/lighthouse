package helpers

import (
	"html/template"
	"time"

	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/util/reflectutil"
	"crdx.org/lighthouse/util/stringutil"
	"crdx.org/lighthouse/util/timeutil"
)

func GetFuncMap() template.FuncMap {
	return template.FuncMap{
		"timeAgoLong": func(v any) string {
			if t, found := reflectutil.GetTime(v); found {
				return timeutil.TimeAgo(int(time.Since(t).Seconds()), true, 1)
			}
			return ""
		},
		"timeAgoShort": func(v any) string {
			if t, found := reflectutil.GetTime(v); found {
				return timeutil.TimeAgo(int(time.Since(t).Seconds()), false, 1)
			}
			return ""
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

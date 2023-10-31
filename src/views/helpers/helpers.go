package helpers

import (
	"html/template"
	"time"

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
				return timeutil.ToLocal(t).Format("2006-01-02 15:04:05 MST")
			}
			return ""
		},
		"formatDateTimeReadable": func(v any) string {
			if t, found := reflectutil.GetTime(v); found {
				return timeutil.ToLocal(t).Format("15:04 on Mon, Jan _2 2006")
			}
			return ""
		},

		"renderMarkdown":   func(s string) template.HTML { return template.HTML(stringutil.RenderMarkdown(s)) },
		"enableLiveReload": env.LiveReload,
		"isProduction":     env.Production,
		"disableAuth":      env.DisableAuth,
	}
}

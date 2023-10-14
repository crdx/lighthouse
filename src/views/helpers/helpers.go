package helpers

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"

	"crdx.org/lighthouse/util/stringutil"
	"crdx.org/lighthouse/util/timeutil"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func GetFuncMap() template.FuncMap {
	return template.FuncMap{
		"timeAgoVerbose": func(t time.Time) string {
			return timeAgo(int(time.Since(t).Seconds()), true, 1)
		},
		"timeAgo": func(t time.Time) string {
			return timeAgo(int(time.Since(t).Seconds()), false, 1)
		},
		"formatDate": func(t time.Time) string {
			return timeutil.ToLocal(t).Format("02/01/2006")
		},
		"formatDateTime": func(t time.Time) string {
			return timeutil.ToLocal(t).Format("02/01/2006 15:04 MST")
		},
		"escape":         escape,
		"nl2br":          nl2br,
		"renderMarkdown": renderMarkdown,
	}
}

func renderMarkdown(s string) template.HTML {
	// Use the typographer extension, but disable some unwanted substitutions.
	// https://github.com/yuin/goldmark#typographerExtension-extension
	typographerExtension := extension.NewTypographer(
		extension.WithTypographicSubstitutions(extension.TypographicSubstitutions{
			extension.LeftSingleQuote:  nil,
			extension.RightSingleQuote: nil,
			extension.LeftDoubleQuote:  nil,
			extension.RightDoubleQuote: nil,
			extension.LeftAngleQuote:   nil,
			extension.RightAngleQuote:  nil,
		}),
	)

	markdownRenderer := goldmark.New(
		// https://github.com/yuin/goldmark#html-renderer-options
		goldmark.WithRendererOptions(html.WithHardWraps()),
		// https://github.com/yuin/goldmark#built-in-extensions
		goldmark.WithExtensions(extension.Linkify, typographerExtension),
	)

	var buf bytes.Buffer
	if err := markdownRenderer.Convert([]byte(s), &buf); err != nil {
		return template.HTML("An error occurred rendering this field's markdown")
	}

	return template.HTML(buf.String())
}

func nl2br(s string) template.HTML {
	s = template.HTMLEscapeString(s)
	s = strings.ReplaceAll(s, "\n", "<br>")
	return template.HTML(s)
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

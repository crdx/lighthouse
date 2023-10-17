package stringutil

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func Pluralise(count int, unit string) string {
	if count == 1 {
		return unit
	}
	return unit + "s"
}

func RenderMarkdown(s string) string {
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
		return "An error occurred rendering this field's markdown"
	}

	return buf.String()
}

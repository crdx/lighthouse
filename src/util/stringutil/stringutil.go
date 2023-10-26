package stringutil

import (
	"bytes"

	"github.com/samber/lo"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"golang.org/x/crypto/bcrypt"
)

// Pluralise pluralises a string based on count.
func Pluralise(count int, singular, plural string) string {
	if count == 1 {
		return singular
	} else {
		return plural
	}
}

// RenderMarkdown renders some markdown.
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
		return "An error occurred rendering this markdown"
	}

	return buf.String()
}

// Hash bcrypt hashes a password using a default cost.
func Hash(value string) string {
	return string(lo.Must(bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)))
}

// VerifyHashAndPassword verifies a bcrypt hash against a password.
func VerifyHashAndPassword(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

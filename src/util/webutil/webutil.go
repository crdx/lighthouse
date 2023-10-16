package webutil

import (
	"bytes"
	"net/url"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

// BuildURL builds a URL out of basePath and the contents of queryParams, or an error if basePath
// could not be parsed.
func BuildURL(basePath string, queryParams map[string]string) (string, error) {
	u, err := url.Parse(basePath)
	if err != nil {
		return "", err
	}

	values := url.Values{}

	for key, value := range queryParams {
		values.Add(key, value)
	}

	u.RawQuery = values.Encode()
	return u.String(), nil
}

// IsHTMLContentType returns true if the value of a Content-Type header is text/html.
func IsHTMLContentType(contentType string) bool {
	mimeType, _, _ := strings.Cut(contentType, ";")
	return strings.TrimSpace(mimeType) == "text/html"
}

// MinifyHTML minifies some HTML.
func MinifyHTML(s []byte) ([]byte, error) {
	htmlMinifier := &html.Minifier{}
	htmlMinifier.KeepComments = false            // Preserve all comments
	htmlMinifier.KeepConditionalComments = false // Preserve all IE conditional comments
	htmlMinifier.KeepDefaultAttrVals = false     // Preserve default attribute values
	htmlMinifier.KeepDocumentTags = false        // Preserve html, head and body tags
	htmlMinifier.KeepEndTags = false             // Preserve all end tags
	htmlMinifier.KeepWhitespace = false          // Preserve whitespace characters but still collapse multiple into one
	htmlMinifier.KeepQuotes = false              // Preserve quotes around attribute values

	var minifiedHTML bytes.Buffer
	err := htmlMinifier.Minify(minify.New(), &minifiedHTML, bytes.NewReader(s), nil)

	if err != nil {
		return nil, err
	} else {
		return minifiedHTML.Bytes(), nil
	}
}

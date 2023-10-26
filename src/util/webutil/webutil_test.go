package webutil_test

import (
	"fmt"
	"testing"

	"crdx.org/lighthouse/util/webutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		basePath    string
		queryParams map[string]string
		expected    string
		expectError bool
	}{
		{"https://example.com", map[string]string{"key1": "value1"}, "https://example.com?key1=value1", false},
		{"https://example.com", map[string]string{"key1": "value1", "key2": "value2"}, "https://example.com?key1=value1&key2=value2", false},
		{"https://example.com", map[string]string{}, "https://example.com", false},
		{"https://example.com", nil, "https://example.com", false},
		{":", map[string]string{"key1": "value1"}, "", true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.basePath, func(t *testing.T) {
			t.Parallel()

			actual, err := webutil.BuildURL(testCase.basePath, testCase.queryParams)

			if testCase.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestIsHTMLContentType(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		contentType string
		expected    bool
	}{
		{"text/html", true},
		{"text/html; charset=utf-8", true},
		{"charset=utf-8", false},
		{"  text/html  ", true},
		{"text/plain", false},
		{"application/json", false},
		{"application/json; charset=utf-8", false},
		{"", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.contentType, func(t *testing.T) {
			t.Parallel()

			actual := webutil.IsHTMLContentType(testCase.contentType)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestMinifyHTML(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input     []byte
		expected  []byte
		expectErr bool
	}{
		{[]byte("<html>  <body>\n   </body>  </html>"), []byte(nil), false},
		{[]byte("<div>  <p>Text</p>   </div>\n\n"), []byte("<div> <p>Text </div>"), false},
		{[]byte("<div>\n</div><!-- comment -->"), []byte("<div>\n</div>"), false},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Case%d", i+1), func(t *testing.T) {
			t.Parallel()

			actual, err := webutil.MinifyHTML(testCase.input)

			if testCase.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

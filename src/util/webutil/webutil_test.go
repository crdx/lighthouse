package webutil_test

import (
	"testing"

	"crdx.org/lighthouse/util/webutil"
	"github.com/stretchr/testify/assert"
)

func TestBuildURL(t *testing.T) {
	testCases := []struct {
		basePath    string
		queryParams map[string]string
		expected    string
		expectError bool
	}{
		{"http://example.com", map[string]string{"key1": "value1"}, "http://example.com?key1=value1", false},
		{"http://example.com", map[string]string{"key1": "value1", "key2": "value2"}, "http://example.com?key1=value1&key2=value2", false},
		{"http://example.com", map[string]string{}, "http://example.com", false},
		{"http://example.com", nil, "http://example.com", false},
		{":", map[string]string{"key1": "value1"}, "", true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.basePath, func(t *testing.T) {
			actual, err := webutil.BuildURL(testCase.basePath, testCase.queryParams)

			if testCase.expectError {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestIsHTMLContentType(t *testing.T) {
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
			actual := webutil.IsHTMLContentType(testCase.contentType)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

package stringutil_test

import (
	"fmt"
	"testing"

	"crdx.org/lighthouse/util/stringutil"
	"github.com/go-playground/assert/v2"
)

func TestPluralise(t *testing.T) {
	testCases := []struct {
		inputCount int
		inputUnit  string
		expected   string
	}{
		{1, "apple", "apple"},
		{2, "apple", "apples"},
		{0, "apple", "apples"},
		{1, "item", "item"},
		{2, "item", "items"},
		{-1, "apple", "apples"},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%d,%s", testCase.inputCount, testCase.inputUnit), func(t *testing.T) {
			actual := stringutil.Pluralise(testCase.inputCount, testCase.inputUnit)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestRenderMarkdown(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"**Bold**", "<p><strong>Bold</strong></p>\n"},
		{"**Bold*", "<p>*<em>Bold</em></p>\n"},
		{"*Italic*", "<p><em>Italic</em></p>\n"},
		{"[Link](https://example.com)", "<p><a href=\"https://example.com\">Link</a></p>\n"},
		{"An error occurred", "<p>An error occurred</p>\n"},
		{"", ""},
		{"<script>alert('1')</script>", "<!-- raw HTML omitted -->\n"},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Case%d", i+1), func(t *testing.T) {
			actual := stringutil.RenderMarkdown(testCase.input)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

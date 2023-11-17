package stringutil_test

import (
	"fmt"
	"testing"

	"crdx.org/lighthouse/pkg/util/stringutil"
	"github.com/stretchr/testify/assert"
)

func TestPluralise(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		inputCount    int
		inputSingular string
		inputPlural   string
		expected      string
	}{
		{1, "apple", "apples", "apple"},
		{2, "apple", "apples", "apples"},
		{0, "apple", "apples", "apples"},
		{1, "item", "items", "item"},
		{2, "item", "items", "items"},
		{-1, "apple", "apples", "apples"},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%d,%s", testCase.inputCount, testCase.inputSingular), func(t *testing.T) {
			t.Parallel()

			actual := stringutil.Pluralise(testCase.inputCount, testCase.inputSingular, testCase.inputPlural)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestRenderMarkdown(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input    string
		expected string
	}{
		{"**Bold**", "<p><strong>Bold</strong></p>\n"},
		{"**Bold*", "<p>*<em>Bold</em></p>\n"},
		{"*Italic*", "<p><em>Italic</em></p>\n"},
		{"[Link](https://example.com)", "<p><a href=\"https://example.com\">Link</a></p>\n"},
		{"Hello, world", "<p>Hello, world</p>\n"},
		{"", ""},
		{"<script>alert('1')</script>", "<!-- raw HTML omitted -->\n"},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Case%d", i+1), func(t *testing.T) {
			t.Parallel()

			actual := stringutil.RenderMarkdown(testCase.input)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestHashAndVerifyHash(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		inputPassword        string
		expectedVerifyResult bool
	}{
		{"hunter2", true},
		{"password", true},
		{"foo", true},
		{"", true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.inputPassword, func(t *testing.T) {
			t.Parallel()

			hashedPassword := stringutil.Hash(testCase.inputPassword)

			assert.NotEqual(t, testCase.inputPassword, hashedPassword)
			assert.True(t, stringutil.VerifyHashAndPassword(hashedPassword, testCase.inputPassword))
			assert.False(t, stringutil.VerifyHashAndPassword(hashedPassword, "incorrectPassword"))
		})
	}
}

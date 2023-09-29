package helpers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscape(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input    string
		expected string
	}{
		{`<`, `&lt;`},
		{`>`, `&gt;`},
		{`&`, `&amp;`},
		{`"`, `&#34;`},
		{`'`, `&#39;`},
		{`<script>alert('hello')</script>`, `&lt;script&gt;alert(&#39;hello&#39;)&lt;/script&gt;`},
		{`plain text`, `plain text`},
		{``, ``},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.input, func(t *testing.T) {
			t.Parallel()
			actual := escape(testCase.input)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestTimeAgo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		inputN         int
		inputVerbose   bool
		inputPrecision int
		expected       string
	}{
		{0, false, 0, "just now"},
		{60, false, 0, "1m ago"},
		{60, true, 0, "1 min ago"},
		{3600, false, 0, "1h ago"},
		{3660, false, 1, "1h ago"},
		{3660, true, 1, "1 hour ago"},
		{3660, true, 2, "1 hour 1 min ago"},
		{90000, false, 2, "1d 1h ago"},
		{90000, true, 2, "1 day 1 hour ago"},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(fmt.Sprintf("%d,%v,%d", testCase.inputN, testCase.inputVerbose, testCase.inputPrecision), func(t *testing.T) {
			t.Parallel()
			actual := timeAgo(testCase.inputN, testCase.inputVerbose, testCase.inputPrecision)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

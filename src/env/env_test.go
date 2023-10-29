package env

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
	r "github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input    string
		expected map[string]string
	}{
		{"A=1\nB=2", map[string]string{"A": "1", "B": "2"}},
		{"A=1\n# ignored\nB=2", map[string]string{"A": "1", "B": "2"}},
		{"A=1\n# ignored\nB=\"with quotes\"", map[string]string{"A": "1", "B": "with quotes"}},
		{"A=1\ninvalid value\nB=\"with quotes\"", map[string]string{"A": "1", "B": "with quotes"}},
		{"", map[string]string{}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			t.Parallel()

			actual := parse(testCase.input)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestRequire(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input func() string
		err   bool
	}{
		{func() string { return "foo" }, false}, //nolint
		{func() string { return "" }, true},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Case%d", i+1), func(t *testing.T) {
			t.Parallel()

			err := require(testCase.input, "")
			if testCase.err {
				r.Error(t, err)
			} else {
				r.NoError(t, err)
			}
		})
	}
}
func TestRequireIn(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		inputF          func() string
		inputValues     []string
		inputCanBeEmpty bool
		err             bool
	}{
		{func() string { return "foo" }, []string{"foo"}, false, false},
		{func() string { return "" }, []string{"foo"}, true, false},
		{func() string { return "foo" }, []string{"bar"}, true, true},
		{func() string { return "foo" }, []string{"bar"}, false, true},
		{func() string { return "" }, []string{"foo"}, false, true},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Case%d", i+1), func(t *testing.T) {
			t.Parallel()

			err := requireIn(testCase.inputF, "", testCase.inputValues, testCase.inputCanBeEmpty)
			if testCase.err {
				r.Error(t, err)
			} else {
				r.NoError(t, err)
			}
		})
	}
}

package transform_test

import (
	"fmt"
	"testing"

	"crdx.org/lighthouse/pkg/transform"
	"github.com/stretchr/testify/assert"
)

func TestStructTrim(t *testing.T) {
	t.Parallel()

	type S struct {
		Field1 string `transform:"trim"`
		Field2 string
	}

	testCases := []struct {
		input    S
		expected S
	}{
		{S{"  foo  ", "  foo "}, S{"foo", "  foo "}},
		{S{"   ", " "}, S{"", " "}},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s_%s", testCase.input.Field1, testCase.input.Field2), func(t *testing.T) {
			t.Parallel()

			transform.Struct(&testCase.input)
			assert.Equal(t, testCase.expected, testCase.input)
		})
	}
}

func TestStructTransform(t *testing.T) {
	t.Parallel()

	type S struct {
		Field1 string `transform:"upper"`
		Field2 string
	}

	testCases := []struct {
		input    S
		expected S
	}{
		{S{"foo", "foo"}, S{"FOO", "foo"}},
		{S{"foo   bar", "foo   bar"}, S{"FOO   BAR", "foo   bar"}},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s_%s", testCase.input.Field1, testCase.input.Field2), func(t *testing.T) {
			t.Parallel()

			transform.Struct(&testCase.input)
			assert.Equal(t, testCase.expected, testCase.input)
		})
	}
}

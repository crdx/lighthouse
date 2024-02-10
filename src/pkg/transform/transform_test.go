package transform_test

import (
	"fmt"
	"testing"

	"crdx.org/lighthouse/pkg/transform"
	"crdx.org/lighthouse/pkg/util/stringutil"
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

func TestPasswordFields(t *testing.T) {
	values := map[string]any{
		"password":         "hunter2",
		"confirm_password": "hunter2",
		"other_field":      "value",
	}

	transform.PasswordFields(values)

	assert.True(t, stringutil.VerifyHashAndPassword(values["password_hash"].(string), "hunter2"))
	assert.Nil(t, nil, values["password"])
	assert.Nil(t, nil, values["confirm_password"])
	assert.Equal(t, "value", values["other_field"])
}

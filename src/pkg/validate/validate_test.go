package validate_test

import (
	"fmt"
	"testing"

	"crdx.org/lighthouse/pkg/validate"
	"github.com/stretchr/testify/assert"
)

func TestStruct(t *testing.T) {
	t.Parallel()

	type S struct {
		Email    string `form:"email" validate:"required,email"`
		Password string `form:"password" validate:"required"`
	}

	testCases := []struct {
		input    S
		expected map[string]validate.Field
		err      bool
	}{
		{
			S{"", ""},
			map[string]validate.Field{
				"Email":    {Error: "required field", Value: "", Name: "email"},
				"Password": {Error: "required field", Value: "", Name: "password"},
			},
			true,
		},
		{
			S{"invalid", "password"},
			map[string]validate.Field{
				"Email":    {Error: "must be a valid email address", Value: "invalid", Name: "email"},
				"Password": {Error: "", Value: "password", Name: "password"},
			},
			true,
		},
		{
			S{"test@example.com", "password"},
			nil,
			false,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s_%s", testCase.input.Email, testCase.input.Password), func(t *testing.T) {
			t.Parallel()

			actual, err := validate.Struct(testCase.input)
			assert.Equal(t, testCase.err, err)

			if !err {
				assert.Equal(t, testCase.expected, actual)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	t.Parallel()

	type S struct {
		Field string `validate:"is-odd"`
	}

	validate.Register("is-odd", "must be an odd-length string", func(value string) bool {
		return len(value)%2 == 1
	})

	testCases := []struct {
		input    string
		expected string
		err      bool
	}{
		{"", "must be an odd-length string", true},
		{"even", "must be an odd-length string", true},
		{"odd", "", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			t.Parallel()

			fields, err := validate.Struct(S{Field: testCase.input})
			assert.Equal(t, testCase.err, err)
			assert.Equal(t, testCase.expected, fields["Field"].Error)
		})
	}
}

func TestFields(t *testing.T) {
	t.Parallel()

	type S struct {
		Field1 string `form:"field_1"`
		Field2 string `form:"field_2"`
	}

	expected := map[string]validate.Field{
		"Field1": {Name: "field_1"},
		"Field2": {Name: "field_2"},
	}

	actual := validate.Fields[S]()

	assert.Equal(t, expected, actual)
}

func TestMailAddrValidator(t *testing.T) {
	t.Parallel()

	type S struct {
		Field1 string `validate:"mailaddr"`
	}

	testCases := []struct {
		input    string
		expected string
		err      bool
	}{
		{"", `must be in the format "xxx <yyy>"`, true},
		{"John", `must be in the format "xxx <yyy>"`, true},
		{"John <john@example.com>", "", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			t.Parallel()

			testStruct := S{Field1: testCase.input}
			fields, err := validate.Struct(testStruct)
			assert.Equal(t, testCase.err, err)
			assert.Equal(t, testCase.expected, fields["Field1"].Error)
		})
	}
}

func TestTimezoneValidator(t *testing.T) {
	t.Parallel()

	type S struct {
		Field1 string `validate:"timezone"`
	}

	testCases := []struct {
		input     string
		expected  string
		expectErr bool
	}{
		{"foo", "must be valid", true},
		{"Invalid/Zone", "must be valid", true},
		{"UTC", "", false},
		{"Europe/London", "", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			t.Parallel()

			testStruct := S{Field1: testCase.input}
			fields, err := validate.Struct(testStruct)
			assert.Equal(t, testCase.expectErr, err)
			assert.Equal(t, testCase.expected, fields["Field1"].Error)
		})
	}
}

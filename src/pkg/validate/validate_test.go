package validate_test

import (
	"fmt"
	"testing"

	"crdx.org/lighthouse/pkg/validate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStruct(t *testing.T) {
	type S struct {
		Email    string `form:"email" validate:"required,email"`
		Password string `form:"password" validate:"required"`
	}

	testCases := []struct {
		input     S
		expected  map[string]validate.Field
		expectErr bool
	}{
		{
			S{"test@example.com", "password"},
			map[string]validate.Field{},
			false,
		},
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
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s_%s", testCase.input.Email, testCase.input.Password), func(t *testing.T) {
			actual, err := validate.Struct(testCase.input)
			if testCase.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestRegister(t *testing.T) {
	type S struct {
		Field string `validate:"is-odd"`
	}

	validate.Register("is-odd", "must be an odd-length string", func(value string) bool {
		return len(value)%2 == 1
	})

	testCases := []struct {
		input     string
		expected  string
		expectErr bool
	}{
		{"odd", "", false},
		{"", "must be an odd-length string", true},
		{"even", "must be an odd-length string", true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			fields, err := validate.Struct(S{Field: testCase.input})
			if testCase.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, testCase.expected, fields["Field"].Error)
		})
	}
}

func TestRegisterWithParam(t *testing.T) {
	type S struct {
		Field string `validate:"is=odd"`
	}

	validate.RegisterWithParam("is", "must be an {0} string", func(value string, param string) bool {
		return len(value)%2 == 1
	})

	testCases := []struct {
		input     string
		expected  string
		expectErr bool
	}{
		{"odd", "", false},
		{"", "must be an odd string", true},
		{"even", "must be an odd string", true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			fields, err := validate.Struct(S{Field: testCase.input})
			if testCase.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, testCase.expected, fields["Field"].Error)
		})
	}
}

func TestFields(t *testing.T) {
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
	type S struct {
		Field1 string `validate:"mailaddr"`
	}

	testCases := []struct {
		input     string
		expected  string
		expectErr bool
	}{
		{"John <john@example.com>", "", false},
		{"", `must be in the format "name <email>"`, true},
		{"John", `must be in the format "name <email>"`, true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			testStruct := S{Field1: testCase.input}
			fields, err := validate.Struct(testStruct)
			if testCase.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, testCase.expected, fields["Field1"].Error)
		})
	}
}

func TestTimezoneValidator(t *testing.T) {
	type S struct {
		Field1 string `validate:"timezone"`
	}

	testCases := []struct {
		input     string
		expected  string
		expectErr bool
	}{
		{"UTC", "", false},
		{"Europe/London", "", false},
		{"foo", "must be a valid timezone", true},
		{"Invalid/Zone", "must be a valid timezone", true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			testStruct := S{Field1: testCase.input}
			fields, err := validate.Struct(testStruct)
			if testCase.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, testCase.expected, fields["Field1"].Error)
		})
	}
}

func TestDurationValidator(t *testing.T) {
	type S struct {
		Field1 string `validate:"duration"`
	}

	testCases := []struct {
		input     string
		expected  string
		expectErr bool
	}{
		{"1 hour", "", false},
		{"invalid", "must be a valid duration", true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			testStruct := S{Field1: testCase.input}
			fields, err := validate.Struct(testStruct)
			if testCase.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, testCase.expected, fields["Field1"].Error)
		})
	}
}

func TestMaxDurationValidator(t *testing.T) {
	type S struct {
		Field1 string `validate:"dmax=1 hour"`
	}

	testCases := []struct {
		input     string
		expected  string
		expectErr bool
	}{
		{"1 hour", "", false},
		{"1 min", "", false},
		{"2 hours", "must be at most 1 hour", true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			testStruct := S{Field1: testCase.input}
			fields, err := validate.Struct(testStruct)
			if testCase.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, testCase.expected, fields["Field1"].Error)
		})
	}
}

func TestMinDurationValidator(t *testing.T) {
	type S struct {
		Field1 string `validate:"dmin=1 hour"`
	}

	testCases := []struct {
		input     string
		expected  string
		expectErr bool
	}{
		{"2 hours", "", false},
		{"1 hour", "", false},
		{"1 min", "must be at least 1 hour", true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			testStruct := S{Field1: testCase.input}
			fields, err := validate.Struct(testStruct)
			if testCase.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, testCase.expected, fields["Field1"].Error)
		})
	}
}

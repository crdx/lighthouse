package validate_test

import (
	"errors"
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
		input      S
		validators validate.ValidatorMap
		expected   map[string]validate.Field
		expectErr  bool
	}{
		{
			S{"test@example.com", "password"},
			validate.ValidatorMap{},
			map[string]validate.Field{},
			false,
		},
		{
			S{"test@example.com", "password"},
			validate.ValidatorMap{
				"Email": func(value string) error {
					return errors.New("invalid email")
				},
			},
			map[string]validate.Field{
				"Email":    {Error: "invalid email", Value: "test@example.com", Name: "email"},
				"Password": {Error: "", Value: "password", Name: "password"},
			},
			true,
		},
		{
			S{"", ""},
			validate.ValidatorMap{},
			map[string]validate.Field{
				"Email":    {Error: "required field", Value: "", Name: "email"},
				"Password": {Error: "required field", Value: "", Name: "password"},
			},
			true,
		},
		{
			S{"invalid", "password"},
			validate.ValidatorMap{},
			map[string]validate.Field{
				"Email":    {Error: "must be a valid email address", Value: "invalid", Name: "email"},
				"Password": {Error: "", Value: "password", Name: "password"},
			},
			true,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s_%s", testCase.input.Email, testCase.input.Password), func(t *testing.T) {
			actual, err := validate.Struct(testCase.input, testCase.validators)
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

func TestRoleValidator(t *testing.T) {
	type S struct {
		Field1 string `validate:"role"`
	}

	testCases := []struct {
		input     string
		expected  string
		expectErr bool
	}{
		{"1", "", false},
		{"2", "", false},
		{"3", "", false},
		{"0", "must be a valid role", true},
		{"4", "must be a valid role", true},
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

func TestIconValidator(t *testing.T) {
	type S struct {
		Field1 string `validate:"icon"`
	}

	testCases := []struct {
		input     string
		expected  string
		expectErr bool
	}{
		{"duotone:foo", "", false},
		{"solid:foo", "", false},
		{"brands:foo", "", false},
		{"foo:bar", "must be a valid icon", true},
		{"foo", "must be a valid icon", true},
		{"bar:foo", "must be a valid icon", true},
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

func TestIPAddressValidator(t *testing.T) {
	type S struct {
		Field1 string `validate:"ip_address"`
	}

	testCases := []struct {
		input     string
		expected  string
		expectErr bool
	}{
		{"127.0.0.1", "", false},
		{"10.0.0.1", "", false},
		{"82.83.84.85", "", false},
		{"foo", "must be a valid IPv4 address", true},
		{"1.2", "must be a valid IPv4 address", true},
		{"12345678", "must be a valid IPv4 address", true},
		{"::1", "must be a valid IPv4 address", true},
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

func TestMACAddressListValidator(t *testing.T) {
	type S struct {
		Field1 string `validate:"mac_address_list"`
	}

	testCases := []struct {
		input     string
		expected  string
		expectErr bool
	}{
		{"AA:AA:AA:AA:AA:AA", "", false},
		{"AA:AA:AA:AA:AA:AA, BB:BB:BB:BB:BB:BB", "", false},
		{"AA:AA:AA:AA:AA:AA, BB:BB:BB:BB:BB", "must be a valid list of MAC addresses", true},
		{"AA:AA:AA:AA:AA:AA; BB:BB:BB:BB:BB:BB", "must be a valid list of MAC addresses", true},
		{"AA:SS:AA:AA:AA:AA, BB:BB:BB:BB:BB:BB", "must be a valid list of MAC addresses", true},
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

func TestConfirmPassword(t *testing.T) {
	testCases := []struct {
		password        string
		confirmPassword string
		expectError     bool
	}{
		{"hunter2", "hunter2", false},
		{"hunter2", "Hunter2", true},
		{"foo123", "foo123", false},
		{"foo123", "foo1234", true},
		{"", "", false},
		{"password", "", true},
		{"", "password", true},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s_%s", testCase.password, testCase.confirmPassword), func(t *testing.T) {
			err := validate.ConfirmPassword(testCase.password)(testCase.confirmPassword)

			if testCase.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCurrentPassword(t *testing.T) {
	testCases := []struct {
		hash        string
		password    string
		expectError bool
	}{
		{`$2a$10$zqWx9ciqybX6FAARj0EDQOSn70FJ6RN.aHdAUNQFzucWnjXI3bnE6`, "password123", false},
		{`$2a$10$zqWx9ciqybX6FAARj0EDQOSn70FJ6RN.aHdAUNQFzucWnjXI3bnE6`, "Password123", true},
		{`$2a$10$2oaf1XqVTLS6iTFKPzZyWOA.ysfyUqptIIv18r8cP/AAIGivNo4su`, "abc123", false},
		{`$2a$10$2oaf1XqVTLS6iTFKPzZyWOA.ysfyUqptIIv18r8cP/AAIGivNo4su`, "abc1234", true},
		{"", "", true},
		{"foo", "", true},
		{"", "password", true},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s_%s", testCase.hash, testCase.password), func(t *testing.T) {
			err := validate.CurrentPassword(testCase.hash)(testCase.password)

			if testCase.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

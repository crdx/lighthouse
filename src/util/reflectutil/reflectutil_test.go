package reflectutil_test

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"crdx.org/lighthouse/util/reflectutil"
	"github.com/stretchr/testify/assert"
)

func TestStructToMap(t *testing.T) {
	type TestStruct struct {
		A string `json:"a"`
		B int    `json:"b"`
		C bool   `json:"c"`
	}

	testCases := []struct {
		inputStruct any
		inputTag    string
		expected    map[string]any
	}{
		{TestStruct{"one", 1, true}, "json", map[string]any{"a": "one", "b": 1, "c": true}},
		{TestStruct{"two", 2, false}, "json", map[string]any{"a": "two", "b": 2, "c": false}},
		{TestStruct{"", 0, false}, "json", map[string]any{"a": "", "b": 0, "c": false}},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%v,%s", testCase.inputStruct, testCase.inputTag), func(t *testing.T) {
			actual := reflectutil.StructToMap(testCase.inputStruct, testCase.inputTag)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestGetValue(t *testing.T) {
	type SimpleStruct struct {
		A int
	}

	testCases := []struct {
		input    any
		expected reflect.Value
	}{
		{123, reflect.ValueOf(123)},
		{"foo", reflect.ValueOf("foo")},
		{SimpleStruct{A: 1}, reflect.ValueOf(SimpleStruct{A: 1})},
		{&SimpleStruct{A: 2}, reflect.ValueOf(SimpleStruct{A: 2})},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%v", testCase.input), func(t *testing.T) {
			actual := reflectutil.GetValue(testCase.input)
			assert.Equal(t, testCase.expected.Interface(), actual.Interface())
			assert.Equal(t, testCase.expected.Kind(), actual.Kind())
		})
	}
}

func TestGetType(t *testing.T) {
	type SimpleStruct struct {
		A int
	}

	testCases := []struct {
		input    any
		expected reflect.Type
	}{
		{123, reflect.TypeOf(123)},
		{"foo", reflect.TypeOf("foo")},
		{SimpleStruct{A: 1}, reflect.TypeOf(SimpleStruct{})},
		{&SimpleStruct{A: 2}, reflect.TypeOf(SimpleStruct{})},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%v", testCase.input), func(t *testing.T) {
			actual := reflectutil.GetType(testCase.input)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestGetName(t *testing.T) {
	type SimpleStruct struct {
		A int
	}

	testCases := []struct {
		input    any
		expected string
	}{
		{123, "int"},
		{"foo", "string"},
		{SimpleStruct{A: 1}, "SimpleStruct"},
		{&SimpleStruct{A: 2}, "SimpleStruct"},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%v", testCase.input), func(t *testing.T) {
			actual := reflectutil.GetName(testCase.input)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestToString(t *testing.T) {
	testCases := []struct {
		input    reflect.Value
		expected string
	}{
		{reflect.ValueOf(123), "123"},
		{reflect.ValueOf(uint(123)), "123"},
		{reflect.ValueOf(123.456), "123.456"},
		{reflect.ValueOf("foo"), "foo"},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%v", testCase.input.Interface()), func(t *testing.T) {
			actual := reflectutil.ToString(testCase.input)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestGetTime(t *testing.T) {
	currentTime := time.Now()

	nullTime := sql.NullTime{
		Time:  currentTime,
		Valid: true,
	}

	invalidNullTime := sql.NullTime{
		Valid: false,
	}

	testCases := []struct {
		input       interface{}
		expected    time.Time
		expectValid bool
	}{
		{currentTime, currentTime, true},
		{nullTime, currentTime, true},
		{invalidNullTime, time.Time{}, false},
		{nil, time.Time{}, false},
		{123, time.Time{}, false},
		{"invalid", time.Time{}, false},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%v", testCase.input), func(t *testing.T) {
			actual, valid := reflectutil.GetTime(testCase.input)
			assert.Equal(t, testCase.expected, actual)
			assert.Equal(t, testCase.expectValid, valid)
		})
	}
}

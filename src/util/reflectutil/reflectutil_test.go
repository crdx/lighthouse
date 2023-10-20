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
		expected    map[string]any
	}{
		{TestStruct{"one", 1, true}, map[string]any{"a": "one", "b": 1, "c": true}},
		{TestStruct{"two", 2, false}, map[string]any{"a": "two", "b": 2, "c": false}},
		{TestStruct{"", 0, false}, map[string]any{"a": "", "b": 0, "c": false}},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%v", testCase.inputStruct), func(t *testing.T) {
			actual := reflectutil.StructToMap(testCase.inputStruct, "json")
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestMapToStruct(t *testing.T) {
	type Person struct {
		Name     string `json:"name"`
		Age      int    `json:"age"`
		Employed bool   `json:"employed"`
	}

	testCases := []struct {
		inputMap map[string]string
		expected Person
	}{
		{
			map[string]string{"name": "Alice", "age": "30", "employed": "1"},
			Person{Name: "Alice", Age: 30, Employed: true},
		},
		{
			map[string]string{"name": "Bob", "age": "40", "employed": "0"},
			Person{Name: "Bob", Age: 40, Employed: false},
		},
		{
			map[string]string{"name": "Charlie", "age": "50"},
			Person{Name: "Charlie", Age: 50, Employed: false},
		},
		{
			map[string]string{},
			Person{},
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%v", testCase.inputMap), func(t *testing.T) {
			actual := reflectutil.MapToStruct[Person](testCase.inputMap, "json")
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestGetStructFieldValue(t *testing.T) {
	type S struct {
		FieldA string `tagA:"FieldA"`
		FieldB int    `tagA:"FieldB"`
		FieldC bool   `tagB:"FieldC"`
	}

	testCases := []struct {
		inputStruct   S
		inputTagValue string
		inputTagName  string
		expected      any
		found         bool
	}{
		{
			S{"valueA", 42, true},
			"FieldA",
			"tagA",
			"valueA",
			true,
		},
		{
			S{"valueA", 42, true},
			"FieldB",
			"tagA",
			42,
			true,
		},
		{
			S{"valueA", 42, true},
			"FieldC",
			"tagB",
			true,
			true,
		},
		{
			S{"valueA", 42, true},
			"InvalidField",
			"tagA",
			nil,
			false,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%v,%s,%s", testCase.inputStruct, testCase.inputTagValue, testCase.inputTagName), func(t *testing.T) {
			value, found := reflectutil.GetStructFieldValue(testCase.inputStruct, testCase.inputTagValue, testCase.inputTagName)
			if found {
				assert.True(t, found)
				assert.Equal(t, testCase.expected, value.Interface())
			} else {
				assert.False(t, found)
				assert.False(t, value.IsValid())
			}
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

func TestToString(t *testing.T) {
	testCases := []struct {
		input    reflect.Value
		expected string
	}{
		{reflect.ValueOf(123), "123"},
		{reflect.ValueOf(uint(123)), "123"},
		{reflect.ValueOf(123.456), "123.456"},
		{reflect.ValueOf("foo"), "foo"},
		{reflect.ValueOf(true), "1"},
		{reflect.ValueOf(false), "0"},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%v", testCase.input.Interface()), func(t *testing.T) {
			actual := reflectutil.ToString(testCase.input.Interface())
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
		input       any
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

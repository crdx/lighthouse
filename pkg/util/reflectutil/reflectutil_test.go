package reflectutil_test

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"crdx.org/lighthouse/pkg/util/reflectutil"
	"github.com/stretchr/testify/assert"
)

func TestStructToMap(t *testing.T) {
	t.Parallel()

	type S struct {
		A string `json:"a"`
		B int    `json:"b"`
		C bool   `json:"c"`
	}

	testCases := []struct {
		inputStruct any
		expected    map[string]any
	}{
		{S{"one", 1, true}, map[string]any{"a": "one", "b": 1, "c": true}},
		{S{"two", 2, false}, map[string]any{"a": "two", "b": 2, "c": false}},
		{S{"", 0, false}, map[string]any{"a": "", "b": 0, "c": false}},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%v", testCase.inputStruct), func(t *testing.T) {
			t.Parallel()

			actual := reflectutil.StructToMap(testCase.inputStruct, "json")
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestMapToStruct(t *testing.T) {
	t.Parallel()

	type S struct {
		Field1 uint    `json:"field1"`
		Field2 float64 `json:"field2"`
		Field3 string  `json:"field3"`
		Field4 int     `json:"field4"`
		Field5 bool    `json:"field5"`
	}

	testCases := []struct {
		inputMap map[string]string
		expected S
	}{
		{
			map[string]string{"field1": "12345", "field2": "3.14", "field3": "bd949968-d6eb-4cb6-8c5f-1e1bf3b380ad", "field4": "30", "field5": "1"},
			S{Field3: "bd949968-d6eb-4cb6-8c5f-1e1bf3b380ad", Field4: 30, Field5: true, Field1: 12345, Field2: 3.14},
		},
		{
			map[string]string{"field3": "ced99b9d-e3a9-4bce-9e12-a8df06d41a8d", "field4": "40", "field5": "0"},
			S{Field3: "ced99b9d-e3a9-4bce-9e12-a8df06d41a8d", Field4: 40, Field5: false},
		},
		{
			map[string]string{"field3": "45e7b0fe-c6b8-4b69-80ca-eeac29cab851", "field4": "50"},
			S{Field3: "45e7b0fe-c6b8-4b69-80ca-eeac29cab851", Field4: 50, Field5: false},
		},
		{
			map[string]string{"field3": "45e7b0fe-c6b8-4b69-80ca-eeac29cab851", "field4": "50", "foo": "bar"},
			S{Field3: "45e7b0fe-c6b8-4b69-80ca-eeac29cab851", Field4: 50, Field5: false},
		},
		{
			map[string]string{},
			S{},
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%v", testCase.inputMap), func(t *testing.T) {
			t.Parallel()

			actual := reflectutil.MapToStruct[S](testCase.inputMap, "json")
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestGetStructFieldValue(t *testing.T) {
	t.Parallel()

	type S struct {
		Field1 string `tagA:"Field1"`
		Field2 int    `tagA:"Field2"`
		Field3 bool   `tagB:"Field3"`
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
			"Field1",
			"tagA",
			"valueA",
			true,
		},
		{
			S{"valueA", 42, true},
			"Field2",
			"tagA",
			42,
			true,
		},
		{
			S{"valueA", 42, true},
			"Field3",
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
			t.Parallel()

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
	t.Parallel()

	type S struct {
		Field int
	}

	testCases := []struct {
		input    any
		expected reflect.Value
	}{
		{123, reflect.ValueOf(123)},
		{"foo", reflect.ValueOf("foo")},
		{S{Field: 1}, reflect.ValueOf(S{Field: 1})},
		{&S{Field: 2}, reflect.ValueOf(S{Field: 2})},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%v", testCase.input), func(t *testing.T) {
			t.Parallel()

			actual := reflectutil.GetValue(testCase.input)
			assert.Equal(t, testCase.expected.Interface(), actual.Interface())
			assert.Equal(t, testCase.expected.Kind(), actual.Kind())
		})
	}
}

func TestGetType(t *testing.T) {
	t.Parallel()

	type S struct {
		Field int
	}

	testCases := []struct {
		input    any
		expected reflect.Type
	}{
		{123, reflect.TypeFor[int]()},
		{"foo", reflect.TypeFor[string]()},
		{S{Field: 1}, reflect.TypeFor[S]()},
		{&S{Field: 2}, reflect.TypeFor[S]()},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%v", testCase.input), func(t *testing.T) {
			t.Parallel()

			actual := reflectutil.GetType(testCase.input)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestToString(t *testing.T) {
	t.Parallel()

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
		{reflect.ValueOf(struct{}{}), "{}"},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%v", testCase.input.Interface()), func(t *testing.T) {
			t.Parallel()

			actual := reflectutil.ToString(testCase.input.Interface())
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestGetTime(t *testing.T) {
	t.Parallel()

	currentTime := time.Now()

	nullTime := sql.Null[time.Time]{
		V:     currentTime,
		Valid: true,
	}

	invalidNullTime := sql.Null[time.Time]{
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
			t.Parallel()

			actual, valid := reflectutil.GetTime(testCase.input)
			assert.Equal(t, testCase.expected, actual)
			assert.Equal(t, testCase.expectValid, valid)
		})
	}
}

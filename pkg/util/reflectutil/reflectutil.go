package reflectutil

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/samber/lo"
)

// StructToMap converts struct T into a map[string]any using values from tag tagName as the keys.
func StructToMap[T any](s T, tagName string) map[string]any {
	structValue := GetValue(s)
	fields := map[string]any{}

	for i := 0; i < structValue.NumField(); i++ {
		tagValue := structValue.Type().Field(i).Tag.Get(tagName)

		fields[tagValue] = structValue.Field(i).Interface()
	}

	return fields
}

// MapToStruct converts a map[string]string into a struct T using values from tag tagName as the
// keys.
func MapToStruct[T any](m map[string]string, tagName string) T {
	var s T

	for name, value := range m {
		fieldValue, ok := GetStructFieldValue(&s, name, tagName)
		if !ok {
			continue
		}

		switch fieldValue.Kind() {
		case reflect.Int:
			fieldValue.SetInt(lo.Must(strconv.ParseInt(value, 10, 64)))
		case reflect.Uint:
			fieldValue.SetUint(lo.Must(strconv.ParseUint(value, 10, 64)))
		case reflect.Float64:
			fieldValue.SetFloat(lo.Must(strconv.ParseFloat(value, 64)))
		case reflect.Bool:
			fieldValue.SetBool(value == "1")
		default:
			fieldValue.SetString(value)
		}
	}

	return s
}

// GetStructFieldValue returns the reflect.Value corresponding to the field of struct T whose field
// has a tag tagName with value tagValue.
//
// For example, given this struct and var declaration:
//
//	type S struct {
//		Foo string `slug:"foo"
//		Bar string `slug:"bar"
//	}
//	var s S
//
// GetStructFieldValue(&s, "foo", "slug") will return the reflect.Value corresponding to s.Foo.
func GetStructFieldValue[T any](s T, tagValue string, tagName string) (reflect.Value, bool) {
	structValue := GetValue(s)

	for i := 0; i < structValue.NumField(); i++ {
		actualValue := structValue.Type().Field(i).Tag.Get(tagName)

		if actualValue == tagValue {
			return structValue.Field(i), true
		}
	}

	return reflect.Value{}, false
}

// GetValue gets the value of s, following pointers.
func GetValue(s any) reflect.Value {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

// GetValue gets the type of s, following pointers.
func GetType(s any) reflect.Type {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// ToString returns value converted into a string.
func ToString(value any) string {
	switch value := value.(type) {
	case int:
		return strconv.FormatInt(int64(value), 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case float64:
		return strconv.FormatFloat(value, 'f', 3, 64)
	case bool:
		if value {
			return "1"
		} else {
			return "0"
		}
	case string:
		return value
	default:
		return fmt.Sprintf("%s", value)
	}
}

// GetTime returns the time.Time that corresponds to the passed value. The value must be either
// time.Time or sql.Null[time.Time].
func GetTime(v any) (time.Time, bool) {
	switch t := v.(type) {
	case sql.Null[time.Time]:
		if t.Valid {
			return t.V, true
		}
	case time.Time:
		return t, true
	}
	return time.Time{}, false
}

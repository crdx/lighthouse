package reflectutil

import (
	"fmt"
	"reflect"
	"strconv"
)

// StructToMap converts struct s into a map[string]any with keys taken from the value of the
// provided tag.
func StructToMap(s any, tag string) map[string]any {
	value := GetValue(s)
	fields := map[string]any{}

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		tag := value.Type().Field(i).Tag.Get(tag)
		fields[tag] = field.Interface()
	}

	return fields
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

// GetName returns the name of s.
//
// For example, for a struct named Foo, it returns "Foo".
func GetName(s any) string {
	return GetType(s).Name()
}

// ToString returns a reflect.Value converted into a string.
func ToString(value reflect.Value) string {
	switch value := value.Interface().(type) {
	case int:
		return strconv.Itoa(value)
	case float64:
		return strconv.FormatFloat(value, 'f', 3, 64)
	case string:
		return value
	default:
		return fmt.Sprintf("%s", value)
	}
}

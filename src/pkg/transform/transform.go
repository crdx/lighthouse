package transform

import (
	"strings"

	"crdx.org/lighthouse/util/reflectutil"
)

// Struct transforms a struct's contents according to the rules set in the "transform" tag.
func Struct[T any](s T) {
	structValue := reflectutil.GetValue(s)

	for i := 0; i < structValue.NumField(); i++ {
		fieldValue := structValue.Field(i)

		if str, ok := fieldValue.Interface().(string); ok {
			tagValue := structValue.Type().Field(i).Tag.Get("transform")
			if tagValue == "trim" {
				fieldValue.SetString(strings.TrimSpace(str))
			}
		}
	}
}

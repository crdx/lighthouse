package transform

import (
	"strings"

	"crdx.org/lighthouse/util/reflectutil"
)

// Struct transforms a struct's contents according to the rules set in the "transform" tag.
func Struct[T any](s T) {
	value := reflectutil.GetValue(s)
	type_ := reflectutil.GetType(s)

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)

		if str, ok := field.Interface().(string); ok {
			tag := type_.Field(i).Tag.Get("transform")
			if tag == "trim" {
				field.SetString(strings.TrimSpace(str))
			}
		}
	}
}

package transform

import (
	"reflect"
	"strings"

	"crdx.org/lighthouse/pkg/util/reflectutil"
)

// Struct transforms a struct's contents according to the rules set in the "transform" tag.
func Struct[T any](s T) {
	structValue := reflectutil.GetValue(s)

	for i := range structValue.NumField() {
		fieldValue := structValue.Field(i)

		if str, ok := reflect.TypeAssert[string](fieldValue); ok {
			tagValue := structValue.Type().Field(i).Tag.Get("transform")

			noTrim := false
			for transformation := range strings.SplitSeq(tagValue, ",") {
				if transformation == "no-trim" {
					noTrim = true
				}

				if transformation == "upper" {
					str = strings.ToUpper(str)
				}
			}

			if !noTrim {
				str = strings.TrimSpace(str)
			}
			fieldValue.SetString(str)
		}
	}
}

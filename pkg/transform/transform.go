package transform

import (
	"strings"

	"crdx.org/lighthouse/pkg/util/reflectutil"
)

// Struct transforms a struct's contents according to the rules set in the "transform" tag.
func Struct[T any](s T) {
	structValue := reflectutil.GetValue(s)

	for i := range structValue.NumField() {
		fieldValue := structValue.Field(i)

		if str, ok := fieldValue.Interface().(string); ok {
			tagValue := structValue.Type().Field(i).Tag.Get("transform")

			noTrim := false
			for _, transformation := range strings.Split(tagValue, ",") {
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

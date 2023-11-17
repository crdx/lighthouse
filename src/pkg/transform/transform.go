package transform

import (
	"strings"

	"crdx.org/lighthouse/pkg/util/reflectutil"
	"crdx.org/lighthouse/pkg/util/stringutil"
)

// Struct transforms a struct's contents according to the rules set in the "transform" tag.
func Struct[T any](s T) {
	structValue := reflectutil.GetValue(s)

	for i := 0; i < structValue.NumField(); i++ {
		fieldValue := structValue.Field(i)

		if str, ok := fieldValue.Interface().(string); ok {
			tagValue := structValue.Type().Field(i).Tag.Get("transform")
			for _, transformation := range strings.Split(tagValue, ",") {
				if transformation == "trim" {
					str = strings.TrimSpace(str)
				}

				fieldValue.SetString(str)
			}
		}
	}
}

func PasswordFields(values map[string]any) {
	if password := values["password"].(string); password != "" {
		values["password_hash"] = stringutil.Hash(password)
	}

	delete(values, "password")
	delete(values, "confirm_password")
}

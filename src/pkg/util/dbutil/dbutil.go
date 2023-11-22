package dbutil

import (
	"reflect"

	"github.com/samber/lo"
)

// MapByID maps a slice of *T to a map[uint]*T using the ID field as the key. If T is not a struct,
// or the ID field does not exist in T, this function will panic.
func MapByID[T any](items []*T) map[uint]*T {
	return MapBy[uint]("ID", items)
}

// MapBy maps a slice of *T to a map[A]*T using keyField as the key. If T is not a struct, or
// keyField does not exist in T, this function will panic.
func MapBy[A comparable, T any](keyField string, items []*T) map[A]*T {
	return lo.SliceToMap(items, func(item *T) (A, *T) {
		return reflect.ValueOf(item).Elem().FieldByName(keyField).Interface().(A), item
	})
}

// MapBy2 maps a slice of *T to a map[A]B using keyField and valueField as the key and value. If T
// is not a struct, or keyField or valueField do not exist in T, this function will panic.
func MapBy2[A comparable, B any, T any](keyField, valueField string, items []*T) map[A]B {
	return lo.SliceToMap(items, func(item *T) (A, B) {
		v := reflect.ValueOf(item).Elem()
		return v.FieldByName(keyField).Interface().(A), v.FieldByName(valueField).Interface().(B)
	})
}

package internal

import (
	"reflect"
)

// RemainingKeys will inspect a struct and compare it to a map. Any key that
// is defined in a JSON tag of the struct will be removed from the map. The
// remaining map and keys are then returned.
//
// This is useful for determining the extra fields returned in response bodies
// for resources that can contain an arbitrary or dynamic number of fields.
func RemainingKeys(s interface{}, m map[string]interface{}) map[string]interface{} {
	valueOf := reflect.ValueOf(s)
	typeOf := reflect.TypeOf(s)
	for i := 0; i < valueOf.NumField(); i++ {
		field := typeOf.Field(i)
		tagValue := field.Tag.Get("json")
		if _, ok := m[tagValue]; ok {
			delete(m, tagValue)
		}
	}

	return m
}

package maptools

import (
	"reflect"
)

// convert a struct to a map
// c: if it is a struct pointer, you must use * to get the actual value of the memory area
func StructToMap(c interface{}) map[string]interface{} {
	ref := reflect.ValueOf(c)
	t := ref.Type()

	result := make(map[string]interface{})

	for i := 0; i < t.NumField(); i++ {
		n := t.Field(i).Name
		v := ref.FieldByName(n)
		result[n] = v.Interface()
	}
	return result
}

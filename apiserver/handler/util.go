package handler

import "reflect"

func IsEmpty(v any) bool {
	if v == nil {
		return true
	}

	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return val.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return val.IsNil()
	case reflect.Bool:
		return !val.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return val.Float() == 0
	case reflect.Struct:
		// For struct: check if all fields are zero
		return reflect.DeepEqual(v, reflect.Zero(val.Type()).Interface())
	default:
		return false
	}
}

func StructToMap(obj any, tag string) map[string]any {
	if tag == "" {
		tag = "db"
	}
	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}
	m := make(map[string]any)
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get(tag)
		if tag == "" || tag == "-" {
			continue
		}
		m[tag] = val.Field(i).Interface()
		// if !IsEmpty(val.Field(i).Interface()) {
		// 	m[tag] = val.Field(i).Interface()
		// }
	}
	return m
}

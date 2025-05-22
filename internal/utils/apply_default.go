// Package utils defines some internal utilities
package utils

import (
	"errors"
	"reflect"
	"strconv"
)

// ApplyDefaults sets default values based on `default` tag for nil pointer fields.
func ApplyDefaults(obj any) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("obj should be a pointer and not a nil pointer")
	}

	v = v.Elem()
	t := v.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		fieldVal := v.Field(i)

		// Get the default tag
		defaultTag := field.Tag.Get("default")
		if defaultTag == "" || fieldVal.Kind() != reflect.Ptr || !fieldVal.IsNil() {
			continue
		}

		switch field.Type.Elem().Kind() {
		case reflect.Bool:
			val, err := strconv.ParseBool(defaultTag)
			if err == nil {
				fieldVal.Set(reflect.ValueOf(&val))
			}
		case reflect.Int:
			val, err := strconv.Atoi(defaultTag)
			if err == nil {
				fieldVal.Set(reflect.ValueOf(&val))
			}
		case reflect.Int8:
			val, err := ParseInt8(defaultTag)
			if err == nil {
				fieldVal.Set(reflect.ValueOf(&val))
			}
		case reflect.Uint32:
			val, err := ParseUInt32(defaultTag)
			if err == nil {
				fieldVal.Set(reflect.ValueOf(&val))
			}
		case reflect.String:
			fieldVal.Set(reflect.ValueOf(&defaultTag))
		}
	}
	return nil
}

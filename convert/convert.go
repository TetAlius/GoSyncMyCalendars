package convert

import (
	"fmt"
	"reflect"
	"strings"
)

type Converter interface {
	// Method use to convert from an interface{} to a struct
	Convert(interface{}, string, string) (Converter, error)
}

// parseTag splits a struct field's tag into its name and
// comma-separated options.
func parseTag(tag string) (string, string) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], string(tag[idx+1:])
	}
	return tag, string("")
}

// Function that starts a conversion between to different models
func Convert(from interface{}, to interface{}) error {
	v := reflect.ValueOf(to)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return fmt.Errorf("nil struct sended")
	}

	m := deconvert(from)
	return conversion(v, m)
}

// function that loops in all attributes and sets the deconverted info
func conversion(val reflect.Value, from map[string]interface{}) (err error) {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag, opts := parseTag(field.Tag.Get("convert"))
		rv := reflect.ValueOf(field)
		v := field.Type
		if tag == "" || tag == "-" {
			continue
		}
		m, ok := val.Field(i).Interface().(Converter)
		if ok {
			f, err := m.Convert(from[tag], tag, opts)
			if err != nil {
				return err
			}
			val.Field(i).Set(reflect.ValueOf(f))
			continue
		}
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
			if err = conversion(rv, from); err != nil {
				return err
			}
			continue
		}
		if v.Kind() == reflect.Struct {
			if err = conversion(rv, from); err != nil {
				return err
			}
			continue
		}

		val.Field(i).Set(reflect.ValueOf(from[tag]))
	}
	return nil
}

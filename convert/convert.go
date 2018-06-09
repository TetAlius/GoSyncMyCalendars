package convert

import (
	"fmt"
	"reflect"
	"strings"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

type Converter interface {
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

func Convert(from interface{}, to interface{}) error {
	log.Debugln("Converting...")

	v := reflect.ValueOf(to)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return fmt.Errorf("nil struct sended")
	}

	m := deconvert(from)
	log.Debugf("%s", m)
	return conversion(v, m)
}

func conversion(val reflect.Value, from map[string]interface{}) (err error) {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag, opts := parseTag(field.Tag.Get("convert"))
		if tag == "" || tag == "-" {
			continue
		}
		rv := reflect.ValueOf(field)
		v := field.Type
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

//
//func structToMap(i interface{}) (values map[string]interface{}) {
//	values = make(map[string]interface{})
//	iVal := reflect.ValueOf(i)
//	if iVal.Kind() == reflect.Ptr {
//		iVal = iVal.Elem()
//	}
//	typ := iVal.Type()
//	for i := 0; i < iVal.NumField(); i++ {
//		f := iVal.Field(i)
//		tag, _ := parseTag(typ.Field(i).Tag.Get("convert"))
//		if tag == "" || tag == "-" {
//			continue
//		}
//		if f.Kind() == reflect.Ptr && f.IsNil() {
//			log.Debugf("nil: %s", tag)
//			continue
//		}
//
//		if f.Kind() == reflect.Ptr {
//			values[tag] = structToMap(f.Interface())
//			continue
//		}
//		if f.Kind() == reflect.Struct {
//			values[tag] = structToMap(f.Interface())
//			continue
//		}
//
//		var v interface{}
//		switch f.Interface().(type) {
//		case int, int8, int16, int32, int64:
//			v = int(f.Int())
//		case uint, uint8, uint16, uint32, uint64:
//			v = uint(f.Uint())
//		case float32, float64:
//			v = f.Float()
//		case []byte:
//			v = f.Bytes()
//		case string, time.Time, time.Location, *time.Time, *time.Location:
//			v = f.String()
//		case bool:
//			v = f.Bool()
//		}
//		values[tag] = v
//	}
//	return
//}

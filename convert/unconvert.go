package convert

import (
	"reflect"
	"time"
)

type Deconverter interface {
	// Method use to deconvert as struct to an interface{}
	Deconvert() interface{}
}

// Function that deconverts a model to an interface{}
func deconvert(i interface{}) (values map[string]interface{}) {
	values = make(map[string]interface{})
	iVal := reflect.ValueOf(i)
	if iVal.Kind() == reflect.Ptr {
		iVal = iVal.Elem()
	}
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		f := iVal.Field(i)
		tag, _ := parseTag(typ.Field(i).Tag.Get("convert"))
		if f.Kind() == reflect.Ptr && f.IsNil() {
			continue
		}
		if tag == "" || tag == "-" {
			continue
		}
		m, ok := iVal.Field(i).Interface().(Deconverter)
		if ok {
			values[tag] = m.Deconvert()
			continue
		}

		if f.Kind() == reflect.Ptr {
			values[tag] = deconvert(f.Interface())
			continue
		}
		if f.Kind() == reflect.Struct {
			values[tag] = deconvert(f.Interface())
			continue
		}

		var v interface{}
		switch f.Interface().(type) {
		case int, int8, int16, int32, int64:
			v = int(f.Int())
		case uint, uint8, uint16, uint32, uint64:
			v = uint(f.Uint())
		case float32, float64:
			v = f.Float()
		case []byte:
			v = f.Bytes()
		case string, time.Time, time.Location, *time.Time, *time.Location:
			v = f.String()
		case bool:
			v = f.Bool()
		}
		values[tag] = v
	}
	return
}

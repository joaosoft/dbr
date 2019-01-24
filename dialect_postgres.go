package dbr

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type dialectPostgres struct{}

func (d *dialectPostgres) Name() string {
	return string(ConstDialectPostgres)
}

func (d *dialectPostgres) Encode(i interface{}) string {
	value := reflect.ValueOf(i)

	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return "NULL"
		}
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.String:
		return d.EncodeString(value.String())
	case reflect.Bool:
		return d.EncodeBool(value.Bool())
	default:
		switch value.Type() {
		case reflect.TypeOf(time.Time{}):
			return d.EncodeTime(i.(time.Time))
		case reflect.TypeOf([]byte{}):
			return d.EncodeBytes(i.([]byte))
		}
	}

	return fmt.Sprintf("%+v", value.Interface())
}

func (d *dialectPostgres) EncodeString(s string) string {
	return `'` + strings.Replace(s, `'`, `''`, -1) + `'`
}

func (d *dialectPostgres) EncodeBool(b bool) string {
	if b {
		return "TRUE"
	}
	return "FALSE"
}

func (d *dialectPostgres) EncodeTime(t time.Time) string {
	return `'` + t.UTC().Format(ConstTimeFormat) + `'`
}

func (d *dialectPostgres) EncodeBytes(b []byte) string {
	return fmt.Sprintf(`E'\\x%x'`, b)
}

func (d *dialectPostgres) Placeholder() string {
	return ConstPostgresPlaceHolder
}

package dbr

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type DialectPostgres struct{}

func (d *DialectPostgres) Name() string {
	return string(ConstDialectPostgres)
}

func (d *DialectPostgres) Encode(i interface{}) string {
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

func (d *DialectPostgres) EncodeString(s string) string {
	return `'` + strings.Replace(s, `'`, `''`, -1) + `'`
}

func (d *DialectPostgres) EncodeBool(b bool) string {
	if b {
		return "TRUE"
	}
	return "FALSE"
}

func (d *DialectPostgres) EncodeTime(t time.Time) string {
	return `'` + t.UTC().Format(ConstTimeFormat) + `'`
}

func (d *DialectPostgres) EncodeBytes(b []byte) string {
	return fmt.Sprintf(`E'\\x%x'`, b)
}

func (d *DialectPostgres) Placeholder() string {
	return ConstPostgresPlaceHolder
}

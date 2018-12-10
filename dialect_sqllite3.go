package dbr

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type DialectSqlLite3 struct{}

func (d *DialectSqlLite3) Name() string {
	return string(ConstDialectSqlLite3)
}

func (d *DialectSqlLite3) Encode(i interface{}) string {
	value := reflect.ValueOf(i)

	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return fmt.Sprintf("%+v", i)
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

	return fmt.Sprintf("%+v", i)
}

// https://www.sqlite.org/faq.html
func (d *DialectSqlLite3) EncodeString(s string) string {
	return `'` + strings.Replace(s, `'`, `''`, -1) + `'`
}

// https://www.sqlite.org/lang_expr.html
func (d *DialectSqlLite3) EncodeBool(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

// https://www.sqlite.org/lang_datefunc.html
func (d *DialectSqlLite3) EncodeTime(t time.Time) string {
	return `'` + t.UTC().Format(ConstTimeFormat) + `'`
}

// https://www.sqlite.org/lang_expr.html
func (d *DialectSqlLite3) EncodeBytes(b []byte) string {
	return fmt.Sprintf(`X'%x'`, b)
}

func (d *DialectSqlLite3) Placeholder() string {
	return ConstSqlLite3PlaceHolder
}
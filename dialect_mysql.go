package dbr

import (
	"bytes"
	"fmt"
	"reflect"
	"time"
)

type DialectMySql struct {}

func (d *DialectMySql) Name() string {
	return string(ConstDialectMysql)
}

func (d *DialectMySql) Encode(i interface{}) string {
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

// https://dev.mysql.com/doc/refman/5.7/en/string-literals.html
func (d *DialectMySql) EncodeString(s string) string {
	buf := new(bytes.Buffer)

	buf.WriteRune('\'')
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case 0:
			buf.WriteString(`\0`)
		case '\'':
			buf.WriteString(`\'`)
		case '"':
			buf.WriteString(`\"`)
		case '\b':
			buf.WriteString(`\b`)
		case '\n':
			buf.WriteString(`\n`)
		case '\r':
			buf.WriteString(`\r`)
		case '\t':
			buf.WriteString(`\t`)
		case 26:
			buf.WriteString(`\Z`)
		case '\\':
			buf.WriteString(`\\`)
		default:
			buf.WriteByte(s[i])
		}
	}

	buf.WriteRune('\'')
	return buf.String()
}

func (d *DialectMySql) EncodeBool(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func (d *DialectMySql) EncodeTime(t time.Time) string {
	return `'` + t.UTC().Format(ConstTimeFormat) + `'`
}

func (d *DialectMySql) EncodeBytes(b []byte) string {
	return fmt.Sprintf(`0x%x`, b)
}

func (d *DialectMySql) Placeholder() string {
	return ConstMysqlPlaceHolder
}

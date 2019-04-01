package dbr

import "fmt"

type command string

const (
	commandAs     command = "as"
	commandIsNull command = "is_null"
	commandOnNull command = "on_null"
)

type function struct {
	field     interface{}
	fieldName string
	value     interface{}
	command   command
	encode    bool
	encoder   *encoder
}

func As(field interface{}, fieldName string) *function {
	return &function{encode: false, encoder: encoderInstance, field: field, fieldName: fieldName, command: commandAs}
}

func IsNull(field interface{}) *function {
	return &function{encode: false, encoder: encoderInstance, field: field, command: commandIsNull}
}

func OnNull(field interface{}, value interface{}, fieldName string) *function {
	return &function{encode: false, encoder: encoderInstance, field: field, fieldName: fieldName, value: value, command: commandOnNull}
}

func (f *function) String() string {
	field := f.field

	switch f.command {
	case commandAs:
		if stmt, ok := f.field.(*StmtSelect); ok {
			var err error
			field, err = stmt.Build()
			if err != nil {
				return ""
			}
			return fmt.Sprintf("(%s) AS %s", field, f.fieldName)
		}

		var value string
		if f.encode {
			value = f.encoder.encode(field)
		} else {
			value = fmt.Sprintf("%+v", field)
		}

		return fmt.Sprintf("%s AS %s", value, f.fieldName)

	case commandIsNull:
		if stmt, ok := f.field.(*StmtSelect); ok {
			var err error
			field, err = stmt.Build()
			if err != nil {
				return ""
			}
			return fmt.Sprintf("(%s) IS NULL", field)
		}

		var value string
		if f.encode {
			value = f.encoder.encode(field)
		} else {
			value = fmt.Sprintf("%+v", field)
		}

		return fmt.Sprintf("%s IS NULL", value)

	case commandOnNull:
		if stmt, ok := f.field.(*StmtSelect); ok {
			var err error
			field, err = stmt.Build()
			if err != nil {
				return ""
			}
			return fmt.Sprintf("(COALESCE(%s)) AS %s", field, f.fieldName)
		}

		var value string
		if f.encode {
			value = f.encoder.encode(field)
		} else {
			value = fmt.Sprintf("%+v", field)
		}

		return fmt.Sprintf("COALESCE(%s) AS %s", value, f.fieldName)
	}

	return ""
}

package dbr

import "fmt"

type command string

const (
	commandAs command = "as"
)

type function struct {
	field   interface{}
	value   string
	command command
}

func As(field interface{}, name string) *function {
	return &function{field: field, value: name, command: commandAs}
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
			return fmt.Sprintf("(%s) AS %s", field, f.value)
		}

		return fmt.Sprintf("%s AS %s", encodeColumn(field), f.value)
	}

	return ""
}

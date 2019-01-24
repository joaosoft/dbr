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
	switch f.command {
	case commandAs:
		return fmt.Sprintf("%s AS %s", f.field, f.value)
	}

	return ""
}

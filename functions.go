package dbr

import (
	"fmt"
)

type Function struct {
	field interface{}
}

func Field(field interface{}) *Function {
	return &Function{field: field}
}
func (f *Function) As(name string) string {
	return fmt.Sprintf("%s AS %s", f.field, name)
}

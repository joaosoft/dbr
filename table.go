package dbr

import (
	"fmt"
)

type table struct {
	data interface{}
}

func newTable(data interface{}) *table {
	return &table{data: data}
}

func (t *table) Build() (string, error) {

	var value string

	switch stmt := t.data.(type) {
	case *function:
		value = stmt.String()
	case *StmtSelect:
		value, err := stmt.Build()
		if err != nil {
			return "", err
		}
		value = fmt.Sprintf("(%s)", value)
	default:
		value = fmt.Sprintf("%+v", stmt)
	}

	return value, nil
}

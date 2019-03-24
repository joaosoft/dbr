package dbr

import (
	"fmt"
)

type columns []interface{}

func (c columns) Build() (string, error) {

	var query string

	lenC := len(c)
	var err error

	for i, item := range c {
		var value string

		switch stmt := item.(type) {
		case *function:
			value = stmt.String()
		case *StmtSelect:
			value, err = stmt.Build()
			if err != nil {
				return "", err
			}
			value = fmt.Sprintf("(%s)", value)

		default:
			value = encodeColumn(item)
		}

		query += value

		if i+1 < lenC {
			query += ", "
		}
	}

	return query, nil
}

func encodeColumn(v interface{}) string {
	return fmt.Sprintf("%+v", v)
}

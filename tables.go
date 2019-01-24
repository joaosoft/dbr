package dbr

import (
	"fmt"
)

type tables []interface{}

func (t tables) Build() (string, error) {

	var query string

	lenT := len(t)
	var err error

	for i, item := range t {
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
			value = fmt.Sprintf("%+v", item)
		}

		query += value

		if i+1 < lenT {
			query += ", "
		}
	}

	return query, nil
}

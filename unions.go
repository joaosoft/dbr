package dbr

import (
	"fmt"
)

type functionType string
type functions []*function

type function struct {
	functionType functionType
	stmt         *StmtSelect
}

func (u functions) Build() (string, error) {

	if len(u) == 0 {
		return "", nil
	}

	var query string

	for _, union := range u {
		stmt, err := union.stmt.Build()
		query += fmt.Sprintf(" %s %s", string(union.functionType), stmt)

		if err != nil {
			return "", err
		}
	}

	return query, nil
}

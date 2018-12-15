package dbr

import (
	"fmt"
)

type unions []*StmtSelect

func (u unions) Build() (string, error) {
	var query string

	for _, stmtSelect := range u {
		stmt, err := stmtSelect.Build()
		query += fmt.Sprintf(" UNION %s", stmt)

		if err != nil {
			return "", err
		}
	}

	return query, nil
}

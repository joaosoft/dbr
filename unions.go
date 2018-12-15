package dbr

import (
	"fmt"
)

type unionType string
type unions []*union
type union struct {
	unionType unionType
	stmt *StmtSelect
}

const(
	unionNormal = "UNION"
	unionIntersect = "INTERSECT"
	unionExcept = "EXCEPT"
)

func (u unions) Build() (string, error) {
	var query string

	for _, union := range u {
		stmt, err := union.stmt.Build()
		query += fmt.Sprintf(" %s %s", string(union.unionType), stmt)

		if err != nil {
			return "", err
		}
	}

	return query, nil
}

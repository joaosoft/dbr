package dbr

import (
	"fmt"
)

type sets struct {
	list []*set
	db   *db
}

func (s sets) Build() (string, error) {

	if len(s.list) == 0 {
		return "", nil
	}

	var query string

	lenS := len(s.list)
	for i, item := range s.list {
		query += fmt.Sprintf("%s = %s", item.column, s.db.dialect.Encode(item.value))

		if i+1 < lenS {
			query += ", "
		}
	}

	return query, nil
}

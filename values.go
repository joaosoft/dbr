package dbr

import (
	"fmt"
)

type values struct {
	list []interface{}
	db  *Db
}

func (v values) Build() (string, error) {
	var query string

	if len(v.list) == 0 {
		return "", ErrorNumberOfConditionValues
	}

	lenV := len(v.list)
	for i, item := range v.list {
		query += fmt.Sprintf("%s", v.db.dialect.Encode(item))

		if i+1 < lenV {
			query += ", "
		}
	}

	return query, nil
}

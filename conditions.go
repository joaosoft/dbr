package dbr

import (
	"strings"
)

type conditions struct {
	list []*condition
	db   *db
}

func (c conditions) Build() (string, error) {

	if len(c.list) == 0 {
		return "", nil
	}

	var query string

	lenC := len(c.list)
	for i, item := range c.list {
		if strings.Count(item.query, c.db.Dialect.Placeholder()) != len(item.values) {
			return "", ErrorNumberOfConditionValues
		}

		query += item.query
		for _, value := range item.values {
			query = strings.Replace(query, c.db.Dialect.Placeholder(), c.db.Dialect.Encode(value), 1)
		}

		if i+1 < lenC {
			query += " AND "
		}
	}

	return query, nil
}

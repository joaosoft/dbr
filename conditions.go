package dbr

import (
	"strings"
)

type conditions struct {
	list []*condition
	db   *db
}

func (c conditions) Build() (string, error) {
	var query string

	if len(c.list) > 0 {
		lenC := len(c.list)
		for i, item := range c.list {
			if strings.Count(item.query, c.db.dialect.Placeholder()) != len(item.values) {
				return "", ErrorNumberOfConditionValues
			}

			query += item.query
			for _, value := range item.values {
				query = strings.Replace(query, c.db.dialect.Placeholder(), c.db.dialect.Encode(value), 1)
			}

			if i+1 < lenC {
				query += " AND "
			}
		}
	}

	return query, nil
}

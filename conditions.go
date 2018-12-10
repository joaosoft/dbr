package dbr

import (
	"strings"
)

type conditions struct {
	list []*condition
	db   *Db
}

func (c conditions) Build() (string, error) {
	var query string

	if len(c.list) > 0 {
		lenC := len(c.list)
		for i, item := range c.list {
			if strings.Count(item.query, c.db.dialect.Placeholder()) != len(item.values.list) {
				return "", ErrorNumberOfConditionValues
			}

			query += item.query
			for _, value := range item.values.list {
				query = strings.Replace(query, ConstPostgresPlaceHolder, c.db.dialect.Encode(value), 1)
			}

			if i+1 < lenC {
				query += " AND "
			}
		}
	}

	return query, nil
}

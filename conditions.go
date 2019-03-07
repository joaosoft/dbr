package dbr

import (
	"fmt"
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
		condition, err := item.Build()
		if err != nil {
			return "", nil
		}

		query += condition

		if i+1 < lenC {
			query += fmt.Sprintf(" %s ", c.list[i+1].operator)
		}
	}

	return query, nil
}

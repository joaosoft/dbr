package dbr

import (
	"fmt"
)

type direction string

const (
	orderAsc  direction = "asc"
	orderDesc direction = "desc"
)

type order struct {
	column    string
	direction direction
}

type orders []*order

func (o orders) Build() (string, error) {
	var query string
	if len(o) > 0 {
		query = " ORDER BY "
	}

	lenO := len(o)
	for i, item := range o {
		query += fmt.Sprintf("%s %s", item.column, item.direction)

		if i+1 < lenO {
			query += ", "
		}
	}

	return query, nil
}

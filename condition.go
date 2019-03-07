package dbr

import (
	"fmt"
	"strings"
)

type condition struct {
	operator operator
	query    string
	values   []interface{}

	db *db
}

func (c *condition) Build() (string, error) {
	var query string

	if strings.Count(c.query, c.db.Dialect.Placeholder()) != len(c.values) {
		return "", ErrorNumberOfConditionValues
	}

	query += c.query

	var value string
	var err error

	for _, v := range c.values {

		switch stmt := v.(type) {
		case *function:
			value = stmt.String()
		case *StmtSelect:
			value, err = stmt.Build()
			if err != nil {
				return "", err
			}
			value = fmt.Sprintf("(%s)", value)
		default:
			value = fmt.Sprintf("%+v", stmt)
		}

		query = strings.Replace(query, c.db.Dialect.Placeholder(), c.db.Dialect.Encode(value), 1)
	}

	return query, nil
}

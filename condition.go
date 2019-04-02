package dbr

import (
	"fmt"
	"strings"
)

type condition struct {
	operator operator
	query    interface{}
	values   []interface{}

	db *db
}

func (c *condition) Build() (string, error) {
	var query string
	var err error

	switch stmt := c.query.(type) {
	case *StmtSelect:
		query, err = stmt.Build()
		if err != nil {
			return "", err
		}
		query = fmt.Sprintf("(%s)", query)
	default:
		if impl, ok := stmt.(ifunction); ok {
			query, err = impl.Build(c.db)
			if err != nil {
				return "", err
			}
		} else {
			query = fmt.Sprintf("%+v", stmt)
		}
	}

	if strings.Count(query, c.db.Dialect.Placeholder()) != len(c.values) {
		return "", ErrorNumberOfConditionValues
	}

	var value string

	for _, v := range c.values {

		switch stmt := v.(type) {
		case *StmtSelect:
			value, err = stmt.Build()
			if err != nil {
				return "", err
			}
			value = fmt.Sprintf("(%s)", value)
		default:
			if impl, ok := stmt.(ifunction); ok {
				value, err = impl.Build(c.db)
				if err != nil {
					return "", err
				}
			} else {
				value = fmt.Sprintf("%+v", stmt)
			}
		}

		query = strings.Replace(query, c.db.Dialect.Placeholder(), c.db.Dialect.Encode(value), 1)
	}

	return query, nil
}

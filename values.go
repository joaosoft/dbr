package dbr

import (
	"database/sql/driver"
	"fmt"
)

type values struct {
	list []interface{}
	db   *db
}

func newValues(db *db) *values {
	return &values {
		db: db,
	}
}

func (v values) Build() (string, error) {
	var query string

	lenV := len(v.list)
	var err error
	var withoutParentheses bool

	for i, item := range v.list {
		var value string

		switch stmt := item.(type) {
		case *StmtSelect:
			value, err = stmt.Build()
			if err != nil {
				return "", err
			}
			value = fmt.Sprintf("(%s)", value)
		case *values:
			withoutParentheses = true
			value, err = stmt.Build()
			if err != nil {
				return "", err
			}
		case driver.Valuer:
			valuer, err := stmt.Value()
			if err != nil {
				return "", err
			}
			value = fmt.Sprintf("%+v", valuer)
		default:
			if item == nil {
				value = fmt.Sprintf("NULL")
			} else {
				value = fmt.Sprintf("%s", v.db.Dialect.Encode(item))
			}
		}

		query += value

		if i+1 < lenV {
			query += ", "
		}
	}

	if withoutParentheses {
		return query, nil
	}

	return fmt.Sprintf("(%s)", query), nil
}

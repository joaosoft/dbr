package dbr

import (
	"fmt"
)

type columns struct {
	list    []interface{}
	encode  bool
	encoder *encoder

	db *db
}

func newColumns(db *db, encode bool) *columns {
	return &columns{
		db:      db,
		list:    make([]interface{}, 0),
		encode:  encode,
		encoder: &encoder{},
	}
}

func (c columns) Build() (string, error) {

	var query string

	lenC := len(c.list)
	var err error

	for i, item := range c.list {
		var value string

		switch stmt := item.(type) {
		case *function:
			value = stmt.String()
		case *StmtSelect:
			value, err = stmt.Build()
			if err != nil {
				return "", err
			}
			value = fmt.Sprintf("(%s)", value)

		default:
			if c.encode {
				value = c.encoder.encode(item)
			} else {
				value = fmt.Sprintf("%+v", item)
			}
		}

		query += value

		if i+1 < lenC {
			query += ", "
		}
	}

	return query, nil
}

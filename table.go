package dbr

import (
	"fmt"
)

type table struct {
	data interface{}
	db   *db
}

func newTable(db *db, data interface{}) *table {
	return &table{db: db, data: data}
}

func (t *table) Build() (string, error) {

	var value string
	var err error

	switch stmt := t.data.(type) {
	case *StmtSelect:
		value, err = stmt.Build()
		if err != nil {
			return "", err
		}
		value = fmt.Sprintf("(%s)", value)
	default:
		if impl, ok := stmt.(ifunction); ok {
			value, err = impl.Build(t.db)
			if err != nil {
				return "", err
			}
		} else {
			value = fmt.Sprintf("%+v", stmt)
		}
	}

	return value, nil
}

func (t *table) String() string {
	var table string
	var err error

	switch stmt := t.data.(type) {
	case *StmtSelect:
		table, err = stmt.Build()
		if err != nil {
			return ""
		}
	default:
		if impl, ok := stmt.(ifunction); ok {
			table, err = impl.Field(t.db)
			if err != nil {
				return ""
			}
		} else {
			table = fmt.Sprintf("%+v", stmt)
		}
	}

	return table
}

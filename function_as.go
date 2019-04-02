package dbr

import (
	"fmt"
)

type FunctionAs struct {
	field  interface{}
	alias  string
	encode bool

	db *db
}

func newFunctionAs(field interface{}, alias string) *FunctionAs {
	return &FunctionAs{field: field, alias: alias}
}

func (c *FunctionAs) Field(db *db) (string, error) {
	var value string

	if stmt, ok := c.field.(*StmtSelect); ok {
		var err error
		value, err = stmt.Build()
		if err != nil {
			return "", nil
		}

		value = fmt.Sprintf("(%s)", value)
	} else {
		if c.encode {
			value = db.Dialect.EncodeColumn(c.field)
		} else {
			value = fmt.Sprintf("%+v", c.field)
		}
	}
	return value, nil
}

func (c *FunctionAs) Build(db *db) (string, error) {

	field, err := c.Field(db)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("%s AS %s", field, c.alias)

	return query, nil
}

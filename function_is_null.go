package dbr

import (
	"fmt"
)

type FunctionIsNull struct {
	field  interface{}
	encode bool

	db *db
}

func newFunctionIsNull(field interface{}) *FunctionIsNull {
	return &FunctionIsNull{field: field}
}

func (c *FunctionIsNull) Field(db *db) (string, error) {
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

func (c *FunctionIsNull) Build(db *db) (string, error) {

	field, err := c.Field(db)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("%s IS NULL", field)

	return query, nil
}

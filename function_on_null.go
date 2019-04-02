package dbr

import (
	"fmt"
)

type FunctionOnNull struct {
	field       interface{}
	onNullValue interface{}
	alias       string
	encode      bool

	db *db
}

func newFunctionOnNull(field interface{}, onNullValue interface{}, alias string) *FunctionOnNull {
	return &FunctionOnNull{field: field, onNullValue: onNullValue, alias: alias}
}

func (c *FunctionOnNull) Field(db *db) (string, error) {
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

func (c *FunctionOnNull) OnNullValue(db *db) (string, error) {
	var value string

	if stmt, ok := c.onNullValue.(*StmtSelect); ok {
		var err error
		value, err = stmt.Build()
		if err != nil {
			return "", nil
		}

		value = fmt.Sprintf("(%s)", value)
	} else {
		if c.encode {
			value = db.Dialect.EncodeColumn(c.onNullValue)
		} else {
			value = fmt.Sprintf("%+v", c.onNullValue)
		}
	}
	return value, nil
}

func (c *FunctionOnNull) Build(db *db) (string, error) {

	// field
	field, err := c.Field(db)
	if err != nil {
		return "", err
	}

	// onNullValue
	onNullValue, err := c.Field(db)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("COALESCE(%s, %s) AS %s", field, onNullValue, c.alias)

	return query, nil
}

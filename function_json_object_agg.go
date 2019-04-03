package dbr

import (
	"fmt"
)

type functionJsonObjectAgg struct {
	name  interface{}
	value interface{}

	*functionBase
}

func newFunctionJsonObjectAgg(name interface{}, value interface{}) *functionJsonObjectAgg {
	return &functionJsonObjectAgg{functionBase: newFunctionBase(false), name: name, value: value}
}

func (c *functionJsonObjectAgg) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.name)
}

func (c *functionJsonObjectAgg) Value() (string, error) {
	return handleExpression(c.functionBase, c.value)
}

func (c *functionJsonObjectAgg) Build(db *db) (string, error) {
	c.db = db

	// name
	name, err := c.Expression(c.db)
	if err != nil {
		return "", err
	}

	// value
	value, err := c.Value()
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("%s(%s, %s)", constFunctionJsonObjectAgg, name, value)

	return query, nil
}

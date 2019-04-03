package dbr

import (
	"fmt"
)

type functionJsonbObjectAgg struct {
	name  interface{}
	value interface{}

	*functionBase
}

func newFunctionJsonbObjectAgg(name interface{}, value interface{}) *functionJsonbObjectAgg {
	return &functionJsonbObjectAgg{functionBase: newFunctionBase(false), name: name, value: value}
}

func (c *functionJsonbObjectAgg) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.name)
}

func (c *functionJsonbObjectAgg) Value() (string, error) {
	return handleExpression(c.functionBase, c.value)
}

func (c *functionJsonbObjectAgg) Build(db *db) (string, error) {
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

	query := fmt.Sprintf("%s(%s, %s)", constFunctionJsonbObjectAgg, name, value)

	return query, nil
}

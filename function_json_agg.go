package dbr

import (
	"fmt"
)

type functionJsonAgg struct {
	expression interface{}

	*functionBase
}

func newFunctionJsonAgg(expression interface{}) *functionJsonAgg {
	return &functionJsonAgg{functionBase: newFunctionBase(false), expression: expression}
}

func (c *functionJsonAgg) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.expression)
}

func (c *functionJsonAgg) Build(db *db) (string, error) {
	c.db = db

	expression, err := c.Expression(db)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("%s(%s)", constFunctionJsonAgg, expression)

	return query, nil
}

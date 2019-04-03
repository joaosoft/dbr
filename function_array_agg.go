package dbr

import (
	"fmt"
)



type functionArrayAgg struct {
	expression interface{}

	*functionBase
}

func newFunctionArrayAgg(expression interface{}) *functionArrayAgg {
	return &functionArrayAgg{functionBase: newFunctionBase(false), expression: expression}
}

func (c *functionArrayAgg) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.expression)
}

func (c *functionArrayAgg) Build(db *db) (string, error) {
	c.db = db

	expression, err := c.Expression(c.db)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("%s(%s)", constFunctionArrayAgg, expression)

	return query, nil
}

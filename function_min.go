package dbr

import (
	"fmt"
)

type functionMin struct {
	expression interface{}

	*functionBase
}

func newFunctionMin(expression interface{}) *functionMin {
	return &functionMin{functionBase: newFunctionBase(false), expression: expression}
}

func (c *functionMin) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.expression)
}

func (c *functionMin) Build(db *db) (string, error) {
	c.db = db

	expression, err := c.Expression(db)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("%s(%s)", constFunctionMin, expression)

	return query, nil
}

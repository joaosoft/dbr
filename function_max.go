package dbr

import (
	"fmt"
)

type functionMax struct {
	expression interface{}

	*functionBase
}

func newFunctionMax(expression interface{}) *functionMax {
	return &functionMax{functionBase: newFunctionBase(false), expression: expression}
}

func (c *functionMax) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.expression)
}

func (c *functionMax) Build(db *db) (string, error) {
	c.db = db

	expression, err := c.Expression(db)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("%s(%s)", constFunctionMax, expression)

	return query, nil
}

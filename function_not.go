package dbr

import (
	"fmt"
)

type functionNot struct {
	expression interface{}

	*functionBase
}

func newFunctionNot(expression interface{}) *functionNot {
	return &functionNot{functionBase: newFunctionBase(false), expression: expression}
}

func (c *functionNot) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.expression)
}

func (c *functionNot) Build(db *db) (string, error) {
	c.db = db

	expression, err := c.Expression(c.db)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("%s %s", constFunctionNot, expression)

	return query, nil
}

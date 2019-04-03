package dbr

import (
	"fmt"
)

type functionSum struct {
	expression interface{}

	*functionBase
}

func newFunctionSum(expression interface{}) *functionSum {
	return &functionSum{functionBase: newFunctionBase(false), expression: expression}
}

func (c *functionSum) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.expression)
}

func (c *functionSum) Build(db *db) (string, error) {
	c.db = db

	expression, err := c.Expression(db)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("%s(%s)", constFunctionSum, expression)

	return query, nil
}

package dbr

import (
	"fmt"
)

type functionEvery struct {
	expression interface{}

	*functionBase
}

func newFunctionEvery(expression interface{}) *functionEvery {
	return &functionEvery{functionBase: newFunctionBase(false), expression: expression}
}

func (c *functionEvery) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.expression)
}

func (c *functionEvery) Build(db *db) (string, error) {
	c.db = db

	expression, err := c.Expression(c.db)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("%s(%s)", constFunctionEvery, expression)

	return query, nil
}

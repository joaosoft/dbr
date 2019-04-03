package dbr

import (
	"fmt"
)

type functionJsonbAgg struct {
	expression interface{}

	*functionBase
}

func newFunctionJsonbAgg(expression interface{}) *functionJsonbAgg {
	return &functionJsonbAgg{functionBase: newFunctionBase(false), expression: expression}
}

func (c *functionJsonbAgg) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.expression)
}

func (c *functionJsonbAgg) Build(db *db) (string, error) {
	c.db = db

	expression, err := c.Expression(db)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("%s(%s)", constFunctionJsonbAgg, expression)

	return query, nil
}

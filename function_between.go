package dbr

import (
	"fmt"
)

type functionBetween struct {
	expression interface{}
	low        interface{}
	operator   operator
	high       interface{}

	*functionBase
}

func newFunctionBetween(expression interface{}, low interface{}, operator operator, high interface{}) *functionBetween {
	return &functionBetween{functionBase: newFunctionBase(false, false), expression: expression, low: low, operator: operator, high: high}
}

func (c *functionBetween) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.expression)
}

func (c *functionBetween) Build(db *db) (string, error) {
	c.db = db

	expression, err := c.Expression(c.db)
	if err != nil {
		return "", err
	}

	low, err := handleExpression(c.functionBase, c.low)
	if err != nil {
		return "", err
	}

	high, err := handleExpression(c.functionBase, c.high)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("%s %s %s %s %s", expression, constFunctionBetween, low, c.operator, high)

	return query, nil
}

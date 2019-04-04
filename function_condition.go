package dbr

import (
	"fmt"
)

type functionCondition struct {
	expression interface{}
	comparator comparator
	value      interface{}

	*functionBase
}

func newFunctionCondition(expression interface{}, comparator comparator, value interface{}) *functionCondition {
	return &functionCondition{functionBase: newFunctionBase(false, false), expression: expression, comparator: comparator, value: value}
}

func (c *functionCondition) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.expression)
}

func (c *functionCondition) Build(db *db) (string, error) {
	c.db = db

	expression, err := c.Expression(c.db)
	if err != nil {
		return "", err
	}

	value, err := handleExpression(c.functionBase, c.value)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("%s %s %s", expression, c.comparator, value)

	return query, nil
}

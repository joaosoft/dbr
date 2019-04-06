package dbr

import (
	"fmt"
)

type functionOperation struct {
	expression interface{}
	operation  operation
	value      interface{}

	*functionBase
}

func newFunctionOperation(expression interface{}, operation operation, value interface{}) *functionOperation {
	return &functionOperation{functionBase: newFunctionBase(false, false), expression: expression, operation: operation, value: value}
}

func (c *functionOperation) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.expression)
}

func (c *functionOperation) Build(db *db) (string, error) {
	c.db = db

	expression, err := handleBuild(c.functionBase, c.expression)
	if err != nil {
		return "", err
	}

	value, err := handleBuild(c.functionBase, c.value)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("%s %s %s", expression, c.operation, value)

	return query, nil
}

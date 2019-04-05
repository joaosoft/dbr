package dbr

import (
	"fmt"
)

type functionIsNull struct {
	expression interface{}

	*functionBase
}

func newFunctionIsNull(expression interface{}) *functionIsNull {
	return &functionIsNull{functionBase: newFunctionBase(false, false), expression: expression}
}

func (c *functionIsNull) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.expression)
}

func (c *functionIsNull) Build(db *db) (string, error) {
	c.db = db

	expression, err := handleBuild(c.functionBase, c.expression)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("%s %s", expression, constFunctionIsNull)

	return query, nil
}

package dbr

import (
	"fmt"
)

type functionCast struct {
	expression interface{}
	dataType   dataType

	*functionBase
}

func newFunctionCast(expression interface{}, dataType dataType) *functionCast {
	return &functionCast{functionBase: newFunctionBase(false, false), expression: expression, dataType: dataType}
}

func (c *functionCast) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.expression)
}

func (c *functionCast) Build(db *db) (string, error) {
	c.db = db

	expression, err := handleBuild(c.functionBase, c.expression)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("%s(%s %s %s)", constFunctionCast, expression, constFunctionAs, c.dataType)

	return query, nil
}

package dbr

import (
	"fmt"
)

type functionOnNull struct {
	expression  interface{}
	onNullValue interface{}
	alias       string

	*functionBase
}

func newFunctionOnNull(expression interface{}, onNullValue interface{}, alias string) *functionOnNull {
	return &functionOnNull{functionBase: newFunctionBase(false, false), expression: expression, onNullValue: onNullValue, alias: alias}
}

func (c *functionOnNull) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.expression)
}

func (c *functionOnNull) Build(db *db) (string, error) {
	c.db = db

	// expression
	expression, err := handleBuild(c.functionBase, c.expression)
	if err != nil {
		return "", err
	}

	// onNullValue
	onNullValue, err := handleBuild(c.functionBase, c.onNullValue)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("%s(%s, %s) AS %s", constFunctionOnNull, expression, onNullValue, c.alias)

	return query, nil
}

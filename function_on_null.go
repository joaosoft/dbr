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
	return &functionOnNull{functionBase: newFunctionBase(false), expression: expression, onNullValue: onNullValue, alias: alias}
}

func (c *functionOnNull) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.expression)
}

func (c *functionOnNull) OnNullValue() (string, error) {
	return handleExpression(c.functionBase, c.onNullValue)
}

func (c *functionOnNull) Build(db *db) (string, error) {
	c.db = db

	// expression
	expression, err := c.Expression(db)
	if err != nil {
		return "", err
	}

	// onNullValue
	onNullValue, err := c.OnNullValue()
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("%s(%s, %s) AS %s", constFunctionOnNull, expression, onNullValue, c.alias)

	return query, nil
}

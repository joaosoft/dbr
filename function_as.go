package dbr

import (
	"fmt"
)

type functionAs struct {
	expression interface{}
	alias      string

	*functionBase
}

func newFunctionAs(expression interface{}, alias string) *functionAs {
	return &functionAs{functionBase: newFunctionBase(false), expression: expression, alias: alias}
}

func (c *functionAs) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.expression)
}

func (c *functionAs) Build(db *db) (string, error) {
	c.db = db

	expression, err := c.Expression(c.db)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("%s %s %s", expression, constFunctionAs, c.alias)

	return query, nil
}
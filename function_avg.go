package dbr

import (
	"fmt"
)

const (
	ConstFunctionAvg = "AVG"
)

type functionAvg struct {
	expression interface{}

	*functionBase
}

func newFunctionAvg(expression interface{}) *functionAvg {
	return &functionAvg{functionBase: newFunctionBase(false), expression: expression}
}

func (c *functionAvg) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.expression)
}

func (c *functionAvg) Build(db *db) (string, error) {
	c.db = db

	expression, err := c.Expression(c.db)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("%s(%s)", ConstFunctionAvg, expression)

	return query, nil
}

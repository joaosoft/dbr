package dbr

import (
	"fmt"
)

type functionXmlAgg struct {
	expression interface{}

	*functionBase
}

func newFunctionXmlAgg(expression interface{}) *functionXmlAgg {
	return &functionXmlAgg{functionBase: newFunctionBase(false), expression: expression}
}

func (c *functionXmlAgg) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.expression)
}

func (c *functionXmlAgg) Build(db *db) (string, error) {
	c.db = db

	expression, err := c.Expression(db)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("%s(%s)", constFunctionXmlAgg, expression)

	return query, nil
}

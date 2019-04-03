package dbr

import (
	"fmt"
)

type functionStringAgg struct {
	expression interface{}
	delimiter  interface{}

	*functionBase
}

func newFunctionStringAgg(expression interface{}, delimiter interface{}) *functionStringAgg {
	return &functionStringAgg{functionBase: newFunctionBase(false), expression: expression, delimiter: delimiter}
}

func (c *functionStringAgg) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.expression)
}

func (c *functionStringAgg) Delimiter() (string, error) {
	return handleExpression(c.functionBase, c.delimiter)
}

func (c *functionStringAgg) Build(db *db) (string, error) {
	c.db = db

	// expression
	expression, err := c.Expression(c.db)
	if err != nil {
		return "", err
	}

	// delimiter
	delimiter, err := c.Delimiter()
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("%s(%s, %s)", constFunctionStringAgg, expression, delimiter)

	return query, nil
}

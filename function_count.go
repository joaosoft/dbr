package dbr

import (
	"fmt"
)

type functionCount struct {
	expression interface{}
	distinct   bool

	*functionBase
}

func newFunctionCount(expression interface{}, distinct bool) *functionCount {
	return &functionCount{functionBase: newFunctionBase(false), expression: expression, distinct: distinct}
}

func (c *functionCount) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.expression)
}

func (c *functionCount) Build(db *db) (string, error) {
	c.db = db

	expression, err := c.Expression(db)
	if err != nil {
		return "", err
	}

	var distinct string
	if c.distinct {
		distinct = "DISTINCT "
	}

	query := fmt.Sprintf("%s(%s%s)", constFunctionCount, distinct, expression)

	return query, nil
}

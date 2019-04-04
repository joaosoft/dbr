package dbr

import (
	"fmt"
)

type functionIn struct {
	expressions []interface{}

	*functionBase
}

func newFunctionIn(expressions ...interface{}) *functionIn {
	return &functionIn{functionBase: newFunctionBase(false), expressions: expressions}
}

func (c *functionIn) Expression(db *db) (string, error) {
	c.db = db

	return "", nil
}

func (c *functionIn) Build(db *db) (string, error) {
	c.db = db

	var arguments string

	lenE := len(c.expressions)
	for i, expression := range c.expressions {
		value, err := handleExpression(c.functionBase, expression)
		if err != nil {
			return "", err
		}

		arguments += value

		if i+1 < lenE {
			arguments += ", "
		}
	}

	query := fmt.Sprintf("%s(%s)", constFunctionIn, arguments)

	return query, nil
}

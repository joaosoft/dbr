package dbr

import (
	"fmt"
)

type functionNotIn struct {
	expressions []interface{}

	*functionBase
}

func newFunctionNotIn(expressions ...interface{}) *functionNotIn {
	return &functionNotIn{functionBase: newFunctionBase(false), expressions: expressions}
}

func (c *functionNotIn) Expression(db *db) (string, error) {
	c.db = db

	return "", nil
}

func (c *functionNotIn) Build(db *db) (string, error) {
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

	query := fmt.Sprintf("%s(%s)", constFunctionNotIn, arguments)

	return query, nil
}

package dbr

import (
	"fmt"
)

type functionGeneric struct {
	name      string
	arguments []interface{}

	*functionBase
}

func newFunctionGeneric(name string, arguments ...interface{}) *functionGeneric {
	return &functionGeneric{functionBase: newFunctionBase(false), name: name, arguments: arguments}
}

func (c *functionGeneric) Expression(db *db) (string, error) {
	c.db = db

	return handleExpression(c.functionBase, c.name)
}

func (c *functionGeneric) Build(db *db) (string, error) {
	c.db = db

	var arguments string

	lenArgs := len(c.arguments)
	for i, argument := range c.arguments {
		expression, err := handleExpression(c.functionBase, argument)
		if err != nil {
			return "", err
		}

		arguments += expression

		if i < lenArgs-1 {
			arguments += ", "
		}
	}

	query := fmt.Sprintf("%s(%s)", c.name, arguments)

	return query, nil
}

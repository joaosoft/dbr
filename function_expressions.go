package dbr

import (
	"fmt"
)

type functionExpressions struct {
	comma       bool
	expressions []interface{}

	*functionBase
}

func newFunctionExpressions(comma bool, expressions ...interface{}) *functionExpressions {
	return &functionExpressions{functionBase: newFunctionBase(false, false), expressions: expressions, comma: comma}
}

func (c *functionExpressions) Expression(db *db) (string, error) {
	c.db = db

	if len(c.expressions) == 0 {
		return "", nil
	}

	return handleExpression(c.functionBase, c.expressions[0])
}

func (c *functionExpressions) Build(db *db) (string, error) {
	c.db = db

	var expressions string
	var addComma bool

	lenArgs := len(c.expressions)
	for i, argument := range c.expressions {
		expression, err := handleBuild(c.functionBase, argument)
		if err != nil {
			return "", err
		}

		expressions += expression

		if c.comma {
			if expression == constFunctionOpenParentheses {
				addComma = true
				continue
			}

			if expression == constFunctionCloseParentheses {
				addComma = false
				goto next
			}

			if i < lenArgs-1 && c.expressions[i+1] == constFunctionCloseParentheses {
				continue
			}

			if addComma {
				expressions += constFunctionComma
			}
		}

	next:

		if i < lenArgs-1 {
			expressions += " "
		}
	}

	query := fmt.Sprintf("%s", expressions)

	return query, nil
}

package dbr

import (
	"fmt"
)

type functionFilter struct {
	value      interface{}
	conditions *conditions

	*functionBase
}

func newFunctionFilter(value interface{}) *functionFilter {
	return &functionFilter{functionBase: newFunctionBase(false, false), value: value, conditions: newConditions(nil)}
}

func (c *functionFilter) Where(query interface{}, values ...interface{}) *functionFilter {
	c.conditions.list = append(c.conditions.list, &condition{operator: OperatorAnd, query: query, values: values})
	return c
}

func (c *functionFilter) Build(db *db) (string, error) {
	c.db = db

	base := newFunctionBase(false, false, db)
	value, err := handleBuild(base, c.value)
	if err != nil {
		return "", err
	}

	conditions, err := c.conditions.Build(db)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s (%s%s)", value, constFunctionFilter, conditions), nil
}

package dbr

import (
	"fmt"
)

type functionNow struct {
	*functionBase
}

func newFunctionNow() *functionNow {
	return &functionNow{functionBase: newFunctionBase(false)}
}

func (c *functionNow) Expression(db *db) (string, error) {
	c.db = db

	return "", nil
}

func (c *functionNow) Build(db *db) (string, error) {
	c.db = db

	query := fmt.Sprintf("%s()", constFunctionNow)

	return query, nil
}

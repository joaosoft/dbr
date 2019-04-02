package dbr

import (
	"fmt"
	"strings"
)

type FunctionCase struct {
	alias   *string
	onWhens caseWhens
	onElse  *caseElse
}

func newFunctionCase(alias ...string) *FunctionCase {
	funcCase := &FunctionCase{onWhens: newCaseWhens()}

	if len(alias) > 0 {
		funcCase.alias = &alias[0]
	}

	return funcCase
}

func (c *FunctionCase) When(condition interface{}, result interface{}) *FunctionCase {
	c.onWhens = append(c.onWhens, newCaseWhen(condition, result))

	return c
}

func (c *FunctionCase) Else(result interface{}) *FunctionCase {
	c.onElse = newCaseElse(result)

	return c
}

func (c *FunctionCase) Field(db *db) (string, error) {
	return "", nil
}

func (c *FunctionCase) Build(db *db) (string, error) {
	var value string
	var query string

	onWhens, err := c.onWhens.Build(db)
	if err != nil {
		return "", err
	}
	value += onWhens

	onElse, err := c.onElse.Build(db)
	if err != nil {
		return "", err
	}

	if len(onElse) > 0 {
		value += fmt.Sprintf(" %s", onElse)
	}

	query = fmt.Sprintf("(CASE %s END)", value)

	if c.alias != nil && len(strings.TrimSpace(*c.alias)) > 0 {
		query += fmt.Sprintf(" AS %s", *c.alias)
	}

	return query, nil
}

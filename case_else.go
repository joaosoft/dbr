package dbr

import (
	"fmt"
)

type onCaseElse struct {
	result interface{}
}

func newCaseElse(result interface{}) *onCaseElse {
	return &onCaseElse{result: result}
}

func (c *onCaseElse) Build(db *db) (string, error) {
	var query string
	var err error
	var result string

	// result
	switch stmt := c.result.(type) {
	case *StmtSelect:
		result, err = stmt.Build()
		if err != nil {
			return "", err
		}
		result = fmt.Sprintf("(%s)", result)
	default:
		if impl, ok := stmt.(iFunction); ok {
			result, err = impl.Build(db)
			if err != nil {
				return "", err
			}
		} else {
			result = fmt.Sprintf("%+v", stmt)
		}
	}

	query = fmt.Sprintf("%s %s", constFunctionElse, result)

	return query, nil
}

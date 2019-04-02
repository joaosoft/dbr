package dbr

import (
	"fmt"
)

type caseWhen struct {
	condition interface{}
	result    interface{}
}

func newCaseWhen(condition interface{}, result interface{}) *caseWhen {
	return &caseWhen{condition: condition, result: result}
}

func (c *caseWhen) Build(db *db) (string, error) {
	var query string
	var err error
	var condition string
	var result string

	// condition
	switch stmt := c.condition.(type) {
	case *StmtSelect:
		condition, err = stmt.Build()
		if err != nil {
			return "", err
		}
		condition = fmt.Sprintf("(%s)", result)
	default:
		if impl, ok := stmt.(ifunction); ok {
			condition, err = impl.Build(db)
			if err != nil {
				return "", err
			}
		} else {
			condition = fmt.Sprintf("%+v", stmt)
		}
	}

	// result
	switch stmt := c.result.(type) {
	case *StmtSelect:
		result, err = stmt.Build()
		if err != nil {
			return "", err
		}
		result = fmt.Sprintf("(%s)", result)
	default:
		if impl, ok := stmt.(ifunction); ok {
			result, err = impl.Build(db)
			if err != nil {
				return "", err
			}
		} else {
			result = fmt.Sprintf("%+v", stmt)
		}
	}

	query = fmt.Sprintf("WHEN %s THEN %s", condition, result)

	return query, nil
}

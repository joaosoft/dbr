package dbr

import (
	"fmt"
	"strings"
)

type columns []interface{}

func (c columns) Build() (string, error) {

	var query string

	lenC := len(c)
	var err error

	for i, item := range c {
		var value string

		switch stmt := item.(type) {
		case *function:
			value = stmt.String()
		case *StmtSelect:
			value, err = stmt.Build()
			if err != nil {
				return "", err
			}
			value = fmt.Sprintf("(%s)", value)

		default:
			value = encodeColumn(item)
		}

		query += value

		if i+1 < lenC {
			query += ", "
		}
	}

	return query, nil
}

func encodeColumn(v interface{}) string {
	value := fmt.Sprintf("%+v", v)

	switch v.(type) {
	case string:
		if !strings.ContainsAny(value, `"*`) {
			value = fmt.Sprintf(`"%s"`, value)
			value = strings.Replace(value, `.`, `"."`, 1)
		}
	}

	return value
}

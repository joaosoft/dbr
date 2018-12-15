package dbr

import (
	"fmt"
	"strings"
)

type groupBy []string

func (g groupBy) Build() (string, error) {
	if len(g) > 0 {
		return fmt.Sprintf(" GROUP BY %s", strings.Join(g, ", ")), nil
	}

	return "", nil
}

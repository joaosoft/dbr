package dbr

import (
	"strings"
)

type columns []string

func (c columns) Build() (string, error) {
	return strings.Join(c, ", "), nil
}

package dbr

import (
	"strings"
)

type tables []string

func (c tables) Build() (string, error) {
	return strings.Join(c, ", "), nil
}

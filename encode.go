package dbr

import (
	"fmt"
	"strings"
)

var(
	encoderInstance = &encoder{}
)

type encoder struct{}

func (e *encoder) encode(v interface{}) string {
	value := fmt.Sprintf("%+v", v)

	switch v.(type) {
	case string:
		if !strings.ContainsAny(value, `*`) {
			value = fmt.Sprintf(`"%s"`, value)
			value = strings.Replace(value, `.`, `"."`, 1)
		}
	}

	return value
}

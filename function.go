package dbr

type ifunction interface {
	Build(db *db) (string, error)
	Field(db *db) (string, error)
}

func As(field interface{}, alias string) *FunctionAs {
	return newFunctionAs(field, alias)
}

func IsNull(field interface{}) *FunctionIsNull {
	return newFunctionIsNull(field)
}

func Case(alias ...string) *FunctionCase {
	return newFunctionCase(alias...)
}

func OnNull(field interface{}, onNullValue interface{}, alias string) *FunctionOnNull {
	return newFunctionOnNull(field, onNullValue, alias)
}

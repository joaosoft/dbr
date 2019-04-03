package dbr

import "fmt"

type ifunction interface {
	Build(db *db) (string, error)
	Expression(db *db) (string, error)
}

type functionBase struct {
	encode bool
	db     *db
}

func newFunctionBase(encode bool) *functionBase {
	return &functionBase{encode: encode}
}

func Function(name string, arguments ...interface{}) *functionGeneric {
	return newFunctionGeneric(name, arguments...)
}

func As(expression interface{}, alias string) *functionAs {
	return newFunctionAs(expression, alias)
}

func Count(expression interface{}, distinct ...bool) *functionCount {
	var isDistinct bool
	if len(distinct) > 0 {
		isDistinct = distinct[0]
	}

	return newFunctionCount(expression, isDistinct)
}

func IsNull(expression interface{}) *functionIsNull {
	return newFunctionIsNull(expression)
}

func Case(alias ...string) *functionCase {
	return newFunctionCase(alias...)
}

func OnNull(expression interface{}, onNullValue interface{}, alias string) *functionOnNull {
	return newFunctionOnNull(expression, onNullValue, alias)
}

func Min(expression interface{}) *functionMin {
	return newFunctionMin(expression)
}

func Max(expression interface{}) *functionMax {
	return newFunctionMax(expression)
}

func Sum(expression interface{}) *functionSum {
	return newFunctionSum(expression)
}

func Avg(expression interface{}) *functionAvg {
	return newFunctionAvg(expression)
}

func Every(expression interface{}) *functionEvery {
	return newFunctionEvery(expression)
}

func Now() *functionNow {
	return newFunctionNow()
}

func StringAgg(expression interface{}, delimiter interface{}) *functionStringAgg {
	return newFunctionStringAgg(expression, delimiter)
}

func XmlAgg(expression interface{}) *functionXmlAgg {
	return newFunctionXmlAgg(expression)
}

func ArrayAgg(expression interface{}) *functionArrayAgg {
	return newFunctionArrayAgg(expression)
}

func JsonAgg(expression interface{}) *functionJsonAgg {
	return newFunctionJsonAgg(expression)
}

func JsonbAgg(expression interface{}) *functionJsonbAgg {
	return newFunctionJsonbAgg(expression)
}

func JsonObjectAgg(name interface{}, value interface{}) *functionJsonObjectAgg {
	return newFunctionJsonObjectAgg(name, value)
}

func JsonbObjectAgg(name interface{}, value interface{}) *functionJsonbObjectAgg {
	return newFunctionJsonbObjectAgg(name, value)
}

func handleExpression(base *functionBase, expression interface{}) (string, error) {
	var value string

	if stmt, ok := expression.(*StmtSelect); ok {
		var err error
		value, err = stmt.Build()
		if err != nil {
			return "", nil
		}
	} else {
		if base.encode {
			value = base.db.Dialect.EncodeColumn(expression)
		} else {
			value = fmt.Sprintf("%+v", expression)
		}
	}
	return value, nil
}

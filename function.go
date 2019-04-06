package dbr

type ifunction interface {
	Build(db *db) (string, error)
	Expression(db *db) (string, error)
}

type functionBase struct {
	isColumn bool
	encode   bool
	db       *db
}

func newFunctionBase(encode bool, isColumn bool, database ...*db) *functionBase {
	var theDb *db
	if len(database) > 0 {
		theDb = database[0]
	}
	return &functionBase{isColumn: isColumn, encode: encode, db: theDb}
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

func Min(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionMax, expression)
}

func Max(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionMin, expression)
}

func Sum(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionSum, expression)
}

func Avg(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionAvg, expression)
}

func Every(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionEvery, expression)
}

func Now() *functionGeneric {
	return newFunctionGeneric(constFunctionNow)
}

func User() *functionGeneric {
	return newFunctionGeneric(constFunctionUser)
}

func StringAgg(expression interface{}, delimiter interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionStringAgg, expression, delimiter)
}

func XmlAgg(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionXmlAgg, expression)
}

func ArrayAgg(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionArrayAgg, expression)
}

func JsonAgg(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionJsonAgg, expression)
}

func JsonbAgg(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionJsonbAgg, expression)
}

func JsonObjectAgg(name interface{}, value interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionJsonObjectAgg, name, value)
}

func JsonbObjectAgg(name interface{}, value interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionJsonbObjectAgg, name, value)
}

func Cast(expression interface{}, dataType dataType) *functionCast {
	return newFunctionCast(expression, dataType)
}

func Not(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionNot, expression)
}

func In(expressions ...interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionIn, expressions...)
}

func NotIn(expressions ...interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionNotIn, expressions...)
}

func Upper(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionUpper, expression)
}

func Lower(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionLower, expression)
}

func Trim(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionTrim, expression)
}

func Concat(expressions ...interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionNotIn, expressions...)
}

func InitCap(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionInitCap, expression)
}

func Length(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionLength, expression)
}

func Left(expression interface{}, n int) *functionGeneric {
	return newFunctionGeneric(constFunctionLeft, expression, n)
}

func Right(expression interface{}, n int) *functionGeneric {
	return newFunctionGeneric(constFunctionRight, expression, n)
}

func Md5(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionMd5, expression)
}

func Replace(expression interface{}, from interface{}, to interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionReplace, expression, from, to)
}

func Repeat(expression interface{}, n int) *functionGeneric {
	return newFunctionGeneric(constFunctionRepeat, expression, n)
}

func Any(expressions ...interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionAny, expressions...)
}

func Some(expressions ...interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionSome, expressions...)
}

func Condition(expression interface{}, comparator comparator, value interface{}) *functionCondition {
	return newFunctionCondition(expression, comparator, value)
}

func Operation(expression interface{}, operation operation, value interface{}) *functionOperation {
	return newFunctionOperation(expression, operation, value)
}

func Abs(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionAbs, expression)
}

func Sqrt(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionSqrt, expression)
}

func Random(expression interface{}) *functionGeneric {
	return newFunctionGeneric(constFunctionRandom, expression)
}

func Between(expression interface{}, low interface{}, high interface{}, operator ...operator) *functionBetween {
	theOperator := OperatorAnd

	if len(operator) > 0 {
		theOperator = operator[0]
	}

	return newFunctionBetween(expression, low, theOperator, high)
}

func BetweenOr(expression interface{}, low interface{}, high interface{}) *functionBetween {
	return newFunctionBetween(expression, low, OperatorOr, high)
}

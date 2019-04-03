package dbr

const (
	constDialectPostgres DialectName = "postgres"
	constDialectMysql    DialectName = "mysql"
	constDialectSqlLite3 DialectName = "sqlite3"

	constPostgresPlaceHolder = "?"
	constMysqlPlaceHolder    = "?"
	constSqlLite3PlaceHolder = "?"
	constTimeFormat          = "2006-01-02 15:04:05.000000"

	constLoadOptionDefault loadOption = "db"
	constLoadOptionRead    loadOption = "db.read"
	constLoadOptionWrite   loadOption = "db.write"

	constFunctionArrayAgg       = "ARRAY_AGG"
	constFunctionAs             = "As"
	constFunctionCase           = "CASE"
	constFunctionEnd            = "END"
	constFunctionCount          = "COUNT"
	constFunctionEvery          = "EVERY"
	constFunctionIsNull         = "IS NULL"
	constFunctionJsonAgg        = "JSON_AGG"
	constFunctionJsonObjectAgg  = "JSON_OBJECT_AGG"
	constFunctionJsonbAgg       = "JSONB_AGG"
	constFunctionJsonbObjectAgg = "JSONB_OBJECT_AGG"
	constFunctionMax            = "MAX"
	constFunctionMin            = "MIN"
	constFunctionOnNull         = "COALESCE"
	constFunctionStringAgg      = "STRING_AGG"
	constFunctionSum            = "SUM"
	constFunctionWhen           = "WHEN"
	constFunctionThen           = "THEN"
	constFunctionElse           = "ELSE"
	constFunctionGroupBy        = "GROUP BY"
	constFunctionDelete         = "DELETE"
	constFunctionFrom           = "FROM"
	constFunctionWhere          = "WHERE"
	constFunctionReturning      = "RETURNING"
	constFunctionValues         = "VALUES"
	constFunctionInsert         = "INSERT"
	constFunctionInto           = "INTO"
	constFunctionRecursive      = "RECURSIVE"
	constFunctionWith           = "WITH"
	constFunctionUnionNormal    = "UNION"
	constFunctionUnionIntersect = "INTERSECT"
	constFunctionUnionExcept    = "EXCEPT"
	constFunctionNull           = "NULL"
	constFunctionDistinct       = "DISTINCT"
	constFunctionDistinctOn     = "DISTINCT ON"
	constFunctionSelect         = "SELECT"
	constFunctionHaving         = "HAVING"
	constFunctionLimit          = "LIMIT"
	constFunctionOffset         = "OFFSET"
	constFunctionUpdate         = "UPDATE"
	constFunctionSet            = "SET"
	constFunctionDoNothing      = "DO NOTHING"
	constFunctionOnConflict     = "ON CONFLICT"
	constFunctionOnConstraint   = "ON CONSTRAINT"
	constFunctionDoUpdateSet    = "DO UPDATE SET"
	constFunctionOrderBy        = "ORDER BY"
	constFunctionNow            = "NOW"
	constFunctionXmlAgg         = "XMLAGG"

	constFunctionJoin      Join = "JOIN"
	constFunctionLeftJoin  Join = "LEFT JOIN"
	constFunctionRightJoin Join = "RIGHT JOIN"
	constFunctionFullJoin  Join = "FULL JOIN"

	OrderAsc  direction = "ASC"
	OrderDesc direction = "DESC"

	ExecuteOperation SqlOperation = "EXECUTE"
	SelectOperation  SqlOperation = "SELECT"
	InsertOperation  SqlOperation = "INSERT"
	UpdateOperation  SqlOperation = "UPDATE"
	DeleteOperation  SqlOperation = "DELETE"

	constOperatorAnd operator = "AND"
	constOperatorOr  operator = "OR"

	constMySqlBoolTrue  = "1"
	constMySqlBoolFalse = "0"

	constPostgresBoolTrue  = "TRUE"
	constPostgresBoolFalse = "FALSE"

	constSqlLite3BoolTrue  = "1"
	constSqlLite3BoolFalse = "0"
)

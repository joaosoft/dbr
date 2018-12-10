package dbr

import "time"

type DialectName string

const (
	ConstDialectPostgres DialectName = "postgres"
	ConstDialectMysql    DialectName = "mysql"
	ConstDialectSqlLite3 DialectName = "sqlite3"
)

type Dialect interface {
	Name() string
	Encode(i interface{}) string
	EncodeString(s string) string
	EncodeBool(b bool) string
	EncodeTime(t time.Time) string
	EncodeBytes(b []byte) string
	Placeholder() string
}

func newDialect(name string) Dialect {
	switch name {
	case string(ConstDialectPostgres):
		return &DialectPostgres{}
	case string(ConstDialectMysql):
		return &DialectMySql{}
	case string(ConstDialectSqlLite3):
		return &DialectSqlLite3{}
	}

	return nil
}

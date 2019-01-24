package dbr

import (
	"database/sql"
	"strings"
)

type StmtExecute struct {
	query  string
	values values

	Db *db
}

func newStmtExecute(db *db, query string) *StmtExecute {
	return &StmtExecute{Db: db, query: query, values: values{db: db}}
}

func (stmt *StmtExecute) Values(valuesList ...interface{}) *StmtExecute {
	stmt.values.list = append(stmt.values.list, valuesList...)
	return stmt
}

func (stmt *StmtExecute) Build() (string, error) {
	query := stmt.query

	if strings.Count(query, stmt.Db.dialect.Placeholder()) != len(stmt.values.list) {
		return "", ErrorNumberOfConditionValues
	}

	for _, value := range stmt.values.list {
		query = strings.Replace(query, stmt.Db.dialect.Placeholder(), stmt.Db.dialect.Encode(value), 1)
	}

	return query, nil
}

func (stmt *StmtExecute) Exec() (sql.Result, error) {
	query, err := stmt.Build()
	if err != nil {
		return nil, err
	}

	return stmt.Db.Exec(query)
}

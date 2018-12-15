package dbr

import (
	"database/sql"
	"strings"
)

type StmtExecute struct {
	query  string
	values values

	db *db
}

func newStmtExecute(db *db, query string) *StmtExecute {
	return &StmtExecute{db: db, query: query, values: values{db: db}}
}

func (stmt *StmtExecute) Values(valuesList ...interface{}) *StmtExecute {
	stmt.values.list = append(stmt.values.list, valuesList...)
	return stmt
}

func (stmt *StmtExecute) Build() (string, error) {
	query := stmt.query

	if strings.Count(query, stmt.db.dialect.Placeholder()) != len(stmt.values.list) {
		return "", ErrorNumberOfConditionValues
	}

	for _, value := range stmt.values.list {
		query = strings.Replace(query, stmt.db.dialect.Placeholder(), stmt.db.dialect.Encode(value), 1)
	}

	return query, nil
}

func (stmt *StmtExecute) Exec() (sql.Result, error) {
	query, err := stmt.Build()
	if err != nil {
		return nil, err
	}

	return stmt.db.Exec(query)
}

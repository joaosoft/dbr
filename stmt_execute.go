package dbr

import (
	"database/sql"
	"strings"
	"time"
)

type StmtExecute struct {
	query  string
	values values

	Dbr *Dbr
	db *db
}

func newStmtExecute(dbr *Dbr, db *db, query string) *StmtExecute {
	return &StmtExecute{Dbr: dbr, db: db, query: query, values: values{db: dbr.Connections.Write}}
}

func (stmt *StmtExecute) Values(valuesList ...interface{}) *StmtExecute {
	stmt.values.list = append(stmt.values.list, valuesList...)
	return stmt
}

func (stmt *StmtExecute) Build() (string, error) {
	query := stmt.query

	if strings.Count(query, stmt.db.Dialect.Placeholder()) != len(stmt.values.list) {
		return "", ErrorNumberOfConditionValues
	}

	for _, value := range stmt.values.list {
		query = strings.Replace(query, stmt.db.Dialect.Placeholder(), stmt.db.Dialect.Encode(value), 1)
	}

	return query, nil
}

func (stmt *StmtExecute) Exec() (sql.Result, error) {

	startTime := time.Now()
	defer func() {
		stmt.db.Duration = time.Since(startTime)
	}()

	query, err := stmt.Build()
	if err != nil {
		return nil, err
	}

	return stmt.db.Exec(query)
}

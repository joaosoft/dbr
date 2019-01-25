package dbr

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"
)

type StmtDelete struct {
	withStmt   *StmtWith
	table      string
	conditions conditions
	returning  columns

	Dbr      *Dbr
	Db       *db
	Duration time.Duration
}

func newStmtDelete(dbr *Dbr, db *db, withStmt *StmtWith) *StmtDelete {
	return &StmtDelete{Dbr: dbr, Db: db, withStmt: withStmt, conditions: conditions{db: dbr.Connections.Write}}
}

func (stmt *StmtDelete) From(table string) *StmtDelete {
	stmt.table = table
	return stmt
}

func (stmt *StmtDelete) Where(query string, values ...interface{}) *StmtDelete {
	stmt.conditions.list = append(stmt.conditions.list, &condition{query: query, values: values})
	return stmt
}

func (stmt *StmtDelete) Build() (string, error) {
	var query string

	// withStmt
	if len(stmt.withStmt.withs) > 0 {
		withStmt, err := stmt.withStmt.Build()
		if err != nil {
			return "", err
		}
		query += withStmt
	}

	query += fmt.Sprintf("DELETE FROM %s", stmt.table)

	if len(stmt.conditions.list) > 0 {
		conds, err := stmt.conditions.Build()
		if err != nil {
			return "", err
		}

		query += fmt.Sprintf(" WHERE %s", conds)
	}

	if len(stmt.returning) > 0 {
		returning, err := stmt.returning.Build()
		if err != nil {
			return "", err
		}

		query += fmt.Sprintf(" RETURNING %s", returning)
	}

	return query, nil
}

func (stmt *StmtDelete) Exec() (sql.Result, error) {

	startTime := time.Now()
	defer func() {
		stmt.Duration = time.Since(startTime)
	}()

	query, err := stmt.Build()
	if err != nil {
		return nil, err
	}

	return stmt.Db.Exec(query)
}

func (stmt *StmtDelete) Return(column ...interface{}) *StmtDelete {
	stmt.returning = append(stmt.returning, column...)
	return stmt
}

func (stmt *StmtDelete) Load(object interface{}) error {

	if !reflect.ValueOf(object).CanAddr() {
		panic("the object is not a pointer the load")
	}

	startTime := time.Now()
	defer func() {
		stmt.Duration = time.Since(startTime)
	}()

	value := reflect.ValueOf(object)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return ErrorInvalidPointer
	}

	query, err := stmt.Build()
	if err != nil {
		return err
	}

	rows, err := stmt.Db.Query(query)
	if err != nil {
		return err
	}

	defer rows.Close()

	_, err = read(stmt.returning, rows, value)

	return err
}

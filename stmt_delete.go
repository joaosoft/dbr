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

	Dbr          *Dbr
	Db           *db
	Duration     time.Duration
	sqlOperation SqlOperation
}

func newStmtDelete(dbr *Dbr, db *db, withStmt *StmtWith) *StmtDelete {
	return &StmtDelete{sqlOperation: DeleteOperation, Dbr: dbr, Db: db, withStmt: withStmt, conditions: conditions{db: dbr.Connections.Write}}
}

func (stmt *StmtDelete) From(table string) *StmtDelete {
	stmt.table = table
	return stmt
}

func (stmt *StmtDelete) Where(query string, values ...interface{}) *StmtDelete {
	stmt.conditions.list = append(stmt.conditions.list, &condition{operator: operatorAnd, query: query, values: values, db: stmt.Db})
	return stmt
}

func (stmt *StmtDelete) WhereOr(query string, values ...interface{}) *StmtDelete {
	stmt.conditions.list = append(stmt.conditions.list, &condition{operator: operatorOr, query: query, values: values, db: stmt.Db})
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

	result, err := stmt.Db.Exec(query)

	if err := stmt.Dbr.eventHandler(stmt.sqlOperation, []string{stmt.table}, query, err, nil, result); err != nil {
		return nil, err
	}

	return result, err
}

func (stmt *StmtDelete) Return(column ...interface{}) *StmtDelete {
	stmt.returning = append(stmt.returning, column...)
	return stmt
}

func (stmt *StmtDelete) Load(object interface{}) error {

	value := reflect.ValueOf(object)
	if value.Kind() != reflect.Ptr {
		panic("the object is not a pointer the load")
	}

	startTime := time.Now()
	defer func() {
		stmt.Duration = time.Since(startTime)
	}()

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

	if err := stmt.Dbr.eventHandler(stmt.sqlOperation, []string{stmt.table}, query, err, rows, nil); err != nil {
		return err
	}

	defer rows.Close()

	_, err = read(stmt.returning, rows, value)

	return err
}

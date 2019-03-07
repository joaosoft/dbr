package dbr

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"
)

type StmtUpdate struct {
	withStmt   *StmtWith
	table      string
	sets       sets
	columns    columns
	conditions conditions
	returning  columns

	Dbr      *Dbr
	Db       *db
	Duration time.Duration
}

func newStmtUpdate(dbr *Dbr, db *db, withStmt *StmtWith, table string) *StmtUpdate {
	return &StmtUpdate{Dbr: dbr, Db: db, withStmt: withStmt, table: table, sets: sets{db: dbr.Connections.Write}, conditions: conditions{db: dbr.Connections.Write}}
}

func (stmt *StmtUpdate) Set(column string, value interface{}) *StmtUpdate {
	stmt.sets.list = append(stmt.sets.list, &set{column: column, value: value})
	return stmt
}

func (stmt *StmtUpdate) From(table string) *StmtUpdate {
	stmt.table = table
	return stmt
}

func (stmt *StmtUpdate) Where(query string, values ...interface{}) *StmtUpdate {
	stmt.conditions.list = append(stmt.conditions.list, &condition{operator: operatorAnd, query: query, values: values, db: stmt.Db})
	return stmt
}

func (stmt *StmtUpdate) WhereOr(query string, values ...interface{}) *StmtUpdate {
	stmt.conditions.list = append(stmt.conditions.list, &condition{operator: operatorOr, query: query, values: values, db: stmt.Db})
	return stmt
}

func (stmt *StmtUpdate) Build() (string, error) {
	var query string

	// withStmt
	withStmt, err := stmt.withStmt.Build()
	if err != nil {
		return "", err
	}
	query += withStmt

	sets, err := stmt.sets.Build()
	if err != nil {
		return "", err
	}

	query += fmt.Sprintf("UPDATE %s SET %s", stmt.table, sets)

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

func (stmt *StmtUpdate) Exec() (sql.Result, error) {

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

func (stmt *StmtUpdate) Record(record interface{}) *StmtUpdate {
	value := reflect.ValueOf(record)

	mappedValues := make(map[interface{}]reflect.Value)

	if len(stmt.columns) == 0 {
		var columns []interface{}
		loadStructValues(loadOptionWrite, value, &columns, mappedValues)
		stmt.columns = columns
	} else {
		loadStructValues(loadOptionWrite, value, nil, mappedValues)
	}

	for _, column := range stmt.columns {
		stmt.sets.list = append(stmt.sets.list, &set{column: column, value: mappedValues[column].Interface()})
	}

	return stmt
}

func (stmt *StmtUpdate) Return(column ...interface{}) *StmtUpdate {
	stmt.returning = append(stmt.returning, column...)
	return stmt
}

func (stmt *StmtUpdate) Load(object interface{}) error {

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

	defer rows.Close()

	_, err = read(stmt.returning, rows, value)

	return err
}

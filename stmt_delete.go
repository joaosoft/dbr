package dbr

import (
	"database/sql"
	"fmt"
	"reflect"
)

type StmtDelete struct {
	withStmt   *StmtWith
	table      string
	conditions conditions
	returning  columns

	db *db
}

func newStmtDelete(db *db, withStmt *StmtWith) *StmtDelete {
	return &StmtDelete{db: db, withStmt: withStmt, conditions: conditions{db: db}}
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
	query, err := stmt.Build()
	if err != nil {
		return nil, err
	}

	return stmt.db.Exec(query)
}

func (stmt *StmtDelete) Return(column ...string) *StmtDelete {
	stmt.returning = append(stmt.returning, column...)
	return stmt
}

func (stmt *StmtDelete) Load(object interface{}) error {
	value := reflect.ValueOf(object)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return ErrorInvalidPointer
	}

	query, err := stmt.Build()
	if err != nil {
		return err
	}

	rows, err := stmt.db.Query(query)
	if err != nil {
		return err
	}

	defer rows.Close()

	_, err = read(stmt.returning, rows, value)

	return err
}

package dbr

import (
	"database/sql"
	"fmt"
	"reflect"
)

type StmtInsert struct {
	table     string
	columns   columns
	values    values
	returning columns

	db *Db
}

func newStmtInsert(db *Db) *StmtInsert {
	return &StmtInsert{db: db, values: values{db: db}}
}

func (stmt *StmtInsert) Into(table string) *StmtInsert {
	stmt.table = table
	return stmt
}

func (stmt *StmtInsert) Columns(columns ...string) *StmtInsert {
	stmt.columns = append(stmt.columns, columns...)
	return stmt
}

func (stmt *StmtInsert) Build() (string, error) {

	columns, err := stmt.columns.Build()
	if err != nil {
		return "", err
	}

	values, err := stmt.values.Build()
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", stmt.table, columns, values)

	if len(stmt.returning) > 0 {
		returning, err := stmt.returning.Build()
		if err != nil {
			return "", err
		}

		query += fmt.Sprintf(" RETURNING %s", returning)
	}

	return query, nil
}

func (stmt *StmtInsert) Exec() (sql.Result, error) {
	query, err := stmt.Build()
	if err != nil {
		return nil, err
	}

	return stmt.db.Exec(query)
}

func (stmt *StmtInsert) Record(object interface{}) *StmtInsert {
	value := reflect.ValueOf(object)

	mappedValues := make(map[string]reflect.Value)
	loadStructValues(value, mappedValues)

	for column, value := range mappedValues {
		stmt.columns = append(stmt.columns, column)
		stmt.values.list = append(stmt.values.list, value.Interface())
	}

	return stmt
}

func (stmt *StmtInsert) Return(column ...string) *StmtInsert {
	stmt.returning = append(stmt.returning, column...)
	return stmt
}

package dbr

import (
	"database/sql"
	"fmt"
	"reflect"
)

type StmtUpdate struct {
	withs      withs
	table      string
	sets       sets
	columns    columns
	conditions conditions
	returning  columns

	db *db
}

func newStmtUpdate(db *db, withs withs, table string) *StmtUpdate {
	return &StmtUpdate{db: db, withs: withs, table: table, sets: sets{db: db}, conditions: conditions{db: db}}
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
	stmt.conditions.list = append(stmt.conditions.list, &condition{query: query, values: values})
	return stmt
}

func (stmt *StmtUpdate) Build() (string, error) {
	var query string

	// withs
	if len(stmt.withs) > 0 {
		withs, err := stmt.withs.Build()
		if err != nil {
			return "", err
		}
		query += fmt.Sprintf("WITH %s ", withs)
	}

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
	query, err := stmt.Build()
	if err != nil {
		return nil, err
	}

	return stmt.db.Exec(query)
}

func (stmt *StmtUpdate) Record(record interface{}) *StmtUpdate {
	value := reflect.ValueOf(record)

	mappedValues := make(map[string]reflect.Value)

	if len(stmt.columns) == 0 {
		var columns []string
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

func (stmt *StmtUpdate) Return(column ...string) *StmtUpdate {
	stmt.returning = append(stmt.returning, column...)
	return stmt
}

func (stmt *StmtUpdate) Load(object interface{}) error {
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

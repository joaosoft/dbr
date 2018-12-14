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

func (stmt *StmtInsert) Values(valuesList ...interface{}) *StmtInsert {
	stmt.values.list = append(stmt.values.list, &values{db: stmt.db, list: valuesList})
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

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", stmt.table, columns, values)

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

func (stmt *StmtInsert) Record(record interface{}) *StmtInsert {
	value := reflect.ValueOf(record)

	mappedValues := make(map[string]reflect.Value)

	if len(stmt.columns) == 0 {
		var columns []string
		loadStructValues(loadOptionWrite, value, &columns, mappedValues)
		stmt.columns = columns
	} else {
		loadStructValues(loadOptionWrite, value, nil, mappedValues)
	}

	var valueList []interface{}
	for _, column := range stmt.columns {
		valueList = append(valueList, mappedValues[column].Interface())
	}

	stmt.values.list = append(stmt.values.list, &values{db: stmt.db, list: valueList})

	return stmt
}

func (stmt *StmtInsert) Records(records []interface{}) *StmtInsert {
	for _, record := range records {
		stmt.Record(record)
	}

	return stmt
}

func (stmt *StmtInsert) Return(column ...string) *StmtInsert {
	stmt.returning = append(stmt.returning, column...)
	return stmt
}

func (stmt *StmtInsert) Load(object interface{}) error {
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

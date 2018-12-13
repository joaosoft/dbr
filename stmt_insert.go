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

func (stmt *StmtInsert) Values(values ...interface{}) *StmtInsert {
	stmt.values.list = append(stmt.values.list, values...)
	return stmt
}

func (stmt *StmtInsert) Line(lineValues ...interface{}) *StmtInsert {
	stmt.values.list = append(stmt.values.list, &values{db: stmt.db, list: lineValues})
	return stmt
}

func (stmt *StmtInsert) LineRecord(lineRecords ...interface{}) *StmtInsert {
	for _, lineRecord := range lineRecords {
		value := reflect.ValueOf(lineRecord)

		mappedValues := make(map[string]reflect.Value)
		loadStructValues(loadOptionWrite, value, mappedValues)

		var valueList []interface{}
		for _, column := range stmt.columns {
			if columnValue, ok := mappedValues[column]; ok {
				valueList = append(valueList, columnValue.Interface())
			}
		}

		stmt.values.list = append(stmt.values.list, &values{db: stmt.db, list: valueList})
	}

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

func (stmt *StmtInsert) Record(object interface{}) *StmtInsert {
	value := reflect.ValueOf(object)

	mappedValues := make(map[string]reflect.Value)
	loadStructValues(loadOptionWrite, value, mappedValues)

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

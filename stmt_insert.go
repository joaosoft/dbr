package dbr

import (
	"database/sql"
	"fmt"
	"reflect"
)

type StmtInsert struct {
	withStmt     *StmtWith
	table        interface{}
	columns      columns
	values       values
	returning    columns
	stmtConflict StmtConflict
	fromSelect   *StmtSelect

	Dbr *Dbr
	db  *db
}

func newStmtInsert(dbr *Dbr, db *db, withStmt *StmtWith) *StmtInsert {
	return &StmtInsert{Dbr: dbr, db: db, withStmt: withStmt, values: values{db: dbr.connections.write}, stmtConflict: StmtConflict{onConflictDoUpdate: sets{db: dbr.connections.write}}}
}

func (stmt *StmtInsert) Into(table interface{}) *StmtInsert {
	stmt.table = table
	return stmt
}

func (stmt *StmtInsert) Columns(columns ...interface{}) *StmtInsert {
	stmt.columns = append(stmt.columns, columns...)
	return stmt
}

func (stmt *StmtInsert) Values(valuesList ...interface{}) *StmtInsert {
	stmt.values.list = append(stmt.values.list, &values{db: stmt.db, list: valuesList})
	return stmt
}

func (stmt *StmtInsert) FromSelect(selectStmt *StmtSelect) *StmtInsert {
	stmt.fromSelect = selectStmt
	return stmt
}

func (stmt *StmtInsert) Build() (string, error) {
	var query string

	// withStmt
	if len(stmt.withStmt.withs) > 0 {
		withStmt, err := stmt.withStmt.Build()
		if err != nil {
			return "", err
		}
		query += withStmt
	}

	// columns
	columns, err := stmt.columns.Build()
	if err != nil {
		return "", err
	}

	// values
	var values string
	if len(stmt.values.list) > 0 {
		values = "VALUES "
		val, err := stmt.values.Build()
		if err != nil {
			return "", err
		}
		values += val
	}

	// from select statement
	var selectStmt string
	if stmt.fromSelect != nil {
		selectStmt, err = stmt.fromSelect.Build()
		if err != nil {
			return "", err
		}
	}

	query += fmt.Sprintf("INSERT INTO %s (%s) %s%s", stmt.table, columns, values, selectStmt)

	// on conflict
	if stmt.stmtConflict.onConflictType != "" {
		onConflictStmt, err := stmt.stmtConflict.Build()
		if err != nil {
			return "", err
		}

		query += onConflictStmt
	}

	// returning
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

	mappedValues := make(map[interface{}]reflect.Value)

	if len(stmt.columns) == 0 {
		var columns []interface{}
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

func (stmt *StmtInsert) OnConflict(column ...interface{}) *StmtInsert {
	stmt.stmtConflict.onConflictType = onConflictColumn
	stmt.stmtConflict.onConflict = append(stmt.stmtConflict.onConflict, column...)
	return stmt
}

func (stmt *StmtInsert) OnConflictConstraint(constraint interface{}) *StmtInsert {
	stmt.stmtConflict.onConflictType = onConflictConstraint
	stmt.stmtConflict.onConflict = []interface{}{constraint}
	return stmt
}

func (stmt *StmtInsert) DoNothing() *StmtInsert {
	stmt.stmtConflict.onConflictDoType = onConflictDoNothing
	return stmt
}

func (stmt *StmtInsert) DoUpdate(fieldValue ...interface{}) *StmtInsert {
	stmt.stmtConflict.onConflictDoType = onConflictDoUpdate

	if len(fieldValue)%2 != 0 {
		return stmt
	}

	lenC := len(fieldValue)
	for i := 0; i < lenC; i += 2 {
		stmt.stmtConflict.onConflictDoUpdate.list = append(stmt.stmtConflict.onConflictDoUpdate.list, &set{column: fieldValue[i].(string), value: fieldValue[i+1]})
	}

	return stmt
}

func (stmt *StmtInsert) Return(column ...interface{}) *StmtInsert {
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

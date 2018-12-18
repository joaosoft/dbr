package dbr

import (
	"database/sql"
	"fmt"
	"reflect"
)

type StmtInsert struct {
	withs        withs
	table        string
	columns      columns
	values       values
	returning    columns
	stmtConflict StmtConflict
	fromSelect   *StmtSelect

	db *db
}

func newStmtInsert(db *db, withs withs) *StmtInsert {
	return &StmtInsert{db: db, withs: withs, values: values{db: db}, stmtConflict: StmtConflict{onConflictDoUpdate: sets{db: db}}}
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

func (stmt *StmtInsert) FromSelect(selectStmt *StmtSelect) *StmtInsert {
	stmt.fromSelect = selectStmt
	return stmt
}

func (stmt *StmtInsert) Build() (string, error) {
	var query string

	// withs
	if len(stmt.withs) > 0 {
		withs, err := stmt.withs.Build()
		if err != nil {
			return "", err
		}
		query += fmt.Sprintf("WITH %s ", withs)
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

func (stmt *StmtInsert) OnConflict(column ...string) *StmtInsert {
	stmt.stmtConflict.onConflictType = onConflictColumn
	stmt.stmtConflict.onConflict = append(stmt.stmtConflict.onConflict, column...)
	return stmt
}

func (stmt *StmtInsert) OnConflictConstraint(constraint string) *StmtInsert {
	stmt.stmtConflict.onConflictType = onConflictConstraint
	stmt.stmtConflict.onConflict = []string{constraint}
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
